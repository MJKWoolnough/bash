# bash
--
    import "vimagination.zapto.org/bash"

Package bash implements a bash tokeniser and AST.

## Usage

```go
const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenComment
	TokenIdentifier
	TokenFunctionIdentifier
	TokenIdentifierAssign
	TokenLetIdentifierAssign
	TokenAssignment
	TokenKeyword
	TokenBuiltin
	TokenWord
	TokenNumberLiteral
	TokenString
	TokenStringStart
	TokenStringMid
	TokenStringEnd
	TokenBraceSequenceExpansion
	TokenBraceExpansion
	TokenBraceWord
	TokenPunctuator
	TokenHeredoc
	TokenHeredocIndent
	TokenHeredocEnd
	TokenOpenBacktick
	TokenCloseBacktick
	TokenPattern
	TokenOperator
	TokenBinaryOperator
)
```

```go
var (
	ErrInvalidCharacter          = errors.New("invalid character")
	ErrInvalidParameterExpansion = errors.New("invalid parameter expansion")
	ErrInvalidNumber             = errors.New("invalid number")
	ErrInvalidAssignment         = errors.New("invalid assignment")
	ErrMissingClosingBracket     = errors.New("missing closing bracket")
	ErrMissingClosingBrace       = errors.New("missing closing brace")
	ErrMissingClosingParen       = errors.New("missing closing paren")
	ErrMissingCloser             = errors.New("missing closer")
	ErrInvalidEndOfStatement     = errors.New("invalid end of statement")
	ErrIncorrectBacktick         = errors.New("incorrect backtick depth")
	ErrMissingWord               = errors.New("missing word")
	ErrMissingClosingIf          = errors.New("missing if closing")
	ErrMissingThen               = errors.New("missing then")
	ErrMissingIn                 = errors.New("missing in")
	ErrMissingDo                 = errors.New("missing do")
	ErrMissingClosingCase        = errors.New("missing case closing")
	ErrMissingClosingPattern     = errors.New("missing pattern closing")
	ErrInvalidKeyword            = errors.New("invalid keyword")
	ErrInvalidIdentifier         = errors.New("invalid identifier")
	ErrMissingOperator           = errors.New("missing operator")
	ErrInvalidOperator           = errors.New("invalid operator")
)
```
Errors.

#### func  SetTokeniser

```go
func SetTokeniser(t *parser.Tokeniser) *parser.Tokeniser
```
SetTokeniser sets the initial tokeniser state of a parser.Tokeniser.

Used if you want to manually tokenise bash code.

#### type ArithmeticExpansion

```go
type ArithmeticExpansion struct {
	Expression        bool
	WordsAndOperators []WordOrOperator
	Tokens            Tokens
}
```

ArithmeticExpansion represents either an expression ('((') or a compound
('$((').

For the expression, the returned number is the exit code, for the compound the
returned value is a word.

#### func (ArithmeticExpansion) Format

```go
func (f ArithmeticExpansion) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type ArrayWord

```go
type ArrayWord struct {
	Word     Word
	Comments [2]Comments
	Tokens   Tokens
}
```

ArrayWord a word in a Values array value.

The first set of comments are from just before the word, and the second set are
from just after.

#### func (ArrayWord) Format

```go
func (f ArrayWord) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Assignment

```go
type Assignment struct {
	Identifier ParameterAssign
	Assignment AssignmentType
	Expression []WordOrOperator
	Value      *Value
	Tokens     Tokens
}
```

Assignment represents a value assignment.

If Assignment is AssignmentAppend, Expression should be used, otherwise Value
should be set.

#### func (Assignment) Format

```go
func (f Assignment) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type AssignmentOrWord

```go
type AssignmentOrWord struct {
	Assignment *Assignment
	Word       *Word
	Tokens     Tokens
}
```

AssignmentOrWord represents either an Assignment or a Word in a command.

One, and only one, of Assignment or Word must be set.

#### func (AssignmentOrWord) Format

```go
func (f AssignmentOrWord) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type AssignmentType

```go
type AssignmentType uint8
```

AssignmentType represents the type of assignment, either a simple set or and
append.

```go
const (
	AssignmentAssign AssignmentType = iota
	AssignmentAppend
)
```
Assignment types.

#### func (AssignmentType) String

```go
func (a AssignmentType) String() string
```

#### type BraceExpansion

```go
type BraceExpansion struct {
	BraceExpansionType
	Words  []Word
	Tokens Tokens
}
```

BraceExpansion represents either a sequence expansion ('{a..b}', '{1..10..2}'),
or a group of words ('{ab,cd,12}').

#### func (BraceExpansion) Format

