emdr-relay-go
=============

:Version: 1.0
:Status: Stable
:Author: Greg Taylor
:License: BSD

This is the EMDR_ relay, written in Go_. If you are wanting to run a relay
of your own (or contribute changes), this is the place!

Install (in a Docker container)
-------------------------------

This is the recommended installation method. It should be a quick and easy
process for most.

* Install Docker_.
* Head on over to the `Docker Hub repo`_ and read the instructions.

Install (directly on your machine)
----------------------------------

* See our `Direct Install`_ instructions.

Changes
-------

1.0
^^^

* **Now requires ZeroMQ 4.x.**
* Changed ZeroMQ_ bindings to pebbe/zmq4.
* Dockerized release.

0.1
^^^

* Initial experimental release.


License
-------

This project, and all contributed code, are licensed under the BSD License.
A copy of the BSD License may be found in the repository.

.. _ZeroMQ: http://zeromq.org/
.. _Go: http://golang.org/
.. _EMDR: http://readthedocs.org/docs/eve-market-data-relay/
.. _Docker: https://docs.docker.com/installation/#installation
.. _Docker Hub repo: https://registry.hub.docker.com/u/gtaylor/emdr-relay-go/
.. _Direct Install: https://github.com/gtaylor/emdr-relay-go/wiki/Direct-Installation-Instructions
