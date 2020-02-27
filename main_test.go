package main

import (
	"testing"
)

func Test_createRulesRPS(t *testing.T) {
	acceptedHands := []string{"ROCK", "PAPER", "SCISSORS"}
	rules := createRules(acceptedHands)

	for hand, rule := range rules {
		if len(rule.WinsAgainst) != 1 {
			t.Errorf("Expected all rules to have the same amount of wins against rules, but %s had %d rules", hand, len(rule.WinsAgainst))
		}
	}

}

func Test_createRulesRPSLS(t *testing.T) {
	acceptedHands := []string{"ROCK", "PAPER", "SCISSORS", "LIZARD", "SPOCK"}
	rules := createRules(acceptedHands)

	for hand, rule := range rules {
		if len(rule.WinsAgainst) != 2 {
			t.Errorf("Expected all rules to have the same amount of wins against rules, but %s had %d rules", hand, len(rule.WinsAgainst))
		}
	}
}
