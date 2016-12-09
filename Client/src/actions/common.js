/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import axios from "axios";
import {sendError} from "./errorActions";

export function parseResponse(response, sendType) {
    return function (dispatch) {
        parse(response,
            function (data) {
                if (Array.isArray(sendType)) {
                    sendType.forEach(function (item, index, array) {
                        dispatch({type: item, payload: data});
                    })
                } else {
                    dispatch({type: sendType, payload: data});
                }
            },
            function (error) {
                dispatch(sendError(error));
            });
    }
}

function parse(response, successAction, errorAction) {
    if (response.data === undefined || response.data === null || response.data === "") {
        errorAction("Empty response");
    } else if (response.data.error !== undefined || response.data.error === null) {
        errorAction(response.data.error);
    } else {
        if (successAction === undefined || typeof(successAction) !== "function") {
            errorAction("Unknown action")
        } else {
            successAction(response.data);
        }
    }
}

export function request(options) {
    /*
     url - url (throw exception if undefined)
     data - post query data (may be undefined)

     preAction - action called before post query
     preActions - collection of actions which called before post query
     preData - data for preAction

     afterAction - action called after post query
     afterActions - collection of actions which called after post query
     */

    return function (dispatch) {
        if (options.preAction) {
            dispatch({type: options.preAction, payload: options.preData ? options.preData : null})
        }
        if (options.preActions) {
            let data = options.preData ? options.preData : null;
            options.preActions.forEach(function (item, index, array) {
                dispatch({type: item, payload: data});
            });
        }
        if (!options.url) {
            throw "Request error - empty url. Options :" + JSON.stringify(options);
        }
        axios.post(options.url, options.data ? JSON.stringify(options.data) : undefined)
            .then((response) => {
                if (options.afterAction) {
                    dispatch(parseResponse(response, options.afterAction));
                }
                if (options.afterActions) {
                    dispatch(parseResponse(response, options.afterActions));
                }
            })
            .catch((err) => {
                dispatch(sendError(err));
            })
    }
}