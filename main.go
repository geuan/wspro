package main

import (
	"log"
	"net/http"
	"wspro/sources/core"
	"wspro/sources/handlers"
)


func main()  {
	http.HandleFunc("/echo",handlers.Echo)
	http.HandleFunc("/sendall", func(w http.ResponseWriter, req *http.Request) {
		//msg := req.URL.Query().Get("msg")
		core.ClientMap.SendAllPods()
		_, _ = w.Write([]byte("OK"))
	})


	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		log.Fatal(err)
	}
}