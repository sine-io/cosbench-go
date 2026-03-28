package s3

import (
	"crypto/tls"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/smithy-go/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	featuremanager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	awstypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/sine-io/cosbench-go/internal/application/ports"
)

type Adapter struct {
	backend string
	config  Config
	client  *awss3.Client
}

var _ ports.StorageAdapter = (*Adapter)(nil)

func NewAdapter(backend, raw string) *Adapter {
	_ = raw
	return &Adapter{backend: backend}
}

func (a *Adapter) Init(cfg map[string]string) error {
	parsed, err := ParseConfigMap(a.backend, cfg)
	if err != nil {
		return err
	}
	httpClient := buildHTTPClient(parsed)
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(parsed.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(parsed.AccessKey, parsed.SecretKey, "")),
		awsconfig.WithHTTPClient(httpClient),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}
	a.client = awss3.NewFromConfig(awsCfg, func(o *awss3.Options) {
		o.UsePathStyle = parsed.PathStyle
		o.BaseEndpoint = aws.String(parsed.Endpoint)
		if parsed.StorageClass != "" {
			o.APIOptions = append(o.APIOptions, setStorageClassMiddleware(parsed.StorageClass))
		}
	})
	a.config = parsed
	return nil
}

func (a *Adapter) Dispose() error {
	a.client = nil
	return nil
}

func (a *Adapter) CreateBucket(ctx context.Context, bucket string) error {
	_, err := a.client.CreateBucket(ctx, &awss3.CreateBucketInput{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket))})
	if err != nil && isBucketAlreadyOwned(err) {
		return nil
	}
	return err
}

func (a *Adapter) DeleteBucket(ctx context.Context, bucket string) error {
	_, err := a.client.DeleteBucket(ctx, &awss3.DeleteBucketInput{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket))})
	if err != nil && isDeleteMissingError(err) {
		return nil
	}
	return err
}

func (a *Adapter) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64) error {
	_, err := a.client.PutObject(ctx, &awss3.PutObjectInput{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket)), Key: aws.String(key), Body: body, ContentLength: aws.Int64(size)})
	return err
}

func (a *Adapter) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	return a.GetObjectWithOptions(ctx, bucket, key, ports.ReadOptions{})
}

func (a *Adapter) GetObjectWithOptions(ctx context.Context, bucket, key string, opts ports.ReadOptions) (io.ReadCloser, error) {
	input := &awss3.GetObjectInput{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket)), Key: aws.String(key)}
	if opts.HasRange {
		input.Range = aws.String(fmt.Sprintf("bytes=%d-%d", opts.RangeStart, opts.RangeEnd))
	}
	var apiOptions []func(*middleware.Stack) error
	if opts.Prefetch {
		apiOptions = append(apiOptions, setHeaderMiddleware("prefetch", "value"))
	}
	resp, err := a.client.GetObject(ctx, input, func(o *awss3.Options) {
		if len(apiOptions) > 0 {
			o.APIOptions = append(o.APIOptions, apiOptions...)
		}
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (a *Adapter) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := a.client.DeleteObject(ctx, &awss3.DeleteObjectInput{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket)), Key: aws.String(key)})
	if err != nil && isDeleteMissingError(err) {
		return nil
	}
	return err
}

func (a *Adapter) HeadObject(ctx context.Context, bucket, key string) (ports.ObjectMeta, error) {
	resp, err := a.client.HeadObject(ctx, &awss3.HeadObjectInput{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket)), Key: aws.String(key)})
	if err != nil {
		return ports.ObjectMeta{}, err
	}
	return ports.ObjectMeta{ContentLength: aws.ToInt64(resp.ContentLength), ContentType: aws.ToString(resp.ContentType), ETag: aws.ToString(resp.ETag)}, nil
}

func (a *Adapter) ListObjects(ctx context.Context, bucket, prefix string, maxKeys int) ([]ports.ObjectEntry, error) {
	input := &awss3.ListObjectsV2Input{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket)), Prefix: aws.String(prefix)}
	if maxKeys > 0 {
		input.MaxKeys = aws.Int32(int32(maxKeys))
	}
	resp, err := a.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}
	items := make([]ports.ObjectEntry, 0, len(resp.Contents))
	for _, item := range resp.Contents {
		items = append(items, ports.ObjectEntry{Key: aws.ToString(item.Key), Size: aws.ToInt64(item.Size)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Key < items[j].Key })
	return items, nil
}

func (a *Adapter) MultipartPut(ctx context.Context, bucket, key string, body io.Reader, size int64, partSize int64) error {
	uploader := featuremanager.NewUploader(a.client, func(u *featuremanager.Uploader) {
		if partSize > 0 {
			u.PartSize = partSize
		}
	})
	_, err := uploader.Upload(ctx, &awss3.PutObjectInput{Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket)), Key: aws.String(key), Body: body})
	return err
}

func (a *Adapter) RestoreObject(ctx context.Context, bucket, key string, days int) error {
	_, err := a.client.RestoreObject(ctx, &awss3.RestoreObjectInput{
		Bucket: aws.String(normalizeBucketName(a.config.Backend, bucket)),
		Key:    aws.String(key),
		RestoreRequest: &awstypes.RestoreRequest{
			Days: aws.Int32(int32(days)),
		},
	})
	return err
}

func isBucketAlreadyOwned(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "bucketalreadyownedbyyou") || strings.Contains(msg, "bucket already exists")
}

func isDeleteMissingError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "nosuchbucket") || strings.Contains(msg, "nosuchkey") || strings.Contains(msg, "404") || strings.Contains(msg, "not found")
}

func normalizeBucketName(backend, bucket string) string {
	kind := strings.ToLower(strings.TrimSpace(backend))
	if kind != "sio" && kind != "siov1" && kind != "gdas" {
		return bucket
	}
	head, _, found := strings.Cut(bucket, "/")
	if !found {
		return bucket
	}
	return head
}

func buildHTTPClient(cfg Config) *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.NoVerifySSL},
		DialContext: (&net.Dialer{}).DialContext,
	}
	if cfg.ProxyHost != "" {
		proxyURL := "http://" + cfg.ProxyHost
		if cfg.ProxyPort != "" {
			proxyURL += ":" + cfg.ProxyPort
		}
		transport.Proxy = func(*http.Request) (*url.URL, error) {
			return url.Parse(proxyURL)
		}
	}
	return &http.Client{Transport: transport}
}

func setStorageClassMiddleware(storageClass string) func(*middleware.Stack) error {
	class := awstypes.StorageClass(storageClass)
	return func(stack *middleware.Stack) error {
		return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("set-storage-class", func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (middleware.InitializeOutput, middleware.Metadata, error) {
			switch input := in.Parameters.(type) {
			case *awss3.PutObjectInput:
				if input.StorageClass == "" {
					input.StorageClass = class
				}
			case *awss3.CreateMultipartUploadInput:
				if input.StorageClass == "" {
					input.StorageClass = class
				}
			}
			return next.HandleInitialize(ctx, in)
		}), middleware.After)
	}
}

func setHeaderMiddleware(name, value string) func(*middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		return stack.Build.Add(middleware.BuildMiddlewareFunc("set-header-"+name, func(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (middleware.BuildOutput, middleware.Metadata, error) {
			req, ok := in.Request.(*awshttp.Request)
			if ok {
				req.Header.Set(name, value)
			}
			return next.HandleBuild(ctx, in)
		}), middleware.After)
	}
}
