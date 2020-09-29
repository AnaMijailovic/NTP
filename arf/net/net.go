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
	http.HandleFunc("/api/fileTypeData", GetChartData)
	http.HandleFunc("/api/deleteFiles", DeleteFiles)
	http.HandleFunc("/api/rename", RenameFiles)
	http.HandleFunc("/api/reorganize", ReorganizeFiles)
	http.HandleFunc("/api/recover", Recover)

	fmt.Println("Starting server ...")
	http.ListenAndServe(":8080", nil)

}

func GetFileTree(w http.ResponseWriter, r *http.Request) {

	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	tree := service.CreateTree(path, true)

	var err = json.NewEncoder(w).Encode(tree)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetChartData(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	keys, _ = r.URL.Query()["chartType"]
	chartType := keys[0]

	data := service.GetFileChartData(path, chartType)
	var err = json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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

	deleteData := model.DeleteData{path, recursive, empty, createdBefore,
		notAccessedAfter}

	filesDeleted := service.DeleteFiles(&deleteData)

	json.NewEncoder(w).Encode("Deleted " + strconv.FormatInt(int64(*filesDeleted),10 )+ " files")

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

	errs := service.ReorganizeFiles(&reorganizeData)

	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		if _, ok := errs[0].(model.UnableToRenameFileError); ok {
			json.NewEncoder(w).Encode("Arf was unable to move all files")
		} else {
			json.NewEncoder(w).Encode(errs[0].Error())
		}
	} else {
		json.NewEncoder(w).Encode("Files were successfully reorganized")
	}

}

func RenameFiles(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	keys, _ = r.URL.Query()["recursive"]
	recursive, _ := strconv.ParseBool(keys[0])

	keys, _ = r.URL.Query()["random"]
	random, _ := strconv.ParseBool(keys[0])

	keys, _ = r.URL.Query()["remove"]
	remove := keys[0]

	keys, _ = r.URL.Query()["replaceWith"]
	replaceWith := keys[0]

	keys, _ = r.URL.Query()["pattern"]
	pattern := keys[0]

	renameData := model.RenameData{path, recursive, random,
		remove, replaceWith, pattern}

	errs := service.Rename(&renameData)

	if len(errs) > 0 {
		if _, ok := errs[0].(model.UnableToRenameFileError); ok {
			json.NewEncoder(w).Encode("Arf was unable to rename all files")
		} else {
			json.NewEncoder(w).Encode(errs[0].Error())
		}
	} else {
		json.NewEncoder(w).Encode("Files were successfully renamed.")
	}

}

func Recover(w http.ResponseWriter, r *http.Request) {

	keys, _ := r.URL.Query()["path"]
	path := keys[0]

	errs := service.Recover(path)

	if len(errs) > 0 {
		if _, ok := errs[0].(model.UnableToRenameFileError); ok {
			json.NewEncoder(w).Encode("Arf was unable to recover all files")
		} else {
			json.NewEncoder(w).Encode(errs[0].Error())
		}
	} else {
		json.NewEncoder(w).Encode("Successfully recovered")
	}

}