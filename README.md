# CockroachSQL Resource Provider

The CockroachSQL Resource Provider lets you manage [CockroachDB](https://www.cockroachlabs.com/) SQL resources.

## Installing

This package is available for several languages/platforms:

### Node.js (JavaScript/TypeScript)

To use from JavaScript or TypeScript in Node.js, install using either `npm`:

```bash
npm install @nellisauction/pulumi-cockroachsql
```

or `yarn`:

```bash
yarn add @nellisauction/pulumi-cockroachsql
```

### Python

To use from Python, install using `pip`:

```bash
pip install pulumi_cockroachsql
```

### Go

To use from Go, use `go get` to grab the latest version of the library:

```bash
go get github.com/nellisauction/pulumi-cockroachsql/sdk/go/...
```

### .NET

To use from .NET, install using `dotnet add package`:

```bash
dotnet add package Pulumi.Cockroachsql
```

## Configuration

The following configuration points are available for the `cockroachsql` provider:

- `cockroachsql:host` (environment: `COCKROACH_HOST`) - CockroachDB host
- `cockroachsql:port` (environment: `COCKROACH_PORT`) - port (default: `26257`)
- `cockroachsql:database` (environment: `COCKROACH_DATABASE`) - database name (default: `defaultdb`)
- `cockroachsql:username` (environment: `COCKROACH_USER`) - username (default: `root`)
- `cockroachsql:password` (environment: `COCKROACH_PASSWORD`) - password
- `cockroachsql:url` (environment: `COCKROACH_URL`) - connection URL; overrides all other connection parameters
- `cockroachsql:sslmode` (environment: `COCKROACH_SSLMODE`) - SSL mode (default: `require`)
- `cockroachsql:connectTimeout` (environment: `COCKROACH_CONNECT_TIMEOUT`) - connection timeout in seconds (default: `180`)

## Reference

For detailed reference documentation, please visit [the Pulumi registry](https://www.pulumi.com/registry/packages/cockroachsql/api-docs/).
