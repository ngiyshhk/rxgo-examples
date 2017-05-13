package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
)

type CsvParser struct {
	Delimiter string
	Columns   []int
}

type CsvParserFactory struct{}

func (_ CsvParserFactory) Create() *CsvParser {
	str_cols := flag.String("C", "", "display columns")
	is_tsv := flag.Bool("t", false, "tsv")
	flag.Parse()

	delimiter := ","
	if *is_tsv {
		delimiter = "\t"
	}

	arr_cols := strings.Split(*str_cols, ",")
	columns := make([]int, 0)
	for _, col := range arr_cols {
		if int_col, err := strconv.Atoi(col); err == nil {
			columns = append(columns, int_col)
		}
	}

	return &CsvParser{Delimiter: delimiter, Columns: columns}
}

// 標準入力をcsvにパースしてobservableに
func (cp CsvParser) Src(r io.Reader) observable.Observable {
	reader := csv.NewReader(r)
	reader.Comma = []rune(cp.Delimiter)[0]

	ch := make(chan interface{})
	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				ch <- err
				break
			}

			ch <- record
		}
	}()
	it, _ := iterable.New(ch)
	return observable.From(it)
}

// 出力カラム指定されていたら、そのカラムだけのレコードに変換
func (cp CsvParser) Convert(v interface{}) interface{} {
	if len(cp.Columns) == 0 {
		return v
	}

	record, ok := v.([]string)
	if !ok || len(record) == 0 {
		return v
	}

	result := make([]string, 0)
	for _, col := range cp.Columns {
		result = append(result, record[col])
	}
	return result
}

// 標準出力
func (cp CsvParser) Sink() observer.Observer {
	writer := csv.NewWriter(os.Stdout)
	return observer.Observer{
		NextHandler: func(v interface{}) {
			item := v.([]string)
			writer.Write(item)
			writer.Flush()
		},

		ErrHandler: func(err error) {
			if err.Error() != "EOF" {
				fmt.Printf("Encountered error: %v\n", err)
			}
		},
	}
}

func main() {
	tg := CsvParserFactory{}.Create()
	<-tg.Src(os.Stdin).Map(tg.Convert).Subscribe(tg.Sink())
}
