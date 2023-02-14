export namespace Common {
  export interface PriceMetadata {
    amount?: number
    currency: string
  }
  export interface Positioning {
    name: string
    address?: string
    continent?: string
    country?: string
    state?: string
    city?: string
    longitude?: string
    latitude?: string
  }
}

export const BASE_URL = "http://localhost:2022";
export const BASE_WS_URL = "ws://localhost:2022/ws";
