package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/ericzhang-cn/givemecontact/models"
	"github.com/ericzhang-cn/givemecontact/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

// Create new encryptor
// @Summary Create new encryptor
// @Produce  json
// @Success 201 {object} models.EncryptorCreateResponse
// @Failure 500 {object} models.HttpError
// @Router /endpoint/v1/encryptor/ [post]
// @Tags Encryptor
func CreateEncryptor(c *gin.Context) {
	var encryptor models.Encryptor

	reader := rand.Reader
	size := viper.GetInt("encryptor.rsa_size")
	key, err := rsa.GenerateKey(reader, size)
	if err != nil {
		log.Printf("generate rsa key pair failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}
	public := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	})
	private := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	encryptor.Phrase = utils.RandStringRunes(viper.GetInt("encryptor.phrase_length"))
	encryptor.PublicKey = string(public)

	db := utils.GetDB()
	if err := db.WithContext(c).Create(&encryptor).Error; err != nil {
		log.Printf("insert rsa public key to sqlite failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	var resp models.EncryptorCreateResponse
	resp.PublicKey = encryptor.PublicKey
	resp.PrivateKey = string(private)
	resp.Phrase = encryptor.Phrase

	c.JSON(http.StatusCreated, resp)
}

// Encrypt plain text
// @Summary Encrypt plain text
// @Accept json
// @Produce  json
// @Param request body models.EncryptRequest true "encrypt request body"
// @Success 200 {object} models.EncryptResponse
// @Failure 400 {object} models.HttpError
// @Failure 404 {object} models.HttpError
// @Failure 500 {object} models.HttpError
// @Router /endpoint/v1/encryptor/enc/ [post]
// @Tags Encryptor
func Encrypt(c *gin.Context) {

}

// Decrypt cipher text
// @Summary Decrypt cipher text
// @Accept json
// @Produce  json
// @Param request body models.DecryptRequest true "decrypt request body"
// @Success 200 {object} models.DecryptResponse
// @Failure 400 {object} models.HttpError
// @Failure 500 {object} models.HttpError
// @Router /endpoint/v1/encryptor/dec/ [post]
// @Tags Encryptor
func Decrypt(c *gin.Context) {

}
