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
