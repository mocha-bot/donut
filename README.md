# Donut Call

## Description

A way to call a function with a donut.

## Example Result

```md
  Start...
  Bob paired with Grace
  Frank paired with Alice
  Eve paired with David
  Charlie paired with Harry

  Do calls...
  Alice called Bob
  Charlie called David
  Eve called Frank

  Add person...
  Added Ivan
  Added Goldi
  Added Samde

  RePair...
  Harry paired with Goldi
  Ivan paired with Grace
  3-way call Samde

  Remove person...
  Removed Ivan

  RePair...
  Goldi paired with Samde
  Grace paired with Harry

  Print...
  People: Alice, Bob, Charlie, David, Eve, Frank, Grace, Harry, Goldi, Samde
  PeopleMap: map[Alice:true Bob:true Charlie:true David:true Eve:true Frank:true Goldi:false Grace:false Harry:false Samde:false]
  MatchMap: map[Alice:Frank Bob:Grace Charlie:Harry David:Eve Eve:David Frank:Alice Goldi:Samde Grace:Harry Harry:Grace Ivan:Grace Samde:Goldi]
  Completed: Alice, Bob, Charlie, David, Eve, Frank
  Remaining: Grace, Harry, Goldi, Samde
```
