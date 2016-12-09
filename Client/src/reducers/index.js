/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {combineReducers} from "redux";
import errors from "./errorReducer";
import routers from "./routerReducer";
import LSPs from "./lspsReducer";
import results from "./resultReduser";
import currentResult from "./currentResultReducer"
import states from "./stateReduser"
import printing from "./printReducer"

export default combineReducers({
    errors,
    routers,
    LSPs,
    results,
    currentResult,
    states,
    printing,
});