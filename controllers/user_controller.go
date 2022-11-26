package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	database "github.com/okoge-kaz/todo-application/db"
)

func NewUserForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new_user_form.html", gin.H{"Title": "Register user"})
}

func RegisterUser(ctx *gin.Context) {
	// フォームデータの受け取り
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Username or Password is not provided", "Username": username, "Password": password})
		return
	}

	if len(password) < 8 {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password must be at least 8 characters", "Username": username, "Password": password})
		return
	}

	// password double check
	password2 := ctx.PostForm("password2")
	if password != password2 {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password and Password confirmation are not same"})
		return
	}

	// DB 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// 重複チェック
	var duplicate int
	err = db.Get(&duplicate, "SELECT COUNT(*) FROM users WHERE name=?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if duplicate > 0 { // count の結果が返却されるので
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Username is already taken", "Username": username, "Password": password})
		return
	}

	// DB への保存
	result, err := db.Exec("INSERT INTO users(name, password) VALUES (?, ?)", username, hash(password))
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// 保存状態の確認
	id, _ := result.LastInsertId()
	var user database.User
	err = db.Get(&user, "SELECT id, name, password FROM users WHERE id = ?", id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// change password
func ChangePasswordForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "change_password_form.html", gin.H{"Title": "Change password"})
}

func ChangeUserInfo(ctx *gin.Context) {
	// フォームデータの受け取り
	username := ctx.PostForm("username")
	newUsername := ctx.PostForm("new_username")
	oldPassword := ctx.PostForm("old_password")
	newPassword := ctx.PostForm("new_password")

	if newPassword == "" {
		ctx.HTML(http.StatusBadRequest, "change_password_form.html", gin.H{"Title": "Change password", "Error": "New password is not provided"})
		return
	}

	if len(newPassword) < 8 {
		ctx.HTML(http.StatusBadRequest, "change_password_form.html", gin.H{"Title": "Change password", "Error": "New password must be at least 8 characters"})
		return
	}

	// db 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// old password is correct?
	var user database.User
	err = db.Get(&user, "SELECT id, name, password FROM users WHERE name = ?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if user.Password == nil || string(user.Password) != string(hash(oldPassword)) {
		ctx.HTML(http.StatusBadRequest, "change_password_form.html", gin.H{"Title": "Change password", "Error": "Old password is incorrect"})
		return
	}

	// update username and password
	_, err = db.Exec("UPDATE users SET name=?, password=? WHERE id=?", newUsername, hash(newPassword), user.ID)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// login form
func LoginForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Login"})
}

// login
const userKey = "user_key"

func Login(ctx *gin.Context) {
	// フォームデータの受け取り
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	// db 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// password is correct?
	var user database.User
	err = db.Get(&user, "SELECT id, name, password FROM users WHERE name = ?", username)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Title": "Login", "Username": username, "Error": "No such user"})
		return
	}
	if hex.EncodeToString(user.Password) != hex.EncodeToString(hash(password)) {
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Title": "Login", "Username": username, "Error": "Username or Password is incorrect"})
		return
	}

	// session に user を保存
	session := sessions.Default(ctx)
	session.Set(userKey, user.ID)
	session.Save()

	ctx.Redirect(http.StatusFound, "/list")
}

// user login check  (auth middleware)
func LoginCheck(ctx *gin.Context) {
	if sessions.Default(ctx).Get(userKey) == nil {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
	} else {
		ctx.Next()
	}
}

// user can access this task?
func TaskAccessCheck(ctx *gin.Context) {
	// task id を取得
	taskID := ctx.Param("id")

	if isValidUser(ctx, taskID) {
		ctx.Next()
	} else {
		ctx.Redirect(http.StatusFound, "/list") // エラーメッセージが欲しい
		ctx.Abort()
	}
}

// logout
func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	ctx.Redirect(http.StatusFound, "/")
}

// delete user
func DeleteUser(ctx *gin.Context) {
	// db 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// session から user id を取得
	session := sessions.Default(ctx)
	userID := session.Get(userKey)

	// user を削除
	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// session を削除
	// user の削除はできているが、session が残っている
	session = sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()

	ctx.HTML(http.StatusOK, "delete_user.html", gin.H{"Title": "Delete user"})
}

// private
func hash(password string) []byte {
	const salt = "todolist.go#"
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(password))
	return h.Sum(nil)
}

// is valid user?
func isValidUser(ctx *gin.Context, taskID string) bool {
	// session から user を取得
	session := sessions.Default(ctx)
	userID := session.Get(userKey)

	// db 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return false
	}

	// user が task を持っているか
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM tasks INNER JOIN ownerships ON tasks.id = ownerships.task_id WHERE tasks.id = ? AND ownerships.user_id = ?", taskID, userID)

	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
