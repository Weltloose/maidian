package main

import (
	"log"
	"maidian/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/emit_pv", handler.EmitPv)
	http.HandleFunc("/query_pv", handler.QueryPv)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
