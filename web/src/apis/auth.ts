import axios from 'axios';
import _get from 'lodash/get';
import { Auth } from '../lib/auth';

import { BASE_URL, makeCommonAxios } from './common';

export interface  LoginResponse {
  jwtToken?: string
  error?: string
}

export interface ReadResponse {
  user?: Auth.User
  error?: string
}

export interface SearchUsersResponse {
  users?: Array<Auth.User>
  error?: string
}

export interface UpdateFilter {
  labels: {[key: string]: string}
}

export const makeUpdateFilter = (labels: {[key: string]: string}): UpdateFilter => {
  return {labels}
}

export interface UpdateResponse {
  error?: string
}

const authLoginPathPrefix = "/api/v1/auth/login";
const authUserPathPrefix = "/api/v1/auth/users";

const AuthAPI = {

  login: async (code: string): Promise<LoginResponse> => {
    const url = `${BASE_URL}${authLoginPathPrefix}`;
    return axios.post(url, { code })
      .then((res) => {
        const token = _get(res, "data.jwtToken", "");
        return {jwtToken:token};
      })
      .catch((err) => {
        return {error: err.message};
      });
  },

  readUser: (usrID: string): Promise<ReadResponse> => {
    const ax = makeCommonAxios();
    return ax.get(`${authUserPathPrefix}/${usrID}`)
      .then((res) => {
        const user = _get(res, "data.user", {}) as Auth.User;
        return {user: user}
      })
      .catch((err) => {
        return {error: err.message};
      })
  },

  searchUsers: (email: string): Promise<SearchUsersResponse> => {
    const ax = makeCommonAxios();
    return ax.get(`${authUserPathPrefix}`, {params: {email}})
      .then((res) => {
        const users = _get(res, "data.users", []);
        return {users}
      })
      .catch((err) => {
        return {error: err.message};
      })
  },

  updateUser: (usrID: string, ff: UpdateFilter): Promise<UpdateResponse> => {
    const ax = makeCommonAxios();
    return ax.put(`${authUserPathPrefix}/${usrID}`, {ff})
      .then((res) => {
        return {error: undefined}
      })
      .catch((err) => {
        return {error: err.message};
      })
  }

};

export default AuthAPI;
