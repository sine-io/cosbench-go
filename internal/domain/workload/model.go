package workload

type Workload struct {
	Name        string
	Description string
	Trigger     string
	Config      string
	Storage     *StorageSpec
	Workflow    Workflow
}

type Workflow struct {
	Config string
	Stages []Stage
}

type Stage struct {
	Name         string
	ClosureDelay int
	Trigger      string
	Config       string
	Storage      *StorageSpec
	Works        []Work
}

type Work struct {
	Name       string
	Type       string
	Workers    int
	Interval   int
	Division   string
	Runtime    int
	RampUp     int
	RampDown   int
	AFR        int
	TotalOps   int
	TotalBytes int64
	Driver     string
	Config     string
	Storage    *StorageSpec
	Operations []Operation
}

type Operation struct {
	Type     string
	Ratio    int
	Division string
	Config   string
	ID       string
}

type StorageSpec struct {
	Type   string
	Config string
}
