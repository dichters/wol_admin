const BASE = '/wol/api'

async function request(url: string, options?: RequestInit) {
  const res = await fetch(`${BASE}${url}`, options)
  return res.json()
}

interface ApiResponse {
  code: number
  message: string
}

interface VersionResponse {
  code: number
  version: string
  arch: string
  build_time: string
}

export function fetchVersion(): Promise<VersionResponse> {
  return request('/version')
}

export function fetchWOL(): Promise<ApiResponse> {
  return request('/wol', { method: 'POST' })
}

export function fetchShutdown(): Promise<ApiResponse> {
  return request('/shutdown', { method: 'POST' })
}
