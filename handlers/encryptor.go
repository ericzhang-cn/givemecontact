package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/ericzhang-cn/givemecontact/models"
	"github.com/ericzhang-cn/givemecontact/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// Create new encryptor
// @Summary Create new encryptor
// @Produce  json
// @Success 201 {object} models.EncryptorCreateResponse
// @Failure 500 {object} models.HttpError
// @Router /endpoint/v1/encryptors/ [post]
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
// @Param id path string true "encryptor id"
// @Param request body models.EncryptRequest true "encrypt request body"
// @Success 200 {object} models.EncryptResponse
// @Failure 400 {object} models.HttpError
// @Failure 404 {object} models.HttpError
// @Failure 500 {object} models.HttpError
// @Router /endpoint/v1/encryptors/{id}/encrypt/ [post]
// @Tags Encryptor
func Encrypt(c *gin.Context) {
	var req models.EncryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("bind request body failed, error message: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.HttpError{Error: err.Error()})
		return
	}

	var encryptor models.Encryptor
	db := utils.GetDB()
	err := db.Model(&models.Encryptor{}).Where("phrase = ?", c.Param("id")).First(&encryptor).Error
	if err == gorm.ErrRecordNotFound {
		log.Printf("can not find phrase: %s, error message: %s", c.Param("id"), err.Error())
		c.JSON(http.StatusNotFound, models.HttpError{Error: err.Error()})
		return
	}
	if err != nil {
		log.Printf("query encryptor failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	block, _ := pem.Decode([]byte(encryptor.PublicKey))
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		log.Printf("decode public key failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(req.Text))
	if err != nil {
		log.Printf("encrypt failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.EncryptResponse{Text: base64.StdEncoding.EncodeToString(cipher)})
}
