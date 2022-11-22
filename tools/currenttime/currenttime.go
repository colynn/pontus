package currenttime

import "time"

// GetCurrentTimeStr ..
func GetCurrentTimeStr() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// GetCurrentTime ..
func GetCurrentTime() time.Time {
	return time.Now()
}
