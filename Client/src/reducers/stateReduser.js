/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {GET_STATE, SET_STATE, RESET_STATE, GENERATE_STATE} from "../constants/actionTypes"

/*
 states: []
 {
 {
 lspId: 0,
 selectedTab: 0,
 tabs:[]
 {
 tabId: 0,
 state: {}
 }
 }
 }

 activeState:
 {
 lspId: 0,
 selectedTab: 0,
 tabStates:[]
 {
 tabId: 0,
 state: {}
 }
 }
 */

export default function reducer(state = {
    activeState: null,
    states: []
}, action) {
    switch (action.type) {
        case GET_STATE: {
            return {...state, activeState: findState(action.payload, state.states)}
        }
        case SET_STATE: {
            return {...state, activeState: action.payload, states: updateState(action.payload, state.states)}
        }
        case RESET_STATE: {
            return {...state, activeState: null, states: []}
        }
        case GENERATE_STATE: {
            return {...state, activeState: null, states: generateStates(action.payload)}
        }
    }
    return state;
}

function generateStates(results) {
    let states = [];
    states = generateByCollection(states, results.p2p);
    states = generateByCollection(states, results.p2mp);
    return states
}

function generateByCollection(states, collection) {
    collection.forEach(function (l_element, l_index, l_array) {
        let tabs = [];
        l_element.results.forEach(function (r_element, r_index, r_array) {
            tabs.push({tabId: r_element.id, state: {}})
        });
        let selectedTab = null;
        if (l_element.results.length > 0) {
            selectedTab = l_element.results[0].id;
        }
        states.push({lspId: l_element.lsp.id, selectedTab: selectedTab, tabs: tabs})
    });
    return states;
}

function findState(options, states) {
    if (options === null) {
        return null;
    }

    let state = states.find(function (item, index) {
        if (item.lspId === options.lsp.id) {
            return item;
        }
    });

    return state ? state : null;
}

function updateState(state, states) {
    if (state == null) {
        return states;
    }

    for (let i = 0; i < states.length; i++) {
        if (states[i].lspId == state.lspId) {
            states[i] = state;
            break;
        }
    }
    return states;
}