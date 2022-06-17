import * as React from 'react';
import {Provider} from 'react-redux';
import {AppState, appStateStore} from '../App/AppStateStore';
import {MemoryRouter} from 'react-router-dom';
import {MockWebServer} from './MockWebServer';
import {AppContext, appContext} from '../App/AppContext';
import * as Redux from 'redux';

type TestAppContextProps = {
    stateStore?: Redux.Store<AppState>
    oauthServer?: MockWebServer
    currentPath?: string
    children: React.ReactElement
}

export const TestAppContext = (props: TestAppContextProps): React.ReactElement => {
    const store = props.stateStore ?? appStateStore.create();

    const context = appContext.defaultEnv;
    if (props.oauthServer) {
        context.tokenUri = props.oauthServer.url('/token');
    }

    return <AppContext.Provider value={context}>
        <Provider store={store}>
            <MemoryRouter initialEntries={[props.currentPath ?? '/']}>
                {props.children}
            </MemoryRouter>
        </Provider>
    </AppContext.Provider>;
};
