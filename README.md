![Simple CI](https://github.com/NamelessOne91/bisturi/actions/workflows/simple_ci.yml/badge.svg)

# bisturi
A toy project network packet analyzer

The compiled binary executable for bisturi will attempt to create a new raw socket and bind it to a network interface using syscalls.
This will fail unless the program is run with root privileges, which is not advisable.

You can use [setcap](https://man7.org/linux/man-pages/man8/setcap.8.html) to grant the binary executable *only* the capability to operate on raw sockets.
This is the default behaviour of the included Makefile's **build** command.


## Usage

You can build the binary executable with the `make build` command.

The following flags are available to customize bisturi's behaviour:

| Flag | Type | Default | Meaning
| :---:|:--:|:--:|:--|
| i | string | eth0 | network interface for which the packets will be analyzed |
| p | string | all  | protocl filter - 'all' equals to no filter |

Running bisturi with the provided `make run` command is functionally equivalent to running bisturi with the following flags, which are its defaults:

`bisturi -i eth0 -p all`

