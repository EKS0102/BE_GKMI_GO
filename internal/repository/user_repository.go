package repository

import (
	"context"

	"BE_GKMI_NTC_GO/internal/model"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	DB *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]model.User, error) {

	rows, err := r.DB.Query(
		ctx,
		`
		SELECT id, username, role, is_active, email
		FROM users
		ORDER BY id
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {

		var user model.User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Role,
			&user.IsActive,
			&user.Email,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {

	err := r.DB.QueryRow(
		ctx,
		`
		INSERT INTO users (
			username,
			password_hash,
			role,
			is_active,
			email
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.IsActive,
		user.Email,
	).Scan(&user.ID)

	return err
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	id int,
) (*model.User, error) {
	var user model.User

	err := r.DB.QueryRow(
		ctx,
		`
		SELECT id, username, role, is_active, email
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Role,
		&user.IsActive,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {
	var user model.User

	err := r.DB.QueryRow(
		ctx,
		`
		SELECT id, username, password_hash, role, is_active, email
		FROM users
		WHERE username = $1
		`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(
	ctx context.Context,
	user *model.User,
) error {
	_, err := r.DB.Exec(
		ctx,
		`
		UPDATE users
		SET username = $1,
		    role = $2,
		    is_active = $3,
		    email = $4
		WHERE id = $5
		`,
		user.Username,
		user.Role,
		user.IsActive,
		user.Email,
		user.ID,
	)

	return err
}

func (r *UserRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	_, err := r.DB.Exec(
		ctx,
		`
		DELETE FROM users
		WHERE id = $1
		`,
		id,
	)

	return err
}
