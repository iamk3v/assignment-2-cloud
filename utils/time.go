package utils

import (
	"assignment-2/config"
	"fmt"
	"time"
)

/*
StartTime initiates the Starttime global variable
*/
func StartTime() {
	config.Starttime = time.Now()
}

/*
GetTime Gets the time and formats it to human-readable
*/
func GetTime() string {
	uptime := time.Since(config.Starttime)
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	return fmt.Sprintf("%dd:%02dh:%02dm:%02ds", days, hours, minutes, seconds)
}
