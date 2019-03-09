package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

var JwtSecret []byte

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if string(JwtSecret) == "" {
		log.Fatal("Jwt secret cannot be null")
	} else {
		log.Println("JWT Secret: " + string(JwtSecret))
	}

	RunServer("5000")

	//for i := 23000; i < 23050; i++ {
	//	port := strconv.Itoa(i)
	//	proxyConfig := common.ProxyConfig{"http", port, "", ""}
	//
	//	serverConfig := common.CreateServerConfig(common.USWest1, common.Micro)
	//
	//	go management.StartProxyAndReturnRecord(proxyConfig, serverConfig, time.Minute*10, 1000)
	//	fmt.Println("Sleeping for 5 seconds before starting next server...")
	//	time.Sleep(time.Second * 5)
	//}
	//
	//for true {
	//	fmt.Println("Sleepy")
	//	time.Sleep(time.Minute * 5)
	//}
}

func RunServer(port string) {
	muxx := makeMuxRouter()

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Server is now running on port " + port)
	log.Fatal(s.ListenAndServe()) //Key and cert are already set in the TLS config
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()

	//TODO Remove "OPTIONS"(?) and all CORS stuff for server deployment
	muxRouter.HandleFunc("/login", Login).Methods("GET", "OPTIONS")

	return muxRouter
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	defer r.Body.Close()

	user, pass, _ := r.BasicAuth()

	if user == "RUSTLED" && pass == "JIMMIES" {
		token := jwt.New(jwt.SigningMethodHS256)
		/* Create a map to store our claims */
		claims := token.Claims.(jwt.MapClaims)

		/* Set token claims */
		claims["name"] = user
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

		/* Sign the token with our secret */
		tokenString, _ := token.SignedString(JwtSecret)
		fmt.Println("Created token: " + tokenString)

		/* Finally, write the token to the browser window */
		w.Write([]byte("{\"Token\": \"" + tokenString + "\"}"))
	} else {
		fmt.Println("Login: " + user + ", " + pass + " FAILED")
		w.Write([]byte("{\"Token\": \"FAIL\"}"))
	}
}
