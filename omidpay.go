package omidpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	terminalID       string
	terminalPassword string
	redirectUrl      string

	httpClient *http.Client

	urlPrefix string
}

func NewClient(terminalID, terminalPassword string) *Client {
	return &Client{
		terminalID:       terminalID,
		terminalPassword: terminalPassword,

		httpClient: &http.Client{},

		urlPrefix: defaultPrefix,
	}
}

func NewClientWithReverseProxy(proxyURL string, terminalID, terminalPassword string) *Client {
	return &Client{
		terminalID:       terminalID,
		terminalPassword: terminalPassword,

		httpClient: &http.Client{},

		urlPrefix: proxyURL,
	}
}

func NewClientWithHttpProxy(terminalID, terminalPassword string, proxyURL string) (*Client, error) {
	proxyUrl, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	return &Client{
		terminalID:       terminalID,
		terminalPassword: terminalPassword,

		httpClient: &http.Client{Transport: transport},

		urlPrefix: defaultPrefix,
	}, nil
}

func (c *Client) GeneratePaymentToken(invoiceID string, amount float64, redirectURL string) (*TokenResponse, error) {
	data := map[string]interface{}{
		"WSContext": map[string]string{
			"UserId":   c.terminalID,
			"Password": c.terminalPassword,
		},
		"TransType":   "EN_GOODS",
		"ReserveNum":  invoiceID,
		"MerchantId":  c.terminalID,
		"Amount":      amount,
		"RedirectUrl": redirectURL,
	}

	j, _ := json.Marshal(data)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", c.urlPrefix, tokenEndpoint),
		bytes.NewReader(j),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(b))
	}

	b, _ := io.ReadAll(resp.Body)

	var tokenResponse *TokenResponse

	err = json.Unmarshal(b, &tokenResponse)
	if err != nil {
		return nil, err
	}

	if tokenResponse.Result != "erSucceed" {
		errMsg := statusCodes[tokenResponse.Result]
		if errMsg == "" {
			errMsg = tokenResponse.Result
		}

		return nil, errors.New(errMsg)
	}

	return tokenResponse, nil
}

func (c *Client) VerifyPayment(token string, refNum string) (*VerifyResponse, error) {
	data := map[string]interface{}{
		"WSContext": map[string]string{
			"UserId":   c.terminalID,
			"Password": c.terminalPassword,
		},
		"Token":  token,
		"RefNum": refNum,
	}

	j, _ := json.Marshal(data)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", c.urlPrefix, verificationEndpoint),
		bytes.NewReader(j),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(b))
	}

	b, _ := io.ReadAll(resp.Body)

	var verifyResponse *VerifyResponse

	err = json.Unmarshal(b, &verifyResponse)
	if err != nil {
		return nil, err
	}

	if verifyResponse.Result != "erSucceed" {
		errMsg := statusCodes[verifyResponse.Result]
		if errMsg == "" {
			errMsg = verifyResponse.Result
		}

		return nil, errors.New(errMsg)
	}

	return verifyResponse, nil
}

func (c *Client) HttpCallback(
	callback func(c *Callback, r *http.Request, w http.ResponseWriter),
	onError func(err error, r *http.Request, w http.ResponseWriter),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		fmt.Println(r.Form)

		mid := r.Form.Get("MID")
		state := r.Form.Get("State")
		ResNum := r.Form.Get("ResNum")
		language := r.Form.Get("language")
		token := r.Form.Get("token")

		if mid == "" || state == "" || ResNum == "" || token == "" {
			fmt.Println(mid, state, ResNum, language, token)

			onError(MissingParams, r, w)

			return
		}

		cback := &Callback{
			MID:      mid,
			State:    state,
			ResNum:   ResNum,
			Language: language,
			Token:    token,
		}

		if state == "OK" {
			if traceNo := r.Form.Get("TraceNo"); traceNo != "" {
				cback.TraceNo = &traceNo
			}

			if customerRefNum := r.Form.Get("CustomerRefNum"); customerRefNum != "" {
				cback.CustomerRefNum = &customerRefNum
			}

			if redirectUrl := r.Form.Get("redirectUrl"); redirectUrl != "" {
				cback.RedirectUrl = &redirectUrl
			}

			if cardHashPan := r.Form.Get("CardHashPan"); cardHashPan != "" {
				cback.CardHashPan = &cardHashPan
			}

			if cardMaskPan := r.Form.Get("CardMaskPan"); cardMaskPan != "" {
				cback.CardMaskPan = &cardMaskPan
			}

			if transactionAmount := r.Form.Get("transactionAmount"); transactionAmount != "" {
				amount, err := strconv.ParseInt(transactionAmount, 10, 64)
				if err != nil {
					onError(err, r, w)

					return
				}

				cback.TransactionAmount = &amount
			}

			if userId := r.Form.Get("userId"); userId != "" {
				cback.UserID = &userId
			}

			if refNum := r.Form.Get("RefNum"); refNum != "" {
				cback.RefNum = &refNum
			} else {
				onError(MissingRefNum, r, w)

				return
			}

			if emailAddress := r.Form.Get("emailAddress"); emailAddress != "" {
				cback.EmailAddress = &emailAddress
			}

			if mobileNo := r.Form.Get("mobileNo"); mobileNo != "" {
				cback.MobileNo = &mobileNo
			}

			callback(cback, r, w)
		}
	}
}
