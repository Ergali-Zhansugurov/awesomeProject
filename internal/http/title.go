package http

import (
	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/models"
	"awesomeProject/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
)

type TitleResourse struct {
	store  store.Store
	broker broker_models.Broker
	cache  *lru.TwoQueueCache
}

func NewTitleResource(store store.Store, cache *lru.TwoQueueCache, broker broker_models.Broker) *TitleResourse {
	return &TitleResourse{store: store, broker: broker, cache: cache}
}
func (tr TitleResourse) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", tr.CreateTitle)
	r.Get("/", tr.AllTitles)
	r.Delete("/", tr.DeleteTitle)
	r.Get("/", tr.GetByid)
	return r
}
func (tr TitleResourse) CreateTitle(writer http.ResponseWriter, request *http.Request) {
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
func (cr TitleResourse) AllTitles(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.Titlesfilter{}
	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	titles, err := cr.store.Title().All(request.Context(), filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, titles)
}
func (tr TitleResourse) GetByid(writer http.ResponseWriter, request *http.Request) {
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
func (tr TitleResourse) GetBycategory_id(writer http.ResponseWriter, request *http.Request) {
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
func (tr TitleResourse) DeleteTitle(writer http.ResponseWriter, request *http.Request) {
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
	tr.broker.Cache().Remove(id)
	render.JSON(writer, request, eror)
}
