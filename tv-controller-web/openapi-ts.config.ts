import {defineConfig} from '@hey-api/openapi-ts';

export default defineConfig({
    client: '@hey-api/client-fetch',
    input: '../oapi/api.yaml',
    output: 'app/vlccontrol',
});
