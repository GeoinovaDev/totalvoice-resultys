package totalvoice

import (
	"github.com/GeoinovaDev/totalvoice-resultys/call"
	"github.com/GeoinovaDev/totalvoice-resultys/client"
	"github.com/GeoinovaDev/totalvoice-resultys/webhook"
)

// TotalVoice ...
type TotalVoice struct {
	client  *client.Client
	webhook *webhook.Server
}

// New ...
func New(accessToken string) *TotalVoice {
	return &TotalVoice{
		client:  client.New(accessToken),
		webhook: webhook.New(":36466"),
	}
}

// Call ...
func (t *TotalVoice) Call(param client.CompostoParameter) *call.Call {
	call := call.New(t.webhook)
	response, err := t.client.Composto(param)

	if err != nil {
		call.Status = "error"
		call.Message = err.Error()
	} else {
		call.ID = response.Data.ID
		call.Status = "calling"
		call.Message = ""
	}

	return call
}
