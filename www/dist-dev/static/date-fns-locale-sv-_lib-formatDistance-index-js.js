"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-sv-_lib-formatDistance-index-js"],{

/***/ "./node_modules/date-fns/locale/sv/_lib/formatDistance/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/sv/_lib/formatDistance/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatDistanceLocale = {\n  lessThanXSeconds: {\n    one: 'mindre än en sekund',\n    other: 'mindre än {{count}} sekunder'\n  },\n  xSeconds: {\n    one: 'en sekund',\n    other: '{{count}} sekunder'\n  },\n  halfAMinute: 'en halv minut',\n  lessThanXMinutes: {\n    one: 'mindre än en minut',\n    other: 'mindre än {{count}} minuter'\n  },\n  xMinutes: {\n    one: 'en minut',\n    other: '{{count}} minuter'\n  },\n  aboutXHours: {\n    one: 'ungefär en timme',\n    other: 'ungefär {{count}} timmar'\n  },\n  xHours: {\n    one: 'en timme',\n    other: '{{count}} timmar'\n  },\n  xDays: {\n    one: 'en dag',\n    other: '{{count}} dagar'\n  },\n  aboutXWeeks: {\n    one: 'ungefär en vecka',\n    other: 'ungefär {{count}} vecka'\n  },\n  xWeeks: {\n    one: 'en vecka',\n    other: '{{count}} vecka'\n  },\n  aboutXMonths: {\n    one: 'ungefär en månad',\n    other: 'ungefär {{count}} månader'\n  },\n  xMonths: {\n    one: 'en månad',\n    other: '{{count}} månader'\n  },\n  aboutXYears: {\n    one: 'ungefär ett år',\n    other: 'ungefär {{count}} år'\n  },\n  xYears: {\n    one: 'ett år',\n    other: '{{count}} år'\n  },\n  overXYears: {\n    one: 'över ett år',\n    other: 'över {{count}} år'\n  },\n  almostXYears: {\n    one: 'nästan ett år',\n    other: 'nästan {{count}} år'\n  }\n};\nvar wordMapping = ['noll', 'en', 'två', 'tre', 'fyra', 'fem', 'sex', 'sju', 'åtta', 'nio', 'tio', 'elva', 'tolv'];\nvar formatDistance = function formatDistance(token, count, options) {\n  var result;\n  var tokenValue = formatDistanceLocale[token];\n  if (typeof tokenValue === 'string') {\n    result = tokenValue;\n  } else if (count === 1) {\n    result = tokenValue.one;\n  } else {\n    if (options && options.onlyNumeric) {\n      result = tokenValue.other.replace('{{count}}', String(count));\n    } else {\n      result = tokenValue.other.replace('{{count}}', count < 13 ? wordMapping[count] : String(count));\n    }\n  }\n  if (options !== null && options !== void 0 && options.addSuffix) {\n    if (options.comparison && options.comparison > 0) {\n      return 'om ' + result;\n    } else {\n      return result + ' sedan';\n    }\n  }\n  return result;\n};\nvar _default = formatDistance;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL3N2L19saWIvZm9ybWF0RGlzdGFuY2UvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQSx3QkFBd0IsUUFBUTtBQUNoQyxHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBO0FBQ0Esd0JBQXdCLFFBQVE7QUFDaEMsR0FBRztBQUNIO0FBQ0E7QUFDQSxjQUFjLFFBQVE7QUFDdEIsR0FBRztBQUNIO0FBQ0E7QUFDQSxzQkFBc0IsUUFBUTtBQUM5QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLHNCQUFzQixRQUFRO0FBQzlCLEdBQUc7QUFDSDtBQUNBO0FBQ0EsY0FBYyxRQUFRO0FBQ3RCLEdBQUc7QUFDSDtBQUNBO0FBQ0Esc0JBQXNCLFFBQVE7QUFDOUIsR0FBRztBQUNIO0FBQ0E7QUFDQSxjQUFjLFFBQVE7QUFDdEIsR0FBRztBQUNIO0FBQ0E7QUFDQSxzQkFBc0IsUUFBUTtBQUM5QixHQUFHO0FBQ0g7QUFDQTtBQUNBLGNBQWMsUUFBUTtBQUN0QixHQUFHO0FBQ0g7QUFDQTtBQUNBLG1CQUFtQixRQUFRO0FBQzNCLEdBQUc7QUFDSDtBQUNBO0FBQ0EscUJBQXFCLFFBQVE7QUFDN0I7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLElBQUk7QUFDSjtBQUNBLElBQUk7QUFDSjtBQUNBLDJDQUEyQyxPQUFPO0FBQ2xELE1BQU07QUFDTiwyQ0FBMkMsT0FBTztBQUNsRDtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsTUFBTTtBQUNOO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLGtCQUFlO0FBQ2YiLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9wcml0dW5sLXplcm8vLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL3N2L19saWIvZm9ybWF0RGlzdGFuY2UvaW5kZXguanM/ZTRkMSJdLCJzb3VyY2VzQ29udGVudCI6WyJcInVzZSBzdHJpY3RcIjtcblxuT2JqZWN0LmRlZmluZVByb3BlcnR5KGV4cG9ydHMsIFwiX19lc01vZHVsZVwiLCB7XG4gIHZhbHVlOiB0cnVlXG59KTtcbmV4cG9ydHMuZGVmYXVsdCA9IHZvaWQgMDtcbnZhciBmb3JtYXREaXN0YW5jZUxvY2FsZSA9IHtcbiAgbGVzc1RoYW5YU2Vjb25kczoge1xuICAgIG9uZTogJ21pbmRyZSDDpG4gZW4gc2VrdW5kJyxcbiAgICBvdGhlcjogJ21pbmRyZSDDpG4ge3tjb3VudH19IHNla3VuZGVyJ1xuICB9LFxuICB4U2Vjb25kczoge1xuICAgIG9uZTogJ2VuIHNla3VuZCcsXG4gICAgb3RoZXI6ICd7e2NvdW50fX0gc2VrdW5kZXInXG4gIH0sXG4gIGhhbGZBTWludXRlOiAnZW4gaGFsdiBtaW51dCcsXG4gIGxlc3NUaGFuWE1pbnV0ZXM6IHtcbiAgICBvbmU6ICdtaW5kcmUgw6RuIGVuIG1pbnV0JyxcbiAgICBvdGhlcjogJ21pbmRyZSDDpG4ge3tjb3VudH19IG1pbnV0ZXInXG4gIH0sXG4gIHhNaW51dGVzOiB7XG4gICAgb25lOiAnZW4gbWludXQnLFxuICAgIG90aGVyOiAne3tjb3VudH19IG1pbnV0ZXInXG4gIH0sXG4gIGFib3V0WEhvdXJzOiB7XG4gICAgb25lOiAndW5nZWbDpHIgZW4gdGltbWUnLFxuICAgIG90aGVyOiAndW5nZWbDpHIge3tjb3VudH19IHRpbW1hcidcbiAgfSxcbiAgeEhvdXJzOiB7XG4gICAgb25lOiAnZW4gdGltbWUnLFxuICAgIG90aGVyOiAne3tjb3VudH19IHRpbW1hcidcbiAgfSxcbiAgeERheXM6IHtcbiAgICBvbmU6ICdlbiBkYWcnLFxuICAgIG90aGVyOiAne3tjb3VudH19IGRhZ2FyJ1xuICB9LFxuICBhYm91dFhXZWVrczoge1xuICAgIG9uZTogJ3VuZ2Vmw6RyIGVuIHZlY2thJyxcbiAgICBvdGhlcjogJ3VuZ2Vmw6RyIHt7Y291bnR9fSB2ZWNrYSdcbiAgfSxcbiAgeFdlZWtzOiB7XG4gICAgb25lOiAnZW4gdmVja2EnLFxuICAgIG90aGVyOiAne3tjb3VudH19IHZlY2thJ1xuICB9LFxuICBhYm91dFhNb250aHM6IHtcbiAgICBvbmU6ICd1bmdlZsOkciBlbiBtw6VuYWQnLFxuICAgIG90aGVyOiAndW5nZWbDpHIge3tjb3VudH19IG3DpW5hZGVyJ1xuICB9LFxuICB4TW9udGhzOiB7XG4gICAgb25lOiAnZW4gbcOlbmFkJyxcbiAgICBvdGhlcjogJ3t7Y291bnR9fSBtw6VuYWRlcidcbiAgfSxcbiAgYWJvdXRYWWVhcnM6IHtcbiAgICBvbmU6ICd1bmdlZsOkciBldHQgw6VyJyxcbiAgICBvdGhlcjogJ3VuZ2Vmw6RyIHt7Y291bnR9fSDDpXInXG4gIH0sXG4gIHhZZWFyczoge1xuICAgIG9uZTogJ2V0dCDDpXInLFxuICAgIG90aGVyOiAne3tjb3VudH19IMOlcidcbiAgfSxcbiAgb3ZlclhZZWFyczoge1xuICAgIG9uZTogJ8O2dmVyIGV0dCDDpXInLFxuICAgIG90aGVyOiAnw7Z2ZXIge3tjb3VudH19IMOlcidcbiAgfSxcbiAgYWxtb3N0WFllYXJzOiB7XG4gICAgb25lOiAnbsOkc3RhbiBldHQgw6VyJyxcbiAgICBvdGhlcjogJ27DpHN0YW4ge3tjb3VudH19IMOlcidcbiAgfVxufTtcbnZhciB3b3JkTWFwcGluZyA9IFsnbm9sbCcsICdlbicsICd0dsOlJywgJ3RyZScsICdmeXJhJywgJ2ZlbScsICdzZXgnLCAnc2p1JywgJ8OldHRhJywgJ25pbycsICd0aW8nLCAnZWx2YScsICd0b2x2J107XG52YXIgZm9ybWF0RGlzdGFuY2UgPSBmdW5jdGlvbiBmb3JtYXREaXN0YW5jZSh0b2tlbiwgY291bnQsIG9wdGlvbnMpIHtcbiAgdmFyIHJlc3VsdDtcbiAgdmFyIHRva2VuVmFsdWUgPSBmb3JtYXREaXN0YW5jZUxvY2FsZVt0b2tlbl07XG4gIGlmICh0eXBlb2YgdG9rZW5WYWx1ZSA9PT0gJ3N0cmluZycpIHtcbiAgICByZXN1bHQgPSB0b2tlblZhbHVlO1xuICB9IGVsc2UgaWYgKGNvdW50ID09PSAxKSB7XG4gICAgcmVzdWx0ID0gdG9rZW5WYWx1ZS5vbmU7XG4gIH0gZWxzZSB7XG4gICAgaWYgKG9wdGlvbnMgJiYgb3B0aW9ucy5vbmx5TnVtZXJpYykge1xuICAgICAgcmVzdWx0ID0gdG9rZW5WYWx1ZS5vdGhlci5yZXBsYWNlKCd7e2NvdW50fX0nLCBTdHJpbmcoY291bnQpKTtcbiAgICB9IGVsc2Uge1xuICAgICAgcmVzdWx0ID0gdG9rZW5WYWx1ZS5vdGhlci5yZXBsYWNlKCd7e2NvdW50fX0nLCBjb3VudCA8IDEzID8gd29yZE1hcHBpbmdbY291bnRdIDogU3RyaW5nKGNvdW50KSk7XG4gICAgfVxuICB9XG4gIGlmIChvcHRpb25zICE9PSBudWxsICYmIG9wdGlvbnMgIT09IHZvaWQgMCAmJiBvcHRpb25zLmFkZFN1ZmZpeCkge1xuICAgIGlmIChvcHRpb25zLmNvbXBhcmlzb24gJiYgb3B0aW9ucy5jb21wYXJpc29uID4gMCkge1xuICAgICAgcmV0dXJuICdvbSAnICsgcmVzdWx0O1xuICAgIH0gZWxzZSB7XG4gICAgICByZXR1cm4gcmVzdWx0ICsgJyBzZWRhbic7XG4gICAgfVxuICB9XG4gIHJldHVybiByZXN1bHQ7XG59O1xudmFyIF9kZWZhdWx0ID0gZm9ybWF0RGlzdGFuY2U7XG5leHBvcnRzLmRlZmF1bHQgPSBfZGVmYXVsdDtcbm1vZHVsZS5leHBvcnRzID0gZXhwb3J0cy5kZWZhdWx0OyJdLCJuYW1lcyI6W10sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/sv/_lib/formatDistance/index.js\n");

/***/ })

}]);