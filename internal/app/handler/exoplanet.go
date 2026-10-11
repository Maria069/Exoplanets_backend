package handler

import (
	"exoplanets/internal/app/ds"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetExoplanets(ctx *gin.Context) {
	var exoplanets []ds.Exoplanet
	var err error

	searchQuery := ctx.Query("query")

	// получаем значение из поля поиска
	if searchQuery == "" { // если поле поиска пусто, то просто получаем из репозитория все записи
		exoplanets, err = h.Repository.GetPublishedExoplanets()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		weight, err := strconv.ParseFloat(
			strings.Replace(searchQuery, ",", ".", 1), 64,
		)
		if err != nil {
			logrus.Error(err)
		}
		exoplanets, err = h.Repository.GetExoplanetsByWeight(weight) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	likesMap := map[int]int{}

	for _, exoplanet := range exoplanets {
		likesNumber, err := h.Repository.GetLikesNumber(exoplanet.ID)

		if err != nil {
			return
		}

		likesMap[exoplanet.ID] = likesNumber
	}

	ctx.HTML(http.StatusOK, "exoplanet.html", gin.H{
		"exoplanets": exoplanets,
		"likesMap":   likesMap,
		"query":      searchQuery,
	})
}

func (h *Handler) GetExoplanet(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		return
	}

	nextParam := ctx.Query("next")
	var exoplanet ds.Exoplanet

	switch nextParam {
	case "":
		exoplanet, err = h.Repository.GetPublishedExoplanet(id)
	case "true":
		exoplanet, err = h.Repository.GetNextPublishedExoplanetTo(id)
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

	likes, err := h.Repository.GetLikesNumber(exoplanet.ID)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"exoplanet": exoplanet,
		"likes":     likes,
	})
}

func (h *Handler) GetDraftExoplanet(ctx *gin.Context) {
	creatorID := uint(1)

	draftExoplanet, hasDraft, err := h.Repository.GetDraftExoplanet(creatorID)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addExoplanet.html", gin.H{
		"exoplanet": draftExoplanet,
		"hasDraft":  hasDraft,
	})
}

func (h *Handler) AddDraftExoplanet(ctx *gin.Context) {
	title := ctx.PostForm("title")
	creatorID := uint(1)

	draftExoplanet := ds.Exoplanet{
		Status:    "draft",
		Title:     title,
		CreatorID: creatorID,
	}

	err := h.Repository.AddDraftExoplanet(draftExoplanet)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addExoplanet.html", gin.H{
		"exoplanet": draftExoplanet,
		"hasDraft":  true,
	})
}

func (h *Handler) PublishExoplanet(ctx *gin.Context) {
	distanceString := ctx.PostForm("distance")
	weightString := ctx.PostForm("weight")
	description := ctx.PostForm("description")
	creatorID := uint(1)

	distance, err := strconv.Atoi(distanceString)

	if err != nil {
		logrus.Error(err)
	}

	weight, err := strconv.ParseFloat(
		strings.Replace(weightString, ",", ".", 1), 64,
	)

	if err != nil {
		logrus.Error(err)
	}

	exoplanet, _, err := h.Repository.GetDraftExoplanet(creatorID)

	if err != nil {
		logrus.Error(err)
	}

	exoplanet.Status = "published"
	exoplanet.Distance = distance
	exoplanet.Weight = weight
	exoplanet.Description = description

	err = h.Repository.UpdateExoplanet(exoplanet)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addExoplanet.html", gin.H{
		"hasDraft": false,
	})
}

func (h *Handler) DeleteExoplanet(ctx *gin.Context) {
	idString := ctx.PostForm("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.DeleteExoplanet(id)

	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusFound, "/Exoplanets")
}
