/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react";
import "./main.css"
import MainPage from "./MainPage"
import {AppLogo, AppHeader} from "./Common/Logo"

export default class App extends React.Component {
    render() {
        return (
			<div id="app-main">
				<div className="app-header">
					<AppLogo fill="#FFFFFF"/>
					<AppHeader />
				</div>
				<div className="app-body">
					<MainPage />
				</div>
			</div>
        )
    }
}