package models

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"openseasync/common/utils"
	"openseasync/database"
	"openseasync/logs"
	"time"
)

type Asset struct {
	//ID                int64  `json:"id" bson:"id"`                                   // 主键
	UserAddress       string `json:"user_address" bson:"user_address"`               // 用户地址
	Title             string `json:"title" bson:"title"`                             // NFT作品标题
	ImageURL          string `json:"image_url" bson:"image_url"`                     // NFT作品图片
	ImagePreviewURL   string `json:"image_preview_url" bson:"image_preview_url"`     // NFT作品原图
	ImageThumbnailURL string `json:"image_thumbnail_url" bson:"image_thumbnail_url"` // NFT作品缩略图
	Description       string `json:"description" bson:"description"`                 // NFT作品描述
	ContractAddress   string `json:"contract_address" bson:"contract_address"`       // 合约地址
	TokenId           string `json:"token_id" bson:"token_id"`                       // NFT token id
	NumSales          int    `json:"num_sales" bson:"num_sales"`                     // NFT售卖次数
	Owner             string `json:"owner" bson:"owner"`                             // NFT拥有者
	OwnerImgURL       string `json:"owner_img_url" bson:"owner_img_url"`             // 拥有者头像
	Creator           string `json:"creator" bson:"creator"`                         // NFT创造者
	CreatorImgURL     string `json:"creator_img_url" bson:"creator_img_url"`         // 创造者头像
	TokenMetadata     string `json:"token_metadata" bson:"token_metadata"`           // NFT元数据

	Slug string `json:"slug" bson:"slug"` // 集合唯一标识符号

	//Contract            Contract             `json:"contract" bson:"contract"`
	//Collection          Collection           `json:"collection" bson:"collection"`
	//AssetsTopOwnerships []AssetsTopOwnership `json:"assets_top_ownership" bson:"assets_top_ownership,omitempty"`
	//Traits              []Trait              `json:"trait" bson:"trait,omitempty"`

	IsDelete    int8 `json:"is_delete" bson:"is_delete"`       // 是否删除 1删除 0未删除 默认为0
	RefreshTime int  `json:"refresh_time" bson:"refresh_time"` // 刷新时间
}

type Contract struct {
	//ID           int64  `json:"id" bson:"id"`                       // 主键
	Address      string `json:"address" bson:"address"`             // 合约地址
	ContractType string `json:"contract_type" bson:"contract_type"` // 合约类型 semi-fungible可替代 non-fungible 不可替代
	ContractName string `json:"contract_name" bson:"contract_name"` // 合约名字
	Symbol       string `json:"symbol" bson:"symbol"`               // 符号
	SchemaName   string `json:"schema_name" bson:"schema_name"`     // 合约类型
	TotalSupply  string `json:"total_supply" bson:"total_supply"`   // 总供应量
	Description  string `json:"description" bson:"description"`     // 合约描述
}

type Trait struct {
	//ID              int64  `json:"id" bson:"id"`                     // 主键
	UserAddress     string `json:"user_address" bson:"user_address"` // 用户地址
	ContractAddress string `json:"_" bson:"contract_address"`        // 合约地址
	TokenId         string `json:"_" bson:"token_id"`                // token id
	TraitType       string `json:"trait_type" bson:"trait_type"`     // 特征类型
	Value           string `json:"value" bson:"value"`               // 特征值
	DisplayType     string `json:"display_type" bson:"display_type"`
	MaxValue        int    `json:"max_value" bson:"max_value"`
	TraitCount      int    `json:"trait_count" bson:"trait_count"` // 数量
	OrderBy         string `json:"order_by" bson:"order_by"`
	IsDelete        int8   `json:"is_delete" bson:"is_delete"`       // 是否删除 1删除 0未删除 默认为0
	RefreshTime     int    `json:"refresh_time" bson:"refresh_time"` // 刷新时间
}

type AssetsTopOwnership struct {
	//ID              int64  `json:"id" bson:"id"`                           // 主键
	UserAddress     string `json:"user_address" bson:"user_address"`       // 用户地址
	ContractAddress string `json:"_" bson:"contract_address"`              // 合约地址
	TokenId         string `json:"_" bson:"token_id"`                      // token id
	Owner           string `json:"owner" bson:"owner"`                     // 所有者地址
	ProfileImgURL   string `json:"profile_img_url" bson:"profile_img_url"` // 所有者头像
	Quantity        string `json:"quantity" bson:"quantity"`               // 数量
	IsDelete        int8   `json:"is_delete" bson:"is_delete"`             // 是否删除 1删除 0未删除 默认为0
	RefreshTime     int    `json:"refresh_time" bson:"refresh_time"`       // 刷新时间
}

