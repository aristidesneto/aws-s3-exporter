package helper

import (
	"fmt"
	"slices"
	"strings"
)

func Contains(list []string, bucket string) bool {
	return slices.Contains(list, bucket)
}

func ExtractDatePrefix(key string) string {
	parts := strings.Split(key, "/")
	if len(parts) >= 3 {
		return fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
	}
	return "unknown"
}
