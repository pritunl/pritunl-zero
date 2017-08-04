tsc
jspm bundle app/App.js
rm -rf dist/static
mkdir -p dist/static
cp styles/global.css dist/static/
cp node_modules/normalize.css/normalize.css dist/static/
cp node_modules/@blueprintjs/core/dist/blueprint.css dist/static/
cp node_modules/@blueprintjs/datetime/dist/blueprint-datetime.css dist/static/
cp node_modules/@blueprintjs/core/resources/icons/icons-16.eot dist/static/
cp node_modules/@blueprintjs/core/resources/icons/icons-16.ttf dist/static/
cp node_modules/@blueprintjs/core/resources/icons/icons-16.woff dist/static/
cp node_modules/@blueprintjs/core/resources/icons/icons-20.eot dist/static/
cp node_modules/@blueprintjs/core/resources/icons/icons-20.ttf dist/static/
cp node_modules/@blueprintjs/core/resources/icons/icons-20.woff dist/static/
cp jspm_packages/system.js dist/static/
sed -i 's|../resources/icons/||g' dist/static/blueprint.css
mv build.js dist/static/app.js
mv build.js.map dist/static/app.js.map
cp index_dist.html dist/index.html
cp login.html dist/login.html

