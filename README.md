<h1 align="center">Lenslocked</h1>

## 📜 Summary
- [About](#About)
- [Libs/Dependencies](#Libs/Dependencies)
- [Setup](#Setup)
- [Run server](#Run-server)
- [Tests](#Tests)

---

<a id="About"></a> 
## 📃 About

This project was developed based on a web development course by Jon Calhoun (https://www.usegolang.com/). It is a web-based image management system, similar to a blog. Although my code differs from what was taught in the course, which primarily focused on MVC code design and lacked testing, my goal with this project was to deepen my understanding of the Go language. I aimed to explore the language's resources, libraries, and practices related to two subjects: Clean architecture and domain-driven design. This project is still on progress, so there's a few features that are missing.

---
<a id="Libs/Dependencies"></a> 
## 🗄 Libs/Dependencies </br>

| Name        | Description | Documentation | Installation |
| ----------- | ----------- | ------------- | ----------- |     
| pgx      | postgres database driver       |  github.com/jackc/pgx/v4 |  go get go get github.com/jackc/pgx/v4      |
| go-cmp   | Test library        | github.com/google/go-cmp     |   go get github.com/google/go-cmp          |
|  sqlite3  |   sqlite3 database driver     | github.com/mattn/go-sqlite3    |   go get github.com/mattn/go-sqlite3          |
|  bcrypt  |    Hash generator. Used for encrypt passwords   | golang.org/x/crypto/bcrypt  |   go get golang.org/x/crypto/bcrypt         |
| goose     | Database migrations      | github.com/pressly/goose/v3/cmd/goose@v3 | go install github.com/pressly/goose/v3/cmd/goose@v3      |   
| go-mail       | email sender library              |  github.com/go-mail/mail/v2 | go get github.com/go-mail/mail/v2     |
| gomock           | Mock library for tests.            | https://github.com/golang/mock                 | go get github.com/golang/mock     | 
| chi               |  http router  lib | https://github.com/go-chi/chi                   | go get github.com/go-chi/chi   |
| godotenv             | .env vars manager              | github.com/joho/godotenv             | go get github.com/joho/godotenv    | 
| gorilla/csrf         | middleware library that provides cross-site request forgery (CSRF)protection             | github.com/gorilla/csrf               | go get github.com/gorilla/csrf               | 
| google/uuid                 | uuid generator                   | github.com/google/uuid                        | go get github.com/google/uuid  |

---

<a id="Setup"></a> 
## 🔧 Setup

After cloning the project, inside the project there's an file called .env.example. You can create a new file called .env and copy and
paste all the content inside the .env.example file or rename it to .env. The .env file contains all the environment variables that is required
to the project. It's neccessary to put values in the variables so that the project can be used.

---
<a id="Run-server"></a> 
## ⚙️ Run

After configuring the .env file, inside the root folder, run one of the commands below:

```bash
    go run main.go 
```

```bash
    make run
```

Open your browser and type at the address bar: http://localhost:SERVER_PORT

*SERVER_PORT is the variable that was configured at the .env file

---
<a id="Tests"></a> 
## 🧪 Tests

Inside the root folder, run one of the commands below:

```bash
   make tests
```

```bash
    go test ./tests -v
```
