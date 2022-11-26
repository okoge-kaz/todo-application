package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/okoge-kaz/todo-application/controllers"
)

func Init() *gin.Engine {
	// initialize engine
	engine := gin.Default()

	// set template directory
	engine.LoadHTMLGlob("views/*.html")

	// initialize session
	store := cookie.NewStore([]byte("my-secret"))
	engine.Use(sessions.Sessions("user-session", store))

	// import static files
	engine.Static("/assets", "./assets")

	// routing
	engine.GET("/", controllers.Home)
	engine.GET("/list", controllers.LoginCheck, controllers.TaskList)

	taskGroup := engine.Group("/task")
	taskGroup.Use(controllers.LoginCheck)

	// Grouping /task/xxx
	{
		// Create, Update, Delete
		taskGroup.GET("/new", controllers.NewTaskForm)
		taskGroup.POST("/new", controllers.NewTask)

		taskGroup.GET("/:id", controllers.TaskAccessCheck, controllers.ShowTask) // ":id" is a parameter
		//:id
		taskIDGroup := taskGroup.Group("/:id")
		taskIDGroup.Use(controllers.TaskAccessCheck)
		{
			taskIDGroup.GET("/edit", controllers.EditTaskForm)
			taskIDGroup.POST("/edit", controllers.EditTask)
			taskIDGroup.GET("/delete", controllers.DeleteTask)
		}
	}

	// user registration
	engine.GET("/user/new", controllers.NewUserForm)
	engine.POST("/user/new", controllers.RegisterUser)

	// logged in user
	userGroup := engine.Group("/user")
	userGroup.Use(controllers.LoginCheck)
	{
		// change password
		userGroup.GET("/change_password", controllers.ChangePasswordForm)
		userGroup.POST("/change_password", controllers.ChangeUserInfo)
		// delete user
		userGroup.GET("/delete", controllers.DeleteUser)
	}
	// login
	engine.GET("/login", controllers.LoginForm)
	engine.POST("/login", controllers.Login)
	// logout
	engine.GET("/logout", controllers.LoginCheck, controllers.Logout)

	return engine
}
