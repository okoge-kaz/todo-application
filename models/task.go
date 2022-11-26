package models

import (
	"github.com/jmoiron/sqlx"
	database "github.com/okoge-kaz/todo-application/db"
)

func GetTaskByUserIDAndKeyword(userID int, keyword string) ([]database.Task, error) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	var tasks []database.Task
	// keyword search (title, description)
	err = db.Select(&tasks, "SELECT  id, title, is_done, description, status, due_date, priority FROM tasks INNER JOIN ownerships ON tasks.id = ownerships.task_id WHERE ownerships.user_id = ? AND (title LIKE ? OR description LIKE ?) ORDER BY due_date ASC, priority DESC", userID, "%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByUserID(userID int) ([]database.Task, error) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	var tasks []database.Task
	err = db.Select(&tasks, "SELECT  id, title, is_done, description, status, due_date, priority FROM tasks INNER JOIN ownerships ON tasks.id = ownerships.task_id WHERE ownerships.user_id = ? ORDER BY due_date ASC, priority DESC", userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByUserIDAndStatus(userID int, status []string) ([]database.Task, error) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	sql, params, err := sqlx.In("SELECT  id, title, is_done, description, status, due_date, priority FROM tasks INNER JOIN ownerships ON tasks.id = ownerships.task_id WHERE ownerships.user_id = ? AND status IN (?) ORDER BY due_date ASC, priority DESC", userID, status)
	if err != nil {
		return nil, err
	}

	var tasks []database.Task
	err = db.Select(&tasks, sql, params...)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByUserIDAndKeywordAndStatus(userID int, keyword string, status []string) ([]database.Task, error) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	sql, params, err := sqlx.In("SELECT  id, title, is_done, description, status, due_date, priority FROM tasks INNER JOIN ownerships ON tasks.id = ownerships.task_id WHERE ownerships.user_id = ? AND (title LIKE ? OR description LIKE ?) AND status IN (?) ORDER BY due_date ASC, priority DESC", userID, "%"+keyword+"%", "%"+keyword+"%", status)
	if err != nil {
		return nil, err
	}

	var tasks []database.Task
	err = db.Select(&tasks, sql, params...)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
