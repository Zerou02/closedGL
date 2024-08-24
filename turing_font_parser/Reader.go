package turingfontparser

import "os"

type Reader struct {
	path string
}

func NewReader(path string) Reader {
	return Reader{
		path: path,
	}
}

func (this *Reader) parse() {
	var f, err = os.Open(this.path)
	if err != nil {
		panic(err.Error())
	}
	_ = f
}
