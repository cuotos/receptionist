# Receptionist

Receptionist is a very simple Golang application that creates a single page of links to the containers you have running 
on your server. It runs on port 8080 in the container, for you to expose anywhere you want.

It's not a proxy of any sort, it just lists the ports that you wanted listed by adding the label RECEPTIONIST to your 
containers.

### Why?

I run a small server at home with a number of docker images running, its purely for prototyping and is no a production 
quality setup. I always forget what tool is running on what port of that server, so I wanted a simple landing page that 
would list the running containers and which port that are listening on.

### Usage

#### Receptionist
`$ docker run -v /var/run/docker.sock:/var/run/docker.sock:ro -p 8080:8080 cuotos/receptionist`

Receptionist will list any containers that contain the label key `RECEPTIONIST`

#### Additional Containers
The `RECEPTIONIST` label can contain 1 or more comma seperated ports.

`docker run --name webserver -l RECEPTIONIST=4567,7654 -p 4567:80 container/image`

If your container exposes multiple ports, an optional *name* can be assigned to the ports by prefixing it with `<name>:`

`docker run --name webserver -l RECEPTIONIST=ui:4567,api:9999 container/image`

### Environment Variables (for Receptionist )

* `WATCHVAR` - The environment variable for Receptionist to look for on running containers (default `RECEPTIONIST`).

### Volumes

* `/var/run/docker.sock:ro` - Receptionist needs to be able to see what containers are running on the Docker host.

