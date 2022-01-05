package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"openseasync/database"
	"openseasync/logs"
	"time"
)

var CONNOT_DELETE_COLLECTION_ERR = errors.New("Cannot delete a collection that has an asset")

type Collection struct {
	//ID              int64   `json:"id" bson:"id"`                               // 主键
	Slug            string  `json:"slug" bson:"slug"`                           // 集合唯一标识符
	UserAddress     string  `json:"user_address" bson:"user_address"`           // 集合拥有者
	Name            string  `json:"name" bson:"name"`                           // 集合名称
	BannerImageURL  string  `json:"banner_image_url" bson:"banner_image_url"`   // 集合背景图
	Description     string  `json:"description" bson:"description"`             // 集合描述
	ImageURL        string  `json:"image_url" bson:"image_url"`                 // 集合头像
	LargeImageURL   string  `json:"large_image_url" bson:"large_image_url"`     // 头像大图
	IsDelete        int8    `json:"is_delete" bson:"is_delete"`                 // 是否删除 1删除 0未删除 默认为0
	CreateDate      string  `json:"create_date" bson:"create_date"`             // 集合创建时间
	RefreshTime     int     `json:"refresh_time" bson:"refresh_time"`           // 刷新时间
	NumOwners       int     `json:"num_owners" bson:"num_owners"`               // 集合中输入自己的NFT个数
	TotalSupply     int     `json:"total_supply" bson:"total_supply"`           // 集合中NFT总数
	TotalVolume     float64 `json:"total_volume" bson:"total_volume"`           // 交易量
	OwnedAssetCount string  `json:"owned_asset_count" bson:"owned_asset_count"` // 所有NFT中属于自己的NFT个数 此地段可能是个big int, 所以采用string存储
}

// InsertOpenSeaCollection find collection through opensea API and insert
func InsertOpenSeaCollection(collections *OwnerCollection, user string) error {
	db := database.GetMongoClient()
	refreshTime := int(time.Now().Unix())
	for _, v := range collections.Collections {
		var collection = Collection{
			Slug:            v.Slug,
			UserAddress:     user,
			Name:            v.Name,
			Description:     v.Description,
			BannerImageURL:  v.BannerImageURL,
			ImageURL:        v.ImageURL,
			LargeImageURL:   v.LargeImageURL,
			CreateDate:      v.CreatedDate,
			RefreshTime:     refreshTime,
			NumOwners:       v.Stats.NumOwners,
			TotalSupply:     int(v.Stats.TotalSupply),
			TotalVolume:     v.Stats.TotalVolume,
			OwnedAssetCount: v.OwnedAssetCount.String(),
		}
		// gorm v1  batch insert is not supported

		count, err := db.Collection("collections").
			CountDocuments(context.TODO(), bson.M{"user_address": user, "slug": v.Slug, "is_delete": 0})
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		if count == 0 {
			if _, err = db.Collection("collections").InsertOne(context.TODO(), &collection); err != nil {
				logs.GetLogger().Error(err)
				return err
			}
		} else {
			// update
			if _, err = db.Collection("collections").UpdateOne(
				context.TODO(),
				bson.M{"user_address": user, "slug": v.Slug, "is_delete": 0},
				bson.M{"$set": bson.M{"name": v.Name, "description": v.Description, "banner_image_url": v.BannerImageURL, "image_url": v.ImageURL,
					"large_image_url": v.LargeImageURL, "refresh_time": refreshTime}}); err != nil {
				logs.GetLogger().Error(err)
				return err
			}

		}

	}

	if _, err := db.Collection("collections").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "refresh_time": bson.M{"$lt": refreshTime}, "is_delete": 0},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	return nil

}

// FindCollectionByOwner find collections by owner
func FindCollectionByOwner(user string) ([]*Collection, error) {
	var collections []*Collection
	db := database.GetMongoClient()
	cursor, err := db.Collection("collections").Find(context.TODO(), bson.M{"user_address": user, "is_delete": 0})
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &collections); err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return collections, nil
}

// DeleteCollectionBySlug delete empty collection
func DeleteCollectionBySlug(user, slug string) error {
	db := database.GetMongoClient()
	row, err := db.Collection("assets").CountDocuments(context.TODO(), bson.M{"user_address": user, "slug": slug, "is_delete": 0})
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	// There cannot be an asset in the collection otherwise deleting failed
	if row >= 1 {
		return CONNOT_DELETE_COLLECTION_ERR
	}

	if _, err := db.Collection("collections").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "slug": slug, "is_delete": 0},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}
