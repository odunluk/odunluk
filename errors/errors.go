package errors

type ValidationError struct {
	Detail interface{}
	Code   string
}

func (v *ValidationError) Error() string {
	return "Validation Error"
}

type SkipField struct {

}

func (s *SkipField) Error() string {
	return "Skip Field"
}