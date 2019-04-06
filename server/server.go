package server

import (
	"log"
	"net/http"
        "io/ioutil"

	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)

var ChanCapacity = 1024

type DashboardMessage struct {
    Date, Source, Message string
}

/*
    URLS:
    /client - we will drink all stdout into channel
    /server - we will pour the channel into ws to be displayed
*/

type RoddyServer struct {
        upgrader websocket.Upgrader
        index []byte
        pipers map[string]chan DashboardMessage
}

func NewRoddyServer(indexPath string) *RoddyServer {
    index, _ := ioutil.ReadFile("static/index.html")
    return &RoddyServer{
        upgrader:  websocket.Upgrader{},
        index: index,
        pipers: make(map[string]chan DashboardMessage),
    }
}


func (this *RoddyServer) handleHome(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "text/html")
    w.Write(this.index)
}

func (thisn *RoddyServer) loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

func (this *RoddyServer) handlePiperWS(w http.ResponseWriter, r *http.Request){
    piperName := r.RemoteAddr

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

        if _, exists := this.pipers[piperName]; !exists {
            this.pipers[piperName] = make(chan DashboardMessage, ChanCapacity)
        }
        ch := this.pipers[piperName]

        // No more space, discard old messages
        if len(ch) == cap(ch) {
            m := <-ch
            log.Printf("Discarding old message. msg: %+v\n", m)
        }

        dashboardMessage := DashboardMessage{
            Date: "now",
            Source: piperName,
            Message: string(message),
        }

        ch <- dashboardMessage
        log.Printf("Done sending. len: %v cap: %v\n", len(ch), cap(ch))
    }

}


func (this *RoddyServer) handleDashboardWS(w http.ResponseWriter, r *http.Request){
    c, err := this.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()

    for {
        for _, ch := range this.pipers {
            message := <-ch
            c.WriteJSON(message)
        }
    }
}

func (this *RoddyServer) Start(addr string) error {
    r := mux.NewRouter()
    r.HandleFunc("/ws/piper", this.handlePiperWS)
    r.HandleFunc("/ws/dashboard", this.handleDashboardWS)
    r.HandleFunc("/", this.handleHome)

    fs := http.FileServer(http.Dir("static"))
    r.PathPrefix("/static").Handler( http.StripPrefix("/static/", fs))

    r.Use(this.loggingMiddleware)
    return http.ListenAndServe(addr, r)
}
