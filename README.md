# pritunl-zero: zero trust system

[![github](https://img.shields.io/badge/github-pritunl-11bdc2.svg?style=flat)](https://github.com/pritunl)
[![twitter](https://img.shields.io/badge/twitter-pritunl-55acee.svg?style=flat)](https://twitter.com/pritunl)
[![medium](https://img.shields.io/badge/medium-pritunl-b32b2b.svg?style=flat)](https://pritunl.medium.com)
[![forum](https://img.shields.io/badge/discussion-forum-ffffff.svg?style=flat)](https://forum.pritunl.com)

[Pritunl-Zero](https://zero.pritunl.com) is a [zero trust](https://techster.wiki/zero-trust-security-model) system
that provides secure authenticated access to internal services from untrusted
networks without the use of a VPN. Documentation and more information can be
found at [docs.pritunl.com](https://docs.pritunl.com/docs/pritunl-zero)

[![pritunl](img/logo_code.png)](https://docs.pritunl.com/docs/pritunl-zero)

## Run from Source

```bash
# Install Go
sudo yum -y install git

wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
echo "e85278e98f57cdb150fe8409e6e5df5343ecb13cebf03a5d5ff12bd55a80264f go1.18.linux-amd64.tar.gz" | sha256sum -c -

sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xf go1.18.linux-amd64.tar.gz
rm -f go1.18.linux-amd64.tar.gz

tee -a ~/.bashrc << EOF
export GO111MODULE=on
export GOPATH=\$HOME/go
export PATH=/usr/local/go/bin:\$PATH
EOF
source ~/.bashrc

# Install MongoDB
sudo tee /etc/yum.repos.d/mongodb-org-5.0.repo << EOF
[mongodb-org-5.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/8/mongodb-org/5.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-5.0.asc
EOF

sudo yum -y install mongodb-org
sudo service mongod start

# Install Pritunl Zero
go get -u github.com/pritunl/pritunl-zero

# Run Pritunl Zero (must be run from source directory)
cd ~/go/src/github.com/pritunl/pritunl-zero
sudo ~/go/bin/pritunl-zero start
```

## License

Please refer to the [`LICENSE`](LICENSE) file for a copy of the license.
