package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/petersizovdev/MEDODS-T.git/models"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	rows, err := h.DB.Query("SELECT id, email FROM Users")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var uuidStr string

		if err := rows.Scan(&uuidStr, &user.Email); err != nil {
			fmt.Println(err)
			return
		}

		user.Id, err = uuid.Parse(uuidStr)
		if err != nil {
			fmt.Println(err)
			return
		}

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		fmt.Println(err)
		return
	}

}
