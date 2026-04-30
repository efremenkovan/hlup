package testdata

// QuerySimple is a plain 4-word expression using and/or only.
var QuerySimple = `Go and language or software`

// QueryMedium is a ~20 word expression with all expression types and 3-4 levels of nesting.
var QueryMedium = `(Go and (concurrency or goroutines or 'garbage collection')) and not ('to hell' or deprecated) and (reliable or efficient or 'open source') or (Google and (Pike or Thompson) and not Java)`

// QueryComplex is a ~40 word expression with 12 levels of nesting and all expression types.
var QueryComplex = `((((Go and (((concurrency or goroutines) and not 'to hell') or ((interfaces or 'error handling') and (testing or benchmarks)))) or (((Google and (Pike or Thompson)) and not (Java or Python)) and ('open source' or community))) and (((reliable or efficient) and not deprecated) or ((simple or straightforward) and ('standard library' or packages)))) and ((('type system' or generics) and not (inheritance or exceptions)) or (((compilation or 'cross compilation') and (Linux or macOS)) or ((performance or optimization) and not 'runtime overhead'))))`
