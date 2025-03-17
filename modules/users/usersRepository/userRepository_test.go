package usersRepository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"kaffein/domain"
)

func TestUsersRepository_Store(t *testing.T) {
	testCases := []struct {
		name      string
		user      *domain.Users
		expectErr bool
	}{
		{
			name: "successfully",
			user: &domain.Users{
				Name:        "xzibitz",
				Password:    "xzibitz123",
				Email:       "xzibitz@gmail.com",
				PhoneNumber: "08811223344",
				Gender:      "male",
			},
			expectErr: false,
		},
		{
			name: "empty name",
			user: &domain.Users{
				Name:        "",
				Password:    "xzibitz123",
				Email:       "xzibitz@gmail.com",
				PhoneNumber: "08811223344",
				Gender:      "male",
			},
			expectErr: true,
		},
		{
			name: "empty password",
			user: &domain.Users{
				Name:        "xzibitz",
				Password:    "",
				Email:       "xzibitz@gmail.com",
				PhoneNumber: "08811223344",
				Gender:      "male",
			},
			expectErr: true,
		},
		{
			name: "empty email",
			user: &domain.Users{
				Name:        "xzibitz",
				Password:    "xzibitz123",
				Email:       "",
				PhoneNumber: "08811223344",
				Gender:      "male",
			},
			expectErr: true,
		},
		{
			name: "empty phone number",
			user: &domain.Users{
				Name:        "xzibitz",
				Password:    "xzibitz123",
				Email:       "",
				PhoneNumber: "",
				Gender:      "male",
			},
		},
		{
			name: "empty or invalid gender",
			user: &domain.Users{
				Name:        "xzibitz",
				Password:    "xzibitz123",
				Email:       "xzibitz@gmail.com",
				PhoneNumber: "08811223344",
				Gender:      "invalid",
			},
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening mock DB: %s", err)
	}
	defer db.Close()

	repo := NewUsersRepository(db)

	rawQuery := `
		INSERT INTO goffy_users (name, password, email, phone_number, gender)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectErr {
				mock.ExpectQuery(rawQuery).WillReturnError(sql.ErrNoRows)
			} else {
				mock.ExpectQuery(rawQuery).
					WithArgs(tc.user.Name, tc.user.Password, tc.user.Email, tc.user.PhoneNumber, tc.user.Gender).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			}

			err := repo.Store(context.Background(), tc.user)

			if tc.expectErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "did not expect an error but got one")
			}

			assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
		})
	}
}

func TestUsersRepository_FindById(t *testing.T) {
	testCases := []struct {
		name      string
		userID    int64
		expectErr bool
	}{
		{
			name:      "successfully",
			userID:    1,
			expectErr: false,
		},
		{
			name:      "empty id",
			userID:    0,
			expectErr: true,
		},
		{
			name:      "not found",
			userID:    int64(rand.New(rand.NewSource(time.Now().Unix())).Uint64()),
			expectErr: true,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening mock DB: %s", err)
	}
	defer db.Close()

	repo := NewUsersRepository(db)

	rawQuery := `
		SELECT id, name, password, email, phone_number, gender, created_at, updated_at
		FROM goffy_users WHERE id = $1
	`

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expectErr {
				rows := sqlmock.NewRows([]string{"id", "name", "password", "email", "phone_number", "gender", "created_at", "updated_at"}).
					AddRow(tc.userID, "xzibitz", "xzibit123", "xzibitz@gmail.com", "08811223344", "male", time.Now(), time.Now())

				mock.ExpectQuery(rawQuery).
					WithArgs(tc.userID).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(rawQuery).
					WithArgs(tc.userID).
					WillReturnError(sql.ErrNoRows)
			}

			user, err := repo.FindById(context.TODO(), tc.userID)

			if tc.expectErr {
				assert.Error(t, err, "expected an error but got none")
				assert.Nil(t, user, "expected user to be nil")
			} else {
				assert.NoError(t, err, "did not expect an error but got one")
				assert.NotNil(t, user, "expected user to be not nil")
				assert.Equal(t, tc.userID, user.Id)
			}

			assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
		})
	}
}

func TestUsersRepository_Update(t *testing.T) {
	testCases := []struct {
		name      string
		user      *domain.Users
		mockExec  func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "successfully updated",
			user: &domain.Users{
				Id:          1,
				Name:        "xzibitz",
				Password:    "xzibitz123",
				Email:       "xzibitz@gmail.com",
				PhoneNumber: "08811223344",
				Gender:      "male",
			},
			mockExec: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE goffy_users SET name = $1, password = $2, email = $3, phone_number = $4, gender = $5 WHERE id = $6`).
					WithArgs("xzibitz", "xzibitz123", "xzibitz@gmail.com", "08811223344", "male", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectErr: false,
		},
		{
			name: "update failed due to no rows affected",
			user: &domain.Users{
				Id:          1,
				Name:        "xzibitz",
				Password:    "xzibitz123",
				Email:       "xzibitz@gmail.com",
				PhoneNumber: "08811223344",
				Gender:      "male",
			},
			mockExec: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE goffy_users SET name = $1, password = $2, email = $3, phone_number = $4, gender = $5 WHERE id = $6`).
					WithArgs("xzibitz", "xzibitz123", "xzibitz@gmail.com", "08811223344", "male", 1).
					WillReturnResult(sqlmock.NewResult(0, 0)) // Tidak ada baris yang diperbarui
			},
			expectErr: true,
		},
		{
			name: "database error on update",
			user: &domain.Users{
				Id:          1,
				Name:        "xzibitz",
				Password:    "xzibitz123",
				Email:       "xzibitz@gmail.com",
				PhoneNumber: "08811223344",
				Gender:      "male",
			},
			mockExec: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE goffy_users SET name = $1, password = $2, email = $3, phone_number = $4, gender = $5 WHERE id = $6`).
					WithArgs("xzibitz", "xzibitz123", "xzibitz@gmail.com", "08811223344", "male", 1).
					WillReturnError(errors.New("database error"))
			},
			expectErr: true,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err, "unexpected error when opening mock DB")
	defer db.Close()

	repo := NewUsersRepository(db)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExec(mock)

			err := repo.Update(context.Background(), tc.user)
			if tc.expectErr {
				assert.Error(t, err, "expected an error but got nil")
			} else {
				assert.NoError(t, err, "unexpected error")
			}

			assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled database expectations")
		})
	}
}

