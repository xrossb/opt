# opt

[![Go Reference](https://pkg.go.dev/badge/github.com/xrossb/opt.svg)](https://pkg.go.dev/github.com/xrossb/opt)
[![build](https://github.com/xrossb/opt/actions/workflows/build.yml/badge.svg)](https://github.com/xrossb/opt/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xrossb/opt/branch/main/graph/badge.svg?token=4T3RMEZA7U)](https://codecov.io/gh/xrossb/opt)

Optional values, sans pointers.

**While the library may be functional, this repo is more of an exercise in writing and publishing a Go module.**

## Install

```sh
go get github.com/xrossb/opt@latest
```

## Usage

Make some optional values:

```go
import "github.com/xrossb/opt"

type UpdateUser struct {
    SetEmail    opt.Opt[string]
    SetPassword opt.Opt[string]
}

update := UpdateUser{
    SetPassword: opt.New("use_something_stronger"),
}
```

Zero-values are empty optional values.

```go
if email, ok := update.SetEmail.Get(); ok {
    // we didn't set this above, so this is skipped.
}

if password, ok := update.SetPassword.Get(); ok {
    // do something with the value.
}
```

Quickly switch back to pointers for serialisation:

```go
type UpdateUserReq struct {
    SetEmail    *string
    SetPassword *string
}

func NewUpdateUserReq(update UpdateUser) UpdateUserReq {
    return UpdateUserReq{
        SetEmail:    update.SetEmail.Ptr(),
        SetPassword: update.SetPassword.Ptr(),
    }
}

func (r UpdateUserReq) ToModel() UpdateUser {
    return UpdateUser{
        SetEmail:    opt.Of(r.SetEmail),
        SetPassword: opt.Of(r.SetPassword),
    }
}
```

Avoid unwanted referencing when copying:

```go
type DeepOpt struct {
    Inner opt.Opt[InnerOpt]
}

type InnerOpt struct {
    Message string
}

original := DeepOpt{
    Inner: New(InnerOpt{}),
}
copy := original
copy.Inner.Value.Message = ":)"

// original.Inner.Value.Message == ""
```

As opposed to using pointers:

```go
type Deep struct {
    Inner *Inner
}

type Inner struct {
    Message string
}

original := Deep{
    Inner: &Inner{},
}
copy := original
copy.Inner.Message = ":("

// original.Inner.Message == ":("
```
