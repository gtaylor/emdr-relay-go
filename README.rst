emdr-relay-go
=============

:Status: Experimentation
:Author: Greg Taylor
:License: BSD

This is an EMDR_ gateway written in Go_. Resource consumption is markedly
lower compared to our Python relay. 

.. _Go: http://golang.org/
.. _EMDR: http://readthedocs.org/docs/eve-market-data-relay/

Install
-------

* Install Go_. If you are on Debian or Ubuntu, you can ``sudo apt-get intall golang``
* Install a recent zeromq 2.x. ZeroMQ 3.x may or may not work, so it's probably best not to use it just yet.
* Install uuid-dev, libtool, and mercurial (Debian/Ubuntu package names)
* ``sudo go get github.com/alecthomas/gozmq``
* ``sudo go get code.google.com/p/vitess/go/cache``
* From within your ``emdr-relay-go`` dir: ``go build emdr-relay-go.go``
* You should now be able to run the relay: ``./emdr-relay-go``

.. note:: You will need to send an email to gtaylor (at) gc-taylor (dot) 
	com before your relay will be allowed to connect to the announcers.

License
-------

This project, and all contributed code, are licensed under the BSD License.
A copy of the BSD License may be found in the repository.
