package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"RIP2026/internal/app/models"
	minioClient "RIP2026/internal/app/minioClient"
	"RIP2026/internal/app/serializer"
)

func (r *Repository) GetBatteryTypesAPI() ([]models.BatteryType, error) {
	var list []models.BatteryType
	err := r.db.Where("is_deleted = ?", false).Find(&list).Error
	return list, err
}

func (r *Repository) GetBatteryTypeAPI(id int) (*models.BatteryType, error) {
	var bt models.BatteryType
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&bt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: тип аккумулятора с id %d", ErrNotFound, id)
		}
		return nil, err
	}
	return &bt, nil
}

func (r *Repository) GetBatteryTypesByTitleAPI(title string) ([]models.BatteryType, error) {
	var list []models.BatteryType
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+title+"%", false).Find(&list).Error
	return list, err
}

func (r *Repository) CreateBatteryTypeAPI(j serializer.BatteryTypeJSON) (models.BatteryType, error) {
	bt := serializer.BatteryTypeFromJSON(j)
	if err := r.db.Create(&bt).Error; err != nil {
		return models.BatteryType{}, err
	}
	return bt, nil
}

func (r *Repository) AddPhoto(ctx *gin.Context, batteryTypeID int, file *multipart.FileHeader) (*models.BatteryType, error) {
	if r.mc == nil {
		return nil, fmt.Errorf("minio not configured")
	}
	bt, err := r.GetBatteryTypeAPI(batteryTypeID)
	if err != nil {
		return nil, err
	}
	if bt.Photo != nil && strings.TrimSpace(*bt.Photo) != "" {
		_ = minioClient.DeleteObject(ctx.Request.Context(), r.mc, minioClient.GetBucket(), *bt.Photo)
	}
	fileName, err := minioClient.UploadImage(ctx.Request.Context(), r.mc, minioClient.GetBucket(), file, bt.ID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.BatteryType{}).Where("id = ?", batteryTypeID).Update("photo", fileName).Error; err != nil {
		return nil, err
	}
	bt.Photo = &fileName
	return bt, nil
}

func (r *Repository) AddVideo(ctx *gin.Context, batteryTypeID int, file *multipart.FileHeader) (*models.BatteryType, error) {
	if r.mc == nil {
		return nil, fmt.Errorf("minio not configured")
	}
	bt, err := r.GetBatteryTypeAPI(batteryTypeID)
	if err != nil {
		return nil, err
	}
	if bt.Video != "" {
		_ = minioClient.DeleteObject(ctx.Request.Context(), r.mc, minioClient.GetBucket(), bt.Video)
	}
	fileName, err := minioClient.UploadVideo(ctx.Request.Context(), r.mc, minioClient.GetBucket(), file, bt.ID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.BatteryType{}).Where("id = ?", batteryTypeID).Update("video", fileName).Error; err != nil {
		return nil, err
	}
	bt.Video = fileName
	return bt, nil
}

func (r *Repository) EnsureMinioBucket(ctx context.Context) error {
	if r.mc == nil {
		return nil
	}
	exists, err := r.mc.BucketExists(ctx, minioClient.GetBucket())
	if err != nil {
		return err
	}
	if !exists {
		return r.mc.MakeBucket(ctx, minioClient.GetBucket(), minio.MakeBucketOptions{})
	}
	return nil
}
