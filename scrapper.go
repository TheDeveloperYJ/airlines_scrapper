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

//Status root object
type Status struct {
	data []map[string]string
}

//StatusHistory structure
type StatusHistory struct {
	trackingProcess string
	stationCode     string
	pieces          int
	weight          int
	actualTime      string
	flight          string
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
	// if doc.Contains(table) {
	// 	fmt.Println("Could not find any record")
	// }
	doc.Find("table").Each(func(index int, tabelHtml *goquery.Selection) {

		if index == 2 {
			tabelHtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
				headings = append(headings, strings.ToLower(tableheading.Text()))
			})

			tabelHtml.Find("tr").Each(func(indexth int, rowData *goquery.Selection) {
				if indexth > 0 {
					//var dataRow []StatusHistory
					//var dataRow []map[string]string
					dataRow := make(map[string]string)
					//var dataRow = make([]map[string]string, 0)
					for j := 0; j < rowData.Children().Length(); j++ {
						dataRow[headings[j]] = rowData.Children().Eq(j).Text()
					}
					finalData = append(finalData, dataRow)
				}

				//headings = append(headings, strings.ToLower(tableheading.Text()))
			})
			//fmt.Println(finalData)
			// m := Status{finalData}
			// b, err := json.Marshal(m)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// tabelHtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			// 	// fmt.Println(rowhtml.Children().Eq(0).Text())
			// 	// rowhtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
			// 	// 	headings = append(headings, strings.ToLower(tableheading.Text()))
			// 	// })
			// 	rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
			// 		fmt.Println(tablecell.Children().Text())
			// 		//fmt.Println(tablecell.Eq(indexth).Text())
			// 		// close := tablecell.Closest("tr").Text()
			// 		// //fmt.Println(close.Text())
			// 		// statusHistory[close] = append(statusHistory["close"], StatusHistory.)
			// 		// fmt.Println(statusHistory)
			// 		row = append(row, tablecell.Text())
			// 		//statusHistory['asd']
			// 	})
			// 	rows = append(rows, row)
			// 	row = nil
			// })

		}
	})

	//fmt.Println("####### rows = ", len(rows), rows)
	m := Status{finalData}
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
