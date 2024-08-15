package services

import (
	"errors"
	"pelith-assignment/database"
	"pelith-assignment/models"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDistributePoints_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	dbAdapter := &database.PgMockDB{PgxPoolIface: mockDB}
	database.DBInstance = dbAdapter

	mockTx := pgxmock.NewTx()
	mockDB.ExpectBegin().WillReturnTx(mockTx, nil)
	mockTx.ExpectExec("update users set points = points + $1").
		WithArgs(50.0, "test_address").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mockTx.ExpectExec("insert into complete_history").
		WithArgs(1, 2, 50.0).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mockTx.ExpectCommit()

	user := &models.User{ID: 1, Address: "test_address", OnboardingCompleted: false, Points: 0}
	task := &models.Task{ID: 2}

	err = database.DistributePoints(user, task, 50.0)
	assert.NoError(t, err)
	assert.True(t, user.OnboardingCompleted)
	assert.Equal(t, 50, user.Points)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDistributePoints_UpdatePoints_Failure(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	dbAdapter := &database.PgMockDB{PgxPoolIface: mockDB}
	database.DBInstance = dbAdapter

	mockTx := pgxmock.NewTx()
	mockDB.ExpectBegin().WillReturnTx(mockTx, nil)
	mockTx.ExpectExec("update users set points = points + $1").
		WithArgs(50.0, "test_address").
		WillReturnError(errors.New("update error"))
	mockTx.ExpectRollback()

	user := &models.User{ID: 1, Address: "test_address", OnboardingCompleted: false, Points: 0}
	task := &models.Task{ID: 2}

	err = database.DistributePoints(user, task, 50.0)
	assert.Error(t, err)
	assert.False(t, user.OnboardingCompleted)
	assert.Equal(t, 0, user.Points)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDistributePoints_InsertHistory_Failure(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	dbAdapter := &database.PgMockDB{PgxPoolIface: mockDB}
	database.DBInstance = dbAdapter

	mockTx := pgxmock.NewTx()
	mockDB.ExpectBegin().WillReturnTx(mockTx, nil)
	mockTx.ExpectExec("update users set points = points + $1").
		WithArgs(50.0, "test_address").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mockTx.ExpectExec("insert into complete_history").
		WithArgs(1, 2, 50.0).
		WillReturnError(errors.New("insert error"))
	mockTx.ExpectRollback()

	user := &models.User{ID: 1, Address: "test_address", OnboardingCompleted: false, Points: 0}
	task := &models.Task{ID: 2}

	err = database.DistributePoints(user, task, 50.0)
	assert.Error(t, err)
	assert.False(t, user.OnboardingCompleted)
	assert.Equal(t, 0, user.Points)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}
