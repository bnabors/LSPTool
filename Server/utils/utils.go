/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package utils

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/WOWLABS/LSPTool/Server/config"
	"github.com/WOWLABS/LSPTool/Server/log"
	"github.com/WOWLABS/LSPTool/Server/models"
	_ "github.com/go-sql-driver/mysql"
)

func ReadFromFile(fileName string) []byte {
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		lspLogger.Errorln(err.Error())
	}

	return bytes
}

func WriteToFile(fileName string, bytes []byte) {
	file, err := os.Create(fileName)

	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.Write(bytes)
	writer.Flush()
}

func WriteToFileAsXML(fileName string, data interface{}) {
	xmlOut, err := xml.MarshalIndent(data, "", "    ")

	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	WriteToFile(fileName, xmlOut)
}

func ConvertToJson(data interface{}) (result []byte, err error) {
	result, err = json.Marshal(&data)

	if err != nil {
		lspLogger.Error(err)
		return
	}
	lspLogger.Debugln(string(result))

	return
}

// GetPhysicalName returns physical name of interface by logical name
func GetPhysicalName(logicalName string) string {
	var lastIndex = strings.LastIndex(logicalName, ".")

	if lastIndex == -1 {
		return logicalName
	}

	return logicalName[:lastIndex]
}

/*main page*/
func GetRoutersFromMysql() (result []*models.Router, err error) {
	db, err := sql.Open("mysql", config.LspConfig.MysqlConnectionsString)
	if err != nil {
		lspLogger.Error(err)
		return
	}

	defer db.Close()

	rows, err := db.Query(config.LspConfig.MysqlRouterQuery)
	if err != nil {
		lspLogger.Error(err)
		return
	}

	result, err = ParseDbAnswer(rows)
	if err != nil {
		lspLogger.Error(err)
	}
	return
}

func GetRoutersFromMysqlByID(routersIdList []string) (result []*models.Router, err error) {
	db, err := sql.Open("mysql", config.LspConfig.MysqlConnectionsString)
	if err != nil {
		lspLogger.Error(err)
		return
	}

	defer db.Close()

	var routersList = ""

	for _, routerId := range routersIdList {
		if routersList != "" {
			routersList += ", "
		}

		routersList += routerId
	}

	rows, err := db.Query(config.LspConfig.MysqlRouterQuery + " WHERE element_id IN (" + routersList + ")")
	if err != nil {
		lspLogger.Error(err)
		return
	}

	result, err = ParseDbAnswer(rows)
	if err != nil {
		lspLogger.Error(err)
	}
	return
}

func GetRouterFromMysqlByID(idRouter string) (*models.Router, error) {
	db, err := sql.Open("mysql", config.LspConfig.MysqlConnectionsString)
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(config.LspConfig.MysqlRouterQuery + " WHERE element_id=" + idRouter)
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}

	result, err := ParseDbAnswer(rows)
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}
	if result == nil || len(result) < 1 {
		return nil, errors.New("Router does not exists")
	}

	return result[0], nil
}

func MapRouterCol(colname string, router *models.Router) interface{} {
	switch colname {
	case "element_id":
		return &router.Id
	case "hostname":
		return &router.Name
	case "loopbackip":
		return &router.Ip
	case "proxyIp4":
		return &router.ProxyIp
	default:
		panic("unknown column " + colname)
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func ParseDbAnswer(rows *sql.Rows) ([]*models.Router, error) {
	// get the column names from the query
	var columns []string
	columns, err := rows.Columns()
	panicOnErr(err)

	colNum := len(columns)
	result := []*models.Router{}

	for rows.Next() {
		router := models.Router{}

		// make references for the cols with the aid of VehicleCol
		cols := make([]interface{}, colNum)
		for i := 0; i < colNum; i++ {
			cols[i] = MapRouterCol(columns[i], &router)
		}

		err = rows.Scan(cols...)
		panicOnErr(err)

		result = append(result, &router)
	}
	return result, nil
}
