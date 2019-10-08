package services

import (
	"context"
	"os"

	"order-service/startup"

	"github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)

var c *maps.Client

func init() {
	log := logrus.WithFields(logrus.Fields{"module": "service/distance", "method": "init"})

	key := startup.Config.GoogleApiKey

	var err error
	c, err = maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		log.WithError(err).Error("Failed to init client")
		os.Exit(1)
	}
}

func getDistance(origins, destinations []string) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service/distance", "method": "getDistance", "origins": origins, "destinations": destinations})

	r := &maps.DistanceMatrixRequest{
		Origins:      origins,
		Destinations: destinations,
		Mode:         maps.TravelModeDriving,
	}

	resp, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.WithError(err).Error("Failed to get distance from google")
		return 0, err
	}

	if len(resp.Rows) != 1 ||
		len(resp.Rows[0].Elements) != 1 ||
		resp.Rows[0].Elements[0].Status != "OK" {
		log.Error("Failed to calculate distance for give location")
		return 0, ErrCannotCalculateDistance
	}

	return resp.Rows[0].Elements[0].Distance.Meters, nil
}
