package main

import (
    "fmt"
    "net/http"
    "text/template"
    "github.com/gorilla/mux"
    "github.com/gorilla/schema"
    "github.com/gorilla/sessions"
    "database/sql"
    _ "github.com/lib/pq"
)

type Login struct {
    Username string
    Password string
}

type GuestbookRow struct {
    Id          int
    User        string
    Written_at  string
    Comment     string
}

type GuestbookRows struct {
    Rows []GuestbookRow
}

var db *sql.DB

var decoder       = schema.NewDecoder()
var store         = sessions.NewCookieStore([]byte("frase-ultra-secreta-de-encriptacao"))

var homeTemplate  = template.Must(template.New("home").ParseFiles("templates/home.html"))
var loginTemplate = template.Must(template.New("login").ParseFiles("templates/login.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "auth-session")

    data := map[string]interface{} {
        "logged_in": false,
        "username": "",
    }

    if username, ok := session.Values["username"]; ok {
        data["logged_in"] = true
        data["username"] = username
    }

    rows, err := db.Query("SELECT id, username, to_char(written_at, 'DD/MM/YYYY às HH24:MI'), comment FROM guestbook")

    if err != nil {
        fmt.Println(err)
    } else {
        guestbook_rows := new(GuestbookRows)

        for rows.Next() {
            var id int
            var user string
            var written_at string
            var comment string

            rows.Scan(&id, &user, &written_at, &comment)

            newRow := GuestbookRow{id,user,written_at,comment}
            guestbook_rows.Rows = append(guestbook_rows.Rows, newRow)
        }

        data["rows"] = guestbook_rows
    }

    if err := homeTemplate.ExecuteTemplate(w, "home", data); err != nil {
        fmt.Println(err)
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    data := make(map[string]interface{})

    data["wrong_password"] = len(query["wrong_password"]) > 0 && query["wrong_password"][0] == "1"

    if err := loginTemplate.ExecuteTemplate(w, "login", data); err != nil {
        fmt.Println(err)
    }
}

func DoLoginHandler(w http.ResponseWriter, r *http.Request) {
    var count int
    page := "/"
    login := new(Login)
    r.ParseForm();
    decoder.Decode(login, r.PostForm)

    db.QueryRow("SELECT COUNT(*) FROM users WHERE username=$1 AND password=$2", login.Username, login.Password).Scan(&count)

    if count == 0 {
        page = "/login?wrong_password=1"
    } else {
        session, _ := store.Get(r, "auth-session")
        session.Values["username"] = login.Username
        session.Save(r, w)
    }

    http.Redirect(w, r, page, http.StatusFound)
}

func GuestbookHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "auth-session")

    if username, ok := session.Values["username"]; ok {
        tx, _ := db.Begin()
        stmt, err := db.Prepare("INSERT INTO guestbook (username, comment) VALUES ( $1, $2 )")
        if err != nil {
            fmt.Println(err)
        } else {
            if comment := r.FormValue("comment"); comment != "" {
                stmt.Exec(username, comment)
                tx.Commit()
            } else {
                tx.Rollback()
            }
        }
    }

    // após gravar no guestbook, redirecione
    http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // após logout, redirecione
    http.SetCookie(w, &http.Cookie{Name: "auth-session", MaxAge: -1, Path: "/"})
    http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
    var err error
    db, err = sql.Open("postgres", "user=andre dbname=guestbook sslmode=disable")

    if err != nil {
        fmt.Println(err)
    }

    r := mux.NewRouter()

    r.HandleFunc("/",          HomeHandler)
    r.HandleFunc("/guestbook", GuestbookHandler)
    r.HandleFunc("/login",     LoginHandler)
    r.HandleFunc("/do_login",  DoLoginHandler)
    r.HandleFunc("/logout",    LogoutHandler)

    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
