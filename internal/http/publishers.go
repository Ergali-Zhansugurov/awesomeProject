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

type PublisherResourse struct {
	store  store.Store
	broker broker_models.Broker
	cache  *lru.TwoQueueCache
}

func NewPublisherResourse(store store.Store, cache *lru.TwoQueueCache, broker broker_models.Broker) *PublisherResourse {
	return &PublisherResourse{store: store, broker: broker, cache: cache}
}
func (Pr PublisherResourse) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", Pr.CreatePublisher)
	r.Get("/", Pr.GetPublishers)
	r.Get("/", Pr.GetByid)
	r.Delete("/", Pr.DeletePublisher)
	return r
}
func (Pr PublisherResourse) CreatePublisher(writer http.ResponseWriter, request *http.Request) {
	Publisher := new(models.Publisher)
	if err := json.NewDecoder(request.Body).Decode(Publisher); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(writer, "unknow err: %v", err)
		return
	}
	if err := Pr.store.Publisher().Create(request.Context(), Publisher); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	Pr.broker.Cache().Purge()
	writer.WriteHeader(http.StatusCreated)
}
func (Pr PublisherResourse) GetPublishers(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.Publisherfilter{}
	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	Publisher, err := Pr.store.Publisher().Get(request.Context(), filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, Publisher)
}
func (Pr PublisherResourse) GetByid(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	Publisher, err := Pr.store.Publisher().ByID(request.Context(), id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, Publisher)
}
func (Pr PublisherResourse) DeletePublisher(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	eror := Pr.store.Publisher().Delete(request.Context(), id)
	if eror != nil {
		fmt.Fprintf(writer, "Unknow db delete err : %s", err)
		return
	}
	Pr.broker.Cache().Remove(id)
	render.JSON(writer, request, eror)
}
