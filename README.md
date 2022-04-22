![Shattered Realms Online](https://github.com/ShatteredRealms/Documentation/raw/main/assets/images/logo/WhiteLogo.png)

# Overview [![Build](https://github.com/ShatteredRealms/Accounts/actions/workflows/build.yml/badge.svg)](https://github.com/ShatteredRealms/Accounts/actions/workflows/build.yml) [![codecov](https://codecov.io/gh/ShatteredRealms/Accounts/branch/main/graph/badge.svg?token=P01UR012I1)](https://codecov.io/gh/ShatteredRealms/Accounts)
The accounts microservice for [Shattered Realms Online](https://github.com/ShatteredRealms/Game). Manages user accounts,
and authentication throughout SRO.

# Development
The `Makefile` is located within the `build` folder within the project root directory. All make commands should be run from there.

## Environments
This project uses environment variables which should be stored within a `.env` file located within the project root direcory. If one is not configured, rename `.env.template` to `.env` and configure the variables for deployment. These variables can be overwritten in the OS, in a docker env file, kubernetes env file, and at runtime by supplying them before the run command.

## Commands
### Building
**Binary:** To build a binary output run `make build`. The output result will be placed in the `bin` folder in the project root directory. \
**Docker:** To build the docker image run `make build-image` and a image called `sro-accounts` will be generated.

### Testing
To run the tests and see the coverage report use `make test`. To view a the HTML results, simply run `make report`.

### Deployment
Deployment is done using docker. If using an AWS docker repository, running `make aws-docker-login` will authenticate with the default aws credential context. To push the images, run `make push`. This will build the image and push them to the docker repository.
