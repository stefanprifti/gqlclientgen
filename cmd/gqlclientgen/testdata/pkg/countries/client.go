package countries

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

func (c *Client) Country(ctx context.Context, req *CountryRequest) (*CountryResponse, error) {
	var resp CountryResponse

	query := `query Country($code: ID!) {
  country(code: $code) {
    name
    native
    languages {
      code
      name
    }
    emoji
    currency
    languages {
      code
      name
    }
  }
}
`

	err := c.gqlclient.Query(ctx, query, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
