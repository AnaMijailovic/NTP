package net

import (
	"encoding/json"
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/service"
	"net/http"
)

func Serve() {

	http.HandleFunc("/api/fileTree", GetFileTree)
	http.HandleFunc("/api/fileTypeData", GetTypeChartData)

	fmt.Println("Starting server ...")
	http.ListenAndServe(":8080", nil)

}

func GetFileTree(w http.ResponseWriter, r *http.Request) {

	/*vars := mux.Vars(r)
	path := vars["path"]
	fmt.Println("Path: ", path)*/

	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	tree := service.CreateTree(path)

	var err = json.NewEncoder(w).Encode(tree)

	if err != nil {
		fmt.Println("Some error happened: ", err)
	}
}

func GetTypeChartData(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	data := service.FileTypeChartData(path)
	var err = json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Println("Some error happened: ", err)
	}
}
