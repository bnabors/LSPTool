/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {
    RUN_TEST,
    RUN_TEST_SUCCESS,
    RUN_TEST_CLEAR,
    REFRESH_RESULT_ALL,
    REFRESH_ROUTER_ALL,
    REFRESH_ROUTER_INTERFACE,
    CLEAR_REFRESH_ROUTER_ALL,
    CLEAR_REFRESH_ROUTER_INTERFACE,
    REFRESH_PING_ALL,
    REFRESH_PING,
    UPDATE_RESULT_STATES,
    APP_ERROR,
    CLEAR_APP_MSG
} from "../constants/actionTypes"

import {ResultType} from "../constants/resultType"

export default function reducer(state = {
    testOptions: null,
    testResults: null,
    fetching: false,
    fetched: false,
    error: null
}, action) {
    switch (action.type) {
        case RUN_TEST: {
            return {...state, fetching: true, testOptions: action.payload, testResults: null}
        }
        case RUN_TEST_SUCCESS: {
            return {...state, fetching: false, fetched: true, testResults: action.payload}
        }
        case RUN_TEST_CLEAR: {
            return {...state, fetching: false, testOptions: null, testResults: null, fetched: false, error: null}
        }
        case REFRESH_RESULT_ALL: {
            return refresh(action, state, refreshLsp);
        }
        case REFRESH_ROUTER_ALL:
        case CLEAR_REFRESH_ROUTER_ALL:
        case REFRESH_PING_ALL: {
            return refresh(action, state, refreshResult);
        }
        case REFRESH_ROUTER_INTERFACE:
        case CLEAR_REFRESH_ROUTER_INTERFACE: {
            return refresh(action, state, refreshRouterInterface);
        }
        case REFRESH_PING: {
            return refresh(action, state, refreshPing);
        }
        case APP_ERROR: {
            if (state.fetching) {
                return {...state, fetching: false, fetched: true}
            } else {
                return {...state}
            }
        }
        case CLEAR_APP_MSG: {
            return {...state, fetching: false}
        }
        case UPDATE_RESULT_STATES: {
            return {...state, fetching: false, error: null}
        }
    }

    return state;
}

function refresh(action, state, refreshFunc) {
    let data = action.payload;
    let p2p = state.testResults.p2p;
    let p2mp = state.testResults.p2mp;
    let groupRouters = state.testResults.groupRouters;

    if (!refreshFunc(p2p, data)) {
        if (!refreshFunc(p2mp, data)) {
            return {...state}
        } else {
            detectLspErrors(p2mp);
        }
    } else {
        detectLspErrors(p2p);
    }
    return {...state, fetching: false, fetched: true, testResults: {p2p: p2p, p2mp: p2mp, groupRouters: groupRouters}}
}

function detectLspErrors(collection) {
    for (let i = 0; i < collection.length; i++) {
        collection[i].isError = detectResultsErrors(collection[i].results);
    }
}
function detectResultsErrors(collection) {
    let isError = false;
    for (let i = 0; i < collection.length; i++) {
        collection[i].isError = detectResultErrors(collection[i]);
        isError = isError || collection[i].isError
    }
    return isError;
}
function detectResultErrors(result) {
    switch (result.type) {
        case ResultType.LSPs: {
            return detectLspResultErrors(result.content);
        }
        case ResultType.ICPM: {
            return detectIcmpResultErrors(result.content);
        }
        case ResultType.ROUTER: {
            return detectRouterResultErrors(result.content);
        }
        default: {
            return false;
        }
    }
}
function detectLspResultErrors(content) {
    return false;
}
function detectIcmpResultErrors(content) {
    if(content !== null){
        for(let i = 0; i < content.length; i ++){
            if(content[i].isError){
                return true;
            }
        }
    }
    return false;
}
function detectRouterResultErrors(content) {
    if(content !== null){
        let isError = false;
        for (let i = 0; i < content.length; i++) {
            isError = findError(content[i]) || isError;
        }
        return isError
    }
    return false;
}

function refreshLsp(collection, data) {
    if (!collection || collection.length == 0) {
        return false;
    }

    for (let i = 0; i < collection.length; i++) {
        if (collection[i].lsp.id === data.lsp.id) {
            collection[i] = data.result;
            return true;
        }
    }
    return false;
}

function refreshResult(collection, data) {
    if (!collection || collection.length == 0) {
        return false;
    }

    let lsp = findLsp(collection, data.lsp);
    if (!lsp) {
        return false;
    }

    for (let i = 0; i < lsp.results.length; i++) {
        if (lsp.results[i].id === data.result.id) {
            lsp.results[i].content = data.result.content;
            return true;
        }
    }

    return false;
}

function refreshRouterInterface(collection, data) {
    if (!collection || collection.length == 0) {
        return false;
    }

    let lsp = findLsp(collection, data.lsp);
    if (!lsp) {
        return false;
    }

    let result = findInCollection(lsp.results, function (item) {
        return item.id === data.result.id
    });
    if (!result) {
        return false;
    }

    let resultUpdate = false;
    for (let i = 0; i < result.content.length; i++) {
        if (searchInterface(result.content[i], data.result.interface, data.result.result)) {
            resultUpdate = true;
            break;
        }
    }

    result.isError = false;
    if (resultUpdate) {
        for (let i = 0; i < result.content.length; i++) {
            result.isError = findError(result.content[i]) || result.isError;
        }
    }

    return resultUpdate;
}

function refreshPing(collection, data) {
    if (!collection || collection.length == 0) {
        return false;
    }

    let lsp = findLsp(collection, data.lsp);
    if (!lsp) {
        return false;
    }

    let result = findInCollection(lsp.results, function (item) {
        return item.id === data.result.id
    });
    if (!result) {
        return false;
    }

    let dataResult = data.result.result;

    for (let i = 0; i < result.content.length; i++) {
        if (result.content[i].id === dataResult.id) {
            result.content[i] = dataResult;
            return true;
        }
    }
}

function findLsp(collection, lsp) {
    return findInCollection(collection, function (item) {
        return item.lsp.id === lsp.id
    });
}

function findInCollection(collection, predicate) {
    return collection.find(function (item, index) {
        if (predicate(item)) {
            return item;
        }
    });
}

function searchInterface(content, routerInterface, result) {
    if (content.id !== routerInterface.interface) {
        return false;
    }

    if (routerInterface.innerInterface === null) {
        content.name = result.name;
        content.statistics = result.statistics;
        content.sub_interfaces = result.sub_interfaces;
        content.isError = result.isError;
        return true;
    }

    for (let i = 0; i < content.sub_interfaces.length; i++) {
        if (searchInterface(content.sub_interfaces[i], routerInterface.innerInterface, result)) { // <<< recursion
            return true;
        }
    }

    return false;
}

function findError(content) {
    content.isError = findErrorStatistics(content.statistics);

    if (content.sub_interfaces !== undefined || content.sub_interfaces != null) {
        for (let i = 0; i < content.sub_interfaces.length; i++) {
            content.isError = findError(content.sub_interfaces[i]) || content.isError; // <<< recursion
        }
    }
    return content.isError;
}

function findErrorStatistics(statistics) {
    if (statistics === undefined || statistics === null) {
        return false;
    }
    for (let i = 0; i < statistics.length; i++) {
        if (findErrorValue(statistics[i].values)) {
            return true;
        }
    }
    return false;
}

function findErrorValue(values) {
    if (values === undefined || values === null) {
        return false;
    }
    for (let i = 0; i < values.length; i++) {
        if (values[i].isError) {
            return true;
        }
    }
    return false;
}