```go
func (f BraceExpansion) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type BraceExpansionType

```go
type BraceExpansionType uint8
```

BraceExpansionType represents which type of BraceExpansion is being represented.

```go
const (
	BraceExpansionWords BraceExpansionType = iota
	BraceExpansionSequence
)
```
Brace Expansion types.

#### func (BraceExpansionType) String

```go
func (b BraceExpansionType) String() string
```

#### type BraceWord

```go
type BraceWord struct {
	Parts  []WordPart
	Tokens Tokens
}
```


#### func (BraceWord) Format

```go
func (f BraceWord) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type CaseCompound

```go
type CaseCompound struct {
	Word     Word
	Matches  []PatternLines
	Comments [3]Comments
	Tokens   Tokens
}
```

CaseCompound represents a case select compound.

The first two comment groups represent comments on either side on the 'in'
keyword, and the third group represents comments from just before the closing
'esac' keyword.

#### func (CaseCompound) Format

```go
func (f CaseCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type CaseTerminationType

```go
type CaseTerminationType uint8
```

CaseTerminationType represents the final punctuation of a case match.

Must be one of CaseTerminationNone, CaseTerminationEnd, CaseTerminationContinue,
or CaseTerminationFallthrough.

```go
const (
	CaseTerminationNone CaseTerminationType = iota
	CaseTerminationEnd
	CaseTerminationContinue
	CaseTerminationFallthrough
)
```
CaseTermination types.

#### func (CaseTerminationType) String

```go
func (c CaseTerminationType) String() string
```

#### type Command

```go
type Command struct {
	Vars               []Assignment
	Redirections       []Redirection
	AssignmentsOrWords []AssignmentOrWord
	Tokens             Tokens
}
```

Command represents an assignment or a call to a command or builtin.

At least one Var, Redirection, or Word must be set.

#### func (Command) Format

```go
func (f Command) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type CommandOrCompound

```go
type CommandOrCompound struct {
	Command  *Command
	Compound *Compound
	Tokens   Tokens
}
```

CommandOrCompound represents either a Command or a Compound; one, and only one
of which must be set.

#### func (CommandOrCompound) Format

```go
func (f CommandOrCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type CommandSubstitution

```go
type CommandSubstitution struct {
	SubstitutionType SubstitutionType
	Backtick         *Token
	Command          File
	Tokens           Tokens
}
```

CommandSubstitution represents a subshell that returns some value.

For a SubstitutionNew or SubstitutionBacktick, the Standard Out is returned; for
a SubstitutionProcessInput or SubstitutionProcessOutput a path is return.

For a SubstitutionBacktick, the Backtick must be set to the escaped backtick
being used for the subshell.

The Command must contain at least one statement.

#### func (CommandSubstitution) Format

```go
func (f CommandSubstitution) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Comments

```go
type Comments []Token
```

Comments is a collection of Comment Tokens.

#### type Compound

```go
type Compound struct {
	IfCompound         *IfCompound
	CaseCompound       *CaseCompound
	LoopCompound       *LoopCompound
	ForCompound        *ForCompound
	SelectCompound     *SelectCompound
	GroupingCompound   *GroupingCompound
	TestCompound       *TestCompound
	ArithmeticCompound *ArithmeticExpansion
	FunctionCompound   *FunctionCompound
	Redirections       []Redirection
	Tokens             Tokens
}
```

Compound represents one of the Bash compound statements. One, and only of the
compounds must be set.

#### func (Compound) Format

```go
func (f Compound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Error

```go
type Error struct {
	Err     error
	Parsing string
	Token   Token
}
```

Error represents a Bash parsing error.

#### func (Error) Error

```go
func (e Error) Error() string
```
Error implements the error interface.

#### func (Error) Unwrap

```go
func (e Error) Unwrap() error
```
Unwrap returns the underlying error.

#### type File

```go
type File struct {
	Lines    []Line
	Comments [2]Comments
	Tokens   Tokens
}
```

File represents a parsed Bash file, a subshell, or a compound body.

The first set of comments are from the start of the file/body, the second set
from the end.

#### func  Parse

```go
func Parse(t Tokeniser) (*File, error)
```
Parse parses Bash input into AST.

#### func (File) Format

```go
func (f File) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type ForCompound

```go
type ForCompound struct {
	Identifier          *Token
	Words               []Word
	ArithmeticExpansion *ArithmeticExpansion
	File                File
	Comments            [2]Comments
	Tokens              Tokens
}
```

ForCompound represents a For loop.

One, and only one, of Identifier and ArithmeticExpansion must be set.

The File must contain at least one statement.

