package files

import (
	"bufio"
	"io"
	"os"
	"sync"
)

type ReaderByLine struct{}

func NewReaderByLine() ReaderByLine {
	return ReaderByLine{}
}

func (r ReaderByLine) Read(filepath string, lineChannel chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	f, errorOpening := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if errorOpening != nil {
		panic(errorOpening)
	}
	defer f.Close()
	defer close(lineChannel)

	var lineRead string
	var errorReading error
	reader := bufio.NewReader(f)
	for {
		lineRead, errorReading = reader.ReadString('\n')
		if errorReading != nil {
			if errorReading == io.EOF {
				break
			}

			panic(errorReading)
		}
		lineChannel <- lineRead
	}
}
