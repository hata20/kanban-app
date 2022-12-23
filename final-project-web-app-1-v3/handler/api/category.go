package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {

	cekId := r.Context().Value("id").(string)
	if cekId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	intVar, _ := strconv.Atoi(cekId)
	result, err := c.categoryService.GetCategories(r.Context(), intVar)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	// result := map[string]interface{}{
	// 	"user_id":     categoryId,
	// 	"category_id": categoryId,
	// 	"message":     "success delete category",
	// }
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	if category.Type == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	id := r.Context().Value("id").(string)
	if id == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	intVar, _ := strconv.Atoi(id)
	result, err := c.categoryService.StoreCategory(r.Context(), &entity.Category{
		Type: category.Type, UserID: intVar,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	hasil := map[string]interface{}{
		"user_id":     result.UserID,
		"category_id": result.ID,
		"message":     "success create new category",
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(hasil)

}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	cekId := r.Context().Value("id").(string)
	categoryID := r.URL.Query().Get("category_id")

	cekIdInt, _ := strconv.Atoi(categoryID)

	err := c.categoryService.DeleteCategory(r.Context(), cekIdInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	cekIdfix, _ := strconv.Atoi(cekId)
	// categoryId := r.URL.Query().Get("category_id")
	result := map[string]interface{}{
		"user_id":     cekIdfix,
		"category_id": cekIdInt,
		"message":     "success delete category",
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
