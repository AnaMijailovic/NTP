package net

import (
	"encoding/json"
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"github.com/AnaMijailovic/NTP/arf/service"
	"net/http"
	"strconv"
	"time"
)

func Serve() {

	http.HandleFunc("/api/fileTree", GetFileTree)
	http.HandleFunc("/api/fileTypeData", GetTypeChartData)
	http.HandleFunc("/api/deleteFiles", DeleteFiles)

	fmt.Println("Starting server ...")
	http.ListenAndServe(":8080", nil)

}

func GetFileTree(w http.ResponseWriter, r *http.Request) {

	/*vars := mux.Vars(r)
	path := vars["path"]
	fmt.Println("Path: ", path)*/

	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	tree := service.CreateTree(path, true)

	var err = json.NewEncoder(w).Encode(tree)

	if err != nil {
		fmt.Println("Some error happened: ", err)
	}
}

func GetTypeChartData(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	keys, _ = r.URL.Query()["chartType"]
	chartType := keys[0]

	data := service.GetFileChartData(path, chartType)
	var err = json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Println("Some error happened: ", err)
	}
}

func DeleteFiles(w http.ResponseWriter, r *http.Request) {

	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	keys, _ = r.URL.Query()["recursive"]
	recursive, _ := strconv.ParseBool(keys[0])

	keys, _ = r.URL.Query()["empty"]
	empty, _ := strconv.ParseBool(keys[0])

	keys, _ = r.URL.Query()["createdBefore"]
	fmt.Println("Created: ", keys[0])
	createdBefore, _ := time.Parse("02-Jan-2006", keys[0])

	keys, _ = r.URL.Query()["notAccessedAfter"]
	fmt.Println("Accessed: ", keys[0])
	notAccessedAfter, _ := time.Parse("02-Jan-2006", keys[0])

	fmt.Println("Path: ", path)
	fmt.Println("Recursive: ", recursive)
	fmt.Println("Empty: ", empty)
	fmt.Println("CB: ", createdBefore)
	fmt.Println("NA: ", notAccessedAfter)

	deleteData := model.DeleteData{path, recursive, empty, createdBefore,
		notAccessedAfter}
	service.DeleteFiles(&deleteData)

	json.NewEncoder(w).Encode("Deleted")

}

func ReorganizeFiles(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["src"]
	src := keys[0]

	keys, _ = r.URL.Query()["dest"]
	dest := keys[0]

	keys, _ = r.URL.Query()["recursive"]
	recursive, _ := strconv.ParseBool(keys[0])

	keys, _ = r.URL.Query()["fileType"]
	fileType, _ := strconv.ParseBool(keys[0])

	keys, _ = r.URL.Query()["fileSize"]
	fileSize, _ := strconv.ParseInt(keys[0], 10, 64 )

	keys, _ = r.URL.Query()["createdDate"]
	createdDate := keys[0]

	reorganizeData := model.ReorganizeData{src, dest, recursive,
		fileType, fileSize, createdDate}
	service.ReorganizeFiles(&reorganizeData)

}
