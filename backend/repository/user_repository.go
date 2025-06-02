package repository

import (
	"context"
	"fmt"
	"gobunker/model"
	"gobunker/model/dto"
)

type userRepository struct {
	exec SQLExecutor
}

type UserRepository interface {
	CreateUser(ctx context.Context, userDTO dto.UserDTO) (model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int) (model.User, error)
	GetUserByName(ctx context.Context, name string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, userDTO dto.UserDTO) (model.User, error)
	DeleteUserByID(ctx context.Context, id int) error
}

func NewUserRepository(exec SQLExecutor) UserRepository {
	return &userRepository{exec: exec}
}

func (u *userRepository) CreateUser(ctx context.Context, userDTO dto.UserDTO) (model.User, error) {
	var user model.User

	query := `INSERT INTO users (name, email, password, role) VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at`
	row := u.exec.QueryRowContext(ctx, query, userDTO.Name, userDTO.Email, userDTO.Password, userDTO.Role)
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("repo : failed to insert user : %w", err)
	}

	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.Password = userDTO.Password
	user.Role = userDTO.Role

	return user, nil
}

func (u *userRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	query := `SELECT id,name,email,password,role,created_at,updated_at FROM users`
	rows, err := u.exec.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo : failed to get all users : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan all users %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	query := `SELECT id,name,email,password,role,created_at,updated_at FROM users WHERE id = $1`
	row := u.exec.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("repo : failed to get user by id : %w", err)
	}

	return user, nil
}

func (u *userRepository) GetUserByName(ctx context.Context, name string) (model.User, error) {
	var user model.User
	query := `SELECT id,name,email,password,role,created_at,updated_at FROM users WHERE name = $1`
	row := u.exec.QueryRowContext(ctx, query, name)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("repo : failed to get user by name : %w", err)
	}

	return user, nil

}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `SELECT id,name,email,password,role,created_at,updated_at FROM users WHERE email = $1`
	row := u.exec.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("repo : failed to get user by email : %w", err)
	}

	return user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, userDTO dto.UserDTO) (model.User, error) {
	var user model.User
	query := `
	UPDATE users 
	SET name = $1, email = $2, password = $3, role = $4, updated_at = NOW()
	WHERE id = $2
	RETURNING id, created_at, updated_at
	`
	row := u.exec.QueryRowContext(
		ctx,
		query,
		userDTO.Name,
		userDTO.Email,
		userDTO.Password,
		userDTO.Role,
	)
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return model.User{}, fmt.Errorf("repo : failed to update user : %w", err)
	}

	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.Password = userDTO.Password
	user.Role = userDTO.Role

	return user, nil
}

func (u *userRepository) DeleteUserByID(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := u.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repo : failed to delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("repo: no user deleted")
	}

	return nil
}
