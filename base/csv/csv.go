package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

func writeCsv(file *os.File) {
	w := csv.NewWriter(file)
	_ = w.Write([]string{"123", "456", "789", "999"})
	w.Flush()
	_ = file.Close()
}

func readCsv(file *os.File) {
	r := csv.NewReader(file)
	strs, _ := r.Read()
	for _, str := range strs {
		fmt.Println(str)
	}
}
