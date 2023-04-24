import http from 'k6/http';
import { describe, expect } from 'k6/x/chai';

export let options = {
  thresholds: {
    checks: [{ threshold: 'rate == 1.00', abortOnFail: true }],
    http_req_failed: [{ threshold: 'rate == 0.00', abortOnFail: true }],
  },
};

export default function testSuite() {
  describe('Get the answer from httpbin', () => {
    const response = http.get('https://httpbin.test.k6.io/get?answer=42');

    expect(response.status, 'response status').to.equal(200);
    expect(response).to.have.validJsonBody();
    expect(response.json().args.answer, 'args.answer').to.equal('42');
  });
}
