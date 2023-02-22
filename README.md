# Debian 11 Setup
All commands were tested as root on a fresh debian instance

### Install dependencies

1. `apt install git curl postgresql nginx python3-certbot-nginx`
1. `curl -Lo monero.tar.gz https://downloads.getmonero.org/cli/linux64`
1. `tar xvf monero.tar.gz && mv monero-x86*/monero-wallet-rpc /usr/bin && rm -rf monero.tar.gz monero-x86*`

### Create user and database

1. `useradd -m -d /srv/shadowchat -s /bin/bash shadowchat`
1. `su - postgres -c 'createuser -P shadowchat && createdb -l C -E UTF8 -T template0 -O shadowchat shadowchat_db'`

### Upload your view only wallet
Copy your viewonly and viewonly.keys files to the /opt/shadowchat/ directory.
1. `scp "C:\Users\user\Documents\Monero\wallets\user\walletname_viewonly*" root@server_ip:/srv/shadowchat`

### Install and configure shadowchat
1. `git clone https://github.com/strangefru1t/sc.git && cd sc`
1. `cp service-files/*.service /etc/systemd/system && systemctl daemon-reload`
1. `cp -r html/ rpc.conf config.json /srv/shadowchat/`
1. `chown -R shadowchat:shadowchat /srv/shadowchat`
1. `wget https://go.dev/dl/go1.20.1.linux-amd64.tar.gz`
1. `tar xvf go1.20.1.linux-amd64.tar.gz`
1. `GOBIN=/usr/bin ./go/bin/go install cmd/sc-api/shadowchat.go`

#### Edit /srv/shadowchat/rpc.conf
#### Edit /srv/shadowchat/config.json

1. `systemctl enable --now shadowchat xmr-rpc`

#### /etc/nginx/sites-enabled/default
    server {
    
        server_name example.com;
        
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        add_header X-Frame-Options DENY;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP  $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
    
        location / {
            proxy_pass http://127.0.0.1:8000;
        }
    }
