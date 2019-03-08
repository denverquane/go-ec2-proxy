package create

import (
	"encoding/base64"
	"github.com/denverquane/go-ec2-proxy/common"
)

func CreateGoProxyScriptString(proxyConfig common.ProxyConfig) string {
	command := ""

	if proxyConfig.User != "" {
		command = proxyConfig.Type + " -t tcp -p '0.0.0.0:" + proxyConfig.Port + "' -a '" + proxyConfig.User + ":" + proxyConfig.Pass + "'"
	} else {
		command = proxyConfig.Type + " -t tcp -p '0.0.0.0:" + proxyConfig.Port + "'"
	}

	cmdString := `#!/bin/bash
yum update -y
curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto.sh | sudo bash

sudo echo "[Unit]
Description=Proxy server
Requires=network.target
[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu
ExecStart=/bin/bash -lc '/usr/bin/proxy ` + command +
		`'
TimeoutSec=15
Restart=always
[Install]
WantedBy=multi-user.target" > /etc/systemd/system/proxy.service

# Enable service
sudo systemctl daemon-reload
sudo systemctl enable proxy.service
sudo systemctl start proxy.service

# Print status
sudo systemctl status proxy.service --no-pager`

	//fmt.Println(cmdString)

	return base64.StdEncoding.EncodeToString([]byte(cmdString))
}
