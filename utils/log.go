package utils

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"push-guard/config"
	"runtime"
)

var Logger *StandardLogger = NewLogger()

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*slog.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	shouldDebug := os.Getenv("PUSH_GUARD_DEBUG")
	lvl := new(slog.LevelVar)
	if shouldDebug != "" {
		lvl.Set(slog.LevelDebug)
	} else {
		lvl.Set(slog.LevelError)
	}
	logger := &StandardLogger{slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))}
	return logger
}

type Message struct {
	User   string `json:"User"`
	Host   string `json:"Host"`
	OS     string `json:"OS"`
	Remote string `json:"Remote"`
	Pushed bool   `json:"Pushed"`
}

func NewMessage(remote string, pushed bool) *Message {
	hostname, _ := os.Hostname()
	return &Message{
		User:   os.Getenv("USER"),
		Host:   hostname,
		OS:     runtime.GOOS,
		Remote: remote,
		Pushed: pushed,
	}
}

func SendMessage(remote string, pushed bool) {
	jsonValue, _ := json.Marshal(NewMessage(remote, pushed))
	http.Post(config.LogCollectorURL, "application/json", bytes.NewBuffer(jsonValue))
}
