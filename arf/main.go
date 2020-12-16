/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

	package main

import "github.com/AnaMijailovic/NTP/arf/cmd"

func main() {
		cmd.Execute()
	/*	renameData := model.RenameData{"C:/Users/anaam/Desktop/del", true, true, "",
			"", "" }
		start := time.Now()
		fmt.Println("Started parallel..")
		errs := service.Rename(&renameData)
		fmt.Println("From main ...")
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Println(err)
			}
		}


		fmt.Printf("\nTime elapsedP: %v\n",  time.Since(start))
	} */






	/*package main

	import (
		"fmt"
		"sync"
	)

	func worker (count *int, lock chan bool, wg *sync.WaitGroup) {
		lock <- true // Zaključaj resurs
		*count ++
		<- lock // Otključaj resurs
		wg.Done() // Smanji brojač go rutina

	}

	func main() {
		// Kreiranje wait group-e, ne mora da bude inicijalizovana
		var wg sync.WaitGroup
		lock := make(chan bool, 1)
		count := 0 // deljeni resurs
		for i := 0; i < 10; i++ {
			wg.Add(1) // Uvećaj brojač go rutina
			go worker(&count, lock, &wg)
		}

		// Sačekaj da se pokrenute go rutine izvrše
		wg.Wait()
		fmt.Println("Count: ", count) */
	}

/*
	package main

	import (
		"fmt"
		"controller/http"
	)

	func helloHandler(writer http.ResponseWriter, request *http.Request){
		fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])
	}

	func main() {
		http.HandleFunc("/", helloHandler)
		http.ListenAndServe(":8080", nil)
	}



*/