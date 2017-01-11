/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./svg.css"
import React from "react";

const DIRECTION = { RIGHT: 0, LEFT: 1, DOWN: 2 };
const SHIFT_TEXT = 37;
const SHIFT_TEXT_CENTER = 10;
const RADIUS = 25;
const LENGTH = 300;
const MAX_STEP = 3;

const INITIAL_POSITION_X = 200;
const INITIAL_POSITION_X_SHORT = 100;

const INITIAL_POSITION_Y = 70;
const INITIAL_POSITION_LABEL_X = 10;
const INITIAL_POSITION_LABEL_Y = INITIAL_POSITION_Y;

const ADDITIONAL_X_LABEL_RIGHT = 175;
const ADDITIONAL_X_LABEL_RIGHT_SHORT = 75;

const SUMMARY_STEP = 18;
const SMALL_SHIFT_TEXT = 7;

const DrawLine = function(option) {
    let sx = option.sx;
    let sy = option.sy;
    let ex = option.ex;
    let ey = option.ey;
    switch (option.direction) {
        case DIRECTION.DOWN:
            sy += RADIUS;
            ey -= RADIUS;
            break;
        case DIRECTION.RIGHT:
            sx += RADIUS;
            ex -= RADIUS;
            break;
        case DIRECTION.LEFT:
            sx -= RADIUS;
            ex += RADIUS;
            break;
    }
    return (<line x1={sx + "px"} y1={sy + "px"} x2={ex + "px"} y2={ey + "px"} className="line-style"/>)
};

const DrawCircle = function(x, y) {
    return (<circle cx={x + "px"} cy={y + "px"} r={RADIUS + "px"} className="circle-style"/>)
};

const DrawCircleText = function(x, y, label) {
    return (<text x={x + "px"} y={y + 9 + "px"} className="circle-text-style text-center">{label}</text>)
};

const DrawTextCenter = function(option, name, className) {
    let x = (option.ex + option.sx) / 2;
    let y = (option.ey + option.sy) / 2;
    let baseClassName = "text-center";
    if (option.direction == DIRECTION.DOWN) {
        baseClassName = "text-left";
        x += SHIFT_TEXT_CENTER;
        y += SHIFT_TEXT_CENTER * 3;
    } else {
        y -= SHIFT_TEXT_CENTER;
    }


    let finishClassName = className ? baseClassName + " " + className : baseClassName;
    return (<text className={finishClassName} x={x + "px"} y={y + "px"}>{name}</text>)
};

const DrawText = function(x, y, className, text) {
    return (<text className={className} x={x + "px"} y={y + "px"}>{text}</text>)
};

const DrawRect = function(x, y, width, height, className){
    return <rect x={x} y={y} width={width} height={height} className={className} />
};

const GetMax = function(max, test1, test2) {
    return Math.max(Math.max(test1, test2), max);
};

const NumberToLocal = function(stringValue) {
    let number = parseFloat(stringValue);
    if (number) {
        return number.toLocaleString("en-US", {maximumFractionDigits: 0});
    }
    return stringValue;
};

export default class RouteDiagram extends React.Component {

    static maxx = 0;
    static maxy = 0;
    static minx = 0;
    static miny = 0;

    constructor(props){
        super(props);

        RouteDiagram.maxx = 0;
        RouteDiagram.maxy = 0;
        RouteDiagram.minx = 0;
        RouteDiagram.miny = 0;
    }

