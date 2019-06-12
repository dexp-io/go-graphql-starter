package dexp

import "time"

func Log(userId int64, messageType, message string) {

	DB.Insert("system_log").Fields(map[string]interface{}{
		"type":       messageType,
		"user_id":    userId,
		"read":       0,
		"message":    message,
		"created_at": time.Now(),
	})
}
