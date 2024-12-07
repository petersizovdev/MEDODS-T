package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"


)

type UserHandler struct{
	DB *sql.DB
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request){
	
}