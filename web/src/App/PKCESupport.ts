import * as Crypto from 'crypto-js';

const generateVerifier = (): string => {
    // todo - improve randomness
    return safeEncode(Crypto.lib.WordArray.random(64));
};

const generateChallenge = (verifier: string): string =>
    safeEncode(Crypto.SHA256(verifier));

const safeEncode = (word: Crypto.lib.WordArray): string => {
    return word.toString(Crypto.enc.Base64url);
};

export const pkceSupport = {
    generateVerifier,
    generateChallenge,
};
