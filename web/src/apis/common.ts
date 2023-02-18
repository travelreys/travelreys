import axios from 'axios';
import { readAuthToken } from '../lib/auth';

export const BASE_URL = "http://localhost:2022";
export const BASE_WS_URL = "ws://localhost:2022/ws";

export const makeCommonAxios = () => {
  const ax = axios.create({
    baseURL: BASE_URL,
  });

  ax.defaults.headers.common['Authorization'] = `Bearer ${readAuthToken()}`;
  return ax;
}
