package repository

import "exoplanets/internal/app/ds"

func (r *Repository) AddUser(user *ds.User) error {
	err := r.db.Create(user).Error

	if err != nil {
		return err
	}

	return nil
}
