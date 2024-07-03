![Simple CI](https://github.com/NamelessOne91/bisturi/actions/workflows/simple_ci.yml/badge.svg)

# bisturi
A TUI network packet analyzer toy project

The compiled binary executable for bisturi will attempt to create a new raw socket and bind it to a network interface using syscalls.
This will fail unless the program is run with root privileges, which is not advisable.

You can use [setcap](https://man7.org/linux/man-pages/man8/setcap.8.html) to grant the binary executable *only* the capability to operate on raw sockets.
This is the default behaviour of the included Makefile's **build** command.


## Usage

You can build the binary executable with the `make build` command or build & run it with `make run`.

A [Bubbletea](https://github.com/charmbracelet/bubbletea) based TUI will ask you to select a network interface and a protocol to filter for - selecting 'all' equals to having no filter.
