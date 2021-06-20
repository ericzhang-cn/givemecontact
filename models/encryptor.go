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
	Text   string `json:"text"`
}

type EncryptResponse struct {
	Text string `json:"text"`
}
