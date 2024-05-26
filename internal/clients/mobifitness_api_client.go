package clients

import (
	"fmt"
	"github.com/ZetoOfficial/agym-cur-load-cron/config"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/models"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type MobiFitnessApi struct {
	client *resty.Client
}

func NewMobiFitnessApi(conf config.MobiFitness) *MobiFitnessApi {
	client := resty.New()
	client.SetAuthToken(conf.AccessToken)
	client.SetBaseURL(conf.ApiURL)
	return &MobiFitnessApi{client: client}
}

func (c *MobiFitnessApi) GetClubsList() ([]*models.ClubListResponse, error) {
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		Get("/api/v8/franchise/clubs.json")
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var list []*models.ClubListResponse
	if err := jsoniter.Unmarshal(resp.Body(), &list); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return list, nil
}

func (c *MobiFitnessApi) GetClubInfo(clubID int) (*models.ClubInfoResponse, error) {
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("/api/v8/clubs/%v.json", clubID))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	info := &models.ClubInfoResponse{}
	if err := jsoniter.Unmarshal(resp.Body(), info); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return info, nil
}
