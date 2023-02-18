import _isEmpty from 'lodash/isEmpty';
import jwt_decode from "jwt-decode";

export namespace Auth {

  export interface Metadata {
    iss: string
    sub: string
    email: string
    iat: number
  }

  export interface User {
    id: string
    email: string
    name: string
    labels: {[key: string]: string}
  }

}

export const LabelCurrency = "currency";
export const LabelLocale = "locale";
export const LabelUserGoogleImage = "google|picture";

const AuthTokenKey = "tiinyplanet.com:auth:token"
const UserTokenKey = "tiinyplanet.com:auth:user"

export const persistAuthToken = (tkn: string) => {
  localStorage.setItem(AuthTokenKey, tkn);
}

export const readAuthToken = () => {
  return localStorage.getItem(AuthTokenKey) || "";
}

export const deleteAuthToken = () => {
  localStorage.removeItem(AuthTokenKey);
}

export const readAuthMetadata = (): Auth.Metadata | undefined => {
  const tkn = readAuthToken();
  if (_isEmpty(tkn)) {
    return undefined;
  }
  return jwt_decode(tkn);
}

export const persistAuthUser = (user: Auth.User) => {
  localStorage.setItem(UserTokenKey, JSON.stringify(user));
}

export const deleteAuthUser = () => {
  localStorage.removeItem(UserTokenKey)
}

export const readAuthUser = (): Auth.User|null => {
  const json = localStorage.getItem(UserTokenKey);
  if (json) {
    return JSON.parse(json);
  }
  return null;
}

export const readUserLocale = (): string => {
  return localStorage.getItem("i18nextLng") || "en";
}
