package main

import "testing"
import "math/rand"

//Test a 3 Clubs
func TestCalcValSingles(t *testing.T) {
	testName := "CalcVal:Singles:"
	a, b, c := calcVal([]int{1})
	if a != 1 {
		t.Error(testName, "Expected combo to be 1, got", a)
	}
	if b != 3 {
		t.Error(testName, "Expected value to be 3, got", b)
	}
	if c != 1 {
		t.Error(testName, "Expected suit to be 1, got", c)
	}
}

//Test a random single card
func TestCalcValRandSingles(t *testing.T) {
	testName := "CalcVal:RandSingles:"
	card := rand.Intn(51)
	a, b, c := calcVal([]int{card})
	if a != 1 {
		t.Error(testName, "Expected combo to be 1, got", a)
	}
	if b != (card/4)+3 {
		t.Error(testName, "CalcVal:Singles: Expected value to be", (card/4)+3 , "got", b)
	}
	if c != card%4 {
		t.Error(testName, "CalcVal:Singles: Expected suit to be", card%4, "got", c)
	}
}
