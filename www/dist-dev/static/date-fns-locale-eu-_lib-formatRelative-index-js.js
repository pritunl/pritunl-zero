"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-eu-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/eu/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/eu/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatRelativeLocale = {\n  lastWeek: \"'joan den' eeee, LT\",\n  yesterday: \"'atzo,' p\",\n  today: \"'gaur,' p\",\n  tomorrow: \"'bihar,' p\",\n  nextWeek: 'eeee, p',\n  other: 'P'\n};\nvar formatRelativeLocalePlural = {\n  lastWeek: \"'joan den' eeee, p\",\n  yesterday: \"'atzo,' p\",\n  today: \"'gaur,' p\",\n  tomorrow: \"'bihar,' p\",\n  nextWeek: 'eeee, p',\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date) {\n  if (date.getUTCHours() !== 1) {\n    return formatRelativeLocalePlural[token];\n  }\n  return formatRelativeLocale[token];\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2V1L19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0Esa0JBQWU7QUFDZiIsInNvdXJjZXMiOlsid2VicGFjazovL3ByaXR1bmwtemVyby8uL25vZGVfbW9kdWxlcy9kYXRlLWZucy9sb2NhbGUvZXUvX2xpYi9mb3JtYXRSZWxhdGl2ZS9pbmRleC5qcz8zMmYxIl0sInNvdXJjZXNDb250ZW50IjpbIlwidXNlIHN0cmljdFwiO1xuXG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgXCJfX2VzTW9kdWxlXCIsIHtcbiAgdmFsdWU6IHRydWVcbn0pO1xuZXhwb3J0cy5kZWZhdWx0ID0gdm9pZCAwO1xudmFyIGZvcm1hdFJlbGF0aXZlTG9jYWxlID0ge1xuICBsYXN0V2VlazogXCInam9hbiBkZW4nIGVlZWUsIExUXCIsXG4gIHllc3RlcmRheTogXCInYXR6bywnIHBcIixcbiAgdG9kYXk6IFwiJ2dhdXIsJyBwXCIsXG4gIHRvbW9ycm93OiBcIidiaWhhciwnIHBcIixcbiAgbmV4dFdlZWs6ICdlZWVlLCBwJyxcbiAgb3RoZXI6ICdQJ1xufTtcbnZhciBmb3JtYXRSZWxhdGl2ZUxvY2FsZVBsdXJhbCA9IHtcbiAgbGFzdFdlZWs6IFwiJ2pvYW4gZGVuJyBlZWVlLCBwXCIsXG4gIHllc3RlcmRheTogXCInYXR6bywnIHBcIixcbiAgdG9kYXk6IFwiJ2dhdXIsJyBwXCIsXG4gIHRvbW9ycm93OiBcIidiaWhhciwnIHBcIixcbiAgbmV4dFdlZWs6ICdlZWVlLCBwJyxcbiAgb3RoZXI6ICdQJ1xufTtcbnZhciBmb3JtYXRSZWxhdGl2ZSA9IGZ1bmN0aW9uIGZvcm1hdFJlbGF0aXZlKHRva2VuLCBkYXRlKSB7XG4gIGlmIChkYXRlLmdldFVUQ0hvdXJzKCkgIT09IDEpIHtcbiAgICByZXR1cm4gZm9ybWF0UmVsYXRpdmVMb2NhbGVQbHVyYWxbdG9rZW5dO1xuICB9XG4gIHJldHVybiBmb3JtYXRSZWxhdGl2ZUxvY2FsZVt0b2tlbl07XG59O1xudmFyIF9kZWZhdWx0ID0gZm9ybWF0UmVsYXRpdmU7XG5leHBvcnRzLmRlZmF1bHQgPSBfZGVmYXVsdDtcbm1vZHVsZS5leHBvcnRzID0gZXhwb3J0cy5kZWZhdWx0OyJdLCJuYW1lcyI6W10sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/eu/_lib/formatRelative/index.js\n");

/***/ })

}]);