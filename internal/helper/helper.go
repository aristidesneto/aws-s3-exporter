package helper

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
)

func Contains(list []string, bucket string) bool {
	return slices.Contains(list, bucket)
}

// Extrai o prefixo da pasta
func ExtractDatePrefix(key string) string {
	prefixDateRegex := regexp.MustCompile(`^\d{4}/\d{2}/\d{2}/.+`)

	if !prefixDateRegex.MatchString(key) {
		return ""
	}

	parts := strings.Split(key, "/")
	if len(parts) >= 3 {
		return fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
	}
	return "unknown"
}

// Extrai o prefixo e verifica se a data é recente
func ExtractDatePrefixAndCheck(key string, lastModified time.Time, daysLimit int) (string, bool) {
	prefix := ExtractDatePrefix(key)
	if prefix == "" {
		return "", false
	}

	limite := time.Now().AddDate(0, 0, -daysLimit)
	if lastModified.After(limite) {
		return prefix, true
	}
	return "", false
}
