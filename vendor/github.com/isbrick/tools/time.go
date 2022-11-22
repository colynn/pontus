// Copyright 2020 tools authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Date Format unix time int64 to string
// eg format: YYYY-MM-DD HH:mm:SS
func Date(ti int64, format string) string {
	t := time.Unix(int64(ti), 0)
	return DateT(t, format)
}

// DateS Format unix time string to string
func DateS(ts string, format string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	return Date(i, format)
}

// DateT ..
// Format time.Time struct to string
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
func DateT(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

// ParseStrToDate ..
// timeStr
// format
// locationName, when locationName is "", defaultLocation rely on system env.
func ParseStrToDate(timestr, format, locationName string) (time.Time, error) {
	timeFormat := timeFormatParse(format)
	loc := gettimeLocation(locationName)
	timeValue, err := time.ParseInLocation(timeFormat, timestr, loc)
	return timeValue, err
}

// ParseDateToStr ..
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
func ParseDateToStr(timeValue time.Time, format string) string {
	formatStr := timeFormatParse(format)
	return timeValue.Format(formatStr)
}

// TimeComparison time
// 0 equal,
// 1 time1 greater than time2
// 2 time1 less than time2
func TimeComparison(time1, time2 time.Time) int {
	subSeconds := time1.Sub(time2).Seconds()
	if subSeconds == 0 {
		return 0
	} else if subSeconds > 0 {
		return 1
	} else {
		return 2
	}
}

// timeFormatParse ..
func timeFormatParse(formatStr string) string {
	res := strings.Replace(formatStr, "MM", "01", -1)
	res = strings.Replace(res, "M", "1", -1)
	res = strings.Replace(res, "DD", "02", -1)
	res = strings.Replace(res, "D", "2", -1)
	res = strings.Replace(res, "YYYY", "2006", -1)
	res = strings.Replace(res, "YY", "06", -1)
	res = strings.Replace(res, "HH", "15", -1)
	res = strings.Replace(res, "H", "15", -1)
	res = strings.Replace(res, "hh", "03", -1)
	res = strings.Replace(res, "h", "3", -1)
	res = strings.Replace(res, "mm", "04", -1)
	res = strings.Replace(res, "m", "4", -1)
	res = strings.Replace(res, "ss", "05", -1)
	res = strings.Replace(res, "s", "5", -1)
	return res
}

// gettimeLocation
func gettimeLocation(locationName string) *time.Location {
	var loc *time.Location
	if locationName == "" {
		loc = time.Now().Location()
	} else {
		loc, _ = time.LoadLocation(locationName)
	}
	return loc
}
