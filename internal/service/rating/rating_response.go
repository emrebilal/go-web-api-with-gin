package rating

type SendRatingServiceResponse struct {
	Error error `json:"-"`
	Info  string
}

type GetAverageRatingServiceResponse struct {
	Error         error `json:"-"`
	AverageRating AverageRatingModel
}

type AverageRatingModel struct {
	ProviderId  string
	AverageRate float64
}
