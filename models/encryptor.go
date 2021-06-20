package models

import "github.com/jinzhu/gorm"

type Encryptor struct {
	gorm.Model
	Phrase    string `gorm:"column:phrase" json:"phrase"`
	PublicKey string `gorm:"column:public_key" json:"publicKey"`
}

type EncryptorCreateResponse struct {
	Phrase     string `json:"phrase"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type EncryptRequest struct {
	Phrase string `json:"phrase"`
	Text   string `json:"text"`
}

type EncryptResponse struct {
	Text string `json:"text"`
}

type DecryptRequest struct {
	PrivateKey string `json:"privateKey"`
	Text       string `json:"text"`
}

type DecryptResponse struct {
	Text string `json:"text"`
}
