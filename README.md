# go-jnlp-launcher
JNLP launcher written in Go, providing a simple command line interface for launching Java applications described by a JNLP descriptor.

## Objective

I started this project because I wanted to learn the Go language and tools and I felt like its self-contained executables were a good fit for a CLI tool like this one.

The goal of this project is to provide a simple CLI for launching Java applications without using conventional JNLP clients such as Web Start. If you frequently use the terminal to get things done, and you are working with different Java applications or versions of the same Java application deployed using JNLP, then you may find this tool useful.

## Building
	$ git clone https://imcdonagh/go-jnlp-launcher
	$ cd go-jnlp-launcher
	$ go build

## Usage
	go-jnlp-launcher <url> <dir> [<args>...]

	url  - JNLP descriptor URL
	dir  - local cache directory where resources will be downloaded
	args - Optional program arguments to pass to the application

Example:

	$ ./go-jnlp-launcher https://example.com/someapp.jnlp arg1 arg2

## Limitations
The launcher does not implement the full JNLP spec:
- Pack200 is not supported
- Versioned resources are not supported
- Jardiff is not supported
- Lazily downloaded resources are ignored (though this can be changed easily)
- Only applications are supported, while nativelibs, extensions and installers are not.

## Roadmap
The aim of this project is to keep things simple, but here are some ideas for improving the tool:
- show download progress
- multi part download for larger files
- allow overriding JVM options, such as -Xmx
- proxy support

Ideas and contributions welcome. Feel free to submit a PR if you have an enhancement or fix.

## License
The code is licensed under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license.
