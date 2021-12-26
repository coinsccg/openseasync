package models

import (
	"openseasync/database"
)

type Collection struct {
	Slug           string `json:"slug"`             // 集合唯一标识符
	Owner          string `json:"owner"`            // 集合拥有者
	Name           string `json:"name"`             // 集合名称
	BannerImageURL string `json:"banner_image_url"` // 集合背景图
	Description    string `json:"description"`      // 集合描述
	ImageURL       string `json:"image_url"`        // 集合头像
	LargeImageURL  string `json:"large_image_url"`  // 头像大图
}

// InsertOpenSeaCollection find collection through opensea API and insert
func InsertOpenSeaCollection(collections *OwnerCollection, owner string) error {
	db := database.GetDB()

	for _, v := range collections.Collections {
		var collection = Collection{
			Slug:           v.Slug,
			Owner:          owner,
			Name:           v.Name,
			Description:    v.Description,
			BannerImageURL: v.BannerImageURL,
			ImageURL:       v.ImageURL,
			LargeImageURL:  v.LargeImageURL,
		}

		// gorm v1  batch insert is not supported
		var tmp Collection
		rows := db.Table("collections").
			Where("owner = ? AND slug = ?", owner, v.Slug).
			Find(&tmp).RowsAffected
		if rows == 0 {
			if err := db.Table("collections").Create(&collection).Error; err != nil {
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
	if err := db.Table("collections").Where("owner = ?", owner).Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}
