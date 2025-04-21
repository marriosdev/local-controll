package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	ip := getLocalIPv4()
	port := "8080"
	address := fmt.Sprintf("%s:%s", ip, port)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Olá! Você acessou %s\n", r.URL.Path)
	})

	fmt.Printf("Servidor rodando em http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func getLocalIPv4() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conn)
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
