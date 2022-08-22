package client

import (
	"github.com/GeoinovaDev/lower-resultys/convert/decode"
	"github.com/GeoinovaDev/lower-resultys/net/request"
)

// Client ...
type Client struct {
	AccessToken string
}

// ProtocolResponse ...
type ProtocolResponse struct {
	Status  int    `json:"status"`
	Sucesso bool   `json:"sucesso"`
	Reason  int    `json:"motivo"`
	Message string `json:"mensagem"`
	Data    struct {
		ID int `json:"id"`
	} `json:"dados"`
}

// CompostoAction ...
type CompostoAction struct {
	Name      string            `json:"acao"`
	Option    string            `json:"opcao"`
	Parameter map[string]string `json:"acao_dados"`
}

type CompostoParameter struct {
	PhoneNumber   string           `json:"numero_destino"`
	IsSaveAudio   bool             `json:"gravar_audio"`
	IsDetectURA   bool             `json:"detecta_caixa"`
	IsLocalNumber bool             `json:"bina_inteligente"`
	Actions       []CompostoAction `json:"dados"`
}

// New ...
func New(accessToken string) *Client {
	return &Client{AccessToken: accessToken}
}

// Composto ...
func (c *Client) Composto(param CompostoParameter) (*ProtocolResponse, error) {
	rq := request.New("https://api2.totalvoice.com.br/composto")
	rq.AddHeader("Accept", "application/json")
	rq.AddHeader("Content-Type", "application/json")
	rq.AddHeader("Access-Token", c.AccessToken)

	response, err := rq.PostJSON(param)
	if err != nil {
		return nil, err
	}

	protocol := &ProtocolResponse{}
	decode.JSON(response, &protocol)

	return protocol, nil
}
