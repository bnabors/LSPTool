/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {
    FETCH_ROUTERS,
    FETCH_ROUTERS_SUCCESS,
    APP_ERROR,
    CLEAR_APP_MSG
} from "../constants/actionTypes"

export default function reducer(state = {
    ingress_routers: [],
    egress_routers: [],
    fetching: false,
    fetched: false,
    error: null
}, action) {
    switch (action.type) {
        case FETCH_ROUTERS: {
            return {...state, fetching: true}
        }
        case FETCH_ROUTERS_SUCCESS: {

            let ingress_routers = action.payload.ingress_routers;
            if (ingress_routers) {
                ingress_routers.sort(function (a, b) {
                    return a.name.localeCompare(b.name);
                })
            }

            let egress_routers = action.payload.egress_routers;
            if (egress_routers) {
                egress_routers.sort(function (a, b) {
                    return a.name.localeCompare(b.name);
                })
            }

            return {
                ...state,
                fetching: false,
                fetched: true,
                ingress_routers: ingress_routers,
                egress_routers: egress_routers
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