# Receptionist

Receptionist is a very simple Golang application that creates a single page of links to the containers you have running 
on your server.

### Why?

I run a small server at home with a number of docker images running, its purely for prototyping and is no a production 
quality setup. I always forget what tool is running on what port of that server, so I wanted a simple landing page that 
would list the running containers and which port that are listening on.

### Usage

`$ docker run -v /var/run/docker.sock:/var/run/docker.sock:ro -p 8080:8080 cuotos/receptionist`

Receptionist will list any containers that contain the environment variable `RECEPTIONIST=<port>`

i.e.

`docker run --name webserver -e RECEPTIONIST=4567 -p 4567:80 nginx`

Will show a link to `http://localhost:4567` on the UI.

### Environment Variables

* `PORT` - The port for receptionist ui to list on (default: `8080`).
* `WATCHVAR` - The environment variable for Receptionist to look for on running containers.

### Volumes

* `/var/run/docker.sock:ro` - Receptionist needs to be able to see what containers are running on the Docker host.