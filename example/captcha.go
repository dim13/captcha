package main

import (
	"flag"
	"github.com/dim13/captcha"
	"html/template"
	"log"
	"net/http"
)

var (
	private = flag.String("private", "", "private key")
	public  = flag.String("public", "", "public key")
	listen = flag.String("listen", ":8000", "listen on")
)

type Page struct {
	Title  string
	Server string
	Public string
	Ok     string
	Error  string
}

var cc *captcha.Captcha

func root(w http.ResponseWriter, r *http.Request) {
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
	flag.Parse()
	if *private == "" || *public == "" {
		flag.PrintDefaults()
		return
	}
	cc = captcha.New(*private, *public)
	http.HandleFunc("/", root)
	log.Println("Listen on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
