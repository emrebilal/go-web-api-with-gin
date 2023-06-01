package rating

type AddRatingResponse struct {
	Error error `json:"-"`
}

type GetAllRatingsResponse struct {
	Error error `json:"-"`
	Rates []int
}
