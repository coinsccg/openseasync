package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"openseasync/common"
	"openseasync/common/constants"
	"openseasync/common/errorinfo"
)

func HostManager(router *gin.RouterGroup) {
	router.GET(constants.URL_HOST_GET_HOST_INFO, GetSwanMinerVersion)
	router.GET(constants.URL_OPENSEA_OWNER_ASSETS, GetOpenSeaOwnerAssets)
	router.GET(constants.URL_OPENSEA_OWNER_Collections, GetOpenSeaOwnerCollections)
	router.GET(constants.URL_FIND_ASSETS_OWNER, GetAssetsByOwner)
	router.GET(constants.URL_FIND_Collections_OWNER, GetCollectionsByOwner)
	router.GET(constants.URL_FIND_ASSETS_SLUG, GetAssetsBySlug)

}

func GetSwanMinerVersion(c *gin.Context) {
	info := getSwanMinerHostInfo()
	c.JSON(http.StatusOK, common.CreateSuccessResponse(info))
}

func GetOpenSeaOwnerAssets(c *gin.Context) {
	owner := c.Param("owner")
	if len(owner) != 42 {
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_CODE, errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_MSG))
		return
	}
	if err := getOpenSeaOwnerAssets(owner); err != nil {
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(errorinfo.OPENSEA_HTTP_REQUEST_ERROR_CODE, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(nil))
}

func GetOpenSeaOwnerCollections(c *gin.Context) {
	owner := c.Param("owner")
	if len(owner) != 42 {
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_CODE, errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_MSG))
		return
	}
	if err := getOpenSeaOwnerCollection(owner); err != nil {
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(errorinfo.OPENSEA_HTTP_REQUEST_ERROR_CODE, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(nil))
}

func GetAssetsByOwner(c *gin.Context) {
	owner := c.Param("owner")
	if len(owner) != 42 {
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_CODE, errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_MSG))
		return
	}
	assets, err := getAssetByOwner(owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(errorinfo.GET_RECORD_lIST_ERROR_CODE, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(assets))
}

func GetCollectionsByOwner(c *gin.Context) {
	owner := c.Param("owner")
	if len(owner) != 42 {
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_CODE, errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_MSG))
		return
	}
	collections, err := getCollectionsByOwner(owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(errorinfo.GET_RECORD_lIST_ERROR_CODE, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(collections))
}

func GetAssetsBySlug(c *gin.Context) {
	owner := c.Param("owner")
	slug := c.Param("slug")
	if len(owner) != 42 && slug != "" {
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_CODE, errorinfo.HTTP_REQUEST_PARAM_VALUE_ERROR_MSG))
		return
	}
	assets, err := getAssetBySlug(owner, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(errorinfo.GET_RECORD_lIST_ERROR_CODE, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(assets))
}
