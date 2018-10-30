# glutton

Glutton is a small HTTP server that can be called with *ANY* data and the data is stored.

In this release the only destination for the data is the file system. 

## How to run glutton

* from source
   * run `make run` - this will spin up glutton on local port 4354
* docker
   * run `docker --rm -it -p 4354:4354 -v glutton:glutton defectus/glutton` - this will spin up glutton on local port 4354

## Configuration

First, command line arguments. At the moment two:

* -f *file.yaml* : use yaml file to configure application
* -d : enable debug messages

Configuration lives either in OS environment, or is provided as a yaml file. Yaml offers greater variablity and more importantly allows you to define more than one route.

Basic structure of the yaml file:
```yaml
debug: true
port: 8080
host: 0.0.0.0
settings:
  - name: default glutton route
    redirect: some_url
    uri: save
    parser: SimpleParser # choice of `SimpleParser`
    notifier: NilNotifier # choice of `NilNotifier`, `SMTPNotifier`
    saver: SimpleFileSystemSaver # choice of `SimpleFileSystemSaver`
    # SimpleFileSystemSaver settings
    output_folder: glutton # lcoation to which request are saved
    base_name: glutton_%d # name of request files (supports single numeric counter variable)
    # SMTPNotifier settings
    smtp_server: smtp.gmail.com
    smtp_port: 25
    smtp_use_tls: true
    smtp_from: your@email.address
    smtp_to: target@email.address
```

As you can see, the settings is fairly straight forward. When using the environment keys are:
* `DEBUG`
* `HOST`
* `PORT`
* `NAME`
* `URI`
* `REDIRECT`
* `PARSER`
* `NOTIFIER`
* `SAVER`

SimpleFileSystemSaver settings

* `OUTPUT_FOLDER`
* `BASE_NAME`

SMTPNotifier settings

* `SMTP_SERVER`
* `SMTP_PORT`
* `SMTP_USE_TLS`
* `SMTP_FROM`
* `SMTP_TO`

Please note that only one route can be defined with environment variables.

## Endpoint

A sample request can be found in the http/save-basic.http file. Effectively you have to do HTTP `POST` on `/v1/glutton/save`. As the payload is in no paricular format any payload will do.

## Output

Requests are stored on a path defined by the `OUTPUT_FOLDER` variable. If ommited it defaults to `glutton`.

## Future

In the future releases you hopefully find the following features

❌ saving to database

✔️ redirect on save 

❌ auth tokens (allow saving with a valid token only)


above all, keep this project low profile, I'm not building an application server here. glutton must be simple, stupid.
