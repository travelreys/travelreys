export namespace Common {
  export interface Price {
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

