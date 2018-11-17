# Tool for making database migrations in golang

[![Build Status](https://travis-ci.org/blackwood-apps/migrater.svg?branch=master)](https://travis-ci.org/blackwood-apps/migrater) 
[![codecov](https://codecov.io/gh/blackwood-apps/migrater/branch/master/graph/badge.svg)](https://codecov.io/gh/blackwood-apps/migrater)
[![Go Report Card](https://goreportcard.com/badge/github.com/blackwood-apps/migrater)](https://goreportcard.com/report/github.com/blackwood-apps/migrater)

#### Installation
```
go get -u github.com/blackwood-apps/migrater
```

#### Usage

First you should declare all your migrations
```gotemplate
s := Set{
	1: Step{
		Up: []string{
			"CREATE TABLE test (id int PRIMARY KEY);",
		},
		Down: []string{
			"DROP TABLE test;",
		},
	},
	2: Step{
		Up: []string{
			"CREATE TABLE test2 (id int PRIMARY KEY);",
		},
		Down: []string{
			"DROP TABLE test2;",
		},
	},
}
```

Then
```gotemplate
from, to, err = s.Upgrade(db)
if err != nil {
	log.Fatal(err)
}
```

`from` will have the version migration was made from
`to` will have the new version
