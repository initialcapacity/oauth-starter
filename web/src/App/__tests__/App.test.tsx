import * as Redux from 'redux';
import {render, screen, waitFor} from '@testing-library/react';
import {App} from '../App';
import {TestAppContext} from '../../TestSupport/TestAppContext';
import {mockWebServer, MockWebServer} from '../../TestSupport/MockWebServer';
import {AppState, appStateStore} from '../AppStateStore';

describe('App', () => {

    let oauthServer: MockWebServer;
    let stateStore: Redux.Store<AppState>;

    beforeEach(() => {
        oauthServer = mockWebServer.create();
        stateStore = appStateStore.create();
    });

    afterEach(async () => {
        await oauthServer.stop();
    });

    test('rendering the App', () => {
        render(
            <TestAppContext>
                <App/>
            </TestAppContext>
        );
        expect(screen.queryByText('OAuth SPA starter')).not.toBeNull();
    });

    test('handling the oauth callback', async () => {
        localStorage.setItem('verifier', 'aVerifier');

        oauthServer
            .register(
                {method: 'POST', path: '/token'},
                {statusCode: 201, body: {access_token: 'anAccessToken'}}
            );

        render(
            <TestAppContext stateStore={stateStore} oauthServer={oauthServer} currentPath={'/callback'}>
                <App/>
            </TestAppContext>
        );

        await waitFor(() => {
            const expectedSignedInState = {signedIn: true, accessToken: 'anAccessToken'};
            expect(stateStore.getState().authStatus).toEqual(expectedSignedInState);
        });

        screen.getByText('count is: 0');
    });
});