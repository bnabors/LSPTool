/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react"
import * as ReactDOM from "react/lib/ReactDOM";

export default class Loader extends React.Component {

    componentDidMount() {
        let element = ReactDOM.findDOMNode(this.refs.animation);
        if (element) {
            element.setAttribute("attributeName", "transform");
            element.setAttribute("type", "rotate");
            element.setAttribute("from", "0 25 25");
            element.setAttribute("to", "360 25 25");
            element.setAttribute("dur", "1s");
            element.setAttribute("repeatCount", "indefinite");
        }

        let root = ReactDOM.findDOMNode(this.refs.root);
        if (root) {
            let width = this.props.cwidth ? this.props.cwidth : 40;
            let height = this.props.cheight ? this.props.cheight : 40;
            let val = "margin-left:-" + width / 2 + "px;margin-top:-" + height / 2 + "px;width:" + width + "px;height:" + height + "px;";
            root.setAttribute("style", val);
        }
    }

    render() {
        let width = this.props.cwidth ? this.props.cwidth : 40;
        let height = this.props.cheight ? this.props.cheight : 40;
        return (
            <div className="loader" ref="root">
                <svg x="0px" y="0px" width={width + "px"} height={height + "px"} viewBox="0 0 50 50">
                    <path
                        d="M43.935,25.145c0-10.318-8.364-18.683-18.683-18.683c-10.318,0-18.683,8.365-18.683,18.683h4.068c0-8.071,6.543-14.615,14.615-14.615c8.072,0,14.615,6.543,14.615,14.615H43.935z">
                        <animateTransform ref="animation" attributeName="transform" type="rotate" from="0 25 25"
                                          to="360 25 25" dur="0.6s" repeatCount="indefinite"/>
                    </path>
                </svg>
            </div>
        )
    }
}


