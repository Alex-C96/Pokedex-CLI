package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResp, error) {
	endpoint := "/location-area"
	offsetLimit := "?offset=0&limit=20"
	fullURL := baseURL + endpoint + offsetLimit
	if pageURL != nil {
		fullURL = *pageURL
	}

	// check the cache
	dat, ok := c.cache.Get(fullURL)
	if ok {
		fmt.Println("Cache Hit!")
		locationAreaResp := LocationAreaResp{}
		err := json.Unmarshal(dat, &locationAreaResp)
		if err != nil {
			return LocationAreaResp{}, err
		}

		return locationAreaResp, nil
	}
	fmt.Println("Cache Miss!")
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreaResp{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return LocationAreaResp{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResp{}, err
	}

	locationAreaResp := LocationAreaResp{}
	err = json.Unmarshal(data, &locationAreaResp)
	if err != nil {
		return LocationAreaResp{}, err
	}

	c.cache.Add(fullURL, data)

	return locationAreaResp, nil
}
