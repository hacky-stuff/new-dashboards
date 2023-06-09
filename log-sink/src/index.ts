import { CloudEvent } from 'cloudevents';
import { Context } from 'faas-js-runtime';

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
const handle = async (context: Context, cloudevent?: CloudEvent<never>): Promise<void> => {
  // The incoming CloudEvent
  if (cloudevent) {
    context.log.info(`CloudEvent: ${cloudevent}, data: ${JSON.stringify(cloudevent.data, null, 2)}`);
  } else {
    context.log.info('No CloudEvent received');
  }

  // eslint-disable-next-line no-console
  console.log('Another test', cloudevent, cloudevent?.data);
};

export { handle };