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
	private = flag.String("private", "", "private key")
	public  = flag.String("public", "", "public key")
	listen  = flag.String("listen", ":8000", "listen on")
)

type Page struct {
	captcha.Captcha
	Result string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{
		Captcha: captcha.New(*private, *public),
	}
	if r.Method == "POST" {
		if ok, err := p.Captcha.Verify(r); ok {
			p.Result = "Ok"
		} else {
			p.Result = err.Error()
		}
	}
	tmpl.ExecuteTemplate(w, "root", p)
}

func main() {
	flag.Parse()
	if *private == "" || *public == "" {
		flag.PrintDefaults()
		return
	}
	http.HandleFunc("/", rootHandler)
	log.Println("Listen on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
