(self.webpackChunkpritunl_zero=self.webpackChunkpritunl_zero||[]).push([[98500,46144,70140],{94461:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return function(t){var a=arguments.length>1&&void 0!==arguments[1]?arguments[1]:{},i=a.width,r=i&&e.matchPatterns[i]||e.matchPatterns[e.defaultMatchWidth],n=t.match(r);if(!n)return null;var u,l=n[0],d=i&&e.parsePatterns[i]||e.parsePatterns[e.defaultParseWidth],s=Array.isArray(d)?function(e,t){for(var a=0;a<e.length;a++)if(t(e[a]))return a;return}(d,(function(e){return e.test(l)})):function(e,t){for(var a in e)if(e.hasOwnProperty(a)&&t(e[a]))return a;return}(d,(function(e){return e.test(l)}));return u=e.valueCallback?e.valueCallback(s):s,{value:u=a.valueCallback?a.valueCallback(u):u,rest:t.slice(l.length)}}},e.exports=t.default},44497:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return function(t){var a=arguments.length>1&&void 0!==arguments[1]?arguments[1]:{},i=t.match(e.matchPattern);if(!i)return null;var r=i[0],n=t.match(e.parsePattern);if(!n)return null;var u=e.valueCallback?e.valueCallback(n[0]):n[0];return{value:u=a.valueCallback?a.valueCallback(u):u,rest:t.slice(r.length)}}},e.exports=t.default},26527:(e,t,a)=>{"use strict";var i=a(1654).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var r=i(a(94461)),n=i(a(44497)),u=["רא","שנ","של","רב","ח","שי","שב","שמ","ת","ע"],l={ordinalNumber:(0,n.default)({matchPattern:/^(\d+|(ראשון|שני|שלישי|רביעי|חמישי|שישי|שביעי|שמיני|תשיעי|עשירי|ראשונה|שנייה|שלישית|רביעית|חמישית|שישית|שביעית|שמינית|תשיעית|עשירית))/i,parsePattern:/^(\d+|רא|שנ|של|רב|ח|שי|שב|שמ|ת|ע)/i,valueCallback:function(e){var t=parseInt(e,10);return isNaN(t)?u.indexOf(e)+1:t}}),era:(0,r.default)({matchPatterns:{narrow:/^ל(ספירה|פנה״ס)/i,abbreviated:/^ל(ספירה|פנה״ס)/i,wide:/^ל(פני ה)?ספירה/i},defaultMatchWidth:"wide",parsePatterns:{any:[/^לפ/i,/^לס/i]},defaultParseWidth:"any"}),quarter:(0,r.default)({matchPatterns:{narrow:/^[1234]/i,abbreviated:/^q[1234]/i,wide:/^רבעון [1234]/i},defaultMatchWidth:"wide",parsePatterns:{any:[/1/i,/2/i,/3/i,/4/i]},defaultParseWidth:"any",valueCallback:function(e){return e+1}}),month:(0,r.default)({matchPatterns:{narrow:/^\d+/i,abbreviated:/^(ינו|פבר|מרץ|אפר|מאי|יוני|יולי|אוג|ספט|אוק|נוב|דצמ)׳?/i,wide:/^(ינואר|פברואר|מרץ|אפריל|מאי|יוני|יולי|אוגוסט|ספטמבר|אוקטובר|נובמבר|דצמבר)/i},defaultMatchWidth:"wide",parsePatterns:{narrow:[/^1$/i,/^2/i,/^3/i,/^4/i,/^5/i,/^6/i,/^7/i,/^8/i,/^9/i,/^10/i,/^11/i,/^12/i],any:[/^ינ/i,/^פ/i,/^מר/i,/^אפ/i,/^מא/i,/^יונ/i,/^יול/i,/^אוג/i,/^ס/i,/^אוק/i,/^נ/i,/^ד/i]},defaultParseWidth:"any"}),day:(0,r.default)({matchPatterns:{narrow:/^[אבגדהוש]׳/i,short:/^[אבגדהוש]׳/i,abbreviated:/^(שבת|יום (א|ב|ג|ד|ה|ו)׳)/i,wide:/^יום (ראשון|שני|שלישי|רביעי|חמישי|שישי|שבת)/i},defaultMatchWidth:"wide",parsePatterns:{abbreviated:[/א׳$/i,/ב׳$/i,/ג׳$/i,/ד׳$/i,/ה׳$/i,/ו׳$/i,/^ש/i],wide:[/ן$/i,/ני$/i,/לישי$/i,/עי$/i,/מישי$/i,/שישי$/i,/ת$/i],any:[/^א/i,/^ב/i,/^ג/i,/^ד/i,/^ה/i,/^ו/i,/^ש/i]},defaultParseWidth:"any"}),dayPeriod:(0,r.default)({matchPatterns:{any:/^(אחר ה|ב)?(חצות|צהריים|בוקר|ערב|לילה|אחה״צ|לפנה״צ)/i},defaultMatchWidth:"any",parsePatterns:{any:{am:/^לפ/i,pm:/^אחה/i,midnight:/^ח/i,noon:/^צ/i,morning:/בוקר/i,afternoon:/בצ|אחר/i,evening:/ערב/i,night:/לילה/i}},defaultParseWidth:"any"})};t.default=l,e.exports=t.default},1654:e=>{e.exports=function(e){return e&&e.__esModule?e:{default:e}},e.exports.__esModule=!0,e.exports.default=e.exports}}]);
//# sourceMappingURL=date-fns-locale-he-_lib-match-index-js.js.map