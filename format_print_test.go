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
			"<<a\nb$c\na",
			"<<a;\nb$c\na",
			"<<a;\nb$c\na",
		},
		{ // 42
			"{\n<<a\nb$c\na\n}",
			"{\n\t<<a;\nb$c\na\n};",
			"{\n\t<<a;\nb$c\na\n};",
		},
		{ // 43
			"{\n<<-a\nb$c\na\n}",
			"{\n\t<<-a;\n\tb$c\n\ta\n};",
			"{\n\t<<-a;\n\tb$c\n\ta\n};",
		},
		{ // 44
			"a | b <<c\nc",
			"a | b <<c;\nc",
			"a | b <<c;\nc",
		},
		{ // 45
			"a && b <<c\nc",
			"a && b <<c;\nc",
			"a && b <<c;\nc",
		},
		{ // 46
			"if a; then b;fi",
			"if a; then\n\tb;\nfi;",
			"if a; then\n\tb;\nfi;",
		},
		{ // 47
			"if a||b; then b\nc\nelif d\nthen\ne\nelse if f\nthen\ng\nfi\nfi",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
		},
		{ // 48
			"if a; then b;elif c; then d;elif e\nthen f;fi",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
		},
		{ // 49
			"while a\ndo\nb\ndone",
			"while a; do\n\tb;\ndone;",
			"while a; do\n\tb;\ndone;",
		},
		{ // 50
			"until a&&b; # A\n# B\ndo\n# C\nb\nc\ndone",
			"until a && b; # A\n# B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
			"until a && b; # A\n# B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
		},
		{ // 51
			"a[b]=",
			"a[b]=;",
			"a[b]=;",
		},
		{ // 52
			"a[b+c]=",
			"a[b+c]=;",
			"a[b + c]=;",
		},
		{ // 53
			"a[$(b\nc)]=",
			"a[$(\n\tb;\n\tc;\n)]=;",
			"a[$(\n\tb;\n\tc;\n)]=;",
		},
		{ // 54
			"${a}",
			"${a};",
			"${a};",
		},
		{ // 55
			"${!a}",
			"${!a};",
			"${!a};",
		},
		{ // 56
			"${#a}",
			"${#a};",
			"${#a};",
		},
		{ // 57
			"${a:=b}",
			"${a:=b};",
			"${a:=b};",
		},
		{ // 58
			"${a=b}",
			"${a=b};",
			"${a=b};",
		},
		{ // 59
			"${a:?b}",
			"${a:?b};",
			"${a:?b};",
		},
		{ // 60
			"${a?b}",
			"${a?b};",
			"${a?b};",
		},
		{ // 61
			"${a:+b}",
			"${a:+b};",
			"${a:+b};",
		},
		{ // 62
			"${a+b}",
			"${a+b};",
			"${a+b};",
		},
		{ // 63
			"${a:-b}",
			"${a:-b};",
			"${a:-b};",
		},
		{ // 64
			"${a-b}",
			"${a-b};",
			"${a-b};",
		},
		{ // 65
			"${a#b}",
			"${a#b};",
			"${a#b};",
		},
		{ // 66
			"${a##b}",
			"${a##b};",
			"${a##b};",
		},
		{ // 67
			"${a%b}",
			"${a%b};",
			"${a%b};",
		},
		{ // 68
			"${a%%b}",
			"${a%%b};",
			"${a%%b};",
		},
		{ // 69
			"${a/b}",
			"${a/b};",
			"${a/b};",
		},
		{ // 70
			"${a/b/c}",
			"${a/b/c};",
			"${a/b/c};",
		},
		{ // 71
			"${a//b}",
			"${a//b};",
			"${a//b};",
		},
		{ // 72
			"${a//b/c}",
			"${a//b/c};",
			"${a//b/c};",
		},
		{ // 73
			"${a/%b}",
			"${a/%b};",
			"${a/%b};",
		},
		{ // 74
			"${a/%b/c}",
			"${a/%b/c};",
			"${a/%b/c};",
		},
		{ // 75
			"${a/#b}",
			"${a/#b};",
			"${a/#b};",
		},
		{ // 76
			"${a/#b/c}",
			"${a/#b/c};",
			"${a/#b/c};",
		},
		{ // 77
			"${a^b}",
			"${a^b};",
			"${a^b};",
		},
		{ // 78
			"${a^^b}",
			"${a^^b};",
			"${a^^b};",
		},
		{ // 79
			"${a,b}",
			"${a,b};",
			"${a,b};",
		},
		{ // 80
			"${a,,b}",
			"${a,,b};",
			"${a,,b};",
		},
		{ // 81
			"${*}",
			"${*};",
			"${*};",
		},
		{ // 82
			"${@}",
			"${@};",
			"${@};",
		},
		{ // 83
			"${a:1}",
			"${a:1};",
			"${a:1};",
		},
		{ // 84
			"${a:1:2}",
			"${a:1:2};",
			"${a:1:2};",
		},
		{ // 85
			"${a: -1:2}",
			"${a: -1:2};",
			"${a: -1:2};",
		},
		{ // 86
			"${a@U}",
			"${a@U};",
			"${a@U};",
		},
		{ // 87
			"${a@u}",
			"${a@u};",
			"${a@u};",
		},
		{ // 88
			"${a@L}",
			"${a@L};",
			"${a@L};",
		},
		{ // 89
			"${a@Q}",
			"${a@Q};",
			"${a@Q};",
		},
		{ // 90
			"${a@E}",
			"${a@E};",
			"${a@E};",
		},
		{ // 91
			"${a@P}",
			"${a@P};",
			"${a@P};",
		},
		{ // 92
			"${a@A}",
			"${a@A};",
			"${a@A};",
		},
		{ // 93
			"${a@K}",
			"${a@K};",
			"${a@K};",
		},
		{ // 94
			"${a@a}",
			"${a@a};",
			"${a@a};",
		},
		{ // 95
			"${a@k}",
			"${a@k};",
			"${a@k};",
		},
	} {
		for m, input := range test {
			if m == 2 && (n == 34 || n == 35) {
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
