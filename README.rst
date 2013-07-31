emdr-relay-go
=============

:Status: Stable
:Author: Greg Taylor
:License: BSD

This is an EMDR_ gateway written in Go_. Resource consumption is markedly
lower compared to our Python relay. 

.. _Go: http://golang.org/
.. _EMDR: http://readthedocs.org/docs/eve-market-data-relay/

Install
-------

* Install Go_. If you are on Debian or Ubuntu, you can ``sudo apt-get intall golang``
* Install a recent ZeroMQ 2.x or ZeroMQ 3.x.
* Install uuid-dev, libtool, and mercurial (Debian/Ubuntu package names)
* ``sudo go get github.com/alecthomas/gozmq`` (for ZeroMQ 3.x: ``sudo go get -tags zmq_3_x github.com/alecthomas/gozmq``)
* ``sudo go get code.google.com/p/vitess/go/cache``
* From within your ``emdr-relay-go`` dir: ``go build emdr-relay-go.go``
* You should now be able to run the relay: ``./emdr-relay-go``
* Before we can list you, your relay will need to be running under a process 
  supervisor like Runit, supervisord, systemd, upstart, or something similar.
* You'll also want to restart your relay process daily. Sometimes ZeroMQ gets
  in a funk, or doesn't pick up an upstream DNS change.

.. note:: You will need to send an email to gtaylor (at) gc-taylor (dot) 
	com before your relay will be allowed to connect to the announcers.

License
-------

This project, and all contributed code, are licensed under the BSD License.
A copy of the BSD License may be found in the repository.