The first set of comments are from after an Identifier; the second set of
comments are from just before the 'do' keyword.

#### func (ForCompound) Format

```go
func (f ForCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type FunctionCompound

```go
type FunctionCompound struct {
	HasKeyword bool
	Identifier *Token
	Body       Compound
	Comments   Comments
	Tokens     Tokens
}
```

Function compound represents a defined function, either with or without the
'function' keyword.

The Comments are from just before the Body.

#### func (FunctionCompound) Format

```go
func (f FunctionCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type GroupingCompound

```go
type GroupingCompound struct {
	SubShell bool
	File
	Tokens Tokens
}
```

GroupingCompound represents either a brace or parenthesized set of statements.

File must contain at least one statement.

#### func (GroupingCompound) Format

```go
func (f GroupingCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Heredoc

```go
type Heredoc struct {
	HeredocPartsOrWords []HeredocPartOrWord
	Tokens              Tokens
}
```

Heredoc represents the parts of a Here Document.

#### func (Heredoc) Format

```go
func (f Heredoc) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type HeredocPartOrWord

```go
type HeredocPartOrWord struct {
	HeredocPart *Token
	Word        *Word
	Tokens      Tokens
}
```

HeredocPartOrWord represents either the string of Word part of a Here Document.

One of HeredocPart or Word must be set.

#### func (HeredocPartOrWord) Format

```go
func (f HeredocPartOrWord) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type IfCompound

```go
type IfCompound struct {
	If     TestConsequence
	ElIf   []TestConsequence
	Else   *File
	Tokens Tokens
}
```

IfCompound represents and if compound with optional elif and else sections.

#### func (IfCompound) Format

```go
func (f IfCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type JobControl

```go
type JobControl uint8
```

JobControl determines whether a job starts in the foreground or background.

```go
const (
	JobControlForeground JobControl = iota
	JobControlBackground
)
```

#### func (JobControl) String

```go
func (j JobControl) String() string
```

#### type Line

```go
type Line struct {
	Statements []Statement
	Comments   [2]Comments
	Tokens     Tokens
}
```

Line represents a logical bash line; it may contain multiple statements.

The first set of Comments are from just before the line, the second set from the
end of the line, and those on following lines not preceded by an empty line.

#### func (Line) Format

```go
func (f Line) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type LogicalOperator

```go
type LogicalOperator uint8
```

LogicalOperator represents how two statements are joined.

```go
const (
	LogicalOperatorNone LogicalOperator = iota
	LogicalOperatorAnd
	LogicalOperatorOr
)
```
Logical Operators.

#### func (LogicalOperator) String

```go
func (l LogicalOperator) String() string
```

#### type LoopCompound

```go
type LoopCompound struct {
	Until     bool
	Statement Statement
	File      File
	Comments  Comments
	Tokens    Tokens
}
```

LoopCompound represents either While or Until loops.

The File must contain at least one statement.

The comments are parsed after statement, before the 'do' keyword.

#### func (LoopCompound) Format

```go
func (f LoopCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Parameter

```go
type Parameter struct {
	Parameter *Token
	Array     []WordOrOperator
	Tokens    Tokens
}
```

Parameter represents the Parameter, an Identifier with a possible Array
subscript, used in a ParameterExpansion.

#### func (Parameter) Format

```go
func (f Parameter) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type ParameterAssign

```go
type ParameterAssign struct {
	Identifier *Token
	Subscript  []WordOrOperator
	Tokens     Tokens
}
```

ParameterAssign represents an identifier being assigned to, with a possible
subscript.

Identifier must be set.

#### func (ParameterAssign) Format

```go
func (f ParameterAssign) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type ParameterExpansion

```go
type ParameterExpansion struct {
	Indirect       bool
	Parameter      Parameter
	Type           ParameterType
	SubstringStart *Token
	SubstringEnd   *Token
	BraceWord      *BraceWord
	Pattern        *Token
	String         *String
	Tokens         Tokens
}
```

ParameterExpansion represents the expansion of a parameter.

#### func (ParameterExpansion) Format

```go
func (f ParameterExpansion) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type ParameterType

```go
type ParameterType uint8
```

ParameterType represents the type of a ParameterExpansion.

```go
const (
	ParameterValue ParameterType = iota
	ParameterLength
	ParameterSubstitution
	ParameterAssignment
	ParameterMessage
	ParameterSetAssign
	ParameterUnsetSubstitution
	ParameterUnsetAssignment
	ParameterUnsetMessage
	ParameterUnsetSetAssign
	ParameterSubstring
	ParameterPrefix
	ParameterPrefixSeperate
	ParameterRemoveStartShortest
	ParameterRemoveStartLongest
	ParameterRemoveEndShortest
	ParameterRemoveEndLongest
	ParameterReplace
	ParameterReplaceAll
	ParameterReplaceStart
	ParameterReplaceEnd
	ParameterLowercaseFirstMatch
	ParameterLowercaseAllMatches
	ParameterUppercaseFirstMatch
	ParameterUppercaseAllMatches
	ParameterUppercase
	ParameterUppercaseFirst
	ParameterLowercase
	ParameterQuoted
	ParameterEscaped
	ParameterPrompt
	ParameterDeclare
	ParameterQuotedArrays
	ParameterQuotedArraysSeperate
	ParameterAttributes
)
```
ParameterExpansion types.

#### func (ParameterType) String

```go
func (p ParameterType) String() string
```

#### type Pattern

```go
type Pattern struct {
	Parts  []WordPart
	Tokens Tokens
}
```

Pattern represents a pattern being matched against in a TestCompound test.

Must contain at least one WordPart.

#### func (Pattern) Format

```go
func (f Pattern) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type PatternLines

```go
type PatternLines struct {
	Patterns []Word
	Lines    File
	CaseTerminationType
	Comments Comments
	Tokens   Tokens
}
```

PatternLines represents a CaseCompound pattern and the code to run for that
match.

The Comments are parsed from just before the pattern.

#### func (PatternLines) Format

```go
func (f PatternLines) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Pipeline

```go
type Pipeline struct {
	PipelineTime      PipelineTime
	Not               bool
	Coproc            bool
	CoprocIdentifier  *Token
	CommandOrCompound CommandOrCompound
	Pipeline          *Pipeline
	Tokens            Tokens
}
```

Pipeline represents a command or compound, possibly connected to another
pipeline by a pipe ('|').

#### func (Pipeline) Format

```go
func (f Pipeline) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type PipelineTime

```go
type PipelineTime uint8
```

PipelineTime represents a potential 'time' keyword prefixed to a pipeline.

```go
const (
	PipelineTimeNone PipelineTime = iota
	PipelineTimeBash
	PipelineTimePosix
)
```
Pipeline Time options.

#### func (PipelineTime) String

```go
func (p PipelineTime) String() string
```

#### type Redirection

```go
type Redirection struct {
	Input      *Token
	Redirector *Token
	Output     Word
	Heredoc    *Heredoc
	Tokens     Tokens
}
```

Redirection presents input/output redirection.

Redirector must be set to the redirection operator.

#### func (Redirection) Format

```go
func (f Redirection) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type SelectCompound

```go
type SelectCompound struct {
	Identifier *Token
	Words      []Word
	File       File
	Comments   [2]Comments
	Tokens     Tokens
}
```

SelectCompound represents a Select loop.

The Identifier must be set and the File must contain at least one statement.

The first set of Comments is from just after the Identifier and the second are
from just before the 'do' keyword.

#### func (SelectCompound) Format

```go
func (f SelectCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Statement

```go
type Statement struct {
	Pipeline        Pipeline
	LogicalOperator LogicalOperator
	Statement       *Statement
	JobControl      JobControl
	Tokens
}
```

Statement represents a statement or statements joined by '||' or '&&' operators.

With a LogicalOperator set to either LogicalOperatorAnd or LogicalOperatorOr,
the Statement must be set.

#### func (Statement) Format

```go
func (f Statement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type String

```go
type String struct {
	WordsOrTokens []WordOrToken
	Tokens        Tokens
}
```

String represents a collection of string or word parts that make up string.

#### func (String) Format

```go
func (f String) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type SubstitutionType

```go
type SubstitutionType uint8
```

SubstitutionType represents the type of a CommandSubstitution.

```go
const (
	SubstitutionNew SubstitutionType = iota
	SubstitutionBacktick
	SubstitutionProcessInput
	SubstitutionProcessOutput
)
```
Substitution types.

#### func (SubstitutionType) String

```go
func (s SubstitutionType) String() string
```

#### type TestCompound

```go
type TestCompound struct {
	Tests    Tests
	Comments [2]Comments
	Tokens   Tokens
}
```

TestCompound represents the wrapping of a '[[ ... ]]' compound.

The first set of comments are from just after the opening '[[' and the second
set are from just before the closing ']]'.

#### func (TestCompound) Format

```go
func (f TestCompound) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type TestConsequence

```go
type TestConsequence struct {
	Test        Statement
	Consequence File
	Comments    Comments
	Tokens
}
```

TestConsequence represents the conditional test and body of an if or elif
section.

The Consequence must contain at least one statement.

The comments are parsed after the test statement, before the 'then' keyword.

#### func (TestConsequence) Format

```go
func (f TestConsequence) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type TestOperator

```go
type TestOperator uint8
```

TestOperator represents the type of test being represented.

```go
const (
	TestOperatorNone TestOperator = iota
	TestOperatorFileExists
	TestOperatorFileIsBlock
	TestOperatorFileIsCharacter
	TestOperatorDirectoryExists
	TestOperatorFileIsRegular
	TestOperatorFileHasSetGroupID
	TestOperatorFileIsSymbolic
	TestOperatorFileHasStickyBit
	TestOperatorFileIsPipe
	TestOperatorFileIsReadable
	TestOperatorFileIsNonZero
	TestOperatorFileIsTerminal
	TestOperatorFileHasSetUserID
	TestOperatorFileIsWritable
	TestOperatorFileIsExecutable
	TestOperatorFileIsOwnedByEffectiveGroup
	TestOperatorFileWasModifiedSinceLastRead
	TestOperatorFileIsOwnedByEffectiveUser
	TestOperatorFileIsSocket
	TestOperatorOptNameIsEnabled
	TestOperatorVarNameIsSet
	TestOperatorVarnameIsRef
	TestOperatorStringIsZero
	TestOperatorStringIsNonZero
	TestOperatorStringsEqual
	TestOperatorStringsMatch
	TestOperatorStringsNotEqual
	TestOperatorStringBefore
	TestOperatorStringAfter
	TestOperatorEqual
	TestOperatorNotEqual
	TestOperatorLessThan
	TestOperatorLessThanEqual
	TestOperatorGreaterThan
	TestOperatorGreaterThanEqual
	TestOperatorFilesAreSameInode
	TestOperatorFileIsNewerThan
	TestOperatorFileIsOlderThan
)
```

#### func (TestOperator) String

```go
func (t TestOperator) String() string
```

#### type Tests

```go
type Tests struct {
	Not             bool
	Test            TestOperator
	Word            *Word
	Pattern         *Pattern
	Parens          *Tests
	LogicalOperator LogicalOperator
	Tests           *Tests
	Comments        [5]Comments
	Tokens          Tokens
}
```

Tests represents the actual test conditions of a TestCompound.

#### func (Tests) Format

```go
func (f Tests) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Token

```go
type Token struct {
	parser.Token
	Pos, Line, LinePos uint64
}
```

Token represents a parser.Token combined with positioning information.

#### type Tokeniser

```go
type Tokeniser interface {
	Iter(func(parser.Token) bool)
	TokeniserState(parser.TokenFunc)
	GetError() error
}
```

Tokeniser represents the methods required by the bash tokeniser.

#### type Tokens

```go
type Tokens []Token
```

Tokens represents a list of tokens that have been parsed.

#### type Type

```go
type Type interface {
	fmt.Formatter
	// contains filtered or unexported methods
}
```

Type is an interface satisfied by all bash structural types.

#### type Value

```go
type Value struct {
	Word     *Word
	Array    []ArrayWord
	Comments [2]Comments
	Tokens   Tokens
}
```

Value represents the value to be assigned in an Assignment.

One, and only one, of Word or Array must be used.

When assigning an array, the first set of comments are from just after the
opening paren, and the second set of comments are from just before the closing
paren.

#### func (Value) Format

```go
func (f Value) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type Word

```go
type Word struct {
	Parts  []WordPart
	Tokens Tokens
}
```

Word represents a collection of WordParts that make up a single word.

#### func (Word) Format

```go
func (f Word) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type WordOrOperator

```go
type WordOrOperator struct {
	Word     *Word
	Operator *Token
	Tokens   Tokens
}
```

WordOrOperator represents either a Word or an Arithmetic Operator, one, and only
one of which must be set.

#### func (WordOrOperator) Format

```go
func (f WordOrOperator) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type WordOrToken

```go
type WordOrToken struct {
	Token  *Token
	Word   *Word
	Tokens Tokens
}
```

WordOrToken represents either a string token or a Word, one and only one of
which must be set.

#### func (WordOrToken) Format

```go
func (f WordOrToken) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.

#### type WordPart

```go
type WordPart struct {
	Part                *Token
	ParameterExpansion  *ParameterExpansion
	CommandSubstitution *CommandSubstitution
	ArithmeticExpansion *ArithmeticExpansion
	BraceExpansion      *BraceExpansion
	Tokens              Tokens
}
```

WordPart represents a single part of a word.

One and only one of Part, ParameterExpansion, CommandSubstitution,
ArithmeticExpansion, or BraceExpansion must be set.

#### func (WordPart) Format

```go
func (f WordPart) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface.
