package parser

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/defectus/glutton/pkg/iface"
	"github.com/pkg/errors"
)

// SimpleParser is the default implementation if the parser interface.
type SimpleParser struct {
}

// Parse reads request and builds a payload from it.
func (s *SimpleParser) Parse(req *http.Request) (*iface.PayloadRecord, error) {
	payload := &iface.PayloadRecord{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading payload")
	}
	payload.Payload = string(body)
	payload.Timestamp = time.Now()
	payload.Remote = req.RemoteAddr
	payload.Meta = req.Header
	return payload, nil
}

// Configure initilizes the instance of parser.
func (s *SimpleParser) Configure(*iface.Settings) error {
	return nil
}
