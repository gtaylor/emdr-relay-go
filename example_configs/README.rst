Example Configs
===============

The files in this directory can be used as a starting point for supervising
and running emdr-relay-go. You'll want to change the paths to fit your
deployment, of course.

* supervisord-relay.conf will go in /etc/supervisor/conf.d/ on Ubuntu, or
  directly in /etc/supervisor/supervisor.conf on other distros.
* example.cron.daily -> /etc/cron.daily/emdr

If you want to limit one connection per IP to your relay, this is iptables rule:

* /sbin/iptables -I INPUT -p tcp --syn --dport 8050 -m connlimit --connlimit-above 1 -j DROP
