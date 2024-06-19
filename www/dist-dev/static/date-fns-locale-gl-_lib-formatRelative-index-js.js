"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-gl-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/gl/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/gl/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatRelativeLocale = {\n  lastWeek: \"'o' eeee 'pasado á' LT\",\n  yesterday: \"'onte á' p\",\n  today: \"'hoxe á' p\",\n  tomorrow: \"'mañá á' p\",\n  nextWeek: \"eeee 'á' p\",\n  other: 'P'\n};\nvar formatRelativeLocalePlural = {\n  lastWeek: \"'o' eeee 'pasado ás' p\",\n  yesterday: \"'onte ás' p\",\n  today: \"'hoxe ás' p\",\n  tomorrow: \"'mañá ás' p\",\n  nextWeek: \"eeee 'ás' p\",\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date, _baseDate, _options) {\n  if (date.getUTCHours() !== 1) {\n    return formatRelativeLocalePlural[token];\n  }\n  return formatRelativeLocale[token];\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2dsL19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0Esa0JBQWU7QUFDZiIsInNvdXJjZXMiOlsid2VicGFjazovL3ByaXR1bmwtemVyby8uL25vZGVfbW9kdWxlcy9kYXRlLWZucy9sb2NhbGUvZ2wvX2xpYi9mb3JtYXRSZWxhdGl2ZS9pbmRleC5qcz8zMGM2Il0sInNvdXJjZXNDb250ZW50IjpbIlwidXNlIHN0cmljdFwiO1xuXG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgXCJfX2VzTW9kdWxlXCIsIHtcbiAgdmFsdWU6IHRydWVcbn0pO1xuZXhwb3J0cy5kZWZhdWx0ID0gdm9pZCAwO1xudmFyIGZvcm1hdFJlbGF0aXZlTG9jYWxlID0ge1xuICBsYXN0V2VlazogXCInbycgZWVlZSAncGFzYWRvIMOhJyBMVFwiLFxuICB5ZXN0ZXJkYXk6IFwiJ29udGUgw6EnIHBcIixcbiAgdG9kYXk6IFwiJ2hveGUgw6EnIHBcIixcbiAgdG9tb3Jyb3c6IFwiJ21hw7HDoSDDoScgcFwiLFxuICBuZXh0V2VlazogXCJlZWVlICfDoScgcFwiLFxuICBvdGhlcjogJ1AnXG59O1xudmFyIGZvcm1hdFJlbGF0aXZlTG9jYWxlUGx1cmFsID0ge1xuICBsYXN0V2VlazogXCInbycgZWVlZSAncGFzYWRvIMOhcycgcFwiLFxuICB5ZXN0ZXJkYXk6IFwiJ29udGUgw6FzJyBwXCIsXG4gIHRvZGF5OiBcIidob3hlIMOhcycgcFwiLFxuICB0b21vcnJvdzogXCInbWHDscOhIMOhcycgcFwiLFxuICBuZXh0V2VlazogXCJlZWVlICfDoXMnIHBcIixcbiAgb3RoZXI6ICdQJ1xufTtcbnZhciBmb3JtYXRSZWxhdGl2ZSA9IGZ1bmN0aW9uIGZvcm1hdFJlbGF0aXZlKHRva2VuLCBkYXRlLCBfYmFzZURhdGUsIF9vcHRpb25zKSB7XG4gIGlmIChkYXRlLmdldFVUQ0hvdXJzKCkgIT09IDEpIHtcbiAgICByZXR1cm4gZm9ybWF0UmVsYXRpdmVMb2NhbGVQbHVyYWxbdG9rZW5dO1xuICB9XG4gIHJldHVybiBmb3JtYXRSZWxhdGl2ZUxvY2FsZVt0b2tlbl07XG59O1xudmFyIF9kZWZhdWx0ID0gZm9ybWF0UmVsYXRpdmU7XG5leHBvcnRzLmRlZmF1bHQgPSBfZGVmYXVsdDtcbm1vZHVsZS5leHBvcnRzID0gZXhwb3J0cy5kZWZhdWx0OyJdLCJuYW1lcyI6W10sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/gl/_lib/formatRelative/index.js\n");

/***/ })

}]);