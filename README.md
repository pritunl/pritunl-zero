# pritunl-zero: zero trust system

[![github](https://img.shields.io/badge/github-pritunl-11bdc2.svg?style=flat)](https://github.com/pritunl)
[![twitter](https://img.shields.io/badge/twitter-pritunl-55acee.svg?style=flat)](https://twitter.com/pritunl)
[![substack](https://img.shields.io/badge/substack-pritunl-ff6719.svg?style=flat)](https://pritunl.substack.com/)
[![forum](https://img.shields.io/badge/discussion-forum-ffffff.svg?style=flat)](https://forum.pritunl.com)

[Pritunl-Zero](https://zero.pritunl.com) is a zero trust system
that provides secure authenticated access to internal services from untrusted
networks without the use of a VPN. Documentation and more information can be
found at [docs.pritunl.com](https://docs.pritunl.com/kb/zero)

[![pritunl](img/logo_code.png)](https://docs.pritunl.com/kb/zero)

## Install from Source

```bash
# Install Required Tools
sudo dnf -y install git-core

sudo rm -rf /usr/local/go
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
echo "9e9b755d63b36acf30c12a9a3fc379243714c1c6d3dd72861da637f336ebb35b go1.25.5.linux-amd64.tar.gz" | sha256sum -c - && sudo tar -C /usr/local -xf go1.25.5.linux-amd64.tar.gz
rm -f go1.25.5.linux-amd64.tar.gz

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

# Build Pritunl Zero (update with latest version from releases)
go install -v github.com/pritunl/pritunl-zero@1.0.3648.46
go install -v github.com/pritunl/pritunl-zero/redirect@1.0.3648.46

# Install Systemd Units
sudo cp $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/tools/pritunl-zero.service /etc/systemd/system/
sudo cp $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/tools/pritunl-zero-redirect.socket /etc/systemd/system/
sudo cp $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/tools/pritunl-zero-redirect.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo useradd -r -s /sbin/nologin -c 'Pritunl web server' pritunl-zero-web

# Install Pritunl Zero
sudo mkdir -p /usr/share/pritunl-zero/www/
sudo cp -r $(ls -d ~/go/pkg/mod/github.com/pritunl/pritunl-zero@v* | sort -V | tail -n 1)/www/dist/. /usr/share/pritunl-zero/www/
sudo cp ~/go/bin/pritunl-zero /usr/bin/pritunl-zero
sudo cp ~/go/bin/redirect /usr/bin/pritunl-zero-redirect
sudo systemctl enable --now pritunl-zero
```

## License

Please refer to the [`LICENSE`](LICENSE) file for a copy of the license.
