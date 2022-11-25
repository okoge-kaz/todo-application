# Todolist

This repository provides a base project to implemet Todolist application.

`docker-compose.yml` provides Go 1.17 build tool, MySQL server and phpMyAdmin.

## Dependencies

- [Gin Web Framework](https://pkg.go.dev/github.com/gin-gonic/gin)
- [Sqlx](https://pkg.go.dev/github.com/jmoiron/sqlx)

## How to run the application

First, you need to start Docker containers.

```sh
$ docker-compose up -d
```

When you finish exercise, please don't forget to stop the containers.

```sh
$ docker-compose down
```

## Advanced: How to initialize database

When you modify the database schema, you will need to discard the current DB volumes for creating a new one.
It will be easier to rebuild everything than to rebuild only DB container.
Following command helps you to do it.

```sh
$ docker-compose down --rmi all --volumes --remove-orphans
$ docker-compose up -d
```

## sqlx

Official site: https://jmoiron.github.io/sqlx/

## Gin

### URL parameters

- Query Parameters: `/user?id=123`

  ```go
  router.GET("/user", func(ctx *gin.Context) {
    id := ctx.Query("id")
  })
  ```

- Path Parameters: `/user/123`

  ```go
  router.GET("/user/:id", func(ctx *gin.Context) {
    id := ctx.Param("id")
  })
  ```

  ちなみに以下のような場合でも path parameter は正しく動作する

  `tasks/:id/edit`

  ```go
  router.GET("/tasks/:id/edit", func(ctx *gin.Context) {
    id := ctx.Param("id")
  })
  ```

## Directory structure

```
.
├── Dockerfile
├── Makefile
├── README.md
├── assets
│   ├── script.js
│   └── style.css
├── db
│   ├── conn.go
│   └── entity.go
├── docker
│   └── db
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── service
│   ├── default.go
│   ├── task.go
│   └── user.go
├── tmp
│   ├── build-errors.log
│   └── main
└── views
    ├── _footer.html
    ├── _header.html
    ├── change_password_form.html
    ├── error.html
    ├── form_edit_task.html
    ├── form_new_task.html
    ├── index.html
    ├── new_user_form.html
    ├── task.html
    └── task_list.html
```
