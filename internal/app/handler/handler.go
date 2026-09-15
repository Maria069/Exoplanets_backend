package handler

import (
	"exoplanets/internal/app/repository"
	"net/http"
	"strconv"

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

func (h *Handler) GetExoplanets(ctx *gin.Context) {
	var exoplanets []repository.Exoplanet
	var err error
	searchQuery := ctx.Query("query")

	if searchQuery == "" {
		exoplanets, err = h.Repository.GetExoplanets()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		temperature, err := strconv.Atoi(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
		exoplanets, err = h.Repository.GetExoplanetsByTemperature(temperature)
		if err != nil {
			logrus.Error(err)
		}
	}
	ctx.HTML(http.StatusOK, "exoplanet.html", gin.H{
		"exoplanets": exoplanets,
		"query":      searchQuery,
	})
}

// func (h *Handler) GetExoplanet(ctx *gin.Context) {
// 	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
// 	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
// 	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	exoplanet, err := h.Repository.GetExoplanet(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	ctx.HTML(http.StatusOK, "exoplanet.html", gin.H{
// 		"exoplanet": exoplanet,
// 	})
// }

func (h *Handler) GetExoplanet(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		return
	}
	nextParam := ctx.Query("next")
	var exoplanet repository.Exoplanet

	switch nextParam {
	case "":
		exoplanet, err = h.Repository.GetExoplanet(id)
	case "true":
		exoplanet, err = h.Repository.GetNextPublishedExoplanet(id)
		if err == nil {
			ctx.Redirect(http.StatusFound, "/Exoplanets/"+strconv.Itoa(exoplanet.ID))
			return
		}
	default:
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "параметр next должен быть true или отсутствовать")
		return
	}
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "запись не найдена")
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"exoplanet": exoplanet,
	})
}

func (h *Handler) GetDraftExoplanet(ctx *gin.Context) {
	exoplanet, err := h.Repository.GetDraftExoplanet()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addExoplanet.html", gin.H{
		"exoplanet": exoplanet,
	})
}
