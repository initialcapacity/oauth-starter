import * as React from 'react';
import {env} from './Env';

const defaultEnv = {
    clientId: env.require('clientId'),
    authUri: env.require('authUri'),
    tokenUri: env.require('tokenUri'),
    callbackUrl: env.require('callbackUrl'),
    oauthResourceServer: env.require('oauthResourceServer'),
};

export const AppContext =
    React.createContext(defaultEnv);

export const appContext = {
    defaultEnv,
};
