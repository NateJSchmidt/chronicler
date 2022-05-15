# Chronicler
Chronicler is a simple, fast, and powerful eventlog that is capable of supporting event sourced systems.  Chronicler provides an intuitive, and easy to use, API for storing and retrieving events while also providing a way for services to subscribe to events via Kafka topics.  It can be horizontally scaled up or down depending on the amount of traffic that it needs to support.

## Documentation
For additional details, please see the [Chronicler Documentation](docs/README.md).

# Development instructions
This section contains instructions for how to go about development of the chronicler project.

## Running the project
Chronicler is a golang program that can be run from the command line.  Chronicler has two dependencies that are expected to be running prior to the chronicler executable:  postgres and kafka.  The easiest way to run these dependencies is via [docker-compose](https://docs.docker.com/compose/).  [deployments/docker-compose-dev.yaml](deployments/docker-compose-dev.yaml) will launch both a fresh postgresql instance and a kafka instance.  This can be done via the command line:

`docker-compose -f deployments/docker-compose-dev.yaml up -d`

Once everything is running, database migrations can be run via the [run-migrations.sh](/deployments/run-migrations.sh) script.  In order for migrations to work correctly, you'll need to set some environment variables.  For details, see the script.  This script can be run like:

`./run-migrations.sh`

Note that migrations will likely not work without the latest PR from sledger:  https://github.com/decadentsoup/sledger/pull/2.  This PR enables variable substitution in the [sledger.yaml](/migrations/sledger.yaml).

After migrations have been run, the easiest way to run chronicler is to run the following command:

`go run cmd/main.go`

## Contributing to the project
PRs are welcome.
