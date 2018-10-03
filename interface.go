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

type Configurable interface {
	Configure(*Settings) error
}

type PayloadParser interface {
	Configurable
	Parse(*http.Request) (*PayloadRecord, error)
}

type PayloadSaver interface {
	Configurable
	Save(*PayloadRecord) error
}

type PayloadNotifier interface {
	Configurable
	Notify(*PayloadRecord) error
}
