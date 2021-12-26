package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"openseasync/common"
	"openseasync/common/constants"
	"openseasync/logs"
	"openseasync/models"
	"runtime"
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
		content, err := requestOpenSeaAssets(owner, 50*(n-1), 50*n)
		var assets models.OwnerAsset
		if err = json.Unmarshal(content, &assets); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		if err = models.InsertOpenSeaAsset(&assets, owner); err != nil {
			return err
		}
		if len(assets.Assets) < 50 {
			break
		}
		n++
	}

	return nil
}

// getOpenSeaOwnerCollection get all collections by owner
func getOpenSeaOwnerCollection(owner string) error {
	var n int64 = 1
	for {
		content, err := requestOpenSeaCollections(owner, 300*(n-1), 300*n)
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
		logs.GetLogger().Error(err)
		return nil, err
	}
	return assets, nil
}

// getAssetBySlug get assets by owner
func getAssetBySlug(owner, slug string) ([]*models.Asset, error) {
	assets, err := models.FindWorksBySlug(owner, slug)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return assets, nil
}

// getCollectionsByOwner get collection by owner
func getCollectionsByOwner(owner string) ([]*models.Collection, error) {
	collections, err := models.FindCollectionByOwner(owner)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return collections, nil
}

func requestOpenSeaAssets(owner string, offset, limit int64) ([]byte, error) {
	url := fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", constants.OPENSEA_ASSETS_URL, owner, offset, limit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	return content, nil
}

func requestOpenSeaCollections(owner string, offset, limit int64) ([]byte, error) {
	url := fmt.Sprintf("%s?asset_owner=%s&offset=%d&limit=%d", constants.OPENSEA_COLLECTION_URL, owner, offset, limit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	return content, nil
}
