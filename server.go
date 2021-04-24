package main

import (
	// "io/ioutil" // to parse JSON
    "fmt"
    "log"
    "net/http" // using this package for 
	"github.com/zmb3/spotify"
	"github.com/joho/godotenv"
    "os"
)

import _ "github.com/joho/godotenv/autoload"

// var counter int
// mutex was needed for count variable example, since the server is async
// var mutex = &sync.Mutex{} 
// func incrementCounter(w http.ResponseWriter, r *http.Request) {
//     mutex.Lock()
//     counter++
//     fmt.Fprintf(w, strconv.Itoa(counter))
//     mutex.Unlock()
// }

func echoString(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello")
}

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)


func main() {

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	
	// get enviornment variables
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	
	// build auth url
	auth.SetAuthInfo(clientId, clientSecret)
	url := auth.AuthURL(state)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

    // http.HandleFunc("/", echoString)
	// //http.HandleFunc("https://accounts.spotify.com/authorize", spotifyAuth)

	// // make the GET request
	// resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	// if err != nil {
	//    log.Fatalln(err)
	// }
	// // read the GET request
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// // convert the body of the response to a string
	// sb := string(body)
	// log.Printf(sb)

    // // http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
    // //     fmt.Fprintf(w, "Hi")
    // // })

	// // may need to kill server manually 
	// // lsof -iTCP:8081 -sTCP:LISTEN -> gets pid
	// // kill _pid_ -> kills server
    // log.Fatal(http.ListenAndServe(":8081", nil)) // server on http://localhost:8081/

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}