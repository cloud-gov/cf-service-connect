# CONTRIBUTING

## Welcome

We're so glad you're thinking about contributing to an 18F open source project! If you're unsure about anything, just ask -- or submit the issue or pull request anyway. The worst that can happen is you'll be politely asked to change something. We love all friendly contributions.

We want to ensure a welcoming environment for all of our projects. Our staff follow the [18F Code of Conduct](https://github.com/18F/code-of-conduct/blob/master/code-of-conduct.md) and all contributors should do the same.

We encourage you to read this project's CONTRIBUTING policy (you are here), its [LICENSE](LICENSE.md), and its [README](README.md).

If you have any questions or want to read more, check out the [18F Open Source Policy GitHub repository](https://github.com/18f/open-source-policy), or just [shoot us an email](mailto:18f@gsa.gov).

## Public domain

This project is in the public domain within the United States, and
copyright and related rights in the work worldwide are waived through
the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).

All contributions to this project will be released under the CC0
dedication. By submitting a pull request, you are agreeing to comply
with this waiver of copyright interest.

## Development

Requires Golang. After modifying the source, run

```sh
./run.sh <app_name> <service_instance_name>
```

This will (re)install then run the plugin, all in one.

### Releasing

1. Update `Version` in [`main.go`](main.go).
1. Run `bin/create-release-binaries.sh` to [create cross-compiled binaries](https://github.com/cloudfoundry-incubator/cli-plugin-repo#cross-compile-to-the-3-different-operating-systems).
1. Commit, tag, and push via Git.
1. Upload the binaries to [the new Release](https://github.com/18F/cf-service-connect/releases).
