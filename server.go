package horloge

import "net/http"

func PingHandler(w http.ResponseWriter, r *http.Request) {

}

func VersionHandler(w http.ResponseWriter, r *http.Request) {

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

}

func RegisterationHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/ping", PingHandler)
	r.HandleFunc("/version", VersionHandler)
	r.HandleFunc("/health_check", HealthCheckHandler)
	r.HandleFunc("/register", RegisterationHandler)

	http.Handle("/", r)
}
