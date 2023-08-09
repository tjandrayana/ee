package engine

type User struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type SecureUser struct {
	Name              string `json:"bmFtZQ=="`
	Email             string `json:"ZW1haWw="`
	Address           string `json:"YWRkcmVzcw=="`
	PhoneNumber       string `json:"cGhvbmVfbnVtYmVy"`
	SecureName        string `json:"c2VjdXJlX25hbWU="`
	SecureEmail       string `json:"c2VjdXJlX2VtYWls"`
	SecureAddress     string `json:"c2VjdXJlX2FkZHJlc3M="`
	SecurePhoneNumber string `json:"c2VjdXJlX3Bob25lX251bWJlcg=="`
}
