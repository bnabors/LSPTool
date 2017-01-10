/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./main.css"
import React from "react";
import {FilterableSelect} from "../Common/Select"

export default class RouterList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            ingress: null,
            egress: null
        }
    }

    componentWillReceiveProps(newProps) {
        //init selected ingress
        let ingressId = null;
        if (this.state.ingress === null) {
            if (newProps.ingress_routers !== this.props.ingress_routers) {
                if (newProps.ingress_routers.length > 0) {
                    /*
                     let router = newProps.ingress_routers[0];
                     ingressId = router.id;
                     this.setState({ingress: router});
                     */
                    this.setState({ingress: null});
                }
            }
        } else {
            ingressId = this.state.ingress.id
        }

        //init selected egress
        if (this.state.egress === null) {
            if (newProps.egress_routers !== this.props.egress_routers) {
                if (newProps.egress_routers.length > 0) {
                    /*
                     if (newProps.egress_routers[0].id != ingressId) {
                     this.setState({egress: newProps.egress_routers [0]});
                     } else if (newProps.egress_routers.length > 1) {
                     this.setState({egress: newProps.egress_routers [1]});
                     }
                     */
                    this.setState({egress: null});
                }
            }
        }
    }

    selectIngressHandle(event) {
        //let router = this.getRouterById(this.props.ingress_routers, event.target.value);

        let router = event ? this.getRouterById(this.props.ingress_routers, event.value) : null;
        this.setState({ingress: router});

        if (this.state.egress !== null && router !== null && this.state.egress.id === router.id) {
            this.setState({egress: null});
            this.refs.egress.resetSelection();
        }
    }

    selectEgressHandle(event) {
        //let router = this.getRouterById(this.props.egress_routers, event.target.value);
        let router = event ? this.getRouterById(this.props.egress_routers, event.value) : null;
        this.setState({egress: router});

        if (this.state.ingress !== null && router !== null && this.state.ingress.id === router.id) {
            this.setState({ingress: null});
            this.refs.ingress.resetSelection();
        }
    }

    getRouterById(collection, id) {
        if (!id || !collection) {
            return null;
        }

        let result = collection.find(function (item, index) {
            if (item.id === parseInt(id)) {
                return item;
            }
        });

        return result ? result : null;
    }

    determineLSps() {
        let error = undefined;
        if (!this.state.ingress || this.state.ingress == null) {
            error = "Select ingress router";
        }
        if (!this.state.egress || this.state.egress == null) {
            if (error == undefined) {
                error = "Select egress router"
            } else {
                error += "\r\nSelect egress router"
            }
        }
        if (error) {
            this.props.sendMessage(error);
            return
        }

        this.props.determineLsps({
            ingress: this.state.ingress,
            egress: this.state.egress
        });
    }

    render() {
        let ingressRouterNodes = this.props.ingress_routers.map(function (router) {
            return { value: router.id, label: router.name };
            /*(
                <option id={router.id} value={router.id}>
                    {router.name}
                </option>
            )*/
        }, this);

        let egressRouterNodes = this.props.egress_routers.map(function (router) {
            return { value: router.id, label: router.name };
            /*(
                <option id={router.id} value={router.id}>
                    {router.name}
                </option>
            )*/
        }, this);

        return (
            <div className="routerLst">
                <div className="text-overflow router-header" title="Ingress and Egress Routers">Ingress and Egress Routers</div>
                <div className="routers-select-container">
                    <div className="router-row ">
                        <div className="cell">
                            <label className="named margin-bottom10">Ingress</label>
                        </div>
                        <div className="cell w100">
                            <FilterableSelect ref="ingress" id="select-ingress" value={this.state.ingress}
                                              onChange={this.selectIngressHandle.bind(this)}
                                              disabled={!this.props.isEnabled} content={ingressRouterNodes}/>
                        </div>
                    </div>
                    <div className="router-row">
                        <div className="cell">
                            <label className="named">Egress</label>
                        </div>
                        <div className="cell w100">
                            <FilterableSelect ref="egress" id="select-egress" value={this.state.egress}
                                              onChange={this.selectEgressHandle.bind(this)}
                                              disabled={!this.props.isEnabled} content={egressRouterNodes}/>
                        </div>
                    </div>
                </div>
                <div>
                    <button className="app-control-gradient" onClick={this.determineLSps.bind(this)}
                            disabled={!this.props.isEnabled } title="Determine LSPs">Determine LSPs
                    </button>
                </div>
                <div className="separator"></div>
            </div>
        )
    }
}