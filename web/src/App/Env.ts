declare global {
    interface Window {
        env: Record<string, string>
    }
}

const requireEnv = (name: string): string => {
    const value = window.env[name];
    if (value === undefined) {
        throw `missing env configuration: ${name}`;
    }

    return value;
};

export const env = {
    require: requireEnv,
};
