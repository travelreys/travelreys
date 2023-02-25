import axios from 'axios';
import _get from 'lodash/get';
import { User } from '../lib/auth';
import { Labels } from '../lib/common';

import { BASE_URL, makeCommonAxios } from './common';

export interface  LoginResponse {
  jwtToken?: string
  error?: string
}

export interface ReadResponse {
  user?: User
  error?: string
}

export interface SearchUsersResponse {
  users?: Array<User>
  error?: string
}

export interface UpdateResponse {
  error?: string
}

const authLoginPathPrefix = "/api/v1/auth/login";
const authUserPathPrefix = "/api/v1/auth/users";


const login = async (code: string): Promise<LoginResponse> => {
  const url = `${BASE_URL}${authLoginPathPrefix}`;
  return axios.post(url, { code })
    .then((res) => {
      const token = _get(res, "data.jwtToken", "");
      return {jwtToken:token};
    })
    .catch((err) => {
      return {error: err.message};
    });
}

const read = (id: string): Promise<ReadResponse> => {
  return  makeCommonAxios().get(`${authUserPathPrefix}/${id}`)
    .then((res) => {
      const user = _get(res, "data.user", {}) as User;
      return {user: user}
    })
    .catch((err) => {
      return {error: err.message};
    })
}

const update = (usrID: string, ff: Labels): Promise<UpdateResponse> => {
  return makeCommonAxios().put(`${authUserPathPrefix}/${usrID}`, {ff})
    .then((res) => {
      return {error: undefined}
    })
    .catch((err) => {
      return {error: err.message};
    })
}

const searchUsers = (email: string): Promise<SearchUsersResponse> => {
  return  makeCommonAxios().get(authUserPathPrefix, {params: {email}})
    .then((res) => {
      const users = _get(res, "data.users", []);
      return {users}
    })
    .catch((err) => {
      return {error: err.message};
    })
}

export default {
  login,
  read,
  searchUsers,
  update
};
