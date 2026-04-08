const LOCAL_LOGIN_KEY = "login";
const LOCAL_NAME_KEY = "name";
const LOCAL_PASSWORD_KEY = "password";

type StoredCredentials = {
  login: string;
  name: string;
  password: string;
};

type ApiErrorResponse = {
  err?: string;
};

const LOGIN_CHARS = "abcdefghijklmnopqrstuvwxyz0123456789";
const PASSWORD_CHARS =
  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

const randomString = (length: number, chars: string): string => {
  const bytes = new Uint8Array(length);
  crypto.getRandomValues(bytes);
  let result = "";
  for (const b of bytes) {
    result += chars[b % chars.length];
  }
  return result;
};

const makeRandomCredentials = (): StoredCredentials => {
  const suffix = randomString(12, LOGIN_CHARS);
  return {
    login: `player${suffix}`.slice(0, 20),
    name: `Player${suffix}`.slice(0, 20),
    password: randomString(16, PASSWORD_CHARS),
  };
};

const readStoredCredentials = (): StoredCredentials | null => {
  const login = localStorage.getItem(LOCAL_LOGIN_KEY);
  const name = localStorage.getItem(LOCAL_NAME_KEY);
  const password = localStorage.getItem(LOCAL_PASSWORD_KEY);
  if (!login || !name || !password) {
    return null;
  }
  return { login, name, password };
};

const saveCredentials = (creds: StoredCredentials) => {
  localStorage.setItem(LOCAL_LOGIN_KEY, creds.login);
  localStorage.setItem(LOCAL_NAME_KEY, creds.name);
  localStorage.setItem(LOCAL_PASSWORD_KEY, creds.password);
};

const parseApiError = async (response: Response): Promise<string> => {
  try {
    const data = (await response.json()) as ApiErrorResponse;
    return data.err ?? `${response.status} ${response.statusText}`;
  } catch {
    return `${response.status} ${response.statusText}`;
  }
};

const checkProtectedOk = async (httpBase: string): Promise<boolean> => {
  const resp = await fetch(`${httpBase}/api/protected/ok`, {
    method: "GET",
    credentials: "include",
  });
  return resp.ok;
};

const loginWithCredentials = async (
  httpBase: string,
  creds: StoredCredentials,
): Promise<void> => {
  const resp = await fetch(`${httpBase}/api/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      login: creds.login,
      password: creds.password,
    }),
  });
  if (!resp.ok) {
    throw new Error(await parseApiError(resp));
  }
};

const registerWithCredentials = async (
  httpBase: string,
  creds: StoredCredentials,
): Promise<void> => {
  const resp = await fetch(`${httpBase}/api/register`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name: creds.name,
      login: creds.login,
      password: creds.password,
    }),
  });
  if (!resp.ok) {
    throw new Error(await parseApiError(resp));
  }
};

export const ensureAuthenticated = async (
  httpBase: string,
): Promise<StoredCredentials> => {
  if (await checkProtectedOk(httpBase)) {
    const existing = readStoredCredentials();
    if (existing) return existing;
  }

  const stored = readStoredCredentials();
  if (stored) {
    await loginWithCredentials(httpBase, stored);
    return stored;
  }

  for (let i = 0; i < 3; i++) {
    const generated = makeRandomCredentials();
    try {
      await registerWithCredentials(httpBase, generated);
      saveCredentials(generated);
      return generated;
    } catch (e) {
      if (i === 2) {
        throw e;
      }
    }
  }

  throw new Error("Failed to initialize account.");
};

