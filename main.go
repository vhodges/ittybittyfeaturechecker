package main

import (
	"encoding/json"

	"fmt"
    "flag"

    "log"
    "net/http"
	"os"
	"strings"
	"time"

	"github.com/antonmedv/expr"
)

var Environment string
var ListendAddress string

func main() {

	flag.StringVar(&ListendAddress, "listen", ":8081", "Address the service should listen on")
	flag.StringVar(&Environment, "env", "production", "Environment we're running in and checking feature switches for")

	flag.Parse()

	log.Printf("The IttyBittyFeatureChecker 1.0\n")
	log.Printf("Listen address: %s\n", ListendAddress)
	log.Printf("Environment: %s\n", Environment)
	
	Features = make(map[string]Feature)

	dat, err := os.ReadFile("features.json")
	if err != nil {
		panic(err)
	}
	
	if err := json.Unmarshal(dat, &Features); err != nil {
        panic(err)
    }

	server()
}

type ProgramEnv map[string]interface{}

func (ProgramEnv) Now() time.Time {
	return time.Now()
}

func (ProgramEnv) ParseDate(d string) (time.Time, error) {
	return time.Parse(time.DateOnly, d)
}

func server() {

	http.HandleFunc("/enabled/", func(w http.ResponseWriter, r *http.Request) {		
		feature_id :=  strings.TrimPrefix(r.URL.Path, "/enabled/")
		
		counters.Incr(feature_id) // always increment the counter

		if feature, ok := Features[feature_id]; ok {
			if feature.Archived {
				fmt.Fprintf(w, "true")
			} else {

				env := ProgramEnv{
					"feature":   	feature,
					"count":     	counters.Get(feature_id),
					"params":    	r.URL.Query(),
					"environment":  Environment,
				}
		
				program, err := expr.Compile(feature.Expression, expr.Env(env))
				if err != nil {
					log.Printf("expression compile error %v\n", err)
					fmt.Fprintf(w, "false")
					return
				} 

				output, err := expr.Run(program, env)
				if err != nil {
					log.Printf("expression eval error %v\n", err)
					fmt.Fprintf(w, "false")
					return
				}

				if output.(bool) {
					fmt.Fprintf(w, "true")
				} else {
					fmt.Fprintf(w, "false")
				}
			}
		} else {
			fmt.Fprintf(w, "false")
		}

    })

    http.HandleFunc("/features/", func(w http.ResponseWriter, r *http.Request){
		json, err := os.ReadFile("features.json")
		if err != nil {
			fmt.Fprintf(w, "Error getting features\n")
		}
		_, _ = w.Write(json)
    })

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "These droids are not the ones you're looking... move along")
    })

    log.Fatal(http.ListenAndServe(ListendAddress, nil))
}