package main

import (
	"views";
	"http";
)

func main() {
	http.Handle("/blog", http.HandlerFunc(views.Index));
	http.Handle("/blog/entry/add", http.HandlerFunc(views.AddEntry));
	http.Handle("/blog/entry/", http.HandlerFunc(views.Entry));
	http.Handle("/blog/comment/add", http.HandlerFunc(views.AddComment));
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic("ListenAndServe: ", err.String());
	}
}
