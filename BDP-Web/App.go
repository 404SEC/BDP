package main

import (
	"log"
	"net/http"
	"os"

	funct "github.com/404SEC/BDP/BDP-Web/function"
	micro "github.com/micro/go-micro"
)

func main() {
	funct.Servicec = micro.NewService(
		micro.Name(funct.Namespace+".Web"),
		micro.Registry(funct.Reg),
	)

	funct.Servicec.Init()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRPCWeb)
	mux.HandleFunc("/api/", handleRPCApi)
	fsh := http.FileServer(http.Dir(wd + "/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fsh))
	log.Println("Listen on :8888")
	http.ListenAndServe(":8888", mux)
}
func handleRPCWeb(w http.ResponseWriter, r *http.Request) {
	funct.HandleJSONRPC(funct.Servicec, r)
	////Here add view Render
}
func handleRPCApi(w http.ResponseWriter, r *http.Request) {
	log.Println("handleRPC coming ....")
	if r.URL.Path == "/" {
		w.Write([]byte("ok,this is the server ..."))
		return
	}
	if origin := r.Header.Get("Origin"); funct.Cors[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else if len(origin) > 0 && funct.Cors["*"] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-Token, X-Client")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write(funct.HandleJSONRPC(funct.Servicec, r))
	return

}
