package rating

type AddRatingModel struct {
	UserName   string `validate:"required"`
	ProviderId string `validate:"required"`
	ServiceId  string `validate:"required"`
	Rate       int    `validate:"required,gte=1,lte=5"`
}

type GetAllRatingsModel struct {
	ProviderId string `validate:"required"`
}
