package main

import "net/http"

type category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (app *application) categoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories := []category{
		{ID: 1, Name: "Food"},
	}
	env := envelope{
		"categories": categories,
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
