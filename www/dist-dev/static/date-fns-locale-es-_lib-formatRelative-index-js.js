"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-es-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/es/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/es/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatRelativeLocale = {\n  lastWeek: \"'el' eeee 'pasado a la' p\",\n  yesterday: \"'ayer a la' p\",\n  today: \"'hoy a la' p\",\n  tomorrow: \"'mañana a la' p\",\n  nextWeek: \"eeee 'a la' p\",\n  other: 'P'\n};\nvar formatRelativeLocalePlural = {\n  lastWeek: \"'el' eeee 'pasado a las' p\",\n  yesterday: \"'ayer a las' p\",\n  today: \"'hoy a las' p\",\n  tomorrow: \"'mañana a las' p\",\n  nextWeek: \"eeee 'a las' p\",\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date, _baseDate, _options) {\n  if (date.getUTCHours() !== 1) {\n    return formatRelativeLocalePlural[token];\n  } else {\n    return formatRelativeLocale[token];\n  }\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2VzL19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLElBQUk7QUFDSjtBQUNBO0FBQ0E7QUFDQTtBQUNBLGtCQUFlO0FBQ2YiLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9wcml0dW5sLXplcm8vLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2VzL19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanM/NTkwYyJdLCJzb3VyY2VzQ29udGVudCI6WyJcInVzZSBzdHJpY3RcIjtcblxuT2JqZWN0LmRlZmluZVByb3BlcnR5KGV4cG9ydHMsIFwiX19lc01vZHVsZVwiLCB7XG4gIHZhbHVlOiB0cnVlXG59KTtcbmV4cG9ydHMuZGVmYXVsdCA9IHZvaWQgMDtcbnZhciBmb3JtYXRSZWxhdGl2ZUxvY2FsZSA9IHtcbiAgbGFzdFdlZWs6IFwiJ2VsJyBlZWVlICdwYXNhZG8gYSBsYScgcFwiLFxuICB5ZXN0ZXJkYXk6IFwiJ2F5ZXIgYSBsYScgcFwiLFxuICB0b2RheTogXCInaG95IGEgbGEnIHBcIixcbiAgdG9tb3Jyb3c6IFwiJ21hw7FhbmEgYSBsYScgcFwiLFxuICBuZXh0V2VlazogXCJlZWVlICdhIGxhJyBwXCIsXG4gIG90aGVyOiAnUCdcbn07XG52YXIgZm9ybWF0UmVsYXRpdmVMb2NhbGVQbHVyYWwgPSB7XG4gIGxhc3RXZWVrOiBcIidlbCcgZWVlZSAncGFzYWRvIGEgbGFzJyBwXCIsXG4gIHllc3RlcmRheTogXCInYXllciBhIGxhcycgcFwiLFxuICB0b2RheTogXCInaG95IGEgbGFzJyBwXCIsXG4gIHRvbW9ycm93OiBcIidtYcOxYW5hIGEgbGFzJyBwXCIsXG4gIG5leHRXZWVrOiBcImVlZWUgJ2EgbGFzJyBwXCIsXG4gIG90aGVyOiAnUCdcbn07XG52YXIgZm9ybWF0UmVsYXRpdmUgPSBmdW5jdGlvbiBmb3JtYXRSZWxhdGl2ZSh0b2tlbiwgZGF0ZSwgX2Jhc2VEYXRlLCBfb3B0aW9ucykge1xuICBpZiAoZGF0ZS5nZXRVVENIb3VycygpICE9PSAxKSB7XG4gICAgcmV0dXJuIGZvcm1hdFJlbGF0aXZlTG9jYWxlUGx1cmFsW3Rva2VuXTtcbiAgfSBlbHNlIHtcbiAgICByZXR1cm4gZm9ybWF0UmVsYXRpdmVMb2NhbGVbdG9rZW5dO1xuICB9XG59O1xudmFyIF9kZWZhdWx0ID0gZm9ybWF0UmVsYXRpdmU7XG5leHBvcnRzLmRlZmF1bHQgPSBfZGVmYXVsdDtcbm1vZHVsZS5leHBvcnRzID0gZXhwb3J0cy5kZWZhdWx0OyJdLCJuYW1lcyI6W10sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/es/_lib/formatRelative/index.js\n");

/***/ })

}]);