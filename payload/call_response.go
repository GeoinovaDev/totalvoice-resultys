package payload

import (
	"strconv"
)

// CallResponse ...
type CallResponse struct {
	ID                      int     `json:"id"`
	PhoneNumber             string  `json:"numero_destino"`
	CreateAt                string  `json:"data_criacao"`
	DateBegin               string  `json:"data_inicio"`
	Type                    string  `json:"tipo"`
	Status                  string  `json:"status"`
	DurationInSeconds       int     `json:"duracao_segundos"`
	Duration                string  `json:"duracao"`
	DurationTaxInSeconds    int     `json:"duracao_cobrada_segundos"`
	DurationTax             string  `json:"duracao_cobrada"`
	DurationSpeechInSeconds int     `json:"duracao_falada_segundos"`
	DurationSpeech          string  `json:"duracao_falada"`
	Price                   float32 `json:"preco"`
	IsUserResponse          bool    `json:"resposta_usuario"`
	UserResponse            string  `json:"resposta"`
	UrlRec                  string  `json:"url_gravacao"`
	Tags                    string  `json:"tags"`
}

func (c CallResponse) KeyPressNumber() int {
	if !c.IsUserResponse {
		return -1
	}

	n, e := strconv.Atoi(c.UserResponse)
	if e != nil {
		return -1
	}

	return n
}
