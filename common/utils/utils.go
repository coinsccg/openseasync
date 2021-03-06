package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	constants2 "openseasync/common/constants"
	"openseasync/config"
	"openseasync/logs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filswan/go-swan-lib/constants"
)

// GetEpochInMillis get current timestamp
func GetEpochInMillis() (millis int64) {
	nanos := time.Now().UnixNano()
	millis = nanos / 1000000
	return
}

func ReadContractAbiJsonFile(aptpath string) (string, error) {
	jsonFile, err := os.Open(aptpath)

	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}
	return string(byteValue), nil
}

func GetRewardPerBlock() *big.Int {
	rewardBig, _ := new(big.Int).SetString("35000000000000000000", 10) // the unit is wei
	return rewardBig
}

func CheckTx(client *ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
retry:
	rp, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		if err == ethereum.NotFound {
			logs.GetLogger().Error("tx %v not found, check it later", tx.Hash().String())
			time.Sleep(1 * time.Second)
			goto retry
		} else {
			logs.GetLogger().Error("TransactionReceipt fail: %s", err)
			return nil, err
		}
	}
	return rp, nil
}

func GetFromAndToAddressByTxHash(client *ethclient.Client, chainID *big.Int, txHash common.Hash) (*addressInfo, error) {
	addrInfo := new(addressInfo)
	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	addrInfo.AddrTo = tx.To().Hex()
	txMsg, err := tx.AsMessage(types.NewEIP155Signer(chainID), nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	addrInfo.AddrFrom = txMsg.From().Hex()
	return addrInfo, nil
}

type addressInfo struct {
	AddrFrom string
	AddrTo   string
}

func GetOffsetByPagenumber(pageNumber, pageSize string) (int64, error) {
	pageNumberInt, err := strconv.ParseInt(pageNumber, 10, 64)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}
	pageSizeInt, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}
	offset := (pageNumberInt - 1) * pageSizeInt
	return offset, nil
}

func HttpUploadFileByStream(uri, filefullpath string) (string, error) {
	fileReader, err := os.Open(filefullpath)
	if err != nil {
		logs.GetLogger().Error(err)
		return constants.EMPTY_STRING, err
	}

	filename := filepath.Base(filefullpath)

	boundary := "MyMultiPartBoundary12345"
	token := "DEPLOY_GATE_TOKEN"
	message := "Uploaded by Nebula"
	releaseNote := "Built by Nebula"
	fieldFormat := "--%s\r\nContent-Disposition: form-data; name=\"%s\"\r\n\r\n%s\r\n"
	tokenPart := fmt.Sprintf(fieldFormat, boundary, "token", token)
	messagePart := fmt.Sprintf(fieldFormat, boundary, "message", message)
	releaseNotePart := fmt.Sprintf(fieldFormat, boundary, "release_note", releaseNote)
	fileName := filename
	fileHeader := "Content-type: application/octet-stream"
	fileFormat := "--%s\r\nContent-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n%s\r\n\r\n"
	filePart := fmt.Sprintf(fileFormat, boundary, fileName, fileHeader)
	bodyTop := fmt.Sprintf("%s%s%s%s", tokenPart, messagePart, releaseNotePart, filePart)
	bodyBottom := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	body := io.MultiReader(strings.NewReader(bodyTop), fileReader, strings.NewReader(bodyBottom))

	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)

	response, err := http.Post(uri, contentType, body)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", nil
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status:%s, code:%d, url:%s", response.Status, response.StatusCode, uri)
		logs.GetLogger().Error(err)
		switch response.StatusCode {
		case http.StatusNotFound:
			logs.GetLogger().Error("please check your url:", uri)
		}
		return constants.EMPTY_STRING, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.GetLogger().Error(err)
		return constants.EMPTY_STRING, err
	}

	responseStr := string(responseBody)
	//logs.GetLogger().Info(responseStr)
	filesInfo := strings.Split(responseStr, "\n")
	if len(filesInfo) < 4 {
		err := fmt.Errorf("not enough files infor returned")
		logs.GetLogger().Error(err)
		return constants.EMPTY_STRING, err
	}
	responseStr = filesInfo[3]
	return responseStr, nil
}

