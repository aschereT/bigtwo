package main

import (
	"fmt"
	"math/rand"
	"time"
)

var CardNames = map[int]string{
	0: "3 Diamonds", 1: "3 Clubs", 2: "3 Hearts", 3: "3 Spades",
	4: "4 Diamonds", 5: "4 Clubs", 6: "4 Hearts", 7: "4 Spades",
	8: "5 Diamonds", 9: "5 Clubs", 10: "5 Hearts", 11: "5 Spades",
	12: "6 Diamonds", 13: "6 Clubs", 14: "6 Hearts", 15: "6 Spades",
	16: "7 Diamonds", 17: "7 Clubs", 18: "7 Hearts", 19: "7 Spades",
	20: "8 Diamonds", 21: "8 Clubs", 22: "8 Hearts", 23: "8 Spades",
	24: "9 Diamonds", 25: "9 Clubs", 26: "9 Hearts", 27: "9 Spades",
	28: "10 Diamonds", 29: "10 Clubs", 30: "10 Hearts", 31: "10 Spades",
	32: "Jack Diamonds", 33: "Jack Clubs", 34: "Jack Hearts", 35: "Jack Spades",
	36: "Queen Diamonds", 37: "Queen Clubs", 38: "Queen Hearts", 39: "Queen Spades",
	40: "King Diamonds", 41: "King Clubs", 42: "King Hearts", 43: "King Spades",
	44: "Ace Diamonds", 45: "Ace Clubs", 46: "Ace Hearts", 47: "Ace Spades",
	48: "2 Diamonds", 49: "2 Clubs", 50: "2 Hearts", 51: "2 Spades",
}

var PlayErrors = map[int]string{
	0: "Successful", 1: "Not this player's turn", 2: "A card is played twice", 3: "You don't have this card",
}

type Player struct {
	ID   int
	Deck []int
}

type GameState struct {
	Players   [4]Player
	Discard   []int
	CurPlayer int
}

func NewGameState() GameState {
	fmt.Println("Main:NewGameState: New game")
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	gs := GameState{Players: [4]Player{Player{ID: 0}, Player{ID: 1}, Player{ID: 2}, Player{ID: 3}}}

	for i := 0; i < 52; i++ {
		for {
			k := r1.Intn(4)
			if len(gs.Players[k].Deck) < 13 {
				gs.Players[k].Deck = append(gs.Players[k].Deck, i)
				break
			}
		}
	}
	return gs
}

func (gs *GameState) Print() {
	fmt.Println("Last played cards:")
	for _, k := range gs.Discard {
		fmt.Print(CardNames[k], "(", k, ")", ", ")
	}
	for i := 0; i < 4; i++ {
		fmt.Println("Player", i, "cards", len(gs.Players[i].Deck), ":")
		for _, k := range gs.Players[i].Deck {
			fmt.Print(CardNames[k], "(", k, ")", ", ")
		}
		fmt.Println("")
	}
}

func (gs *GameState) getCurPlayer() int {
	return gs.CurPlayer
}

func (gs *GameState) getNextPlayer() int {
	return (gs.CurPlayer + 1) % 4
}

func (gs *GameState) play(player int, play []int) int {
	//check if player's turn
	if gs.CurPlayer != player {
		return 1
	}
	//check if player passes
	if len(play) == 0 {

	}
	//check if 3 other players passes

	for i := 0; i < len(play); i++ {
		for j := i + 1; j < len(play); j++ {
			//check all cards are unique
			if i == j {
				return 2
			}
		}
	}
	//check if player has this card on hand
	// for card := range gs.Players[player].Deck {

	// }
	//Checks to see if play is valid:
	//check if bigger than last played
	//remove old discard
	//move from player hand to discard
	//check if player wins (0 cards in hand)
	//change to next player
	return 0
}

func main() {
	fmt.Println("Big2 Game started")
	gs := NewGameState()
	fmt.Println("Sample gamestate:", gs)
	gs.Print()
	DebugPlay(gs)
}

//https://stackoverflow.com/questions/15413469/how-to-make-fmt-scanln-read-into-a-slice-of-integers
func intScanln(n int) ([]int, error) {
	x := make([]int, n)
	y := make([]interface{}, len(x))
	for i := range x {
		y[i] = &x[i]
	}
	n, err := fmt.Scanln(y...)
	x = x[:n]
	return x, err
}

func DebugPlay(gs GameState) {
	for {
		fmt.Println("Currently turn of player", gs.getCurPlayer())
		fmt.Println("Which cards to play? Enter values for cards, space separated")

		cardsPlayed, err := intScanln(5)
		for err != nil {
			fmt.Println("Can't read that, try again. Error:", err)
			cardsPlayed, err = intScanln(5)
			//os.Exit(1)
		}
		fmt.Println("Cards played:", cardsPlayed)
		playErr := gs.play(gs.getCurPlayer(), cardsPlayed)
		fmt.Println("Results:", playErr, PlayErrors[playErr])
	}
}
