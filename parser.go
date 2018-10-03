package glutton

import (
	"time"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type SimpleParser struct {
}

func (s *SimpleParser) Parse(req *http.Request) (*PayloadRecord, error) {
	payload := &PayloadRecord{}
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

func (n *SimpleParser) Configure(*Settings) error {
	return nil
}
