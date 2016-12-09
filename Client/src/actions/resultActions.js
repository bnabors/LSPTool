/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {RUN_TEST, RUN_TEST_SUCCESS, RUN_TEST_CLEAR, GENERATE_STATE} from "../constants/actionTypes";
import {request} from "./common";

export function runTests(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/runTests',
            data: data,
            preAction: RUN_TEST,
            preData: data,
            afterActions: [GENERATE_STATE, RUN_TEST_SUCCESS]
        }));
    }
}

export function clearResults() {
    return function (dispatch) {
        dispatch({type: RUN_TEST_CLEAR, payload: null});
    }
}