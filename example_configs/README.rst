Example Configs
===============

The files in this directory can be used as a starting point for supervising
and running emdr-relay-go. You'll want to change the paths to fit your
deployment, of course.

* supervisord-relay.conf will go in /etc/supervisor/conf.d/ on Ubuntu, or
  directly in /etc/supervisor/supervisor.conf on other distros.
* example.cron.daily -> /etc/cron.daily/emdr
