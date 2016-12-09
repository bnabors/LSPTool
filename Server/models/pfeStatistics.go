/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

// PfeStatistic данные для интрефейса PFE
type PfeStatistic struct {
	Timeout          string
	TruncatedKey     string
	BitsToTest       string
	DataError        string
	StackUnderflow   string
	StackOverflow    string
	ExtendedDiscard  string
	InvalidInterface string
	InfoCellDrops    string
	FabricDrops      string
	OutputMTU        string
}

// ToRouterStatisticsContent создает модель для веба
func (stat PfeStatistic) ToRouterStatisticsContent(router *Router) RouterStatisticsContent {
	rows := []*StatisticsValue{}
	AddErrorRow(&rows, "Timeout", stat.Timeout)
	AddErrorRow(&rows, "Truncated key", stat.TruncatedKey)
	AddErrorRow(&rows, "Bits to test", stat.BitsToTest)
	AddErrorRow(&rows, "Data error", stat.DataError)
	AddErrorRow(&rows, "Stack underflow", stat.StackUnderflow)
	AddErrorRow(&rows, "Stack overflow", stat.StackOverflow)
	AddErrorRow(&rows, "Extended discard", stat.ExtendedDiscard)
	AddErrorRow(&rows, "Invalid interface", stat.InvalidInterface)
	AddErrorRow(&rows, "Info cell drops", stat.InfoCellDrops)
	AddErrorRow(&rows, "Fabric drops", stat.FabricDrops)
	AddErrorRow(&rows, "Output MTU", stat.OutputMTU)
	statistic := Statistics{Title: "PFE Statistics", Values: rows}

	return RouterStatisticsContent{
		Id:            "pfe",
		Name:          "PFE",
		Statistics:    []*Statistics{&statistic},
		SubInterfaces: []*RouterStatisticsContent{}}
}
