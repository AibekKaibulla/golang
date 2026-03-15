package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
 
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
 
	filter := UserFilter{
		Page:      page,
		PageSize:  pageSize,
		OrderBy:   q.Get("order_by"),
		OrderDir:  q.Get("order_dir"),
		Name:      q.Get("name"),
		Email:     q.Get("email"),
		Gender:    q.Get("gender"),
		BirthDate: q.Get("birth_date"),
	}
 
	if idStr := q.Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}
		filter.ID = &id
	}
 
	result, err := h.repo.GetPaginatedUsers(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
 
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) GetCommonFriendsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
 
	q := r.URL.Query()
 
	uid1Str := q.Get("user_id_1")
	uid2Str := q.Get("user_id_2")
 
	if uid1Str == "" || uid2Str == "" {
		writeError(w, http.StatusBadRequest, "user_id_1 and user_id_2 are required")
		return
	}
 
	uid1, err := strconv.Atoi(uid1Str)
	if err != nil {
		writeError(w, http.StatusBadRequest, "user_id_1 must be an integer")
		return
	}
	uid2, err := strconv.Atoi(uid2Str)
	if err != nil {
		writeError(w, http.StatusBadRequest, "user_id_2 must be an integer")
		return
	}
 
	friends, err := h.repo.GetCommonFriends(uid1, uid2)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
 
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id_1":      uid1,
		"user_id_2":      uid2,
		"common_friends": friends,
		"count":          len(friends),
	})
}
