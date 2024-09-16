package omidpay

type Callback struct {
	MID               string  `json:"MID"`
	State             string  `json:"State"`
	ResNum            string  `json:"ResNum"`
	RefNum            *string `json:"RefNum"`
	Language          string  `json:"language"`
	Token             string  `json:"token"`
	TraceNo           *string `json:"traceNo"`
	CustomerRefNum    *string `json:"CustomerRefNum"`
	RedirectUrl       *string `json:"redirectUrl"`
	CardHashPan       *string `json:"CardHashPan"`
	CardMaskPan       *string `json:"CardMaskPan"`
	TransactionAmount *int64  `json:"transactionAmount"`
	UserID            *string `json:"userId"`
	EmailAddress      *string `json:"emailAddress"`
	MobileNo          *string `json:"mobileNo"`
}
