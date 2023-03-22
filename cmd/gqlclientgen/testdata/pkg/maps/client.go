package maps

import (
	"context"

	"github.com/stefanprifti/gqlclient"
)

type Client struct {
	URL       string
	gqlclient *gqlclient.Client
}

func NewClient(url string) *Client {
	return &Client{
		URL: url,
		gqlclient: gqlclient.New(gqlclient.Options{
			Endpoint: url,
		}),
	}
}

func (c *Client) CalculateTimeTravelList(ctx context.Context, req *CalculateTimeTravelListRequest) (*CalculateTimeTravelListResponse, error) {
	query := `# CalculateTimeTravelList - Calculate travel time between multiple pairs of locations
query CalculateTimeTravelList($pairs: [LocationPair!]!, $roundOff: Boolean!) {
  calculateTravelTimeList(request: { pairs: $pairs, roundOff: $roundOff }) {
    travelTimeMinutes
  }
}
`
	var response struct {
		CalculateTravelTimeList *CalculateTimeTravelListResponse `json:"calculateTravelTimeList"`
	}
	err := c.gqlclient.Query(ctx, query, req, &response)
	if err != nil {
		return nil, err
	}

	return response.CalculateTravelTimeList, nil
}
