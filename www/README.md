### pritunl-zero-www

Requires [jspm](https://www.npmjs.com/package/jspm)

```
npm install
jspm install
sed -i 's|lib/node/index.js|lib/client.js|g' jspm_packages/npm/superagent@*.js
```

#### lint

```
tslint -c tslint.json app/*.ts*
tslint -c tslint.json app/**/*.ts*
tslint -c tslint.json uapp/*.ts*
tslint -c tslint.json uapp/**/*.ts*
```

### development

```
tsc
jspm depcache app/App.js
jspm depcache uapp/App.js
tsc --watch
```

#### production

```
sh build.sh
```

### clean

```
rm -rf app/*.js*
rm -rf app/**/*.js*
rm -rf uapp/*.js*
rm -rf uapp/**/*.js*
```
