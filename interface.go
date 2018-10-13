package glutton

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PayloadRecord struct {
	Payload   string
	Timestamp time.Time
	Meta      map[string][]string
	Remote    string
}

func (p *PayloadRecord) String() string {
	builder := new(strings.Builder)
	builder.WriteString(p.Timestamp.String())
	builder.WriteString(":")
	builder.WriteString(p.Remote)
	builder.WriteString("\n\n")
	builder.WriteString(p.Payload)
	builder.WriteString("\n\n")
	builder.WriteString(fmt.Sprintf("%+v\n", p.Meta))
	return builder.String()
}

// Configurable is anything that can be configured with settings.
type Configurable interface {
	Configure(*Settings) error
}

// PayloadParser is anything that can parse http request and return some payload
type PayloadParser interface {
	Configurable
	Parse(*http.Request) (*PayloadRecord, error)
}

// PayloadSaver can save payload (e.g. to filesystem)
type PayloadSaver interface {
	Configurable
	Save(*PayloadRecord) error
}

// PayloadNotifier is anything that can notify (e.g. send an email, slack) of payload received. 
type PayloadNotifier interface {
	Configurable
	Notify(*PayloadRecord) error
}
