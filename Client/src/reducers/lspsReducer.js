/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {
    DETERMINE_LSP,
    DETERMINE_LSP_SUCCESS,
    APP_ERROR,
    CLEAR_APP_MSG,
} from "../constants/actionTypes"

export default function reducer(state = {
    routers: null,
    p2pLSPs: [],
    p2mpLSPs: [],
    downLSPs: [],
    lspGroups: [],
    fetching: false,
    fetched: false
}, action) {
    switch (action.type) {
        case DETERMINE_LSP: {
            return {...state, fetching: true, fetched: false, p2pLSPs: [], p2mpLSPs: [], downLSPs: [], lspGroups: [], routers: action.payload}
        }
        case DETERMINE_LSP_SUCCESS: {
            return {
                ...state,
                fetching: false,
                fetched: true,
                p2pLSPs: action.payload.p2p,
                p2mpLSPs: action.payload.p2mp,
                downLSPs: action.payload.down,
                lspGroups: action.payload.lspGroup
            }
        }
        case APP_ERROR: {
            if (state.fetching) {
                return {...state, fetching: false, fetched: false}
            } else {
                return {...state}
            }
        }
        case CLEAR_APP_MSG: {
            return {...state, fetching: false}
        }
    }

    return state;
}
