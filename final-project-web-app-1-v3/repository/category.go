package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	hasilCategory := []entity.Category{}
	hasil, err := r.db.WithContext(ctx).
		Table("categories").
		Where("user_id= ?", id).
		Select("*").Rows()
	if err != nil {
		return []entity.Category{}, err
	}
	defer hasil.Close()
	for hasil.Next() {
		r.db.ScanRows(hasil, &hasilCategory)
	}
	return hasilCategory, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	r.db.WithContext(ctx)
	result := r.db.Create(&category)
	if result.Error != nil {
		return 0, result.Error
	}
	categoryId = category.ID
	return categoryId, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	r.db.WithContext(ctx)
	result := r.db.Create(&categories)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	var hasilCategory entity.Category
	err := r.db.WithContext(ctx).Model(&entity.Category{}).Where("id= ?", id).Take(&hasilCategory)
	if err.Error != nil {
		return entity.Category{}, err.Error
	}
	return hasilCategory, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	r.db.WithContext(ctx)
	err := r.db.Where("id= ?", category.ID).Updates(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	r.db.WithContext(ctx)
	err := r.db.Where("id = ?", id).Delete(&entity.Category{})
	if err.Error != nil {
		return err.Error
	}
	return nil
}
