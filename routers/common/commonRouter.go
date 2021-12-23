package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"openseasync/common"
	"openseasync/common/constants"
	"openseasync/common/errorinfo"
	"strconv"
)

func HostManager(router *gin.RouterGroup) {
	router.GET(constants.URL_HOST_GET_HOST_INFO, GetSwanMinerVersion)
	router.GET(constants.URL_OPENSEA_OWNER_ASSETS, GetOpenSeaOwnerAssets)
	router.GET(constants.URL_OPENSEA_SINGLE_ASSETS, GetOpenSeaSingleAsset)

}

func GetSwanMinerVersion(c *gin.Context) {
	info := getSwanMinerHostInfo()
	c.JSON(http.StatusOK, common.CreateSuccessResponse(info))
}

func GetOpenSeaOwnerAssets(c *gin.Context) {
	var param struct {
		Owner  string `form:"owner"`
		Offset int64  `form:"offset"`
		Limit  int64  `form:"limit"`
	}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(errorinfo.HTTP_REQUEST_PARSER_STRUCT_TO_REQUEST_ERROR_CODE, err.Error()))
	}

	assets, errCode, errMsg := getOpenSeaOwnerAssets(param.Owner, param.Offset, param.Limit)
	if errMsg != "" {
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(errCode, errMsg))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(assets))
}

func GetOpenSeaSingleAsset(c *gin.Context) {
	contractAddress := c.Param("contract_address")
	tokenId, err := strconv.ParseInt(c.Param("token_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(errorinfo.HTTP_REQUEST_PARAM_TYPE_ERROR_CODE, err.Error()))
	}
	asset, errCode, errMsg := getOpenSeaSingleAsset(contractAddress, tokenId)
	if errCode != "" {
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(errCode, errMsg))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(asset))
}