// InsertOpenSeaAsset query Aseets through opensea API and insert
func InsertOpenSeaAsset(assets *OwnerAsset, user string) error {
	db := database.GetMongoClient()
	refreshTime := int(time.Now().Unix())

	// No blocking query opensea assets_top_ownerships
	go queryAssetsTopOwnerShip(db, assets, refreshTime, user)

	for _, v := range assets.Assets {
		owner := user
		if v.Owner.Address != "0x0000000000000000000000000000000000000000" {
			owner = v.Owner.Address
		}

		var asset = Asset{
			UserAddress:       user,
			Title:             v.Name,
			ImageURL:          v.ImageURL,
			ImagePreviewURL:   v.ImagePreviewURL,
			ImageThumbnailURL: v.ImageThumbnailURL,
			Description:       v.Description,
			ContractAddress:   v.AssetContract.Address,
			TokenId:           v.TokenID,
			NumSales:          v.NumSales,
			Owner:             owner,
			OwnerImgURL:       v.Owner.ProfileImgURL,
			Creator:           v.Creator.Address,
			CreatorImgURL:     v.Creator.ProfileImgURL,
			Slug:              v.Collection.Slug,
			TokenMetadata:     v.TokenMetadata,
			RefreshTime:       refreshTime,
		}
		var contract = Contract{
			Address:      v.AssetContract.Address,
			ContractName: v.AssetContract.Name,
			ContractType: v.AssetContract.AssetContractType,
			Symbol:       v.AssetContract.Symbol,
			SchemaName:   v.AssetContract.SchemaName,
			TotalSupply:  v.AssetContract.TotalSupply,
			Description:  v.AssetContract.Description,
		}

		count, err := db.Collection("assets").
			CountDocuments(context.TODO(), bson.M{"user_address": user, "contract_address": v.AssetContract.Address, "token_id": v.TokenID, "is_delete": 0})
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		if count == 0 {
			if _, err = db.Collection("assets").InsertOne(context.TODO(), &asset); err != nil {
				logs.GetLogger().Error(err)
				return err
			}
		} else {
			// Refresh synchronization time
			if _, err = db.Collection("assets").UpdateOne(
				context.TODO(),
				bson.M{"user_address": user, "contract_address": v.AssetContract.Address, "token_id": v.TokenID, "is_delete": 0},
				bson.M{"$set": bson.M{"refresh_time": refreshTime}}); err != nil {
				logs.GetLogger().Error(err)
				return err
			}

			// update NFT without locked metadata
			if v.TokenMetadata == "" {
				if _, err = db.Collection("assets").UpdateOne(
					context.TODO(),
					bson.M{"user_address": user, "contract_address": v.AssetContract.Address, "token_id": v.TokenID, "is_delete": 0},
					bson.M{"$set": bson.M{"title": v.Name, "image_url": v.ImageURL, "image_preview_url": v.ImagePreviewURL, "image_thumbnail_url": v.ImageThumbnailURL,
						"description": v.Description, "refresh_time": refreshTime}}); err != nil {
					logs.GetLogger().Error(err)
					return err
				}
			}

		}

		// insert contract
		count, err = db.Collection("contracts").
			CountDocuments(context.TODO(), bson.M{"address": v.AssetContract.Address})
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		if count == 0 {
			if _, err = db.Collection("contracts").InsertOne(context.TODO(), &contract); err != nil {
				logs.GetLogger().Error(err)
			}
		}

		// update ro insert traits
		for _, v1 := range v.Traits {
			var trait = Trait{
				UserAddress:     user,
				ContractAddress: v.AssetContract.Address,
				TokenId:         v.TokenID,
				TraitType:       v1.TraitType,
				Value:           v1.Value,
				DisplayType:     v1.DisplayType,
				MaxValue:        v1.MaxValue,
				TraitCount:      v1.TraitCount,
				OrderBy:         v1.Order,
				RefreshTime:     refreshTime,
			}

			count, err = db.Collection("traits").CountDocuments(
				context.TODO(),
				bson.M{"user_address": user, "contract_address": v.AssetContract.Address, "token_id": v.TokenID,
					"trait_type": v1.TraitType, "value": v1.Value, "is_delete": 0})
			if err != nil {
				logs.GetLogger().Error(err)
				return err
			}
			if count == 0 {
				if _, err = db.Collection("traits").InsertOne(context.TODO(), &trait); err != nil {
					logs.GetLogger().Error(err)
				}
			} else {
				if _, err = db.Collection("traits").UpdateOne(
					context.TODO(),
					bson.M{"user_address": user, "contract_address": v.AssetContract.Address, "token_id": v.TokenID,
						"trait_type": v1.TraitType, "value": v1.Value, "is_delete": 0},
					bson.M{"$set": bson.M{"refresh_time": refreshTime}}); err != nil {
					logs.GetLogger().Error(err)
					return err
				}
			}

		}
	}

	// Delete opensea deleted asset
	if _, err := db.Collection("assets").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "refresh_time": bson.M{"$lt": refreshTime}, "is_delete": 0},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	// Delete opensea deleted traits
	if _, err := db.Collection("traits").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "refresh_time": bson.M{"$lt": refreshTime}, "is_delete": 0},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil

}

