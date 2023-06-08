"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.handle = void 0;
const cloudevents_1 = require("cloudevents");
const jira_client_1 = __importDefault(require("jira-client"));
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
    const jira = new jira_client_1.default({
        protocol: 'https',
        host: 'issues.redhat.com',
        bearer: cloudevent?.data?.['accessToken'],
        apiVersion: '2',
        strictSSL: true
    });
    const result = await jira.getAllBoards();
    // respond with a new CloudEvent
    return new cloudevents_1.CloudEvent({ ...meta, data: { result } });
};
exports.handle = handle;
