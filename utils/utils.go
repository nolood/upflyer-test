package utils

import (
	"regexp"
	"strings"
)

func Median(arr []int) int {
	result := 0
	for _, v := range arr {
		result += v
	}
	return result / 7
}

func IsValidChannelURL(url string) bool {
	regex := regexp.MustCompile(`^https?:\/\/t\.me\/([a-zA-Z0-9_-]+)$`)
	return regex.MatchString(url)
}

func ExtractChannelName(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
