import { type AxiosResponse } from 'axios'

export function unwrapData(response: AxiosResponse) {
  return response.data?.data ?? response.data
}
