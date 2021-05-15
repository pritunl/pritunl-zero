### pritunl-zero-www

```
npm install
rm node_modules/react-stripe-checkout/index.d.ts
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
tsc --watch
webpack-cli --config webpack.dev.config --progress --color --watch
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

### internal

```
# desktop
rsync --human-readable --archive --xattrs --progress --delete --exclude "/node_modules/*" --exclude "/jspm_packages/*" --exclude "app/*.js" --exclude "app/*.js.map" --exclude "app/**/*.js" --exclude "app/**/*.js.map" --exclude "uapp/*.js" --exclude "uapp/*.js.map" --exclude "uapp/**/*.js" --exclude "uapp/**/*.js.map" /home/cloud/go/src/github.com/pritunl/pritunl-zero/www/ $NPM_SERVER:/home/cloud/pritunl-zero-www/

# npm-server
cd /home/cloud/pritunl-zero-www/
rm package-lock.json
rm -rf node_modules
npm install
rm node_modules/react-stripe-checkout/index.d.ts

# desktop
scp $NPM_SERVER:/home/cloud/pritunl-zero-www/package.json /home/cloud/go/src/github.com/pritunl/pritunl-zero/www/package.json
scp $NPM_SERVER:/home/cloud/pritunl-zero-www/package-lock.json /home/cloud/go/src/github.com/pritunl/pritunl-zero/www/package-lock.json
rsync --human-readable --archive --xattrs --progress --delete $NPM_SERVER:/home/cloud/pritunl-zero-www/node_modules/ /home/cloud/go/src/github.com/pritunl/pritunl-zero/www/node_modules/

# npm-server
rsync --human-readable --archive --xattrs --progress --delete --exclude "/node_modules/*" --exclude "/jspm_packages/*" --exclude "app/*.js" --exclude "app/*.js.map" --exclude "app/**/*.js" --exclude "app/**/*.js.map" --exclude "uapp/*.js" --exclude "uapp/*.js.map" --exclude "uapp/**/*.js" --exclude "uapp/**/*.js.map" /home/cloud/go/src/github.com/pritunl/pritunl-zero/www/ $NPM_SERVER:/home/cloud/pritunl-zero-www/
sh build.sh

# desktop
rsync --human-readable --archive --xattrs --progress --delete $NPM_SERVER:/home/cloud/pritunl-zero-www/dist/ /home/cloud/go/src/github.com/pritunl/pritunl-zero/www/dist/
rsync --human-readable --archive --xattrs --progress --delete $NPM_SERVER:/home/cloud/pritunl-zero-www/dist-dev/ /home/cloud/go/src/github.com/pritunl/pritunl-zero/www/dist-dev/
```
