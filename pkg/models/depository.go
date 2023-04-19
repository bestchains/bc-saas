package models

// Depository defines valuable fields for a depository
type Depository struct {
	Index       string `json:"index" pg:"index"`
	KID         string `json:"kid" pg:"kid,pk"`
	Platform    string `json:"platform" pg:"platform"`
	Operator    string `json:"operator" pg:"operator"`
	Owner       string `json:"owner" pg:"owner"`
	BlockNumber string `json:"blockNumber" pg:"blockNumber"`

	// Content related
	ContentName      string `json:"contentName" pg:"contentName"`
	ContentID        string `json:"contentID" pg:"contentID"`
	ContentType      string `json:"contentType" pg:"contentType"`
	TrustedTimestamp string `json:"trustedTimestamp" pg:"trustedTimestamp"`
}
