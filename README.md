[![Go Report Card](https://goreportcard.com/badge/github.com/szkiba/xk6-chai)](https://goreportcard.com/report/github.com/szkiba/xk6-chai)
[![GitHub Actions](https://github.com/szkiba/xk6-chai/workflows/Test/badge.svg)](https://github.com/szkiba/xk6-chai/actions?query=workflow%3ATest+branch%3Amaster)
[![codecov](https://codecov.io/gh/szkiba/xk6-chai/branch/master/graph/badge.svg?token=D43BZDXASS)](https://codecov.io/gh/szkiba/xk6-chai)


# xk6-chai

A [k6](https://go.k6.io/k6) extension that embeds [Chai.js](https://www.chaijs.com/) (actually [k6chaijs](https://k6.io/docs/javascript-api/jslib/k6chaijs/) 4.3.4.3) into the k6 binary.

To use the extension, simply change the k6chaijs import path to `k6/x/chai`. This way, the test will not have runtime external JavaScript dependencies.

```js
import http from 'k6/http';
import { describe, expect } from 'k6/x/chai';

export let options = {
  thresholds: {
    checks: [{ threshold: 'rate == 1.00', abortOnFail: true }],
    http_req_failed: [{ threshold: 'rate == 0.00', abortOnFail: true }],
  },
};

export default function () {
  describe('Get the answer from httpbin', () => {
    const response = http.get('https://httpbin.test.k6.io/get?answer=42');

    expect(response.status, 'response status').to.equal(200);
    expect(response).to.have.validJsonBody();
    expect(response.json().args.answer, 'args.answer').to.equal('42');
  });
}
```

> **Note**
> This is a simple wrapper plugin. Making Chai.js k6 compatible is the credit of the k6 team.

## Download

You can download pre-built k6 binaries from [Releases](https://github.com/szkiba/xk6-chai/releases/) page. Check [Packages](https://github.com/szkiba/xk6-chai/pkgs/container/xk6-chai) page for pre-built k6 Docker images.

## Build

You can build the k6 binary on various platforms, each with its requirements. The following shows how to build k6 binary with this extension on GNU/Linux distributions.

### Prerequisites

You must have the latest Go version installed to build the k6 binary. The latest version should match [k6](https://github.com/grafana/k6#build-from-source) and [xk6](https://github.com/grafana/xk6#requirements).

- [Git](https://git-scm.com/) for cloning the project
- [xk6](https://github.com/grafana/xk6) for building k6 binary with extensions

### Install and build the latest tagged version

1. Install `xk6`:

   ```shell
   go install go.k6.io/xk6/cmd/xk6@latest
   ```

2. Build the binary:

   ```shell
   xk6 build --with github.com/szkiba/xk6-chai@latest
   ```

### Build for development

If you want to add a feature or make a fix, clone the project and build it using the following commands. The xk6 will force the build to use the local clone instead of fetching the latest version from the repository. This process enables you to update the code and test it locally.

```bash
git clone git@github.com:szkiba/xk6-chai.git && cd xk6-chai
xk6 build --with github.com/szkiba/xk6-chai@latest=.
```

## Docker

You can also use pre-built k6 image within a Docker container. In order to do that, you will need to execute something like the following:

**Linux**

```plain
docker run -v $(pwd):/scripts -it --rm ghcr.io/szkiba/xk6-chai:latest run /scripts/script.js
```

**Windows**

```plain
docker run -v %cd%:/scripts -it --rm ghcr.io/szkiba/xk6-chai:latest run /scripts/script.js
```

## Example scripts

There are many examples in the [scripts](https://github.com/szkiba/xk6-chai/tree/master/scripts) directory that show how to use various features of the extension.
