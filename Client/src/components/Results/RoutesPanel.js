/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react";
import "./main.css"
import * as ReactDOM from "react/lib/ReactDOM";
//import Loader from "../Common/Loader"
import CssLoader from "../Common/CssLoader"

export default class RoutesPanel extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            selectedId: null
        }
    }

    selectTab(route) {
        this.setState({selectedId: route.id});
        this.props.onSelected(route);
    }

    componentWillReceiveProps(newProps) {
        if (this.state.selectedId != null) {
            return
        }

        if (newProps.data !== this.props.data) {
            if (newProps.data.p2p != null && newProps.data.p2p.length > 0) {
                this.selectTab(newProps.data.p2p[0]);
            } else if (newProps.data.p2mp != null && newProps.data.p2mp.length > 0) {
                this.selectTab(newProps.data.p2mp[0]);
            }
        }
    }

    render() {
        let clsName = "tab margin";
        let content = null;

        if (this.props.data != null && (this.props.data.p2p != null || this.props.data.p2mp != null )) {
            let p2pContent = [];
            if (this.props.data.p2p != null && this.props.data.p2p.length > 0) {
                p2pContent.push(<div className="result-header">P2P LSPs</div>);
                let p2p = [];
                for (let i = 0; i < this.props.data.p2p.length; i++) {
                    let item = this.props.data.p2p[i];
                    p2p.push(<li><RouteItem id={item.id} data={item} onClicked={this.selectTab.bind(this)}
                                            selectedId={this.state.selectedId}/></li>);
                }
                p2pContent.push(<ul className={clsName}>{p2p}</ul>)
            }
            let p2mpContent = [];
            if (this.props.data.p2mp != null && this.props.data.p2mp.length > 0) {
                p2pContent.push(<div className="result-header">P2MP LSPs</div>);
                let p2mp = [];
                for (let i = 0; i < this.props.data.p2mp.length; i++) {
                    let item = this.props.data.p2mp[i];
                    p2mp.push(<li><RouteItem id={item.id} data={item} onClicked={this.selectTab.bind(this)}
                                             selectedId={this.state.selectedId}/></li>);
                }
                p2pContent.push(<ul className={clsName}>{p2mp}</ul>)
            }

            content = (
                <div>
                    <div className="routers-header">Selected LSPs</div>
                    {p2pContent}
                    {p2mpContent}
                </div>
            );
        } else {
            clsName += " height";
            content = (<CssLoader cwidth="49" cheight="49"/>);
        }

        return (<div className={clsName}>{content}</div>)
    }
}

class RouteItem extends React.Component {

    clickHandler() {
        this.props.onClicked(this.props.data);
    }

    componentWillReceiveProps(newProps) {
        this.selectElement(newProps.selectedId);
    }

    selectElement(id) {
        let element = ReactDOM.findDOMNode(this.refs.test);
        if (element) {
            if (this.props.data.id === id) {
                if (element.className.indexOf("active") < 0) {
                    element.className += " active";
                }

            } else {
                element.className = element.className.replace(" active", "");
            }
        }
    }

    componentDidMount() {
        this.selectElement(this.props.selectedId);
    }

    render() {
        return (<a ref="test" className="routelinks border-right" onClick={this.clickHandler.bind(this)}>{this.props.data.name}<div/></a>)
    }
}