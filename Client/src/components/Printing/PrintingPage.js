/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react";
import {closePrintDialog, sendError, refreshPings} from "../../actions/bundleActions"
import * as ReactDOM from "react/lib/ReactDOM";
import PrintingContent from "./Content/PrintContent"
import PrintingOneContent from "./Content/PrintingOneContent"
import {connect} from "react-redux";
//import {ResultType} from "../../constants/resultType";

@connect((store) => {
    return {
        printing: store.printing,
        LSPs: store.LSPs,
        results: store.results,
        routers: store.routers,
    }
})

export default class PrintingPage extends React.Component {
    openWindow(url, name, props) {
        if (/*@cc_on!@*/false) { //do this only in IE
            let windowRef = window.open("", name, props);
            windowRef.close();
        }
        let windowRef = window.open(url, name, props);
        if (!windowRef) {
            return windowRef;
        }

        if (!windowRef.opener) {
            windowRef.opener = self;
        }
        windowRef.focus();
        return windowRef;
    }

    print() {
        let element = ReactDOM.findDOMNode(this.refs.printingContent);
        if (element) {
            try {
                //let win = this.openWindow('/print', 'Print Summary');
                let win = this.openWindow('', this.getTitle(), 'scrollbars=yes,resizable=yes,menubar=yes,location=yes', false);
                if (win) {
                    win.document.open();
                    let header= this.getTitle();
                    win.document.write('<html><head><title>'+header+'</title>' +
                        '<link rel="stylesheet" type="text/css" href="//fonts.googleapis.com/css?family=Open+Sans" />' +
                        '<link rel="stylesheet" type="text/css" href="/src/app.bundle.css"></head><body>');
                    win.document.write(element.innerHTML);
                    win.document.write('</body></html>');
                    win.document.close();
                    /*
                     win.document.onload(function () {
                     let entry = document.getElementById("print-entry");
                     if(entry){
                     entry.appendChild(element);
                     }
                     });
                     */

                    win.focus();
                } else {
                    throw "Can't open printing tab."
                }
            } catch (ex) {
                this.props.dispatch(sendError(ex));
            }
        }
        this.props.dispatch(closePrintDialog());
    }

    getTitle() {
        if(this.props.printing.data !== null){
            return this.props.LSPs.routers.ingress.name + " - " + this.props.LSPs.routers.egress.name + " - " + new Date().toLocaleString();
            //return this.props.printing.data.lsp.name + Date.now().toString();
        }
        if(this.props.printing.content !== null) {
            return this.props.LSPs.routers.ingress.name + " - " + this.props.LSPs.routers.egress.name + " - " + new Date().toLocaleString();
            //return this.props.printing.content.name + Date.now().toString();
        }

        return Date.now().toString();
    }

    isLoaded(data) {
        /*
        let icmp = data.results.find(function (item) {
            if (item.type == ResultType.ICPM) {
                return item;
            }
        });

        if (icmp) {
            if (icmp.content == null || icmp.content.length < 1) {
                return false;
            }
        }
        */

        return true;
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.props.printing.data != null) {
            if (this.isLoaded(this.props.printing.data)) {
                this.print();
            }
        } else if (this.props.printing.content != null) {
            if (this.isLoaded(this.props.printing.content)) {
                this.print();
            }
        }
    }

    refreshPings(options) {
        let data = {
            lsp: this.props.printing.data.lsp,
            groupRouters: this.props.results.testResults.groupRouters,
            options: JSON.stringify(options).toString()
        };
        this.props.dispatch(refreshPings(data));
    }

    render() {
        let content = null;
        let className = "printing-dialog hidden";
        if (this.props.printing.data !== null) {
            content = (<PrintingContent routers={this.props.LSPs.routers} data={this.props.printing.data}
                                        onRefreshPings={this.refreshPings.bind(this)}/>);
        } else if (this.props.printing.content !== null) {
            content = (<PrintingOneContent routers={this.props.LSPs.routers} content={this.props.printing.content}
                                           onRefreshPings={this.refreshPings.bind(this)}/>);
        }

        return (
            <div className={className}>
                <form ref="printingContent" className="print-content">
                    {content}
                </form>
            </div>
        )
    }
}
