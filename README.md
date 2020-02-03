# Simple Web Port Knocking Implementation

## What is the motivation of this repository?

Accepting TCP connections from the open internet is something I'm fairly paranoid about
once packets are parsed by anything above the kernel. Especially if those services aren't
web or SSH servers. Like a Minecraft server. I have no idea what the Minecraft wire protocol
looks like, so I don't know exactly what I'm getting into when I open a Minecraft port to the
open internet, but suffice to say I'd really rather not do that.

But I would like to let people selectively connect if they want. If I could ask my friends for
their IP addresses, I could add their addresses to my firewall config and I'd be fine. But that's
more work than I want to put in. Also many other people in the history of networked computing have
had this exact same desire.

Port knocking is one way people used to do this, but that requires an understanding of how network
protocols work and I'd rather not have to explain to each one of my friends what a TCP port in the
Discord before they can get on my server. Port knocking was basically an implementation of
authorization based on "something you know," and a URL with some entropy is a thing I can
give my friends that can serve the same purpose as a sequence of ports to connect to.

I didn't find anything else that really satisfied this use case...there's a node package out
there but it seemed like overkill and maybe only for web traffic? Also I want to keep the
amount of parsing that needs to be done as small as possible given that this service is going
to be completely open to the internet. So I made one myself.

## What this service does not do

Logging, periodic removal of whitelisted IP addresses, keep you from shooting yourself in the foot,
HTTPS, multiple authentication keys. It's barebones and works for my specific threat model. Be guaranteed
to work on anything other than a Linux distro similar to Ubuntu. IPv6 support.

The main [server.go](./server.go) basicaly listens on a port you define in an environment variable,
waits for a GET request to a route you define in an environment variable, and runs a shell command
hardcoded in the file. Also, in my estimation, you probably don't want any of that processing done
by the service listening on the open internet. You might want HTTP logging and I might add that at
some point, but you could also use tcpdump to log connections too and not lose out on any data.
And it probably makes sense to do IP address cleanup with a separate service running on a cron job
(I might add some logic for that too).

I have a simple systemd unit for running the service at startup, might add it at some point.

Port knocking is done in the clear, so not having TLS doesn't seem like a grave sin.

## Usage:

### Prerequsites:

You have ufw installed and enabled. If you want you can change the go file to use straight iptables but that's
up to you to implement if you want. Also the user the server executes as must be able to run ufw with passwordless
sudo. At least that's how I did it. I tried giving the binary CAP\_NET\_ADMIN but that doesn't let you mess with
ufw rules. Maybe it works with iptables if you do it that way. Might test that and update it later. If you think
it's a bad idea to give a binary listening on the open internet the ability to execute ufw with passwordless sudo
drop an issue.

Build the server with `go build server.go`. Execute the server, providing the following 4
environment variables to the runtime:

 * `LISTEN\_SPEC`: An IP address colon port string, e.g. "127.0.0.1:8080", "10.0.0.2:1111". This is the
   what the go HTTP server will bind to. The thing open to the internet. If you're behind a router you'll
   probably wanna forward a port in.
 * `HANDLER\_ROUTE`: A string containing at least 12 characters. The HTTP path that when browsed to will
   cause a command to be run to add a firewall rule. Probably should be hard to guess.
 * `LOCAL\_SERVICE\_ADDRESS:` The IP address the service you are protecting is running on. Specifically,
   this is what will be added to the "to" directive when a ufw rule is added.
 * `SERVICE\_PORT`: The port the service you are protecting is running on.


Then give your friends the HTTP url that corresponds to your WAN IP address, HTTP listening port, and route.

## License:

Dual licensed under the MIT license and the AntiLicense, see [LICENSE.md](./LICENSE.md).
