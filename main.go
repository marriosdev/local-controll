package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
)

var ServerIp string
var lastX, lastY int

type HomeData struct {
	ServerIp string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	ServerIp = getLocalIPv4()
	port := "8080"
	address := fmt.Sprintf("%s:%s", ServerIp, port)

	http.HandleFunc("/", serveTouchpad)
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Printf("Acesse o touchpad em http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func serveTouchpad(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, HomeData{ServerIp: ServerIp})
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao subir WebSocket:", err)
		return
	}
	defer conn.Close()

	screenWidth, screenHeight := robotgo.GetScreenSize()

	for {
		var pos struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		}
		err := conn.ReadJSON(&pos)
		if err != nil {
			log.Println("Erro ao ler dados WebSocket:", err)
			break
		}

		normalizedX := int(pos.X * float64(screenWidth) / 1000)
		normalizedY := int(pos.Y * float64(screenHeight) / 1000)

		moveX := normalizedX - lastX
		moveY := normalizedY - lastY

		robotgo.MoveMouse(lastX+moveX, lastY+moveY)

		lastX, lastY = lastX+moveX, lastY+moveY
	}
}

func getLocalIPv4() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
