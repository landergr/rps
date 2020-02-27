package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	PORT    int = 4567
	UNKNOWN     = "UNKNOWN"
	WIN         = "WIN"
	DRAW        = "DRAW"
	LOST        = "LOST"
)

type Game struct {
	AcceptedHands []string
	Rules         map[string]Rule
	Score         *Score
	sync.Mutex
}

func (game Game) play(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var response Response
	playerHand, err := parseHand(r)
	if err != nil {
		response = Response{Result: UNKNOWN, ComputerHand: nil}
	} else {
		computerHand := createComputerHand(game.AcceptedHands)
		response = evaluateHand(playerHand, computerHand, game.Rules)
	}
	game.Score = game.updateScore(response, game.Score)
	b, err := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (game Game) updateScore(response Response, currentScore *Score) *Score {
	game.Lock()
	defer game.Unlock()
	if response.Result == WIN {
		currentScore.Wins = currentScore.Wins + 1
	}
	if response.Result == LOST {
		currentScore.Losses = currentScore.Losses + 1
	}
	if response.Result == DRAW {
		currentScore.Draws = currentScore.Draws + 1
	}
	return currentScore
}

func (game Game) getScore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, _ := json.Marshal(game.Score)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func main() {
	rand.Seed(time.Now().Unix())
	acceptedHands := []string{"ROCK", "PAPER", "SCISSORS", "SPOCK", "LIZARD"}
	rules := createRules(acceptedHands)
	game := Game{
		AcceptedHands: acceptedHands,
		Rules:         rules,
		Score: &Score{
			Wins:   0,
			Losses: 0,
			Draws:  0,
		},
	}

	router := httprouter.New()
	router.POST("/game", game.play)
	router.GET("/score", game.getScore)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), router))
	return
}

func parseHand(r *http.Request) (Hand, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Hand{Hand: UNKNOWN}, err
	}
	var hand Hand
	err = json.Unmarshal(b, &hand)
	if err != nil {
		return Hand{Hand: UNKNOWN}, err
	}
	return hand, nil
}

func createComputerHand(acceptedHands []string) Hand {
	index := rand.Int31n(int32(len(acceptedHands)))
	hand := acceptedHands[index]
	return Hand{Hand: hand}
}

func evaluateHand(playerHand Hand, computerHand Hand, rules map[string]Rule) Response {
	rule, found := rules[playerHand.Hand]
	if !found {
		return Response{
			Result:       UNKNOWN,
			ComputerHand: nil,
		}
	}
	if playerHand.Hand == computerHand.Hand {
		return Response{
			Result:       DRAW,
			ComputerHand: &computerHand.Hand,
		}
	}
	if contains(rule.WinsAgainst, computerHand.Hand) {
		return Response{
			Result:       WIN,
			ComputerHand: &computerHand.Hand,
		}
	}

	return Response{
		Result:       LOST,
		ComputerHand: &computerHand.Hand,
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/**
Assumes that acceptedHands is a list of hands with and odd number of hands
*/
func createRules(acceptedHands []string) map[string]Rule {
	rules := make(map[string]Rule)
	ruleSize := (len(acceptedHands) - 1) / 2
	for i, hand := range acceptedHands {
		winAgainst := []string{}
		for j := 1; j <= ruleSize; j++ {
			index := (i + (2 * j)) % len(acceptedHands)
			winAgainst = append(winAgainst, acceptedHands[index])
		}

		rules[hand] = Rule{WinsAgainst: winAgainst}
	}

	return rules
}

type Score struct {
	Wins   int `json:"Wins"`
	Losses int `json:"Losses"`
	Draws  int `json:"Draws"`
}

type Rule struct {
	WinsAgainst []string
}

type Hand struct {
	Hand string `json:"hand"`
}

type Response struct {
	Result       string  `json:"result"`
	ComputerHand *string `json:"computerHand,omitempty"`
}
