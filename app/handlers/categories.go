package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ungame/timetrack/app/service"
	"github.com/ungame/timetrack/httpext"
	"net/http"
	"strconv"
)

type categoriesHandler struct {
	categoriesService service.CategoriesService
}

func NewCategoriesHandler(categoriesService service.CategoriesService) Handler {
	return &categoriesHandler{categoriesService: categoriesService}
}

func (c *categoriesHandler) Register(router *mux.Router) {
	router.Path("/categories/{id}").HandlerFunc(c.GetCategory).Methods(http.MethodGet)
	router.Path("/categories").HandlerFunc(c.GetCategories).Methods(http.MethodGet)
}

func (c *categoriesHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	category, err := c.categoriesService.GetCategory(r.Context(), id)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, category)
}

func (c *categoriesHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.categoriesService.GetCategories(r.Context())
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, categories)
}
