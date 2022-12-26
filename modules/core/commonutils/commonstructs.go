package commonutils

import (
	"time"
)

type Data struct {
	Token           string
	AdminServer     string
	AdminChannel    string
	InfoChannel     string
	WarnChannel     string
	ErrChannel      string
	UpdateChannel   string
	FeedbackChannel string
	Debugging       bool
}

type LogEntry struct {
	Time    time.Time
	Message string
	Level   string
}

var AuthKeyUpdater []func(GID string, AuthKey string, OldKey string)
