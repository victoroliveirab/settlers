package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type LogMessageStruct struct {
	Timestamp string `json:"timestamp"`
	Direction string `json:"direction"`
	ClientID  int64  `json:"client_id"`
	Type      string `json:"type"`
	Message   any    `json:"message"`
}

var logFile *os.File
var persist bool = false

func Init(shouldPersist bool) {
	if !shouldPersist {
		return
	}
	persist = true
	var err error
	logFile, err = os.OpenFile("logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// log.SetOutput(logFile)
}

func LogHttpRequest(method, path, remoteAddr, userAgent string, duration, statusCode int) {
	logEntry := map[string]interface{}{
		"method":      method,
		"path":        path,
		"remote_addr": remoteAddr,
		"user_agent":  userAgent,
		"status":      statusCode,
		"duration":    duration,
	}

	logData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}
	log.Println(string(logData))
	if persist {
		writeLog(logData)
	}
}

func LogWSMessage(direction string, clientID int64, msgType string, msg any) {
	logEntry := LogMessageStruct{
		Timestamp: time.Now().Format(time.RFC3339),
		Direction: direction,
		ClientID:  clientID,
		Type:      msgType,
		Message:   msg,
	}

	logData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	log.Println(string(logData))
	if persist {
		writeLog(logData)
	}
}

func Log(msg string) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   msg,
	}

	logData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	log.Println(string(logData))
	if persist {
		writeLog(logData)
	}

}

func LogError(clientID int64, action string, errorType int, err error) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"severity":  "error",
		"client_id": clientID,
		"action":    action,
		"type":      errorType,
		"error":     err.Error(),
	}

	logData, _ := json.Marshal(logEntry)
	log.Println(string(logData)) // Print error log in JSON format
	if persist {
		writeLog(logData)
	}
}

func LogSystemError(action string, errorType int, err error) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"severity":  "error",
		"client_id": "system",
		"action":    action,
		"type":      errorType,
		"error":     err.Error(),
	}

	logData, _ := json.Marshal(logEntry)
	log.Println(string(logData)) // Print error log in JSON format
	if persist {
		writeLog(logData)
	}
}

func LogMessage(clientID int64, action, message string) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"severity":  "info",
		"client_id": clientID,
		"action":    action,
		"message":   message,
	}

	logData, _ := json.Marshal(logEntry)
	log.Println(string(logData)) // Print error log in JSON format
	if persist {
		writeLog(logData)
	}
}

func LogSystemMessage(action, message string) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"severity":  "info",
		"client_id": "system",
		"action":    action,
		"message":   message,
	}

	logData, _ := json.Marshal(logEntry)
	log.Println(string(logData)) // Print error log in JSON format
	if persist {
		writeLog(logData)
	}
}

func writeLog(logData []byte) {
	// Ensure each log entry is on a new line
	_, err := logFile.Write(append(logData, '\n'))
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}
