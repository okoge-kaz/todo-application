package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	database "github.com/okoge-kaz/todo-application/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// get login user
	userID := sessions.Default(ctx).Get("user_key")

	// Get query parameter
	keyword := ctx.Query("keyword")

	// Get tasks in DB
	var tasks []database.Task

	// rails active record のようには書けないらしい
	query := "SELECT  id, title, is_done, description, created_at FROM tasks INNER JOIN ownerships ON tasks.id = ownerships.task_id WHERE ownerships.user_id = ?"

	switch {
	case keyword != "":
		// キーワード検索
		err = db.Select(&tasks, query+" AND ( title LIKE ? OR description LIKE ? )", userID, "%"+keyword+"%", "%"+keyword+"%")
	default:
		// 全件取得
		// pixiv で ikumin さんに指摘されたが、SQLのパフォーマンス的にもこうするほうがよい
		err = db.Select(&tasks, query, userID)
	}
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks, "Keyword": keyword})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "task.html", task)
}

// form to create new task
func NewTaskForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form_new_task.html", gin.H{"Title": "New task registration"})
}

// create new task
func NewTask(ctx *gin.Context) {
	// get login user
	userID := sessions.Default(ctx).Get("user_key")

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get form data
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")

	// Insert a task with transaction
	transaction := db.MustBegin()
	// tasks table
	result, err := transaction.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", title, description)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	// ownerships table
	taskID, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	_, err = transaction.Exec("INSERT INTO ownerships (user_id, task_id) VALUES (?, ?)", userID, taskID)
	if err != nil {
		transaction.Rollback()
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	transaction.Commit()

	// Render status
	path := "/list" // デフォルトではタスク一覧ページへ戻る
	if id, err := result.LastInsertId(); err == nil {
		path = fmt.Sprintf("/task/%d", id) // 正常にIDを取得できた場合は /task/<id> へ戻る
	}
	ctx.Redirect(http.StatusFound, path)
}

// form to edit task
func EditTaskForm(ctx *gin.Context) {
	// ID の取得
	// /task/:id/edit
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	// Get target task
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Render edit form
	ctx.HTML(http.StatusOK, "form_edit_task.html",
		gin.H{"Title": fmt.Sprintf("Edit task %d", task.ID), "Task": task})
}

// edit task
func EditTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// /task/:id/edit
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Get form data
	title := ctx.PostForm("title")
	is_done := ctx.PostForm("is_done")
	description := ctx.PostForm("description")

	// 型変換
	is_done_bool, _ := strconv.ParseBool(is_done)

	// Update a task
	result, err := db.Exec("UPDATE tasks SET title=?, is_done=?, description=? WHERE id=?", title, is_done_bool, description, id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render status
	path := "/list" // デフォルトではタスク一覧ページへ戻る
	if rows, err := result.RowsAffected(); err == nil && rows == 1 {
		path = fmt.Sprintf("/task/%d", id) // 正常に1行更新できた場合は /task/<id> へ戻る
	}
	ctx.Redirect(http.StatusFound, path)
}

// delete task
func DeleteTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// /task/:id/delete
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Delete a task
	_, err = db.Exec("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render status
	path := "/list" // デフォルトではタスク一覧ページへ戻る
	ctx.Redirect(http.StatusFound, path)
}
