tsc
rm -rf dist/static
mkdir -p dist/static
cp styles/global.css dist/static/
cp node_modules/normalize.css/normalize.css dist/static/
cp node_modules/@blueprintjs/core/lib/css/blueprint.css dist/static/
cp node_modules/@blueprintjs/datetime/lib/css/blueprint-datetime.css dist/static/
cp node_modules/@blueprintjs/icons/lib/css/blueprint-icons.css dist/static/
cp node_modules/@blueprintjs/icons/resources/icons/icons-16.eot dist/static/
cp node_modules/@blueprintjs/icons/resources/icons/icons-16.ttf dist/static/
cp node_modules/@blueprintjs/icons/resources/icons/icons-16.woff dist/static/
cp node_modules/@blueprintjs/icons/resources/icons/icons-20.eot dist/static/
cp node_modules/@blueprintjs/icons/resources/icons/icons-20.ttf dist/static/
cp node_modules/@blueprintjs/icons/resources/icons/icons-20.woff dist/static/
cp jspm_packages/system.js dist/static/
sed -i 's|../../resources/icons/||g' dist/static/blueprint-icons.css
jspm bundle app/App.js
mv build.js dist/static/app.js
mv build.js.map dist/static/app.js.map
cp index_dist.html dist/index.html
jspm bundle uapp/App.js
mv build.js dist/static/uapp.js
mv build.js.map dist/static/uapp.js.map
cp uindex_dist.html dist/uindex.html
cp login.html dist/login.html

APP_HASH=`md5sum dist/static/app.js | cut -c1-6`
UAPP_HASH=`md5sum dist/static/uapp.js | cut -c1-6`

mv dist/static/app.js dist/static/app.${APP_HASH}.js
mv dist/static/app.js.map dist/static/app.${APP_HASH}.js.map

mv dist/static/uapp.js dist/static/uapp.${UAPP_HASH}.js
mv dist/static/uapp.js.map dist/static/uapp.${UAPP_HASH}.js.map

sed -i -e "s|static/app.js|static/app.${APP_HASH}.js|g" dist/index.html
sed -i -e "s|static/uapp.js|static/uapp.${UAPP_HASH}.js|g" dist/uindex.html
