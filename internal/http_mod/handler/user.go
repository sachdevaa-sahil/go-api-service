package handler

import "net/http"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// dummyUsers supplies both endpoints until persistent storage is added.
func dummyUsers() []User {
	return []User{
		{ID: "1", Name: "Sahil"},
		{ID: "2", Name: "Guggu"},
		{ID: "3", Name: "Advik"},
		{ID: "4", Name: "Adira"},
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, dummyUsers())
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, user := range dummyUsers() {
		if user.ID == id {
			writeJSON(w, http.StatusOK, user)
			return
		}
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
}
