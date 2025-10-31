# pritunl-zero: package

Build test package

## install pacur

```bash
sudo dnf -y install git-core podman

sudo rm -rf /usr/local/go
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
echo "7716a0d940a0f6ae8e1f3b3f4f36299dc53e31b16840dbd171254312c41ca12e go1.25.1.linux-amd64.tar.gz" | sha256sum -c -

sudo tar -C /usr/local -xf go1.25.1.linux-amd64.tar.gz
rm -f go1.25.1.linux-amd64.tar.gz

tee -a ~/.bashrc << EOF
export GO111MODULE=on
export GOPATH=\$HOME/go
export GOROOT=/usr/local/go
export PATH=/usr/local/go/bin:\$PATH:\$HOME/go/bin
EOF
chown cloud:cloud ~/.bashrc
source ~/.bashrc

go install github.com/pacur/pacur@latest
cd "$(ls -d ~/go/pkg/mod/github.com/pacur/pacur@*/podman/ | sort -V | tail -n 1)"
find . -maxdepth 1 -type d -name "*" ! -name "." ! -name ".." ! -name "oraclelinux-10" -exec rm -rf {} +
sh clean.sh
sh build.sh
```

## build package

```bash
git clone github.com/pritunl/pritunl-zero
cd pritunl-zero/tools/package

sudo podman run --rm -t -v `pwd`:/pacur:Z localhost/pacur/oraclelinux-10
```
