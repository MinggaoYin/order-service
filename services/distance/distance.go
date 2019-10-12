package distance

import (
	"context"
	"os"

	"order-service/services"

	"github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)

type distanceService struct {
	c *maps.Client
}

func NewDistanceService() (services.DistanceCalculator, error) {
	key := os.Getenv("GOOGLE_API_KEY")

	c, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		return nil, err
	}

	return &distanceService{c}, nil
}

func (s *distanceService) GetDistance(origins, destinations []string) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service/distance", "method": "GetDistance", "origins": origins, "destinations": destinations})

	r := &maps.DistanceMatrixRequest{
		Origins:      origins,
		Destinations: destinations,
		Mode:         maps.TravelModeDriving,
	}

	resp, err := s.c.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.WithError(err).Error("Failed to get distance from google")
		return 0, err
	}

	if len(resp.Rows) != 1 ||
		len(resp.Rows[0].Elements) != 1 ||
		resp.Rows[0].Elements[0].Status != "OK" {
		log.Error("Failed to calculate distance for give location")
		return 0, services.ErrCannotCalculateDistance
	}

	return resp.Rows[0].Elements[0].Distance.Meters, nil
}
