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

// Create new message
// @Summary Create new message
// @Accept  json
// @Produce  json
// @Param request body models.MessageCreateRequest true "create message request body"
// @Success 201 {object} models.MessageCreateResponse
// @Failure 400 {object} models.HttpError
// @Failure 500 {object} models.HttpError
// @Router /endpoint/v1/messages/ [post]
// @Tags Message
func CreateMessage(c *gin.Context) {
	var req models.MessageCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("bind request body failed, error message: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.HttpError{Error: err.Error()})
		return
	}

	var message models.Message
	message.CipherText = req.Text
	message.Phrase = utils.RandStringRunes(viper.GetInt("encryptor.phrase_length"))

	db := utils.GetDB()
	if err := db.WithContext(c).Create(&message).Error; err != nil {
		log.Printf("insert message to sqlite failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	var resp models.MessageCreateResponse
	resp.Phrase = message.Phrase
	resp.Text = message.CipherText

	c.JSON(http.StatusCreated, resp)
}

// Decrypt cipher text
// @Summary Decrypt cipher text
// @Accept json
// @Produce  json
// @Param id path string true "encryptor id"
// @Param request body models.DecryptRequest true "decrypt request body"
// @Success 200 {object} models.DecryptResponse
// @Failure 400 {object} models.HttpError
// @Failure 404 {object} models.HttpError
// @Failure 500 {object} models.HttpError
// @Router /endpoint/v1/messages/{id}/decrypt/ [post]
// @Tags Message
func Decrypt(c *gin.Context) {
	var req models.DecryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("bind request body failed, error message: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.HttpError{Error: err.Error()})
		return
	}

	var message models.Message
	db := utils.GetDB()
	err := db.Model(&models.Message{}).Where("phrase = ?", c.Param("id")).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		log.Printf("can not find message: %s, error message: %s", c.Param("id"), err.Error())
		c.JSON(http.StatusNotFound, models.HttpError{Error: err.Error()})
		return
	}
	if err != nil {
		log.Printf("query message failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	block, _ := pem.Decode([]byte(req.PrivateKey))
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("decode private key failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	decodedText, err := base64.StdEncoding.DecodeString(message.CipherText)
	if err != nil {
		log.Printf("base64 decode failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}
	plain, err := rsa.DecryptPKCS1v15(rand.Reader, pri, decodedText)
	if err != nil {
		log.Printf("decrypt failed, error message: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.HttpError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.DecryptResponse{Text: string(plain)})
}
