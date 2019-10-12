package repositories

import (
	"testing"

	"order-service/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepo_List(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoErrorf(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "distance", "status"}).
		AddRow(1, 100, "UNASSIGNED").
		AddRow(2, 200, "UNASSIGNED")

	query := "SELECT id, distance, status FROM orders ORDER BY id ASC LIMIT ?, ?"

	mock.ExpectQuery(query).
		WithArgs(0, 10).
		WillReturnRows(rows)

	orderRepo := NewMysqlOrderRepo(db)
	orders, err := orderRepo.List(0, 10)
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Equal(t, 2, len(orders))
}

func TestOrderRepo_GetById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoErrorf(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"distance",
		"status"}).AddRow(
		1,
		100,
		"UNASSIGNED")

	query := "SELECT id, distance, status FROM orders WHERE id = ?"

	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(rows)

	orderRepo := NewMysqlOrderRepo(db)
	order, err := orderRepo.GetById(1)
	assert.NoError(t, err)
	assert.NotNil(t, order)
}

func TestOrderRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoErrorf(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	o := &models.Order{
		Origins:      []float64{22.780247, 113.687473},
		Destinations: []float64{22.217851, 114.207989},
		Distance:     100,
		Status:       models.StatusUnassigned,
	}

	query := `INSERT INTO orders (origin_lat, origin_lng, destination_lat, destination_lng, distance, status) VALUES (?, ?, ?, ?, ?, ?)`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(o.Origins[0], o.Origins[1], o.Destinations[0], o.Destinations[1], o.Distance, o.Status).
		WillReturnResult(sqlmock.NewResult(123, 1))

	orderRepo := NewMysqlOrderRepo(db)
	order, err := orderRepo.Create(o)
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, int64(123), order.Id)
}

func TestOrderRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoErrorf(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	o := &models.Order{
		Origins:      []float64{22.780247, 113.687473},
		Destinations: []float64{22.217851, 114.207989},
		Distance:     100,
		Status:       models.StatusTaken,
	}

	query := "UPDATE orders SET status = ? where id = ? AND status = ?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(models.StatusTaken, o.Id, models.StatusUnassigned).
		WillReturnResult(sqlmock.NewResult(1, 1))

	orderRepo := NewMysqlOrderRepo(db)
	order, err := orderRepo.Update(o, models.StatusUnassigned)
	assert.NoError(t, err)
	assert.NotNil(t, order)
}

func TestOrderRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoErrorf(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "DELETE FROM orders WHERE id = ?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	orderRepo := NewMysqlOrderRepo(db)
	order, err := orderRepo.Delete(1)
	assert.NoError(t, err)
	assert.NotNil(t, order)
}
