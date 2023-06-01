module github.com/blizzy78/varnamelen

go 1.19

require (
	github.com/matryer/is v1.4.1
	golang.org/x/tools v0.9.2
)

require (
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
)

retract v0.6.1 // see https://github.com/blizzy78/varnamelen/issues/13, use 0.6.2 or later instead
