package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/denverquane/go-ec2-proxy/management"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var JwtSecret []byte

type ProxyServerListing struct {
	IP   string
	Port string
}

var JWTStatus = make(map[string][]*management.ServerRecord) // connected clients
var lock = sync.RWMutex{}

const MinPort = 32768
const MaxPort = 61000

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

	muxRouter.HandleFunc("/token", StartServerWithJWT).Methods("POST", "OPTIONS")
	muxRouter.HandleFunc("/tokenStatus/{jwt}", handleGetStatus).Methods("GET", "OPTIONS")
	muxRouter.HandleFunc("/test", handleTest).Methods("GET", "OPTIONS")

	return muxRouter
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	w.Write([]byte(`[{"InstanceId":"i-083f83365ee655125","Region":"us-west-1","User":"","Pass":"","PublicIp":"13.57.209.45","PublicPort":"56473","Constraints":{"DestructionTime":"2019-03-16T19:23:19.2045263-06:00","TotalByteCap":1000000000},"Status":{"CreationTime":"2019-03-15T19:23:19.2045263-06:00","InboundBytesUsed":0,"OutboundBytesUsed":0,"IsRunning":true,"IsDestroyed":false}},{"InstanceId":"i-083f83365ee655125","Region":"us-west-1","User":"","Pass":"","PublicIp":"13.57.209.45","PublicPort":"56473","Constraints":{"DestructionTime":"2019-03-16T19:23:19.2045263-06:00","TotalByteCap":1000000000},"Status":{"CreationTime":"2019-03-15T19:23:19.2045263-06:00","InboundBytesUsed":0,"OutboundBytesUsed":0,"IsRunning":true,"IsDestroyed":false}}]`))
}

func handleGetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	vars := mux.Vars(r)
	jwtt := vars["jwt"]

	lock.RLock()
	defer lock.RUnlock()

	if val, ok := JWTStatus[jwtt]; !ok {
		w.Write([]byte("{\"Status\": \"Invalid JWT!\"}"))
	} else {
		byt, err := json.Marshal(val)
		if err == nil {
			w.Write(byt)
		} else {
			log.Println(err)
		}
	}
}

//Simple function to test the JWT against our signage key
func StartServerWithJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Token")

	time.Sleep(2 * time.Second) // Allow for the creation time of the JWT to come into play

	claims, err := GetStructuredClaimsFromRequest(JwtSecret, r)
	if err != nil {
		w.Write([]byte("{\"Status\": \"INVALID\"}"))
	} else {
		serverType := common.Nano
		proxyConfig := common.ProxyConfig{"http", "99999", "", ""}

		if claims.Type == "t2.micro" {
			serverType = common.Micro
		}

		serverConfig := common.CreateServerConfig(common.USWest1, serverType)

		token, err := request.HeaderExtractor{"Token"}.ExtractToken(r)
		if err != nil {
			log.Println(err)
		}

		gb := claims.Data
		//go waitForServerAndBcast(proxyConfig, serverConfig, time.Hour*time.Duration(claims.Duration), float64(gb)*1000000000.0, token)

		for i := 0; i < int(claims.Quantity); i++ {
			port := rand.Intn(MaxPort-MinPort) + MinPort

			proxyConfig.Port = strconv.FormatInt(int64(port), 10)

			go waitForServerAndBcast(proxyConfig, serverConfig, time.Hour*time.Duration(claims.Duration), float64(gb)*1000000000.0, token)
		}

		//str := "Username:" + claims.Username + ", Expiration: " + claims.Expiration.String()
		w.Write([]byte("{\"Status\": \"VALID\"}"))
	}
	fmt.Println(claims)
}

func waitForServerAndBcast(pConfig common.ProxyConfig, sConfig common.ServerConfig, duration time.Duration, bytes float64, token string) {
	log.Println("Starting server with duration: " + duration.String() + " and data= " + strconv.FormatFloat(bytes, 'f', -1, 64) + " bytes")

	record, err := management.StartProxyAndReturnRecord(pConfig, sConfig, duration, bytes)
	if err != nil {
		fmt.Println(err)
	}

	lock.Lock()
	defer lock.Unlock()

	JWTStatus[token] = append(JWTStatus[token], record)
}

//simple test struct of basic info (expand)
type JWTClaimFields struct {
	//Receipt_url   string
	DiscordID string
	Type      string
	Duration  int64
	Quantity  int64
	Data      int64
}

// extract the claims from the token
func extractStructuredClaimsFromToken(jwtsecret []byte, tokenString string) (JWTClaimFields, error) {
	ReturnClaims := JWTClaimFields{}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtsecret, nil
	})
	if token == nil {
		return JWTClaimFields{}, errors.New("Token is invalid!")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//ReturnClaims.Receipt_url = claims["Receipt_url"].(string)
		ReturnClaims.DiscordID = claims["DiscordID"].(string)
		ReturnClaims.Type = claims["Type"].(string)
		ReturnClaims.Duration = int64(claims["Duration"].(float64))
		ReturnClaims.Quantity = int64(claims["Quantity"].(float64))
		ReturnClaims.Data = int64(claims["Data"].(float64))
		return ReturnClaims, nil
	} else {
		fmt.Println(err)
		return JWTClaimFields{}, err
	}
}

// Get the claims from the http request
func GetStructuredClaimsFromRequest(secret []byte, req *http.Request) (JWTClaimFields, error) {
	tokenString, err2 := request.HeaderExtractor{"Token"}.ExtractToken(req)
	if err2 == nil {
		return extractStructuredClaimsFromToken(secret, tokenString)
	} else {
		fmt.Println(err2)
		return JWTClaimFields{}, err2
	}
}
