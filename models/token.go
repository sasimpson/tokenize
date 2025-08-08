package models

type Token struct {
	BaseModel
	CreateToken
	Token string `json:"token"`
}

func (t Token) Encrypt() {
	t.Payload = "encrypted"
}

func (t Token) Decrypt() {
	t.Payload = "decrypted"
}

type CreateToken struct {
	Payload   string         `json:"payload"`
	TokenType string         `json:"token_type"`
	TTL       int64          `json:"ttl"`
	Metadata  map[string]any `json:"metadata"`
}
