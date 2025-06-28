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
			"let a=1+2",
			"let a=1+2;",
			"let a=1+2;",
		},
		{ // 7
			"let a=1+(2+3)",
			"let a=1+(2+3);",
			"let a=1+( 2 + 3 );",
		},
		{ // 8
			"let a=b?(c?d:e):f",
			"let a=b?(c?d:e):f;",
			"let a=b?( c ? d : e ):f;",
		},
		{ // 9
			"case a in b)c\nesac",
			"case a in\nb)\n\tc;;\nesac;",
			"case a in\nb)\n\tc;;\nesac;",
		},
		{ // 10
			"case a in b);;& c) d;&\nesac",
			"case a in\nb)\n\t;;&\nc)\n\td;&\nesac;",
			"case a in\nb)\n\t;;&\nc)\n\td;&\nesac;",
		},
		{ // 11
			"case a #A\nin #B\n#C\nb)c;;\n#D\nesac",
			"case a #A\nin #B\n#C\nb)\n\tc;;\n#D\nesac;",
			"case a #A\nin #B\n#C\nb)\n\tc;;\n#D\nesac;",
		},
		{ // 12
			"a=1 b=2 c d >e <f",
			"a=1 b=2 c d >e <f;",
			"a=1 b=2 c d >e <f;",
		},
		{ // 13
			"a b >c <d",
			"a b >c <d;",
			"a b >c <d;",
		},
		{ // 14
			">a <b",
			">a <b;",
			">a <b;",
		},
		{ // 15
			"a=1 b=2 >c <d",
			"a=1 b=2 >c <d;",
			"a=1 b=2 >c <d;",
		},
		{ // 16
			"a <<b\nheredoc\ncontents\nb",
			"a <<b;\nheredoc\ncontents\nb",
			"a <<b;\nheredoc\ncontents\nb",
		},
		{ // 17
			"a <<-b\n\theredoc\n\tcontents\nb",
			"a <<-b;\nheredoc\ncontents\nb",
			"a <<-b;\nheredoc\ncontents\nb",
		},
		{ // 18
			"$(a)",
			"$(a);",
			"$(a);",
		},
		{ // 19
			"$(a\nb)",
			"$(\n\ta;\n\tb;\n);",
			"$(\n\ta;\n\tb;\n);",
		},
		{ // 20
			"# A\na # B\n# C\n\n# D",
			"# A\na; # B\n# C\n\n# D",
			"# A\na; # B\n# C\n\n# D",
		},
		{ // 21
			"case a in\nesac <a 2>&1;",
			"case a in\nesac <a 2>&1;",
			"case a in\nesac <a 2>&1;",
		},
		{ // 22
			"a\nb\n\nc\nd\n\n\n\n\ne",
			"a;\nb;\n\nc;\nd;\n\ne;",
			"a;\nb;\n\nc;\nd;\n\ne;",
		},
		{ // 23
			"a\nfor b; do\nc\ndone",
			"a;\nfor b; do\n\tc;\ndone;",
			"a;\nfor b; do\n\tc;\ndone;",
		},
		{ // 24
			"for a;do c\nd\ndone",
			"for a; do\n\tc;\n\td;\ndone;",
			"for a; do\n\tc;\n\td;\ndone;",
		},
		{ // 25
			"for a in b\ndo c\ndone",
			"for a in b; do\n\tc;\ndone;",
			"for a in b; do\n\tc;\ndone;",
		},
		{ // 26
			"for a in b c\ndo d\ndone",
			"for a in b c; do\n\td;\ndone;",
			"for a in b c; do\n\td;\ndone;",
		},
		{ // 27
			"for ((a=0;a<1;a++));do b\ndone",
			"for ((a=0;a<1;a++)); do\n\tb;\ndone;",
			"for (( a = 0; a < 1; a ++ )); do\n\tb;\ndone;",
		},
		{ // 28
			"function a() { b; }",
			"function a() { b; };",
			"function a() {\n\tb;\n};",
		},
		{ // 29
			"function a { b; }",
			"function a() { b; };",
			"function a() {\n\tb;\n};",
		},
		{ // 30
			"a() { b; }",
			"a() { b; };",
			"a() {\n\tb;\n};",
		},
		{ // 31
			"function a() # A\n# B\n{ b; }",
			"function a() # A\n# B\n{ b; };",
			"function a() # A\n# B\n{\n\tb;\n};",
		},
		{ // 32
			"a() # A\n# B\n{ b; }",
			"a() # A\n# B\n{ b; };",
			"a() # A\n# B\n{\n\tb;\n};",
		},
		{ // 33
			"{ a; }",
			"{ a; };",
			"{\n\ta;\n};",
		},
		{ // 34
			"( a; )",
			"( a; );",
			"(\n\ta;\n);",
		},
		{ // 35
			"{ a; b; }",
			"{ a; b; };",
			"{\n\ta;\n\tb;\n};",
		},
		{ // 36
			"{ a || b; }",
			"{ a || b; };",
			"{\n\ta || b;\n};",
		},
		{ // 37
			"( a;\nb; )",
			"(\n\ta;\n\tb;\n);",
			"(\n\ta;\n\tb;\n);",
		},
		{ // 38
			"{ # A\na; # B\n}",
			"{ # A\n\ta; # B\n};",
			"{ # A\n\ta; # B\n};",
		},
		{ // 39
			"{ a; # A\n}",
			"{\n\ta; # A\n};",
			"{\n\ta; # A\n};",
		},
		{ // 40
			"{ a | b $(c\nd)\n}",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n};",
			"{\n\ta | b $(\n\t\tc;\n\t\td;\n\t);\n};",
		},
		{ // 41
			"{ declare a=$(b); }",
			"{ declare a=$(b); };",
			"{\n\tdeclare a=$(b);\n};",
		},
		{ // 42
			"{ declare a=$(b;c); }",
			"{ declare a=$(b; c;); };",
			"{\n\tdeclare a=$(\n\t\tb;\n\t\tc;\n\t);\n};",
		},
		{ // 43
			"{ a=1 b; }",
			"{ a=1 b; };",
			"{\n\ta=1 b;\n};",
		},
		{ // 44
			"{ a=$(b\nc) d; }",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n};",
			"{\n\ta=$(\n\t\tb;\n\t\tc;\n\t) d;\n};",
		},
		{ // 45
			"{ let a[$(b)]=c; }",
			"{ let a[$(b)]=c; };",
			"{\n\tlet a[$(b)]=c;\n};",
		},
		{ // 46
			"{ let a[$(b\nc)]=d; }",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n};",
			"{\n\tlet a[$(\n\t\tb;\n\t\tc;\n\t)]=d;\n};",
		},
		{ // 47
			"{ let a=$(b); }",
			"{ let a=$(b); };",
			"{\n\tlet a=$(b);\n};",
		},
		{ // 48
			"{ let a=$(b\nc); }",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n};",
			"{\n\tlet a=$(\n\t\tb;\n\t\tc;\n\t);\n};",
		},
		{ // 48
			"{ ${a[$(b\nc)]}; }",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n};",
			"{\n\t${a[$(\n\t\tb;\n\t\tc;\n\t)]};\n};",
		},
		{ // 49
			"<<a\nb$c\na",
			"<<a;\nb$c\na",
			"<<a;\nb$c\na",
		},
		{ // 50
			"{\n<<a\nb$c\na\n}",
			"{\n\t<<a;\nb$c\na\n};",
			"{\n\t<<a;\nb$c\na\n};",
		},
		{ // 51
			"{\n<<-a\nb$c\na\n}",
			"{\n\t<<-a;\n\tb$c\n\ta\n};",
			"{\n\t<<-a;\n\tb$c\n\ta\n};",
		},
		{ // 52
			"{ function a() { # A\nb; } }",
			"{\n\tfunction a() { # A\n\t\tb;\n\t};\n};",
			"{\n\tfunction a() { # A\n\t\tb;\n\t};\n};",
		},
		{ // 53
			"{ function a() { b;\nc; } }",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t};\n};",
			"{\n\tfunction a() {\n\t\tb;\n\t\tc;\n\t};\n};",
		},
		{ // 54
			"a | b <<c\nc",
			"a | b <<c;\nc",
			"a | b <<c;\nc",
		},
		{ // 55
			"a && b <<c\nc",
			"a && b <<c;\nc",
			"a && b <<c;\nc",
		},
		{ // 56
			"if a; then b;fi",
			"if a; then\n\tb;\nfi;",
			"if a; then\n\tb;\nfi;",
		},
		{ // 57
			"if a||b; then b\nc\nelif d\nthen\ne\nelse if f\nthen\ng\nfi\nfi",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
		},
		{ // 58
			"if a; then b;elif c; then d;elif e\nthen f;fi",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
		},
		{ // 59
			"while a\ndo\nb\ndone",
			"while a; do\n\tb;\ndone;",
			"while a; do\n\tb;\ndone;",
		},
		{ // 60
			"until a&&b; # A\n# B\ndo\n# C\nb\nc\ndone",
			"until a && b; # A\n# B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
			"until a && b; # A\n# B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
		},
		{ // 61
			"a[b]=",
			"a[b]=;",
			"a[b]=;",
		},
		{ // 62
			"a[b+c]=",
			"a[b+c]=;",
			"a[b + c]=;",
		},
		{ // 63
			"a[$(b\nc)]=",
			"a[$(\n\tb;\n\tc;\n)]=;",
			"a[$(\n\tb;\n\tc;\n)]=;",
		},
		{ // 64
			"${a}",
			"${a};",
			"${a};",
		},
		{ // 65
			"${!a}",
			"${!a};",
			"${!a};",
		},
		{ // 66
			"${#a}",
			"${#a};",
			"${#a};",
		},
		{ // 67
			"${a:=b}",
			"${a:=b};",
			"${a:=b};",
		},
		{ // 68
			"${a=b}",
			"${a=b};",
			"${a=b};",
		},
		{ // 69
			"${a:?b}",
			"${a:?b};",
			"${a:?b};",
		},
		{ // 70
			"${a?b}",
			"${a?b};",
			"${a?b};",
		},
		{ // 71
			"${a:+b}",
			"${a:+b};",
			"${a:+b};",
		},
		{ // 72
			"${a+b}",
			"${a+b};",
			"${a+b};",
		},
		{ // 73
			"${a:-b}",
			"${a:-b};",
			"${a:-b};",
		},
		{ // 74
			"${a-b}",
			"${a-b};",
			"${a-b};",
		},
		{ // 75
			"${a#b}",
			"${a#b};",
			"${a#b};",
		},
		{ // 76
			"${a##b}",
			"${a##b};",
			"${a##b};",
		},
		{ // 77
			"${a%b}",
			"${a%b};",
			"${a%b};",
		},
		{ // 78
			"${a%%b}",
			"${a%%b};",
			"${a%%b};",
		},
		{ // 79
			"${a/b}",
			"${a/b};",
			"${a/b};",
		},
		{ // 80
			"${a/b/c}",
			"${a/b/c};",
			"${a/b/c};",
		},
		{ // 81
			"${a//b}",
			"${a//b};",
			"${a//b};",
		},
		{ // 82
			"${a//b/c}",
			"${a//b/c};",
			"${a//b/c};",
		},
		{ // 83
			"${a/%b}",
			"${a/%b};",
			"${a/%b};",
		},
		{ // 84
			"${a/%b/c}",
			"${a/%b/c};",
			"${a/%b/c};",
		},
		{ // 85
			"${a/#b}",
			"${a/#b};",
			"${a/#b};",
		},
		{ // 86
			"${a/#b/c}",
			"${a/#b/c};",
			"${a/#b/c};",
		},
		{ // 87
			"${a^b}",
			"${a^b};",
			"${a^b};",
		},
		{ // 88
			"${a^^b}",
			"${a^^b};",
			"${a^^b};",
		},
		{ // 89
			"${a,b}",
			"${a,b};",
			"${a,b};",
		},
		{ // 90
			"${a,,b}",
			"${a,,b};",
			"${a,,b};",
		},
		{ // 91
			"${*}",
			"${*};",
			"${*};",
		},
		{ // 92
			"${@}",
			"${@};",
			"${@};",
		},
		{ // 93
			"${a:1}",
			"${a:1};",
			"${a:1};",
		},
		{ // 94
			"${a:1:2}",
			"${a:1:2};",
			"${a:1:2};",
		},
		{ // 95
			"${a: -1:2}",
			"${a: -1:2};",
			"${a: -1:2};",
		},
		{ // 96
			"${a@U}",
			"${a@U};",
			"${a@U};",
		},
		{ // 97
			"${a@u}",
			"${a@u};",
			"${a@u};",
		},
		{ // 98
			"${a@L}",
			"${a@L};",
			"${a@L};",
		},
		{ // 99
			"${a@Q}",
			"${a@Q};",
			"${a@Q};",
		},
		{ // 100
			"${a@E}",
			"${a@E};",
			"${a@E};",
		},
		{ // 101
			"${a@P}",
			"${a@P};",
			"${a@P};",
		},
		{ // 102
			"${a@A}",
			"${a@A};",
			"${a@A};",
		},
		{ // 103
			"${a@K}",
			"${a@K};",
			"${a@K};",
		},
		{ // 104
			"${a@a}",
			"${a@a};",
			"${a@a};",
		},
		{ // 105
			"${a@k}",
			"${a@k};",
			"${a@k};",
		},
		{ // 106
			"${!a@}",
			"${!a@};",
			"${!a@};",
		},
		{ // 107
			"${!a*}",
			"${!a*};",
			"${!a*};",
		},
		{ // 108
			"[[ a == b ]]",
			"[[ a == b ]];",
			"[[ a == b ]];",
		},
		{ // 109
			"[[ a == b$c ]]",
			"[[ a == b$c ]];",
			"[[ a == b$c ]];",
		},
		{ // 110
			"[[ a == b\"c\" ]]",
			"[[ a == b\"c\" ]];",
			"[[ a == b\"c\" ]];",
		},
		{ // 111
			"case a in a|b) a;\nb\nesac",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;",
			"case a in\na|b)\n\ta;\n\tb;;\nesac;",
		},
		{ // 112
			"case a in a|b|\"c\") a;\nb\nesac",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;",
			"case a in\na|b|\"c\")\n\ta;\n\tb;;\nesac;",
		},
		{ // 113
			"! a",
			"! a;",
			"! a;",
		},
		{ // 114
			"coproc a",
			"coproc a;",
			"coproc a;",
		},
		{ // 115
			"coproc a if b; then c\nfi",
			"coproc a if b; then\n\tc;\nfi;",
			"coproc a if b; then\n\tc;\nfi;",
		},
		{ // 116
			"select a; do b; done",
			"select a; do\n\tb;\ndone;",
			"select a; do\n\tb;\ndone;",
		},
		{ // 117
			"select a in b c; do b; done",
			"select a in b c; do\n\tb;\ndone;",
			"select a in b c; do\n\tb;\ndone;",
		},
		{ // 118
			"a&",
			"a&",
			"a &",
		},
	} {
		for m, input := range test {
			if m == 2 && (n == 41 || n == 34) {
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
