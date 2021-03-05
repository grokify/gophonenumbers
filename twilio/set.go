package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/gophonenumbers/common"
	"github.com/grokify/simplego/io/ioutilmore"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/rs/zerolog/log"
)

var MultiLimit = 0 // test limit to gracefully exit process early.

// MultiResults is designed to handle large volumes
// of requests.
type MultiResults struct {
	CountsByStatusCode map[string]int
	Responses          map[string]*NumberInfo
}

func NewMultiResults() MultiResults {
	return MultiResults{
		CountsByStatusCode: map[string]int{},
		Responses:          map[string]*NumberInfo{}}
}

func (mr *MultiResults) Inflate() {
	counts := map[string]int{}
	for _, resp := range mr.Responses {
		scStr := strconv.Itoa(resp.ApiResponseInfo.StatusCode)
		if _, ok := counts[scStr]; !ok {
			counts[scStr] = 0
		}
		counts[scStr]++
	}
	counts["all"] = len(mr.Responses)
	mr.CountsByStatusCode = counts
}

func (mr *MultiResults) AddResponses(resps map[string]*NumberInfo) {
	for k, v := range resps {
		existing, ok := mr.Responses[k]
		if !ok ||
			(existing.ApiResponseInfo.StatusCode >= 300 && v.ApiResponseInfo.StatusCode < 300) {
			mr.Responses[k] = v
		}
	}
}

func (mr *MultiResults) GetNumberInfo(e164Number string) (*NumberInfo, error) {
	e164Number = strings.TrimSpace(e164Number)
	if ni, ok := mr.Responses[e164Number]; ok {
		return ni, nil
	}
	return nil, fmt.Errorf("number [%s] not found", e164Number)
}

func (mr *MultiResults) NumbersSuccess() []string {
	numbers := []string{}
	for _, resp := range mr.Responses {
		pn := strings.TrimSpace(resp.PhoneNumber)
		if len(pn) > 0 {
			numbers = append(numbers, pn)
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
	return ioutil.WriteFile(filename, bytes, 0644)
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
	for _, e164Number := range uniquesRequests {
		i++
		if _, ok := uniqueSkipsMap[e164Number]; ok {
			continue
		}
		validate, _ := client.Validate(
			e164Number, &Params{Type: "carrier"})
		resps.Responses[e164Number] = &validate
		if logAt > 0 && i%int(logAt) == 0 {
			/*apiStatus := "S"
			if validate.ApiResponseInfo.StatusCode >= 300 {
				apiStatus = "F"
			}*/
			log.Info().
				Int("num", i).
				Int("count", count).
				Str("e164number", e164Number).
				Int("httpStatus", validate.ApiResponseInfo.StatusCode).
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

func NewMultiResultsFiles(dir string, rxPattern string) (MultiResults, error) {
	dir = strings.TrimSpace(dir)
	if len(dir) == 0 {
		dir = "."
	}
	all := NewMultiResults()
	rx, err := regexp.Compile(rxPattern)
	if err != nil {
		return all, err
	}
	files, _, err := ioutilmore.ReadDirRx(dir, rx, true)
	if err != nil {
		return all, err
	}
	for _, fi := range files {
		mResults := NewMultiResults()
		err := ioutilmore.ReadFileJSON(
			filepath.Join(dir, fi.Name()), &mResults)
		if err != nil {
			return all, err
		}
		fileModTime := fi.ModTime()
		for key, ni := range mResults.Responses {
			if timeutil.TimeIsZeroAny(ni.ApiResponseInfo.Time) {
				ni.ApiResponseInfo.Time = fileModTime
				mResults.Responses[key] = ni
			}
		}
		all.AddResponses(mResults.Responses)
	}
	return all, nil
}