func TestUsersRepository_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		userID    int
		mockExec  func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:   "successfully deleted",
			userID: 1,
			mockExec: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM goffy_users WHERE id = $1`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected (berhasil dihapus)
			},
			expectErr: false,
		},
		{
			name:   "delete failed - no rows affected",
			userID: 2,
			mockExec: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM goffy_users WHERE id = $1`).
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(0, 0)) // Tidak ada baris yang terhapus
			},
			expectErr: true,
		},
		{
			name:   "database error on delete",
			userID: 3,
			mockExec: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM goffy_users WHERE id = $1`).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expectErr: true,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err, "unexpected error when opening mock DB")
	defer db.Close()

	repo := NewUsersRepository(db)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExec(mock)

			err := repo.Delete(context.Background(), int64(tc.userID))
			if tc.expectErr {
				assert.Error(t, err, "expected an error but got nil")
			} else {
				assert.NoError(t, err, "unexpected error")
			}

			assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled database expectations")
		})
	}
}

func TestUsersRepository_FetchAll(t *testing.T) {
	testCases := []struct {
		name       string
		mockQuery  func(mock sqlmock.Sqlmock)
		expectData []*domain.Users
		expectErr  bool
	}{
		{
			name: "successfully fetch all users",
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "password", "email", "phone_number", "gender", "created_at", "updated_at"}).
					AddRow(1, "xzibitz", "xzibitz123", "xzibitz@gmail.com", "08811223344", "male", time.Now(), time.Now()).
					AddRow(2, "alice", "alice123", "alice@gmail.com", "08855667788", "female", time.Now(), time.Now())

				mock.ExpectQuery(`SELECT id, name, password, email, phone_number, gender, created_at, updated_at FROM goffy_users`).
					WillReturnRows(rows)
			},
			expectData: []*domain.Users{
				{
					Id:          1,
					Name:        "xzibitz",
					Password:    "xzibitz123",
					Email:       "xzibitz@gmail.com",
					PhoneNumber: "08811223344",
					Gender:      "male",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Id:          2,
					Name:        "alice",
					Password:    "alice123",
					Email:       "alice@gmail.com",
					PhoneNumber: "08855667788",
					Gender:      "female",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			expectErr: false,
		},
		{
			name: "fetch returns empty result",
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "password", "email", "phone_number", "gender", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT id, name, password, email, phone_number, gender, created_at, updated_at FROM goffy_users`).
					WillReturnRows(rows)
			},
			expectData: []*domain.Users{},
			expectErr:  false,
		},
		{
			name: "fetch query error",
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, password, email, phone_number, gender, created_at, updated_at FROM goffy_users`).
					WillReturnError(sql.ErrNoRows)
			},
			expectData: nil,
			expectErr:  true,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err, "unexpected error when opening mock DB")
	defer db.Close()

	repo := NewUsersRepository(db)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockQuery(mock)

			data, err := repo.FetchAll(context.Background())

			if tc.expectErr {
				assert.Error(t, err, "expected an error but got nil")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.Equal(t, len(tc.expectData), len(data), "expected and actual data length mismatch")
			}

			assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled database expectations")
		})
	}
}
