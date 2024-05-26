package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

type MobiFitnessApiClient interface {
	GetClubsList() ([]*models.ClubListResponse, error)
	GetClubInfo(clubID int) (*models.ClubInfoResponse, error)
}

type Storage interface {
	SaveClubLoad(ctx context.Context, load *models.ClubInfoResponse) error
}

type Club struct {
	mobiFitnessApiClient MobiFitnessApiClient
	storage              Storage
}

func NewClub(mobiFitnessApiClient MobiFitnessApiClient, storage Storage) *Club {
	return &Club{mobiFitnessApiClient: mobiFitnessApiClient, storage: storage}
}

func (s *Club) StartCron(ctx context.Context) error {
	clubs, err := s.mobiFitnessApiClient.GetClubsList()
	if err != nil {
		return fmt.Errorf("failed to get clubs list: %w", err)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Print("Shutting down gracefully...")
			return nil
		case <-ticker.C:
			if err := s.SaveClubsLoad(ctx, clubs); err != nil {
				logrus.Errorf("error saving clubs load: %v", err)
				time.Sleep(15 * time.Second)
			} else {
				logrus.Print("Successfully saved clubs load")
			}
		}
	}
}

func (s *Club) SaveClubsLoad(ctx context.Context, clubs []*models.ClubListResponse) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, club := range clubs {
		club := club
		if club.City != "Тюмень" {
			continue
		}

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				clubInfo, err := s.mobiFitnessApiClient.GetClubInfo(club.Id)
				if err != nil {
					return fmt.Errorf("failed to get club info for club ID %d: %w", club.Id, err)
				}

				if clubInfo.CurrentLoad == "" {
					return errors.New("current load is empty")
				}

				if err := s.storage.SaveClubLoad(ctx, clubInfo); err != nil {
					return fmt.Errorf("failed to save club load for club ID %d: %w", club.Id, err)
				}
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
