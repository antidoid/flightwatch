package notify

import (
	"os"

	"github.com/antidoid/flightwatch/initializers"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func init() {
	initializers.LoadEnvVars()
}

func SendSMS(to string, message string) error {
	params := &api.CreateMessageParams{}
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetTo(to)
	params.SetBody(message)

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACC_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	_, err := client.Api.CreateMessage(params)
	return err
}
