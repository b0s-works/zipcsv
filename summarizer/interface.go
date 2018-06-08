package summarizer

type DataProcessor interface {
	AddValue(string) error
	Result() string
}
