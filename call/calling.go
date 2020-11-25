package call

import (
	"time"

	"git.resultys.com.br/sdk/totalvoice-golang/payload"
	"git.resultys.com.br/sdk/totalvoice-golang/webhook"
)

// Call ...
type Call struct {
	ID      int
	Status  string
	Message string

	webhook *webhook.Server
}

// New ...
func New(webhook *webhook.Server) *Call {
	return &Call{
		webhook: webhook,
	}
}

// Wait ...
func (c *Call) Wait() (response payload.CallResponse) {
	if c.Status == "calling" {
		count := 0
		isWebhookResponsed := false

		c.webhook.AddHook(c.ID).Ok(func(proto interface{}) {
			response = proto.(payload.CallResponse)
			isWebhookResponsed = true
			c.Status = "success"
		})

		for {
			time.Sleep(1 * time.Second)
			if isWebhookResponsed {
				break
			}
			count++

			// espera até 2 min
			if count == 120 {
				c.Status = "timeout"
				break
			}
		}
	}

	return
}
