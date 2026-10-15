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

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/api/Exoplanets", h.GetExoplanets)
	router.GET("/api/Exoplanets/:id", h.GetExoplanet)
	router.GET("/api/DraftExoplanets", h.GetDraftExoplanet)
	router.POST("/api/AddExoplanets", h.AddDraftExoplanet)
	router.PUT("/api/PublishExoplanets", h.PublishExoplanet)
	router.DELETE("/api/DeleteExoplanets/:id", h.DeleteExoplanet)
	router.POST("/api/LikeExoplanets/:id/like", h.LikeExoplanet)

	router.POST("/api/SignUp", h.SignUp)
	router.POST("/api/SignIn", h.SignIn)
	router.POST("/api/SignOut", h.SignOut)
}

// func (h *Handler) RegisterStatic(router *gin.Engine) {
// 	router.LoadHTMLGlob("./templates/*")
// 	router.Static("/static", "./resources")
// }

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
