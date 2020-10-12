# Rocket Pool - Smart Node Package - Proof of concept for Raspberry Pi 4


---

Fork of the rocketpool x86_64 version, see: [rocketpool smartnode](https://github.com/rocket-pool/smartnode). Proof of concept which works with raspberry pi 4, running [ubuntu 20.04.1 LTS 64-bit](https://ubuntu.com/download/raspberry-pi). It loads the docker images from docker hub mragicl user. Only compiled for lighthouse! When you configure your rocketpool node, you have to choose infura as eth1 client and lighthouse as eth2 client. I will add the other clients at a later time.


## Installation

Download and run this script: [install.sh](https://raw.githubusercontent.com/mragicl/smartnode/v0.0.4-rpi4/downloads/install.sh) with:
```bash
cd $HOME
wget https://raw.githubusercontent.com/mragicl/smartnode/v0.0.4-rpi4/downloads/install.sh
chmod +x install.sh
./install.sh
```

I adjusted the install script from rocketpool (see: [smartnode-install](https://github.com/rocket-pool/smartnode-install)) to use the appropiate packages for rpi4 (docker-compose), and to configure the docker-compose.yml and config.yml in your .rocketpool to use the docker images compiled for arm64.


Once you ran above install script, open a new shell (or source $HOME/.profile), and you have all rocketpool commands available. You can start immediatly with
```
rocketpool service config
```
Please select infura as eth1 and lighthouse as eth2 client. Then you can start the eth1 and eth2 clients with:
```
rocketpool service start
```
and do all the steps as described here: https://medium.com/rocket-pool/rocket-pool-v2-5-beta-node-operators-guide-77859891766b


## Thanks
Let me know about your impressions and how it works: mragicl@protonmail.com
https://etherscan.io/address/0xcc6AEAA7703C4993B8B8FE6bDe2c4814a5a373B9

