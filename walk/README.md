# walk
--
    import "vimagination.zapto.org/bash/walk"

Package walk provides a bash type walker.

## Usage

#### func  Walk

```go
func Walk(t bash.Type, fn Handler) error
```
Walk calls the Handle function on the given interface for each non-nil,
non-Token field of the given bash type.

#### type Handler

```go
type Handler interface {
	Handle(bash.Type) error
}
```

Handler is used to process bash types.

#### type HandlerFunc

```go
type HandlerFunc func(bash.Type) error
```

Handle implements the Handler interface.

#### func (HandlerFunc) Handle

```go
func (h HandlerFunc) Handle(t bash.Type) error
```
