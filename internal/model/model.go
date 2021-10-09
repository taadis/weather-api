package model

type baseResponse struct {
	Code string `json:"code,omitempty"` // 状态码
	//Refer struct {
	//	Sources []string `json:"sources,omitempty"` // 原始数据来源，或数据源说明，可能为空
	//	License []string `json:"license,omitempty"` // 数据许可或版权声明，可能为空
	//} `json:"refer,omitempty"` //
}
