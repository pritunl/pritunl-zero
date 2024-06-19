"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(self["webpackChunkpritunl_zero"] = self["webpackChunkpritunl_zero"] || []).push([["date-fns-locale-cs-_lib-formatRelative-index-js"],{

/***/ "./node_modules/date-fns/locale/cs/_lib/formatRelative/index.js":
/*!**********************************************************************!*\
  !*** ./node_modules/date-fns/locale/cs/_lib/formatRelative/index.js ***!
  \**********************************************************************/
/***/ ((module, exports) => {

eval("\n\nObject.defineProperty(exports, \"__esModule\", ({\n  value: true\n}));\nexports[\"default\"] = void 0;\nvar accusativeWeekdays = ['neděli', 'pondělí', 'úterý', 'středu', 'čtvrtek', 'pátek', 'sobotu'];\nvar formatRelativeLocale = {\n  lastWeek: \"'poslední' eeee 've' p\",\n  yesterday: \"'včera v' p\",\n  today: \"'dnes v' p\",\n  tomorrow: \"'zítra v' p\",\n  nextWeek: function nextWeek(date) {\n    var day = date.getUTCDay();\n    return \"'v \" + accusativeWeekdays[day] + \" o' p\";\n  },\n  other: 'P'\n};\nvar formatRelative = function formatRelative(token, date) {\n  var format = formatRelativeLocale[token];\n  if (typeof format === 'function') {\n    return format(date);\n  }\n  return format;\n};\nvar _default = formatRelative;\nexports[\"default\"] = _default;\nmodule.exports = exports.default;//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9ub2RlX21vZHVsZXMvZGF0ZS1mbnMvbG9jYWxlL2NzL19saWIvZm9ybWF0UmVsYXRpdmUvaW5kZXguanMiLCJtYXBwaW5ncyI6IkFBQWE7O0FBRWIsOENBQTZDO0FBQzdDO0FBQ0EsQ0FBQyxFQUFDO0FBQ0Ysa0JBQWU7QUFDZjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxHQUFHO0FBQ0g7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxrQkFBZTtBQUNmIiwic291cmNlcyI6WyJ3ZWJwYWNrOi8vcHJpdHVubC16ZXJvLy4vbm9kZV9tb2R1bGVzL2RhdGUtZm5zL2xvY2FsZS9jcy9fbGliL2Zvcm1hdFJlbGF0aXZlL2luZGV4LmpzPzkyYzgiXSwic291cmNlc0NvbnRlbnQiOlsiXCJ1c2Ugc3RyaWN0XCI7XG5cbk9iamVjdC5kZWZpbmVQcm9wZXJ0eShleHBvcnRzLCBcIl9fZXNNb2R1bGVcIiwge1xuICB2YWx1ZTogdHJ1ZVxufSk7XG5leHBvcnRzLmRlZmF1bHQgPSB2b2lkIDA7XG52YXIgYWNjdXNhdGl2ZVdlZWtkYXlzID0gWyduZWTEm2xpJywgJ3BvbmTEm2zDrScsICfDunRlcsO9JywgJ3N0xZllZHUnLCAnxI10dnJ0ZWsnLCAncMOhdGVrJywgJ3NvYm90dSddO1xudmFyIGZvcm1hdFJlbGF0aXZlTG9jYWxlID0ge1xuICBsYXN0V2VlazogXCIncG9zbGVkbsOtJyBlZWVlICd2ZScgcFwiLFxuICB5ZXN0ZXJkYXk6IFwiJ3bEjWVyYSB2JyBwXCIsXG4gIHRvZGF5OiBcIidkbmVzIHYnIHBcIixcbiAgdG9tb3Jyb3c6IFwiJ3rDrXRyYSB2JyBwXCIsXG4gIG5leHRXZWVrOiBmdW5jdGlvbiBuZXh0V2VlayhkYXRlKSB7XG4gICAgdmFyIGRheSA9IGRhdGUuZ2V0VVRDRGF5KCk7XG4gICAgcmV0dXJuIFwiJ3YgXCIgKyBhY2N1c2F0aXZlV2Vla2RheXNbZGF5XSArIFwiIG8nIHBcIjtcbiAgfSxcbiAgb3RoZXI6ICdQJ1xufTtcbnZhciBmb3JtYXRSZWxhdGl2ZSA9IGZ1bmN0aW9uIGZvcm1hdFJlbGF0aXZlKHRva2VuLCBkYXRlKSB7XG4gIHZhciBmb3JtYXQgPSBmb3JtYXRSZWxhdGl2ZUxvY2FsZVt0b2tlbl07XG4gIGlmICh0eXBlb2YgZm9ybWF0ID09PSAnZnVuY3Rpb24nKSB7XG4gICAgcmV0dXJuIGZvcm1hdChkYXRlKTtcbiAgfVxuICByZXR1cm4gZm9ybWF0O1xufTtcbnZhciBfZGVmYXVsdCA9IGZvcm1hdFJlbGF0aXZlO1xuZXhwb3J0cy5kZWZhdWx0ID0gX2RlZmF1bHQ7XG5tb2R1bGUuZXhwb3J0cyA9IGV4cG9ydHMuZGVmYXVsdDsiXSwibmFtZXMiOltdLCJzb3VyY2VSb290IjoiIn0=\n//# sourceURL=webpack-internal:///./node_modules/date-fns/locale/cs/_lib/formatRelative/index.js\n");

/***/ })

}]);