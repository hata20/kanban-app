package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	result := []entity.Task{}
	hasil, err := r.db.WithContext(ctx).Table("tasks").Where("user_id= ?", id).Select("*").Rows()
	if err != nil {
		return []entity.Task{}, err
	}
	defer hasil.Close()
	for hasil.Next() {
		r.db.ScanRows(hasil, &result)
	}
	return result, nil

}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	r.db.WithContext(ctx)
	res := r.db.Create(&task)
	if res.Error != nil {
		return 0, res.Error
	}
	taskId = task.ID
	return taskId, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	var result entity.Task
	r.db.WithContext(ctx)
	err := r.db.Model(&entity.Task{}).Where("id= ?", id).Take(&result)
	if err.Error != nil {
		return entity.Task{}, err.Error
	}
	return result, nil
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	var result []entity.Task
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("category_id= ?", catId).Find(&result)
	if err.Error != nil {
		return []entity.Task{}, err.Error
	}
	return result, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	r.db.WithContext(ctx)
	err := r.db.Model(&entity.Task{}).Where("id= ?", task.ID).Updates(task).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	r.db.WithContext(ctx)
	err := r.db.Delete(&entity.Task{}, id)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
