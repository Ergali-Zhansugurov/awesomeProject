package http

import (
	"encoding/json"
	"fmt"

	"awesomeProject/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func (cr BaseResource) CategoryRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", cr.CreateCategory)
	r.Get("/", cr.GetCategories)
	r.Delete("/", cr.DeleteCategory)
	return r
}
func (cr BaseResource) CreateCategory(writer http.ResponseWriter, request *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(request.Body).Decode(category); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(writer, "unknow err: %v", err)
		return
	}
	if err := cr.store.Category().Create(request.Context(), category); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}

	cr.broker.Cache().Purge()
	writer.WriteHeader(http.StatusCreated)
}
func (cr BaseResource) GetCategories(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.Filter{}
	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	categories, err := cr.store.Category().Get(request.Context(), filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, categories)
}
func (cr BaseResource) DeleteCategory(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	eror := cr.store.Category().Delete(request.Context(), id)
	if eror != nil {
		fmt.Fprintf(writer, "Unknow db delete err : %s", err)
		return
	}
	cr.broker.Cache().Purge()
}
