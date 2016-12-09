/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

'use strict';
import {parseResponse} from "./common";
import {
    showResult,
    refreshAll,
    refreshRouter,
    refreshInterface,
    clearRefreshRouter,
    clearRefreshInterface,
    refreshPings,
    refreshPing
} from "./currentResultActions";
import {sendError, sendInfo, clearMessage, sendCopy} from "./errorActions";
import {determineLsps} from "./lspsActions";
import {runTests, clearResults} from "./resultActions";
import {fetchRouters} from "./routerActions";
import {getSate, setSate} from "./stateActions";
import {openPrintDialog, closePrintDialog, printContent} from "./printActions";


module.exports = {
    parseResponse: parseResponse,
    showResult: showResult,
    refreshAll: refreshAll,
    refreshRouter: refreshRouter,
    refreshInterface: refreshInterface,
    clearRefreshRouter: clearRefreshRouter,
    clearRefreshInterface: clearRefreshInterface,
    refreshPings: refreshPings,
    refreshPing: refreshPing,
    sendError: sendError,
    sendInfo: sendInfo,
    sendCopy: sendCopy,
    clearMessage: clearMessage,
    determineLsps: determineLsps,
    runTests: runTests,
    clearResults: clearResults,
    fetchRouters: fetchRouters,
    getSate: getSate,
    setSate: setSate,
    openPrintDialog: openPrintDialog,
    closePrintDialog: closePrintDialog,
    printContent: printContent
};
