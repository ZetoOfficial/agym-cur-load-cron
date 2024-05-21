package models

type ClubListResponse struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
	Address  string `json:"address"`
}

type ClubInfoResponse struct {
	Id                 int      `json:"id"`
	Title              string   `json:"title"`
	City               string   `json:"city"`
	CoordinatesLatLong []string `json:"coordinatesLatLong"`
	CurrentLoad        string   `json:"currentLoad"`
}
