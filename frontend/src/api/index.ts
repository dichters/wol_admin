/**
 * API layer — unified request wrapper with client identification.
 * All requests include an X-Client-ID header for backend lock differentiation.
 */

const CLIENT_ID = generateClientID()

function generateClientID(): string {
  // Try to persist a stable ID in sessionStorage
  const KEY = 'wol_admin_client_id'
  let id = sessionStorage.getItem(KEY)
  if (!id) {
    id = crypto.randomUUID?.() ?? Math.random().toString(36).slice(2)
    sessionStorage.setItem(KEY, id)
  }
  return id
}

interface ApiResponse {
  code: number
  message: string
}

async function post(endpoint: string): Promise<ApiResponse> {
  const res = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Client-ID': CLIENT_ID,
    },
  })
  return res.json()
}

/** POST /api/wol — send Wake-on-LAN packet */
export const wol = () => post('/api/wol')

/** POST /api/shutdown — send SSH poweroff command */
export const shutdown = () => post('/api/shutdown')
