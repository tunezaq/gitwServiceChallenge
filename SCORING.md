GITW Service Challenge Scoring
===

Scoring Criteria
--------------------
Each challenge is graded separately. You must satisfy the basic requirements in Challenge 1 in order to implement the subsequent challenges, but the subsequent challenges are independent of each other. Use your time wisely and fail where appropriate!

Challenge 1 Scoring
--------------------

The first portion of scoring challenge 1 is verifying that the happy path for section 1 requirements are fulfilled. This portion will be scored according to the table below.

| *Endpoint* | *Verb* | *Call Flow* | *Expectation* | *Point Value* |
| ------------- | ------------- | ------------- | ------------- | ------------- |
| /cache/ | POST | Populate a single key / value pair via a POST call.  | Status code 201 returned with a link to new cache item's key.  | 1 |
| /cache/ | GET | Populate multiple cache items via POST calls.</br><br/> Make GET call.  | Status code 201 in response to POST.<br/><br/>Status code 200 returned in response to GET, with each cache item present in array under 'cache' key.  | 1 |
| /cache/ | DELETE | Populate multiple cache items via POST calls.<br/><br/> Make DELETE call.<br/><br/> Make GET call.<br/><br/>  | Status code 201 in response to POST.<br/><br/>Status code 204 in response to DELETE.<br/><br/> Status code 200 in response to GET, with an empty array returned under the 'cache' key.  | 1 |
| /cache/{key} | GET | Populate a single cache item via a POST call.<br/><br/>Make GET call using key from POST.<br/><br/>  | Status code 201 in response to POST.<br/><br/> Status code 200 returned in response to GET, with key and value matching incoming parameters.  | 1 |
| /cache/{key} | PUT | Populate a single cache item via a POST call.<br/><br/>Make PUT call using key from POST, altering key and value.<br/><br/>Make GET call using new key.<br/><br/>  | Status code 201 in response to POST.<br/><br/>Status code 204 in response to PUT.</br><br/>Status code 200 in response to GET, with updated key and value shown.   | 1 |
| /cache/{key} | DELETE | Populate a single cache item via a POST call.<br/><br/>Make DELETE call using key from POST.<br/><br/>Make GET call using key from POST.<br/><br/>  | Status code 201 in response to POST.<br/><br/>Status code 204 in response to DELETE.<br/><br/>Status code of 404 in response to GET.  | 1 |

The second portion of scoring challenge 1 is validating the explicit unhappy paths. This portion will be scored according to the table below.

| *Endpoint* | *Verb* | *Call Flow* | *Expectation* | *Point Value* |
| ------------- | ------------- | ------------- | ------------- | ------------- |
| /cache/ | POST | Populate a single cache item via a POST call.<br/><br/>Duplicate the POST call.  | Status code 201 returned for the first POST.<br/><br/>Status code 409 returned for the second call.<br/><br/>  | 1 |
| /cache/{key} | GET | Make GET call using non-existent key  | Status code 404 in response to GET.  | 1 |
| /cache/{key} | PUT | Make PUT call using non-existent key  | Status code 404 in response to PUT.  | 1 |
| /cache/{key} | DELETE | Make DELETE call using non-existent key  | Status code 404 in response to DELETE.  | 1 |

The final portion of scoring challenge 3 is testing around input length, not-quite-valid JSON, and character sets. We will not ruin the surprise by sharing these test cases, but this portion is worth 5 points.

Challenge 2 Scoring
--------------------

Challenge 3 Scoring
--------------------

Challenge 4 Scoring
--------------------
