package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	ans := calculate("3+5/(5-4)")
	fmt.Println(ans)
}

func calculate(s string) int {
	tokens := tokenize(&s)
	fmt.Printf("%v\n", tokenize(&s))
	return NewScanner(tokens).parseTerm()
}

type Scanner struct {
	tokens []string
	ind    int
}

func NewScanner(tokens []string) *Scanner {
	return &Scanner{
		tokens: append(tokens, "#"),
		ind:    0,
	}
}

func (s *Scanner) parseTerm() int {
	val := s.parseFactor()

loop:
	for {
		switch s.lookup() {
		case "+":
			s.consume()
			val += s.parseFactor()
		case "-":
			s.consume()
			val -= s.parseFactor()
		default:
			break loop
		}
	}
	return val
}

func (s *Scanner) parseFactor() int {
	val := s.parseItem()

loop:
	for {
		switch s.lookup() {
		case "*":
			s.consume()
			val *= s.parseItem()
		case "/":
			s.consume()
			val /= s.parseItem()
		default:
			break loop
		}
	}
	return val
}

func (s *Scanner) parseItem() int {
	fmt.Printf("%+v %v %v\n", s, s.ind, s.lookup())

	if s.lookup() == "(" {
		s.consume()
		val := s.parseTerm()
		if s.lookup() == ")" {
			s.consume()
		} else {
			panic(fmt.Sprintf("Expect ), but got %v in ind %v", s.lookup(), s.ind))
		}
		return val
	} else {
		val, err := strconv.Atoi(s.lookup())
		if err != nil {
			panic(fmt.Sprintf("Expect digits, but got %v in ind %v", s.lookup(), s.ind))
		}
		s.consume()
		return val
	}
}

func (s *Scanner) lookup() string {
	return s.tokens[s.ind]
}

func (s *Scanner) consume() {
	s.ind++
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

	return strings.Fields(strings.Join(tokens, ""))
}
