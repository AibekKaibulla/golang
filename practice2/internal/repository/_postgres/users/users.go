package users

import (
	"fmt"
	"practice2/internal/repository/_postgres"
	"practice2/pkg/modules"
	"time"
)

type Repository struct {
	db *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{db: db, executionTimeout: time.Second * 5}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "select * from users")
	return users, err
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "select * from users where id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	var newID int
	q := "insert into users (name, email, age) values ($1, $2, $3) returning id"
	err := r.db.DB.QueryRow(q, user.Name, user.Email, user.Age).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return newID, nil
}

func (r *Repository) UpdateUser(id int, user modules.User) error {
    q := `UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4`
    res, err := r.db.DB.Exec(q, user.Name, user.Email, user.Age, id)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("user not found: no rows were updated")
    }
    return nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
    res, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
    if err != nil {
        return 0, fmt.Errorf("failed to delete user: %w", err)
    }
    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return 0, fmt.Errorf("user with id %d not found", id)
    }
    return rowsAffected, nil
}

