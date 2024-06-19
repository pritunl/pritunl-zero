"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-lb-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/lb/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/lb/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatRelativeLocale = {\n  lastWeek: function lastWeek(date) {\n    var day = date.getUTCDay();\n    var result = \"'läschte\";\n    if (day === 2 || day === 4) {\n      // Eifeler Regel: Add an n before the consonant d; Here \"Dënschdeg\" \"and Donneschde\".\n      result += 'n';\n    }\n    result += \"' eeee 'um' p\";\n    return result;\n  },\n  yesterday: \"'gëschter um' p\",\n  today: \"'haut um' p\",\n  tomorrow: \"'moien um' p\",\n  nextWeek: \"eeee 'um' p\",\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date, _baseDate, _options) {\n  var format = formatRelativeLocale[token];\n  if (typeof format === 'function') {\n    return format(date);\n  }\n  return format;\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2xiL19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EseURBQXlEO0FBQ3pEO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsR0FBRztBQUNIO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxrQkFBZTtBQUNmIiwic291cmNlcyI6WyJ3ZWJwYWNrOi8vcHJpdHVubC16ZXJvLy4vbm9kZV9tb2R1bGVzL2RhdGUtZm5zL2xvY2FsZS9sYi9fbGliL2Zvcm1hdFJlbGF0aXZlL2luZGV4LmpzP2UyYzQiXSwic291cmNlc0NvbnRlbnQiOlsiXCJ1c2Ugc3RyaWN0XCI7XG5cbk9iamVjdC5kZWZpbmVQcm9wZXJ0eShleHBvcnRzLCBcIl9fZXNNb2R1bGVcIiwge1xuICB2YWx1ZTogdHJ1ZVxufSk7XG5leHBvcnRzLmRlZmF1bHQgPSB2b2lkIDA7XG52YXIgZm9ybWF0UmVsYXRpdmVMb2NhbGUgPSB7XG4gIGxhc3RXZWVrOiBmdW5jdGlvbiBsYXN0V2VlayhkYXRlKSB7XG4gICAgdmFyIGRheSA9IGRhdGUuZ2V0VVRDRGF5KCk7XG4gICAgdmFyIHJlc3VsdCA9IFwiJ2zDpHNjaHRlXCI7XG4gICAgaWYgKGRheSA9PT0gMiB8fCBkYXkgPT09IDQpIHtcbiAgICAgIC8vIEVpZmVsZXIgUmVnZWw6IEFkZCBhbiBuIGJlZm9yZSB0aGUgY29uc29uYW50IGQ7IEhlcmUgXCJEw6tuc2NoZGVnXCIgXCJhbmQgRG9ubmVzY2hkZVwiLlxuICAgICAgcmVzdWx0ICs9ICduJztcbiAgICB9XG4gICAgcmVzdWx0ICs9IFwiJyBlZWVlICd1bScgcFwiO1xuICAgIHJldHVybiByZXN1bHQ7XG4gIH0sXG4gIHllc3RlcmRheTogXCInZ8Orc2NodGVyIHVtJyBwXCIsXG4gIHRvZGF5OiBcIidoYXV0IHVtJyBwXCIsXG4gIHRvbW9ycm93OiBcIidtb2llbiB1bScgcFwiLFxuICBuZXh0V2VlazogXCJlZWVlICd1bScgcFwiLFxuICBvdGhlcjogJ1AnXG59O1xudmFyIGZvcm1hdFJlbGF0aXZlID0gZnVuY3Rpb24gZm9ybWF0UmVsYXRpdmUodG9rZW4sIGRhdGUsIF9iYXNlRGF0ZSwgX29wdGlvbnMpIHtcbiAgdmFyIGZvcm1hdCA9IGZvcm1hdFJlbGF0aXZlTG9jYWxlW3Rva2VuXTtcbiAgaWYgKHR5cGVvZiBmb3JtYXQgPT09ICdmdW5jdGlvbicpIHtcbiAgICByZXR1cm4gZm9ybWF0KGRhdGUpO1xuICB9XG4gIHJldHVybiBmb3JtYXQ7XG59O1xudmFyIF9kZWZhdWx0ID0gZm9ybWF0UmVsYXRpdmU7XG5leHBvcnRzLmRlZmF1bHQgPSBfZGVmYXVsdDtcbm1vZHVsZS5leHBvcnRzID0gZXhwb3J0cy5kZWZhdWx0OyJdLCJuYW1lcyI6W10sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/lb/_lib/formatRelative/index.js\n");

/***/ })

}]);