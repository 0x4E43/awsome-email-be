package awmail

type EmailConfig struct {
	SmtpHOST  string `json:"smtpHost"`
	SmtpPass  string `json:"smtpPass"`
	SmtpPort  int    `json:"smtpPort"`
	SmtpFrom  string `json:"smtpFrom"`
	IsDefault bool   `json:"isDefault"`
}
