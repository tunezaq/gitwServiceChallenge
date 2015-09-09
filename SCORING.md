GITW Service Challenge Scoring
===

Scoring Criteria
--------------------
Each challenge is graded separately. You must satisfy the basic requirements in Challenge 1 in order to implement the subsequent challenges, but the subsequent challenges are independent of each other. Use your time wisely and fail where appropriate!

Challenge 1 Scoring
--------------------

| *Endpoint* | *Verb* | *Call Flow* | *Expectation* | *Point Value* |
| ------------- | ------------- | ------------- | ------------- | ------------- |
| /cache/ | POST | * Pass in a valid key / value pair  | * Status code 201 returned with a link to new cache item's key.  | 1 |
| /cache/ | GET | * Populate multiple cache items via POST calls. * Make GET call.  | * Status code 200 returned, with each cache item present in response. Should match the contract for returning multiple items.  | 1 |
| /cache/ | DELETE | Populate multiple cache items via POST calls. * Make DELETE call. * Make GET call.  | * Status code 204 in response to DELETE. * Status code 200 in response to GET, with an empty array returned under the 'cache' key.  | 1 |
| /cache/{key} | GET |   |   | 1 |
| /cache/{key} | PUT |   |   | 1 |
| /cache/{key} | DELETE |   |   | 1 |


Challenge 2 Scoring
--------------------

Challenge 3 Scoring
--------------------

Challenge 4 Scoring
--------------------
