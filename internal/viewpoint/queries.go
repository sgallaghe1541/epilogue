package viewpoint

type ViewpointResult interface {
	Headers() []string
	Data() [][]string
}

type ExcelableResult interface {
	ToExcel(string) error
}
