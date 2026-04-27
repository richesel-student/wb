package bybit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetPrice() float64 {
	url := "https://api-testnet.bybit.com/v5/market/tickers?category=spot&symbol=BTCUSDT"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP ERROR:", err)
		return 70000
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("BAD STATUS:", resp.Status)
		return 70000
	}

	var data struct {
		Result struct {
			List []struct {
				LastPrice string `json:"lastPrice"`
			} `json:"list"`
		} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("JSON ERROR:", err)
		return 70000
	}

	if len(data.Result.List) == 0 {
		return 70000
	}

	var price float64
	fmt.Sscanf(data.Result.List[0].LastPrice, "%f", &price)

	return price
}
