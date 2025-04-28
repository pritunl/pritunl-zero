# pritunl-zero: zero trust system

[![github](https://img.shields.io/badge/github-pritunl-11bdc2.svg?style=flat)](https://github.com/pritunl)
[![twitter](https://img.shields.io/badge/twitter-pritunl-55acee.svg?style=flat)](https://twitter.com/pritunl)
[![medium](https://img.shields.io/badge/medium-pritunl-b32b2b.svg?style=flat)](https://pritunl.medium.com)
[![forum](https://img.shields.io/badge/discussion-forum-ffffff.svg?style=flat)](https://forum.pritunl.com)

[Pritunl-Zero](https://zero.pritunl.com) is a zero trust system
that provides secure authenticated access to internal services from untrusted
networks without the use of a VPN. Documentation and more information can be
found at [docs.pritunl.com](https://docs.pritunl.com/docs/pritunl-zero)

[![pritunl](img/logo_code.png)](https://docs.pritunl.com/docs/pritunl-zero)

## Install from Source

```bash
# Install Go
sudo dnf -y install git-core

wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
echo "68097bd680839cbc9d464a0edce4f7c333975e27a90246890e9f1078c7e702ad go1.24.2.linux-amd64.tar.gz" | sha256sum -c -

sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xf go1.24.2.linux-amd64.tar.gz
rm -f go1.24.2.linux-amd64.tar.gz

tee -a ~/.bashrc << EOF
export GOPATH=\$HOME/go
export GOROOT=/usr/local/go
export PATH=/usr/local/go/bin:\$PATH
EOF
source ~/.bashrc

# Install MongoDB
sudo tee /etc/yum.repos.d/mongodb-org.repo << EOF
[mongodb-org]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/9/mongodb-org/8.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://pgp.mongodb.com/server-8.0.asc
EOF

sudo dnf -y install mongodb-org
sudo systemctl enable --now mongod

# Install Pritunl Zero
go install -v github.com/pritunl/pritunl-zero@latest

# Setup systemd units
sudo cp $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/tools/pritunl-zero.service /etc/systemd/system/
sudo cp $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/tools/pritunl-zero-redirect.socket /etc/systemd/system/
sudo cp $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/tools/pritunl-zero-redirect.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo useradd -r -s /sbin/nologin -c 'Pritunl web server' pritunl-zero-web

# Install Pritunl Zero
sudo mkdir -p /usr/share/pritunl-zero/www/
sudo cp -r $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/www/dist/. /usr/share/pritunl-zero/www/
sudo cp ~/go/bin/pritunl-zero /usr/bin/pritunl-zero
sudo systemctl enable --now pritunl-zero
```

## License

Please refer to the [`LICENSE`](LICENSE) file for a copy of the license.
