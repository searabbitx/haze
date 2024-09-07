# Haze

## Run all tests

```bash
make test
```

## Format code

```bash
make format
```

## Todos
- [x] introduce configurable matchers (not only the 500s)
- [x] add filters (same as matchers)
- [x] match/filter response strings
- [x] add a probe mechanism
- [x] specify the output dir
- [x] make some sane logging and general output look
- [x] make a custom help message to group flags 
- [x] error handling for conection refused etc
- [x] parallel requests!
- [ ] parse HARs
- [ ] handle `multipart/form-data`
- [ ] add a parameter to array mutation ( foo=bar -> foo[]=bar )
- [ ] split mutation and mutable into multiple files
