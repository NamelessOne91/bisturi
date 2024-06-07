# bisturi
Poor men's network analyzer

The compiled binary executable for bisturi will attempt to create a new raw socket and bind it to a network interface using syscalls.
This will fail unless the program is run with root privileges, which is not advisable.

You can use [setcap](https://man7.org/linux/man-pages/man8/setcap.8.html) to grant the binary executable *only* the capability to operate on raw sockets.
This is the default behaviour of the included Makefile's **build** command.