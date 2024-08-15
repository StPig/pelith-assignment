package services

import (
	"database/sql"
	"pelith-assignment/database"
	"pelith-assignment/models"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	database.DBInstance = mockDB

	mockDB.ExpectExec("insert into users").
		WithArgs("test_address").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = CreateUser("test_address")
	assert.NoError(t, err)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateUser_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	database.DBInstance = mockDB

	mockDB.ExpectExec("insert into users").
		WithArgs("test_address").
		WillReturnError(assert.AnError)

	err = CreateUser("test_address")
	assert.Error(t, err)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserDetail_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	database.DBInstance = &database.PgMockDB{PgxPoolIface: mockDB}

	mockDB.ExpectQuery("select id, address, onboarding_completed, points").
		WithArgs("test_address").
		WillReturnRows(pgxmock.NewRows([]string{"id", "address", "onboarding_completed", "points"}).
			AddRow(1, "test_address", true, 100))

	user, err := getUserDetail("test_address")
	assert.NoError(t, err)
	assert.Equal(t, models.User{
		ID:                  1,
		Address:             "test_address",
		OnboardingCompleted: true,
		Points:              100,
	}, user)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserDetail_NoRows(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	database.DBInstance = &database.PgMockDB{PgxPoolIface: mockDB}

	mockDB.ExpectQuery("select id, address, onboarding_completed, points").
		WithArgs("test_address").
		WillReturnRows(pgxmock.NewRows([]string{"id", "address", "onboarding_completed", "points"}))

	user, err := getUserDetail("test_address")
	assert.NoError(t, err)
	assert.Equal(t, models.User{}, user)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserDetail_QueryError(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	database.DBInstance = &database.PgMockDB{PgxPoolIface: mockDB}

	mockDB.ExpectQuery("select id, address, onboarding_completed, points").
		WithArgs("test_address").
		WillReturnError(sql.ErrConnDone)

	user, err := getUserDetail("test_address")
	assert.Error(t, err)
	assert.Equal(t, models.User{}, user)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}
