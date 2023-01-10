package notify

import (
	"os"
	"strings"

	"github.com/antidoid/flightwatch/initializers"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func SendEmail(to string, message string) error {
    from := mail.NewEmail(os.Getenv("SENDGRID_NAME"), os.Getenv("SENDGRID_EMAIL"))
    subject := "Regarding your Tracked Flight"
    reciept := mail.NewEmail(strings.Split(to, "@")[0], to)
    // message may look like, hey <username>, your tracked flight <flight_no> just hit the threshold, its currently priced at ____,
    // Check it out: <linktobook>
    email := mail.NewSingleEmail(from, subject, reciept, message, message)

    client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
    _, err := client.Send(email)

    return err
}

