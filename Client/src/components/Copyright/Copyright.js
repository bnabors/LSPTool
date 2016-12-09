/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./copyright.css";
import React from "react";
import {sendCopy} from "../../actions/bundleActions"
import {connect} from "react-redux";

const CopyrightContent = (<div className="copyright-content">
    <p>Copyright 2016 Juniper Networks, Inc. All rights reserved.</p>
    <p>Licensed under the Juniper Networks Script Software License (the "License").<br/>You may not use this script file
        except in compliance with the License, which is located at</p>
        <p><a href="http://www.juniper.net/support/legal/scriptlicense/">http://www.juniper.net/support/legal/scriptlicense/</a></p>
        <p>Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under
        the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
        implied.</p>
</div>);

@connect((store) => {
    return {}
})

export default class Copyright extends React.Component {
    onCopyright() {
        this.props.dispatch(sendCopy(CopyrightContent));
    }

    render() {
        return (
            <div className="copyright-panel">
                <a onClick={this.onCopyright.bind(this)}>Copyright 2016 Juniper Networks, Inc. All rights reserved.</a>
            </div>
        )
    }
}
