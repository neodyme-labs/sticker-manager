package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/matrix-org/gomatrixserverlib/fclient"
	"github.com/matrix-org/gomatrixserverlib/spec"
)

type RespOpenIDToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	MatrixServerName string `json:"matrix_server_name"`
	TokenType        string `json:"token_type"`
}

type Scalar struct {
	ScalarToken string `json:"scalar_token"`
}

var stickerUrl string

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	var req RespOpenIDToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := fclient.NewClient(fclient.WithWellKnownSRVLookups(true))
	info, err := client.LookupUserInfo(r.Context(), spec.ServerName(req.MatrixServerName), req.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := Scalar{
		ScalarToken: info.Sub,
	}

	json.NewEncoder(w).Encode(res)
}

func account(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()

	if r.Method == http.MethodOptions {
		return
	}

	res := make(map[string]string)
	res["user_id"] = query.Get("scalar_token")

	json.NewEncoder(w).Encode(res)
}

func success(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func ui(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	query := r.URL.Query()
	screen := query.Get("screen")

	switch screen {
	case "type_m.stickerpicker":
		fmt.Fprintf(w, `
		<script>
			setInterval(() => {
				window.parent.postMessage({
					"action": "set_widget",
					"widget_id": "stickerpicker",
					"url": "%s",
					"type": "m.stickerpicker",
					"userWidget": true
				}, "*");

				window.parent.postMessage({
					"action": "close_scalar"
				}, "*");
			}, 1000);
		</script>
		`, stickerUrl)
	default:
		http.Error(w, "Unkown screen", http.StatusBadRequest)
		return
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	port := flag.String("port", "8080", "Port to run the web server on")
	host := flag.String("host", "127.0.0.1", "Host to run the web server on")
	flag.StringVar(&stickerUrl, "sticker_url", "https://example.com/?theme=$theme", "Server which serves the stickers")
	flag.Parse()

	http.HandleFunc("/", ui)
	http.HandleFunc("/api/register", register)
	http.HandleFunc("/api/account", account)
	http.HandleFunc("/api/widgets/set_assets_state", success)
	http.Handle("/*", http.NotFoundHandler())

	fmt.Printf("Starting server on %s:%s...\n", *host, *port)
	err := http.ListenAndServe(*host+":"+*port, logRequest(http.DefaultServeMux))
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
