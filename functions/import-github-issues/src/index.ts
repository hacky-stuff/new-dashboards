import { CloudEvent } from 'cloudevents';
import { Context } from 'faas-js-runtime';
import { Octokit } from 'octokit';

/**
 * Your CloudEvents function, invoked with each request. This
 * is an example function which logs the incoming event and echoes
 * the received event data to the caller.
 *
 * It can be invoked with 'func invoke'.
 * It can be tested with 'npm test'.
 *
 * @param {Context} context a context object.
 * @param {object} context.body the request body if any
 * @param {object} context.query the query string deserialzed as an object, if any
 * @param {object} context.log logging object with methods for 'info', 'warn', 'error', etc.
 * @param {object} context.headers the HTTP request headers
 * @param {string} context.method the HTTP request method
 * @param {string} context.httpVersion the HTTP protocol version
 * See: https://github.com/knative/func/blob/main/docs/guides/nodejs.md#the-context-object
 * @param {CloudEvent} cloudevent the CloudEvent
 */
const handle = async (context: Context, cloudevent?: CloudEvent<Response>): Promise<CloudEvent<Response | string>> => {
  const meta = {
    source: 'function.eventViewer',
    type: 'echo'
  };

  // The incoming CloudEvent
  if (!cloudevent) {
    const response: CloudEvent<string> = new CloudEvent<string>({
      ...meta,
      ...{ type: 'error', data: 'No event received' }
    });
    context.log.info(response.toString());
    return response;
  }

  const octokit = new Octokit();
  const iterator = octokit.paginate.iterator(octokit.rest.issues.listForRepo, {
    owner: 'openshift',
    repo: 'console',
    per_page: 100,
    state: 'open'
  });

  // iterate through each response
  let issueCount = 0;
  for await (const { data: issues } of iterator) {
    for (const issue of issues) {
      console.log('Issue #%d: %s', issue.number, issue.title);
      issueCount++;
    }
  }

  console.log('issueCount', issueCount);

  console.log('environment variables', process.env);

  context.log.info(`
-----------------------------------------------------------
CloudEvent:
${cloudevent}

Data:
${JSON.stringify(cloudevent.data)}
-----------------------------------------------------------
`);
  // respond with a new CloudEvent
  return new CloudEvent<Response>({ ...meta, data: { issueCount: `123123${issueCount}` } });
};

export interface Response {
  issueCount: string;
}

export { handle };
