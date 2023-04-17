# Geny Test
The tool that magically generates test cases to ensure `1.0` [mutation score](https://www.guru99.com/mutation-testing.html#-mutation-score).

The tool can be configured to work with any test that have incremental score such as coverage tests


The tool is based on 2 technologies:
- [approval tests](https://approvaltests.com/)
- [mutation tests](https://en.wikipedia.org/wiki/Mutation_testing)


# Why ?
Before dealing with a complicated legacy code, we need to ensure that we have the required tests that prevents regressions.
Coverage tests are not enough. With `1.0` mutation score we can modify the code without fear of hidden regression of a code that we don't know.

A function with 4 parameters and 5 possible values each parameter will have `5 ** 4 = 225` testcases.
And the number of tests can reach thounds easilly if possible values are more than 5.

Guessing the test cases is hard if the function is complex. Given the example of [GildedRose-Refactoring-Kata](https://github.com/emilybache/GildedRose-Refactoring-Kata), I spent more then hour to find the suitable testcases and I failed.


In a perfect world mutation score is `1.0` :)


The challenge in this tool that the execution of mutation tests consume a lot CPU and memory and take some time (sometimes 60 seconds). And the more the testcases are, the more the execution consumes ressources and time. Plus it require some kind of sandboxed environment because the approval test file will change for each test.
So a good parallel strategy and a good algorithm should be used.



# Example
## GildedRose possible_values.yaml
Given the example of [GildedRose-Refactoring-Kata](https://github.com/emilybache/GildedRose-Refactoring-Kata)


Given the `possible_values.yaml` config file:
``` yaml
name:
  type: string
  values:
    static:
      - Aged Brie
      - Sulfuras, Hand of Ragnaros
      - Backstage passes to a TAFKAL80ETC concert
      - Some title
    ranges: []
sellIn:
  type: int
  values:
    static: []
    ranges:
      -
        min: -2
        max: 2
      -
        min: 9
        max: 11
      -
        min: 5
        max: 7

quality:
  type: int
  values:
    static: []
    ranges:
      -
        min: -2
        max: 2
      -
        min: 47
        max: 52

```

## Execution

That generates 484 test cases

<details>
  <summary>
  Phase 1: Eliminate useless scores with parallel execution of go routinse
  This process is executed in parallel.

  At the end the number of test cases are reduced from `484` to `23`
  With basic maths it eliminates `95.25%` of test cases

  </summary>

```bash
[test]$ go run main.go
---> Mutation
2023/04/17 19:46:10 scores[0] = 14.47
2023/04/17 19:46:35 scores[483] = 100.00
2023/04/17 19:46:58 scores[241] = 86.84
2023/04/17 19:47:24 scores[120] = 44.74
2023/04/17 19:47:26 scores[362] = 97.37
2023/04/17 19:47:57 scores[60] = 44.74
2023/04/17 19:47:57 ---------> marked usesless form 61 to 120
2023/04/17 19:47:57 --------------------------------------------> calculated 13.43%
2023/04/17 19:47:59 scores[180] = 84.21
2023/04/17 19:48:04 scores[301] = 93.42
2023/04/17 19:48:05 scores[422] = 98.68
2023/04/17 19:48:38 scores[30] = 44.74
2023/04/17 19:48:38 ---------> marked usesless form 31 to 60
2023/04/17 19:48:38 --------------------------------------------> calculated 20.25%
2023/04/17 19:48:41 scores[150] = 73.68
2023/04/17 19:48:42 scores[210] = 86.84
2023/04/17 19:48:48 scores[271] = 90.79
2023/04/17 19:48:48 ---------> marked usesless form 211 to 241
2023/04/17 19:48:48 --------------------------------------------> calculated 27.07%
2023/04/17 19:48:48 scores[331] = 97.37
2023/04/17 19:49:23 scores[392] = 97.37
2023/04/17 19:49:23 ---------> marked usesless form 363 to 392
2023/04/17 19:49:23 --------------------------------------------> calculated 33.47%
2023/04/17 19:49:23 scores[15] = 26.32

....


2023/04/17 20:00:26 --------------------------------------------> calculated 98.97%
2023/04/17 20:00:29 scores[313] = 94.74
2023/04/17 20:00:29 --------------------------------------------> calculated 99.17%
2023/04/17 20:00:29 ---------> marked usesless form 313 to 313
2023/04/17 20:00:29 --------------------------------------------> calculated 99.17%
2023/04/17 20:01:04 scores[398] = 97.37
2023/04/17 20:01:04 ---------> marked usesless form 398 to 398
2023/04/17 20:01:04 --------------------------------------------> calculated 99.38%
2023/04/17 20:01:05 scores[445] = 98.68
2023/04/17 20:01:05 ---------> marked usesless form 447 to 483
2023/04/17 20:01:05 --------------------------------------------> calculated 99.59%
2023/04/17 20:01:05 ---------> marked usesless form 445 to 445
2023/04/17 20:01:05 --------------------------------------------> calculated 99.59%
2023/04/17 20:01:05 scores[179] = 84.21
2023/04/17 20:01:05 ---------> marked usesless form 180 to 180
2023/04/17 20:01:05 --------------------------------------------> calculated 99.79%
2023/04/17 20:01:07 scores[270] = 90.79
2023/04/17 20:01:07 ---------> marked usesless form 271 to 271
2023/04/17 20:01:07 --------------------------------------------> calculated 100.00%
```
</details>


---

<details>
  <summary>
  Phase 2: Fine tune the end results and elimnates testcases that does not increment mutation score with an other order of testcases.

  This phase is optional

  At the end the number of test cases are reduced from `23` to `20`
  With basic maths it eliminates `13%` of test cases
  </summary>

``` logs
2023/04/17 20:01:07 ------------------------------> start sequentialFilter
--- individual socres
2023/04/17 20:02:18 scores[1] = 14.47
2023/04/17 20:02:18 scores[3] = 17.11
2023/04/17 20:02:18 scores[0] = 14.47
2023/04/17 20:02:19 scores[5] = 18.42
2023/04/17 20:02:19 scores[4] = 15.79
2023/04/17 20:02:19 scores[7] = 19.74
2023/04/17 20:02:19 scores[2] = 18.42
2023/04/17 20:02:20 scores[6] = 30.26
2023/04/17 20:03:28 scores[9] = 35.53
2023/04/17 20:03:29 scores[13] = 13.16
2023/04/17 20:03:29 scores[11] = 28.95
2023/04/17 20:03:29 scores[8] = 22.37
2023/04/17 20:03:29 scores[10] = 30.26
2023/04/17 20:03:29 scores[15] = 30.26
2023/04/17 20:03:30 scores[14] = 15.79
2023/04/17 20:03:30 scores[12] = 28.95
2023/04/17 20:04:29 scores[16] = 22.37
2023/04/17 20:04:30 scores[20] = 22.37
2023/04/17 20:04:30 scores[17] = 27.63
2023/04/17 20:04:30 scores[21] = 19.74
2023/04/17 20:04:30 scores[19] = 21.05
2023/04/17 20:04:30 scores[18] = 28.95
2023/04/17 20:04:30 scores[22] = 34.21
--- sort test cases and recalculate scores in parallel for each subcombo of testcases
2023/04/17 20:05:41 scores[4] = 57.89
2023/04/17 20:05:41 scores[5] = 63.16
2023/04/17 20:05:41 scores[0] = 28.95
2023/04/17 20:05:41 scores[6] = 65.79
2023/04/17 20:05:41 scores[1] = 36.84
2023/04/17 20:05:41 scores[7] = 73.68
2023/04/17 20:05:42 scores[3] = 52.63
2023/04/17 20:05:42 scores[2] = 47.37
2023/04/17 20:06:50 scores[9] = 75.00
2023/04/17 20:06:50 scores[10] = 77.63
2023/04/17 20:06:51 scores[11] = 85.53
2023/04/17 20:06:51 scores[8] = 75.00
2023/04/17 20:06:52 scores[12] = 88.16
2023/04/17 20:06:52 scores[14] = 94.74
2023/04/17 20:06:52 scores[13] = 92.11
2023/04/17 20:06:53 scores[15] = 96.05
2023/04/17 20:07:50 scores[18] = 97.37
2023/04/17 20:07:50 scores[16] = 96.05
2023/04/17 20:07:50 scores[17] = 96.05
2023/04/17 20:07:51 scores[19] = 98.68
2023/04/17 20:07:51 scores[20] = 100.00
2023/04/17 20:07:51 scores[21] = 100.00
2023/04/17 20:07:51 scores[22] = 100.00
2023/04/17 20:07:51 ------------------------------> end sequentialFilter

```

</details>


## Results

<details>
<summary>
Example of Results:
</summary>

``` yaml
-
  comboHash: "323b0843c996705319c9645513f9617483cac124"
  input:
    -
        name: Aged Brie
        sellIn: 0
        quality: -1

-
  comboHash: "18196499088da78902887f2b47ae2dcf25e464ec"
  input:
    -
        name: Some title
        sellIn: 11
        quality: 50

-
  comboHash: "6c7d85b57aeea6ba99aa3e561d23b849c4ef6f4e"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 10
        quality: -1

-
  comboHash: "9bc6b38ba6bfb2d9905428adff6ca05194096bb6"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 1
        quality: 50

-
  comboHash: "5a1d09011c5a7012f208dd420f89b6fa1aa90d3a"
  input:
    -
        name: Some title
        sellIn: 0
        quality: 0

-
  comboHash: "1a420cee1ce21fe5a12d402b47b11f78793c3073"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 0
        quality: 50

-
  comboHash: "6f3db4d50ef034fda030ac8b811549889bffdc30"
  input:
    -
        name: Aged Brie
        sellIn: 0
        quality: 50

-
  comboHash: "c8ce3f73efaa9cd1963a3ca52173ab414f03da08"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 5
        quality: -1

-
  comboHash: "88dc5dd974a57b84141dfa17123f0181d134b288"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 6
        quality: 48

-
  comboHash: "551d1a04e58ee41f969f120fbd41e2df34b8096f"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 6
        quality: -1

-
  comboHash: "f2fbb545ce825406d4bdc57e7150c5678eeb92cd"
  input:
    -
        name: Some title
        sellIn: 0
        quality: 50

-
  comboHash: "8a982d7e3253df2a307878374711fd69d56b8ea7"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 5
        quality: 48

-
  comboHash: "339e331552c9bd8bbfd0389ea2d1be49043d6dad"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 5
        quality: 49

-
  comboHash: "2732f9769bb026ef5a49d8cc6316b2d362c82592"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 11
        quality: -1

-
  comboHash: "f2dac14ea6977115803ccf2ec33626849162b8cc"
  input:
    -
        name: Aged Brie
        sellIn: 0
        quality: 48

-
  comboHash: "58068af86be0110470da14a6de97ee52d8ee684f"
  input:
    -
        name: Some title
        sellIn: 11
        quality: 1

-
  comboHash: "518c17d4042a49327768ef2d9423a5a6f2f53e0f"
  input:
    -
        name: Backstage passes to a TAFKAL80ETC concert
        sellIn: 5
        quality: 47

-
  comboHash: "5b3d16808ee26e70555edbe8bd189ff3fd6f4ec5"
  input:
    -
        name: Some title
        sellIn: 0
        quality: 2
```
</details>

# Todo
[todo.md](./todo.md)
