# Entando Upgrade CLI

Currently this app depends on some packages of [upgrade-operator](https://github.com/entgigi/upgrade-operator), that are referenced in the file go.mod using a relative file path (`../upgrade-operator`). Be sure to have put the two directories at the same level.

## Environment variables

Following environment variables must be set:

* `ENTANDO_KUBECTL_BASE_COMMAND`: base `kubectl` command; for testing purposes it can be set to `kubectl -n entando`
* `ENTANDO_APPNAME`: name of the Entando app

These variable will be passed to the app by the `ent` wrapper.
