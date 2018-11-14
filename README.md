# Tool for making database migrations in golang

[![Build Status](https://travis-ci.org/thenixan/migrater.svg?branch=master)](https://travis-ci.org/thenixan/migrater) 
[![codecov](https://codecov.io/gh/thenixan/migrater/branch/master/graph/badge.svg)](https://codecov.io/gh/thenixan/migrater)
[![Go Report Card](https://goreportcard.com/badge/github.com/thenixan/migrater)](https://goreportcard.com/report/github.com/thenixan/migrater)

#### Installation
```
go get -u github.com/thenixan/migrater
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