package omidpay

import (
	"html/template"
	"io"
)

type TokenResponse struct {
	Result         string `json:"Result"`
	ExpirationDate int64  `json:"ExpirationDate"`
	Token          string `json:"Token"`
	ChannelID      string `json:"ChannelId"`
	UserID         string `json:"UserId"`
	ResultCode     int    `json:"ResultCode"`
}

func (t *TokenResponse) HttpWriteSampleRedirectForm(rw io.Writer) error {
	data := struct {
		PaymentUrl   string
		PaymentToken string
	}{
		PaymentUrl:   paymentURL,
		PaymentToken: t.Token,
	}

	tmp, err := template.New("redirect").Parse(redirectionFormTemplate)
	if err != nil {
		return err
	}

	err = tmp.Execute(rw, data)
	if err != nil {
		return err
	}

	return nil
}
