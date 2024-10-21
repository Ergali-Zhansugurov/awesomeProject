package http

import (
	"encoding/json"
	"fmt"

	"awesomeProject/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (Ur BaseResource) UserRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", Ur.CreateUser)
	r.Delete("/", Ur.DeleteUser)
	r.Post("/users/login", Ur.UserAuthenticate)
	r.Post("/users/{id}/titles", Ur.UserAddTitleToLibrary)
	r.Delete("/users/{id}/titles/{titleID}", Ur.UserRemoveTitleFromLibrary)
	return r
}
func (ur *BaseResource) UserAuthenticate(w http.ResponseWriter, r *http.Request) {
	var login models.LoginID
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ur.store.User().Authenticate(r.Context(), login.Username, login.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Вернуть токен или пользовательские данные
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (ur *BaseResource) UserAddTitleToLibrary(writer http.ResponseWriter, request *http.Request) {
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

func (ur *BaseResource) UserRemoveTitleFromLibrary(writer http.ResponseWriter, request *http.Request) {
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
func (Ur BaseResource) CreateUser(writer http.ResponseWriter, request *http.Request) {
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
func (Ur BaseResource) GetUsers(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.Filter{}
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
func (Ur BaseResource) GetByid(writer http.ResponseWriter, request *http.Request) {
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
func (Ur BaseResource) DeleteUser(writer http.ResponseWriter, request *http.Request) {
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
