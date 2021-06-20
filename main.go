package main

import (
	_ "github.com/ericzhang-cn/givemecontact/api/docs"
	"github.com/ericzhang-cn/givemecontact/handlers"
	"github.com/ericzhang-cn/givemecontact/models"
	"github.com/ericzhang-cn/givemecontact/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title GiveMeContact Endpoint
// @version 1.0
// @description GiveMeContact Endpoint
// @termsOfService http://swagger.io/terms/

// @host givemecontact.io
// @BasePath /
func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("conf/")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read configuare file failed, error message: %s", err.Error())
	}

	db := utils.GetDB()
	if err := db.AutoMigrate(&models.Encryptor{}); err != nil {
		log.Fatalf("migrate database failed, error message: %s", err.Error())
	}

	r := gin.Default()
	r.POST("/endpoint/v1/encryptor/", handlers.CreateEncryptor)
	api := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, api))
	if err := r.Run(); err != nil {
		log.Fatalf("http server start failed, error message: %s", err.Error())
	}
}
