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
	params.SetBody(message)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetTo(to)

    client := twilio.NewRestClient()
    _, err := client.Api.CreateMessage(params)
    return err
}

func SendWhatsAppMessage(to string, message string) error {
    params := &api.CreateMessageParams{}
	params.SetBody(message)
    params.SetFrom(os.Getenv("TWILIO_WHATSAPP_NUMBER"))
    params.SetTo("whatsapp:" + to)

    client := twilio.NewRestClient()
    _, err := client.Api.CreateMessage(params)
    return err
}

