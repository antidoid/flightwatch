package notify

import (
	"fmt"
	"os"

	"github.com/antidoid/flightwatch/initializers"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func init() {
	initializers.LoadEnvVars()
}

func SendSMS(to string, message string) error {
	fmt.Println("I was triggered to send a sms")
	params := &api.CreateMessageParams{}
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetTo(to)
	params.SetBody(message)

	fmt.Println(os.Getenv("TWILIO_PHONE_NUMBER"))
	fmt.Println(os.Getenv("TWILIO_ACC_SID"))
	fmt.Println(os.Getenv("TWILIO_AUTH_TOKEN"))

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACC_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	_, err := client.Api.CreateMessage(params)
	return err
}
