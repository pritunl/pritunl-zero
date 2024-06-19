"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-hu-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/hu/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/hu/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar accusativeWeekdays = ['vasárnap', 'hétfőn', 'kedden', 'szerdán', 'csütörtökön', 'pénteken', 'szombaton'];\nfunction week(isFuture) {\n  return function (date) {\n    var weekday = accusativeWeekdays[date.getUTCDay()];\n    var prefix = isFuture ? '' : \"'múlt' \";\n    return \"\".concat(prefix, \"'\").concat(weekday, \"' p'-kor'\");\n  };\n}\nvar formatRelativeLocale = {\n  lastWeek: week(false),\n  yesterday: \"'tegnap' p'-kor'\",\n  today: \"'ma' p'-kor'\",\n  tomorrow: \"'holnap' p'-kor'\",\n  nextWeek: week(true),\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date) {\n  var format = formatRelativeLocale[token];\n  if (typeof format === 'function') {\n    return format(date);\n  }\n  return format;\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2h1L19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxrQkFBZTtBQUNmIiwic291cmNlcyI6WyJ3ZWJwYWNrOi8vcHJpdHVubC16ZXJvLy4vbm9kZV9tb2R1bGVzL2RhdGUtZm5zL2xvY2FsZS9odS9fbGliL2Zvcm1hdFJlbGF0aXZlL2luZGV4LmpzP2UxOTQiXSwic291cmNlc0NvbnRlbnQiOlsiXCJ1c2Ugc3RyaWN0XCI7XG5cbk9iamVjdC5kZWZpbmVQcm9wZXJ0eShleHBvcnRzLCBcIl9fZXNNb2R1bGVcIiwge1xuICB2YWx1ZTogdHJ1ZVxufSk7XG5leHBvcnRzLmRlZmF1bHQgPSB2b2lkIDA7XG52YXIgYWNjdXNhdGl2ZVdlZWtkYXlzID0gWyd2YXPDoXJuYXAnLCAnaMOpdGbFkW4nLCAna2VkZGVuJywgJ3N6ZXJkw6FuJywgJ2Nzw7x0w7ZydMO2a8O2bicsICdww6ludGVrZW4nLCAnc3pvbWJhdG9uJ107XG5mdW5jdGlvbiB3ZWVrKGlzRnV0dXJlKSB7XG4gIHJldHVybiBmdW5jdGlvbiAoZGF0ZSkge1xuICAgIHZhciB3ZWVrZGF5ID0gYWNjdXNhdGl2ZVdlZWtkYXlzW2RhdGUuZ2V0VVRDRGF5KCldO1xuICAgIHZhciBwcmVmaXggPSBpc0Z1dHVyZSA/ICcnIDogXCInbcO6bHQnIFwiO1xuICAgIHJldHVybiBcIlwiLmNvbmNhdChwcmVmaXgsIFwiJ1wiKS5jb25jYXQod2Vla2RheSwgXCInIHAnLWtvcidcIik7XG4gIH07XG59XG52YXIgZm9ybWF0UmVsYXRpdmVMb2NhbGUgPSB7XG4gIGxhc3RXZWVrOiB3ZWVrKGZhbHNlKSxcbiAgeWVzdGVyZGF5OiBcIid0ZWduYXAnIHAnLWtvcidcIixcbiAgdG9kYXk6IFwiJ21hJyBwJy1rb3InXCIsXG4gIHRvbW9ycm93OiBcIidob2xuYXAnIHAnLWtvcidcIixcbiAgbmV4dFdlZWs6IHdlZWsodHJ1ZSksXG4gIG90aGVyOiAnUCdcbn07XG52YXIgZm9ybWF0UmVsYXRpdmUgPSBmdW5jdGlvbiBmb3JtYXRSZWxhdGl2ZSh0b2tlbiwgZGF0ZSkge1xuICB2YXIgZm9ybWF0ID0gZm9ybWF0UmVsYXRpdmVMb2NhbGVbdG9rZW5dO1xuICBpZiAodHlwZW9mIGZvcm1hdCA9PT0gJ2Z1bmN0aW9uJykge1xuICAgIHJldHVybiBmb3JtYXQoZGF0ZSk7XG4gIH1cbiAgcmV0dXJuIGZvcm1hdDtcbn07XG52YXIgX2RlZmF1bHQgPSBmb3JtYXRSZWxhdGl2ZTtcbmV4cG9ydHMuZGVmYXVsdCA9IF9kZWZhdWx0O1xubW9kdWxlLmV4cG9ydHMgPSBleHBvcnRzLmRlZmF1bHQ7Il0sIm5hbWVzIjpbXSwic291cmNlUm9vdCI6IiJ9\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/hu/_lib/formatRelative/index.js\n");

/***/ })

}]);