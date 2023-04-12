
# TODO list
- json schema to `possible_values.yaml` 
- handle errors:
    - when no values exist
    - when files does not exist
- fix remove approved file
- use golang rather than filesystem commands
- use config file. this allows testing in different languages
- [x] `id` should be reserved input key. Maybe change the structure of yaml template ?
- error when min > max
- store combos in file and load them so that the process may be stopped and continued
- parallel execution in sandboxed env until one of go routines increment mutation score:
    - if more than test increments mutation score. they should be re-executed in parallel like the traditional way
    - tests that does not increment mutation score should be deleted
    - then update test_cases.yaml file and then create go routines again
    - parallel execution should be optional. for some cases, parallel tests are not possible
- sandboxed env should be optional
    - `test_cases.yaml` can be passed as env variable
- use map to store combos. combos hash is the key
- get value args should return only the value. extracting value from output should be within commands

