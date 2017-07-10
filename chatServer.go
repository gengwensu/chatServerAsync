package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

//MAXOUTPUT max number of entries output to user
const MAXOUTPUT = 100

//ReqEntry User is the user name,Text is the text message for the chat
type ReqEntry struct {
	User string
	Text string
}

type chatEntry struct {
	Timestamp int64
	ReqEntry
}

type chatDatabase []chatEntry

func main() {
	//create chatDB
	chatDB := chatDatabase{}
	log.Fatal(http.ListenAndServe("localhost:8081", &chatDB))
}

func (db *chatDatabase) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/message": //message POST
		if req.Method == "POST" {
			var re ReqEntry
			err := json.NewDecoder(req.Body).Decode(&re)
			if err != nil {
				http.Error(w, "Error decoding JSON request body",
					http.StatusInternalServerError)
			}
			c := chatEntry{time.Now().Unix(), re}

			*db = append(*db, c) //adding UNIX timestamp

			fmt.Fprint(w, "{\n 'ok': true\n}\n")
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	case "/messages": //messages GET
		if req.Method == "GET" {
			s := *db
			//sort reverse by Timestamp
			sort.Slice(s, func(i, j int) bool { return s[j].Timestamp < s[i].Timestamp })
			o := make([]chatEntry, MAXOUTPUT)
			n := copy(o, s) //copy the min of len(s) or MAXOUTPUT
			mout, err := json.MarshalIndent(o[:n], "", " ")
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}
			fmt.Fprintf(w, "{\n\"messages\":   %+v\n}\n", string(mout))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	case "/users": //users GET
		if req.Method == "GET" {
			d := map[string]bool{}
			for _, re := range *db { //get rid of duplicates in user names
				d[re.User] = true
			}
			userList := []string{}
			for key := range d {
				userList = append(userList, key)
			}
			sort.Strings(userList) // sort it alphabetically
			uout, err := json.MarshalIndent(userList, "", " ")
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}
			fmt.Fprintf(w, "{\n \"users\":   %+v\n}\n", string(uout))
			//fmt.Fprintf(w, "db content: %+v\n", db.d)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "http 404, %s invalid. Only messge,messges,users are allowed.\n", req.URL)
	}
}
