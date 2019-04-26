// make_http_request_with_timeout.go
//235-31634525
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

//TrackStatus root object
type TrackStatus struct {
	Status string              `json:"status"`
	Data   []map[string]string `json:"data"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	awbNumber := vars["awbNumber"]
	turkishCargoURL := fmt.Sprintf("%s%s", "https://www.turkishcargo.com.tr/en/online-services/shipment-tracking?quick=True&awbInput=", awbNumber)
	getDetails(turkishCargoURL, w, r)

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/{awbNumber}", getStatus)

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func main() {
	fmt.Println("Server started...")
	handleRequests()
}

func getDetails(turkishCargoURL string, w http.ResponseWriter, r *http.Request) {
	var headings []string
	var finalData []map[string]string

	client := &http.Client{
		Timeout: 300 * time.Second,
	}

	// Make requests
	response, err := client.Get(turkishCargoURL)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}
	if doc.Find("table").Length() <= 0 {
		fmt.Println("Error")
		m := TrackStatus{"error", finalData}
		b, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	} else {
		doc.Find("table").Each(func(index int, tabelHtml *goquery.Selection) {

			if index == 2 {
				tabelHtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
					headings = append(headings, strings.ToLower(tableheading.Text()))
				})

				tabelHtml.Find("tr").Each(func(indexth int, rowData *goquery.Selection) {
					if indexth > 0 {
						dataRow := make(map[string]string)
						for j := 0; j < rowData.Children().Length(); j++ {
							dataRow[headings[j]] = rowData.Children().Eq(j).Text()
						}
						finalData = append(finalData, dataRow)
					}

				})
				m := TrackStatus{"success", finalData}

				b, err := json.Marshal(m)
				if err != nil {
					log.Fatal(err)
				}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		})

	}

}
