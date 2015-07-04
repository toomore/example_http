package utils

import (
	"encoding/csv"
	"io"
)

type csvData map[string]string

func CSV2map(r io.Reader) ([]csvData, error) {
	var (
		csvalldata [][]string
		csvmap     []csvData
		err        error
	)
	if csvalldata, err = csv.NewReader(r).ReadAll(); err == nil {
		csvmap = make([]csvData, len(csvalldata)-1)
		for i, v := range csvalldata[1:len(csvalldata)] {
			csvmap[i] = make(csvData)
			for mi, mv := range csvalldata[0] {
				if mv != "" {
					csvmap[i][mv] = v[mi]
				}
			}
		}
	}
	return csvmap, err
}
