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
			"for a;do c\nd\ndone",
			"for a; do\n\tc;\n\td;\ndone;",
			"for a; do\n\tc;\n\td;\ndone;",
		},
		{ // 24
			"for a in b\ndo c\ndone",
			"for a in b; do\n\tc;\ndone;",
			"for a in b; do\n\tc;\ndone;",
		},
		{ // 25
			"for a in b c\ndo d\ndone",
			"for a in b c; do\n\td;\ndone;",
			"for a in b c; do\n\td;\ndone;",
		},
		{ // 26
			"for ((a=0;a<1;a++));do b\ndone",
			"for ((a=0;a<1;a++)); do\n\tb;\ndone;",
			"for (( a = 0; a < 1; a ++ )); do\n\tb;\ndone;",
		},
		{ // 27
			"function a() { b; }",
			"function a() { b; };",
			"function a() {\n\tb;\n};",
		},
		{ // 28
			"function a { b; }",
			"function a() { b; };",
			"function a() {\n\tb;\n};",
		},
		{ // 29
			"a() { b; }",
			"a() { b; };",
			"a() {\n\tb;\n};",
		},
		{ // 30
			"function a() # A\n# B\n{ b; }",
			"function a() # A\n# B\n{ b; };",
			"function a() # A\n# B\n{\n\tb;\n};",
		},
		{ // 31
			"a() # A\n# B\n{ b; }",
			"a() # A\n# B\n{ b; };",
			"a() # A\n# B\n{\n\tb;\n};",
		},
		{ // 32
			"{ a; }",
			"{ a; };",
			"{\n\ta;\n};",
		},
		{ // 33
			"( a; )",
			"( a; );",
			"(\n\ta;\n);",
		},
		{ // 34
			"{ a; b; }",
			"{ a; b; };",
			"{\n\ta;\n\tb;\n};",
		},
		{ // 35
			"( a;\nb; )",
			"(\n\ta;\n\tb;\n);",
			"(\n\ta;\n\tb;\n);",
		},
		{ // 36
			"{ # A\na; # B\n}",
			"{ # A\n\ta; # B\n};",
			"{ # A\n\ta; # B\n};",
		},
		{ // 37
			"<<a\nb$c\na",
			"<<a;\nb$c\na",
			"<<a;\nb$c\na",
		},
		{ // 38
			"{\n<<a\nb$c\na\n}",
			"{\n\t<<a;\nb$c\na\n};",
			"{\n\t<<a;\nb$c\na\n};",
		},
		{ // 39
			"{\n<<-a\nb$c\na\n}",
			"{\n\t<<-a;\n\tb$c\n\ta\n};",
			"{\n\t<<-a;\n\tb$c\n\ta\n};",
		},
		{ // 40
			"a | b <<c\nc",
			"a | b <<c;\nc",
			"a | b <<c;\nc",
		},
		{ // 41
			"a && b <<c\nc",
			"a && b <<c;\nc",
			"a && b <<c;\nc",
		},
		{ // 42
			"if a; then b;fi",
			"if a; then\n\tb;\nfi;",
			"if a; then\n\tb;\nfi;",
		},
		{ // 43
			"if a||b; then b\nc\nelif d\nthen\ne\nelse if f\nthen\ng\nfi\nfi",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
			"if a || b; then\n\tb;\n\tc;\nelif d; then\n\te;\nelse\n\tif f; then\n\t\tg;\n\tfi;\nfi;",
		},
		{ // 44
			"if a; then b;elif c; then d;elif e\nthen f;fi",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
			"if a; then\n\tb;\nelif c; then\n\td;\nelif e; then\n\tf;\nfi;",
		},
		{ // 45
			"while a\ndo\nb\ndone",
			"while a; do\n\tb;\ndone;",
			"while a; do\n\tb;\ndone;",
		},
		{ // 46
			"until a&&b; # A\n# B\ndo\n# C\nb\nc\ndone",
			"until a && b; # A\n# B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
			"until a && b; # A\n# B\ndo\n\t# C\n\tb;\n\tc;\ndone;",
		},
	} {
		for m, input := range test {
			if m == 2 && n == 33 {
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
