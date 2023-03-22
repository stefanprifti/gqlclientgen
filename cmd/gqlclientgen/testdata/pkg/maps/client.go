package maps

import (
	"context"
	"encoding/json"

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
	var resp CalculateTimeTravelListResponse

	query := `# CalculateTimeTravelList - Calculate travel time between multiple pairs of locations
query CalculateTimeTravelList($pairs: [LocationPair!]!, $roundOff: Boolean!) {
  calculateTravelTimeList(request: { pairs: $pairs, roundOff: $roundOff }) {
    travelTimeMinutes
  }
}
`

	reqMap, err := structToMap(req)
	if err != nil {
		return nil, err
	}

	
	err = c.gqlclient.Query(ctx, query, reqMap, &resp)
	
	if err != nil {
		return nil, err
	}

	return &resp, nil
}


func structToMap(s interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}