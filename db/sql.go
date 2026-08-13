package db

import (
	"sync"
	"test/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	db    *sqlx.DB
	mutex sync.Mutex
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		db: db,
	}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) CreateTask(task *entity.Task) (entity.Task, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	task.CreatedAt = time.Now().Format(time.DateTime)
	task.UpdatedAt = time.Now().Format(time.DateTime)

	result, err := db.db.Exec(`INSERT INTO tasks (name, status, created_at, updated_at) VALUES (?, ?, ?, ?)`, task.Name, task.Status, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return entity.Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entity.Task{}, err
	}

	task.ID = id
	return *task, nil
}

func (db *DB) GetTask(id int64) (*entity.Task, error) {
	var task entity.Task
	err := db.db.Get(&task, `SELECT id, name, status, created_at, updated_at FROM tasks WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (db *DB) GetTasks() (entity.TaskList, error) {
	var tasks entity.TaskList
	err := db.db.Select(&tasks, `SELECT id, name, status, created_at, updated_at FROM tasks`)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (db *DB) UpdateTask(task *entity.Task) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	updatedAt := time.Now().Format(time.DateTime)

	_, err := db.db.Exec(`UPDATE tasks SET updated_at = ? WHERE id = ?`, updatedAt, task.ID)
	return err
}

func (db *DB) DeleteTask(id int64) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	_, err := db.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}
