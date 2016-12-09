/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controller

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/commands"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/models"
)

func DetermineLsps(o models.DetermineOptions) (models.LspCollection, error) {
	address := o.Ingress.GetAddress()

	result := models.LspCollection{P2P: []*models.LspItem{}, P2MP: []*models.LspItem{}, Down: []*models.LspItem{}, LspGroups: []*models.LspGroup{}}
	mplsLspInfos, err := command.LoadMplsLspInfo(address, o.Egress.Ip)
	if err != nil {
		return result, err
	}

	for _, mplsInfo := range mplsLspInfos {
		var path = filterIpAdreses(mplsInfo.MplsLspPath.ReceivedRro)

		var item = models.LspItem{
			Id:             mplsInfo.Name,
			Name:           mplsInfo.Name,
			IngressIp:      o.Ingress.Ip,
			IngressPuttyip: o.Ingress.GetAddress(),
			EgressIp:       o.Egress.Ip,
			EgressPuttyip:  o.Egress.GetAddress(),
			GroupId:        getGroupId(&result.LspGroups, path),
			Bandwidth:      mplsInfo.MplsLspPath.Bandwidth}

		switch strings.ToLower(mplsInfo.LspState) {
		case "up":
			if mplsInfo.IsP2MP() {
				result.P2MP = append(result.P2MP, &item)
			} else {
				result.P2P = append(result.P2P, &item)
			}
		case "dn":
			result.Down = append(result.Down, &item)
		}
	}

	return result, nil
}

func getGroupId(groups *[]*models.LspGroup, path []string) string {
	for _, group := range *groups {
		if len(group.ReceivedRro) != len(path) {
			continue
		}

		var pos int
		for pos = 0; pos < len(path); pos++ {
			if group.ReceivedRro[pos] != path[pos] {
				break
			}
		}

		if pos == len(path) {
			return group.Id
		}
	}

	var newGroup = models.LspGroup{Id: "lspGroup" + strconv.Itoa(len(*groups)), ReceivedRro: path}
	*groups = append(*groups, &newGroup)

	return newGroup.Id
}

func filterIpAdreses(receivedRro string) []string {
	result := make([]string, 0)

	for _, param := range getParams(`(?P<IP>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s*(?P<Info>\(.*?\))?`, receivedRro) {
		var ip = param["IP"]
		var info = param["Info"]

		if ip == "" || strings.Contains(info, "flag=0x2") {
			continue
		}
		result = append(result, ip)
	}

	return result
}

func getParams(regEx string, str string) (paramsMap []map[string]string) {
	var re = regexp.MustCompile(regEx)

	result := make([]map[string]string, 0)

	groupNames := re.SubexpNames()

	for _, match := range re.FindAllStringSubmatch(str, -1) {
		var groupsMap = make(map[string]string)

		for groupNumber, groupName := range groupNames {
			if groupName == "" {
				continue
			}

			groupsMap[groupName] = match[groupNumber]
		}

		result = append(result, groupsMap)
	}

	return result
}
