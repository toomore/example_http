package utils

import (
	"encoding/csv"
	"io"
	"net/url"
)

type CSVData map[string]string

func CSV2map(r io.Reader) ([]CSVData, error) {
	var (
		csvalldata [][]string
		csvmap     []CSVData
		err        error
	)
	if csvalldata, err = csv.NewReader(r).ReadAll(); err == nil {
		csvmap = make([]CSVData, len(csvalldata)-1)
		for i, v := range csvalldata[1:len(csvalldata)] {
			csvmap[i] = make(CSVData)
			for mi, mv := range csvalldata[0] {
				if mv != "" {
					csvmap[i][mv] = v[mi]
				}
			}
		}
	}
	return csvmap, err
}

func Map2Values(m []CSVData) []url.Values {
	r := make([]url.Values, len(m))
	for i, d := range m {
		r[i] = url.Values{}
		for k, v := range d {
			r[i].Set(k, v)
		}
	}
	return r
}

func Map2ValuesMust(m []CSVData, err error) []url.Values {
	return Map2Values(m)
}
