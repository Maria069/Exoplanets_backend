package repository

import (
	"errors"
	"exoplanets/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetPublishedExoplanet(id int) (ds.Exoplanet, error) {
	exoplanet := ds.Exoplanet{}
	err := r.db.Where("id = ? AND status = ?", id, "published").First(&exoplanet).Error

	if err != nil {
		return ds.Exoplanet{}, err
	}

	return exoplanet, nil
}

func (r *Repository) GetPublishedExoplanets() ([]ds.Exoplanet, error) {
	var exoplanets []ds.Exoplanet
	err := r.db.Where("status = ?", "published").Find(&exoplanets).Error

	if err != nil {
		return nil, err
	}

	return exoplanets, nil
}

func (r *Repository) GetExoplanetsByWeight(weight float64) ([]ds.Exoplanet, error) {
	var exoplanets []ds.Exoplanet
	err := r.db.Where("status = ? AND weight = ?", "published", weight).Find(&exoplanets).Error

	if err != nil {
		return nil, err
	}

	return exoplanets, nil
}

func (r *Repository) GetDraftExoplanet(creatorID uint) (ds.Exoplanet, bool, error) {
	draftExoplanet := ds.Exoplanet{}
	err := r.db.Where("status = ? AND creator_id = ?", "draft", creatorID).Take(&draftExoplanet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Exoplanet{}, false, nil
		} else {
			return ds.Exoplanet{}, false, err
		}
	}

	return draftExoplanet, true, nil
}

func (r *Repository) GetNextPublishedExoplanetTo(id int) (ds.Exoplanet, error) {
	var next ds.Exoplanet
	err := r.db.Where("id > ?", id).Order("id ASC").First(&next).Error

	if err != nil {
		return ds.Exoplanet{}, err
	}

	return next, nil
}

func (r *Repository) AddDraftExoplanet(exoplanet ds.Exoplanet) error {
	err := r.db.Create(&exoplanet).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateExoplanet(exoplanet ds.Exoplanet) error {
	err := r.db.Save(&exoplanet).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLikesNumber(id int) (int, error) {
	var count int64
	err := r.db.Model(&ds.Likes{}).
		Where("exoplanet_id = ?", id).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) DeleteExoplanet(id int) error {
	query := "UPDATE exoplanets SET status = $1 WHERE id = $2"

	err := r.db.Exec(query, "deleted", id).Error

	if err != nil {
		return err
	}

	return nil
}
