package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/go-mail/mail"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var router *chi.Mux

func main() {
	godotenv.Load(".env")
	// Change to templates https://github.com/sgulics/go-chi-example/blob/master/cmd/server/main.go
	router = chi.NewRouter()
	router.Use(middleware.Logger)

	staticDir := http.Dir("static")
	FileServer(router, "/static", staticDir)

	// Define a route for the homepage
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Load and parse the HTML template
		tmpl, _ := template.ParseFiles("views/partials/base.html", "views/index.html")
		err := tmpl.Execute(w, nil)
		if err != nil {
			return
		}
	})

	router.Get("/order", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("views/partials/base.html", "views/order.html")
		err := tmpl.Execute(w, nil)
		if err != nil {
			return
		}
	})

	router.Get("/sidebar", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("views/sidebar.html")
		err := tmpl.Execute(w, nil)
		if err != nil {
			return
		}
	})

	router.Post("/order", func(w http.ResponseWriter, r *http.Request) {
		m := mail.NewMessage()

		m.SetHeader("From", "noreply@sign_wave_solutions.pt")

		m.SetHeader("To", "ricardonunosr@gmail.com", "noah.doe@example.com")

		m.SetHeader("Subject", "Pedido de Orcamento")

		m.SetBody("text/html", "Hello <b>Kate</b> and <i>Noah</i>!")

		d := mail.NewDialer("smtp.gmail.com", 587, "ricardonunosr@gmail.com", "bbtx fqlb hntz rxly")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}

	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
	if err != nil {
		return
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve static files.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
