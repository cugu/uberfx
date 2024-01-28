# uberfx

uberfx is a framework for building and deploying applications on https://uberspace.de/.

## Installation

`uberfx` is a Go program. You can install it with:

```shell
go install github.com/cugu/uberfx
```

## Usage

### Create a new project

New projects can be created with the `init` command:

```shell
uberfx init myproject
```

This will create a new directory `myproject` with a basic project structure.

### Configure the project

The project can be configured in the `uberfx.hcl` file, which is created by the `init` command. 
The following example shows the basic configuration:

```hcl
# An input variable, which can be set with the
# --var flag, e.g. uberfx deploy --var 'password=1234567890'
# or by setting the environment variable UBERFX_VAR_password
var secret password {
  name = "password"
}

# A build step, which builds the go binary as a wasm module
build go www {
  path = "./server"
}

# A deploy step, which deploys the binary to an uberspace
deploy uberspace www {
  source   = build.go.www.output
  username = "fx"
  password = var.secret.password.value
  address  = "tucana.uberspace.de:22"
  domain   = "www.fx.uber.space"
}
```

This configuration will build the go binary in the `server` directory and deploy it to the uberspace `tucana.uberspace.de` with the username `fx` and the password from the `password` variable. 
The binary will be deployed to the domain `www.fx.uber.space`.

### Build and deploy the project

To build and deploy a project, use the `deploy` command:

```shell
uberfx deploy --var 'password=1234567890'
```

## Disclaimer

This project is not affiliated to uberspace.de, Uber or any other company with "uber" in its name.
