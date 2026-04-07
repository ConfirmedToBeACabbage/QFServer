package FR

import (
	"fmt"
	"os"
)

// Reading a file means
/*
	1. Get file input
	2. Check if access to file is possible
	3. Hold pointer to file
	4. Read in file buffer and store it into inner buffer
	5. Hold the information in a cache
	6. Clear the buffer
*/

type filereader struct {
	inputFile   *os.File
	inputBuffer []byte
	inputCache  [][]byte
}

func ReadFromFile(filePath string) []byte {

	fileReadObj := &filereader{}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("File can't be opened")
	} else {
		fileReadObj.inputFile = file

		var countKB int = 0
		var number int64 = 0
		EOF := false
		for !EOF {
			fileReadObj.inputBuffer = make([]byte, 1000)
			readInt, err := fileReadObj.inputFile.ReadAt(fileReadObj.inputBuffer, number)
			countKB += 1

			if err != nil {
				EOF = true
			}

			fileReadObj.inputCache[number] = fileReadObj.inputBuffer
			number = int64(readInt)
		}

		out := make([]byte, 0, countKB+1)
		for _, p := range fileReadObj.inputCache {
			out = append(out, p...)
		}

		return out

	}

	return nil
}
