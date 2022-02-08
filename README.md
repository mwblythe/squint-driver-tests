# Squint SQL Driver Tests

This is a set of compatibility tests for the [squint/driver](https://github.com/mwblythe/squint/driver) package. The tests are in a separate repository so that the database driver dependencies are not part of the main [squint](https://github.com/mwblythe/squint) package.

## Running the Tests

First, clone this repository:

```shell
git clone https://github.com/mwblythe/squint-driver-tests.git
```

Next, make sure you have [Docker](https://www.docker.com/products/docker-desktop) and `docker-compose`. Then, simply run:

```shell
cd squint-driver-tests
docker-compose up
go test -v
```

There are currently tests for sqlite and mysql.