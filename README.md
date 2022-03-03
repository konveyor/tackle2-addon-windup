# tackle2-addon-windup
Tackle (2nd generation) addon for Windup.


Task data.

*=_optional_

```
{
  application: int,
  mode: {
    binary: bool,
    withDeps: bool,
    artifact* {
      bucket: int,
      path:   str
    }
  },
  targets: [str,],
  scope: {
    withKnown: bool,
    packages: {
      included: [str,],
      excluded: [str,]
    }
  },
  rules*: {
    directory: {
      bucket: int,
      path:   str
    },
    tags: {
      included: [str,],
      excluded: [str,]
    }
  }
}
```
