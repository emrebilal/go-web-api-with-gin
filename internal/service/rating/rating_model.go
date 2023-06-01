package rating

type SendRatingServiceModel struct {
	UserName   string `validate:"required"`
	ProviderId string `validate:"required"`
	ServiceId  string `validate:"required"`
	Rate       int    `validate:"required,gte=1,lte=5"`
}

type GetAverageRatingServiceModel struct {
	ProviderId string `validate:"required"`
}
