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
			"$(a)",
			"$(a);\n",
			"$(a);\n",
		},
		{ // 20
			"$(a\nb)",
			"$(\n\ta;\n\tb;\n);\n",
			"$(\n\ta;\n\tb;\n);\n",
		},
		{ // 21
			"# A\na # B\n  # C\n\n# D",
			"# A\n\na; # B\n   # C\n\n# D\n",
			"# A\n\na; # B\n   # C\n\n# D\n",
		},
		{ // 22
			"case a in\nesac < a 2>&1;",
			"case a in\nesac <a 2>&1;\n",
			"case a in\nesac < a 2>&1;\n",
		},
		{ // 23
			"a\nb\n\nc\nd\n\n\n\n\ne",
			"a;\nb;\n\nc;\nd;\n\ne;\n",
			"a;\nb;\n\nc;\nd;\n\ne;\n",
		},
		{ // 24
			"a\nfor b; do\nc\ndone",
			"a;\nfor b; do\n\tc;\ndone;\n",
			"a;\nfor b; do\n\tc;\ndone;\n",
		},
		{ // 25
			"for a;do c\nd\ndone",
			"for a; do\n\tc;\n\td;\ndone;\n",
			"for a; do\n\tc;\n\td;\ndone;\n",
		},
		{ // 26
			"for a in b\ndo c\ndone",
			"for a in b; do\n\tc;\ndone;\n",
			"for a in b; do\n\tc;\ndone;\n",
		},
		{ // 27
			"for a in b c\ndo d\ndone",
			"for a in b c; do\n\td;\ndone;\n",
			"for a in b c; do\n\td;\ndone;\n",
		},
		{ // 28
			"for ((a=0;a<1;a++));do b\ndone",
			"for ((a=0;a<1;a++)); do\n\tb;\ndone;\n",
			"for (( a = 0; a < 1; a ++ )); do\n\tb;\ndone;\n",
		},
		{ // 29
			"function a() { b; }",
			"function a() { b; }\n",
			"function a() {\n\tb;\n}\n",
		},
		{ // 30
			"function a { b; }",
			"function a() { b; }\n",
			"function a() {\n\tb;\n}\n",
		},
		{ // 31
			"a() { b; }",
			"a() { b; }\n",
			"a() {\n\tb;\n}\n",
		},
		{ // 32
			"function a() # A\n# B\n{ b; }",
			"function a() # A\n             # B\n{ b; }\n",
			"function a() # A\n             # B\n{\n\tb;\n}\n",
		},
		{ // 33
			"a() # A\n# B\n{ b; }",
			"a() # A\n    # B\n{ b; }\n",
			"a() # A\n    # B\n{\n\tb;\n}\n",
		},
		{ // 34
			"{ a; }",
			"{ a; }\n",
			"{\n\ta;\n}\n",
		},
		{ // 35
			"( a; )",
			"( a; )\n",
			"(\n\ta;\n)\n",
		},
		{ // 36
			"{ a; b; }",
			"{ a; b; }\n",
			"{\n\ta;\n\tb;\n}\n",
		},
		{ // 37
			"{ a || b; }",
			"{ a || b; }\n",
			"{\n\ta || b;\n}\n",
		},
		{ // 38
			"( a;\nb; )",
			"(\n\ta;\n\tb;\n)\n",
			"(\n\ta;\n\tb;\n)\n",
		},
		{ // 39
			"{ # A\na; # B\n}",
			"{ # A\n\ta; # B\n}\n",
			"{ # A\n\ta; # B\n}\n",
		},
		{ // 40
			"{ a; # A\n}",
			"{\n\ta; # A\n}\n",
			"{\n\ta; # A\n}\n",
		},
		{ // 41
			"{ a | b $(c\nd)\n}",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n}\n",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n}\n",
		},
		{ // 42
			"{ declare a=$(b); }",
			"{ declare a=$(b); }\n",
			"{\n\tdeclare a=$(b);\n}\n",
		},
		{ // 43
			"{ declare a=$(b;c); }",
			"{ declare a=$(b; c;); }\n",
			"{\n\tdeclare a=$(\n\t\tb;\n\t\tc;\n\t);\n}\n",
		},
		{ // 44
			"{ a=1 b; }",
			"{ a=1 b; }\n",
			"{\n\ta=1 b;\n}\n",
		},
		{ // 45
			"{ a=$(b\nc) d; }",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n}\n",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n}\n",
		},
		{ // 46
			"{ let a[$(b)]=c; }",
			"{ let a[$(b)]=c; }\n",
			"{\n\tlet a[$(b)]=c;\n}\n",
		},
		{ // 47
			"{ let a[$(b\nc)]=d; }",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n}\n",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n}\n",
		},
		{ // 48
			"{ let a=$(b); }",
			"{ let a=$(b); }\n",
			"{\n\tlet a=$(b);\n}\n",
		},
		{ // 49
			"{ let a=$(b\nc); }",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n}\n",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n}\n",
		},
		{ // 50
			"{ ${a[$(b\nc)]}; }",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n}\n",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n}\n",
		},
		{ // 51
			"<<a\nb$c\na",
			"<<a\nb$c\na\n",
			"<<a\nb$c\na\n",
		},
		{ // 52
			"{\n<<a\nb$c\na\n}",
			"{\n\t<<a\nb$c\na\n}\n",
			"{\n\t<<a\nb$c\na\n}\n",
		},
		{ // 53
			"{\n<<-a\nb$c\na\n}",
			"{\n\t<<-a\n\tb$c\n\ta\n}\n",
			"{\n\t<<-a\n\tb$c\n\ta\n}\n",
		},
		{ // 54
			"{ function a() { # A\nb; } }",
			"{\n\tfunction a() { # A\n\t\tb;\n\t}\n}\n",
			"{\n\tfunction a() { # A\n\t\tb;\n\t}\n}\n",
		},
		{ // 55
			"{ function a() { b;\nc; } }",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t}\n}\n",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t}\n}\n",
		},
		{ // 56
			"{ if a; then b\nc\nfi; }",
			"{\n\tif a; then\n\t\tb;\n\t\tc;\n\tfi;\n}\n",
			"{\n\tif a; then\n\t\tb;\n\t\tc;\n\tfi;\n}\n",
		},
		{ // 57
			"{ until a; do b\nc\ndone; }",
			"{\n\tuntil a; do\n\t\tb;\n\t\tc;\n\tdone;\n}\n",
			"{\n\tuntil a; do\n\t\tb;\n\t\tc;\n\tdone;\n}\n",
		},
		{ // 58
			"{ case a in b)c;;esac; }",
			"{\n\tcase a in\n\tb)\n\t\tc;;\n\tesac;\n}\n",
			"{\n\tcase a in\n\tb)\n\t\tc;;\n\tesac;\n}\n",
		},
		{ // 59
			"{ for a in b; do c\nd\ndone; }",
			"{\n\tfor a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
			"{\n\tfor a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
		},
		{ // 60
			"{ select a in b; do c\nd\ndone; }",
			"{\n\tselect a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
			"{\n\tselect a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}\n",
		},
		{ // 61
			"{ [[ a = b ]]; }",
			"{ [[ a == b ]]; }\n",
			"{\n\t[[ a == b ]];\n}\n",
		},
		{ // 62
			"{ [[ a = $(a\nb) ]]; }",
			"{\n\t[[\n\t\ta == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
			"{\n\t[[\n\t\ta == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
		},
		{ // 63
			"{ [[ a = b || c = $(a\nb) ]]; }",
			"{\n\t[[\n\t\ta == b || c == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
			"{\n\t[[\n\t\ta == b || c == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
		},
		{ // 64
			"{ [[ $(a\nb) ]]; }",
			"{\n\t[[\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
			"{\n\t[[\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}\n",
		},
		{ // 65
			"{ ((a)); }",
			"{ ((a)); }\n",
			"{\n\t(( a ));\n}\n",
		},
		{ // 66
			"{ (($(a\nb))); }",
			"{\n\t(($(\n\t\ta;\n\t\tb;\n\t)));\n}\n",
			"{\n\t(( $(\n\t\ta;\n\t\tb;\n\t) ));\n}\n",
		},
		{ // 67
			"{ (($(a\nb))) >&2; }",
			"{\n\t(($(\n\t\ta;\n\t\tb;\n\t))) >&2;\n}\n",
			"{\n\t(( $(\n\t\ta;\n\t\tb;\n\t) )) >&2;\n}\n",
		},
		{ // 68
			"{ ((a)) >$(a\nb); }",
			"{\n\t((a)) >$(\n\t\ta;\n\t\tb;\n\t);\n}\n",
			"{\n\t(( a )) > $(\n\t\ta;\n\t\tb;\n\t);\n}\n",
		},
		{ // 69
			"{ a=( # A\n); }",
			"{\n\ta=( # A\n\t);\n}\n",
			"{\n\ta=( # A\n\t);\n}\n",
		},
		{ // 70
			"{ a=( # A\n# B\n\n# C\n); }",
			"{\n\ta=( # A\n\t    # B\n\n\t# C\n\t);\n}\n",
			"{\n\ta=( # A\n\t    # B\n\n\t# C\n\t);\n}\n",
		},
		{ // 71
			"{ a=( # A\n# B\n\nb\n# C\n); }",
			"{\n\ta=( # A\n\t    # B\n\n\t\tb\n\t# C\n\t);\n}\n",
			"{\n\ta=( # A\n\t    # B\n\n\t\tb\n\t# C\n\t);\n}\n",
		},
		{ // 72
			"{ a=(b); }",
			"{ a=(b); }\n",
			"{\n\ta=( b );\n}\n",
		},
		{ // 73
			"{ a=($(a\nb)); }",
			"{\n\ta=(\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t);\n}\n",
			"{\n\ta=(\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t);\n}\n",
		},
		{ // 74
			"{ $(($(a\nb))); }",
			"{\n\t$(($(\n\t\ta;\n\t\tb;\n\t)));\n}\n",
			"{\n\t$(( $(\n\t\ta;\n\t\tb;\n\t) ));\n}\n",
		},
		{ // 75
			"{ ${a:=$(b)}; }",
			"{ ${a:=$(b)}; }\n",
			"{\n\t${a:=$(b)};\n}\n",
		},
		{ // 76
			"{ ${a:=$(b\nc)}; }",
			"{\n\t${a:=$(\n\t\tb;\n\t\tc;\n\t)};\n}\n",
			"{\n\t${a:=$(\n\t\tb;\n\t\tc;\n\t)};\n}\n",
		},
		{ // 77
			"{ ${a/b/$(c)}; }",
			"{ ${a/b/$(c)}; }\n",
			"{\n\t${a/b/$(c)};\n}\n",
		},
		{ // 78
			"{ ${a/b/$(c\nd)}; }",
			"{\n\t${a/b/$(\n\t\tc;\n\t\td;\n\t)};\n}\n",
			"{\n\t${a/b/$(\n\t\tc;\n\t\td;\n\t)};\n}\n",
		},
		{ // 79
			"{ {a,\"$(b\nc)\"}; }",
			"{\n\t{a,\"$(\n\t\tb;\n\t\tc;\n\t)\"};\n}\n",
			"{\n\t{a,\"$(\n\t\tb;\n\t\tc;\n\t)\"};\n}\n",
		},
		{ // 80
			"{ {a,b}; }",
			"{ {a,b}; }\n",
			"{\n\t{a,b};\n}\n",
		},
		{ // 81
			"a | b <<c\nc",
			"a | b <<c\nc\n",
			"a | b <<c\nc\n",
		},
		{ // 82
			"a && b <<c\nc",
			"a && b <<c\nc\n",
			"a && b <<c\nc\n",
		},
		{ // 83
			"a <<b\nb\n\nc;",
			"a <<b\nb\n\nc;\n",
			"a <<b\nb\n\nc;\n",
		},
		{ // 84
			"if a; then b;fi",
			"if a; then\n\tb;\nfi;\n",
			"if a; then\n\tb;\nfi;\n",
		},
		{ // 85
			"if a||b; then b\nc\nelif d\nthen\ne\nelse if f\nthen\ng\nfi\nfi",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;\n",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;\n",
		},
		{ // 86
			"if a; then b;elif c; then d;elif e\nthen f;fi",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;\n",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;\n",
		},
		{ // 87
			"if [ \"$a\" = \"b\" ]; then c;fi",
			"if [ \"$a\" = \"b\" ]; then\n\tc;\nfi;\n",
			"if [ \"$a\" = \"b\" ]; then\n\tc;\nfi;\n",
		},
		{ // 88
			"while a\ndo\nb\ndone",
			"while a; do\n\tb;\ndone;\n",
			"while a; do\n\tb;\ndone;\n",
		},
		{ // 89
			"until a&&b; # A\n# B\ndo\n# C\nb\nc\ndone",
			"until a && b; # A\n              # B\ndo\n\t# C\n\tb;\n\tc;\ndone;\n",
			"until a && b; # A\n              # B\ndo\n\t# C\n\tb;\n\tc;\ndone;\n",
		},
		{ // 90
			"a[b]=",
			"a[b]=;\n",
			"a[b]=;\n",
		},
		{ // 91
			"a[b+c]=",
			"a[b+c]=;\n",
			"a[b + c]=;\n",
		},
		{ // 92
			"a[$(b\nc)]=",
			"a[$(\n\tb;\n\tc;\n)]=;\n",
			"a[$(\n\tb;\n\tc;\n)]=;\n",
		},
		{ // 93
			"${a}",
			"${a};\n",
			"${a};\n",
		},
		{ // 94
			"${!a}",
			"${!a};\n",
			"${!a};\n",
		},
		{ // 95
			"${#a}",
			"${#a};\n",
			"${#a};\n",
		},
		{ // 96
			"${a:=b}",
			"${a:=b};\n",
			"${a:=b};\n",
		},
		{ // 97
			"${a=b}",
			"${a=b};\n",
			"${a=b};\n",
		},
		{ // 98
			"${a:?b}",
			"${a:?b};\n",
			"${a:?b};\n",
		},
		{ // 99
			"${a?b}",
			"${a?b};\n",
			"${a?b};\n",
		},
		{ // 100
			"${a:+b}",
			"${a:+b};\n",
			"${a:+b};\n",
		},
		{ // 101
			"${a+b}",
			"${a+b};\n",
			"${a+b};\n",
		},
		{ // 102
			"${a:-b}",
			"${a:-b};\n",
			"${a:-b};\n",
		},
		{ // 103
			"${a-b}",
			"${a-b};\n",
			"${a-b};\n",
		},
		{ // 104
			"${a#b}",
			"${a#b};\n",
			"${a#b};\n",
		},
		{ // 105
			"${a##b}",
			"${a##b};\n",
			"${a##b};\n",
		},
		{ // 106
			"${a%b}",
			"${a%b};\n",
			"${a%b};\n",
		},
		{ // 107
			"${a%%b}",
			"${a%%b};\n",
			"${a%%b};\n",
		},
		{ // 108
			"${a/b}",
			"${a/b};\n",
			"${a/b};\n",
		},
		{ // 109
			"${a/b/c}",
			"${a/b/c};\n",
			"${a/b/c};\n",
		},
		{ // 110
			"${a//b}",
			"${a//b};\n",
			"${a//b};\n",
		},
		{ // 111
			"${a//b/c}",
			"${a//b/c};\n",
			"${a//b/c};\n",
		},
		{ // 112
			"${a/%b}",
			"${a/%b};\n",
			"${a/%b};\n",
		},
		{ // 113
			"${a/%b/c}",
			"${a/%b/c};\n",
			"${a/%b/c};\n",
		},
		{ // 114
			"${a/#b}",
			"${a/#b};\n",
			"${a/#b};\n",
		},
		{ // 115
			"${a/#b/c}",
			"${a/#b/c};\n",
			"${a/#b/c};\n",
		},
		{ // 116
			"${a^b}",
			"${a^b};\n",
			"${a^b};\n",
		},
		{ // 117
			"${a^^b}",
			"${a^^b};\n",
			"${a^^b};\n",
		},
		{ // 118
			"${a,b}",
			"${a,b};\n",
			"${a,b};\n",
		},
		{ // 119
			"${a,,b}",
			"${a,,b};\n",
			"${a,,b};\n",
		},
		{ // 120
			"${*}",
			"${*};\n",
			"${*};\n",
		},
		{ // 121
			"${@}",
			"${@};\n",
			"${@};\n",
		},
		{ // 122
			"${a:1}",
			"${a:1};\n",
			"${a:1};\n",
		},
		{ // 123
			"${a:1:2}",
			"${a:1:2};\n",
			"${a:1:2};\n",
		},
		{ // 124
			"${a: -1:2}",
			"${a: -1:2};\n",
			"${a: -1:2};\n",
		},
		{ // 125
			"${a@U}",
			"${a@U};\n",
			"${a@U};\n",
		},
		{ // 126
			"${a@u}",
			"${a@u};\n",
			"${a@u};\n",
		},
		{ // 127
			"${a@L}",
			"${a@L};\n",
			"${a@L};\n",
		},
		{ // 128
			"${a@Q}",
			"${a@Q};\n",
			"${a@Q};\n",
		},
		{ // 129
			"${a@E}",
			"${a@E};\n",
			"${a@E};\n",
		},
		{ // 130
			"${a@P}",
			"${a@P};\n",
			"${a@P};\n",
		},
		{ // 131
			"${a@A}",
			"${a@A};\n",
			"${a@A};\n",
		},
		{ // 132
			"${a@K}",
			"${a@K};\n",
			"${a@K};\n",
		},
		{ // 133
			"${a@a}",
			"${a@a};\n",
			"${a@a};\n",
		},
		{ // 134
			"${a@k}",
			"${a@k};\n",
			"${a@k};\n",
		},
		{ // 135
			"${!a@}",
			"${!a@};\n",
			"${!a@};\n",
		},
		{ // 136
			"${!a*}",
			"${!a*};\n",
			"${!a*};\n",
		},
		{ // 137
			"[[ a == b ]]",
			"[[ a == b ]];\n",
			"[[ a == b ]];\n",
		},
		{ // 138
			"[[ a == b$c ]]",
			"[[ a == b$c ]];\n",
			"[[ a == b$c ]];\n",
		},
		{ // 139
			"[[ a == b\"c\" ]]",
			"[[ a == b\"c\" ]];\n",
			"[[ a == b\"c\" ]];\n",
		},
		{ // 140
			"case a in a|b) a;\nb\nesac",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;\n",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;\n",
		},
		{ // 141
			"case a in a|b|\"c\") a;\nb\nesac",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;\n",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;\n",
		},
		{ // 142
			"! a",
			"! a;\n",
			"! a;\n",
		},
		{ // 143
			"coproc a",
			"coproc a;\n",
			"coproc a;\n",
		},
		{ // 144
			"! coproc a",
			"! coproc a;\n",
			"! coproc a;\n",
		},
		{ // 145
			"time a",
			"time a;\n",
			"time a;\n",
		},
		{ // 146
			"time -p a",
			"time -p a;\n",
			"time -p a;\n",
		},
		{ // 147
			"time coproc a",
			"time coproc a;\n",
			"time coproc a;\n",
		},
		{ // 148
			"time -p coproc a",
			"time -p coproc a;\n",
			"time -p coproc a;\n",
		},
		{ // 149
			"time ! coproc a",
			"time ! coproc a;\n",
			"time ! coproc a;\n",
		},
		{ // 150
			"time -p ! coproc a",
			"time -p ! coproc a;\n",
			"time -p ! coproc a;\n",
		},
		{ // 151
			"coproc a if b; then c\nfi",
			"coproc a if b; then\n\tc;\nfi;\n",
			"coproc a if b; then\n\tc;\nfi;\n",
		},
		{ // 152
			"select a; do b; done",
			"select a; do\n\tb;\ndone;\n",
			"select a; do\n\tb;\ndone;\n",
		},
		{ // 153
			"select a in b c; do b; done",
			"select a in b c; do\n\tb;\ndone;\n",
			"select a in b c; do\n\tb;\ndone;\n",
		},
		{ // 154
			"a&",
			"a&\n",
			"a &\n",
		},
		{ // 155
			"[[ # A\na == b\n# B\n]]",
			"[[ # A\n\ta == b\n# B\n]];\n",
			"[[ # A\n\ta == b\n# B\n]];\n",
		},
		{ // 156
			"[[ # A\na == b\n]]",
			"[[ # A\n\ta == b\n]];\n",
			"[[ # A\n\ta == b\n]];\n",
		},
		{ // 157
			"[[\n\t! # A\n\ta -ot b ]]",
			"[[\n\t! # A\n\ta -ot b\n]];\n",
			"[[\n\t! # A\n\ta -ot b\n]];\n",
		},
		{ // 158
			"[[\n\t! # A\n\ta -nt b ]]",
			"[[\n\t! # A\n\ta -nt b\n]];\n",
			"[[\n\t! # A\n\ta -nt b\n]];\n",
		},
		{ // 159
			"[[ (a -ef b) ]]",
			"[[ ( a -ef b ) ]];\n",
			"[[ ( a -ef b ) ]];\n",
		},
		{ // 160
			"[[ (# A\na -ge b) ]]",
			"[[\n\t( # A\n\t\ta -ge b\n\t)\n]];\n",
			"[[\n\t( # A\n\t\ta -ge b\n\t)\n]];\n",
		},
		{ // 161
			"[[ (a -gt b\n# A\n) ]]",
			"[[\n\t(\n\t\ta -gt b\n\t# A\n\t)\n]];\n",
			"[[\n\t(\n\t\ta -gt b\n\t# A\n\t)\n]];\n",
		},
		{ // 162
			"[[ (\n# A\na -le b # B\n) ]]",
			"[[\n\t(\n\t\t# A\n\t\ta -le b # B\n\t)\n]];\n",
			"[[\n\t(\n\t\t# A\n\t\ta -le b # B\n\t)\n]];\n",
		},
		{ // 163
			"[[\n# A\na =~ b #B\n]]",
			"[[\n\t# A\n\ta =~ b #B\n]];\n",
			"[[\n\t# A\n\ta =~ b #B\n]];\n",
		},
		{ // 164
			"[[ # A\n\n# B\na != b # C\n\n# D\n]]",
			"[[ # A\n\n\t# B\n\ta != b # C\n\n# D\n]];\n",
			"[[ # A\n\n\t# B\n\ta != b # C\n\n# D\n]];\n",
		},
		{ // 165
			"[[ # A\n\n# B\na < b # C\n||# D\n\n# E\nd>e # F\n\n# G\n]]",
			"[[ # A\n\n\t# B\n\ta < b # C\n\t|| # D\n\n\t   # E\n\td > e # F\n\n# G\n]];\n",
			"[[ # A\n\n\t# B\n\ta < b # C\n\t|| # D\n\n\t   # E\n\td > e # F\n\n# G\n]];\n",
		},
		{ // 166
			"[[ # A\n# B\n\n# C\n# D\n( # E\n# F\n\n# G\n# H\na -eq b # I\n# J\n&& # K\n# L\nc -ne d # M\n# N\n\n# O\n\n# P\n) # Q\n# R\n\n# S\n# T\n&& # U\n# V\n\n# W\n# X\ne -lt f # Y\n# Z\n]]",
			"[[ # A\n   # B\n\n\t# C\n\t# D\n\t( # E\n\t  # F\n\n\t\t# G\n\t\t# H\n\t\ta -eq b # I\n\t\t        # J\n\t\t&& # K\n\t\t   # L\n\t\tc -ne d # M\n\t\t        # N\n\n\t# O\n\n\t# P\n\t) # Q\n\t  # R\n\n\t  # S\n\t  # T\n\t&& # U\n\t   # V\n\n\t   # W\n\t   # X\n\te -lt f # Y\n\t        # Z\n]];\n",
			"[[ # A\n   # B\n\n\t# C\n\t# D\n\t( # E\n\t  # F\n\n\t\t# G\n\t\t# H\n\t\ta -eq b # I\n\t\t        # J\n\t\t&& # K\n\t\t   # L\n\t\tc -ne d # M\n\t\t        # N\n\n\t# O\n\n\t# P\n\t) # Q\n\t  # R\n\n\t  # S\n\t  # T\n\t&& # U\n\t   # V\n\n\t   # W\n\t   # X\n\te -lt f # Y\n\t        # Z\n]];\n",
		},
		{ // 167
			"[[ -e a ]]",
			"[[ -e a ]];\n",
			"[[ -e a ]];\n",
		},
		{ // 168
			"[[ -b a || -c b ]]",
			"[[ -b a || -c b ]];\n",
			"[[ -b a || -c b ]];\n",
		},
		{ // 169
			"[[ -d a && -f b ]]",
			"[[ -d a && -f b ]];\n",
			"[[ -d a && -f b ]];\n",
		},
		{ // 170
			"[[ (-g a && -L b) || -k c ]]",
			"[[ ( -g a && -L b ) || -k c ]];\n",
			"[[ ( -g a && -L b ) || -k c ]];\n",
		},
		{ // 171
			"[[ (-p a && -r b) || (-s c && (-t d || -u e)) ]]",
			"[[ ( -p a && -r b ) || ( -s c && ( -t d || -u e ) ) ]];\n",
			"[[ ( -p a && -r b ) || ( -s c && ( -t d || -u e ) ) ]];\n",
		},
		{ // 172
			"[[ (-w a && -x b) || (-G c && (-N d || -O e)) ]]",
			"[[ ( -w a && -x b ) || ( -G c && ( -N d || -O e ) ) ]];\n",
			"[[ ( -w a && -x b ) || ( -G c && ( -N d || -O e ) ) ]];\n",
		},
		{ // 173
			"[[ (-S a && -o b) || (-v c && (-R d || -z e || -n f)) ]]",
			"[[ ( -S a && -o b ) || ( -v c && ( -R d || -z e || -n f ) ) ]];\n",
			"[[ ( -S a && -o b ) || ( -v c && ( -R d || -z e || -n f ) ) ]];\n",
		},
		{ // 174
			"if a; then b; fi",
			"if a; then\n\tb;\nfi;\n",
			"if a; then\n\tb;\nfi;\n",
		},
		{ // 175
			"if a # A\nthen b; fi",
			"if a; # A\nthen\n\tb;\nfi;\n",
			"if a; # A\nthen\n\tb;\nfi;\n",
		},
		{ // 176
			"if a # A\n# B\nthen b; fi",
			"if a; # A\n      # B\nthen\n\tb;\nfi;\n",
			"if a; # A\n      # B\nthen\n\tb;\nfi;\n",
		},
		{ // 177
			"${a/b/ }",
			"${a/b/ };\n",
			"${a/b/ };\n",
		},
		{ // 178
			"<(a)",
			"<(a);\n",
			"<(a);\n",
		},
		{ // 179
			">(a;b)",
			">(a; b;);\n",
			">(\n\ta;\n\tb;\n);\n",
		},
		{ // 180
			"`a`",
			"`a`;\n",
			"`a`;\n",
		},
		{ // 181
			"`a``b`",
			"`a``b`;\n",
			"`a``b`;\n",
		},
		{ // 182
			"`a \\`b\\``",
			"`a \\`b\\``;\n",
			"`a \\`b\\``;\n",
		},
		{ // 183
			"while read a; do\n\tb \"$a\";\ndone <<EOF\nA\nB\nC\nEOF",
			"while read a; do\n\tb \"$a\";\ndone <<EOF\nA\nB\nC\nEOF\n",
			"while read a; do\n\tb \"$a\";\ndone <<EOF\nA\nB\nC\nEOF\n",
		},
		{ // 184
			"while read a; do\n\tb \"$a\";\ndone < <(c)",
			"while read a; do\n\tb \"$a\";\ndone < <(c);\n",
			"while read a; do\n\tb \"$a\";\ndone < <(c);\n",
		},
		{ // 185
			"{a,b}",
			"{a,b};\n",
			"{a,b};\n",
		},
		{ // 186
			"{a,,bcd,\"\"}",
			"{a,,bcd,\"\"};\n",
			"{a,,bcd,\"\"};\n",
		},
		{ // 187
			"{1..2}",
			"{1..2};\n",
			"{1..2};\n",
		},
		{ // 188
			"{-10..-100..-2}",
			"{-10..-100..-2};\n",
			"{-10..-100..-2};\n",
		},
		{ // 189
			"{a..z}",
			"{a..z};\n",
			"{a..z};\n",
		},
		{ // 190
			"{A..Z..3}",
			"{A..Z..3};\n",
			"{A..Z..3};\n",
		},
	} {
		for m, input := range test {
			if m == 2 && (n == 42 || n == 35 || n == 178) {
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
