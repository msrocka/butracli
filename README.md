# butracli

This is a small CLI tool for querying the
[BuildingTransparency API](https://etl-api.cqd.io/). Simply start the REPL
using your user crendentials:

```
$ butracli -u USER -p PASSWORD
```

Currently, only `GET` get requests are supported. You just enter the last part
of the request URL, e.g.:

```
##> GET https://etl-api.cqd.io/api/ orgs?page_size=1
```

See the [API documentation](https://etl-api.cqd.io/) for the request details.
Entering `q`, `quit`, `exit`, or `halt` will logout and stop the REPL.
