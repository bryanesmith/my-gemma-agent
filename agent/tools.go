package agent

import (
	"strings"
	"time"
)

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func IsDateQuery(input string) bool {
	lower := strings.ToLower(input)
	return strings.Contains(lower, "date") || strings.Contains(lower, "time")
}
