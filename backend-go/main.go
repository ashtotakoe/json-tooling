package main

import (
	parser "back/parser2"
	"fmt"
	"net/http"
)

func setCors(w http.ResponseWriter) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

}


func main() {
	// http.HandleFunc("POST /parse", func(w http.ResponseWriter, r *http.Request) {
	// 	setCors(w)

	// 	data, err := io.ReadAll(r.Body)

	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 	}

	// 	w.WriteHeader(http.StatusOK)

	// 	err = parser.Parse(string(data))
	// 	if err != nil {
	// 		fmt.Fprintf(w, "parsing failed. Error: %s", err.Error())
	// 	}

	// 	w.Write([]byte("parsing success"))
	// })

	// http.ListenAndServe(":3000", nil)

	// err := parser.Parse(`
	// 	{
	// 		"aaao": true,
	// 		"bb:b": "uw:w\"au:",
	// 		"u:u": 223
	// 	}
	// 	`)
	//


	err := parser.Parse((`


		{
			"aaao": true,
			"bbb": 5,
			"bro": [
				"test",
				"bimbo",
				{
						"aaa": "uuu"
				}
			]
		}




		`))
	fmt.Println(err)
}
