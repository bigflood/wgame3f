package view

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"

	//	"gae_wgame/auth"
	"gae_wgame/db"
)

func gameListHandler(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := appengine.NewContext(r)
	games, err := db.DB.ListGames(ctx)
	if err != nil {
		return AppErrorf(err, "could not list games: %v", err)
	}

	return gameListTmpl.Execute(w, r, games)
}

func gameFromRequest(r *http.Request) (*db.Game, error) {
	id := mux.Vars(r)["id"]
	ctx := appengine.NewContext(r)
	game, err := db.DB.GetGame(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find game: %v", err)
	}
	return game, nil
}

func gameDetailHandler(w http.ResponseWriter, r *http.Request) *AppError {
	game, err := gameFromRequest(r)
	if err != nil {
		return AppErrorf(err, "%v", err)
	}

	return gameDetailTmpl.Execute(w, r, game)
}

func gameAddFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	return gameEditTmpl.Execute(w, r, nil)
}

func gameEditFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	game, err := gameFromRequest(r)
	if err != nil {
		return AppErrorf(err, "%v", err)
	}

	return gameEditTmpl.Execute(w, r, game)
}

// gameFromForm populates the fields of a Game from form values
// (see templates/edit.html).
func gameFromForm(r *http.Request) (*db.Game, error) {

	game := &db.Game{
		ID:          r.FormValue("name"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		EditKey:     r.FormValue("editkey"),
	}

	return game, nil
}

// updateHandler updates the details of a given game.
func gameUpdateHandler(w http.ResponseWriter, r *http.Request) *AppError {
	game, err := gameFromForm(r)
	if err != nil {
		return AppErrorf(err, "could not parse game from form: %v", err)
	}

	if game.ID == "" {
		return AppErrorf(fmt.Errorf("empty ID"), "empty ID")
	}

	ctx := appengine.NewContext(r)

	err = db.DB.UpdateGame(ctx, game)
	if err != nil {
		return AppErrorf(err, "could not save game: %v", err)
	}
	//go publishUpdate(game.ID)
	http.Redirect(w, r, fmt.Sprintf("/games/%s", game.ID), http.StatusFound)
	return nil
}

func gameDeleteHandler(w http.ResponseWriter, r *http.Request) *AppError {
	id := mux.Vars(r)["id"]

	ctx := appengine.NewContext(r)
	err := db.DB.DeleteGame(ctx, id)
	if err != nil {
		return AppErrorf(err, "could not delete game: %v", err)
	}
	http.Redirect(w, r, "/games", http.StatusFound)
	return nil
}
