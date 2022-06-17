import * as http from 'http';
import net from 'net';
import {createHttpTerminator} from 'http-terminator';

export declare namespace MockWebServer {
    type RequestMatcher = {
        method: 'GET' | 'POST' | 'PUT' | 'DELETE'
        path: string
    }

    type ServerResponse = {
        statusCode: number
        body: Record<string, unknown>
    }

    type RecordedRequest = {
        method: string | undefined
        path: string | undefined
        body: string
    }
}

export type MockWebServer = {
    stop: () => Promise<void>
    url: (path?: string) => string
    register: (matcher: MockWebServer.RequestMatcher, response: MockWebServer.ServerResponse) => MockWebServer
    lastRequest: () => MockWebServer.RecordedRequest | undefined
}

const create = (): MockWebServer => {

    const registeredRequests: { matcher: MockWebServer.RequestMatcher, response: MockWebServer.ServerResponse }[] = [];
    let lastRequest: MockWebServer.RecordedRequest | undefined;

    const setCorsHeaders = (res: http.ServerResponse) => {
        res.setHeader('Access-Control-Allow-Credentials', 'true');
        res.setHeader('Access-Control-Allow-Methods', '*');
        res.setHeader('Access-Control-Allow-Origin', '*');
        res.setHeader('Access-Control-Allow-Headers', '*');
    };

    const noMatchFoundResponse: MockWebServer.ServerResponse =
        {statusCode: 503, body: {message: 'No response could be found for the request'}};

    const findMatchingResponse = (req: http.IncomingMessage): MockWebServer.ServerResponse => {
        const matchingRequest = registeredRequests.find(({matcher}) => {
            return matcher.method === req.method && matcher.path === req.url;
        });
        return matchingRequest?.response ?? noMatchFoundResponse;
    };

    const requestListener = (req: http.IncomingMessage, res: http.ServerResponse) => {
        const recordedRequest: MockWebServer.RecordedRequest = {
            method: req.method,
            path: req.url,
            body: '',
        };

        setCorsHeaders(res);

        if (req.method === 'OPTIONS') {
            res.writeHead(200);
            res.end();
            return;
        }

        const response = findMatchingResponse(req);

        res.setHeader('Content-Type', 'application/json');
        res.writeHead(response.statusCode);

        req.on('data', chunk => {
            recordedRequest.body += chunk;
        });

        req.on('end', () => {
            lastRequest = recordedRequest;
            res.end(JSON.stringify(response.body));
        });
    };

    const server = http.createServer(requestListener);
    const terminator = createHttpTerminator({server});
    server.listen(0);

    const url = (path?: string) => {
        const address = server.address() as net.AddressInfo;
        return `http://localhost:${address.port}${path ?? ''}`;
    };

    const mockWebServer = {
        stop: () => terminator.terminate(),
        url,
        register: (matcher: MockWebServer.RequestMatcher, response: MockWebServer.ServerResponse) => {
            registeredRequests.push({matcher, response});
            return mockWebServer;
        },
        lastRequest: () => lastRequest,
    };
    return mockWebServer;
};

export const mockWebServer = {
    create,
};