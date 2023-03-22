package maps

type CalculateListRequest struct {
	Pairs []LocationPair `json:"pairs"`
	RoundOff bool `json:"roundOff"`
}
type CalculateListResponse struct {
	TravelTimeMinutes []int `json:"travelTimeMinutes"`
}
type CalculateRequest struct {
	Source Location `json:"source"`
	Destination Location `json:"destination"`
	RoundOff bool `json:"roundOff"`
}
type CalculateResponse struct {
	TravelTimeMinutes int `json:"travelTimeMinutes"`
}
type Country string
const (
	Germany Country = "Germany"
	France Country = "France"
	Austria Country = "Austria"
)
type DateTime string
type Location struct {
	PostalCode string `json:"postalCode"`
	Street string `json:"street,omitempty"`
	City string `json:"city,omitempty"`
	Country Country `json:"country"`
}
type LocationPair struct {
	Source Location `json:"source"`
	Destination Location `json:"destination"`
}
type CalculateTimeTravelListRequest struct {
	Pairs []LocationPair `json:"pairs"`
	RoundOff bool `json:"roundOff"`
}
type CalculateTimeTravelListResponse struct {
	CalculateTravelTimeList struct {
		TravelTimeMinutes []int `json:"travelTimeMinutes"`
	} `json:"calculateTravelTimeList"`
}
