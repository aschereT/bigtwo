package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
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

//3-10, 11 = Jack, 12 = Queen, 13 = King, 14 = Ace, 15 = 2
func getCardValue(c int) int {
	return (c / 4) + 3
}

//0 = Diamonds, 1 = Clubs, 2 = Hearts, 3 = Spades
func getCardSuit(c int) int {
	return c % 4
}

var PlayErrors = map[int]string{
	0: "Successful", 1: "Not this player's turn", 2: "A card is played twice", 3: "You don't have this card", 7: "Winner!",
}

type Player struct {
	ID   int   `json:"ID,omitempty"`
	deck []int `json:"deck,omitempty"`
}

type GameState struct {
	GameId    int       `json:"gameid,omitempty`
	players   [4]Player `json:"players,omitempty`
	discard   []int     `json:"discard,omitempty"`
	curPlayer int       `json:"curPlayer,omitempty"`
	finished  bool      `json:"finished,omitempty"`
}

func NewGameState() GameState {
	fmt.Println("Main:NewGameState: New game")
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	gs := GameState{players: [4]Player{Player{ID: 0}, Player{ID: 1}, Player{ID: 2}, Player{ID: 3}}}

	for i := 0; i < 52; i++ {
		for {
			k := r1.Intn(4)
			if len(gs.players[k].deck) < 13 {
				gs.players[k].deck = append(gs.players[k].deck, i)
				break
			}
		}
	}
	return gs
}

func (gs *GameState) Print() {
	fmt.Println("Last played cards:")
	for _, k := range gs.discard {
		fmt.Print(CardNames[k], "(", k, ")", ", ")
	}
	fmt.Print("\n")
	for i := 0; i < 4; i++ {
		fmt.Println("Player", i, "cards", len(gs.players[i].deck), ":")
		for _, k := range gs.players[i].deck {
			fmt.Print(CardNames[k], "(", k, ")", ", ")
		}
		fmt.Println("")
	}
}

func (gs *GameState) getGameState() string {
	return fmt.Sprint(gs)
}

func (gs *GameState) getCurPlayer() int {
	return gs.curPlayer
}

func (gs *GameState) getNextPlayer() int {
	return (gs.curPlayer + 1) % 4
}

//given a set of cards, calculate value, suit, and combo
//compare cards based on combo and value
//value = highest card in a play - sort cards from lowest to highest and return value of last card
//combos:
//0 = invalid
//1 = singles
//2 = pair
//3 = three of a kind
//4 = straight
//5 = flush
//6 = full house
//7 = four of a kind
//8 = straight flush
//9 = royal flush
func calcVal(play []int) (combo int, value int, suit int) {
	combo = 0 //assume invalid combo to begin with
	if len(play) == 4 {
		//invalid af, no 4 card plays
		return 0, 0, 0
	}

	//reorder cards from smallest to largest
	sort.Ints(play)

	cardVal := make([]int, len(play))
	for i := 0; i < len(play); i++ {
		cardVal[i] = play[i]/4 + 3
	}

	cardSuit := make([]int, len(play))
	for i := 0; i < len(play); i++ {
		cardSuit[i] = play[i] % 4
	}

	//singles
	if len(play) == 1 {
		return 1, cardVal[0], cardSuit[0]
	}

	//pair
	if len(play) == 2 {
		if cardVal[0] == cardVal[1] {
			return 2, cardVal[0], cardSuit[1]
		}
	}

	//triples
	if len(play) == 3 {
		if cardVal[0] == cardVal[1] && cardVal[1] == cardVal[2] {
			return 3, cardVal[2], cardSuit[2]
		}
	}

	//straight
	straight := false
	for i := 1; i < len(play); i++ {
		if cardVal[i]-cardVal[i-1] != 1 {
			straight = true
			combo = 4
		}
	}

	//flush
	flush := false
	for i := 1; i < len(play); i++ {
		if cardSuit[i] == cardSuit[i-1] {
			flush = true
			combo = 5
		}
	}

	//full house(refactored)
	if cardVal[0] == cardVal[1] && cardVal[1] == cardVal[2] {
		//first 3
		if cardVal[3] == cardVal[4] {
			return 6, cardVal[2], cardSuit[2]

		}
	} else if cardVal[1] == cardVal[2] && cardVal[2] == cardVal[3] {
		//this case is probably impossible;
		//needs proof
		if cardVal[0] == cardVal[4] {
			return 6, cardVal[2], cardSuit[3]
		}

	} else if cardVal[2] == cardVal[3] && cardVal[3] == cardVal[4] {
		//last 3
		if cardVal[0] == cardVal[1] {
			return 6, cardVal[2], cardSuit[4]
		}
	}

	//full house
	// fullHouse := false
	// threeOfAKind := 0
	// threeOfAKind_index := make([]int, len(play))
	// for i := 1; i < len(play); i++ {
	// 	if cardVal[i] == cardVal[i-1] {
	// 		threeOfAKind++
	// 		threeOfAKind_index[i-1] = 1
	// 		threeOfAKind_index[i] = 1
	// 	} else {
	// 		threeOfAKind_index[i] = 0
	// 	}
	// }
	// // check if remaining two cards are the same - depends on the card values being sorted
	// if threeOfAKind == 3 {
	// 	for i := 0; i < len(play)-1; i++ {
	// 		if threeOfAKind_index[i] != 1 {
	// 			if cardVal[i] == cardVal[i+1] {
	// 				fullHouse = true //full house valid
	// 				combo = 6
	// 			}
	// 		}
	// 	}
	// }
	// if fullHouse {
	// 	combo = 6
	// }

	//four of a kind
	fourOfAKind := false
	countFour := 0
	for i := 1; i < len(play); i++ {
		if cardSuit[i] == cardSuit[i-1] {
			countFour++
		}
	}
	if countFour == 4 {
		fourOfAKind = true
		combo = 7
	}
	if fourOfAKind {
		combo = 7
	}

	//straight flush
	if straight && flush {
		combo = 8
	}

	value = cardVal[4]

	return combo, value, cardSuit[4]
}

