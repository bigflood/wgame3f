package view

import (
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
)

type iapValidationRequest struct {
	Data string
}

type iapValidationResult struct {
	ID, Platform string
	URL1         string
	Error        string
	Data         string
}

func iapValidationFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	id := mux.Vars(r)["id"]
	platform := mux.Vars(r)["platform"]

	info := struct {
		ID, Platform string
	}{
		ID:       id,
		Platform: platform,
	}

	return iapValidationFormTmpl.Execute(w, r, info)
}

func iapValidationHandler(w http.ResponseWriter, r *http.Request) *AppError {
	info, err := iapValidationRequestFromForm(r)
	if err != nil {
		return AppErrorf(err, "could not parse request from form: %v", err)
	}

	ctx := appengine.NewContext(r)

	id := mux.Vars(r)["id"]
	platform := mux.Vars(r)["platform"]

	result := &iapValidationResult{}
	result.ID = id
	result.Platform = platform
	result.Data = info.Data

	if id == "" || platform == "" {
		result.Error = "incorrect id or platform"
	} else if info.Data == "" {
		result.Error = "empty data"
	} else if platform == "google" {
		msg, err := iapValidationGoogle(ctx, id, info.Data, result)
		if err != nil {
			result.Error = err.Error()
		}
		result.Data = msg
	} else if platform == "apple" {
		msg, err := iapValidationApple(ctx, id, info.Data, result)
		if err != nil {
			result.Error = err.Error()
		}
		result.Data = msg
	} else {
		result.Error = "invalid platform"
	}

	return iapValidationResultTmpl.Execute(w, r, result)
}

func iapValidationRequestFromForm(r *http.Request) (*iapValidationRequest, error) {

	info := &iapValidationRequest{
		Data: r.FormValue("data"),
	}

	return info, nil
}
