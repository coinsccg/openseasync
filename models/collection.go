package models

import (
	"context"
	"errors"
	"openseasync/common/utils"
	"openseasync/database"
	"openseasync/logs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var CONNOT_DELETE_COLLECTION_ERR = errors.New("Cannot delete a collection that has an asset")

// InsertOpenSeaCollection find collection through opensea API and insert
func InsertOpenSeaCollection(collections *OwnerCollection, user string, refreshTime int64) error {
	db := database.GetMongoClient()

	for _, v := range collections.Collections {
		var collection = Collection{
			ID:                 v.Slug,
			UserMetamaskID:     user,
			CreatorMetamaskId:  v.PayoutAddress,
			CollectionName:     v.Name,
			Description:        v.Description,
			UserCoverUrl:       v.BannerImageURL,
			CoverImageUrl:      v.ImageURL,
			CoverLargeImageUrl: v.LargeImageURL,
			CreateDate:         utils.ParseTime(v.CreatedDate),
			RefreshTime:        refreshTime,
			OwnersCount:        v.Stats.NumOwners,
			ItemsCount:         int(v.Stats.TotalSupply),
			TotalVolume:        v.Stats.TotalVolume,
		}
		if collection.CreatorMetamaskId == "" {
			collection.CreatorMetamaskId = user
		}
		count, err := db.Collection("collections").
			CountDocuments(context.TODO(), bson.M{"userMetamaskId": user, "id": v.Slug, "isDelete": 0})
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
			collectionByte, err := bson.Marshal(collection)
			if err != nil {
				logs.GetLogger().Error(err)
				return err
			}
			var tmpCollection bson.M
			if err := bson.Unmarshal(collectionByte, &tmpCollection); err != nil {
				logs.GetLogger().Error(err)
				return err
			}
			if _, err = db.Collection("collections").UpdateOne(
				context.TODO(),
				bson.M{"userMetamaskId": user, "id": v.Slug, "isDelete": 0},
				bson.M{"$set": tmpCollection}); err != nil {
				logs.GetLogger().Error(err)
				return err
			}
		}
	}

	if _, err := db.Collection("collections").UpdateMany(
		context.TODO(),
		bson.M{"userMetamaskId": user, "refreshTime": bson.M{"$lt": refreshTime}, "isDelete": 0},
		bson.M{"$set": bson.M{"isDelete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	return nil

}

// FindCollectionByUserMetamaskID find collections by usermetamaskid
func FindCollectionByUserMetamaskID(userMetamaskId string, page, pageSize int64) (map[string]interface{}, error) {
	var (
		collections = make([]ResponseColl, 0)
		result      = make(map[string]interface{})
	)
	db := database.GetMongoClient()
	total, err := db.Collection("collections").CountDocuments(context.TODO(), bson.M{"userMetamaskId": userMetamaskId, "isDelete": 0})
	if err != nil && err != mongo.ErrNoDocuments {
		logs.GetLogger().Error(err)
		return nil, err
	}
	totalPage := total / pageSize
	if total%pageSize != 0 {
		totalPage++
	}

	pipe := mongo.Pipeline{
		{{"$match", bson.M{"userMetamaskId": userMetamaskId, "isDelete": 0}}},
		{{"$skip", (page - 1) * pageSize}},
		{{"$limit", pageSize}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "creatorMetamaskId",
			"foreignField": "userMetamaskId",
			"as":           "item",
		}}},
		{{
			"$addFields", bson.M{"user_item": bson.M{"$arrayElemAt": bson.A{"$item", 0}}},
		}},
		{{
			"$addFields", bson.M{"userId": "$user_item.id", "userName": "$user_item.userName", "avatarUrl": "$user_item.avatarUrl"},
		}},
		{{"$project",
			bson.M{
				"_id": 0, "id": 1, "userId": 1, "userMetamaskId": 1, "coverImageUrl ": 1, "avatarUrl": 1, "userName": 1,
				"collectionName": 1, "description": 1,
			},
		}},
	}

	cursor, err := db.Collection("collections").Aggregate(context.TODO(), pipe)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &collections); err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	result["data"] = collections
	result["metadata"] = map[string]int64{"page": page, "pageSize": pageSize, "total": total, "totalPage": totalPage}

	return result, nil
}

// FindCollectionByCollectionID find collections by collectionId
func FindCollectionByCollectionID(collectionId string) (map[string]interface{}, error) {
	var (
		collection  bson.M
		collections = make([]bson.M, 0)
		result      = make(map[string]interface{})
	)
	db := database.GetMongoClient()
	pipe := mongo.Pipeline{
		{{"$match", bson.M{"id": collectionId, "isDelete": 0}}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "creatorMetamaskId",
			"foreignField": "userMetamaskId",
			"as":           "item",
		}}},
		{{
			"$addFields", bson.M{"user_item": bson.M{"$arrayElemAt": bson.A{"$item", 0}}},
		}},
		{{
			"$addFields", bson.M{"userId": "$user_item.id", "userName": "$user_item.userName", "avatarUrl": "$user_item.avatarUrl"},
		}},
		{{"$project",
			bson.M{"_id": 0, "id": 1, "userId": 1, "userMetamaskId": 1, "userCoverUrl": 1, "avatarUrl": 1,
				"userName": 1, "itemsCount": 1, "ownersCount": 1, "floorPrice": 1, "highestPrice": 1,
				"collectionName": 1, "likesCount": 1, "viewsCount": 1, "description": 1}}},
	}

	cursor, err := db.Collection("collections").Aggregate(context.TODO(), pipe)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &collections); err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if len(collections) >= 1 {
		collection = collections[0]
		result["data"] = collection
	} else {
		result["data"] = ResponseCollection{}
	}

	return result, nil
}

// FindUserMediaByUserId find user media by userId
func FindUserMediaByUserId(userMetamaskId string) (interface{}, error) {
	var user User
	db := database.GetMongoClient()
	if err := db.Collection("users").FindOne(context.TODO(), bson.M{"userMetamaskId": userMetamaskId}).Decode(&user); err != nil && err != mongo.ErrNoDocuments {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return user, nil
}

// DeleteCollectionByCollectionId delete empty collection
func DeleteCollectionByCollectionId(user, slug string) error {
	db := database.GetMongoClient()
	row, err := db.Collection("assets").CountDocuments(
		context.TODO(), bson.M{"userMetamaskId": user, "collectionId": slug, "is_delete": 0})
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
		bson.M{"userMetamaskId": user, "id": slug, "isDelete": 0},
		bson.M{"$set": bson.M{"isDelete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}
