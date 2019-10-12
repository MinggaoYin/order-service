package distance

import (
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"order-service/services"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDistanceService_GetDistance(t *testing.T) {
	loadEnv()

	service, err := NewDistanceService()
	assert.NoError(t, err)
	assert.NotNil(t, service)

	t.Run("success", func(t *testing.T) {
		origins := []string{"22.286681", "114.193260"}
		destinations := []string{"22.279707", "114.186301"}

		_, err = service.GetDistance(
			[]string{strings.Join(origins, ",")},
			[]string{strings.Join(destinations, ",")})
		assert.NoError(t, err)
	})

	t.Run("on-failed cannot get distance between two location", func(t *testing.T) {
		origins := []string{"22.545285", "114.125790"}
		destinations := []string{"22.279707", "114.186301"}

		_, err = service.GetDistance(
			[]string{strings.Join(origins, ",")},
			[]string{strings.Join(destinations, ",")})
		assert.EqualError(t, err, services.ErrCannotCalculateDistance.Error())
	})

}

const projectDirName = "order-service"

// loadEnv loads env vars from .env
func loadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Problem loading .env file")

		os.Exit(-1)
	}
}
