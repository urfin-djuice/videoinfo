package main

import "flag"

func getListenAddr() string {
	la := flag.String("addr", ":80", "Listen address (:8080, localhost:80, 192.168.10.10:8787)")
	flag.Parse()
	return *la
}

func main() {
	app := NewInfoServer(getListenAddr())
	app.Up()
}
