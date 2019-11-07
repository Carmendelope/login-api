# login-api
​
This component exposes a gRPC API to introduce basic credentials and get a JWT authenticated token with the roles asigned
to the specified user.
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
In order to build and compile this repository use the provided Makefile:
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
Tests are executed using Ginkgo. To run all the available tests:
​
```
make test
```
​
### Update dependencies
​
Dependencies are managed using Godep. For an automatic dependencies download use:
​
```
make dep
```
​
In order to have all dependencies up-to-date run:
​
```
dep ensure -update -v
```
​
## User client interface
Explain the main features for the user client interface. Explaining the whole
CLI is never required. If you consider relevant to explain certain aspects of
this client, please provided the users with them.
​
Ignore this entry if it does not apply.
​
## Known Issues
​
Explain any relevant issues that may affect this repo.
​
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
See also the list of [contributors](https://github.com/nalej/grpc-utils/contributors) who participated in this project.
​
## License

This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
