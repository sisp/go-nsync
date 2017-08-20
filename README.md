# Named synchronization primitives for Golang

``NamedRWMutex`` allows to acquire read-write mutexes based on a name (string). A possible used case is to dynamically acquire and release mutexes for files.

Comments, suggestions, and contributions are welcome!

## Tests

Running tests requires [``github.com/stretchr/testify``](https://github.com/stretchr/testify) which can be installed as follows:

```bash
$ go get -u github.com/stretchr/testify
```

To execute the tests and check code coverage, run the following commands:

```bash
$ go test -v --race
$ go test -cover fmt
```
