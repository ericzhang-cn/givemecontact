package models

import "github.com/jinzhu/gorm"

type Message struct {
	gorm.Model
	Phrase     string `gorm:"column:phrase" json:"phrase"`
	CipherText string `gorm:"column:cipher_text" json:"cipherText"`
}

type MessageCreateRequest struct {
	Text string `json:"text"`
}

type MessageCreateResponse struct {
	Phrase string `json:"phrase"`
	Text   string `json:"text"`
}

type DecryptRequest struct {
	PrivateKey string `json:"privateKey"`
}

type DecryptResponse struct {
	Text string `json:"text"`
}
