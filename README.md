# htmlquery

Command line tool for parsing html with xpath.

This simple program reads html from stdin and parse with the given xpath expression and print it to stdout.

[github.com/antchfx/xpath](https://github.com/antchfx/xpath) package is used for xpath.

# Installation
```
go get -u github.com/pubblic/htmlquery
```

# Examples
```
$ htmlquery "//span[contains(@class, 'tag1') or contains(@class, 'tag2')"
```

```
$ curl -s "https://google.com/" | htmlquery "//h1"
<h1>301 Moved</h1>
```

```
$ curl -s "https://google.com/" | htmlquery "///h1/text()"
301 Moved
```

