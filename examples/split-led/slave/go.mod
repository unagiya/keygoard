module github.com/unagiya/keygoard/examples/split-led/slave

go 1.25.6

require github.com/unagiya/keygoard v0.0.0

require (
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	tinygo.org/x/drivers v0.34.0 // indirect
)

replace github.com/unagiya/keygoard => ../../..
