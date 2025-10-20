// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/spkg/bom"
// )

// func main() {
// 	f, err := os.Open("messages.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	tempLine := ""
// 	var lines []string

// 	for {
// 		buff := make([]byte, 8)
// 		n, err := f.Read(buff)
// 		if err == io.EOF {
// 			break
// 		}

// 		if n > 0 {
// 			tempLine += string(bom.Clean(buff[:n]))
// 			normtemplLine := strings.ReplaceAll(tempLine, "\r\n", "\n")
// 			lines = strings.Split(normtemplLine, "\n")
// 		}
// 	}

// 	for _, v := range lines {
// 		fmt.Printf("read: %s\n", v)
// 	}
// }

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spkg/bom"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChan := make(chan string)

	go func() {
		defer close(lineChan)
		defer f.Close()

		var tempLine string
		buf := make([]byte, 8)

		for {
			n, err := f.Read(buf)
			if n > 0 {
				tempLine += string(bom.Clean(buf[:n]))
				normTempLine := strings.ReplaceAll(tempLine, "\r\n", "\n")
				lines := strings.Split(normTempLine, "\n")

				for i := 0; i < len(lines)-1; i++ {
					lineChan <- lines[i]
				}

				tempLine = lines[len(lines)-1]
			}

			if err == io.EOF {
				if len(tempLine) > 0 {
					lineChan <- tempLine
				}
				break
			}

			if err != nil {
				log.Printf("error reading file: %v", err)
				break
			}
		}
	}()

	return lineChan
}

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := getLinesChannel(f)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
