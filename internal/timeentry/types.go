package timeentry

type SelectOption interface {
	SelectValue() string
	SelectString() string
}
