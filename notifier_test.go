package glutton

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSMTPNotifier_Notify1(t *testing.T) {
	if len(os.Getenv("SMTP_SERVER")) == 0 {
		t.Skip("make sure required env. variables are set before running this script, skipping")
	}
	settings := createSettings(new(Settings))
	smtpNotifier := new(SMTPNotifier)
	err := smtpNotifier.Configure(settings)
	assert.NoError(t, err)
	log.Printf("smtpNotifier:%+v", smtpNotifier)
	err = smtpNotifier.Notify(&PayloadRecord{
		Payload:   "test payload",
		Timestamp: time.Now(),
		Remote:    "0.0.0.0",
	})
	assert.NoError(t, err)
}
