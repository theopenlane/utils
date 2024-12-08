// Package contextx is a helper package for managing context values
// Most **request-scoped data** is a singleton per request
// That is, it doesn't make sense  for a request to carry around multiple loggers, users, traces
// you want to carry the _same one_ with you from function call to function call
// the way we've handled this historically is a separate context key per type you want to carry in the struct
// but with generics, instead of having to make a new zero-sized type for every struct
// we can just make a single generic type and use it for everything which is what this helper package is intended to do
package contextx
