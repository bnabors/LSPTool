/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {
    SHOW_RESULT,
    REFRESH_RESULT_ALL,
    REFRESH_REQUEST,
    APP_ERROR,
    CLEAR_APP_MSG,
    REFRESH_PING_ALL,
    REFRESH_PING,
    REFRESH_ROUTER_ALL,
    REFRESH_ROUTER_INTERFACE,
    CLEAR_REFRESH_ROUTER_ALL,
    CLEAR_REFRESH_ROUTER_INTERFACE
} from "../constants/actionTypes"

export default function reducer(state = {
    currentResult: null,
    fetching: false,
    fetched: false,
    error: null
}, action) {
    switch (action.type) {
        case SHOW_RESULT: {
            return {...state, fetching: false, fetched: true, currentResult: action.payload}
        }
        case REFRESH_REQUEST: {
            return {...state, fetching: true, fetched: false}
        }
        case APP_ERROR: {
            if (state.fetching) {
                return {...state, fetching: false, fetched: true}
            } else {
                return {...state}
            }
        }
        case CLEAR_APP_MSG: {
            return {...state, fetching: false, fetched: true}
        }
        case REFRESH_RESULT_ALL: {
            let data = action.payload;
            if (data) {
                return {...state, fetching: false, fetched: true, currentResult: data.result}
            } else {
                return {...state, fetching: false, fetched: false, error: "Invalid response"}
            }
        }
        case REFRESH_PING_ALL:
        case REFRESH_PING:
        case REFRESH_ROUTER_ALL:
        case REFRESH_ROUTER_INTERFACE:
        case CLEAR_REFRESH_ROUTER_ALL:
        case CLEAR_REFRESH_ROUTER_INTERFACE:{
            let data = action.payload;
            if (data) {
                return {...state, fetching: false, fetched: true}
            } else {
                return {...state, fetching: false, fetched: false, error: "Invalid response"}
            }
        }
    }

    return state;
}
