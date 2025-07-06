package validators

type Validator interface {
	Validate(link string) bool
}
