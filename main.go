package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: "TBD",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.tdameritrade.com/auth",
			TokenURL: "https://api.tdameritrade.com/v1/oauth2/token",
		},
	}
)

//Homepage
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Homepage Hit")

	u := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, u, http.StatusFound)
}

//Authorize
func Authorize(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Auth reached")
	r.ParseForm()
	state := r.Form.Get("state")
	if state != "state" {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)
}

func main() {
	// RefreshAccessCode(refreshCode)

	// return

	// requestBody, err := json.Marshal(map[string]string{
	// 	"apikey": clientID,
	// })

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// timeout := time.Duration(10 * time.Second)
	// client := http.Client{
	// 	Timeout: timeout,
	// }

	// request, err := http.NewRequest("GET", "https://api.tdameritrade.com/v1/marketdata/AAPL/quotes", bytes.NewBuffer(requestBody))
	// request.Header.Set("Content-type", "application/json")
	// request.Header.Add("Authorization", "Bearer "+authCode)

	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// resp, err := client.Do(request)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println(string(body))

	// return

	// acsToken := authCode
	// config := &oauth2.Config{
	// 	ClientID:     clientID,
	// 	ClientSecret: "",
	// 	Endpoint: oauth2.Endpoint{
	// 		TokenURL: urlAccessToken,
	// 	},
	// }

	// restoredToken := &oauth2.Token{
	// 	AccessToken:  acsToken,
	// 	RefreshToken: refreshCode,
	// 	Expiry:       time.Now(),
	// 	TokenType:    "Bearer",
	// }
	// tokenSource := config.TokenSource(oauth2.NoContext, restoredToken)
	// client := oauth2.NewClient(oauth2.NoContext, tokenSource)
	// savedToken, err := tokenSource.Token()
	// acsToken = savedToken.AccessToken
	// handleFatalErr("Main/savedToken", err)

	client, code, expiry, err := GetClient(authCode, clientID, refreshCode)
	handleFatalErr("Main/Client", err)
	fmt.Printf("Authorized Code: %v\r\n", code)
	fmt.Printf("Expiry: %v\r\n", expiry)
	fmt.Printf("%v\r\n", GetOptionChain(client, "NIO"))

	// requestBody, err := json.Marshal(map[string]string{
	// 	"apikey":        clientID,
	// 	"symbol":        "AAPL",
	// 	"includeQuotes": "true",
	// })
	// handleFatalErr("GetOptionChain/requestBody", err)
	// request, err := http.NewRequest("GET", "https://api.tdameritrade.com/v1/marketdata/chains", bytes.NewBuffer(requestBody))
	// request.Header.Set("Content-type", "application/json")
	// request.Header.Add("Authorization", "Bearer "+acsToken)

	// handleFatalErr("GetOptionChain/request", err)
	// resp, err := client.Do(request)
	// handleFatalErr("GetOptionChain/resp", err)
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// handleFatalErr("GetOptionChain/body", err)
	// log.Println(string(body))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/oauth2", Authorize)
	log.Println("Client is running at 9094 port.")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
