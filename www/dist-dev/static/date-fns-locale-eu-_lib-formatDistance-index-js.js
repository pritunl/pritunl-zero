"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-eu-_lib-formatDistance-index-js"],{

/***/ "./node_modules/date-fns/locale/eu/_lib/formatDistance/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/eu/_lib/formatDistance/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatDistanceLocale = {\n  lessThanXSeconds: {\n    one: 'segundo bat baino gutxiago',\n    other: '{{count}} segundo baino gutxiago'\n  },\n  xSeconds: {\n    one: '1 segundo',\n    other: '{{count}} segundo'\n  },\n  halfAMinute: 'minutu erdi',\n  lessThanXMinutes: {\n    one: 'minutu bat baino gutxiago',\n    other: '{{count}} minutu baino gutxiago'\n  },\n  xMinutes: {\n    one: '1 minutu',\n    other: '{{count}} minutu'\n  },\n  aboutXHours: {\n    one: '1 ordu gutxi gorabehera',\n    other: '{{count}} ordu gutxi gorabehera'\n  },\n  xHours: {\n    one: '1 ordu',\n    other: '{{count}} ordu'\n  },\n  xDays: {\n    one: '1 egun',\n    other: '{{count}} egun'\n  },\n  aboutXWeeks: {\n    one: 'aste 1 inguru',\n    other: '{{count}} aste inguru'\n  },\n  xWeeks: {\n    one: '1 aste',\n    other: '{{count}} astean'\n  },\n  aboutXMonths: {\n    one: '1 hilabete gutxi gorabehera',\n    other: '{{count}} hilabete gutxi gorabehera'\n  },\n  xMonths: {\n    one: '1 hilabete',\n    other: '{{count}} hilabete'\n  },\n  aboutXYears: {\n    one: '1 urte gutxi gorabehera',\n    other: '{{count}} urte gutxi gorabehera'\n  },\n  xYears: {\n    one: '1 urte',\n    other: '{{count}} urte'\n  },\n  overXYears: {\n    one: '1 urte baino gehiago',\n    other: '{{count}} urte baino gehiago'\n  },\n  almostXYears: {\n    one: 'ia 1 urte',\n    other: 'ia {{count}} urte'\n  }\n};\nvar formatDistance = function formatDistance(token, count, options) {\n  var result;\n  var tokenValue = formatDistanceLocale[token];\n  if (typeof tokenValue === 'string') {\n    result = tokenValue;\n  } else if (count === 1) {\n    result = tokenValue.one;\n  } else {\n    result = tokenValue.other.replace('{{count}}', String(count));\n  }\n  if (options !== null && options !== void 0 && options.addSuffix) {\n    if (options.comparison && options.comparison > 0) {\n      return 'en ' + result;\n    } else {\n      return 'duela ' + result;\n    }\n  }\n  return result;\n};\nvar _default = formatDistance;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2V1L19saWIvZm9ybWF0RGlzdGFuY2UvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQSxjQUFjLFFBQVE7QUFDdEIsR0FBRztBQUNIO0FBQ0E7QUFDQSxjQUFjLFFBQVE7QUFDdEIsR0FBRztBQUNIO0FBQ0E7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGlCQUFpQixRQUFRO0FBQ3pCO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsSUFBSTtBQUNKO0FBQ0EsSUFBSTtBQUNKLHlDQUF5QyxPQUFPO0FBQ2hEO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsTUFBTTtBQUNOO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLGtCQUFlO0FBQ2YiLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9wcml0dW5sLXplcm8vLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2V1L19saWIvZm9ybWF0RGlzdGFuY2UvaW5kZXguanM/Yzk4ZCJdLCJzb3VyY2VzQ29udGVudCI6WyJcInVzZSBzdHJpY3RcIjtcblxuT2JqZWN0LmRlZmluZVByb3BlcnR5KGV4cG9ydHMsIFwiX19lc01vZHVsZVwiLCB7XG4gIHZhbHVlOiB0cnVlXG59KTtcbmV4cG9ydHMuZGVmYXVsdCA9IHZvaWQgMDtcbnZhciBmb3JtYXREaXN0YW5jZUxvY2FsZSA9IHtcbiAgbGVzc1RoYW5YU2Vjb25kczoge1xuICAgIG9uZTogJ3NlZ3VuZG8gYmF0IGJhaW5vIGd1dHhpYWdvJyxcbiAgICBvdGhlcjogJ3t7Y291bnR9fSBzZWd1bmRvIGJhaW5vIGd1dHhpYWdvJ1xuICB9LFxuICB4U2Vjb25kczoge1xuICAgIG9uZTogJzEgc2VndW5kbycsXG4gICAgb3RoZXI6ICd7e2NvdW50fX0gc2VndW5kbydcbiAgfSxcbiAgaGFsZkFNaW51dGU6ICdtaW51dHUgZXJkaScsXG4gIGxlc3NUaGFuWE1pbnV0ZXM6IHtcbiAgICBvbmU6ICdtaW51dHUgYmF0IGJhaW5vIGd1dHhpYWdvJyxcbiAgICBvdGhlcjogJ3t7Y291bnR9fSBtaW51dHUgYmFpbm8gZ3V0eGlhZ28nXG4gIH0sXG4gIHhNaW51dGVzOiB7XG4gICAgb25lOiAnMSBtaW51dHUnLFxuICAgIG90aGVyOiAne3tjb3VudH19IG1pbnV0dSdcbiAgfSxcbiAgYWJvdXRYSG91cnM6IHtcbiAgICBvbmU6ICcxIG9yZHUgZ3V0eGkgZ29yYWJlaGVyYScsXG4gICAgb3RoZXI6ICd7e2NvdW50fX0gb3JkdSBndXR4aSBnb3JhYmVoZXJhJ1xuICB9LFxuICB4SG91cnM6IHtcbiAgICBvbmU6ICcxIG9yZHUnLFxuICAgIG90aGVyOiAne3tjb3VudH19IG9yZHUnXG4gIH0sXG4gIHhEYXlzOiB7XG4gICAgb25lOiAnMSBlZ3VuJyxcbiAgICBvdGhlcjogJ3t7Y291bnR9fSBlZ3VuJ1xuICB9LFxuICBhYm91dFhXZWVrczoge1xuICAgIG9uZTogJ2FzdGUgMSBpbmd1cnUnLFxuICAgIG90aGVyOiAne3tjb3VudH19IGFzdGUgaW5ndXJ1J1xuICB9LFxuICB4V2Vla3M6IHtcbiAgICBvbmU6ICcxIGFzdGUnLFxuICAgIG90aGVyOiAne3tjb3VudH19IGFzdGVhbidcbiAgfSxcbiAgYWJvdXRYTW9udGhzOiB7XG4gICAgb25lOiAnMSBoaWxhYmV0ZSBndXR4aSBnb3JhYmVoZXJhJyxcbiAgICBvdGhlcjogJ3t7Y291bnR9fSBoaWxhYmV0ZSBndXR4aSBnb3JhYmVoZXJhJ1xuICB9LFxuICB4TW9udGhzOiB7XG4gICAgb25lOiAnMSBoaWxhYmV0ZScsXG4gICAgb3RoZXI6ICd7e2NvdW50fX0gaGlsYWJldGUnXG4gIH0sXG4gIGFib3V0WFllYXJzOiB7XG4gICAgb25lOiAnMSB1cnRlIGd1dHhpIGdvcmFiZWhlcmEnLFxuICAgIG90aGVyOiAne3tjb3VudH19IHVydGUgZ3V0eGkgZ29yYWJlaGVyYSdcbiAgfSxcbiAgeFllYXJzOiB7XG4gICAgb25lOiAnMSB1cnRlJyxcbiAgICBvdGhlcjogJ3t7Y291bnR9fSB1cnRlJ1xuICB9LFxuICBvdmVyWFllYXJzOiB7XG4gICAgb25lOiAnMSB1cnRlIGJhaW5vIGdlaGlhZ28nLFxuICAgIG90aGVyOiAne3tjb3VudH19IHVydGUgYmFpbm8gZ2VoaWFnbydcbiAgfSxcbiAgYWxtb3N0WFllYXJzOiB7XG4gICAgb25lOiAnaWEgMSB1cnRlJyxcbiAgICBvdGhlcjogJ2lhIHt7Y291bnR9fSB1cnRlJ1xuICB9XG59O1xudmFyIGZvcm1hdERpc3RhbmNlID0gZnVuY3Rpb24gZm9ybWF0RGlzdGFuY2UodG9rZW4sIGNvdW50LCBvcHRpb25zKSB7XG4gIHZhciByZXN1bHQ7XG4gIHZhciB0b2tlblZhbHVlID0gZm9ybWF0RGlzdGFuY2VMb2NhbGVbdG9rZW5dO1xuICBpZiAodHlwZW9mIHRva2VuVmFsdWUgPT09ICdzdHJpbmcnKSB7XG4gICAgcmVzdWx0ID0gdG9rZW5WYWx1ZTtcbiAgfSBlbHNlIGlmIChjb3VudCA9PT0gMSkge1xuICAgIHJlc3VsdCA9IHRva2VuVmFsdWUub25lO1xuICB9IGVsc2Uge1xuICAgIHJlc3VsdCA9IHRva2VuVmFsdWUub3RoZXIucmVwbGFjZSgne3tjb3VudH19JywgU3RyaW5nKGNvdW50KSk7XG4gIH1cbiAgaWYgKG9wdGlvbnMgIT09IG51bGwgJiYgb3B0aW9ucyAhPT0gdm9pZCAwICYmIG9wdGlvbnMuYWRkU3VmZml4KSB7XG4gICAgaWYgKG9wdGlvbnMuY29tcGFyaXNvbiAmJiBvcHRpb25zLmNvbXBhcmlzb24gPiAwKSB7XG4gICAgICByZXR1cm4gJ2VuICcgKyByZXN1bHQ7XG4gICAgfSBlbHNlIHtcbiAgICAgIHJldHVybiAnZHVlbGEgJyArIHJlc3VsdDtcbiAgICB9XG4gIH1cbiAgcmV0dXJuIHJlc3VsdDtcbn07XG52YXIgX2RlZmF1bHQgPSBmb3JtYXREaXN0YW5jZTtcbmV4cG9ydHMuZGVmYXVsdCA9IF9kZWZhdWx0O1xubW9kdWxlLmV4cG9ydHMgPSBleHBvcnRzLmRlZmF1bHQ7Il0sIm5hbWVzIjpbXSwic291cmNlUm9vdCI6IiJ9\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/eu/_lib/formatDistance/index.js\n");

/***/ })

}]);