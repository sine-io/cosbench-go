package config

import "strings"

func ParseKVConfig(input string) map[string]string {
	result := map[string]string{}
	input = strings.TrimSpace(input)
	if input == "" {
		return result
	}
	for _, entry := range strings.Split(input, ";") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		pos := strings.Index(entry, "=")
		if pos < 0 {
			continue
		}
		key := strings.TrimSpace(entry[:pos])
		value := strings.TrimSpace(entry[pos+1:])
		result[key] = value
	}
	return result
}
