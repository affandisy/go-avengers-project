package debug

import "log"

func LogDebug(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func ErrorDebug(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}
