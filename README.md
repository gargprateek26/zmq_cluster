ZeroMQ Cluster
=============
zmq_cluster is yet another implementation of a cluster of servers communicating using the well known ZeroMQ sockets. 
Code maturity is considered experimental.


Installation
------------

Use `go get github.com/gargprateek26/zmq_cluster`.  Or alternatively, download or clone the repository.

Usage
-----

The Demo Directory contains the files corresponding to the implementation of a single server. It may be replicated across various terminals with different ServerIds to simulate a cluster. 

Testing
-------
The `cluster_test.go` file shows usage examples and also, performs testing of the implementation. The following test cases have been considered:

1. Broadcast of multiple message from one server to all it's peers
2. Unicast of multiple messages from one server to one of it's peers
3. Broadcast of multiple messages from a randomly chosen server to all it's peers
4. Broadcast of Large Size messages from one server to all it's peers.


Maintainer
----------
Prateek Garg ( garg_prateek26[AT]yahoo[DOT]com )

Coding Style
------------
The source code is automatically formatted to follow `go fmt` by the [IDE]
(https://code.google.com/p/liteide/).  And where pragmatic, the source code
follows this general [coding style]
(http://slamet.neocities.org/coding-style.html).

