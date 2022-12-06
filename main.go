package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

var dev bool

var tmpls *template.Template

func init() {
	flag.BoolVar(&dev, "dev", false, "Development mode")
	tmpls = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))

	mux.HandleFunc("/", getMain)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Starting server on port 1313")
	if err := http.ListenAndServe(":1313", mux); err != nil {
		log.Fatal(err)
	}
}

func getMain(w http.ResponseWriter, r *http.Request) {
	// Get bowlGames
	bowlGames, err := getBowlsData()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html")

		w.WriteHeader(http.StatusOK)
		log.Println("Rendering the main page")
		if err := tmpls.ExecuteTemplate(w, "index.gohtml", bowlGames); err != nil {
			log.Println(err)
		}
		return
	}
	// POST method
	// parse form
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: add a template for this error
		return
	}

	// create a Selections Type from the form data
	ss := NewSelctions(r.Form)

	// make the file of the selections with user name as the file name
	err = ss.MakeFile(r.Form.Get("submitter"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: add a template for this error
		return
	}

	w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)
	log.Println("Rendering the main page after POST")
	if err := tmpls.ExecuteTemplate(w, "submit.gohtml", r.Form); err != nil {
		log.Println(err)
	}
	return
}
