/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./main.css"
import "../main.css"
import React from "react";
import LspsTab from "./Tabs/LspsTab"
import RouterTab from "./Tabs/RouterTab"
import IcmpTab from "./Tabs/IcmpTab"
import HeaderItem from "./HeaderItem"
import {ResultType} from "../../../constants/resultType"

export default class Table extends React.Component {
    constructor(props){
        super(props);
        this.state={
            selected: null,
        };
    }

    selectTab(tab){
        this.saveState(this.props.states);

        this.setState({selected: tab});
        let state = this.props.states;
        state.selectedTab = tab.id;
        this.saveState(state);
    }

    updateState(){
        this.saveState(this.props.states);
    }

    saveState(state){
        for(let i = 0; i < state.tabs.length; i++){
            if(state.tabs[i].tabId === state.selectedTab){
                let element = this.refs[state.selectedTab];
                if(element) {
                    state.tabs[i] = element.getState();
                }
            }
        }
        this.props.onStateChanged(state);
    }

    routerSelected(routerId){
        let tab = this.props.data.results.find(function (item, index) {
            if(item.id === routerId){
                return item;
            }
        });

        this.selectTab(tab)
    }

    reloadAll(){
        this.props.commands.reloadAll();
    }

    printSummary(){
        this.props.commands.printSummary();
    }

    componentWillReceiveProps(newProps){
        if(newProps.data !== this.props.data){
            if(this.props.states != null){
                this.saveState(this.props.states);
            }
            if(newProps.states.selectedTab != null){
                let tab = newProps.data.results.find(function (item, index) {
                    if(item.id === newProps.states.selectedTab){
                        return item;
                    }
                });
                this.setState({selected: tab});
                let state = newProps.states;
                state.selectedTab = tab.id;
                this.props.onStateChanged(state);
            } else {
                if (newProps.data.results != null && newProps.data.results.length > 0) {
                    this.selectTab(newProps.data.results[0]);
                } else {
                    this.selectTab(null);
                }
            }
        }
    }

    componentDidMount() {
        if(this.props.states.selectedTab != null){
            let findId = this.props.states.selectedTab;
            let tab = this.props.data.results.find(function (item, index) {
                if(item.id === findId){
                    return item;
                }
            });
            this.selectTab(tab);
        } else if (this.props.data != null && this.props.data.results != null && this.props.data.results.length > 0) {
            this.selectTab(this.props.data.results[0]);
        } else {
            this.selectTab(null);
        }
    }

    onStateChanged(id, state){
        let states = this.props.states;
        for(let i=0; i< states.tabs.length; i++){
            if(states.tabs[i].tabId === id){
                states.tabs[i].state = state
            }
        }
        this.props.onStateChanged(states);
    }

    render(){
        if(this.props.data.results == null){
            return (<div>ERROR!</div>)
        }

        let selectedId = this.state.selected ? this.state.selected.id: null;
        let headers = this.props.data.results.map(function(item){
            switch (item.type) {
                case ResultType.LSPs:
                case ResultType.ROUTER:
                case ResultType.ICPM:
                    return (<HeaderItem id={item.id} data={item} onClicked={this.selectTab.bind(this)} selectedId={selectedId} />);
                default:
                    return (<div>ERROR!</div>)
            }
        }, this);

        let tab = null;
        if(this.state.selected){
            let findId = this.state.selected.id;
            let state = this.props.states.tabs.find(function (item, index) {
                if(item.tabId === findId){
                    return item.state
                }
            });

            switch (this.state.selected.type) {
                case ResultType.LSPs:
                    tab = (<LspsTab id={this.state.selected.id} ref={this.state.selected.id} data={this.state.selected} tabState={state}
                                    onClicked={this.routerSelected.bind(this)} onStateChanged={this.onStateChanged.bind(this)}
                                    printLsp={this.props.commands.printLsp}/>);
                    break;
                case ResultType.ROUTER:
                    tab = (<RouterTab id={this.state.selected.id} ref={this.state.selected.id} data={this.state.selected} tabState={state}
                                      onStateChanged={this.onStateChanged.bind(this)} commands={this.props.commands}/>);
                    break;
                case ResultType.ICPM:
                    tab = (<IcmpTab id={this.state.selected.id} ref={this.state.selected.id} data={this.state.selected} tabState={state}
                                    onStateChanged={this.onStateChanged.bind(this)} commands={this.props.commands}/>);
                    break;
                default:
                    tab = (<div>ERROR!</div>);
                    break;
            }
        }
        return (
            <div id="results" className="app-div-result">
                <div className="app-panel app-inner-panel">
                    <div className="right-border"></div>
                    <div className="container-ex">
                        <div className="result-buttons">
                            <div>
                                <button className="app-control-gradient text-overflow" onClick={this.props.commands.reloadAll} title="Reload All">Reload All</button>
                            </div>
                            <div>
                                <button className="app-control-gradient text-overflow" onClick={this.props.commands.printSummary} title="Print Summary">Print Summary</button>
                            </div>
                        </div>
                        <div className="ul-wrapper">
                            <ul className="tab">
                                {headers}
                            </ul>
                        </div>
                    </div>
                </div>
                <div className="result-panel-content">
                    <div className="app-result-panel">
                        {tab}
                    </div>
                </div>
            </div>
        )
    }
}
