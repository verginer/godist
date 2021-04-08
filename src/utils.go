package godist

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"
)


// NewTransactionDate returns a time.Time object created from a string
func NewTransactionDate(date, dateFormat string) time.Time {
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// lineCounter Counts the lines a the given file source:
// https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently
func lineCounter(path string) (int, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	r := io.Reader(file)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
