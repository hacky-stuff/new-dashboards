import { CloudEvent } from 'cloudevents';
import { Context, LogLevel } from 'faas-js-runtime';

const handle = async (context: Context, cloudevent?: CloudEvent<never>): Promise<CloudEvent<never>> => {
  if (cloudevent) {
    context.log.info(`CloudEvent: ${cloudevent}, data: ${JSON.stringify(cloudevent.data, null, 2)}`);
    context.log.info(cloudevent);
  } else {
    context.log.info('No CloudEvent received');
  }
  // eslint-disable-next-line no-console
  console.log('Another test', cloudevent, cloudevent?.data);

  const response: CloudEvent<never> = new CloudEvent<never>({});
  context.log.info(response.toString());
  return response;
};

const logLevel = 'warn' as never as LogLevel; // TypeError when using LogLevel.warn directly m(

export { handle, logLevel };
