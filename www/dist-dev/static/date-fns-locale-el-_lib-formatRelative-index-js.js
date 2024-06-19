"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-el-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/el/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/el/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar formatRelativeLocale = {\n  lastWeek: function lastWeek(date) {\n    switch (date.getUTCDay()) {\n      case 6:\n        //Σάββατο\n        return \"'το προηγούμενο' eeee 'στις' p\";\n      default:\n        return \"'την προηγούμενη' eeee 'στις' p\";\n    }\n  },\n  yesterday: \"'χθες στις' p\",\n  today: \"'σήμερα στις' p\",\n  tomorrow: \"'αύριο στις' p\",\n  nextWeek: \"eeee 'στις' p\",\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date) {\n  var format = formatRelativeLocale[token];\n  if (typeof format === 'function') return format(date);\n  return format;\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2VsL19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxHQUFHO0FBQ0g7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0Esa0JBQWU7QUFDZiIsInNvdXJjZXMiOlsid2VicGFjazovL3ByaXR1bmwtemVyby8uL25vZGVfbW9kdWxlcy9kYXRlLWZucy9sb2NhbGUvZWwvX2xpYi9mb3JtYXRSZWxhdGl2ZS9pbmRleC5qcz8yM2Y4Il0sInNvdXJjZXNDb250ZW50IjpbIlwidXNlIHN0cmljdFwiO1xuXG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgXCJfX2VzTW9kdWxlXCIsIHtcbiAgdmFsdWU6IHRydWVcbn0pO1xuZXhwb3J0cy5kZWZhdWx0ID0gdm9pZCAwO1xudmFyIGZvcm1hdFJlbGF0aXZlTG9jYWxlID0ge1xuICBsYXN0V2VlazogZnVuY3Rpb24gbGFzdFdlZWsoZGF0ZSkge1xuICAgIHN3aXRjaCAoZGF0ZS5nZXRVVENEYXkoKSkge1xuICAgICAgY2FzZSA2OlxuICAgICAgICAvL86jzqzOss6yzrHPhM6/XG4gICAgICAgIHJldHVybiBcIifPhM6/IM+Az4HOv863zrPOv8+NzrzOtc69zr8nIGVlZWUgJ8+Dz4TOuc+CJyBwXCI7XG4gICAgICBkZWZhdWx0OlxuICAgICAgICByZXR1cm4gXCInz4TOt869IM+Az4HOv863zrPOv8+NzrzOtc69zrcnIGVlZWUgJ8+Dz4TOuc+CJyBwXCI7XG4gICAgfVxuICB9LFxuICB5ZXN0ZXJkYXk6IFwiJ8+HzrjOtc+CIM+Dz4TOuc+CJyBwXCIsXG4gIHRvZGF5OiBcIifPg86uzrzOtc+BzrEgz4PPhM65z4InIHBcIixcbiAgdG9tb3Jyb3c6IFwiJ86xz43Pgc65zr8gz4PPhM65z4InIHBcIixcbiAgbmV4dFdlZWs6IFwiZWVlZSAnz4PPhM65z4InIHBcIixcbiAgb3RoZXI6ICdQJ1xufTtcbnZhciBmb3JtYXRSZWxhdGl2ZSA9IGZ1bmN0aW9uIGZvcm1hdFJlbGF0aXZlKHRva2VuLCBkYXRlKSB7XG4gIHZhciBmb3JtYXQgPSBmb3JtYXRSZWxhdGl2ZUxvY2FsZVt0b2tlbl07XG4gIGlmICh0eXBlb2YgZm9ybWF0ID09PSAnZnVuY3Rpb24nKSByZXR1cm4gZm9ybWF0KGRhdGUpO1xuICByZXR1cm4gZm9ybWF0O1xufTtcbnZhciBfZGVmYXVsdCA9IGZvcm1hdFJlbGF0aXZlO1xuZXhwb3J0cy5kZWZhdWx0ID0gX2RlZmF1bHQ7XG5tb2R1bGUuZXhwb3J0cyA9IGV4cG9ydHMuZGVmYXVsdDsiXSwibmFtZXMiOltdLCJzb3VyY2VSb290IjoiIn0=\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/el/_lib/formatRelative/index.js\n");

/***/ })

}]);