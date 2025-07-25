package bash

import (
	"fmt"
	"testing"

	"vimagination.zapto.org/parser"
)

func TestPrintSource(t *testing.T) {
	for n, test := range [...][3]string{
		{ // 1
			"(( a ))",
			"((a));\n",
			"(( a ));\n",
		},
		{ // 2
			"$(( a + b ))",
			"$((a+b));\n",
			"$(( a + b ));\n",
		},
		{ // 3
			"a=(1)",
			"a=(1);\n",
			"a=( 1 );\n",
		},
		{ // 4
			"a=(\n# word comment\nb # post-word comment\n)",
			"a=(\n\t# word comment\n\tb # post-word comment\n);\n",
			"a=(\n\t# word comment\n\tb # post-word comment\n);\n",
		},
		{ // 5
			"a=1",
			"a=1;\n",
			"a=1;\n",
		},
		{ // 6
			"a+=1",
			"a+=1;\n",
			"a+=1;\n",
		},
		{ // 7
			"let a=1+2",
			"let a=1+2;\n",
			"let a=1+2;\n",
		},
		{ // 8
			"let a=1+(2+3)",
			"let a=1+(2+3);\n",
			"let a=1+( 2 + 3 );\n",
		},
		{ // 9
			"let a=b?(c?d:e):f",
			"let a=b?(c?d:e):f;\n",
			"let a=b?( c ? d : e ):f;\n",
		},
		{ // 10
			"case a in b)c\nesac",
			"case a in\nb)\n\tc;;\nesac;\n",
			"case a in\nb)\n\tc;;\nesac;\n",
		},
		{ // 11
			"case a in b);;& c) d;&\nesac",
			"case a in\nb)\n\t;;&\nc)\n\td;&\nesac;\n",
			"case a in\nb)\n\t;;&\nc)\n\td;&\nesac;\n",
		},
		{ // 12
			"case a #A\nin #B\n   #C\nb)c;;\n#D\nesac",
			"case a #A\nin #B\n   #C\nb)\n\tc;;\n#D\nesac;\n",
			"case a #A\nin #B\n   #C\nb)\n\tc;;\n#D\nesac;\n",
		},
		{ // 13
			"a=1 b=2 c d >e <f",
			"a=1 b=2 c d >e <f;\n",
			"a=1 b=2 c d > e < f;\n",
		},
		{ // 14
			"a b >c <d",
			"a b >c <d;\n",
			"a b > c < d;\n",
		},
		{ // 15
			">a <b",
			">a <b;\n",
			"> a < b;\n",
		},
		{ // 16
			"a=1 b=2 >c <d",
			"a=1 b=2 >c <d;\n",
			"a=1 b=2 > c < d;\n",
		},
		{ // 17
			"a <<b\nheredoc\ncontents\nb",
			"a <<b\nheredoc\ncontents\nb\n",
			"a <<b\nheredoc\ncontents\nb\n",
		},
		{ // 18
			"a <<-b\n\theredoc\n\tcontents\nb",
			"a <<-b\nheredoc\ncontents\nb\n",
			"a <<-b\nheredoc\ncontents\nb\n",
		},
		{ // 19
			"a <<ABC; b <<DEF\n123\nABC\n456\nDEF",
			"a <<ABC; b <<DEF;\n123\nABC\n456\nDEF\n",
			"a <<ABC;\n123\nABC\nb <<DEF;\n456\nDEF\n",
		},
		{ // 20
			"$(a)",
			"$(a);\n",
			"$(a);\n",
		},
		{ // 21
			"$(a\nb)",
			"$(\n\ta;\n\tb;\n);\n",
			"$(\n\ta;\n\tb;\n);\n",
		},
		{ // 22
			"# A\na # B\n  # C\n\n# D",
			"# A\n\na; # B\n   # C\n\n# D\n",
			"# A\n\na; # B\n   # C\n\n# D\n",
		},
		{ // 23
			"case a in\nesac < a 2>&1;",
			"case a in\nesac <a 2>&1;\n",
			"case a in\nesac < a 2>&1;\n",
		},
		{ // 24
			"a\nb\n\nc\nd\n\n\n\n\ne",
			"a;\nb;\n\nc;\nd;\n\ne;\n",
			"a;\nb;\n\nc;\nd;\n\ne;\n",
		},
		{ // 25
			"a\nfor b; do\nc\ndone",
			"a;\nfor b; do\n\tc;\ndone;\n",
			"a;\nfor b; do\n\tc;\ndone;\n",
		},
		{ // 26
			"for a;do c\nd\ndone",
			"for a; do\n\tc;\n\td;\ndone;\n",
			"for a; do\n\tc;\n\td;\ndone;\n",
		},
		{ // 27
			"for a in b\ndo c\ndone",
			"for a in b; do\n\tc;\ndone;\n",
			"for a in b; do\n\tc;\ndone;\n",
		},
		{ // 28
			"for a in b c\ndo d\ndone",
			"for a in b c; do\n\td;\ndone;\n",
			"for a in b c; do\n\td;\ndone;\n",
		},
		{ // 29
			"for ((a=0;a<1;a++));do b\ndone",
			"for ((a=0;a<1;a++)); do\n\tb;\ndone;\n",
			"for (( a = 0; a < 1; a ++ )); do\n\tb;\ndone;\n",
		},
		{ // 30
			"function a() { b; }",
			"function a() { b; }\n",
			"function a() {\n\tb;\n}\n",
		},
		{ // 31
			"function a { b; }",
			"function a() { b; }\n",
			"function a() {\n\tb;\n}\n",
		},
		{ // 32
			"a() { b; }",
			"a() { b; }\n",
			"a() {\n\tb;\n}\n",
		},
		{ // 33
			"function a() # A\n# B\n{ b; }",
			"function a() # A\n             # B\n{ b; }\n",
			"function a() # A\n             # B\n{\n\tb;\n}\n",
		},
		{ // 34
			"a() # A\n# B\n{ b; }",
			"a() # A\n    # B\n{ b; }\n",
			"a() # A\n    # B\n{\n\tb;\n}\n",
		},
		{ // 35
			"{ a; }",
			"{ a; }\n",
			"{\n\ta;\n}\n",
		},
		{ // 36
			"( a; )",
			"( a; )\n",
			"(\n\ta;\n)\n",
		},
		{ // 37
			"{ a; b; }",
			"{ a; b; }\n",
			"{\n\ta;\n\tb;\n}\n",
		},
		{ // 38
			"{ a || b; }",
			"{ a || b; }\n",
			"{\n\ta || b;\n}\n",
		},
		{ // 39
			"( a;\nb; )",
			"(\n\ta;\n\tb;\n)\n",
			"(\n\ta;\n\tb;\n)\n",
		},
		{ // 40
			"{ # A\na; # B\n}",
			"{ # A\n\ta; # B\n}\n",
			"{ # A\n\ta; # B\n}\n",
		},
		{ // 41
			"{ a; # A\n}",
			"{\n\ta; # A\n}\n",
			"{\n\ta; # A\n}\n",
		},
		{ // 42
			"{ a | b $(c\nd)\n}",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n}\n",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n}\n",
		},
		{ // 43
			"{ declare a=$(b); }",
			"{ declare a=$(b); }\n",
			"{\n\tdeclare a=$(b);\n}\n",
		},
		{ // 44
			"{ declare a=$(b;c); }",
			"{ declare a=$(b; c;); }\n",
			"{\n\tdeclare a=$(\n\t\tb;\n\t\tc;\n\t);\n}\n",
		},
		{ // 45
			"{ a=1 b; }",
			"{ a=1 b; }\n",
			"{\n\ta=1 b;\n}\n",
		},
		{ // 46
			"{ a=$(b\nc) d; }",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n}\n",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n}\n",
		},
		{ // 47
			"{ let a[$(b)]=c; }",
			"{ let a[$(b)]=c; }\n",
			"{\n\tlet a[$(b)]=c;\n}\n",
		},
		{ // 48
			"{ let a[$(b\nc)]=d; }",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n}\n",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n}\n",
		},
		{ // 49
			"{ let a=$(b); }",
			"{ let a=$(b); }\n",
			"{\n\tlet a=$(b);\n}\n",
		},
		{ // 50
			"{ let a=$(b\nc); }",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n}\n",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n}\n",
		},
		{ // 51
			"{ ${a[$(b\nc)]}; }",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n}\n",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n}\n",
		},
		{ // 52
			"<<a\nb$c\na",
			"<<a\nb$c\na\n",
			"<<a\nb$c\na\n",
		},
		{ // 53
			"{\n<<a\nb$c\na\n}",
			"{\n\t<<a\nb$c\na\n}\n",
			"{\n\t<<a\nb$c\na\n}\n",
		},
		{ // 54
			"{\n<<-a\nb$c\na\n}",
			"{\n\t<<-a\n\tb$c\n\ta\n}\n",
			"{\n\t<<-a\n\tb$c\n\ta\n}\n",
		},
		{ // 55
			"{ function a() { # A\nb; } }",
			"{\n\tfunction a() { # A\n\t\tb;\n\t}\n}\n",
			"{\n\tfunction a() { # A\n\t\tb;\n\t}\n}\n",
		},
		{ // 56
			"{ function a() { b;\nc; } }",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t}\n}\n",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t}\n}\n",
		},
		{ // 57
			"{ if a; then b\nc\nfi; }",
			"{\n\tif a; then\n\t\tb;\n\t\tc;\n\tfi;\n}\n",
			"{\n\tif a; then\n\t\tb;\n\t\tc;\n\tfi;\n}\n",
		},
		{ // 58
			"{ until a; do b\nc\ndone; }",
			"{\n\tuntil a; do\n\t\tb;\n\t\tc;\n\tdone;\n}\n",
			"{\n\tuntil a; do\n\t\tb;\n\t\tc;\n\tdone;\n}\n",
		},
		{ // 59
			"{ case a in b)c;;esac; }",
			"{\n\tcase a in\n\tb)\n\t\tc;;\n\tesac;\n}\n",
			"{\n\tcase a in\n\tb)\n\t\tc;;\n\tesac;\n}\n",
		},
		{ // 60
			"{ for a in b; do c\nd\ndone; }",
			"{\n\tfor a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
			"{\n\tfor a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
		},
		{ // 61
			"{ select a in b; do c\nd\ndone; }",
			"{\n\tselect a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
			"{\n\tselect a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
		},
		{ // 62
			"{ [[ a = b ]]; }",
			"{ [[ a == b ]]; }\n",
			"{\n\t[[ a == b ]];\n}\n",
		},
		{ // 63
			"{ [[ a = $(a\nb) ]]; }",
			"{\n\t[[\n\t\ta == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
			"{\n\t[[\n\t\ta == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
		},
		{ // 64
			"{ [[ a = b || c = $(a\nb) ]]; }",
			"{\n\t[[\n\t\ta == b || c == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
			"{\n\t[[\n\t\ta == b || c == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
		},
		{ // 65
			"{ [[ $(a\nb) ]]; }",
			"{\n\t[[\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
			"{\n\t[[\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
		},
		{ // 66
			"{ ((a)); }",
			"{ ((a)); }\n",
			"{\n\t(( a ));\n}\n",
		},
		{ // 67
			"{ (($(a\nb))); }",
			"{\n\t(($(\n\t\ta;\n\t\tb;\n\t)));\n}\n",
			"{\n\t(( $(\n\t\ta;\n\t\tb;\n\t) ));\n}\n",
		},
		{ // 68
			"{ (($(a\nb))) >&2; }",
			"{\n\t(($(\n\t\ta;\n\t\tb;\n\t))) >&2;\n}\n",
			"{\n\t(( $(\n\t\ta;\n\t\tb;\n\t) )) >&2;\n}\n",
		},
		{ // 69
			"{ ((a)) >$(a\nb); }",
			"{\n\t((a)) >$(\n\t\ta;\n\t\tb;\n\t);\n}\n",
			"{\n\t(( a )) > $(\n\t\ta;\n\t\tb;\n\t);\n}\n",
		},
		{ // 70
			"{ a=( # A\n); }",
			"{\n\ta=( # A\n\t);\n}\n",
			"{\n\ta=( # A\n\t);\n}\n",
		},
		{ // 71
			"{ a=( # A\n# B\n\n# C\n); }",
			"{\n\ta=( # A\n\t    # B\n\n\t# C\n\t);\n}\n",
			"{\n\ta=( # A\n\t    # B\n\n\t# C\n\t);\n}\n",
		},
		{ // 72
			"{ a=( # A\n# B\n\nb\n# C\n); }",
			"{\n\ta=( # A\n\t    # B\n\n\t\tb\n\t# C\n\t);\n}\n",
			"{\n\ta=( # A\n\t    # B\n\n\t\tb\n\t# C\n\t);\n}\n",
		},
		{ // 73
			"{ a=(b); }",
			"{ a=(b); }\n",
			"{\n\ta=( b );\n}\n",
		},
		{ // 74
			"{ a=($(a\nb)); }",
			"{\n\ta=(\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t);\n}\n",
			"{\n\ta=(\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t);\n}\n",
		},
		{ // 75
			"{ $(($(a\nb))); }",
			"{\n\t$(($(\n\t\ta;\n\t\tb;\n\t)));\n}\n",
			"{\n\t$(( $(\n\t\ta;\n\t\tb;\n\t) ));\n}\n",
		},
		{ // 76
			"{ ${a:=$(b)}; }",
			"{ ${a:=$(b)}; }\n",
			"{\n\t${a:=$(b)};\n}\n",
		},
		{ // 77
			"{ ${a:=$(b\nc)}; }",
			"{\n\t${a:=$(\n\t\tb;\n\t\tc;\n\t)};\n}\n",
			"{\n\t${a:=$(\n\t\tb;\n\t\tc;\n\t)};\n}\n",
		},
		{ // 78
			"{ ${a/b/$(c)}; }",
			"{ ${a/b/$(c)}; }\n",
			"{\n\t${a/b/$(c)};\n}\n",
		},
		{ // 79
			"{ ${a/b/$(c\nd)}; }",
			"{\n\t${a/b/$(\n\t\tc;\n\t\td;\n\t)};\n}\n",
			"{\n\t${a/b/$(\n\t\tc;\n\t\td;\n\t)};\n}\n",
		},
		{ // 80
			"{ {a,\"$(b\nc)\"}; }",
			"{\n\t{a,\"$(\n\t\tb;\n\t\tc;\n\t)\"};\n}\n",
			"{\n\t{a,\"$(\n\t\tb;\n\t\tc;\n\t)\"};\n}\n",
		},
		{ // 81
			"{ {a,b}; }",
			"{ {a,b}; }\n",
			"{\n\t{a,b};\n}\n",
		},
		{ // 82
			"a | b <<c\nc",
			"a | b <<c\nc\n",
			"a | b <<c\nc\n",
		},
		{ // 83
			"a && b <<c\nc",
			"a && b <<c\nc\n",
			"a && b <<c\nc\n",
		},
		{ // 84
			"a <<b\nb\n\nc;",
			"a <<b\nb\n\nc;\n",
			"a <<b\nb\n\nc;\n",
		},
		{ // 85
			"if a; then b;fi",
			"if a; then\n\tb;\nfi;\n",
			"if a; then\n\tb;\nfi;\n",
		},
		{ // 86
			"if (a) then b;fi",
			"if ( a; ); then\n\tb;\nfi;\n",
			"if (\n\ta;\n); then\n\tb;\nfi;\n",
		},
		{ // 87
			"if a||b; then b\nc\nelif d\nthen\ne\nelse if f\nthen\ng\nfi\nfi",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;\n",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;\n",
		},
		{ // 88
			"if a; then b;elif c; then d;elif e\nthen f;fi",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;\n",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;\n",
		},
		{ // 89
			"if [ \"$a\" = \"b\" ]; then c;fi",
			"if [ \"$a\" = \"b\" ]; then\n\tc;\nfi;\n",
			"if [ \"$a\" = \"b\" ]; then\n\tc;\nfi;\n",
		},
		{ // 90
			"while a\ndo\nb\ndone",
			"while a; do\n\tb;\ndone;\n",
			"while a; do\n\tb;\ndone;\n",
		},
		{ // 91
			"while (a) do\nb\ndone",
			"while ( a; ); do\n\tb;\ndone;\n",
			"while (\n\ta;\n); do\n\tb;\ndone;\n",
		},
		{ // 92
			"until a&&b; # A\n# B\ndo\n# C\nb\nc\ndone",
			"until a && b; # A\n              # B\ndo\n\t# C\n\tb;\n\tc;\ndone;\n",
			"until a && b; # A\n              # B\ndo\n\t# C\n\tb;\n\tc;\ndone;\n",
		},
		{ // 93
			"a[b]=",
			"a[b]=;\n",
			"a[b]=;\n",
		},
		{ // 94
			"a[b+c]=",
			"a[b+c]=;\n",
			"a[b + c]=;\n",
		},
		{ // 95
			"a[$(b\nc)]=",
			"a[$(\n\tb;\n\tc;\n)]=;\n",
			"a[$(\n\tb;\n\tc;\n)]=;\n",
		},
		{ // 96
			"${a}",
			"${a};\n",
			"${a};\n",
		},
		{ // 97
			"${!a}",
			"${!a};\n",
			"${!a};\n",
		},
		{ // 98
			"${#a}",
			"${#a};\n",
			"${#a};\n",
		},
		{ // 99
			"${a:=b}",
			"${a:=b};\n",
			"${a:=b};\n",
		},
		{ // 100
			"${a=b}",
			"${a=b};\n",
			"${a=b};\n",
		},
		{ // 101
			"${a:?b}",
			"${a:?b};\n",
			"${a:?b};\n",
		},
		{ // 102
			"${a?b}",
			"${a?b};\n",
			"${a?b};\n",
		},
		{ // 103
			"${a:+b}",
			"${a:+b};\n",
			"${a:+b};\n",
		},
		{ // 104
			"${a:+b c}",
			"${a:+b c};\n",
			"${a:+b c};\n",
		},
		{ // 105
			"${a+b}",
			"${a+b};\n",
			"${a+b};\n",
		},
		{ // 106
			"${a:-b}",
			"${a:-b};\n",
			"${a:-b};\n",
		},
		{ // 107
			"${a-b}",
			"${a-b};\n",
			"${a-b};\n",
		},
		{ // 108
			"${a#b}",
			"${a#b};\n",
			"${a#b};\n",
		},
		{ // 109
			"${a##b}",
			"${a##b};\n",
			"${a##b};\n",
		},
		{ // 110
			"${a%b}",
			"${a%b};\n",
			"${a%b};\n",
		},
		{ // 111
			"${a%%b}",
			"${a%%b};\n",
			"${a%%b};\n",
		},
		{ // 112
			"${a/b}",
			"${a/b};\n",
			"${a/b};\n",
		},
		{ // 113
			"${a/b/c}",
			"${a/b/c};\n",
			"${a/b/c};\n",
		},
		{ // 114
			"${a//b}",
			"${a//b};\n",
			"${a//b};\n",
		},
		{ // 115
			"${a//b/c}",
			"${a//b/c};\n",
			"${a//b/c};\n",
		},
		{ // 116
			"${a/%b}",
			"${a/%b};\n",
			"${a/%b};\n",
		},
		{ // 117
			"${a/%b/c}",
			"${a/%b/c};\n",
			"${a/%b/c};\n",
		},
		{ // 118
			"${a/#b}",
			"${a/#b};\n",
			"${a/#b};\n",
		},
		{ // 119
			"${a/#b/c}",
			"${a/#b/c};\n",
			"${a/#b/c};\n",
		},
		{ // 120
			"${a^b}",
			"${a^b};\n",
			"${a^b};\n",
		},
		{ // 121
			"${a^^b}",
			"${a^^b};\n",
			"${a^^b};\n",
		},
		{ // 122
			"${a,b}",
			"${a,b};\n",
			"${a,b};\n",
		},
		{ // 123
			"${a,,b}",
			"${a,,b};\n",
			"${a,,b};\n",
		},
		{ // 124
			"${*}",
			"${*};\n",
			"${*};\n",
		},
		{ // 125
			"${@}",
			"${@};\n",
			"${@};\n",
		},
		{ // 126
			"${a:1}",
			"${a:1};\n",
			"${a:1};\n",
		},
		{ // 127
			"${a:1:2}",
			"${a:1:2};\n",
			"${a:1:2};\n",
		},
		{ // 128
			"${a: -1:2}",
			"${a: -1:2};\n",
			"${a: -1:2};\n",
		},
		{ // 129
			"${a@U}",
			"${a@U};\n",
			"${a@U};\n",
		},
		{ // 130
			"${a@u}",
			"${a@u};\n",
			"${a@u};\n",
		},
		{ // 131
			"${a@L}",
			"${a@L};\n",
			"${a@L};\n",
		},
		{ // 132
			"${a@Q}",
			"${a@Q};\n",
			"${a@Q};\n",
		},
		{ // 133
			"${a@E}",
			"${a@E};\n",
			"${a@E};\n",
		},
		{ // 134
			"${a@P}",
			"${a@P};\n",
			"${a@P};\n",
		},
		{ // 135
			"${a@A}",
			"${a@A};\n",
			"${a@A};\n",
		},
		{ // 136
			"${a@K}",
			"${a@K};\n",
			"${a@K};\n",
		},
		{ // 137
			"${a@a}",
			"${a@a};\n",
			"${a@a};\n",
		},
		{ // 138
			"${a@k}",
			"${a@k};\n",
			"${a@k};\n",
		},
		{ // 139
			"${!a@}",
			"${!a@};\n",
			"${!a@};\n",
		},
		{ // 140
			"${!a*}",
			"${!a*};\n",
			"${!a*};\n",
		},
		{ // 141
			"[[ a == b ]]",
			"[[ a == b ]];\n",
			"[[ a == b ]];\n",
		},
		{ // 142
			"[[ a == b$c ]]",
			"[[ a == b$c ]];\n",
			"[[ a == b$c ]];\n",
		},
		{ // 143
			"[[ a == b\"c\" ]]",
			"[[ a == b\"c\" ]];\n",
			"[[ a == b\"c\" ]];\n",
		},
		{ // 144
			"case a in a|b) a;\nb\nesac",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;\n",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;\n",
		},
		{ // 145
			"case a in a|b|\"c\") a;\nb\nesac",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;\n",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;\n",
		},
		{ // 146
			"! a",
			"! a;\n",
			"! a;\n",
		},
		{ // 147
			"coproc a",
			"coproc a;\n",
			"coproc a;\n",
		},
		{ // 148
			"! coproc a",
			"! coproc a;\n",
			"! coproc a;\n",
		},
		{ // 149
			"time a",
			"time a;\n",
			"time a;\n",
		},
		{ // 150
			"time -p a",
			"time -p a;\n",
			"time -p a;\n",
		},
		{ // 151
			"time coproc a",
			"time coproc a;\n",
			"time coproc a;\n",
		},
		{ // 152
			"time -p coproc a",
			"time -p coproc a;\n",
			"time -p coproc a;\n",
		},
		{ // 153
			"time ! coproc a",
			"time ! coproc a;\n",
			"time ! coproc a;\n",
		},
		{ // 154
			"time -p ! coproc a",
			"time -p ! coproc a;\n",
			"time -p ! coproc a;\n",
		},
		{ // 155
			"coproc a if b; then c\nfi",
			"coproc a if b; then\n\tc;\nfi;\n",
			"coproc a if b; then\n\tc;\nfi;\n",
		},
		{ // 156
			"select a; do b; done",
			"select a; do\n\tb;\ndone;\n",
			"select a; do\n\tb;\ndone;\n",
		},
		{ // 157
			"select a in b c; do b; done",
			"select a in b c; do\n\tb;\ndone;\n",
			"select a in b c; do\n\tb;\ndone;\n",
		},
		{ // 158
			"a&",
			"a&\n",
			"a &\n",
		},
		{ // 159
			"[[ # A\na == b\n# B\n]]",
			"[[ # A\n\ta == b\n# B\n]];\n",
			"[[ # A\n\ta == b\n# B\n]];\n",
		},
		{ // 160
			"[[ # A\na == b\n]]",
			"[[ # A\n\ta == b\n]];\n",
			"[[ # A\n\ta == b\n]];\n",
		},
		{ // 161
			"[[\n\t! # A\n\ta -ot b ]]",
			"[[\n\t! # A\n\ta -ot b\n]];\n",
			"[[\n\t! # A\n\ta -ot b\n]];\n",
		},
		{ // 162
			"[[\n\t! # A\n\ta -nt b ]]",
			"[[\n\t! # A\n\ta -nt b\n]];\n",
			"[[\n\t! # A\n\ta -nt b\n]];\n",
		},
		{ // 163
			"[[ (a -ef b) ]]",
			"[[ ( a -ef b ) ]];\n",
			"[[ ( a -ef b ) ]];\n",
		},
		{ // 164
			"[[ (# A\na -ge b) ]]",
			"[[\n\t( # A\n\t\ta -ge b\n\t)\n]];\n",
			"[[\n\t( # A\n\t\ta -ge b\n\t)\n]];\n",
		},
		{ // 165
			"[[ (a -gt b\n# A\n) ]]",
			"[[\n\t(\n\t\ta -gt b\n\t# A\n\t)\n]];\n",
			"[[\n\t(\n\t\ta -gt b\n\t# A\n\t)\n]];\n",
		},
		{ // 166
			"[[ (\n# A\na -le b # B\n) ]]",
			"[[\n\t(\n\t\t# A\n\t\ta -le b # B\n\t)\n]];\n",
			"[[\n\t(\n\t\t# A\n\t\ta -le b # B\n\t)\n]];\n",
		},
		{ // 167
			"[[\n# A\na =~ b #B\n]]",
			"[[\n\t# A\n\ta =~ b #B\n]];\n",
			"[[\n\t# A\n\ta =~ b #B\n]];\n",
		},
		{ // 168
			"[[ # A\n\n# B\na != b # C\n\n# D\n]]",
			"[[ # A\n\n\t# B\n\ta != b # C\n\n# D\n]];\n",
			"[[ # A\n\n\t# B\n\ta != b # C\n\n# D\n]];\n",
		},
		{ // 169
			"[[ # A\n\n# B\na < b # C\n||# D\n\n# E\nd>e # F\n\n# G\n]]",
			"[[ # A\n\n\t# B\n\ta < b # C\n\t|| # D\n\n\t   # E\n\td > e # F\n\n# G\n]];\n",
			"[[ # A\n\n\t# B\n\ta < b # C\n\t|| # D\n\n\t   # E\n\td > e # F\n\n# G\n]];\n",
		},
		{ // 170
			"[[ # A\n# B\n\n# C\n# D\n( # E\n# F\n\n# G\n# H\na -eq b # I\n# J\n&& # K\n# L\nc -ne d # M\n# N\n\n# O\n\n# P\n) # Q\n# R\n\n# S\n# T\n&& # U\n# V\n\n# W\n# X\ne -lt f # Y\n# Z\n]]",
			"[[ # A\n   # B\n\n\t# C\n\t# D\n\t( # E\n\t  # F\n\n\t\t# G\n\t\t# H\n\t\ta -eq b # I\n\t\t        # J\n\t\t&& # K\n\t\t   # L\n\t\tc -ne d # M\n\t\t        # N\n\n\t# O\n\n\t# P\n\t) # Q\n\t  # R\n\n\t  # S\n\t  # T\n\t&& # U\n\t   # V\n\n\t   # W\n\t   # X\n\te -lt f # Y\n\t        # Z\n]];\n",
			"[[ # A\n   # B\n\n\t# C\n\t# D\n\t( # E\n\t  # F\n\n\t\t# G\n\t\t# H\n\t\ta -eq b # I\n\t\t        # J\n\t\t&& # K\n\t\t   # L\n\t\tc -ne d # M\n\t\t        # N\n\n\t# O\n\n\t# P\n\t) # Q\n\t  # R\n\n\t  # S\n\t  # T\n\t&& # U\n\t   # V\n\n\t   # W\n\t   # X\n\te -lt f # Y\n\t        # Z\n]];\n",
		},
		{ // 171
			"[[ -e a ]]",
			"[[ -e a ]];\n",
			"[[ -e a ]];\n",
		},
		{ // 172
			"[[ -b a || -c b ]]",
			"[[ -b a || -c b ]];\n",
			"[[ -b a || -c b ]];\n",
		},
		{ // 173
			"[[ -d a && -f b ]]",
			"[[ -d a && -f b ]];\n",
			"[[ -d a && -f b ]];\n",
		},
		{ // 174
			"[[ (-g a && -L b) || -k c ]]",
			"[[ ( -g a && -L b ) || -k c ]];\n",
			"[[ ( -g a && -L b ) || -k c ]];\n",
		},
		{ // 175
			"[[ (-p a && -r b) || (-s c && (-t d || -u e)) ]]",
			"[[ ( -p a && -r b ) || ( -s c && ( -t d || -u e ) ) ]];\n",
			"[[ ( -p a && -r b ) || ( -s c && ( -t d || -u e ) ) ]];\n",
		},
		{ // 176
			"[[ (-w a && -x b) || (-G c && (-N d || -O e)) ]]",
			"[[ ( -w a && -x b ) || ( -G c && ( -N d || -O e ) ) ]];\n",
			"[[ ( -w a && -x b ) || ( -G c && ( -N d || -O e ) ) ]];\n",
		},
		{ // 177
			"[[ (-S a && -o b) || (-v c && (-R d || -z e || -n f)) ]]",
			"[[ ( -S a && -o b ) || ( -v c && ( -R d || -z e || -n f ) ) ]];\n",
			"[[ ( -S a && -o b ) || ( -v c && ( -R d || -z e || -n f ) ) ]];\n",
		},
		{ // 178
			"if a; then b; fi",
			"if a; then\n\tb;\nfi;\n",
			"if a; then\n\tb;\nfi;\n",
		},
		{ // 179
			"if a # A\nthen b; fi",
			"if a; # A\nthen\n\tb;\nfi;\n",
			"if a; # A\nthen\n\tb;\nfi;\n",
		},
		{ // 180
			"if a # A\n# B\nthen b; fi",
			"if a; # A\n      # B\nthen\n\tb;\nfi;\n",
			"if a; # A\n      # B\nthen\n\tb;\nfi;\n",
		},
		{ // 181
			"${a/b/ }",
			"${a/b/ };\n",
			"${a/b/ };\n",
		},
		{ // 182
			"<(a)",
			"<(a);\n",
			"<(a);\n",
		},
		{ // 183
			">(a;b)",
			">(a; b;);\n",
			">(\n\ta;\n\tb;\n);\n",
		},
		{ // 184
			"`a`",
			"`a`;\n",
			"`a`;\n",
		},
		{ // 185
			"`a``b`",
			"`a``b`;\n",
			"`a``b`;\n",
		},
		{ // 186
			"`a \\`b\\``",
			"`a \\`b\\``;\n",
			"`a \\`b\\``;\n",
		},
		{ // 187
			"while read a; do\n\tb \"$a\";\ndone <<EOF\nA\nB\nC\nEOF",
			"while read a; do\n\tb \"$a\";\ndone <<EOF\nA\nB\nC\nEOF\n",
			"while read a; do\n\tb \"$a\";\ndone <<EOF\nA\nB\nC\nEOF\n",
		},
		{ // 188
			"while read a; do\n\tb \"$a\";\ndone < <(c)",
			"while read a; do\n\tb \"$a\";\ndone < <(c);\n",
			"while read a; do\n\tb \"$a\";\ndone < <(c);\n",
		},
		{ // 189
			"{a,b}",
			"{a,b};\n",
			"{a,b};\n",
		},
		{ // 190
			"{a,,bcd,\"\"}",
			"{a,,bcd,\"\"};\n",
			"{a,,bcd,\"\"};\n",
		},
		{ // 191
			"{1..2}",
			"{1..2};\n",
			"{1..2};\n",
		},
		{ // 192
			"{-10..-100..-2}",
			"{-10..-100..-2};\n",
			"{-10..-100..-2};\n",
		},
		{ // 193
			"{a..z}",
			"{a..z};\n",
			"{a..z};\n",
		},
		{ // 194
			"{A..Z..3}",
			"{A..Z..3};\n",
			"{A..Z..3};\n",
		},
		{ // 195
			"if a; then\n\ta=\"\n\";fi",
			"if a; then\n\ta=\"\n\";\nfi;\n",
			"if a; then\n\ta=\"\n\";\nfi;\n",
		},
	} {
		for m, input := range test {
			if m == 2 && (n == 18 || n == 43 || n == 36 || n == 182) {
				continue
			}

			tk := parser.NewStringTokeniser(input)

			if f, err := Parse(&tk); err != nil {
				t.Errorf("test %d.%d: unexpected error: %s", n+1, m+1, err)
			} else if simple := fmt.Sprintf("%s", f); simple != test[1] {
				t.Errorf("test %d.%d.1: expecting output %q, got %q", n+1, m+1, test[1], simple)
			} else if verbose := fmt.Sprintf("%+s", f); verbose != test[2] {
				t.Errorf("test %d.%d.2: expecting output %q, got %q", n+1, m+1, test[2], verbose)
			}
		}
	}
}

func TestLineSource(t *testing.T) {
	l := Line{
		Statements: []Statement{
			{
				Pipeline: Pipeline{
					CommandOrCompound: CommandOrCompound{
						Command: &Command{
							AssignmentsOrWords: []AssignmentOrWord{
								{
									Word: &Word{
										Parts: []WordPart{
											{
												Part: &Token{
													Token: parser.Token{
														Data: "a",
													},
												},
											},
										},
									},
								},
							},
							Redirections: []Redirection{
								{
									Redirector: &Token{
										Token: parser.Token{
											Data: "<",
										},
									},
									Output: Word{
										Parts: []WordPart{
											{
												Part: &Token{
													Token: parser.Token{
														Data: "b",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Comments: [2]Comments{nil, {
			Token{
				Token: parser.Token{
					Data: "A",
				},
			},
		}},
	}

	const expected = "a < b; #A"

	if got := fmt.Sprintf("%+s", l); got != expected {
		t.Errorf("test: expecting output %q, got %q", expected, got)
	}
}
