package models

type Asset struct {
	Id                uint64 `json:"id"`
	TokenId           string `json:"token_id"`
	ImageURL          string `json:"image_url"`
	ImagePreviewURL   string `json:"image_preview_url"`
	ImageThumbnailURL string `json:"image_thumbnail_url"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	AssetContract     `json:"asset_contract"`
}

type AssetContract struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	SchemaName  string `json:"schema_name"`
	TotalSupply string `json:"total_supply"`
	Payout      string `json:"payout"`
	CreateDate  string `json:"create_date"`
}

type Collection struct {
	BannerImageURL string `json:"banner_image_url"`
	Description    string `json:"description"`
	ImageURL       string `json:"image_url"`
	LargeImageURL  string `json:"large_image_url"`
	Name           string `json:"name"`
}

type SellOrder struct {
	OrderHash string `json:"order_hash"`
	Exchange  string `json:"exchange"`
}

type Creator struct {
	Address  string `json:"address"`
	Username string `json:"username"`
}
