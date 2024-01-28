> [!WARNING]  
> Experimental software. Use at your own risk.

# uberfx

uberfx is a framework for building and deploying applications on https://uberspace.de/.
It utilizes the [uberfx-server](https://github.com/cugu/uberfx-server) to host the applications.

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

service uberspace_mysql mysql {
  username = "fx"
  password = var.secret.password.value
  address  = "tucana.uberspace.de:22"
}

# A build step, which builds the go binary as a wasm module
build go www {
  path = "."
}

# A deploy step, which deploys the binary to an uberspace
deploy uberspace www {
  source   = build.go.www.output
  username = "fx"
  password = var.secret.password.value
  address  = "tucana.uberspace.de:22"
  domain   = "www.fx.uber.space"
  env      = {
    "MYSQL_PASSWORD" = service.uberspace_mysql.mysql.password
  }
}
```

This configuration will build the go binary in the `server` directory and deploy it to the uberspace `tucana.uberspace.de` with the username `fx` and the password from the `password` variable. 
The binary will be deployed to the domain `www.fx.uber.space`.

### Run locally

To run the project locally, use the `go run` command:

```shell
go run . localhost:8080
```

The first argument is the address to listen on, in this case `http://localhost:8080`.

### Install the uberfx-server on the uberspace

To install the `uberfx-server` on the uberspace, use the `install` command:

```shell
uberfx install --address "tucana.uberspace.de:22" --username "fx" --password "$PASSWORD"
```

### Build and deploy the project

To build and deploy a project, use the `deploy` command:

```shell
uberfx deploy --var password=foobar
```
