package api

import (
	"exoplanets/internal/app/handler"
	"exoplanets/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("../../templates/*")
	r.Static("/static", "../../resources")
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/Exoplanets", handler.GetExoplanets)
	r.GET("/Exoplanets/:id", handler.GetExoplanet)
	r.GET("/AddExoplanet", handler.GetDraftExoplanet)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
