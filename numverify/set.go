package numverify

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gophonenumbers/common"
	"github.com/grokify/mogo/io/ioutilmore"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/rs/zerolog/log"
)

var MultiLimit = 0 // test limit to gracefully exit process early.

// MultiResults is designed to handle large volumes of requests.
type MultiResults struct {
	Counts    map[string]int
	Responses map[string]*Response
}

func NewMultiResults() MultiResults {
	return MultiResults{
		Counts:    map[string]int{},
		Responses: map[string]*Response{}}
}

func (mr *MultiResults) Inflate() {
	counts := map[string]int{}
	for _, resp := range mr.Responses {
		scStr := strconv.Itoa(resp.StatusCode)
		if _, ok := counts[scStr]; !ok {
			counts[scStr] = 0
		}
		counts[scStr]++
	}
	counts["all"] = len(mr.Responses)
	mr.Counts = counts
}

func (mr *MultiResults) AddResponses(resps map[string]*Response) {
	for k, v := range resps {
		existing, ok := mr.Responses[k]
		if !ok ||
			(existing.StatusCode >= 300 && v.StatusCode < 300) {
			mr.Responses[k] = v
		}
	}
}

func (mr *MultiResults) NumbersSuccess() []string {
	numbers := []string{}
	for _, resp := range mr.Responses {
		if resp.Success != nil {
			pn := strings.TrimSpace(resp.Success.InternationalFormat)
			if len(pn) > 0 {
				numbers = append(numbers, pn)
			}
		}
	}
	return numbers
}

func (mr *MultiResults) Write(filename string) error {
	mr.Inflate()
	bytes, err := json.MarshalIndent(mr, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0600)
}

func GetWriteValidationMulti(client *Client, requestNumbers, skipNumbers []string, filenameBase string, logAt, fileAt uint) MultiResults {
	uniquesRequests := stringsutil.SliceCondenseSpace(requestNumbers, true, true)
	uniqueSkips := stringsutil.SliceCondenseSpace(skipNumbers, true, true)
	uniqueSkipsMap := map[string]int{}
	for _, pnSkip := range uniqueSkips {
		uniqueSkipsMap[pnSkip] = 1
	}

	resps := NewMultiResults()
	count := len(uniquesRequests)
	i := 0
	for _, num := range uniquesRequests {
		i++
		if _, ok := uniqueSkipsMap[num]; ok {
			continue
		}
		validate, _, _ := client.Validate(
			Params{Number: num})
		resps.Responses[num] = validate
		if logAt > 0 && i%int(logAt) == 0 {
			/*apiStatus := "S"
			if validate.StatusCode >= 300 || validate.Success == nil {
				apiStatus = "F"
			}
			log.Infof("[%v/%v][%v][%v][%s]", i, count, num, apiStatus,
				time.Now().Format(time.RFC3339))*/
			log.Info().
				Int("num", i).
				Int("count", count).
				Str("e164number", num).
				Int("httpStatus", validate.StatusCode).
				Msg("logAt")
		}
		if fileAt > 0 && i%int(fileAt) == 0 && len(resps.Responses) > 0 {
			err := resps.Write(common.BuildFilename(filenameBase, i, count))
			if err != nil {
				log.Error().Err(err)
			}
			resps = NewMultiResults()
		}
		if MultiLimit > 0 && i > MultiLimit {
			break
		}
	}
	if len(resps.Responses) > 0 {
		err := resps.Write(common.BuildFilename(filenameBase, i, count))
		if err != nil {
			log.Error().Err(err)
		}
	}
	return resps
}

func ReadFilesMultiResults(dir string, rxPattern string) (MultiResults, error) {
	dir = strings.TrimSpace(dir)
	if len(dir) == 0 {
		dir = "."
	}
	all := NewMultiResults()
	rx, err := regexp.Compile(rxPattern)
	if err != nil {
		return all, err
	}
	files, err := osutil.ReadDirMore(dir, rx, false, true, false)
	if err != nil {
		return all, err
	}
	for _, entry := range files {
		file := filepath.Join(dir, entry.Name())
		mResults := NewMultiResults()
		err := ioutilmore.ReadFileJSON(file, &mResults)
		if err != nil {
			return all, err
		}
		fi, err := entry.Info()
		if err != nil {
			return all, err
		}
		fileModTime := fi.ModTime()
		for key, ni := range mResults.Responses {
			if timeutil.TimeIsZeroAny(ni.Time) {
				ni.Time = fileModTime
				mResults.Responses[key] = ni
			}
		}
		all.AddResponses(mResults.Responses)
	}
	return all, nil
}

func GetNumbers(nvClient Client, filebase string, byNumber *histogram.Histogram) error {
	existing, err := ReadFilesMultiResults(".", filebase+`_\d+\-\d+\.json$`)
	if err != nil {
		return err
	}
	skipNumbers := existing.NumbersSuccess()

	wantNumbers := []string{}
	if byNumber != nil {
		for number := range byNumber.Bins {
			wantNumbers = append(wantNumbers, number)
		}
	}
	GetWriteValidationMulti(&nvClient, wantNumbers, skipNumbers, filebase, uint(20), uint(5000))

	return nil
}
