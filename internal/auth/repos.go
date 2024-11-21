package auth

import (
	"errors"
	"time"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type UserRepository interface {
	Create(user *CreateUser) (*User, error)
	GetByUsername(username string) (*User, error)
	GetById(id int64) (*User, error)
	DeleteById(id int64) error
}

type PgUserRepository struct {
	db     *pgx.Conn
	logger *zap.Logger
}

func NewPgUserRepository(db *pgx.Conn, logger *zap.Logger) UserRepository {
	return &PgUserRepository{db, logger}
}

const createUserSql = `
	INSERT INTO users (first_name, last_name, username, role, password_hash, password_salt, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at
`

func (r *PgUserRepository) Create(user *CreateUser) (*User, error) {
	passwordSalt, err := GenerateSalt(16)
	if err != nil {
		return nil, nil
	}

	passwordHash := HashPassword(user.Password, passwordSalt)
	r.logger.Debug("Executing query", zap.String("query", createUserSql))

	var newUser User
	newUser.FirstName = user.FirstName
	newUser.LastName = user.LastName
	newUser.Username = user.Username
	newUser.Role = user.Role
	newUser.PasswordHash = passwordHash
	newUser.PasswordSalt = passwordSalt

	err = r.db.QueryRow(createUserSql,
		newUser.FirstName,
		newUser.LastName,
		newUser.Username,
		newUser.Role,
		newUser.PasswordHash,
		newUser.PasswordSalt,
		time.Now(),
	).Scan(&newUser.Id, &newUser.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

const getByIdSql = `
	SELECT id, first_name, last_name, username, role, password_hash
	FROM users
	WHERE id = $1
`

func (r *PgUserRepository) GetById(id int64) (*User, error) {
	r.logger.Debug("Executing query", zap.String("query", getByIdSql))

	var user User
	row := r.db.QueryRow(getByIdSql, id)

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Role,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

const getByUsernameSql = `
	SELECT id, first_name, last_name, username, role, password_hash
	FROM users
	WHERE username = $1
`

func (r *PgUserRepository) GetByUsername(username string) (*User, error) {
	r.logger.Debug("Executing query", zap.String("query", getByUsernameSql))

	var user User
	row := r.db.QueryRow(getByUsernameSql, username)

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Role,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

const deleteByIdSql = `
	DELETE FROM users
	WHERE id = $1
`

func (r *PgUserRepository) DeleteById(id int64) error {
	r.logger.Debug("Executing query", zap.String("query", deleteByIdSql))
	_, err := r.db.Exec(deleteByIdSql, id)
	return err
}
