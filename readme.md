## gqlclientgen

GQLClientGen is a Golang client generator for GraphQL APIs. It automatically generates Golang clients for GraphQL services, which can be used for interacting with GraphQL APIs in a simple and efficient way.

For the actual GQL Client check this repo [gqlclient](https://github.com/stefanprifti/gqlclient).

### Installation
Run the following command to install `gqlclientgen`.

```go install github.com/stefanprifti/gqlclientgen/cmd/gqlclientgen@latest```

### Usage
In order to use the generator, the project must declare a configuration file named `gqlclientgen.yml`. Here is a [sample project](https://github.com/stefanprifti/gqlclientgen/tree/main/cmd/gqlclientgen/testdata).

Sample config:

```
version: 1

services:
  - name: Countries API
    package: countries
    url: https://countries.trevorblades.com/graphql
    operations:
      root: gql/countries
    client:
      root: pkg/countries   
```    

The `version` field specifies the version of the config file.

The `services` field is a list of the GraphQL services that the generator will process. Each service has a unique `name` and `package` name. The `url` field is the GraphQL endpoint that the generator will use to retrieve the GraphQL schema. The `operations` field is the path to the directory containing the GraphQL queries that will be used to generate the client. The `client` field is the path to the directory where the generated client code will be stored. 

Once the config file is created, the generator can be run using the `gqlclientgen` command. This will generate the Golang client code for the specified GraphQL services.

Based on this configuration the generator will create [this package](https://github.com/stefanprifti/gqlclientgen/tree/main/cmd/gqlclientgen/testdata/pkg/countries).
- `client.go`: This file contains the code of the generated client. It defines queries and mutations as methods of the client.
- `model.go`: This file contains the GoLang equivlent types of GraphQL schema.
- `schema.graphql`: This file contains the GraphQL schema of the API. It defines the types, fields, and operations that are exposed to the client via the GraphQL API.
- `schema.introspect.json`: This file contains the introspection query result for the GraphQL schema. It can be used to generate client-side code for the GraphQL API.

The generated package can be imported and used in any GoLang application.
