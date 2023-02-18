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

