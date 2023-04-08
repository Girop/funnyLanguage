package lexer

import (
	"fmt"
	"testing"
)

func TestTokenizing(t *testing.T) {
	chars := "fn main() {\n}" 
    fmt.Println(Tokenize(chars))
    
}

func TestComments(t *testing.T) {
	chars := "# This is comment"
    res := Tokenize(chars)
    if len(res) > 0 {
        t.Error("Comment parsed")
    }
}
