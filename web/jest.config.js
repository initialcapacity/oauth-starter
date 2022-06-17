/** @type {import('ts-jest/dist/types').InitialOptionsTsJest} */
module.exports = {
    preset: 'ts-jest',
    testEnvironment: 'jsdom',
    globals: {'ts-jest': {useESM: true}},
    moduleNameMapper: {
        '\\.css$': '<rootDir>/src/TestSupport/AssetsStubs.js',
        '\\.svg$': '<rootDir>/src/TestSupport/AssetsStubs.js',
    },
    setupFilesAfterEnv: [
        '<rootDir>/src/TestSupport/JestSetup.ts',
    ]
};
