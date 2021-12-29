package common

import (
	"encoding/json"
	"openseasync/common"
	"openseasync/common/utils"
	"openseasync/logs"
	"openseasync/models"
	"runtime"
	"time"
)

func getSwanMinerHostInfo() *common.HostInfo {
	info := new(common.HostInfo)
	info.SwanMinerVersion = common.GetVersion()
	info.OperatingSystem = runtime.GOOS
	info.Architecture = runtime.GOARCH
	info.CPUnNumber = runtime.NumCPU()
	return info
}

// getOpenSeaOwnerAssets get all assets by owner
func getOpenSeaOwnerAssets(owner string) error {
	var n int64 = 1
	for {
		// If the number of requests is too many, a 429 error code will be thrown
		content, err := utils.RequestOpenSeaAssets(owner, 50*(n-1), 50)
		if err != nil {
			return err
		}
		var assets models.OwnerAsset
		if err = json.Unmarshal(content, &assets); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		if len(assets.Assets) < 1 {
			break
		}
		if err = models.InsertOpenSeaAsset(&assets, owner); err != nil {
			return err
		}
		if len(assets.Assets) < 50 {
			break
		}
		n++
		time.Sleep(time.Second / time.Duration(2))
	}

	return nil
}

// getOpenSeaOwnerCollection get all collections by owner
func getOpenSeaOwnerCollection(owner string) error {
	var n int64 = 1
	for {
		content, err := utils.RequestOpenSeaCollections(owner, 300*(n-1), 300*n)
		if err != nil {
			return err
		}
		var collections models.OwnerCollection
		if err = json.Unmarshal(content, &collections.Collections); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		if err = models.InsertOpenSeaCollection(&collections, owner); err != nil {
			return err
		}
		if len(collections.Collections) < 300 {
			break
		}
		n++
	}

	return nil
}

// getAssetByOwner get assets by owner
func getAssetByOwner(owner string) ([]*models.Asset, error) {
	assets, err := models.FindAssetByOwner(owner)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// getAssetBySlug get assets by owner
func getAssetBySlug(owner, slug string) ([]*models.Asset, error) {
	assets, err := models.FindWorksBySlug(owner, slug)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// getCollectionsByOwner get collection by owner
func getCollectionsByOwner(owner string) ([]*models.Collection, error) {
	collections, err := models.FindCollectionByOwner(owner)
	if err != nil {
		return nil, err
	}
	return collections, nil
}
