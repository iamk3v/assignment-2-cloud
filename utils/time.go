package utils

import (
	"assignment-2/config"
	"fmt"
	"time"
)

// logs the time on start
func StartTime() {
	config.Starttime = time.Now()
}

func Gettime() string {
	uptime := time.Since(config.Starttime)
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	return fmt.Sprintf("%dd:%02dh:%02dm:%02ds", days, hours, minutes, seconds)
}
