package http

import (
	"encoding/json"
	"fmt"

	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/models"
	"awesomeProject/internal/store"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
)

type UserResourse struct {
	store  store.Store
	broker broker_models.Broker
	cache  *lru.TwoQueueCache
}

func NewUserResourse(store store.Store, cache *lru.TwoQueueCache, broker broker_models.Broker) *UserResourse {
	return &UserResourse{store: store, broker: broker, cache: cache}
}
func (Ur UserResourse) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", Ur.CreateUser)
	r.Get("/", Ur.GetUsers)
	r.Get("/", Ur.GetByid)
	r.Delete("/", Ur.DeleteUser)
	return r
}
func (Ur UserResourse) CreateUser(writer http.ResponseWriter, request *http.Request) {
	User := new(models.User)
	if err := json.NewDecoder(request.Body).Decode(User); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(writer, "unknow err: %v", err)
		return
	}
	if err := Ur.store.User().Create(request.Context(), User); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}
func (Ur UserResourse) GetUsers(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.UserFilter{}
	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	titles, err := Ur.store.User().Get(request.Context(), filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, titles)
}
func (Ur UserResourse) GetByid(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	User, err := Ur.store.User().ByID(request.Context(), id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, User)
}
func (Ur UserResourse) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	eror := Ur.store.User().Delete(request.Context(), id)
	if eror != nil {
		fmt.Fprintf(writer, "Unknow db delete err : %s", err)
		return
	}
	render.JSON(writer, request, eror)
}
