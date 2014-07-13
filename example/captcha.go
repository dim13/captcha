package main

import (
	"flag"
	"github.com/dim13/captcha"
	"html/template"
	"log"
	"net/http"
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
		Title:  "reCAPTCHA test",
		Server: cc.Server,
		Public: cc.Public,
	}

	challenge := r.PostFormValue("recaptcha_challenge_field")
	response := r.PostFormValue("recaptcha_response_field")
	if challenge != "" && response != "" {
		if ok, err := cc.Verify(r.RemoteAddr, challenge, response); ok {
			p.Ok = "valid"
		} else {
			p.Error = err.Error()
		}
	}

	t := template.Must(template.ParseFiles("captcha.tmpl"))
	t.Execute(w, p)
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
