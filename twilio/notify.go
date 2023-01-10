package twilio

import (
	"os"

	"github.com/antidoid/flightwatch/initializers"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

var Client *twilio.RestClient

func init() {
    initializers.LoadEnvVars()
    Client = twilio.NewRestClient()
}

func SendSMS(to string, message string) error {
    params := &api.CreateMessageParams{}
	params.SetBody(message)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetTo(to)

    _, err := Client.Api.CreateMessage(params)
    return err
}

func SendWhatsAppMessage(to string, message string) error {
    params := &api.CreateMessageParams{}
	params.SetBody(message)
    params.SetFrom(os.Getenv("TWILIO_WHATSAPP_NUMBER"))
    params.SetTo("whatsapp:" + to)

    _, err := Client.Api.CreateMessage(params)
    return err
}