    static DrawRouterNames(options, text1, text2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (text2) {
                    res.push(DrawText(options.ex, options.ey - SHIFT_TEXT, "text-center text-bold", text2));
                }
                if (text1) {
                    res.push(DrawText(options.sx, options.sy - SHIFT_TEXT, "text-center text-bold", text1));
                }
                break;
            case (DIRECTION.RIGHT):
                if (text1) {
                    res.push(DrawText(options.sx, options.sy - SHIFT_TEXT, "text-center text-bold", text1));
                }
                if (text2) {
                    res.push(DrawText(options.ex, options.ey - SHIFT_TEXT, "text-center text-bold", text2));
                }
                break;
            case (DIRECTION.DOWN):
                if (text1) {
                    res.push(DrawText(options.sx, options.sy - SHIFT_TEXT, "text-center text-bold", text1));
                }
                if (text2) {
                    res.push(DrawText(options.ex, options.ey - SHIFT_TEXT, "text-center text-bold", text2));
                }
                break;
        }
        return res;
    }

    static DrawIps(options, text1, text2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (text1) {
                    res.push(DrawText(options.sx - SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-right text-decorator-green", text1));
                }
                if (text2) {
                    res.push(DrawText(options.ex + SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-left text-decorator-green", text2));
                }
                break;
            case (DIRECTION.RIGHT):
                if (text1) {
                    res.push(DrawText(options.sx + SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-left text-decorator-green", text1));
                }
                if (text2) {
                    res.push(DrawText(options.ex - SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-right text-decorator-green", text2));
                }
                break;
            case (DIRECTION.DOWN):
                if (options.nextDirection === DIRECTION.LEFT) {
                    if (text1) {
                        res.push(DrawText(options.sx + SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-left text-decorator-green", text1));
                    }
                    if (text2) {
                        res.push(DrawText(options.ex + SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-left text-decorator-green", text2));
                    }
                } else {
                    if (text1) {
                        res.push(DrawText(options.sx - SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-right text-decorator-green", text1));
                    }
                    if (text2) {
                        res.push(DrawText(options.ex - SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-right text-decorator-green", text2));
                    }
                }
                break;
        }
        return res;
    }

    static DrawInterfaces(options, text1, text2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (text1) {
                    res.push(DrawText(options.sx - SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-right", text1));
                }
                if (text2) {
                    res.push(DrawText(options.ex + SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-left", text2));
                }
                break;
            case (DIRECTION.RIGHT):
                if (text1) {
                    res.push(DrawText(options.sx + SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-left", text1));
                }
                if (text2) {
                    res.push(DrawText(options.ex - SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-right", text2));
                }
                break;
            case (DIRECTION.DOWN):
                if (options.nextDirection === DIRECTION.LEFT) {
                    if (text1) {
                        res.push(DrawText(options.sx + SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-left", text1));
                    }
                    if (text2) {
                        res.push(DrawText(options.ex + SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-left", text2));
                    }
                } else {
                    if (text1) {
                        res.push(DrawText(options.sx - SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-right", text1));
                    }
                    if (text2) {
                        res.push(DrawText(options.ex - SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-right", text2));
                    }
                }
                break;
        }
        return res;
    }

    static DrawRouterSummary(x, y, router, className, drawUp) {
        let res = [];

        if (drawUp === true) {
            res.push(DrawText(x, y, className, NumberToLocal(router.outputp)));

            y -= SUMMARY_STEP;
            res.push(DrawText(x, y, className, NumberToLocal(router.inputp)));

            y -= SUMMARY_STEP;
            res.push(DrawText(x, y, className, NumberToLocal(router.outputb)));

            y -= SUMMARY_STEP;
            res.push(DrawText(x, y, className, NumberToLocal(router.inputb)));

        } else {
            res.push(DrawText(x, y, className, NumberToLocal(router.inputb)));

            y += SUMMARY_STEP;
            res.push(DrawText(x, y, className, NumberToLocal(router.outputb)));

            y += SUMMARY_STEP;
            res.push(DrawText(x, y, className, NumberToLocal(router.inputp)));

            y += SUMMARY_STEP;
            res.push(DrawText(x, y, className, NumberToLocal(router.outputp)));
        }
        RouteDiagram.checkPoint(x, y);
        return res;
    }

    static DrawSummary(options, router1, router2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (router1) {
                    res.push(RouteDiagram.DrawRouterSummary(options.sx - SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-right", false));
                }
                if (router2) {
                    res.push(RouteDiagram.DrawRouterSummary(options.ex + SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-left", false));
                }
                break;
            case (DIRECTION.RIGHT):
                if (router1) {
                    res.push(RouteDiagram.DrawRouterSummary(options.sx + SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-left", false));
                }
                if (router2) {
                    res.push(RouteDiagram.DrawRouterSummary(options.ex - SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-right", false));
                }
                break;
            case (DIRECTION.DOWN):
                if (options.nextDirection === DIRECTION.LEFT) {
                    if (router1) {
                        res.push(RouteDiagram.DrawRouterSummary(options.sx + SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-left", false));
                    }
                    if (router2) {
                        res.push(RouteDiagram.DrawRouterSummary(options.ex + SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-left", false));
                    }
                } else {
                    if (router1) {
                        res.push(RouteDiagram.DrawRouterSummary(options.sx - SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-right", false));
                    }
                    if (router2) {
                        res.push(RouteDiagram.DrawRouterSummary(options.ex - SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-right", false));
                    }
                }
                break;
        }
        return res;
    }

    static CheckAndDrawSummaryLabels(options) {
        if (options.direction !== DIRECTION.DOWN) {
            return [];
        }
        options.labelYs.push(options.ey);
        return RouteDiagram.DrawSummaryLabels(INITIAL_POSITION_LABEL_X, options.ey, false);
    }

    static DrawSummaryLabels(x, y, isRight = false, isShort = false) {
        y += SHIFT_TEXT + SMALL_SHIFT_TEXT;
        let res = [];

        let txtClass = "text-left text-label";
        if (isRight === true) {
            txtClass = "text-right text-label";
            x += isShort === true ? ADDITIONAL_X_LABEL_RIGHT_SHORT : ADDITIONAL_X_LABEL_RIGHT;
        }
        res.push(DrawText(x, y, txtClass, "BPS In"));

        y += SUMMARY_STEP;
        res.push(DrawText(x, y, txtClass, "BPS Out"));

        y += SUMMARY_STEP;
        res.push(DrawText(x, y, txtClass, "PPS In"));

        y += SUMMARY_STEP;
        res.push(DrawText(x, y, txtClass, "PPS Out"));

        RouteDiagram.checkPoint(x, y);

        return res
    }

    static DrawSummaryLabelRectangles(x, y, isShort = false) {
        let startX = INITIAL_POSITION_LABEL_X - 2;
        y += SHIFT_TEXT + SMALL_SHIFT_TEXT;
        let startY = y - SMALL_SHIFT_TEXT - 5;

        let res = [];

        x += isShort === true ? ADDITIONAL_X_LABEL_RIGHT_SHORT: ADDITIONAL_X_LABEL_RIGHT;
        let width = x;
        let height = SUMMARY_STEP;
        res.push(DrawRect(startX, startY, width, height, "rect-fill"));

        y += SUMMARY_STEP + SUMMARY_STEP;
        startY = y - SMALL_SHIFT_TEXT - 5;
        res.push(DrawRect(startX, startY, width, height, "rect-fill"));

        return res
    }

    static getNextOptions(options, lastId, nextId) {
        if (lastId === nextId) {
            options.render = false;
            return options;
        }

        options.render = true;

        if (options.step >= MAX_STEP) {
            let direction = null;
            let nextDirection = options.nextDirection;
            options.step = 1;

            switch (options.direction) {
                case DIRECTION.RIGHT:
                    direction = DIRECTION.DOWN;
                    nextDirection = DIRECTION.LEFT;
                    options.step++;
                    break;
                case DIRECTION.LEFT:
                    direction = DIRECTION.DOWN;
                    nextDirection = DIRECTION.RIGHT;
                    options.step++;
                    break;
                case DIRECTION.DOWN:
                    direction = options.nextDirection;
                    nextDirection = options.nextDirection;
                    break;
            }

            options.direction = direction;
            options.nextDirection = nextDirection;
        }
        options.sx = options.ex;
        options.sy = options.ey;
        if (options.direction === DIRECTION.RIGHT) {
            options.ex += LENGTH;
        } else if (options.direction === DIRECTION.LEFT) {
            options.ex -= LENGTH;
        } else if (options.direction === DIRECTION.DOWN) {
            options.ey += LENGTH * 3 / 4;
        }
        options.step++;
        return options;
    }

    static checkPoint(x, y) {
        if (x < 0) {
            RouteDiagram.minx = Math.min(RouteDiagram.minx, x);
        } else {
            RouteDiagram.maxx = Math.max(RouteDiagram.maxx, x);
        }
        if (y < 0) {
            RouteDiagram.miny = Math.min(RouteDiagram.miny , y);
        } else {
            RouteDiagram.maxy = Math.max(RouteDiagram.maxy, y);
        }
    }

    static generate(routes) {
        let initialX = routes.length < 6 ? INITIAL_POSITION_X_SHORT : INITIAL_POSITION_X;

        let options = {
            sx: initialX,
            ex: initialX,
            sy: INITIAL_POSITION_Y,
            ey: INITIAL_POSITION_Y,
            direction: DIRECTION.RIGHT,
            nextDirection: DIRECTION.RIGHT,
            step: 0,
            render: false,
            labelYs: []
        };
        RouteDiagram.maxx = options.sx;
        RouteDiagram.maxy = options.sy;

        let lastRouterId = undefined;
        let rectangles = [];
        let lines = [];
        let texts = [];
        let circles = [];

        texts.push(RouteDiagram.DrawSummaryLabels(INITIAL_POSITION_LABEL_X, INITIAL_POSITION_LABEL_Y, false));
        options.labelYs.push(INITIAL_POSITION_LABEL_Y);

        for (let i = 0; i < routes.length; i++) {
            let route = routes[i];

            options = RouteDiagram.getNextOptions(options, lastRouterId, route.router1.id);
            RouteDiagram.maxx = GetMax(RouteDiagram.maxx, options.sx, options.ex);
            RouteDiagram.maxy = GetMax(RouteDiagram.maxy, options.sy, options.ey);

            if (options.render === true) {
                lines.push(DrawLine(options));

                circles.push(DrawCircle(options.sx, options.sy));
                circles.push(DrawCircle(options.ex, options.ey));

                texts.push(DrawTextCenter(options, route.baseIp, "text-decorator-green"));
                texts.push(DrawCircleText(options.sx, options.sy, i == 0 ? "I": "T"));
                texts.push(DrawCircleText(options.ex, options.ey, i + 1 == routes.length ? "E": "T"));
                texts.push(RouteDiagram.DrawRouterNames(options, route.router1.name, route.router2.name));
                texts.push(RouteDiagram.DrawIps(options, route.router1.ip, route.router2.ip));
                texts.push(RouteDiagram.DrawInterfaces(options, route.router1.interface, route.router2.interface));
                texts.push(RouteDiagram.DrawSummary(options, route.router1, route.router2));

                texts.push(RouteDiagram.CheckAndDrawSummaryLabels(options));

                options.step++;
            } else {
                options = RouteDiagram.getNextOptions(options, route.router1.id, route.router2.id);
                RouteDiagram.maxx = GetMax(RouteDiagram.maxx, options.sx, options.ex);
                RouteDiagram.maxy = GetMax(RouteDiagram.maxy, options.sy, options.ey);

                if (options.render === true) {
                    lines.push(DrawLine(options));

                    circles.push(DrawCircle(options.ex, options.ey));

                    texts.push(DrawTextCenter(options, route.baseIp, "text-decorator-green"));
                    texts.push(DrawCircleText(options.ex, options.ey, i + 1 == routes.length ? "E": "T"));
                    texts.push(RouteDiagram.DrawRouterNames(options, undefined, route.router2.name));
                    texts.push(RouteDiagram.DrawIps(options, route.router1.ip, route.router2.ip));
                    texts.push(RouteDiagram.DrawInterfaces(options, route.router1.interface, route.router2.interface));
                    texts.push(RouteDiagram.DrawSummary(options, route.router1, route.router2));
                    texts.push(RouteDiagram.CheckAndDrawSummaryLabels(options));
                }
            }

            lastRouterId = route.router2.id;
        }

        let maxx = RouteDiagram.maxx;
        for (let i = 0; i < options.labelYs.length; i++) {
            texts.push(RouteDiagram.DrawSummaryLabels(maxx, options.labelYs[i], true, routes.length < 3));

            let rects = RouteDiagram.DrawSummaryLabelRectangles(maxx, options.labelYs[i], routes.length < 3);
            for (let j = 0; j < rects.length; j++) {
                rectangles.push(rects[j]);
            }
        }

        let ex = RouteDiagram.maxx + INITIAL_POSITION_LABEL_X;
        let ey = RouteDiagram.maxy + 30;

        let svgStyle = {
            minWidth: (ex / 2) + 'px',
            maxWidth: ex + 'px',
            minHeight: (ey / 2) + 'px',
            maxHeight: ey + 'px',
            height: ey + 'px',
            width: ex + 'px',
        };

        return (
            <svg className="dia" style={svgStyle} x="0px" y="0px" width={ex + "px"} height={ey + "px"}>
                {rectangles}
                {lines}
                {circles}
                {texts}
            </svg>
        )
    }

    render() {
        return RouteDiagram.generate(this.props.data.paths);
    }
}
