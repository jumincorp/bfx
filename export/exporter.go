package export

type Exporter interface {
	Setup()
	Export(ctype string, gpu string, name string, value string) error
}
