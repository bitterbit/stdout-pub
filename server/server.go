package server

import (
	"log"
	"net/http"
        "io/ioutil"

	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)

/*
    URLS:
    /client - we will drink all stdout into channel
    /server - we will pour the channel into ws to be displayed
*/

type RoddyServer struct {
        upgrader websocket.Upgrader
        index []byte
}

func NewRoddyServer(indexPath string) *RoddyServer {
    index, _ := ioutil.ReadFile("static/index.html")
    return &RoddyServer{
        upgrader:  websocket.Upgrader{},
        index: index,
    }
}


func (this *RoddyServer) handleHome(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "text/html")
    w.Write(this.index)
}

func (thisn *RoddyServer) loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}

func (this *RoddyServer) handlePiper(w http.ResponseWriter, r *http.Request){
    c, err := this.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()
    for {
        mt, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s mt: %v", message, mt)

        // err = c.WriteMessage(mt, message)
        // if err != nil {
        //    log.Println("write:", err)
        //     break
        // }
    }

}

func (this *RoddyServer) handleDashboard(){
    // TODO
}

func (this *RoddyServer) Start(addr string) error {
    r := mux.NewRouter()
    r.HandleFunc("/echo", this.handlePiper)
    r.HandleFunc("/", this.handleHome)

    fs := http.FileServer(http.Dir("static"))
    r.PathPrefix("/static").Handler( http.StripPrefix("/static/", fs))

    r.Use(this.loggingMiddleware)
    return http.ListenAndServe(addr, r)
}
