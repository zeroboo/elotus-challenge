package middleware

import (
	"net/http"
)

// GetLogContext retrieves the log context from request
func GetLogContext(r *http.Request) *LogContext {
	if logCtx, ok := r.Context().Value(logContextKey).(*LogContext); ok {
		return logCtx
	}
	return &LogContext{Fields: make(map[string]interface{})}
}

const (
	logContextKey = "logContext"
)

// LogContext holds fields that can be appended to during request processing
type LogContext struct {
	Fields map[string]interface{}
}

// Add appends a field to the log context
func (lc *LogContext) Add(key string, value interface{}) {
	lc.Fields[key] = value
}

// AddString appends a string field
func (lc *LogContext) AddString(key, value string) {
	lc.Fields[key] = value
}

// AddInt appends an int field
func (lc *LogContext) AddInt(key string, value int) {
	lc.Fields[key] = value
}

// AddLogEntries adds multiple key-value pairs to the log context
func AddLogEntries(r *http.Request, keyValues ...interface{}) {
	logCtx := GetLogContext(r)
	for i := 0; i < len(keyValues)-1; i += 2 {
		if key, ok := keyValues[i].(string); ok {
			logCtx.Add(key, keyValues[i+1])
		}
	}
}

// AddLogString adds a string field to the log context
func AddLogString(r *http.Request, key, value string) {
	GetLogContext(r).AddString(key, value)
}

// AddLogInt adds an int field to the log context
func AddLogInt(r *http.Request, key string, value int) {
	GetLogContext(r).AddInt(key, value)
}
