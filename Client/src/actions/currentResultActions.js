/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {
    REFRESH_REQUEST,
    SHOW_RESULT,
    REFRESH_RESULT_ALL,
    REFRESH_ROUTER_ALL,
    REFRESH_ROUTER_INTERFACE,
    CLEAR_REFRESH_ROUTER_ALL,
    CLEAR_REFRESH_ROUTER_INTERFACE,
    REFRESH_PING_ALL,
    REFRESH_PING,
    GET_STATE

} from "../constants/actionTypes";
import {request} from "./common";


export function showResult(data) {
    return function (dispatch) {
        dispatch({type: SHOW_RESULT, payload: data});
        dispatch({type: GET_STATE, payload: data});
    }
}

export function refreshAll(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/refreshLsp',
            data: data,
            preAction: REFRESH_REQUEST,
            preData: data.lsp,
            afterAction: REFRESH_RESULT_ALL
        }));
    }
}

export function refreshRouter(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/refreshRouterInfo',
            data: data,
            preAction: REFRESH_REQUEST,
            afterAction: REFRESH_ROUTER_ALL
        }));
    }
}

export function refreshInterface(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/refreshRouterInterface',
            data: data,
            preAction: REFRESH_REQUEST,
            afterAction: REFRESH_ROUTER_INTERFACE
        }));
    }
}

export function clearRefreshRouter(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/clearRefreshRouterInfo',
            data: data,
            preAction: REFRESH_REQUEST,
            afterAction: CLEAR_REFRESH_ROUTER_ALL
        }));
    }
}

export function clearRefreshInterface(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/clearRefreshRouterInterface',
            data: data,
            preAction: REFRESH_REQUEST,
            afterAction: CLEAR_REFRESH_ROUTER_INTERFACE
        }));
    }
}

export function refreshPings(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/refreshPings',
            data: data,
            preAction: REFRESH_REQUEST,
            afterAction: REFRESH_PING_ALL
        }));
    }
}

export function refreshPing(data) {
    return function (dispatch) {
        dispatch(request({
            url: '/refreshPing',
            data: data,
            preAction: REFRESH_REQUEST,
            afterAction: REFRESH_PING
        }));
    }
}


