package models

import (
	"github.com/jinzhu/gorm"
	"openseasync/database"
	"openseasync/logs"
)

type Collection struct {
	Slug           string `json:"slug"`             // 集合唯一标识符
	Owner          string `json:"owner"`            // 集合拥有者
	Name           string `json:"name"`             // 集合名称
	BannerImageURL string `json:"banner_image_url"` // 集合背景图
	Description    string `json:"description"`      // 集合描述
	ImageURL       string `json:"image_url"`        // 集合头像
	LargeImageURL  string `json:"large_image_url"`  // 头像大图
	CreateDate     string `json:"create_date"`      // 集合创建时间
}

// InsertOpenSeaCollection find collection through opensea API and insert
func InsertOpenSeaCollection(collections *OwnerCollection, user string) error {
	db := database.GetDB()

	for _, v := range collections.Collections {
		var collection = Collection{
			Slug:           v.Slug,
			Owner:          user,
			Name:           v.Name,
			Description:    v.Description,
			BannerImageURL: v.BannerImageURL,
			ImageURL:       v.ImageURL,
			LargeImageURL:  v.LargeImageURL,
			CreateDate:     v.CreatedDate,
		}

		// gorm v1  batch insert is not supported
		var tmp Collection
		rows := db.Table("collections").
			Where("owner = ? AND slug = ?", user, v.Slug).
			Find(&tmp).RowsAffected
		if rows == 0 {
			if err := db.Table("collections").Create(&collection).Error; err != nil && err != gorm.ErrRecordNotFound {
				logs.GetLogger().Error(err)
				return err
			}
		}

	}
	return nil

}

// FindCollectionByOwner find collections by owner
func FindCollectionByOwner(owner string) ([]*Collection, error) {
	var collections []*Collection
	db := database.GetDB()
	if err := db.Table("collections").Where("owner = ?", owner).Find(&collections).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return collections, nil
}
