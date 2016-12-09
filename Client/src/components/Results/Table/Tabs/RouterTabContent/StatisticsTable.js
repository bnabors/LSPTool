/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react"

export default class StatisticsTable extends React.Component {
    render() {
        let isError = false;
        let rows = [];
        for (let i = 0; i < this.props.data.values.length; i++) {
            let item = this.props.data.values[i];
            if (item.isError) {
                isError = true;
            }
            let className = item.isError ? " error" : " no-error";
            let headerClassName = className + (i == 0 ? "" : " can-hide");
            rows.push((
                <tr className={headerClassName}>
                    <td>{i == 0 ? this.props.data.title : ""}</td>
                    <td className={className}>{item.header}</td>
                    <td className={className}>{item.value}</td>
                </tr>))
        }
        return (
            <tbody className={isError ? " error" : " no-error"}>
            {rows}
            </tbody>
        )
    }
}
