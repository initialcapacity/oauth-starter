import {pkceSupport} from '../PKCESupport';

describe('PKCE Support', () => {
    test('challenge', () => {
        const verifier = pkceSupport.generateVerifier();
        const challenge = pkceSupport.generateChallenge(verifier);
        expect(challenge).toEqual(pkceSupport.generateChallenge(verifier));
    });
});
