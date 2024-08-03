package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-mail/mail"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var router *chi.Mux

type Order struct {
	Duration      string `json:"duration"`
	Date          string `json:"date"`
	About         string `json:"about"`
	NameEntity    string `json:"name_entity"`
	Email         string `json:"email"`
	StreetAddress string `json:"street_address"`
	PostalCode    string `json:"postal_code"`
}

func parseEmailAddresses(emails string) []string {
	return strings.Split(emails, ",")
}

func main() {
	godotenv.Load(".env")

	emails := os.Getenv("EMAIL_ADDRESSES")
	emailList := parseEmailAddresses(emails)
	log.Printf("EMAIL_ADDRESSES: %s", emails)
	log.Printf("PORT: %s", os.Getenv("PORT"))

	tmpl, err := template.ParseFiles("views/email.html")
	if err != nil {
		panic(err)
	}

	// Change to templates https://github.com/sgulics/go-chi-example/blob/master/cmd/server/main.go
	router = chi.NewRouter()
	router.Use(middleware.Logger)

	staticDir := http.Dir("static")
	FileServer(router, "/static", staticDir)

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
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

	router.Post("/order", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		new_order := Order{
			Duration:      r.FormValue("duration"),
			Date:          r.FormValue("date"),
			About:         r.FormValue("about"),
			NameEntity:    r.FormValue("name-entity"),
			Email:         r.FormValue("email"),
			StreetAddress: r.FormValue("street-address"),
			PostalCode:    r.FormValue("postal-code"),
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, new_order)
		if err != nil {
			panic(err)
		}
		htmlContent := buf.String()

		m := mail.NewMessage()
		m.SetHeader("From", "noreply@sign_wave_solutions.pt")
		m.SetHeader("To", emailList...)
		m.SetHeader("Subject", fmt.Sprintf("[%s] Pedido de Orçamento", new_order.NameEntity))
		m.SetBody("text/html", htmlContent)

		d := mail.NewDialer("smtp.gmail.com", 587, "ricardonunosr@gmail.com", "bbtx fqlb hntz rxly")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}

		w.Write([]byte("<p>Succeso</p>"))
	})

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
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
