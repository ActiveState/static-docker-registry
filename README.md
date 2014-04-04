static-docker-registry
======================

This is a demonstration of making 'docker pull' fetch layers (layers in
particular, but also the metadata) directly from a CDN store instead of having
to go throughly the comparably unreliable registry middleman.

We achieve this using:

1. [static registry support](https://github.com/dotcloud/docker/pull/4607) PR
   in docker to make it ignore the special X-* http headers if not present.
   though this is not strictly required if we add the headers to the endpoint app.

2. [patching
   docker-registry](https://github.com/ActiveState/docker-registry/commit/002bf256e7)
   to store 'tags' and 'images' files under each repositority. this allows the
   docker client to directly request /tags and /images from the CDN. the /tags
   file in particular is required as it is not possible to enumerate the /tag_*
   files from the CDN.

3. the [endpoint web
   app](https://github.com/ActiveState/static-docker-registry/blob/master/endpoint)
   that merely responds to registry API request with a redirect to the underlying
   CDN.

4. Use a storage backend like Swift (or s3) in docker-registry (to which you
   normally 'docker push'), and enable CDN syncing for use with the endpoint app.

What happens when you run 'docker pull' is:

* /v1/_ping is requested, and endpoint returns a 200 response

* docker proceeds to fetch metadata and layers

* for every request therein, the endpoint will respond with a http redirect to
  the corresponding CDN file.

* docker respects the redirect response, and fetches from the CDN.

As a result, 'docker pull' is made more reliable as it doesn't rely on a
registry middleman at all. The endpoint app is still useful as a federation
endpoint so that "docker pull <myendpoint.com>/foo/myimg" knows where to go to
get the images from.


