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
			"((a));",
			"(( a ));",
		},
		{ // 2
			"$(( a + b ))",
			"$((a+b));",
			"$(( a + b ));",
		},
		{ // 3
			"a=(1)",
			"a=(1);",
			"a=( 1 );",
		},
		{ // 4
			"a=(\n# word comment\nb # post-word comment\n)",
			"a=(\n\t# word comment\n\tb # post-word comment\n);",
			"a=(\n\t# word comment\n\tb # post-word comment\n);",
		},
		{ // 5
			"a=1",
			"a=1;",
			"a=1;",
		},
		{ // 6
			"a+=1",
			"a+=1;",
			"a+=1;",
		},
		{ // 7
			"let a=1+2",
			"let a=1+2;",
			"let a=1+2;",
		},
		{ // 8
			"let a=1+(2+3)",
			"let a=1+(2+3);",
			"let a=1+( 2 + 3 );",
		},
		{ // 9
			"let a=b?(c?d:e):f",
			"let a=b?(c?d:e):f;",
			"let a=b?( c ? d : e ):f;",
		},
		{ // 10
			"case a in b)c\nesac",
			"case a in\nb)\n\tc;;\nesac;",
			"case a in\nb)\n\tc;;\nesac;",
		},
		{ // 11
			"case a in b);;& c) d;&\nesac",
			"case a in\nb)\n\t;;&\nc)\n\td;&\nesac;",
			"case a in\nb)\n\t;;&\nc)\n\td;&\nesac;",
		},
		{ // 12
			"case a #A\nin #B\n   #C\nb)c;;\n#D\nesac",
			"case a #A\nin #B\n   #C\nb)\n\tc;;\n#D\nesac;",
			"case a #A\nin #B\n   #C\nb)\n\tc;;\n#D\nesac;",
		},
		{ // 13
			"a=1 b=2 c d >e <f",
			"a=1 b=2 c d >e <f;",
			"a=1 b=2 c d >e <f;",
		},
		{ // 14
			"a b >c <d",
			"a b >c <d;",
			"a b >c <d;",
		},
		{ // 15
			">a <b",
			">a <b;",
			">a <b;",
		},
		{ // 16
			"a=1 b=2 >c <d",
			"a=1 b=2 >c <d;",
			"a=1 b=2 >c <d;",
		},
		{ // 17
			"a <<b\nheredoc\ncontents\nb",
			"a <<b;\nheredoc\ncontents\nb",
			"a <<b;\nheredoc\ncontents\nb",
		},
		{ // 18
			"a <<-b\n\theredoc\n\tcontents\nb",
			"a <<-b;\nheredoc\ncontents\nb",
			"a <<-b;\nheredoc\ncontents\nb",
		},
		{ // 19
			"$(a)",
			"$(a);",
			"$(a);",
		},
		{ // 20
			"$(a\nb)",
			"$(\n\ta;\n\tb;\n);",
			"$(\n\ta;\n\tb;\n);",
		},
		{ // 21
			"# A\na # B\n  # C\n\n# D",
			"# A\na; # B\n   # C\n\n# D",
			"# A\na; # B\n   # C\n\n# D",
		},
		{ // 22
			"case a in\nesac <a 2>&1;",
			"case a in\nesac <a 2>&1;",
			"case a in\nesac <a 2>&1;",
		},
		{ // 23
			"a\nb\n\nc\nd\n\n\n\n\ne",
			"a;\nb;\n\nc;\nd;\n\ne;",
			"a;\nb;\n\nc;\nd;\n\ne;",
		},
		{ // 24
			"a\nfor b; do\nc\ndone",
			"a;\nfor b; do\n\tc;\ndone;",
			"a;\nfor b; do\n\tc;\ndone;",
		},
		{ // 25
			"for a;do c\nd\ndone",
			"for a; do\n\tc;\n\td;\ndone;",
			"for a; do\n\tc;\n\td;\ndone;",
		},
		{ // 26
			"for a in b\ndo c\ndone",
			"for a in b; do\n\tc;\ndone;",
			"for a in b; do\n\tc;\ndone;",
		},
		{ // 27
			"for a in b c\ndo d\ndone",
			"for a in b c; do\n\td;\ndone;",
			"for a in b c; do\n\td;\ndone;",
		},
		{ // 28
			"for ((a=0;a<1;a++));do b\ndone",
			"for ((a=0;a<1;a++)); do\n\tb;\ndone;",
			"for (( a = 0; a < 1; a ++ )); do\n\tb;\ndone;",
		},
		{ // 29
			"function a() { b; }",
			"function a() { b; }",
			"function a() {\n\tb;\n}",
		},
		{ // 30
			"function a { b; }",
			"function a() { b; }",
			"function a() {\n\tb;\n}",
		},
		{ // 31
			"a() { b; }",
			"a() { b; }",
			"a() {\n\tb;\n}",
		},
		{ // 32
			"function a() # A\n# B\n{ b; }",
			"function a() # A\n             # B\n{ b; }",
			"function a() # A\n             # B\n{\n\tb;\n}",
		},
		{ // 33
			"a() # A\n# B\n{ b; }",
			"a() # A\n    # B\n{ b; }",
			"a() # A\n    # B\n{\n\tb;\n}",
		},
		{ // 34
			"{ a; }",
			"{ a; }",
			"{\n\ta;\n}",
		},
		{ // 35
			"( a; )",
			"( a; )",
			"(\n\ta;\n)",
		},
		{ // 36
			"{ a; b; }",
			"{ a; b; }",
			"{\n\ta;\n\tb;\n}",
		},
		{ // 37
			"{ a || b; }",
			"{ a || b; }",
			"{\n\ta || b;\n}",
		},
		{ // 38
			"( a;\nb; )",
			"(\n\ta;\n\tb;\n)",
			"(\n\ta;\n\tb;\n)",
		},
		{ // 39
			"{ # A\na; # B\n}",
			"{ # A\n\ta; # B\n}",
			"{ # A\n\ta; # B\n}",
		},
		{ // 40
			"{ a; # A\n}",
			"{\n\ta; # A\n}",
			"{\n\ta; # A\n}",
		},
		{ // 41
			"{ a | b $(c\nd)\n}",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n}",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n}",
		},
		{ // 42
			"{ declare a=$(b); }",
			"{ declare a=$(b); }",
			"{\n\tdeclare a=$(b);\n}",
		},
		{ // 43
			"{ declare a=$(b;c); }",
			"{ declare a=$(b; c;); }",
			"{\n\tdeclare a=$(\n\t\tb;\n\t\tc;\n\t);\n}",
		},
		{ // 44
			"{ a=1 b; }",
			"{ a=1 b; }",
			"{\n\ta=1 b;\n}",
		},
		{ // 45
			"{ a=$(b\nc) d; }",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n}",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n}",
		},
		{ // 46
			"{ let a[$(b)]=c; }",
			"{ let a[$(b)]=c; }",
			"{\n\tlet a[$(b)]=c;\n}",
		},
		{ // 47
			"{ let a[$(b\nc)]=d; }",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n}",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n}",
		},
		{ // 48
			"{ let a=$(b); }",
			"{ let a=$(b); }",
			"{\n\tlet a=$(b);\n}",
		},
		{ // 49
			"{ let a=$(b\nc); }",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n}",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n}",
		},
		{ // 50
			"{ ${a[$(b\nc)]}; }",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n}",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n}",
		},
		{ // 51
			"<<a\nb$c\na",
			"<<a;\nb$c\na",
			"<<a;\nb$c\na",
		},
		{ // 52
			"{\n<<a\nb$c\na\n}",
			"{\n\t<<a;\nb$c\na\n}",
			"{\n\t<<a;\nb$c\na\n}",
		},
		{ // 53
			"{\n<<-a\nb$c\na\n}",
			"{\n\t<<-a;\n\tb$c\n\ta\n}",
			"{\n\t<<-a;\n\tb$c\n\ta\n}",
		},
		{ // 54
			"{ function a() { # A\nb; } }",
			"{\n\tfunction a() { # A\n\t\tb;\n\t}\n}",
			"{\n\tfunction a() { # A\n\t\tb;\n\t}\n}",
		},
		{ // 55
			"{ function a() { b;\nc; } }",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t}\n}",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t}\n}",
		},
		{ // 56
			"{ if a; then b\nc\nfi; }",
			"{\n\tif a; then\n\t\tb;\n\t\tc;\n\tfi;\n}",
			"{\n\tif a; then\n\t\tb;\n\t\tc;\n\tfi;\n}",
		},
		{ // 57
			"{ until a; do b\nc\ndone; }",
			"{\n\tuntil a; do\n\t\tb;\n\t\tc;\n\tdone;\n}",
			"{\n\tuntil a; do\n\t\tb;\n\t\tc;\n\tdone;\n}",
		},
		{ // 58
			"{ case a in b)c;;esac; }",
			"{\n\tcase a in\n\tb)\n\t\tc;;\n\tesac;\n}",
			"{\n\tcase a in\n\tb)\n\t\tc;;\n\tesac;\n}",
		},
		{ // 59
			"{ for a in b; do c\nd\ndone; }",
			"{\n\tfor a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}",
			"{\n\tfor a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}",
		},
		{ // 60
			"{ select a in b; do c\nd\ndone; }",
			"{\n\tselect a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}",
			"{\n\tselect a in b; do\n\t\tc;\n\t\td;\n\tdone;\n}",
		},
		{ // 61
			"{ [[ a = b ]]; }",
			"{ [[ a == b ]]; }",
			"{\n\t[[ a == b ]];\n}",
		},
		{ // 62
			"{ [[ a = $(a\nb) ]]; }",
			"{\n\t[[\n\t\ta == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}",
			"{\n\t[[\n\t\ta == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}",
		},
		{ // 63
			"{ [[ a = b || c = $(a\nb) ]]; }",
			"{\n\t[[\n\t\ta == b || c == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}",
			"{\n\t[[\n\t\ta == b || c == $(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}",
		},
		{ // 64
			"{ [[ $(a\nb) ]]; }",
			"{\n\t[[\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}",
			"{\n\t[[\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t]];\n}",
		},
		{ // 65
			"{ ((a)); }",
			"{ ((a)); }",
			"{\n\t(( a ));\n}",
		},
		{ // 66
			"{ (($(a\nb))); }",
			"{\n\t(($(\n\t\ta;\n\t\tb;\n\t)));\n}",
			"{\n\t(( $(\n\t\ta;\n\t\tb;\n\t) ));\n}",
		},
		{ // 67
			"{ (($(a\nb))) >&2; }",
			"{\n\t(($(\n\t\ta;\n\t\tb;\n\t))) >&2;\n}",
			"{\n\t(( $(\n\t\ta;\n\t\tb;\n\t) )) >&2;\n}",
		},
		{ // 68
			"{ ((a)) >$(a\nb); }",
			"{\n\t((a)) >$(\n\t\ta;\n\t\tb;\n\t);\n}",
			"{\n\t(( a )) >$(\n\t\ta;\n\t\tb;\n\t);\n}",
		},
		{ // 69
			"{ a=( # A\n); }",
			"{\n\ta=( # A\n\t);\n}",
			"{\n\ta=( # A\n\t);\n}",
		},
		{ // 70
			"{ a=( # A\n# B\n\n# C\n); }",
			"{\n\ta=( # A\n\t    # B\n\n\t# C\n\t);\n}",
			"{\n\ta=( # A\n\t    # B\n\n\t# C\n\t);\n}",
		},
		{ // 71
			"{ a=( # A\n# B\n\nb\n# C\n); }",
			"{\n\ta=( # A\n\t    # B\n\n\t\tb\n\t# C\n\t);\n}",
			"{\n\ta=( # A\n\t    # B\n\n\t\tb\n\t# C\n\t);\n}",
		},
		{ // 72
			"{ a=(b); }",
			"{ a=(b); }",
			"{\n\ta=( b );\n}",
		},
		{ // 73
			"{ a=($(a\nb)); }",
			"{\n\ta=(\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t);\n}",
			"{\n\ta=(\n\t\t$(\n\t\t\ta;\n\t\t\tb;\n\t\t)\n\t);\n}",
		},
		{ // 74
			"{ $(($(a\nb))); }",
			"{\n\t$(($(\n\t\ta;\n\t\tb;\n\t)));\n}",
			"{\n\t$(( $(\n\t\ta;\n\t\tb;\n\t) ));\n}",
		},
		{ // 75
			"{ ${a:=$(b)}; }",
			"{ ${a:=$(b)}; }",
			"{\n\t${a:=$(b)};\n}",
		},
		{ // 76
			"{ ${a:=$(b\nc)}; }",
			"{\n\t${a:=$(\n\t\tb;\n\t\tc;\n\t)};\n}",
			"{\n\t${a:=$(\n\t\tb;\n\t\tc;\n\t)};\n}",
		},
		{ // 77
			"{ ${a/b/$(c)}; }",
			"{ ${a/b/$(c)}; }",
			"{\n\t${a/b/$(c)};\n}",
		},
		{ // 78
			"{ ${a/b/$(c\nd)}; }",
			"{\n\t${a/b/$(\n\t\tc;\n\t\td;\n\t)};\n}",
			"{\n\t${a/b/$(\n\t\tc;\n\t\td;\n\t)};\n}",
		},
		{ // 79
			"a | b <<c\nc",
			"a | b <<c;\nc",
			"a | b <<c;\nc",
		},
		{ // 80
			"a && b <<c\nc",
			"a && b <<c;\nc",
			"a && b <<c;\nc",
		},
		{ // 81
			"if a; then b;fi",
			"if a; then\n\tb;\nfi;",
			"if a; then\n\tb;\nfi;",
		},
		{ // 82
			"if a||b; then b\nc\nelif d\nthen\ne\nelse if f\nthen\ng\nfi\nfi",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
		},
		{ // 83
			"if a; then b;elif c; then d;elif e\nthen f;fi",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
		},
		{ // 84
			"while a\ndo\nb\ndone",
			"while a; do\n\tb;\ndone;",
			"while a; do\n\tb;\ndone;",
		},
		{ // 85
			"until a&&b; # A\n# B\ndo\n# C\nb\nc\ndone",
			"until a && b; # A\n              # B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
			"until a && b; # A\n              # B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
		},
		{ // 86
			"a[b]=",
			"a[b]=;",
			"a[b]=;",
		},
		{ // 87
			"a[b+c]=",
			"a[b+c]=;",
			"a[b + c]=;",
		},
		{ // 88
			"a[$(b\nc)]=",
			"a[$(\n\tb;\n\tc;\n)]=;",
			"a[$(\n\tb;\n\tc;\n)]=;",
		},
		{ // 89
			"${a}",
			"${a};",
			"${a};",
		},
		{ // 90
			"${!a}",
			"${!a};",
			"${!a};",
		},
		{ // 91
			"${#a}",
			"${#a};",
			"${#a};",
		},
		{ // 92
			"${a:=b}",
			"${a:=b};",
			"${a:=b};",
		},
		{ // 93
			"${a=b}",
			"${a=b};",
			"${a=b};",
		},
		{ // 94
			"${a:?b}",
			"${a:?b};",
			"${a:?b};",
		},
		{ // 95
			"${a?b}",
			"${a?b};",
			"${a?b};",
		},
		{ // 96
			"${a:+b}",
			"${a:+b};",
			"${a:+b};",
		},
		{ // 97
			"${a+b}",
			"${a+b};",
			"${a+b};",
		},
		{ // 98
			"${a:-b}",
			"${a:-b};",
			"${a:-b};",
		},
		{ // 99
			"${a-b}",
			"${a-b};",
			"${a-b};",
		},
		{ // 100
			"${a#b}",
			"${a#b};",
			"${a#b};",
		},
		{ // 101
			"${a##b}",
			"${a##b};",
			"${a##b};",
		},
		{ // 102
			"${a%b}",
			"${a%b};",
			"${a%b};",
		},
		{ // 103
			"${a%%b}",
			"${a%%b};",
			"${a%%b};",
		},
		{ // 104
			"${a/b}",
			"${a/b};",
			"${a/b};",
		},
		{ // 105
			"${a/b/c}",
			"${a/b/c};",
			"${a/b/c};",
		},
		{ // 106
			"${a//b}",
			"${a//b};",
			"${a//b};",
		},
		{ // 107
			"${a//b/c}",
			"${a//b/c};",
			"${a//b/c};",
		},
		{ // 108
			"${a/%b}",
			"${a/%b};",
			"${a/%b};",
		},
		{ // 109
			"${a/%b/c}",
			"${a/%b/c};",
			"${a/%b/c};",
		},
		{ // 110
			"${a/#b}",
			"${a/#b};",
			"${a/#b};",
		},
		{ // 111
			"${a/#b/c}",
			"${a/#b/c};",
			"${a/#b/c};",
		},
		{ // 112
			"${a^b}",
			"${a^b};",
			"${a^b};",
		},
		{ // 113
			"${a^^b}",
			"${a^^b};",
			"${a^^b};",
		},
		{ // 114
			"${a,b}",
			"${a,b};",
			"${a,b};",
		},
		{ // 115
			"${a,,b}",
			"${a,,b};",
			"${a,,b};",
		},
		{ // 116
			"${*}",
			"${*};",
			"${*};",
		},
		{ // 117
			"${@}",
			"${@};",
			"${@};",
		},
		{ // 118
			"${a:1}",
			"${a:1};",
			"${a:1};",
		},
		{ // 119
			"${a:1:2}",
			"${a:1:2};",
			"${a:1:2};",
		},
		{ // 120
			"${a: -1:2}",
			"${a: -1:2};",
			"${a: -1:2};",
		},
		{ // 121
			"${a@U}",
			"${a@U};",
			"${a@U};",
		},
		{ // 122
			"${a@u}",
			"${a@u};",
			"${a@u};",
		},
		{ // 123
			"${a@L}",
			"${a@L};",
			"${a@L};",
		},
		{ // 124
			"${a@Q}",
			"${a@Q};",
			"${a@Q};",
		},
		{ // 125
			"${a@E}",
			"${a@E};",
			"${a@E};",
		},
		{ // 126
			"${a@P}",
			"${a@P};",
			"${a@P};",
		},
		{ // 127
			"${a@A}",
			"${a@A};",
			"${a@A};",
		},
		{ // 128
			"${a@K}",
			"${a@K};",
			"${a@K};",
		},
		{ // 129
			"${a@a}",
			"${a@a};",
			"${a@a};",
		},
		{ // 130
			"${a@k}",
			"${a@k};",
			"${a@k};",
		},
		{ // 131
			"${!a@}",
			"${!a@};",
			"${!a@};",
		},
		{ // 132
			"${!a*}",
			"${!a*};",
			"${!a*};",
		},
		{ // 133
			"[[ a == b ]]",
			"[[ a == b ]];",
			"[[ a == b ]];",
		},
		{ // 134
			"[[ a == b$c ]]",
			"[[ a == b$c ]];",
			"[[ a == b$c ]];",
		},
		{ // 135
			"[[ a == b\"c\" ]]",
			"[[ a == b\"c\" ]];",
			"[[ a == b\"c\" ]];",
		},
		{ // 136
			"case a in a|b) a;\nb\nesac",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;",
		},
		{ // 137
			"case a in a|b|\"c\") a;\nb\nesac",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;",
		},
		{ // 138
			"! a",
			"! a;",
			"! a;",
		},
		{ // 139
			"coproc a",
			"coproc a;",
			"coproc a;",
		},
		{ // 140
			"coproc a if b; then c\nfi",
			"coproc a if b; then\n\tc;\nfi;",
			"coproc a if b; then\n\tc;\nfi;",
		},
		{ // 141
			"select a; do b; done",
			"select a; do\n\tb;\ndone;",
			"select a; do\n\tb;\ndone;",
		},
		{ // 142
			"select a in b c; do b; done",
			"select a in b c; do\n\tb;\ndone;",
			"select a in b c; do\n\tb;\ndone;",
		},
		{ // 143
			"a&",
			"a&",
			"a &",
		},
		{ // 144
			"[[ # A\na == b\n# B\n]]",
			"[[ # A\n\ta == b\n# B\n]];",
			"[[ # A\n\ta == b\n# B\n]];",
		},
		{ // 145
			"[[ # A\na == b\n]]",
			"[[ # A\n\ta == b\n]];",
			"[[ # A\n\ta == b\n]];",
		},
		{ // 146
			"[[\n\t! # A\n\ta == b ]]",
			"[[\n\t! # A\n\ta == b\n]];",
			"[[\n\t! # A\n\ta == b\n]];",
		},
		{ // 147
			"[[\n\t! # A\n\ta == b ]]",
			"[[\n\t! # A\n\ta == b\n]];",
			"[[\n\t! # A\n\ta == b\n]];",
		},
		{ // 148
			"[[ (a == b) ]]",
			"[[ ( a == b ) ]];",
			"[[ ( a == b ) ]];",
		},
		{ // 149
			"[[ (# A\na -ge b) ]]",
			"[[\n\t( # A\n\t\ta -ge b\n\t)\n]];",
			"[[\n\t( # A\n\t\ta -ge b\n\t)\n]];",
		},
		{ // 150
			"[[ (a -gt b\n# A\n) ]]",
			"[[\n\t(\n\t\ta -gt b\n\t# A\n\t)\n]];",
			"[[\n\t(\n\t\ta -gt b\n\t# A\n\t)\n]];",
		},
		{ // 151
			"[[ (\n# A\na -le b # B\n) ]]",
			"[[\n\t(\n\t\t# A\n\t\ta -le b # B\n\t)\n]];",
			"[[\n\t(\n\t\t# A\n\t\ta -le b # B\n\t)\n]];",
		},
		{ // 152
			"[[\n# A\na =~ b #B\n]]",
			"[[\n\t# A\n\ta =~ b #B\n]];",
			"[[\n\t# A\n\ta =~ b #B\n]];",
		},
		{ // 153
			"[[ # A\n\n# B\na != b # C\n\n# D\n]]",
			"[[ # A\n\n\t# B\n\ta != b # C\n\n# D\n]];",
			"[[ # A\n\n\t# B\n\ta != b # C\n\n# D\n]];",
		},
		{ // 154
			"[[ # A\n\n# B\na < b # C\n||# D\n\n# E\nd>e # F\n\n# G\n]]",
			"[[ # A\n\n\t# B\n\ta < b # C\n\t|| # D\n\n\t   # E\n\td > e # F\n\n# G\n]];",
			"[[ # A\n\n\t# B\n\ta < b # C\n\t|| # D\n\n\t   # E\n\td > e # F\n\n# G\n]];",
		},
		{ // 155
			"[[ # A\n# B\n\n# C\n# D\n( # E\n# F\n\n# G\n# H\na -eq b # I\n# J\n&& # K\n# L\nc -ne d # M\n# N\n\n# O\n\n# P\n) # Q\n# R\n\n# S\n# T\n&& # U\n# V\n\n# W\n# X\ne -lt f # Y\n# Z\n]]",
			"[[ # A\n   # B\n\n\t# C\n\t# D\n\t( # E\n\t  # F\n\n\t\t# G\n\t\t# H\n\t\ta -eq b # I\n\t\t        # J\n\t\t&& # K\n\t\t   # L\n\t\tc -ne d # M\n\t\t        # N\n\n\t# O\n\n\t# P\n\t) # Q\n\t  # R\n\n\t  # S\n\t  # T\n\t&& # U\n\t   # V\n\n\t   # W\n\t   # X\n\te -lt f # Y\n\t        # Z\n]];",
			"[[ # A\n   # B\n\n\t# C\n\t# D\n\t( # E\n\t  # F\n\n\t\t# G\n\t\t# H\n\t\ta -eq b # I\n\t\t        # J\n\t\t&& # K\n\t\t   # L\n\t\tc -ne d # M\n\t\t        # N\n\n\t# O\n\n\t# P\n\t) # Q\n\t  # R\n\n\t  # S\n\t  # T\n\t&& # U\n\t   # V\n\n\t   # W\n\t   # X\n\te -lt f # Y\n\t        # Z\n]];",
		},
	} {
		for m, input := range test {
			if m == 2 && (n == 42 || n == 35) {
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
