package bash

import "io"

func (ArithmeticExpansion) printSource(w io.Writer, v bool) {}
func (Assignment) printSource(w io.Writer, v bool)          {}
func (Command) printSource(w io.Writer, v bool)             {}
func (CommandSubstitution) printSource(w io.Writer, v bool) {}
func (File) printSource(w io.Writer, v bool)                {}
func (ParameterAssign) printSource(w io.Writer, v bool)     {}
func (ParameterExpansion) printSource(w io.Writer, v bool)  {}
func (Parameter) printSource(w io.Writer, v bool)           {}
func (Pipeline) printSource(w io.Writer, v bool)            {}
func (Redirection) printSource(w io.Writer, v bool)         {}
func (Statement) printSource(w io.Writer, v bool)           {}
func (String) printSource(w io.Writer, v bool)              {}
func (Value) printSource(w io.Writer, v bool)               {}
func (WordPart) printSource(w io.Writer, v bool)            {}
func (Word) printSource(w io.Writer, v bool)                {}
