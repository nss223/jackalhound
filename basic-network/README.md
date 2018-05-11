### One-peer-testing Network Usage

*. Edit `./fabric.conf` at your need; add chaincode and its channel under var `CHANNEL`.

*. Run `./generate.sh` if new channel is added.

*. To (re-)start the network, run `./restart.sh`;
   you may restart it every time when adding new chaincode or channel is reconfigured.

*. You may want to upgrade the chaincode without restarting the network:
```
./upgrade.sh <chaincode_name> <version>
```
NOTE:
    - the chaincode is located at `/chaincode/<chaincode_name>`, you must build before upgrade
    - the `version` must be different, the init version is `1.0`

*. Test the chaincode using `./query` or `./invoke`

```
./query | invoke <chaincode_name> <function> <parameter>
```

the `parameter` must in the form of `'"par1", "par2", ...'`, i.e.
`./query test myfunction '"par1", "par2", "par3"'`

*. `./teardown.sh` to shutdown.
