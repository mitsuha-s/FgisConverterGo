package xls

type ParseError struct {
	message string
}

func (p ParseError) Error() string {
	return p.message
}
