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

* Install Go_.
* ``sudo go get github.com/alecthomas/gozmq``
* ``sudo go get github.com/bradfitz/gomemcache``
* From within your ``emdr-relay-go`` dir: ``go build emdr-relay-go.go``
* You should now be able to run the relay: ``./emdr-relay-go``

.. note:: You will need to send an email to gtaylor (at) gc-taylor (dot) 
	com before your relay will be allowed to connect to the announcers.

License
-------

This project, and all contributed code, are licensed under the BSD License.
A copy of the BSD License may be found in the repository.
