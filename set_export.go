package gophonenumbers

import (
	"encoding/json"
	"strings"

	"github.com/grokify/gocharts/data/table"
	"github.com/pkg/errors"
)

func NumbersSetToTable(numsSet *NumbersSet) table.Table {
	tbl := table.NewTable()
	tbl.Columns = []string{
		"number", "lineType", "carrier", "lineTypes", "carriers", "raw"}
	for _, num := range numsSet.NumbersMap {
		tbl.Rows = append(tbl.Rows, numberInfoToRow(num))
	}
	return tbl
}

func numberInfoToRow(num NumberInfo) []string {
	raw, err := json.Marshal(num)
	if err != nil {
		panic(errors.Wrap(err, "gophonenumbers.numberInfoToRow"))
	}
	return []string{
		num.NumberE164,
		num.Carrier.LineType,
		num.Carrier.Name,
		strings.Join(num.LineTypesEach, ","),
		strings.Join(num.CarrierNamesEach, ","),
		string(raw)}
}
