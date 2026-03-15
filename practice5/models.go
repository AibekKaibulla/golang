package main
import "time"
 
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}
 
type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}// "asc" | "desc"
 
type UserFilter struct {
	Page     int
	PageSize int
 
	// sorting  (column name + direction)
	OrderBy  string
	OrderDir string 
 
	ID        *int
	Name      string
	Email     string
	Gender    string
	BirthDate string 
}