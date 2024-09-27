## What is this?
This is the website made for the course of Web And Software Architecture to which i scored [30L/30](https://en.wikipedia.org/wiki/Academic_grading_in_Italy#University).
The website provides the following functions according to the [homework specifications](https://github.com/Alessio-Olivieri/WasaPhoto/blob/master/Homework%20specifications/Project%20specifications.pdf):
* Login
* Change username
* Upload/Delete photo
* Like/Unlike photo
* Comment/Uncomment photo
* Follow/Unfollow user
* Ban/Unban user
* User Profile


## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is an example of a daemon for checking the health of servers daemons; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the documentation (usually, for APIs, this means an OpenAPI file)
* `service/` has all packages for implementing project-specific functionalities
	* `service/api` contains an example of an API server
	* `service/globaltime` contains a wrapper package for `time.Time` (useful in unit testing)
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is an example of a web frontend in Vue.js; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

Other project files include:
* `open-npm.sh` starts a new (temporary) container using `node:lts` image for safe web frontend development (you don't want to use `npm` in your system, do you?)

## Go vendoring

This project uses [Go Vendoring](https://go.dev/ref/mod#vendoring). You must use `go mod vendor` after changing some dependency (`go get` or `go mod tidy`) and add all files under `vendor/` directory in your commit.

For more information about vendoring:

* https://go.dev/ref/mod#vendoring
* https://www.ardanlabs.com/blog/2020/04/modules-06-vendoring.html

## Node/NPM vendoring

This repository contains the `webui/node_modules` directory with all dependencies for Vue.JS. You should commit the content of that directory and both `package.json` and `package-lock.json`.

## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```shell
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-embed
exit
# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

You can launch the backend only using:

```shell
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run dev
```

## How to build (docker)
### Build images
```shell
docker build -t wasaphoto-frontend:latest -f Dockerfile.frontend .
docker build -t wasaphoto-backend:latest -f Dockerfile.backend .
```
### Run images
```shell
docker run -it --rm -p 3000:3000 wasaphoto-backend:latest
docker run -it --rm -p 8080:80 wasaphoto-frontend:latest
```

## Known issues

### Apple M1 / ARM: `failed to load config from`...

If you use Apple M1/M2 hardware, or other ARM CPUs, you may encounter an error message saying that `esbuild` (or some other tool) has been built for another platform.

If so, you can fix issuing these commands **only the first time**:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm install
exit
# Now you can continue as indicated in "How to build/run"
```


**Use these instructions only if you get an error. Do not use it if your build is OK**.

### My build works when I use `npm run dev`, however there is a Javascript crash in production/grading

Some errors in the code are somehow not shown in `vite` development mode. To preview the code that will be used in production/grading settings, use the following commands:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-prod
npm run preview
```

## License

See [LICENSE](LICENSE).
