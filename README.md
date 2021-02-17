# rebed
[![codecov](https://codecov.io/gh/soypat/rebed/branch/main/graph/badge.svg)](https://codecov.io/gh/soypat/rebed)
[![Build Status](https://travis-ci.org/soypat/rebed.svg?branch=main)](https://travis-ci.org/soypat/rebed)
[![Go Report Card](https://goreportcard.com/badge/github.com/soypat/rebed)](https://goreportcard.com/report/github.com/soypat/rebed)
[![go.dev reference](https://pkg.go.dev/badge/github.com/soypat/rebed)](https://pkg.go.dev/github.com/soypat/rebed)

Recreate embedded filesystems from embed.FS type

## Four actions available:

```go

//go:embed someFS/*
var bdFS embed.FS

// Just replicate folder Structure
rebed.Tree(bdFS)

// Make empty files
rebed.Touch(bdFS)

// Recreate entire FS
rebed.Create(bdFS)

// Recreate FS without modifying existing files
rebed.Patch(bdFS)
```
