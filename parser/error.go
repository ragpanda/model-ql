package parser

type SemanticError struct {
	Err      string
	Position int
}

// complete error info
func NewSemanticError(errStr string) *SemanticError {
	return &SemanticError{Err: errStr}
}

func (self *SemanticError) Error() string {
	return self.Err
}
