# rebed
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
