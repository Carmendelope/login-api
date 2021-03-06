# login-api
​
This component exposes a gRPC API to introduce basic credentials and get a JWT authenticated token with the roles assigned
to the specified user.

Notice that the Login API supports REST and gRPC request by means of the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
that is launched when the server starts.
​
## Getting Started
​
This component receives login requests crafted with `grpc_authx_go.LoginWithBasicCredentialRequest` and answers with
`grpc_authx_go.LoginResponse`.
​
### Prerequisites
​
To run this component you should have deployed at least the following components:
​
* [authx](https://github.com/nalej/authx)
​
### Build and compile
​
To build and compile this repository use the provided Makefile:
​
```
make all
```
​
This operation generates the binaries for this repo, download dependencies,
run existing tests and generate ready-to-deploy Kubernetes files.
​
### Run tests
​
Tests are executed using Ginkgo. To run all the available unit tests:
​
```
make test
```
​
#### Integration tests

There are integration tests that need some configuration and appliance to be ready to run:

* Have a running and accessible instance of [authx](https://github.com/nalej/authx).
* Set the following environment variables:

 | Variable             | Example Value  | Description           |
 | -------------------- | -------------- | --------------------- |
 | RUN_INTEGRATION_TEST | true           | Run integration tests |
 | IT_AUTHX_ADDRESS     | localhost:8800 | Authx Address         |

### Update dependencies
​
Dependencies are managed using Godep. For an automatic dependencies download use:
​
```
make dep
```
​
To have all dependencies up-to-date run:
​
```
dep ensure -update -v
```
​
## User client interface

To interact with this component, you can leverage the [public-api-cli](https://github.com/nalej/public-api).
The command that interacts with this component is `login`.

Example:
```shell script
./public-api-cli login --email example@nalej.com --password change_me
```
​
## Known Issues
​
## Contributing
​
Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.
​
## Versioning
​
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/nalej/login-api/tags). 
​
## Authors
​
See also the list of [contributors](https://github.com/nalej/login-api/contributors) who participated in this project.
​
## License

This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
