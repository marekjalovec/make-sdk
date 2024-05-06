# Make SDK for Go

Can be used to query and control your Make Scenarios, Connections, Variables, Users, and more.


## Quick start

```shell
go get github.com/marekjalovec/make-sdk
```

```go
  var config, err = makesdk.NewConfig(apiToken, environmentUrl, rateLimit)
  if err != nil {
      return nil, err
  }

  var client = makesdk.GetClient(config)
```

```go
// load a specific resource
var c = client.GetConnection(id)

// or iterate through a list
var clp = client.NewConnectionListPaginator(-1, teamId)
for clp.HasMorePages() {
    connections, err := clp.NextPage()
    if err != nil {
        return nil, err
    }

    for _, i := range connections {
        log.Println(i)
    }
}
```


For more use-cases check:

  - [steampipe-plugin-make](https://github.com/marekjalovec/steampipe-plugin-make)


Further reading:

  - [Make API documentation](https://www.make.com/en/api-documentation)


Get involved:

  - [Issues](https://github.com/marekjalovec/make-sdk/issues)
