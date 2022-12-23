package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Password == "" || user.Email == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}

	res, err := u.userService.Login(r.Context(), &entity.User{
		Password: user.Password,
		Email:    user.Email,
	})

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		json.NewEncoder(w).
			Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	responLoginUser := map[string]interface{}{
		"message": "login success",
		"user_id": res,
	}
	resInteger := strconv.Itoa(res)
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   resInteger,
		Path:    "/",
		Expires: time.Now().Add(5 * time.Hour),
	})
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responLoginUser)
}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Fullname == "" || user.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	hasil, err := u.userService.Register(r.Context(), &entity.User{
		Fullname: user.Fullname,
		Password: user.Password,
		Email:    user.Email,
	})

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	responRegister := map[string]interface{}{
		"message": "register success",
		"user_id": hasil.ID,
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(responRegister)
}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "logout success"})
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	idUserLogin := r.URL.Query().Get("user_id")

	if idUserLogin == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	hapusUserID, _ := strconv.Atoi(idUserLogin)

	err := u.userService.Delete(r.Context(), int(hapusUserID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
