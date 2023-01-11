package config

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func validateConf(conf CsvConf) error {
	if conf.Order == "random" || conf.Order == "sequential" {
		return nil
	}
	return fmt.Errorf("unsupported order %s, should be random|sequential", conf.Order)
}

func readCsv(conf CsvConf) ([]map[string]interface{}, error) {
	err := validateConf(conf)
	if err != nil {
		return nil, err
	}

	if conf.Src == "local" {
		f, err := os.Open(conf.Path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		csvReader.Comma = []rune(conf.Delimiter)[0]
		csvReader.TrimLeadingSpace = true
		csvReader.LazyQuotes = conf.AllowQuota

		data, err := csvReader.ReadAll()
		if err != nil {
			return nil, err
		}

		if conf.SkipFirstLine {
			data = data[1:]
		}

		rt := make([]map[string]interface{}, 0) // unclear how many empty line exist

		for _, row := range data {
			if conf.SkipEmptyLine && emptyLine(row) {
				continue
			}
			x := map[string]interface{}{}
			for index, tag := range conf.Vars { // "0":"name", "1":"city","2":"team"
				i, err := strconv.Atoi(index)
				if err != nil {
					return nil, err
				}
				// convert
				var val interface{}
				switch tag.Type {
				case "json":
					json.Unmarshal([]byte(row[i]), &val)
				case "int":
					var err error
					val, err = strconv.Atoi(row[i])
					if err != nil {
						return nil, err
					}
				case "float":
					var err error
					val, err = strconv.ParseFloat(row[i], 64)
					if err != nil {
						return nil, err
					}
				case "bool":
					var err error
					val, err = strconv.ParseBool(row[i])
					if err != nil {
						return nil, err
					}
				default:
					val = row[i]
				}
				x[tag.Tag] = val
			}
			rt = append(rt, x)
		}

		return rt, nil

	} else if conf.Src == "remote" {
		// TODOcorr, http call
	}

	return nil, fmt.Errorf("csv read error")
}

func emptyLine(row []string) bool {
	for _, field := range row {
		if field != "" {
			return false
		}
	}
	return true
}
