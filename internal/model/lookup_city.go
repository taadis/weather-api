package model

// LookupCityRequest 城市信息查询请求
type LookupCityRequest struct {
	Location string `json:"location,omitempty"`
	Adm      string `json:"adm,omitempty"`
}

// LookupCityResponse 城市信息查询响应
type LookupCityResponse map[string]interface{}
