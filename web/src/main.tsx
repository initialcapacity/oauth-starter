import React from 'react';
import ReactDOM from 'react-dom/client';
import {App} from './App/App';
import {appContext, AppContext} from './App/AppContext';
import {Provider} from 'react-redux';
import {appStateStore} from './App/AppStateStore';
import {HashRouter} from 'react-router-dom';

ReactDOM.createRoot(document.getElementById('root') || document.body).render(
    <React.StrictMode>
        <AppContext.Provider value={appContext.defaultEnv}>
            <Provider store={appStateStore.create()}>
                <HashRouter>
                    <App/>
                </HashRouter>
            </Provider>
        </AppContext.Provider>
    </React.StrictMode>
);
