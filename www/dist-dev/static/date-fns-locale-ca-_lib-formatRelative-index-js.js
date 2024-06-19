"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-ca-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/ca/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/ca/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatRelativeLocale = {\n  lastWeek: \"'el' eeee 'passat a la' LT\",\n  yesterday: \"'ahir a la' p\",\n  today: \"'avui a la' p\",\n  tomorrow: \"'demà a la' p\",\n  nextWeek: \"eeee 'a la' p\",\n  other: 'P'\n};\nvar formatRelativeLocalePlural = {\n  lastWeek: \"'el' eeee 'passat a les' p\",\n  yesterday: \"'ahir a les' p\",\n  today: \"'avui a les' p\",\n  tomorrow: \"'demà a les' p\",\n  nextWeek: \"eeee 'a les' p\",\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date, _baseDate, _options) {\n  if (date.getUTCHours() !== 1) {\n    return formatRelativeLocalePlural[token];\n  }\n  return formatRelativeLocale[token];\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2NhL19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0Esa0JBQWU7QUFDZiIsInNvdXJjZXMiOlsid2VicGFjazovL3ByaXR1bmwtemVyby8uL25vZGVfbW9kdWxlcy9kYXRlLWZucy9sb2NhbGUvY2EvX2xpYi9mb3JtYXRSZWxhdGl2ZS9pbmRleC5qcz8wNGI1Il0sInNvdXJjZXNDb250ZW50IjpbIlwidXNlIHN0cmljdFwiO1xuXG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgXCJfX2VzTW9kdWxlXCIsIHtcbiAgdmFsdWU6IHRydWVcbn0pO1xuZXhwb3J0cy5kZWZhdWx0ID0gdm9pZCAwO1xudmFyIGZvcm1hdFJlbGF0aXZlTG9jYWxlID0ge1xuICBsYXN0V2VlazogXCInZWwnIGVlZWUgJ3Bhc3NhdCBhIGxhJyBMVFwiLFxuICB5ZXN0ZXJkYXk6IFwiJ2FoaXIgYSBsYScgcFwiLFxuICB0b2RheTogXCInYXZ1aSBhIGxhJyBwXCIsXG4gIHRvbW9ycm93OiBcIidkZW3DoCBhIGxhJyBwXCIsXG4gIG5leHRXZWVrOiBcImVlZWUgJ2EgbGEnIHBcIixcbiAgb3RoZXI6ICdQJ1xufTtcbnZhciBmb3JtYXRSZWxhdGl2ZUxvY2FsZVBsdXJhbCA9IHtcbiAgbGFzdFdlZWs6IFwiJ2VsJyBlZWVlICdwYXNzYXQgYSBsZXMnIHBcIixcbiAgeWVzdGVyZGF5OiBcIidhaGlyIGEgbGVzJyBwXCIsXG4gIHRvZGF5OiBcIidhdnVpIGEgbGVzJyBwXCIsXG4gIHRvbW9ycm93OiBcIidkZW3DoCBhIGxlcycgcFwiLFxuICBuZXh0V2VlazogXCJlZWVlICdhIGxlcycgcFwiLFxuICBvdGhlcjogJ1AnXG59O1xudmFyIGZvcm1hdFJlbGF0aXZlID0gZnVuY3Rpb24gZm9ybWF0UmVsYXRpdmUodG9rZW4sIGRhdGUsIF9iYXNlRGF0ZSwgX29wdGlvbnMpIHtcbiAgaWYgKGRhdGUuZ2V0VVRDSG91cnMoKSAhPT0gMSkge1xuICAgIHJldHVybiBmb3JtYXRSZWxhdGl2ZUxvY2FsZVBsdXJhbFt0b2tlbl07XG4gIH1cbiAgcmV0dXJuIGZvcm1hdFJlbGF0aXZlTG9jYWxlW3Rva2VuXTtcbn07XG52YXIgX2RlZmF1bHQgPSBmb3JtYXRSZWxhdGl2ZTtcbmV4cG9ydHMuZGVmYXVsdCA9IF9kZWZhdWx0O1xubW9kdWxlLmV4cG9ydHMgPSBleHBvcnRzLmRlZmF1bHQ7Il0sIm5hbWVzIjpbXSwic291cmNlUm9vdCI6IiJ9\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/ca/_lib/formatRelative/index.js\n");

/***/ })

}]);