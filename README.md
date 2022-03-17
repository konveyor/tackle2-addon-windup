# tackle2-addon-windup
Tackle (2nd generation) addon for Windup.


Task data.

*=_optional_

```
{
  application: int,
  output: string,
  mode: {
    binary: bool,
    withDeps: bool,
    artifact: string,
  },
  sources: [str,],
  targets: [str,],
  scope: {
    withKnown: bool,
    packages: {
      included: [str,],
      excluded: [str,]
    }
  },
  rules*: {
    path: str, 
    tags: {
      included: [str,],
      excluded: [str,]
    }
  }
}
```
