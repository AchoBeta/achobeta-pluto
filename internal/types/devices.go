package types

type Device struct {
	LoginId    string `json:"login_id"`
	Ip         string `json:"ip"`
	OnlineTime string `json:"online_time"`
	DeviceName string `json:"device_name"`
}
type DevicesReq struct {
	PageNumber int `json:"page_number"`
	LineNumber int `json:"line_number"`
}
type DevicesResp struct {
	Token   string   `json:"token"`
	Devices []Device `json:"devices"`
}
