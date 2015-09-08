GITW Service Challenge
===

Introduction
--------------------
Your company has gone SOA, all the way. It started off with so much promise: microservices as far as the eye could see, fully decoupled teams, and reusability out the wazoo.

Slowly, chaos begin to creep in. These services were not quite similar enough to make sense. Service A was written in Go and UTF-8 friendly, while Service B was written in Node and ASCII only. You had n services and 2n + 1 databases. It's time to add some sanity!

Your job is to bring up a simple caching service. Your service will exist solely to get, put, and delete data into and out of a cache. All these methods will support a clean and transparent contract to make cache interactions as simple as possible.

Requirements will trickle your way throughout the physical challenge, but here are a few hints to get started. Your clients will push you different types of data (all valid JSON), they'll do it concurrently, and they expect a response within 100 milliseconds.

Endpoints
--------------------
Your caching service will operate at http://localhost:8088/. Please provide the following endpoints which handle the following HTTP verbs:
* /cache/
 * POST - creates a new item in the cache. The body of the POST should match the contract specified in the Contract section. 
 * GET - gets the entire cache.
 * DELETE - deletes the cache.

* /cache/{key}
 * GET - gets the cache item with matching key.
 * PUT - updates the cache item with matching key. The body of the PUT should match the contract specified in the Contract section.
 * DELETE - deletes the cache item with matching key.

Contract
--------------------
Your cache service will expect input in valid JSON and emit output in valid JSON. A single cache item will always be described thusly:
```{
    "key": "{key}",
    "value": "{value}"
}```

An example payload:
```{
    "key": "problem_free_philosophy",
    "value": "Hakuna Matata"
}```

Keys and values may be any string, boolean, integer, or decimal value.

Any service call that returns more than one cache item will return a "cache" key with an array of the JSON object above:
```{
    "cache": [
        {
            "key": "foo",
            "value": 3.9999
        },
        {
            "key": "bar",
            "value": true
        }
    ]
}```

