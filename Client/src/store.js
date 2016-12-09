/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import {applyMiddleware, createStore} from "redux";
import logger from "redux-logger";
import promise from "redux-promise-middleware";
import thunk from "redux-thunk";
import reducer from "./reducers/index";

const middleware = applyMiddleware(promise(), thunk, logger());
const store = createStore(reducer, middleware);

export default store;

