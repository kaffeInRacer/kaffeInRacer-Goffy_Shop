package usersRepository

import (
	"context"
	"database/sql"
	"errors"
	"kaffein/domain"
)

type UsersRepository struct {
	Conn *sql.DB
}

func NewUsersRepository(conn *sql.DB) *UsersRepository {
	return &UsersRepository{Conn: conn}
}

// Store - Menyimpan pengguna baru
func (c *UsersRepository) Store(ctx context.Context, u *domain.Users) error {
	rawQuery := `
		INSERT INTO goffy_users (name, password, email, phone_number, gender)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := c.Conn.QueryRowContext(ctx, rawQuery, u.Name, u.Password, u.Email, u.PhoneNumber, u.Gender).Scan(&u.Id)
	if err != nil {
		return err
	}

	return nil
}

// FindById - Mencari pengguna berdasarkan ID
func (c *UsersRepository) FindById(ctx context.Context, id int64) (*domain.Users, error) {
	rawQuery := `
		SELECT id, name, password, email, phone_number, gender, created_at, updated_at
		FROM goffy_users WHERE id = $1
	`

	user := &domain.Users{}
	err := c.Conn.QueryRowContext(ctx, rawQuery, id).Scan(
		&user.Id, &user.Name, &user.Password, &user.Email, &user.PhoneNumber,
		&user.Gender, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return user, nil
}

// Delete - Menghapus pengguna berdasarkan ID
func (c *UsersRepository) Delete(ctx context.Context, id int64) error {
	rawQuery := `DELETE FROM goffy_users WHERE id = $1`

	res, err := c.Conn.ExecContext(ctx, rawQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Update - Memperbarui data pengguna
func (c *UsersRepository) Update(ctx context.Context, user *domain.Users) error {
	rawQuery := `
		UPDATE goffy_users 
		SET name = $1, password = $2, email = $3, phone_number = $4, gender = $5 
		WHERE id = $6
	`

	res, err := c.Conn.ExecContext(ctx, rawQuery,
		user.Name, user.Password, user.Email, user.PhoneNumber, user.Gender, user.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FetchAll - Mengambil semua pengguna
func (c *UsersRepository) FetchAll(ctx context.Context) ([]*domain.Users, error) {
	rawQuery := `
		SELECT id, name, password, email, phone_number, gender, created_at, updated_at FROM goffy_users
	`

	rows, err := c.Conn.QueryContext(ctx, rawQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.Users
	for rows.Next() {
		user := &domain.Users{}
		err := rows.Scan(
			&user.Id, &user.Name, &user.Password, &user.Email, &user.PhoneNumber,
			&user.Gender, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
