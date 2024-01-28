var secret password {
  name = "password"
}

build go docs {
  path = "."
}

deploy uberspace docs {
  source   = build.go.docs.output
  username = "fx"
  password = var.secret.password.value
  address  = "tucana.uberspace.de:22"
  domain   = "docs.fx.uber.space"
}
