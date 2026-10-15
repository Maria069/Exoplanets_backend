package handler

import (
	"exoplanets/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(ctx *gin.Context) {
	newUser := ds.User{
		Login:    ctx.Request.FormValue("login"),
		Password: ctx.Request.FormValue("password"),
	}

	err := h.Repository.AddUser(&newUser)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user": newUser,
	})
}

func (h *Handler) SignIn(ctx *gin.Context) {
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "авторизация",
	// })
}

func (h *Handler) SignOut(ctx *gin.Context) {
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "деавторизация",
	// })
}