// FindAssetByOwner find assets by owner
func FindAssetByOwner(user string) ([]map[string]interface{}, error) {
	var assets []*Asset
	var assetsList []map[string]interface{}
	db := database.GetMongoClient()
	cursor, err := db.Collection("assets").Find(context.TODO(), bson.M{"user_address": user, "is_delete": 0})
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &assets); err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	for _, v := range assets {
		var collection Collection
		var contract Contract
		var assetsTopOwnerships []AssetsTopOwnership
		var traits []Trait
		var tmp map[string]interface{}
		err = db.Collection("collections").FindOne(context.TODO(), bson.M{"user_address": user, "slug": v.Slug}).Decode(&collection)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}

		err = db.Collection("contracts").FindOne(context.TODO(), bson.M{"address": v.ContractAddress}).Decode(&contract)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}

		cursor, err = db.Collection("assets_top_ownerships").
			Find(context.TODO(), bson.M{"user_address": user, "contract_address": v.ContractAddress, "token_id": v.TokenId})
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		if err = cursor.All(context.TODO(), &assetsTopOwnerships); err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}

		cursor, err = db.Collection("traits").
			Find(context.TODO(), bson.M{"user_address": user, "contract_address": v.ContractAddress, "token_id": v.TokenId})
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		if err = cursor.All(context.TODO(), &traits); err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		bytes, err := json.Marshal(v)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		if err = json.Unmarshal(bytes, &tmp); err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		tmp["contract"] = contract
		tmp["collection"] = collection
		tmp["assets_top_ownership"] = assetsTopOwnerships
		tmp["trait"] = traits
		assetsList = append(assetsList, tmp)
	}

	return assetsList, nil
}

// FindWorksBySlug find assets by collection
func FindWorksBySlug(user, slug string) ([]*Asset, error) {
	var assets []*Asset
	db := database.GetMongoClient()

	cursor, err := db.Collection("assets").Find(context.TODO(), bson.M{"user_address": user, "slug": slug, "is_delete": 0})
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &assets); err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return assets, nil
}

// DeleteAssetByTokenID delete asset by tokenId
func DeleteAssetByTokenID(user, contractAddress, tokenID string) error {
	db := database.GetMongoClient()
	if _, err := db.Collection("assets").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "contract_address": contractAddress, "token_id": tokenID},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	if _, err := db.Collection("traits").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "contract_address": contractAddress, "token_id": tokenID},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	if _, err := db.Collection("assets_top_ownerships").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "contract_address": contractAddress, "token_id": tokenID},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	return nil
}

func queryAssetsTopOwnerShip(db *mongo.Database, assets *OwnerAsset, refreshTime int, user string) error {
	for _, v := range assets.Assets {
		time.Sleep(time.Second)
		// If the number of requests is too many, a 429 error code will be thrown
		resp, err := utils.RequestOpenSeaSingleAsset(v.AssetContract.Address, v.TokenID)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		var autoAsset AutoAsset
		if err = json.Unmarshal(resp, &autoAsset); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		for _, a := range autoAsset.TopOwnerships {
			var assetsTopOwnership = AssetsTopOwnership{
				UserAddress:     user,
				ContractAddress: v.AssetContract.Address,
				TokenId:         v.TokenID,
				Owner:           a.Owner.Address,
				ProfileImgURL:   a.Owner.ProfileImgURL,
				Quantity:        a.Quantity,
				RefreshTime:     refreshTime,
			}

			count, err := db.Collection("assets_top_ownerships").
				CountDocuments(context.TODO(), bson.M{"user_address": user, "contract_address": v.AssetContract.Address, "token_id": v.TokenID, "is_delete": 0})
			if err != nil {
				logs.GetLogger().Error(err)
				return err
			}
			if count == 0 {
				if _, err = db.Collection("assets_top_ownerships").InsertOne(context.TODO(), &assetsTopOwnership); err != nil {
					logs.GetLogger().Error(err)
				}
			} else {
				if _, err = db.Collection("assets_top_ownerships").UpdateOne(
					context.TODO(),
					bson.M{"user_address": user, "contract_address": v.AssetContract.Address, "token_id": v.TokenID, "is_delete": 0},
					bson.M{"$set": bson.M{"refresh_time": refreshTime}}); err != nil {
					logs.GetLogger().Error(err)
					return err
				}
			}
		}

	}
	// Delete opensea deleted traits
	if _, err := db.Collection("assets_top_ownerships").UpdateMany(
		context.TODO(),
		bson.M{"user_address": user, "refresh_time": bson.M{"$lt": refreshTime}},
		bson.M{"$set": bson.M{"is_delete": 1}}); err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}
