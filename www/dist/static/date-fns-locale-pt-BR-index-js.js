(self.webpackChunkpritunl_zero=self.webpackChunkpritunl_zero||[]).push([[45649,36108,23534,46144,70140,33135,29006,42468,87236,99818],{86273:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:{},a=t.width?String(t.width):e.defaultWidth;return e.formats[a]||e.formats[e.defaultWidth]}},e.exports=t.default},53347:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return function(t,a){var r;if("formatting"===(null!=a&&a.context?String(a.context):"standalone")&&e.formattingValues){var n=e.defaultFormattingWidth||e.defaultWidth,i=null!=a&&a.width?String(a.width):n;r=e.formattingValues[i]||e.formattingValues[n]}else{var o=e.defaultWidth,d=null!=a&&a.width?String(a.width):e.defaultWidth;r=e.values[d]||e.values[o]}return r[e.argumentCallback?e.argumentCallback(t):t]}},e.exports=t.default},94461:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return function(t){var a=arguments.length>1&&void 0!==arguments[1]?arguments[1]:{},r=a.width,n=r&&e.matchPatterns[r]||e.matchPatterns[e.defaultMatchWidth],i=t.match(n);if(!i)return null;var o,d=i[0],u=r&&e.parsePatterns[r]||e.parsePatterns[e.defaultParseWidth],s=Array.isArray(u)?function(e,t){for(var a=0;a<e.length;a++)if(t(e[a]))return a;return}(u,(function(e){return e.test(d)})):function(e,t){for(var a in e)if(e.hasOwnProperty(a)&&t(e[a]))return a;return}(u,(function(e){return e.test(d)}));return o=e.valueCallback?e.valueCallback(s):s,{value:o=a.valueCallback?a.valueCallback(o):o,rest:t.slice(d.length)}}},e.exports=t.default},44497:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return function(t){var a=arguments.length>1&&void 0!==arguments[1]?arguments[1]:{},r=t.match(e.matchPattern);if(!r)return null;var n=r[0],i=t.match(e.parsePattern);if(!i)return null;var o=e.valueCallback?e.valueCallback(i[0]):i[0];return{value:o=a.valueCallback?a.valueCallback(o):o,rest:t.slice(n.length)}}},e.exports=t.default},8148:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a={lessThanXSeconds:{one:"menos de um segundo",other:"menos de {{count}} segundos"},xSeconds:{one:"1 segundo",other:"{{count}} segundos"},halfAMinute:"meio minuto",lessThanXMinutes:{one:"menos de um minuto",other:"menos de {{count}} minutos"},xMinutes:{one:"1 minuto",other:"{{count}} minutos"},aboutXHours:{one:"cerca de 1 hora",other:"cerca de {{count}} horas"},xHours:{one:"1 hora",other:"{{count}} horas"},xDays:{one:"1 dia",other:"{{count}} dias"},aboutXWeeks:{one:"cerca de 1 semana",other:"cerca de {{count}} semanas"},xWeeks:{one:"1 semana",other:"{{count}} semanas"},aboutXMonths:{one:"cerca de 1 mês",other:"cerca de {{count}} meses"},xMonths:{one:"1 mês",other:"{{count}} meses"},aboutXYears:{one:"cerca de 1 ano",other:"cerca de {{count}} anos"},xYears:{one:"1 ano",other:"{{count}} anos"},overXYears:{one:"mais de 1 ano",other:"mais de {{count}} anos"},almostXYears:{one:"quase 1 ano",other:"quase {{count}} anos"}},r=function(e,t,r){var n,i=a[e];return n="string"==typeof i?i:1===t?i.one:i.other.replace("{{count}}",String(t)),null!=r&&r.addSuffix?r.comparison&&r.comparison>0?"em "+n:"há "+n:n};t.default=r,e.exports=t.default},32677:(e,t,a)=>{"use strict";var r=a(1654).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var n=r(a(86273)),i={date:(0,n.default)({formats:{full:"EEEE, d 'de' MMMM 'de' y",long:"d 'de' MMMM 'de' y",medium:"d MMM y",short:"dd/MM/yyyy"},defaultWidth:"full"}),time:(0,n.default)({formats:{full:"HH:mm:ss zzzz",long:"HH:mm:ss z",medium:"HH:mm:ss",short:"HH:mm"},defaultWidth:"full"}),dateTime:(0,n.default)({formats:{full:"{{date}} 'às' {{time}}",long:"{{date}} 'às' {{time}}",medium:"{{date}}, {{time}}",short:"{{date}}, {{time}}"},defaultWidth:"full"})};t.default=i,e.exports=t.default},5255:(e,t)=>{"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a={lastWeek:function(e){var t=e.getUTCDay();return"'"+(0===t||6===t?"último":"última")+"' eeee 'às' p"},yesterday:"'ontem às' p",today:"'hoje às' p",tomorrow:"'amanhã às' p",nextWeek:"eeee 'às' p",other:"P"},r=function(e,t,r,n){var i=a[e];return"function"==typeof i?i(t):i};t.default=r,e.exports=t.default},87495:(e,t,a)=>{"use strict";var r=a(1654).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var n=r(a(53347)),i={ordinalNumber:function(e,t){var a=Number(e);return"week"===(null==t?void 0:t.unit)?a+"ª":a+"º"},era:(0,n.default)({values:{narrow:["AC","DC"],abbreviated:["AC","DC"],wide:["antes de cristo","depois de cristo"]},defaultWidth:"wide"}),quarter:(0,n.default)({values:{narrow:["1","2","3","4"],abbreviated:["T1","T2","T3","T4"],wide:["1º trimestre","2º trimestre","3º trimestre","4º trimestre"]},defaultWidth:"wide",argumentCallback:function(e){return e-1}}),month:(0,n.default)({values:{narrow:["j","f","m","a","m","j","j","a","s","o","n","d"],abbreviated:["jan","fev","mar","abr","mai","jun","jul","ago","set","out","nov","dez"],wide:["janeiro","fevereiro","março","abril","maio","junho","julho","agosto","setembro","outubro","novembro","dezembro"]},defaultWidth:"wide"}),day:(0,n.default)({values:{narrow:["D","S","T","Q","Q","S","S"],short:["dom","seg","ter","qua","qui","sex","sab"],abbreviated:["domingo","segunda","terça","quarta","quinta","sexta","sábado"],wide:["domingo","segunda-feira","terça-feira","quarta-feira","quinta-feira","sexta-feira","sábado"]},defaultWidth:"wide"}),dayPeriod:(0,n.default)({values:{narrow:{am:"a",pm:"p",midnight:"mn",noon:"md",morning:"manhã",afternoon:"tarde",evening:"tarde",night:"noite"},abbreviated:{am:"AM",pm:"PM",midnight:"meia-noite",noon:"meio-dia",morning:"manhã",afternoon:"tarde",evening:"tarde",night:"noite"},wide:{am:"a.m.",pm:"p.m.",midnight:"meia-noite",noon:"meio-dia",morning:"manhã",afternoon:"tarde",evening:"tarde",night:"noite"}},defaultWidth:"wide",formattingValues:{narrow:{am:"a",pm:"p",midnight:"mn",noon:"md",morning:"da manhã",afternoon:"da tarde",evening:"da tarde",night:"da noite"},abbreviated:{am:"AM",pm:"PM",midnight:"meia-noite",noon:"meio-dia",morning:"da manhã",afternoon:"da tarde",evening:"da tarde",night:"da noite"},wide:{am:"a.m.",pm:"p.m.",midnight:"meia-noite",noon:"meio-dia",morning:"da manhã",afternoon:"da tarde",evening:"da tarde",night:"da noite"}},defaultFormattingWidth:"wide"})};t.default=i,e.exports=t.default},51113:(e,t,a)=>{"use strict";var r=a(1654).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var n=r(a(94461)),i={ordinalNumber:(0,r(a(44497)).default)({matchPattern:/^(\d+)[ºªo]?/i,parsePattern:/\d+/i,valueCallback:function(e){return parseInt(e,10)}}),era:(0,n.default)({matchPatterns:{narrow:/^(ac|dc|a|d)/i,abbreviated:/^(a\.?\s?c\.?|d\.?\s?c\.?)/i,wide:/^(antes de cristo|depois de cristo)/i},defaultMatchWidth:"wide",parsePatterns:{any:[/^ac/i,/^dc/i],wide:[/^antes de cristo/i,/^depois de cristo/i]},defaultParseWidth:"any"}),quarter:(0,n.default)({matchPatterns:{narrow:/^[1234]/i,abbreviated:/^T[1234]/i,wide:/^[1234](º)? trimestre/i},defaultMatchWidth:"wide",parsePatterns:{any:[/1/i,/2/i,/3/i,/4/i]},defaultParseWidth:"any",valueCallback:function(e){return e+1}}),month:(0,n.default)({matchPatterns:{narrow:/^[jfmajsond]/i,abbreviated:/^(jan|fev|mar|abr|mai|jun|jul|ago|set|out|nov|dez)/i,wide:/^(janeiro|fevereiro|março|abril|maio|junho|julho|agosto|setembro|outubro|novembro|dezembro)/i},defaultMatchWidth:"wide",parsePatterns:{narrow:[/^j/i,/^f/i,/^m/i,/^a/i,/^m/i,/^j/i,/^j/i,/^a/i,/^s/i,/^o/i,/^n/i,/^d/i],any:[/^ja/i,/^fev/i,/^mar/i,/^abr/i,/^mai/i,/^jun/i,/^jul/i,/^ago/i,/^set/i,/^out/i,/^nov/i,/^dez/i]},defaultParseWidth:"any"}),day:(0,n.default)({matchPatterns:{narrow:/^(dom|[23456]ª?|s[aá]b)/i,short:/^(dom|[23456]ª?|s[aá]b)/i,abbreviated:/^(dom|seg|ter|qua|qui|sex|s[aá]b)/i,wide:/^(domingo|(segunda|ter[cç]a|quarta|quinta|sexta)([- ]feira)?|s[aá]bado)/i},defaultMatchWidth:"wide",parsePatterns:{short:[/^d/i,/^2/i,/^3/i,/^4/i,/^5/i,/^6/i,/^s[aá]/i],narrow:[/^d/i,/^2/i,/^3/i,/^4/i,/^5/i,/^6/i,/^s[aá]/i],any:[/^d/i,/^seg/i,/^t/i,/^qua/i,/^qui/i,/^sex/i,/^s[aá]b/i]},defaultParseWidth:"any"}),dayPeriod:(0,n.default)({matchPatterns:{narrow:/^(a|p|mn|md|(da) (manhã|tarde|noite))/i,any:/^([ap]\.?\s?m\.?|meia[-\s]noite|meio[-\s]dia|(da) (manhã|tarde|noite))/i},defaultMatchWidth:"any",parsePatterns:{any:{am:/^a/i,pm:/^p/i,midnight:/^mn|^meia[-\s]noite/i,noon:/^md|^meio[-\s]dia/i,morning:/manhã/i,afternoon:/tarde/i,evening:/tarde/i,night:/noite/i}},defaultParseWidth:"any"})};t.default=i,e.exports=t.default},5598:(e,t,a)=>{"use strict";var r=a(1654).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var n=r(a(8148)),i=r(a(32677)),o=r(a(5255)),d=r(a(87495)),u=r(a(51113)),s={code:"pt-BR",formatDistance:n.default,formatLong:i.default,formatRelative:o.default,localize:d.default,match:u.default,options:{weekStartsOn:0,firstWeekContainsDate:1}};t.default=s,e.exports=t.default},1654:e=>{e.exports=function(e){return e&&e.__esModule?e:{default:e}},e.exports.__esModule=!0,e.exports.default=e.exports}}]);
//# sourceMappingURL=date-fns-locale-pt-BR-index-js.js.map