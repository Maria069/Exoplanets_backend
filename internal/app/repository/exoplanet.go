package repository

import (
	"context"
	"errors"
	"exoplanets/internal/app/ds"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/minio/minio-go/v7"

	"gorm.io/gorm"
)

func extensionFor(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	case "video/quicktime":
		return ".mov"
	case "video/x-msvideo":
		return ".avi"
	case "video/x-matroska":
		return ".mkv"
	}
	return ""
}

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

func (r *Repository) GetDraftExoplanet(creatorID uint) (*ds.Exoplanet, error) {
	draftExoplanet := ds.Exoplanet{}
	err := r.db.Where("status = ? AND creator_id = ?", "draft", creatorID).Take(&draftExoplanet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return &ds.Exoplanet{}, err
		}
	}

	return &draftExoplanet, nil
}

func (r *Repository) GetNextPublishedExoplanetTo(id int) (ds.Exoplanet, error) {
	var next ds.Exoplanet
	err := r.db.Where("id > ?", id).Order("id ASC").First(&next).Error

	if err != nil {
		return ds.Exoplanet{}, err
	}

	return next, nil
}

func (r *Repository) AddDraftExoplanet(exoplanet *ds.Exoplanet) error {
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

func (r *Repository) DeleteExoplanet(id int) error {
	query := "UPDATE exoplanets SET status = $1 WHERE id = $2"

	err := r.db.Exec(query, "deleted", id).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddOrReplaceExoplanetImage(ExoplanetID uint, header *multipart.FileHeader, ctx context.Context) error {
	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	contentType := http.DetectContentType(buffer)
	filename := strconv.FormatUint(uint64(ExoplanetID), 10) + "_image" + extensionFor(contentType)

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	_, err = r.minio.PutObject(
		ctx,
		r.minio_bucket_name,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		})

	if err != nil {
		return fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	err = r.db.Model(&ds.Exoplanet{}).Where("id = ?", ExoplanetID).UpdateColumn("Image", "http://127.0.0.1:9000/"+r.minio_bucket_name+"/"+filename).Error
	if err != nil {
		r.minio.RemoveObject(
			ctx,
			r.minio_bucket_name,
			filename,
			minio.RemoveObjectOptions{})

		return fmt.Errorf("ошибка сохранения пути к изображению: %w", err)
	}

	return nil
}

func (r *Repository) AddOrReplaceExoplanetVideo(ExoplanetID uint, header *multipart.FileHeader, ctx context.Context) error {
	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	filename := strconv.FormatUint(uint64(ExoplanetID), 10) + "_video" + extensionFor(contentType)

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	if _, err := r.minio.PutObject(
		ctx,
		r.minio_bucket_name,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	); err != nil {
		return fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	url := "http://127.0.0.1:9000/" + r.minio_bucket_name + "/" + filename
	err = r.db.Model(&ds.Exoplanet{}).
		Where("id = ?", ExoplanetID).
		UpdateColumn("video", url).Error
	if err != nil {
		r.minio.RemoveObject(ctx, r.minio_bucket_name, filename, minio.RemoveObjectOptions{})
		return fmt.Errorf("ошибка сохранения пути к видео: %w", err)
	}

	return nil
}
