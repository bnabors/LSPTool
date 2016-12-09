/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./main.css"
import React from "react"
import * as ReactDOM from "react/lib/ReactDOM";

export default class CssLoader extends React.Component{

    componentDidMount() {
        let width = this.props.cwidth ? this.props.cwidth : 49;
        let height = this.props.cheight ? this.props.cheight : 49;

        let root = ReactDOM.findDOMNode(this.refs.root);
        if (root) {

            let val = "margin-top: -" + height / 2 + "px;height: " + height + "px;";
            root.setAttribute("style", val);
        }
        let animation = ReactDOM.findDOMNode(this.refs.animation);
        if (animation) {
            let val = "width: " + width + "px;height: " + height + "px;";
            animation.setAttribute("style", val);
        }
    }

    render(){
        return (
            <div className="cssload-container" ref="root">
                <div className="cssload-speeding-wheel" ref="animation"></div>
            </div>
        )
    }
}
