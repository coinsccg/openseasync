package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"openseasync/common"
	"openseasync/common/constants"
	"openseasync/common/errorinfo"
	"openseasync/logs"
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

// getOpenSeaOwnerAssets get all assets
func getOpenSeaOwnerAssets(owner string, offset, limit int64) (map[string]interface{}, string, string) {

	url := fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", constants.OPENSEA_ASSETS_URL, owner, offset, limit)
	fmt.Println(url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		logs.GetLogger().Error(err)
		return nil, errorinfo.OPENSEA_HTTP_REQUEST_ERROR_CODE, err.Error()
	}
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	var assets map[string]interface{}
	json.Unmarshal(content, &assets)
	return assets, "", ""
}

// getOpenSeaSingleAsset get single asset
func getOpenSeaSingleAsset(contractAddress string, tokenId int64) (map[string]interface{}, string, string) {

	url := fmt.Sprintf("%s/%s/%d/", constants.OPENSEA_SINGLE_ASSET_URL, contractAddress, tokenId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, errorinfo.OPENSEA_HTTP_REQUEST_ERROR_CODE, err.Error()
	}
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	var assets map[string]interface{}
	json.Unmarshal(content, &assets)
	return assets, "", ""
}