func UrlJoin(root string, parts ...string) string {
	url := root

	for _, part := range parts {
		url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(part, "/")
	}
	url = strings.TrimRight(url, "/")

	return url
}

func GetBoolPointer(b bool) *bool {
	return &b
}

func MD5(v string) string {
	h := md5.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}

func RequestOpenSeaAssets(owner string, offset, limit int64) ([]byte, error) {
	openSeaAssetsURL := constants2.OPENSEA_DEV_ASSETS_URL
	if !config.GetConfig().Dev {
		openSeaAssetsURL = constants2.OPENSEA_PROD_ASSETS_URL
	}

	url := fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", openSeaAssetsURL, owner, offset, limit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if !config.GetConfig().Dev {
		req.Header.Add("X-API-KEY", constants2.OPENSEA_API_KEY)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(content))
	}
	return content, nil
}

func RequestOpenSeaCollections(owner string, offset, limit int64) ([]byte, error) {
	openSeaCollectionsURL := constants2.OPENSEA_DEV_COLLECTION_URL
	if !config.GetConfig().Dev {
		openSeaCollectionsURL = constants2.OPENSEA_PROD_COLLECTIONS_URL
	}
	url := fmt.Sprintf("%s?asset_owner=%s&offset=%d&limit=%d", openSeaCollectionsURL, owner, offset, limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if !config.GetConfig().Dev {
		req.Header.Add("X-API-KEY", constants2.OPENSEA_API_KEY)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
		fmt.Println("+++++++++++++++++++++ Collections Begin ++++++++++++++++++++++++++")

	fmt.Println(string(content))
	fmt.Println("+++++++++++++++++++++ Collections End ++++++++++++++++++++++++++")
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(string(content))
		return nil, err
	}
	return content, nil
}

func RequestOpenSeaSingleAsset(contractAddress, tokenId string) ([]byte, error) {
	openSeaSingleAssetURL := constants2.OPENSEA_DEV_SINGLE_ASSET_URL
	if !config.GetConfig().Dev {
		openSeaSingleAssetURL = constants2.OPENSEA_PROD_SINGLE_ASSET_URL
	}
	url := fmt.Sprintf("%s/%s/%s", openSeaSingleAssetURL, contractAddress, tokenId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if !config.GetConfig().Dev {
		req.Header.Add("X-API-KEY", constants2.OPENSEA_API_KEY)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(string(content))
		return nil, err
	}
	return content, nil
}

func RequestOpenSeaEvent(contractAddress, tokenId string) ([]byte, error) {
	openSeaEventURL := constants2.OPENSEA_DEV_EVENT_URL
	if !config.GetConfig().Dev {
		openSeaEventURL = constants2.OPENSEA_PROD_EVENT_URL
	}
	url := fmt.Sprintf("%s?asset_contract_address=%s&token_id=%s", openSeaEventURL, contractAddress, tokenId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	if !config.GetConfig().Dev {
		req.Header.Add("X-API-KEY", constants2.OPENSEA_API_KEY)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		err = errors.New(string(content))
		return nil, err
	}
	return content, nil
}

func ParseTime(timeStr string) int64 {
	str := strings.Split(timeStr, ".")[0]
	str = strings.ReplaceAll(str, "T", " ")
	i, _ := time.Parse("2006-01-02 15:04:05", str)
	i = i.Add(-time.Hour * 8)
	tm := i.UnixMilli()
	return tm
}

func EthToWei(value float64) *big.Int {
	val := big.NewFloat(value)
	base := new(big.Float)
	base.SetInt(big.NewInt(int64(math.Pow10(18))))

	val.Mul(val, base)
	result := new(big.Int)
	val.Int(result)
	

	return result
}
