package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var dev bool

var tmpls *template.Template

func init() {
	flag.BoolVar(&dev, "dev", false, "Development mode")
	flag.Parse()
	tmpls = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	if dev {
		log.Println("Running in development mode")
		f, err := os.Open("data/bowls.json")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		var bowlGames []bowlGame
		err = json.NewDecoder(f).Decode(&bowlGames)
		if err != nil {
			log.Fatal(err)
		}

		for _, game := range bowlGames {
			fmt.Printf("\"%s\",\n", game.Name)
		}
		os.Exit(0)
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))

	mux.HandleFunc("/", getMain)
	mux.HandleFunc("/selections", userSelections)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/favicon.ico", http.NotFound) // TODO: add a favicon

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
	ss := NewSelections(r.Form)

	submitter := r.Form.Get("submitter")
	// make the file of the selections with user name as the file name
	err = ss.MakeFile(submitter)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: add a template for this error
		return
	}

	submitData := map[string]interface{}{
		"submitter":  submitter,
		"selections": ss,
		"order":      SORTED_GAMES,
	}

	w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)
	if err := tmpls.ExecuteTemplate(w, "submit.gohtml", submitData); err != nil {
		log.Println(err)
	}
	return
}

func userSelections(w http.ResponseWriter, r *http.Request) {
	// grab query params from url
	qv := r.URL.Query()
	_, ok := qv["name"]
	// get the user selections
	us, err := NewUserSelections()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		// TODO: add a template for this error
		err = tmpls.ExecuteTemplate(w, "selections.gohtml", us)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if !ok {
		w.Header().Set("Content-Type", "text/html")
		err = tmpls.ExecuteTemplate(w, "selections.gohtml", struct {
			Us    UserSelections
			Dates []string
		}{
			us,
			DATES,
		})
		if err != nil {
			log.Println(err)
		}
		return
	} else {
		name := qv.Get("name")
		selections, ok := us[name]
		if !ok {
			log.Println("User not found")
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpls.ExecuteTemplate(w, "single_selector.gohtml", map[string]interface{}{
			"Name":       name,
			"Selections": selections,
		})
	}
}
