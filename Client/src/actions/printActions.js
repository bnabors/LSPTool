/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {PRINT_RESULT, CLEAR_PRINT_RESULT, PRINT_CONTENT} from "../constants/actionTypes";

export function openPrintDialog(data) {
    return function (dispatch) {
        dispatch({type: PRINT_RESULT, payload: data});
    }
}

export function printContent(data) {
    return function (dispatch) {
        dispatch({type: PRINT_CONTENT, payload: data});
    }
}

export function closePrintDialog() {
    return function (dispatch) {
        dispatch({type: CLEAR_PRINT_RESULT, payload: null});
    }
}