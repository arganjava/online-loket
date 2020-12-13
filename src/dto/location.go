package dto

type LocationRequest struct {
	Country  string `uri:"country" binding:"required" validate:"min=1,max=30"`
	CityName string `uri:"cityName" binding:"required" validate:"min=1,max=50"`
	Village  string `uri:"village" binding:"required" validate:"min=1,max=50"`
	Address  string `uri:"address" binding:"required" validate:"min=1,max=200"`
}

type LocationResponse struct {
	ID       string `json:"id"`
	Country  string `json:"country"`
	CityName string `json:"cityName"`
	Village  string `json:"village"`
	Address  string `json:"address"`
}
