module stone

go 1.13

replace stone => ./

replace go-test-2/errorMyself => /Users/tonny/source/go/src/go-test-2/errorMyself

require (
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	go-test-2/errorMyself v0.0.0-00010101000000-000000000000
	gotest.tools v2.2.0+incompatible
)
