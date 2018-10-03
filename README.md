# glutton

Glutton is a small HTTP server that can be called with *ANY* data and the data is stored.

In this release the only destination for the data is the file system. 

## How to run glutton

* from source
   * run `make run` - this will spin up glutton on local port 4354
* docker
   * run `docker --rm -it -p 4354:4354 -v glutton:glutton defectus/glutton` - this will spin up glutton on local port 4354

## Endpoint

A sample request can be found in the http/save-basic.http file. Effectively you have to do HTTP `POST` on `/v1/glutton/save`. As the payload is in no paricular format any payload will do.

## Output

Requests are stored on a path defined by the `OUTPUT_FOLDER` variable. If ommited it defaults to `glutton`.

## Future

In the future releases you hopefully find the following feature

* saving to database
* redirect on save
* auth tokens (allow saving with a valid token only)
* notifications (e.g. email)

above all, keep this project low profile, I'm not building an application server here. glutton must be simple, stupid.
