/// <reference path="./References.d.ts"/>
import * as MobileDetect from 'mobile-detect';

let md = new MobileDetect(window.navigator.userAgent);

export const mobile = !!md.mobile();
