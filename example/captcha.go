package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/dim13/captcha"
)

var (
	tmpl    = template.Must(template.ParseFiles("captcha.tmpl"))
	private = flag.String("private", "none", "private key")
	public  = flag.String("public", "none", "public key")
	listen  = flag.String("listen", ":8000", "listen on")
)

type Page struct {
	Title  string
	Server string
	Public string
	Ok     string
	Error  string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c := captcha.New(*private, *public)
	p := &Page{
		Title:  "reCAPTCHA 1.0",
		Server: c.Server,
		Public: c.Public,
	}
	if r.Method == "POST" {
		if ok, err := c.Verify(r); ok {
			p.Ok = "Valid"
		} else {
			p.Error = err.Error()
		}
	}
	tmpl.ExecuteTemplate(w, "root", p)
}

func main() {
	flag.Parse()
	if *private == "none" || *public == "none" {
		flag.PrintDefaults()
		return
	}
	http.HandleFunc("/", rootHandler)
	log.Println("Listen on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
