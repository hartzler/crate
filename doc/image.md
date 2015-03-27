## create file

A tar file containing:
/manifest
/cargo/
/cargo/<path>/123abc.tar

A yaml file that describes how a component is packed into a crate.

### example: crates.armada.io/armada/haproxy-1.0.0.crate

    name: haproxy

    connectors:
      in:
        - name: http
          proto: tcp/http
      out:
        - name: pool
          proto: tcp/http
          cardinality: many

    hooks:
      install: /app/install.sh
      pre-start: /app/pre-start.sh
      start: /app/start.sh
      stop: /app/stop.sh
      cleanup: /app/cleanup.sh
      connection-http-up:
      connection-http-down:

    cargo:
    - http://cargo.armada.io/5581dc785e49e0898aa2e38f947144c956604a51efeca2297c414c05349de1a5.tar


### example: crates.armada.io/mysql.crate

    name: mysql

    connectors:
      in:
        - name: db
          proto: tcp/sql

    volumes:
      - name: data
        mount: /var/lib/mysql/data

   cargo:
      - 123abc.tar


### example: crates.armada.io/rails.crate

    name: rails

    connectors:
      in:
        - name: http
          proto: tcp/http
      out:
        - name: db
          type: tcp/sql
          cardinality: 1

    hooks:
      start: /app/start.sh
      stop: /app/stop.sh

    payload:
      app: /app

    cargo:
      - 123abc.tar


### example: example.com/my-rails-app.crate

    name: my-rails-app

    base: crates.armada.io/armada/rails.crate
