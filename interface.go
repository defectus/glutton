package glutton

import (
	"time"
)

type PayloadRecord struct {
	Payload   string
	Timestamp time.Time
	Meta      map[string][]string
	Remote    string
}

type PayloadParser interface {
	Parse([]byte) *PayloadRecord
}

type PayloadSaver interface {
	Save(*PayloadRecord) error
}

type PayloadNotifier interface {
	Notify(*PayloadRecord) error
}
