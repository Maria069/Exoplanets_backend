package repository

import (
	"exoplanets/internal/app/ds"
)

func (r *Repository) AddLike(like *ds.Likes) error {
	err := r.db.Create(like).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteLike(userID uint, exoplanetID uint) error {
	err := r.db.
		Where("user_id = ? AND exoplanet_id = ?", userID, exoplanetID).
		Delete(&ds.Likes{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) IsLikedByUser(exoplanetId int, userId int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM likes WHERE user_id = $1 AND exoplanet_id = $2)"
	row := r.db.Raw(query, userId, exoplanetId).Row()

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
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
