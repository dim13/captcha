package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/dim13/captcha"
)

var (
	private string
	public  string
	listen  string
)

type Page struct {
	Title  string
	Server string
	Public string
	Ok     string
	Error  string
}

func init() {
	flag.StringVar(&private, "private", "none", "private key")
	flag.StringVar(&public, "public", "none", "public key")
	flag.StringVar(&listen, "listen", ":8000", "listen on")
	flag.Parse()
}

var cc captcha.Captcha

func rootHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		Title:  "reCAPTCHA 1.0",
		Server: cc.Server,
		Public: cc.Public,
	}

	if ok, err := cc.Verify(r); ok {
		p.Ok = "Valid"
	} else {
		p.Error = err.Error()
	}

	t := template.Must(template.ParseFiles("captcha.tmpl"))
	t.ExecuteTemplate(w, "root", p)
}

func main() {
	if private == "none" || public == "none" {
		flag.PrintDefaults()
		return
	}
	cc = captcha.New(private, public)
	http.HandleFunc("/", rootHandler)
	log.Println("Listen on", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
