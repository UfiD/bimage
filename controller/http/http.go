package controller

import (
	"bimage/controller/http/types"
	"bimage/usecase"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	uc usecase.Object
}

func New(uc usecase.Object) *Controller {
	return &Controller{
		uc: uc,
	}
}

func (c *Controller) WithObjectHandler() chi.Router {
	r := chi.NewRouter()
	r.Post("/", c.Post)
	return r
}

func CreateAndRunServer(addr string, r chi.Router) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQUEST")
	task, err := types.GetObjectPostRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := c.uc.Do(task)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
