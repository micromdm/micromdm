When you're just starting to use micromdm on a server, you might run 
```
micromdm serve -server-url=https://mdm.acme.co
``` 
and follow the output directly. But as soon as you close your terminal/ssh session the server will stop running. 

A standard way to run services on linux is using [`systemd`](https://coreos.com/os/docs/latest/getting-started-with-systemd.html). Systemd has a number of benefits, but mainly:
- it will keep your process running.
- it will restart your process after a failure.
- it will restart your process after a server restart.
- it will pass the logs written to stout/stderr to `journalctl` or `syslog`. 

Getting started is easy. 
First, create a file like called `micromdm.service` on your linux host. 

```
[Unit]
Description=MicroMDM MDM Server
Documentation=https://github.com/micromdm/micromdm
After=network.target

[Service]
ExecStart=/usr/local/bin/micromdm serve \
    -server-url=https://mdm.acme.co \
    -api-key=secret \
    -filerepo /home/acme/mdmrepo
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Note that the `ExecStart` should have the `micromdm serve` command with the configuration appropriate for your server. Change/add/remove CLI flags as needed, but keep them with the `ExecStart=` line.

Once you created the file, you need to move it to `/etc/systemd/system/micromdm.service` and start the service.

```
sudo mv micromdm.service /etc/systemd/system/micromdm.service
sudo systemctl start micromdm.service
sudo systemctl status micromdm.service
```

If your `ExecStart` line is all correct you should see the service running. 

Use `sudo journalctl -u micromdm.service -f` to tail the server logs. 

# Making changes

Sometimes you'll need to update the systemd unit file defining the service. To do that, first open `/etc/systemd/system/micromdm.service` in a text editor, and apply your changes. 

Then, run 
```
sudo systemctl daemon-reload
sudo systemctl restart micromdm.service
```

# References

- https://coreos.com/os/docs/latest/getting-started-with-systemd.html
- https://www.digitalocean.com/community/tutorials/systemd-essentials-working-with-services-units-and-the-journal
- https://www.freedesktop.org/software/systemd/man/systemd.service.html

# Showing micromdm logs using journalctl
```
sudo journalctl -u micromdm.service -f
```