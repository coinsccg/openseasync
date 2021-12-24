package models

import (
	"openseasync/database"
)

type Asset struct {
	Id                uint64 `json:"id"`                  // NFT作品唯一ID
	Title             string `json:"title"`               // NFT作品标题
	ImageURL          string `json:"image_url"`           // NFT作品图片
	ImagePreviewURL   string `json:"image_preview_url"`   // NFT作品原图
	ImageThumbnailURL string `json:"image_thumbnail_url"` // NFT作品缩略图
	Description       string `json:"description"`         // NFT作品描述

	Address       string `json:"address"`         // 合约地址
	TokenId       string `json:"token_id"`        // NFT token id
	ContractName  string `json:"contract_name"`   // 合约名字
	Symbol        string `json:"symbol"`          // 符号
	SchemaName    string `json:"schema_name"`     // 合约类型
	TotalSupply   string `json:"total_supply"`    // 总供应量
	Owner         string `json:"owner"`           // NFT拥有者
	OwnerImgURL   string `json:"owner_img_url"`   // 拥有者头像
	Creator       string `json:"creator"`         // NFT创造者
	CreatorImgURL string `json:"creator_img_url"` // 创造者头像
	CreateDate    string `json:"create_date"`     // 合约创建时间

	Slug                   string `json:"slug"`                      // 集合唯一标识符号
	CollectionImgURL       string `json:"collection_img_url"`        // 集合中的图片
	CollectionBannerImgURL string `json:"collection_banner_img_url"` // 集合背景图片
	CollectionDescription  string `json:"collection_description"`    // 集合描述
	CollectionLargeImgURL  string `json:"collection_large_img_url"`  // 集合头像
}

// InsertOpenSeaAsset query Aseets through opensea API and insert
func InsertOpenSeaAsset(assets *OwnerAsset, owner string) error {
	db := database.GetDB()

	for _, v := range assets.Assets {
		var asset = Asset{
			Title:                  v.Name,
			ImageURL:               v.ImageURL,
			ImagePreviewURL:        v.ImagePreviewURL,
			ImageThumbnailURL:      v.ImageThumbnailURL,
			Description:            v.Description,
			Address:                v.AssetContract.Address,
			TokenId:                v.TokenID,
			ContractName:           v.AssetContract.Name,
			Symbol:                 v.AssetContract.Symbol,
			SchemaName:             v.AssetContract.SchemaName,
			TotalSupply:            v.AssetContract.TotalSupply,
			Owner:                  owner,
			OwnerImgURL:            v.Owner.ProfileImgURL,
			Creator:                v.Creator.Address,
			CreatorImgURL:          v.Creator.ProfileImgURL,
			CreateDate:             v.AssetContract.CreatedDate,
			Slug:                   v.Collection.Slug,
			CollectionImgURL:       v.Collection.ImageURL,
			CollectionBannerImgURL: v.Collection.BannerImageURL,
			CollectionDescription:  v.Collection.Description,
			CollectionLargeImgURL:  v.Collection.LargeImageURL,
		}

		// gorm v1  batch insert is not supported
		var tmp Asset
		rows := db.Table("assets").
			Where("address = ? AND token_id = ? AND schema_name = ?", v.AssetContract.Address, v.TokenID, v.AssetContract.SchemaName).
			Find(&tmp).RowsAffected
		if rows == 0 {
			if err := db.Table("assets").Create(&asset).Error; err != nil {
				return err
			}
		}

	}
	return nil

}

// FindAssetByOwner find assets by owner
func FindAssetByOwner(owner string) ([]*Asset, error) {
	var assets []*Asset
	db := database.GetDB()

	if err := db.Table("assets").Where("owner = ?", owner).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

// FindWorksBySlug find assets by collection
func FindWorksBySlug(owner, slug string) ([]*Asset, error) {
	var assets []*Asset
	db := database.GetDB()

	if err := db.Table("assets").Where("owner = ? AND slug = ?", owner, slug).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}
