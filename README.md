# GoAutoYT
<p align="center">GoAutoYT makes it easy for you to automatically download videos from as many YouTube channels as you'd like.</p>
<p align="center"><a href="https://travis-ci.com/xiovv/go-auto-yt"><img src="https://travis-ci.org/XiovV/go-auto-yt.svg?branch=master" alt="Build Status"></a><p>
<p align="center"><img src="https://raw.githubusercontent.com/XiovV/go-auto-yt/master/demo.png" width=700 alt="Screenshot of Example Documentation created with Slate"></p>

Features
------------
* **Clean, very simple design** - The dashboard only contains an input form where you can add a channel and configure checking intervals and what to download, and a little list of all your channels where you can delete them or tell the server to check for new uploads immediately.

* **Everything is on a single page** - You can view and control everything from just one page. 

* **Makes downloading videos/audio automatically very easy** - Just paste a link of a channel you want to download, set a checking interval and that's it, the server will keep checking for new uploads and download if necessary.

Getting Started (without Docker)
------------
### Prerequisites
* **Windows, Mac or Linux** - Only tested on Linux, but should work on Mac and Windows
* **Go, version 1.13.4 or newer**
* **youtube-dl**

### Setting Up (Tested on Linux, but should work on Mac. Windows - not sure)
```
git clone https://github.com/XiovV/go-auto-yt.git
cd go-auto-yt
go build
./go-auto-yt
```

You can now go to https://localhost:8080 and start using GoAutoYT.

Getting Started (with Docker)
------------
### Prerequisites
* **Docker** - Built on 19.03.3
* **Docker Compose** - Build on 1.24.0

### Configuring The Container
[Docker Hub Image](https://hub.docker.com/r/xiovv/go-auto-yt)

The `docker-compose.yml` file in the repository can be used as a guide for setting up your own containers. the only thing that needs to be checked is the `volumes` section. If you keep the default, you need to ensure that there are _downloads_ and _config_ folders in the same directory as the `docker-compose.yml` file. Otherwise, feel free to modify those mapping to your own local directories. [Docker Bind Mount Docs](https://docs.docker.com/storage/bind-mounts/)

```YAML
volumes:
    # Choose ONE of these
      - ./downloads:/app/downloads # local folder mapping
    #  - downloads:/app/downloads # docker volume mapping
    # And ONE of these
      - ./config:/app/config # local folder mapping
    #  - config:/app/config # docker volume mapping
```

If you wish to use Docker volume mappings, comment out the local folder lines and uncomment the docker volume lines. You must also uncomment the root `volumes` section at the bottom of the sample YAML file. This will instruct docker to create the volumes. 

```YAML
# uncomment this if using the docker volume mapping above
volumes:
  downloads:
  config:
```

You can view these volumes with the `docker volumes` command and view further information with `docker volume inspect VOLUME_NAME`. [Docker Volume Docs](https://docs.docker.com/storage/volumes/)

The port can also be changed by modifying the `ports:` item. By default, this is set to 8080 but if you wanted to change it to port 9000, for example, modify the docker-compose.yml to reflect the below. This follows the standard [docker container networking model](https://docs.docker.com/config/containers/container-networking/) of `host:container.`

```YAML
ports:
  - 9000:8080
```

The container runs with PID and GID of 1000 by default. This can be changed by passing environment variables to the container like so:

```YAML
environment:
  - PUID=1001
  - PGID=1001
```

### Running The Container
Once the configuration is complete, `docker-compose up -d` will pull and run the container for you in the background. The container will now be accessible from http://localhost:8080 (or whichever port you've modified it to) on that machine. Using `docker logs` will show the container logs for troubleshooting.

## Known Issues
* Checking interval is hardcoded to 5 hours, at the moment you cannot change it.

## Built With
* [Go](https://golang.org/) - Go Language
* [Gorilla Mux](https://github.com/gorilla/mux) - Go Multiplexer
* [Bootstrap](https://getbootstrap.com/) - CSS Framework
