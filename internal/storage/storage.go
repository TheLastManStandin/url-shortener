package storage

import "errors"

var ErrAliasExists = errors.New("alias already exists")
var ErrAliasNotFound = errors.New("alias not found")
