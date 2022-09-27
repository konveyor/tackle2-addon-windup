# tackle2-addon-windup

[![Windup Addon Repository on Quay](https://quay.io/repository/konveyor/tackle2-addon-windup/status "Windup Addon Repository on Quay")](https://quay.io/repository/konveyor/tackle2-addon-windup) [![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/konveyor/tackle2-addon-windup/pulls) [![Test Windup Addon](https://github.com/mundra-ankur/tackle2-addon-windup/actions/workflows/test-windup.yml/badge.svg?branch=main)](https://github.com/mundra-ankur/tackle2-addon-windup/actions/workflows/test-windup.yml)

Tackle (2nd generation) addon for Windup.


Task data.

*=_optional_

```
{
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
