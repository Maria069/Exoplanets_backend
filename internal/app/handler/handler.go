package handler

import (
	"exoplanets/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
//func (h *Handler) RegisterHandler(router *gin.Engine) {
//router.GET("/exoplanets", h.GetAllExoplanets)
// router.GET("/", h.GetExoplanets)
// router.GET("/exoplanet/:id", h.GetExoplanet)
//router.GET("/", h.GetDraftExoplanet)
//}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/Exoplanets", h.GetExoplanets)
	router.GET("/Exoplanets/:id", h.GetExoplanet)
	router.GET("/AddExoplanet", h.GetDraftExoplanet)
	router.POST("/AddExoplanet", h.AddDraftExoplanet)
	router.POST("/PublishExoplanet", h.PublishExoplanet)
	router.POST("/DeleteExoplanet", h.DeleteExoplanet)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("./templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
