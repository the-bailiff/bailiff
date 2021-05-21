# bailiff 

Sidecar for distributed session for microservices

## Problem

When heaving deal with storing user's data between requests there are two basic variants how to do that:

- passing data in every request (`JWT`, raw values in `Cookie` header)
- storing data in some storage (files, memory, external db, etc) by some ID of _session_, with passing __only__ this ID in request data

There are tons of articles on the internet with pros and cons of both variants, so we will not stop on this.
The one thing that need to be said is that implementing session pattern in microservices environment could be quite tricky.
So here comes __bailiff__.

## Solution

__Bailiff__ is a sidecar for your backend apps which takes the routine of saving and restoring session data.
On every income request it checks if there is a session.
If so bailiff enriches request by passing new headers with session data (see example below).
Also, bailiff checks every response if there are any data to save or update in session.

All bailiff sidecars are connected with single storage, so all session data is shared between all microservices.
So if microservice _Foo_ saved `userID` in session, microservice _Bar_ will get it in next request.

Please see example below:

<img src="https://raw.githubusercontent.com/the-bailiff/bailiff/master/docs/concept.png" width="100%" alt="concept" />

## Usage

As it was described before it should be run as sidecar for existing backends. 

### Docker image

Basic example is:

```sh
docker run \
    -e BAILIFF_STORE_REDIS_ADDR=redis:6379 \
    -e BAILIFF_PROXY=http://app \
    -e BAILIFF_COOKIE_MAXAGE=3600 \
    bailiff/bailiff:latest
```

You can find more complex example with multiple backends in _examples_ folder.

### Storages

Bailiff is built with support of different stores in mind but only __redis__ storage is supported so far.

### Configuration

TO DO...
