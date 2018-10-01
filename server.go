package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Dependency
//go get github.com/gorilla/mux

type PlayAHand struct {
	PlayerId int   `json:"PlayerId,omitempty"`
	Hand     []int `json:"Hand,omitempty"`
}

var gss map[int]GameState

//Generate an ID for games
func generateId() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return r1.Intn(65535)
}

//play a hand for a game's player
func servPlayHand(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server:PlayHand called")
	//Grab PlayAHand struct from POST body
	var hand PlayAHand
	_ = json.NewDecoder(r.Body).Decode(&hand)
	//grab gameid from POST URL
	gameid, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Server:PlayHand: error parsing hand")
		return
	}
	fmt.Println("Server:PlayHand: Game", gameid, "Player", hand.PlayerId, "is playing", hand.Hand)
	gs, found := gss[gameid]
	if found {
		//we found such a game, play the hand
		gs.play(hand.PlayerId, hand.Hand)
		json.NewEncoder(w).Encode(gs.getGameState())
	} else {
		//no such game found
		fmt.Println("Server:PlayHand: Game", gameid, " not found")
	}
}

//get the current game state
//TODO: only show cards for a given player
func servGetGameState(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server:GetGameState: called")
	gameid, err := strconv.Atoi(mux.Vars(r)["gameid"])
	if err != nil {
		fmt.Println("Server:GetGameState: parsing gameid", mux.Vars(r)["gameid"], "failed")
		return
	}
	gs, found := gss[gameid]
	if !found {
		fmt.Println("Server:GetGameState: gameid", gameid, "can't be found")
		return
	}
	json.NewEncoder(w).Encode(gs.getGameState())
}

//generate a new game
func servNewGame(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server:NewGame: called")
	gs := NewGameState()
	for {
		//keep trying to generate ID until an unused one is found
		id := generateId()
		_, ex := gss[id]
		if !ex {
			//we found one
			gs.GameId = id
			gss[id] = gs
			fmt.Println("Gamestate ID", id, "Gamestate", gs)
			break
		}
	}
	json.NewEncoder(w).Encode(gs.getGameState())
}

func TestStuff() {
	json.NewEncoder(os.Stdout).Encode(PlayAHand{Hand: []int{1, 2, 3}})
}

func StartServer() {
	fmt.Println("StartServer: called")
	gss = map[int]GameState{}
	router := mux.NewRouter()

	router.HandleFunc("/newgame", servNewGame).Methods("GET")
	router.HandleFunc("/gamestate/{gameid}", servGetGameState).Methods("GET")
	//where id is player's id
	router.HandleFunc("/play/{gameid}", servPlayHand).Methods("POST")

	//TestStuff() //{"Hand":[1,2,3]}

	log.Fatal(http.ListenAndServe(":80", router))
}
