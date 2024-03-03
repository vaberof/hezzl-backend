package views

import (
	"github.com/go-chi/render"
	"github.com/vaberof/hezzl-backend/pkg/http/protocols/apiv1"
	"net/http"
)

type httpStatus int

// RenderJSON is a wrapper for go-chi json render.
func RenderJSON(w http.ResponseWriter, r *http.Request, status httpStatus, payload *apiv1.Response) {
	w.WriteHeader(int(status))
	render.JSON(w, r, payload)
}
