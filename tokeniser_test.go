package bash

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestTokeniser(t *testing.T) {
	for n, test := range [...]struct {
		Input  string
		Output []parser.Token
	}{
		{ // 1
			"",
			[]parser.Token{
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 2
			" ",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 3
			" \t",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " \t"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 4
			" \n\n \n",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenLineTerminator, Data: "\n\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 5
			"#A comment\n# B comment",
			[]parser.Token{
				{Type: TokenComment, Data: "#A comment"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenComment, Data: "# B comment"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 6
			"$ident $name abc=a",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "$ident"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "$name"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "abc"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWord, Data: "a"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 7
			"if then else elif fi case esac while for in do done time until coproc select function",
			[]parser.Token{
				{Type: TokenKeyword, Data: "if"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "then"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "else"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "elif"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "fi"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "case"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "esac"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "while"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "for"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "in"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "do"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "done"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "time"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "until"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "coproc"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "select"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "function"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 8
			"ident ${name} ab\\nc=a",
			[]parser.Token{
				{Type: TokenWord, Data: "ident"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "${"},
				{Type: TokenWord, Data: "name"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "ab\\nc"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWord, Data: "a"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 9
			"$(( 0 1 29 0xff 0xDeAdBeEf 0755 2#5 ))",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "$(("},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumberLiteral, Data: "0"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumberLiteral, Data: "1"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumberLiteral, Data: "29"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumberLiteral, Data: "0xff"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumberLiteral, Data: "0xDeAdBeEf"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumberLiteral, Data: "0755"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumberLiteral, Data: "2#5"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "))"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 10
			"\"abc\" \"de\\nf\" \"stuff`command`more stuff\" \"text $ident end\" \"text $(command) end\" \"text ${ident} end\"",
			[]parser.Token{
				{Type: TokenString, Data: "\"abc\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"de\\nf\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"stuff"},
				{Type: TokenPunctuator, Data: "`"},
				{Type: TokenWord, Data: "command"},
				{Type: TokenPunctuator, Data: "`"},
				{Type: TokenString, Data: "more stuff\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"text "},
				{Type: TokenIdentifier, Data: "$ident"},
				{Type: TokenString, Data: " end\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"text "},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenWord, Data: "command"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenString, Data: " end\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"text "},
				{Type: TokenPunctuator, Data: "${"},
				{Type: TokenWord, Data: "ident"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenString, Data: " end\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
	} {
		p := parser.NewStringTokeniser(test.Input)

		SetTokeniser(&p)

		for m, tkn := range test.Output {
			tk, _ := p.GetToken()
			if tk.Type != tkn.Type {
				if tk.Type == parser.TokenError {
					t.Errorf("test %d.%d: unexpected error: %s", n+1, m+1, tk.Data)
				} else {
					t.Errorf("test %d.%d: Incorrect type, expecting %d, got %d", n+1, m+1, tkn.Type, tk.Type)
				}

				break
			} else if tk.Data != tkn.Data {
				t.Errorf("test %d.%d: Incorrect data, expecting %q, got %q", n+1, m+1, tkn.Data, tk.Data)

				break
			}
		}
	}
}
