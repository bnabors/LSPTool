/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {APP_INFO, APP_ERROR, APP_COPY, CLEAR_APP_MSG} from "../constants/actionTypes"

export function sendError(error) {
    return function (dispatch) {
        dispatch({type: APP_ERROR, payload: error});
    }
}

export function sendInfo(message) {
    return function (dispatch) {
        dispatch({type: APP_INFO, payload: message});
    }
}

export function clearMessage() {
    return function (dispatch) {
        dispatch({type: CLEAR_APP_MSG, payload: null});
    }
}
export function sendCopy(massage) {
    return function (dispatch) {
        dispatch({type: APP_COPY, payload: massage});
    }
}