package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Exoplanet struct {
	ID          int
	Title       string
	Status      string
	Distance    float32
	Temperature int
	Description string
	Likes       []int
	ImagePath   string
	VideoPath   string
}

func (r *Repository) GetExoplanets() ([]Exoplanet, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	exoplanets := []Exoplanet{ // массив элементов из наших структур
		{
			ID:          1,
			Title:       "Горячий юпитер",
			Status:      "published",
			Distance:    1,
			Temperature: 5500,
			Description: "Масса такой экзопланеты сравнима с массой Юпитера, но расстояние до звезды на порядок меньше.",
			Likes:       []int{},
			ImagePath:   "hot_jupiter_cut.png",
			VideoPath:   "hot_jupiter.mp4",
		},
		{
			ID:          2,
			Title:       "Суперземля",
			Status:      "published",
			Distance:    2.8,
			Temperature: 6000,
			Description: "Масса суперземли превышает Земную массу, но меньше массы Нептуна.",
			Likes:       []int{},
			ImagePath:   "super_earth_cut.png",
			VideoPath:   "super_earth.mp4",
		},
		{
			ID:          3,
			Title:       "Ледяной гигант",
			Status:      "published",
			Distance:    19.9,
			Temperature: 5500,
			Description: "Класс планет-гигантов, состоящих из веществ тяжелее водорода и галия.",
			Likes:       []int{},
			ImagePath:   "ice_giant_cut.png",
			VideoPath:   "ice_giant.mp4",
		},
		{
			ID:          4,
			Title:       "Пустынная планета",
			Status:      "deleted",
			Distance:    0.5,
			Temperature: 5000,
			Description: "Планета с одним типом биома по всей поверхности, чаще всего имеет пустынный биом и минимальное количество осадков.",
			Likes:       []int{},
			ImagePath:   "desert_planet_cut.png",
			VideoPath:   "desert_planet.mp4",
		},
		{
			ID:          5,
			Title:       "Железная планета",
			Status:      "draft",
			Distance:    4.5,
			Temperature: 3000,
			Description: "Разновидность планет Земного типа, у которых есть насыщенное железом ядро.",
			Likes:       []int{},
			ImagePath:   "iron_planet_cut.png",
			VideoPath:   "iron_planet.mp4",
		},
	}

	if len(exoplanets) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return exoplanets, nil
}

func (r *Repository) GetExoplanet(id int) (Exoplanet, error) {
	exoplanets, err := r.GetExoplanets()
	if err != nil {
		return Exoplanet{}, err
	}

	for _, exoplanet := range exoplanets {
		if exoplanet.ID == id {
			return exoplanet, nil
		}
	}
	return Exoplanet{}, fmt.Errorf("планета не найдена")
}

func (r *Repository) GetExoplanetsByTitle(title string) ([]Exoplanet, error) {
	exoplanets, err := r.GetExoplanets()
	if err != nil {
		return []Exoplanet{}, err
	}

	var result []Exoplanet
	for _, exoplanet := range exoplanets {
		if strings.Contains(strings.ToLower(exoplanet.Title), strings.ToLower(title)) {
			result = append(result, exoplanet)
		}
	}

	return result, nil
}

func (r *Repository) GetExoplanetsByTemperature(temperature int) ([]Exoplanet, error) {
	exoplanets, err := r.GetExoplanets()
	if err != nil {
		return []Exoplanet{}, err
	}

	var result []Exoplanet
	for _, exoplanet := range exoplanets {
		if exoplanet.Temperature == temperature {
			result = append(result, exoplanet)
		}
	}

	return result, nil
}

func (r *Repository) GetPublishedExoplanets() ([]Exoplanet, error) {
	exoplanets, err := r.GetExoplanets()
	if err != nil {
		return []Exoplanet{}, err
	}

	var result []Exoplanet
	for _, exoplanet := range exoplanets {
		if exoplanet.Status == "published" {
			result = append(result, exoplanet)
		}
	}

	return result, nil
}

func (r *Repository) GetDraftExoplanet() (Exoplanet, error) {
	exoplanets, err := r.GetExoplanets()
	if err != nil {
		return Exoplanet{}, err
	}

	var result Exoplanet
	for _, exoplanet := range exoplanets {
		if exoplanet.Status == "draft" {
			result = exoplanet
		}
	}

	return result, nil
}

func (r *Repository) GetNextPublishedExoplanet(id int) (Exoplanet, error) {
	exoplanets, err := r.GetPublishedExoplanets()
	if err != nil {
		return Exoplanet{}, err
	}

	var currentIndex int
	for i, exoplanet := range exoplanets {
		if exoplanet.ID == id {
			currentIndex = i
		}
	}

	if (currentIndex + 1) == len(exoplanets) {
		return exoplanets[0], nil
	}

	var result Exoplanet = exoplanets[currentIndex+1]

	return result, nil
}
