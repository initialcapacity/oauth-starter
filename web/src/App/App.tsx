import * as ReactRouter from 'react-router';
import * as ReactRedux from 'react-redux';
import * as React from 'react';
import './App.css';
import {AppState, appStateStore} from './AppStateStore';
import {Layout} from './Layout';
import {Route, Routes} from 'react-router-dom';
import {AppContext} from './AppContext';
import {pkceSupport} from './PKCESupport';
import {isMatching, P} from 'ts-pattern';

const SignedOutApp = (): React.ReactElement => {
    return <section>
        <div className="container">
            <Routes>
                <Route path={'/'} element={<Redirect to={'/login'}/>}/>
                <Route path={'/login'} element={<LoginPage/>}/>
                <Route path={'/callback'} element={<CallbackPage/>}/>
            </Routes>
        </div>
    </section>;
};

const SignedInApp = (): React.ReactElement => {
    const dispatch = ReactRedux.useDispatch();
    const count = ReactRedux.useSelector((app: AppState) => app.count);
    return <section>
        <div className="container">
            <p>
                <button type="button" onClick={() => dispatch(appStateStore.incrementCount)}>
                    count is: {count}
                </button>
            </p>
        </div>
    </section>;
};

const Redirect = ({to}: { to: string }): React.ReactElement => {
    const navigate = ReactRouter.useNavigate();
    React.useEffect(() => {
        navigate(to);
    });
    return <></>;
};

const LoginPage = (): React.ReactElement => {
    const context = React.useContext(AppContext);
    const verifier = pkceSupport.generateVerifier();
    const challenge = pkceSupport.generateChallenge(verifier);
    const query = new URLSearchParams([
        ['client_id', context.clientId],
        ['redirect_url', context.callbackUrl],
        ['response_type', 'code'],
        ['code_challenge', challenge],
        ['code_challenge_method', 'S256'],
        ['scope', 'openid email'],
    ]).toString();

    React.useEffect(() => {
        localStorage.setItem('verifier', verifier);
    }, [verifier, challenge]);

    return <a href={`${context.authUri}?${query}`}>
        Sign in with the Authorization Server
    </a>;
};

type TokenResponseJson = {
    access_token: string
}

const tokenResponseJsonPattern: P.Pattern<TokenResponseJson> = {
    access_token: P.string,
};

const CallbackPage = (): React.ReactElement => {
    const context = React.useContext(AppContext);
    const dispatch = ReactRedux.useDispatch();
    const navigate = ReactRouter.useNavigate();

    React.useEffect(() => {
        const code_verifier = localStorage.getItem('verifier');
        if (code_verifier === null) {
            return;
        }
        localStorage.removeItem('verifier');
        const query = new URLSearchParams([
            ['grant_type', 'authorization_code'],
            ['code', '42'], // todo - get the code
            ['client_id', context.clientId],
            ['redirect_url', context.callbackUrl],
            ['code_verifier', code_verifier],
        ]).toString();

        const requestInit: RequestInit = {
            method: 'POST',
            headers: {'Content-Type': 'application/x-www-form-urlencoded'},
            body: query,
        };

        fetch(context.tokenUri, requestInit)
            .then(r => r.json())
            .catch(() => null)
            .then(json => {
                if (!isMatching(tokenResponseJsonPattern, json)) {
                    return;
                }
                const accessToken = json.access_token;
                dispatch(appStateStore.signIn(accessToken));
                navigate('/');
            });
    }, []);

    return <></>;
};

export const App = () => {
    const authStatus = ReactRedux.useSelector((app: AppState) => app.authStatus);
    if (authStatus.signedIn) {
        return <Layout><SignedInApp/></Layout>;
    }
    return <Layout><SignedOutApp/></Layout>;
};
