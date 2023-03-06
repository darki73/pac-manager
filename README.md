# Proxy Auto-Configuration Manager
This application allows you to automatically generate PAC files based on the stored rules in the SQLite database.  
It provides a simple CLI as well as a web interface to manage the rules.

**Table of Contents**
- [Installation](#installation)
- [Configuration](#configuration)
  - [Main Configuration](#main-configuration)
    - [Database Configuration](#database-configuration)
    - [PAC Configuration](#pac-configuration)
    - [Server Configuration](#server-configuration)
  - [Systemd Service](#systemd-service)
  - [Using a Reverse Proxy](#using-a-reverse-proxy)
- [Usage](#usage)
  - [CLI](#cli)

## Installation

### Building from Source

As this package is using SQLite as a database, [mattn/go-sqlite3](github.com/mattn/go-sqlite3) is used as a driver.  
The "problem" with this package, is that it requires CGO to be enabled.  
This means that you need to have the `gcc` installed on your system in order to be able to build the application from source.

The process of building the application from source is the same as for any other Go application.

## Configuration
If you would like application to automatically discover configuration, you need to:
1. Create a directory `/etc/pac.d` - this is where the application will look for the configuration file
2. Create a file `main.yaml` inside the `/etc/pac.d` directory - this is the name of the configuration file that the application will look for

However, you can also use the following application arguments:
1. `--config-path` - path to where configuration file is located
2. `--config` - name of the configuration file

### Main Configuration
Basic example configuration required to run application can be found in the [examples/main.yaml](https://github.com/darki73/pac-manager/blob/main/examples/main.yaml).

There are 3 sections available for you to configure:
1. `database` - configuration for the SQLite database
2. `server` - configuration for the web server
3. `pac` - configuration for the PAC file

#### Database Configuration
- `name` - name of the database file
- `path` - path where the database file will be stored

#### PAC Configuration
- `name` - name of the generated PAC file
- `path` - path where the generated PAC file will be stored

#### Server Configuration
- `host` - host address on which the web server will be listening
- `port` - port on which the web server will be listening

### Systemd Service
In order to be able to access web interface at all times without explicitly running `pacm run`, you can leverage the systemd service.

Provided example in [examples/pacm.service](https://github.com/darki73/pac-manager/blob/main/examples/pacm.service) will provide you with working service out of the box.

To install and enable service you need to do the following:
1. Copy the contents of the `pacm.service` to `/lib/systemd/system/pacm.service`
2. Run `systemctl daemon-reload`
3. Run `systemctl enable pacm.service`
4. Run `systemctl start pacm.service`

Assuming you have configured the service correctly, you should be able to access the web interface at `http://0.0.0.0:8080`.

### Using a Reverse Proxy
You can find the example configuration for the Traefk reverse proxy in the [examples/traefik-service.yaml](https://github.com/darki73/pac-manager/blob/main/examples/traefik-service.yaml)  

Assuming you've just installed the NGINX on the host machine with the default rules, only thing you need to do is to tell Traefik that static files will be served from the `/files` path.  

Configuration provided in the example will allow you to access the web interface on the `/` path, and the generated PAC file on the `/files/configuration.pac` path.

## Usage

**Things worth noting:**
1. Whenever you create or remove a proxy, PAC file will be regenerated
2. Whenever you create/update/remove a domain, PAC file will be regenerated
3. This logic applies to both CLI and Web interface.

### CLI
The application provides a simple CLI that allows you to manage the rules and list of proxies.

To see the list of available commands, run `pacm --help`.

To see the list of available options for a specific command, run `pacm <command> --help`.

Most of the following commands provide you with interactive way of working with them.   
Same functionality is present through the web interface.

- `pacm run` - runs the web server
- `pacm version` - prints the version of the application
- `pacm proxy` - lists commands for managing proxies
  - `pacm proxy add` - adds a new proxy
  - `pacm proxy delete` - removes a proxy
- `pacm domain` - lists commands for managing domains
  - `pacm domain add` - adds a new domain
  - `pacm domain delete` - removes a domain
  - `pacm domain update` - updates a domain