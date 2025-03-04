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
			" \t\\\n",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " \t\\\n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 4
			"\\\n \t",
			[]parser.Token{
				{Type: TokenWhitespace, Data: "\\\n \t"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 5
			" \n\n \n",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenLineTerminator, Data: "\n\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 6
			"#A comment\n# B comment",
			[]parser.Token{
				{Type: TokenComment, Data: "#A comment"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenComment, Data: "# B comment"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 7
			"$ident $name a\\nbc=a $0 $12 a$b a${b}c",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "$ident"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "$name"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "a\\nbc=a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "$0"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "$1"},
				{Type: TokenWord, Data: "2"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "a"},
				{Type: TokenIdentifier, Data: "$b"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "a"},
				{Type: TokenPunctuator, Data: "${"},
				{Type: TokenWord, Data: "b"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenWord, Data: "c"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 8
			"abc=a def[0]=b ghi[$i]=c jkl+=d",
			[]parser.Token{
				{Type: TokenIdentifierAssign, Data: "abc"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWord, Data: "a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifierAssign, Data: "def"},
				{Type: TokenPunctuator, Data: "["},
				{Type: TokenWord, Data: "0"},
				{Type: TokenPunctuator, Data: "]"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWord, Data: "b"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifierAssign, Data: "ghi"},
				{Type: TokenPunctuator, Data: "["},
				{Type: TokenIdentifier, Data: "$i"},
				{Type: TokenPunctuator, Data: "]"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWord, Data: "c"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifierAssign, Data: "jkl"},
				{Type: TokenPunctuator, Data: "+="},
				{Type: TokenWord, Data: "d"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 9
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
		{ // 10
			"ident ${name} ab\\nc=a ${6} a$ ",
			[]parser.Token{
				{Type: TokenWord, Data: "ident"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "${"},
				{Type: TokenWord, Data: "name"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "ab\\nc=a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "${"},
				{Type: TokenWord, Data: "6"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "a$"},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 11
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
		{ // 12
			"\"abc\" \"de\\nf\" \"stuff`command`more stuff\" \"text $ident $another end\" \"text $(command) end - text ${ident} end\"",
			[]parser.Token{
				{Type: TokenString, Data: "\"abc\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"de\\nf\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenStringStart, Data: "\"stuff"},
				{Type: TokenPunctuator, Data: "`"},
				{Type: TokenWord, Data: "command"},
				{Type: TokenPunctuator, Data: "`"},
				{Type: TokenStringEnd, Data: "more stuff\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenStringStart, Data: "\"text "},
				{Type: TokenIdentifier, Data: "$ident"},
				{Type: TokenStringMid, Data: " "},
				{Type: TokenIdentifier, Data: "$another"},
				{Type: TokenStringEnd, Data: " end\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenStringStart, Data: "\"text "},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenWord, Data: "command"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenStringMid, Data: " end - text "},
				{Type: TokenPunctuator, Data: "${"},
				{Type: TokenWord, Data: "ident"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenStringEnd, Data: " end\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 13
			"< <<< <<- <& <> > >> >& &>> >| | |& || & && () {} = `` $() $(())",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "<<<"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "<<-"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "<&"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "<>"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ">"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ">>"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ">&"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "&>>"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ">|"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "|"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "|&"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "||"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "&"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "&&"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "`"},
				{Type: TokenPunctuator, Data: "`"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "$(("},
				{Type: TokenPunctuator, Data: "))"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 14
			"$(( + += - -= & &= | |= < <= > >= = == ! != * *= / /= % %= ^ ^= ~ ? : , ))",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "$(("},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "+"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "+="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "-"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "-="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "&"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "&="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "|"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "|="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "<"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "<="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ">"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ">="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "=="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "!"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "!="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "*"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "*="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "/"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "/="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "%"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "%="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "^"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "^="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "~"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "?"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ":"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "))"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 15
			"{ ",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 16
			"{ )",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenError, Data: "invalid character"},
			},
		},
		{ // 17
			"(",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 18
			"$(",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "$("},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 19
			"$(}",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "$("},
				{Type: parser.TokenError, Data: "invalid character"},
			},
		},
		{ // 20
			"<<abc\n123\n456\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "123\n456\n"},
				{Type: TokenHeredocEnd, Data: "abc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 21
			"<<a'b 'c\n123\n456\nab c\n",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "a'b 'c"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "123\n456\n"},
				{Type: TokenHeredocEnd, Data: "ab c"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 22
			"<<def\n123\n456\ndef\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "123\n456\n"},
				{Type: TokenHeredocEnd, Data: "def"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWord, Data: "abc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 23
			"<<def cat\n123\n456\ndef\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cat"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "123\n456\n"},
				{Type: TokenHeredocEnd, Data: "def"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWord, Data: "abc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 24
			"<<abc cat;<<def cat\n123\nabc\n456\ndef",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cat"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cat"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "123\n"},
				{Type: TokenHeredocEnd, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "456\n"},
				{Type: TokenHeredocEnd, Data: "def"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 25
			"<<abc cat;echo $(<<def cat\n456\ndef\n)\n123\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cat"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenWord, Data: "echo"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cat"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "456\n"},
				{Type: TokenHeredocEnd, Data: "def"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "123\n"},
				{Type: TokenHeredocEnd, Data: "abc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 26
			"<<abc\na$abc\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "a"},
				{Type: TokenIdentifier, Data: "$abc"},
				{Type: TokenHeredoc, Data: "\n"},
				{Type: TokenHeredocEnd, Data: "abc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 27
			"<<abc\na${abc} $99\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "a"},
				{Type: TokenPunctuator, Data: "${"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenPunctuator, Data: "}"},
				{Type: TokenHeredoc, Data: " "},
				{Type: TokenIdentifier, Data: "$9"},
				{Type: TokenHeredoc, Data: "9\n"},
				{Type: TokenHeredocEnd, Data: "abc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 28
			"<<abc\na$(\necho abc;\n) 1\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "a"},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWord, Data: "echo"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenHeredoc, Data: " 1\n"},
				{Type: TokenHeredocEnd, Data: "abc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 29
			"<<abc\na$(<<def) 1\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "a"},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: parser.TokenError, Data: "invalid character"},
			},
		},
		{ // 30
			"<<abc\na$(<<def cat) 1\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "a"},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cat"},
				{Type: parser.TokenError, Data: "invalid character"},
			},
		},
		{ // 31
			"<<abc;$(<<def cat)\nabc\ndef\nabc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cat"},
				{Type: parser.TokenError, Data: "invalid character"},
			},
		},
		{ // 32
			"<<abc;<<def;$(<<ghi;<<jkl\nghi\njkl\n)\nabc\ndef",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "abc"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "def"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenPunctuator, Data: "$("},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "ghi"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "jkl"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: ""},
				{Type: TokenHeredocEnd, Data: "ghi"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: ""},
				{Type: TokenHeredocEnd, Data: "jkl"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: ""},
				{Type: TokenHeredocEnd, Data: "abc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: ""},
				{Type: TokenHeredocEnd, Data: "def"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 33
			"<<a\\\nbc\nabc\ndef\na\nbc",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "<<"},
				{Type: TokenWord, Data: "a\\\nbc"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenHeredoc, Data: "abc\ndef\n"},
				{Type: TokenHeredocEnd, Data: "a\nbc"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 34
			"2>1 word",
			[]parser.Token{
				{Type: TokenWord, Data: "2"},
				{Type: TokenPunctuator, Data: ">"},
				{Type: TokenWord, Data: "1"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "word"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 35
			"time -p cmd",
			[]parser.Token{
				{Type: TokenKeyword, Data: "time"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "-p"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "cmd"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 36
			"{a..b..2} {a,b,d} a{b,c,d}e a{1..4} {2..10..-1} {-1..-100..5} {a..z..-1}",
			[]parser.Token{
				{Type: TokenBraceExpansion, Data: "{a..b..2}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBraceExpansion, Data: "{a,b,d}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "a"},
				{Type: TokenBraceExpansion, Data: "{b,c,d}"},
				{Type: TokenWord, Data: "e"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "a"},
				{Type: TokenBraceExpansion, Data: "{1..4}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBraceExpansion, Data: "{2..10..-1}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBraceExpansion, Data: "{-1..-100..5}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBraceExpansion, Data: "{a..z..-1}"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 37
			"a={123",
			[]parser.Token{
				{Type: TokenIdentifierAssign, Data: "a"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWord, Data: "{123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 37
			"word{ word{a} word{\nword{",
			[]parser.Token{
				{Type: TokenWord, Data: "word{"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "word{a}"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "word{"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWord, Data: "word{"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 38
			"{ echo 123; echo 456; }",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "echo"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "123"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "echo"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "456"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "}"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 39
			"(echo 123; echo 456)",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenWord, Data: "echo"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "123"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "echo"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenWord, Data: "456"},
				{Type: TokenPunctuator, Data: ")"},
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
