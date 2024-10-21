package http

import (
	"awesomeProject/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (tr BaseResource) TitleRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", tr.CreateTitle)
	r.Get("/", tr.AllTitles)
	r.Delete("/", tr.DeleteTitle)
	r.Get("/{id}", tr.GetTitleByid)
	r.Get("/{categoryid}", tr.GetTitleBycategory_id)
	return r
}
func (tr BaseResource) CreateTitle(writer http.ResponseWriter, request *http.Request) {
	title := new(models.Title)
	if err := json.NewDecoder(request.Body).Decode(title); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(writer, "unknow err: %v", err)
		return
	}
	if err := tr.store.Title().Create(request.Context(), title); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	tr.broker.Cache().Purge()
	writer.WriteHeader(http.StatusCreated)
}
func (tr BaseResource) AllTitles(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.Filter{}
	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	titles, err := tr.store.Title().All(request.Context(), filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, titles)
}
func (tr BaseResource) GetTitleByid(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	title, err := tr.store.Title().ByID(request.Context(), id)
	if err != nil {
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	tr.broker.Cache().Add(id)
	render.JSON(writer, request, title)
}
func (tr BaseResource) GetTitleBycategory_id(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	title, err := tr.store.Title().ByCategoryId(request.Context(), id)
	if err != nil {
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	tr.broker.Cache().Add(id)
	render.JSON(writer, request, title)
}
func (tr BaseResource) DeleteTitle(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	eror := tr.store.Title().Delete(request.Context(), id)
	if eror != nil {
		fmt.Fprintf(writer, "Unknow db delete err : %s", err)
		return
	}
	tr.broker.Cache().Purge()
	render.JSON(writer, request, eror)
}
