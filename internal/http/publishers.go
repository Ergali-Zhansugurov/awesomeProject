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

func (Pr BaseResource) PublisherRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", Pr.CreatePublisher)
	r.Delete("/", Pr.DeletePublisher)
	r.Post("/users/login", Pr.PublisherAuthenticate)
	r.Post("/users/{id}/titles", Pr.PublisherAddTitleToLibrary)
	r.Delete("/users/{id}/titles/{titleID}", Pr.PublisherRemoveTitleFromLibrary)
	return r
}
func (ur *BaseResource) PublisherAuthenticate(w http.ResponseWriter, r *http.Request) {
	var login models.LoginID
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ur.store.Publisher().Authenticate(r.Context(), login.Username, login.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Вернуть токен или пользовательские данные
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (ur *BaseResource) PublisherAddTitleToLibrary(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	title, err := ur.store.Title().ByID(request.Context(), id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	useridStr := chi.URLParam(request, "userid")
	userid, err := strconv.Atoi(useridStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	er := ur.store.User().AddTitleToLibrary(request.Context(), userid, title)
	if er != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	render.JSON(writer, request, er)
}

func (ur *BaseResource) PublisherRemoveTitleFromLibrary(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	useridStr := chi.URLParam(request, "userid")
	userid, err := strconv.Atoi(useridStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	er := ur.store.User().RemoveTitleFromLibrary(request.Context(), userid, id)
	if er != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "unknow err: %s", err)
		return
	}
	render.JSON(writer, request, er)
}
func (Pr BaseResource) CreatePublisher(writer http.ResponseWriter, request *http.Request) {
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
func (Pr BaseResource) GetPublishers(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.Filter{}
	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	publisher, err := Pr.store.Publisher().Get(request.Context(), filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "DB err : %v", err)
		return
	}
	render.JSON(writer, request, publisher)
}
func (Pr BaseResource) GetPublisherByid(writer http.ResponseWriter, request *http.Request) {
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
func (Pr BaseResource) DeletePublisher(writer http.ResponseWriter, request *http.Request) {
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
	Pr.broker.Cache().Purge()
	render.JSON(writer, request, eror)
}
