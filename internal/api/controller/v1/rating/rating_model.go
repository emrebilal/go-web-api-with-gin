package rating

type AddRatingModel struct {
	UserName   string `json:"UserName"`
	ProviderId string `json:"ProviderId"`
	ServiceId  string `json:"ServiceId"`
	Rate       int    `json:"Rate"`
}
