package main

import (
	"github.com/deepilla/gokismet"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"time"
)

var config ConfigSpec

func checkComment() {
	comment := gokismet.Comment{
		UserIP:        "127.0.0.1",
		UserAgent:     "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.2) Gecko/20100115 Firefox/3.6",
		Page:          "https://ochronus.com/posts/feliz-cinco-de-mayo/",
		PageTimestamp: time.Date(2016, time.May, 5, 10, 30, 0, 0, time.UTC),
		//Author:        "A. Commenter",
		//AuthorEmail:   "buyviagra@gmail.com",
		Content: "How nice it is to celebrate",
		// etc...
	}
	akismetAPI := gokismet.NewAPI(config.AkismetKey, "ochronus.com")

	// Call the Check method, passing in your content.
	status, err := akismetAPI.CheckComment(comment.Values())

	// Handle the results.
	switch status {
	case gokismet.StatusHam:
		log.Println("This is legit content")
	case gokismet.StatusProbableSpam, gokismet.StatusDefiniteSpam:
		log.Println("This is spam")
	case gokismet.StatusUnknown:
		log.Println("Something went wrong:", err)
	}
}

func main() {

	err := envconfig.Process("comentador", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	checkComment()

	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, map[string]string{"hello": "json"})
	})

	n := negroni.New()
	recovery := negroni.NewRecovery()
	//recovery.PanicHandlerFunc = reportToSentry
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://foo.com"},
	})
	n.Use(corsMiddleware)
	n.Use(recovery)
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	http.ListenAndServe(":3003", n)
}
