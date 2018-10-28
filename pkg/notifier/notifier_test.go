package notifier

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/defectus/glutton/pkg/iface"

	"github.com/stretchr/testify/assert"
)

func TestSMTPNotifier_Notify1(t *testing.T) {
	if len(os.Getenv("SMTP_SERVER")) == 0 {
		t.Skip("make sure required env. variables are set before running this script, skipping")
	}
	settings := &iface.Configuration{}
	smtpNotifier := new(SMTPNotifier)
	err := smtpNotifier.Configure(&settings.Settings[0])
	assert.NoError(t, err)
	log.Printf("smtpNotifier:%+v", smtpNotifier)
	err = smtpNotifier.Notify(&iface.PayloadRecord{
		Payload:   "test payload",
		Timestamp: time.Now(),
		Remote:    "0.0.0.0",
	})
	assert.NoError(t, err)
}
