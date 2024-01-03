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
		locationAreaResp := LocationAreaResp{}
		err := json.Unmarshal(dat, &locationAreaResp)
		if err != nil {
			return LocationAreaResp{}, err
		}

		return locationAreaResp, nil
	}
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

func (c *Client) ExploreLocation(area string) (ExploreAreaResp, error) {
	endpoint := "/location-area"
	fullURL := fmt.Sprintf("%s%s/%s", baseURL, endpoint, area)

	// check the cache
	dat, ok := c.cache.Get(fullURL)
	if ok {
		exploreAreaResp := ExploreAreaResp{}
		err := json.Unmarshal(dat, &exploreAreaResp)
		if err != nil {
			return ExploreAreaResp{}, err
		}

		return exploreAreaResp, nil
	}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return ExploreAreaResp{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ExploreAreaResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return ExploreAreaResp{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExploreAreaResp{}, err
	}

	exploreAreaResp := ExploreAreaResp{}
	err = json.Unmarshal(data, &exploreAreaResp)
	if err != nil {
		return ExploreAreaResp{}, err
	}

	c.cache.Add(fullURL, data)

	return exploreAreaResp, nil
}