func (gs *GameState) play(player int, play []int) int {
	//check if player's turn
	if gs.curPlayer != player {
		return 1
	}
	//check if player passes
	if len(play) == 0 {
		//change to next player
		gs.curPlayer = gs.getNextPlayer()
		return 0
	}
	//check player is not playing too many cards
	if len(play) > 5 {

	}
	//check if 3 other players passes

	for i := 0; i < len(play); i++ {
		for j := i + 1; j < len(play); j++ {
			//check all cards are unique
			if play[i] == play[j] {
				return 2
			}
		}
	}

	//check if player has this card on hand
	var validcard = 0
	for card := range gs.players[player].deck {
		for i := 0; i < len(play); i++ {
			if play[i] == card {
				validcard++
			}
		}
	}
	if validcard != len(play) {
		return 3
	}
	//Checks to see if play is valid:
	//check if bigger than last player

	//remove old discard, and
	//move from player hand to discard
	gs.discard = play
	gs.players[player].removeCards(play)

	//check if player wins (0 cards in hand)
	if (len(gs.players[player].deck)) == 0 {
		//WEINER WEINER CHEAPER DINERS
		return 7
	}
	//change to next player
	gs.curPlayer = gs.getNextPlayer()
	return 0
}

//remove a card from a player's hands
func (pl *Player) removeCard(card int) int {
	found := 0
	newDeck := []int{}
	for _, n := range pl.deck {
		if n != card {
			newDeck = append(newDeck, n)
		} else {
			found = 1
		}
	}
	pl.deck = newDeck
	return found
}

//remove a set of cards from player's hands
func (pl *Player) removeCards(cards []int) int {
	found := 0
	newDeck := []int{}
	for _, n := range pl.deck {
		fd := true
		for _, k := range cards {
			if n == k {
				found++
				fd = false
				break
			}
		}
		if fd {
			newDeck = append(newDeck, n)
		}
	}
	pl.deck = newDeck
	return found
}

func main() {
	fmt.Println("Big2 Game started")
	gs := NewGameState()
	DebugPlay(gs)

	//gs.Print()
	//StartServer()
	select {}

}

//https://stackoverflow.com/questions/43599253/read-space-separated-integers-from-stdin-into-int-slice?rq=1
func readHand() []int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var n []int
	for _, f := range strings.Fields(scanner.Text()) {
		i, err := strconv.Atoi(f)
		if err == nil {
			n = append(n, i)
		}
	}
	return n
}

func DebugPlay(gs GameState) {
	for {
		gs.Print()
		fmt.Println("Currently turn of player", gs.getCurPlayer())

		fmt.Println("Which cards to play? Enter values for cards, space separated")
		cardsPlayed := readHand()
		fmt.Println("Cards played:", cardsPlayed)

		playErr := gs.play(gs.getCurPlayer(), cardsPlayed)
		fmt.Println("Results:", playErr, PlayErrors[playErr])

		playCombo, playVal, playSuit := calcVal(cardsPlayed)
		if playVal == 0 {
			fmt.Println("Not a valid play")
		} else {
			fmt.Println("Play: Combo", playCombo, "Value ", playVal, "Suit", playSuit)
		}
	}
}
