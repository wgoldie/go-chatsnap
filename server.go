package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"github.com/bndr/gopencils"
)

type ImageManager struct {
	Api	*gopencils.Resource	
}

type Message struct {
    Message string `json:"message"`
    Handle string `json:"handle"`
    Channel string `json:"channel"`
}

func send(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var m Message
	err := decoder.Decode(&m)
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Println(m.Handle + " " + m.Message + " " + m.Channel)
}

func main() {
	im := ImageManager{Api: gopencils.Api("https://api.datamarket.azure.com/Bing/Search/v1/")}

	fmt.Println("API root:" + im.Api.Url)
	
        http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/send", send)
        http.ListenAndServe(":3000", nil)
}
