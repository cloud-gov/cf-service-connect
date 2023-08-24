# Cloud Foundry CLI Service Connection Plugin [![Code Climate](https://codeclimate.com/github/18F/cf-service-connect/badges/gpa.svg)](https://codeclimate.com/github/18F/cf-service-connect)

This plugin makes it easy to connect to your databases or other Cloud Foundry service instances from your local machine. This condenses the steps listed in [Accessing Services with SSH](https://docs.cloudfoundry.org/devguide/deploy-apps/ssh-services.html) to a single command.

![demo screencast](demo.gif)

Requires Diego architecture with [SSH enabled](https://docs.cloudfoundry.org/running/config-ssh.html).

## Support

Currently supports (most) service brokers for the following:

* MongoDB (requires [`mongo` shell](https://docs.mongodb.com/getting-started/shell/installation/))
* MySQL (requires [`mysql` CLI](https://dev.mysql.com/doc/refman/8.0/en/installing.html))
* PostgreSQL (requires [`psql` CLI](https://postgresapp.com/documentation/cli-tools.html))
* Redis (requires [`redis-cli`](https://redis.io/topics/quickstart))

## Local installation

1. Install the [Cloud Foundry CLI](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html) v6.15.0 or later.
2. Install this plugin, using the appropriate binary URL from [the Releases page](https://github.com/18F/cf-service-connect/releases).

    ```sh
    cf install-plugin <binary_url>
    # will be of the format
    # https://github.com/cloud-gov/cf-service-connect/releases/download/<version>/cf-service-connect_<os>-<arch>
    # For non-M1 Macs, use `cf-service-connect_darwin_amd64`
    # For M1 Macs, use `cf-service-connect_darwin_arm64`
    ```

3. Install the CLI corresponding to your service type (see above).

## Usage

> **Note**
> If you are using this tool to connect to a service on cloud.gov, your space must be configured with the `trusted_local_networks_egress` security group. Do this by running `cf bind-security-group trusted_local_networks_egress ORG --space SPACE` with your organization and space. Skipping this step will result in a `connection refused` error. For more, see [cloud.gov: Controlling egress traffic](https://cloud.gov/docs/management/space-egress/).

* `app_name` is the name of the app in your space you want to tunnel through.
* `service_instance_name` is the service instance you wish to connect to.

```shell
$ cf target --organization <org> --space <space>
$ cf connect-to-service <app_name> <service_instance_name>
Finding the service instance details...
Setting up SSH tunnel...
...
mysql>
```

If you get an error such as "connection refused", "error opening SSH connection", or "psql: could not connect to server: Connection refused" this is usually caused by being on a network that blocks the SSH port that this tool is trying to use. Try using a different network, or consider asking your network administrator to unblock the port (typically 22 and/or 2222).

### Optional: overriding `cf` CLI binary name

If you are in Windows or another environment where the Cloud Foundry CLI was installed as `cf7` or `cf8`, you can set an environment to tell the plugin what binary name to use for the Cloud Foundry CLI:

```shell
CF_BINARY_NAME=cf7 cf connect-to-service <app_name> <service_instance_name>
```

### Manual client connection

If you're using a non-default client (such as a GUI), run with the `-no-client` option to set up your client connection on your own.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)
