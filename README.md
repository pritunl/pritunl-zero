# pritunl-zero: zero trust system

[![github](https://img.shields.io/badge/github-pritunl-11bdc2.svg?style=flat)](https://github.com/pritunl)
[![twitter](https://img.shields.io/badge/twitter-pritunl-55acee.svg?style=flat)](https://twitter.com/pritunl)

[Pritunl-Zero](https://zero.pritunl.com) is a zero trust system
that provides secure authenticated access to internal services from untrusted
networks without the use of a VPN. Documentation and more information can be
found at [docs.pritunl.com](https://docs.pritunl.com/docs/pritunl-zero)

[![pritunl](img/logo_code.png)](https://docs.pritunl.com/docs/pritunl-zero)

## Run from Source

```bash
# Install Go
sudo yum -y install git
curl -L https://golang.org/dl/go1.15.6.linux-amd64.tar.gz | sudo tar -C /usr/local -xz
tee -a ~/.bashrc << EOF
export GOPATH=\$HOME/go
export PATH=/usr/local/go/bin:\$PATH
EOF
source ~/.bashrc

# Install MongoDB
sudo tee /etc/yum.repos.d/mongodb-org-4.2.repo << EOF
[mongodb-org-4.2]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/7/mongodb-org/4.2/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-4.2.asc
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
