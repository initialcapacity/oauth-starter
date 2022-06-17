import * as Redux from 'redux';
import {Store} from 'redux';

export declare namespace AppState {
    type AuthStatus =
        | { signedIn: true, accessToken: string }
        | { signedIn: false }

    type Action =
        | { type: 'count/increment' }
        | { type: 'auth/sign in', accessToken: string }
}

export type AppState = {
    authStatus: AppState.AuthStatus;
    count: number
}

const incrementCount: AppState.Action =
    ({type: 'count/increment'});

const signIn = (accessToken: string): AppState.Action =>
    ({type: 'auth/sign in', accessToken});

const initialState: AppState = {
    authStatus: {signedIn: false},
    count: 0,
};

const reducer: Redux.Reducer<AppState, AppState.Action> =
    (state = initialState, action): AppState => {
        switch (action.type) {
            case 'count/increment':
                return {...state, count: state.count + 1};
            case 'auth/sign in':
                return {...state, authStatus: {signedIn: true, accessToken: action.accessToken}};
        }
        return state;
    };

const create = (): Store<AppState> =>
    Redux.createStore(reducer);

export const appStateStore = {
    incrementCount,
    signIn,
    create,
};
