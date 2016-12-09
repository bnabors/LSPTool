/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {APP_ERROR, CLEAR_APP_MSG, APP_INFO, APP_COPY} from "../constants/actionTypes"

export default function reducer(state = {
    hasError: false,
    error: null,

    hasInfo:false,
    info:null,

    hasCopy:false,
    copy:null
}, action) {
    switch (action.type) {
        case APP_ERROR: {
            return {...state, hasError: true, error: action.payload}
        }
        case CLEAR_APP_MSG: {
            return {...state, hasError: false, error: null, hasInfo: false, info: null, hasCopy: false, copy: null}
        }
        case APP_INFO:{
            return {...state, hasInfo: true, info: action.payload}
        }
        case APP_COPY:{
            return {...state, hasCopy: true, copy: action.payload}
        }
    }

    return state;
}
