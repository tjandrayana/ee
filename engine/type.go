package engine

type User struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	SecureAddress string `json:"secure_address"`
}
