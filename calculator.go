package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	calculate("3+5/5*4/(4-32)")
}

func calculate(s string) int {
	fmt.Printf("%v\n", tokenize(&s))
	return 0
}

type Scanner struct {
	tokens []string
	ind    int
}

func NewScanner(tokens []string) *Scanner {
	return &Scanner{
		tokens: tokens,
		ind:    0,
	}
}

func (s *Scanner) parseTerm() int {
	if s[ind] == 
}

func (s *Scanner) parseFactor() int {

}

func (s *Scanner) parseItem() int {
	
}

func tokenize(s *string) []string {
	tokens := make([]string, 0)
	for _, r := range *s {
		if unicode.IsDigit(r) {
			tokens = append(tokens, string(r))
		} else {
			tokens = append(tokens, " "+string(r)+" ")
		}
	}

	return strings.Split(strings.Join(tokens, ""), " ")
}
