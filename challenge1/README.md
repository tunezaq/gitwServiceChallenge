Requirements Dump 1
===

A successful POST should return status code 201 with a link to the new cache item's ID. It should return 409 if a resource with that key already exists.

GET, PUT, and DELETE should return 404 if passed in a cache item key that does not exist in the cache.

If invalid JSON is passed in to a PUT or POST, return a status code 406.


