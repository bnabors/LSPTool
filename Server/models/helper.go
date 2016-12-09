/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

import (
	"strconv"
	"strings"
	"time"
)

func AddValueRow(rows *[]*StatisticsValue, header string, val string) {
	*rows = append(*rows, &StatisticsValue{Header: header, Value: val, IsError: false})
}

func AddErrorRow(rows *[]*StatisticsValue, header string, val string) {
	*rows = append(*rows, &StatisticsValue{Header: header, Value: val, IsError: checkIsError(val)})
}

func AddLastFlappedErrorRow(rows *[]*StatisticsValue, value string) {
	temp := strings.TrimSpace(value)
	index := strings.Index(temp, "(") - 1
	if index < 0 {
		*rows = append(*rows, &StatisticsValue{Header: value, Value: "", IsError: false})
		return
	}
	temp = temp[:index]

	layout := "2006-01-02 15:04:05 MST"
	datetime, err := time.Parse(layout, temp)
	if err != nil {
		*rows = append(*rows, &StatisticsValue{Header: value, Value: "", IsError: false})
		return
	}
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	*rows = append(*rows, &StatisticsValue{Header: value, Value: "", IsError: datetime.After(yesterday)})
}

func checkIsError(countStr string) bool {
	var countStrTrimmed = strings.TrimSpace(countStr)

	if val, err := strconv.Atoi(countStrTrimmed); err == nil {
		return val > 0
	}

	return false
}
