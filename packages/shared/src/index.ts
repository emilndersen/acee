// shared types между фронтом и бэком
export type Photo = {
  id: string
  title: string
  album: string
  description: string
  url: string
  sort_order: number
  created_at: string
}

export type Booking = {
  id: string
  name: string
  contact: string
  shoot_type: string
  date: string
  idea: string
  created_at: string
}
