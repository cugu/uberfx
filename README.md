> [!WARNING]  
> Experimental software. Use at your own risk.

# uberfx

uberfx is a command line tool to build and deploy serverless Go applications to [uberspace](https://uberspace.de/).

The uberfx-server works similar to other FaaS (Function as a Service) providers
like AWS Lambda or Google Cloud Functions.

## Features

- 📦 Build serverless Go applications into WASI modules
- 🧪 [Run applications locally for testing](https://docs.fx.uber.space/test-locally.html)
- ✨ [Bootstrap new uberfx projects](https://docs.fx.uber.space/uberfx-cli/init.html)
- 🗃️ [MySQL support](https://docs.fx.uber.space/examples.html#posts)

## Installation

Installation instructions can be found in the [docs](https://docs.fx.uber.space/install-uberfx.html).

## Quickstart

The [quickstart guide](https://docs.fx.uber.space/quickstart.html) describes
how to build and deploy a simple application to an uberspace.

## Examples

There are two example applications available:

- A static website: https://github.com/cugu/uberfx-example-docs
- A simple Go server with MySQL connection: https://github.com/cugu/uberfx-example-posts
