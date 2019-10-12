package services

type DistanceCalculator interface {
	GetDistance(origins, destinations []string) (int, error)
}
