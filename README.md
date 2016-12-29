# Cloud Foundry CLI Database Connection Plugin

This plugin makes it easy to connect to your databases in Cloud Foundry from your local machine. This condenses the steps listed in [Accessing Services with SSH](https://docs.cloudfoundry.org/devguide/deploy-apps/ssh-services.html) to a single command. Requires Diego architecture with [SSH enabled](https://docs.cloudfoundry.org/running/config-ssh.html).

## Local installation

1. Install the Cloud Foundry CLI v6.15.0 or later.
1. Install the plugin. [TODO]
1. Install `psql` or `mysql` (depending on which database you need to connect to).

## Usage

```
$ cf target --organization <org> --space <space>
$ cf connect-to-db <app_name> <service_instance_name>
Finding the service instance details...
Setting up SSH tunnel...
Initiating connection...
mysql>
```

## Development

```sh
./run.sh
```
