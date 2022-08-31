# boltz-poc

# 0. Requirements
Go (golang) is required [installation instructions](https://go.dev/doc/install)

Recent version of docker is also needed.
## Boltz/RSK Docker images
 - Clone Boltz repo: https://github.com/rootstock/boltz-backend
 - checkout the `rsk_integration` branch
 - install the dependencies with `npm install` (e.g this will pull `boltz-core`)
 - Copy contracts artifacts and hardhat config

```
mkdir -p contracts && cp -R ./hardhat.config.ts node_modules/boltz-core/scripts node_modules/boltz-core/artifacts contracts
```
 - Build Boltz image -> Run: `npm run docker:build`
 - Build RSK node image -> Run: `npm run docker:rskj:build`

Note: On M1 (Apple silicon) mac computers, there may be some npm installation errors (e.g. for grpc-tools). Since Rosetta is now included in Mac OS, it helps to run docker in "x86 mode". This can be done via setting environment variable to x86 linux (e.g. add to ~/.zprofile)

```
export DOCKER_DEFAULT_PLATFORM=linux/amd64     
```
# 1. Setup environment

Return to directory of this repo i.e. `boltz-poc`

- Run: `npm run docker:volumes:reset`
- Run: `npm run docker:start`

If there are any errors of the form `[lncli] rpc error: code = Unknown desc = server is still in the process of starting`, then wait for a minute and then run `npm run docker:connect`

We should see something like

```
Mine a BTC block.
{
  "address": "bcrt1qt5hu3sukzvggmnx0ap6r5mjvsghqryp8w63kqg",
  "blocks": [
    "32b362fc4d1b76a4112e0b919ab23ffcc3f7f8e165506c71509c92d57a896de1"
  ]
}
Connect LN-BTC nodes.
Connect LN-BTC alice->bob.
{

}
Connect LN-BTC alice->carol.
{

}
Connect LN-BTC bob->carol.
{

}
```

## Start boltz daemon 
Starting boltzd manually to avoid error on M1 where container keeps exiting.

Ideally, the `boltzd` service should be an entrypoint for the container. But on M1 macs, the container kept exiting. 

The following should be handled via a new npm run call by using `docker compose` to execute the call from the context of the container (without having to gain direct access). 

```
# change containerID with name
docker exec -it boltz-backend bash

#then from within the container
root@e71c65b55701:/# cd boltz-backend
root@e71c65b55701:/boltz-backend# ./bin/boltzd --configpath ./config.toml
```


# 2. Start server
- Run: `npm run server`

# 3. Start client

- Run `npm install` in the client directory for dependendencies

- Run: `npm run client`

- Go to: `http://localhost:3000` (or 3001?)

# 4. Stop environment
- Run: `npm run docker:stop`

# X. Miscellaneous

### Generate contract.go
- install protoc ([here](https://grpc.io/docs/protoc-installation/))
- install solc ([here](https://docs.soliditylang.org/en/v0.8.9/installing-solidity.html) and [here](https://www.educative.io/answers/how-to-install-solidity-in-mac))
- Clone Geth / make / make devtools ([doc](https://goethereumbook.org/smart-contract-compile/))
- `abigen --abi=server/abi/EtherSwap.json --out=server/abi/EtherSwap.go --pkg=abi`

### Fix common error: https://github.com/lightningnetwork/lnd/issues/1177
* Generate btc blocks

`docker exec -it polar-n1-backend1 bitcoin-cli -chain=regtest -rpcuser=polaruser -rpcpassword=polarpass -generate 1`

### Connect to docker images
* Connect to bitcoin node

`docker exec -it --user lnd polar-n1-backend1 /bin/bash`

* Connect to a LN node

`docker exec -it --user lnd polar-n1-alice /bin/bash`

### Connect btw LN nodes
* `docker exec -it --user lnd polar-n1-alice lncli --network=regtest connect 02045b28b45f0b8efdac9381287b728d7c6897aa8c4d26fa5e9570078dc949e11b@host.docker.internal:9739`
* `docker exec -it --user lnd polar-n1-alice lncli --network=regtest connect 03b45c2032206051af4e130a8901575107bc1441287700ec5a8b6421d70a863a24@host.docker.internal:9740`
* `docker exec -it --user lnd polar-n1-bob lncli --network=regtest connect 03b45c2032206051af4e130a8901575107bc1441287700ec5a8b6421d70a863a24@host.docker.internal:9740`

### LN nodes URLs:

* alice: `02e7dd429d9148b6fde3e026d39ac63d2d76a3d3a5469ad8a64e0310fc6f62f302@host.docker.internal:9735`
* bob: `02045b28b45f0b8efdac9381287b728d7c6897aa8c4d26fa5e9570078dc949e11b@host.docker.internal:9739`
* carol: `03b45c2032206051af4e130a8901575107bc1441287700ec5a8b6421d70a863a24@host.docker.internal:9740`

