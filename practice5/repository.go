package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

var allowedOrderColumns = map[string]bool{
	"id":         true,
	"name":       true,
	"email":      true,
	"gender":     true,
	"birth_date": true,
}

func (r *Repository) GetPaginatedUsers(f UserFilter) (PaginatedResponse, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 10
	}
	offset := (f.Page - 1) * f.PageSize

	orderCol := "id" 
	if allowedOrderColumns[f.OrderBy] {
		orderCol = f.OrderBy
	}
	orderDir := "ASC"
	if strings.ToLower(f.OrderDir) == "desc" {
		orderDir = "DESC"
	}

	conditions := []string{}
	args := []interface{}{}
	argIdx := 1 
 
	if f.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, *f.ID)
		argIdx++
	}
	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+f.Name+"%")
		argIdx++
	}
	if f.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argIdx))
		args = append(args, "%"+f.Email+"%")
		argIdx++
	}
	if f.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argIdx))
		args = append(args, f.Gender)
		argIdx++
	}
	if f.BirthDate != "" {
		conditions = append(conditions, fmt.Sprintf("birth_date = $%d", argIdx))
		args = append(args, f.BirthDate)
		argIdx++
	}
 
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date, COUNT(*) OVER() AS total_count
		FROM   users
		%s
		ORDER  BY %s %s
		LIMIT  $%d OFFSET $%d`,
		whereClause,
		orderCol, orderDir,
		argIdx, argIdx+1,
	)
 
	args = append(args, f.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return PaginatedResponse{}, fmt.Errorf("GetPaginatedUsers query: %w", err)
	}
	defer rows.Close()
 
	var users []User
	totalCount := 0
 
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate, &totalCount); err != nil {
			return PaginatedResponse{}, fmt.Errorf("GetPaginatedUsers scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return PaginatedResponse{}, fmt.Errorf("GetPaginatedUsers rows: %w", err)
	}
 
	if users == nil {
		users = []User{}
	}
 
	return PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       f.Page,
		PageSize:   f.PageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(userID1, userID2 int) ([]User, error) {
	if userID1 == userID2 {
		return nil, fmt.Errorf("user IDs must be different")
	}
 
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM   users u
		JOIN   user_friends uf1 ON uf1.friend_id = u.id AND uf1.user_id = $1
		JOIN   user_friends uf2 ON uf2.friend_id = u.id AND uf2.user_id = $2`
 
	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("GetCommonFriends query: %w", err)
	}
	defer rows.Close()
 
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, fmt.Errorf("GetCommonFriends scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCommonFriends rows: %w", err)
	}
 
	if users == nil {
		users = []User{}
	}
	return users, nil
}