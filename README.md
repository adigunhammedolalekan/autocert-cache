### autocert-cache
Database backed certs cache for https://github.com/golang/crypto/tree/master/acme/autocert. This package implements `autocert.Cache` interface to provide a database backed cache. It works well when app is deployed in a container and you want LetsEncrypt certificates to be around while you deploy your apps steadily.


## Usage
`go get github.com/adigunhammedolalekan/autocert-cache`

Then use `cache.NewDbCache(dbDialect, dbConnectionUrl)` to initialize your `autocert.Manager`, see example below.


```Go
    m := autocert.Manager{
	    Prompt: autocert.AcceptTOS,
	    Cache: cache.NewDbCache(dbDialect, dbConnectionUrl),
	    HostPolicy: autocert.HostWhitelist("your.domain.com"),
    }
```

You can then use `autocert.Manager` however you want, the Database cache would be utilize underneath.

Note: This package is using https://github.com/jinzhu/gorm to interact with the underlying, `dialect` and `connectUri` from `cache.NewDbCache(dialect, connectUri)` could be any string supported by GORM