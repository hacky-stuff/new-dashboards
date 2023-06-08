"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.handle = void 0;
const cloudevents_1 = require("cloudevents");
const octokit_1 = require("octokit");
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
const handle = async (context, cloudevent) => {
    const meta = {
        source: 'function.eventViewer',
        type: 'echo'
    };
    // The incoming CloudEvent
    if (!cloudevent) {
        const response = new cloudevents_1.CloudEvent({
            ...meta,
            ...{ type: 'error', data: 'No event received' }
        });
        context.log.info(response.toString());
        return response;
    }
    const octokit = new octokit_1.Octokit();
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
    context.log.info(`
-----------------------------------------------------------
CloudEvent:
${cloudevent}

Data:
${JSON.stringify(cloudevent.data)}
-----------------------------------------------------------
`);
    // respond with a new CloudEvent
    return new cloudevents_1.CloudEvent({ ...meta, data: { issueCount } });
};
exports.handle = handle;
