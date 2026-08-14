export type Addon = {
  name: string
  displayName: string
  creatorName?: string
  description?: string
  iconPath?: string
  curseforgeUrl?: string
  mcpedlUrl?: string
  currentVersion?: string
  minecraftVersionNote?: string
  dependencies?: AddonDependency[]
}

export type AddonDependency = {
  name: string
  displayName: string
}

export type PackAddon = {
  addon: Addon
  installOrder: number
  note?: string
}

export type ExperiencePack = {
  name: string
  displayName: string
  creatorName: string
  creatorUserId?: string
  description?: string
  setupNotes?: string
  addons?: PackAddon[]
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${import.meta.env.VITE_API_BASE_URL ?? ''}${path}`, {
    credentials: 'same-origin',
    headers: { 'Content-Type': 'application/json', ...init?.headers },
    ...init,
  })
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { error?: string } | null
    throw new Error(body?.error || `The catalog service returned ${response.status}.`)
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}

export async function listAddons() {
  return (await request<{ addons: Addon[] }>('/v1/addons')).addons
}

export async function listPacks() {
  return (await request<{ experiencePacks: ExperiencePack[] }>('/v1/packs')).experiencePacks
}

export function getPack(id: string) {
  return request<ExperiencePack>(`/v1/packs/${id}`)
}

export function downloadPack(pack: ExperiencePack) {
  const fileName = `${pack.displayName.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '') || 'experience-pack'}.json`
  const content = JSON.stringify({ ...pack, addons: pack.addons ?? [] }, null, 2)
  const url = URL.createObjectURL(new Blob([content], { type: 'application/json' }))
  const link = document.createElement('a')
  link.href = url
  link.download = fileName
  link.click()
  URL.revokeObjectURL(url)
}

export type User = {
  name: string
  username: string
}

export function getCurrentUser() {
  return request<User>('/v1/auth/me')
}

export function login(username: string, password: string) {
  return request<User>('/v1/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) })
}

export function logout() {
  return request<void>('/v1/auth/logout', { method: 'POST' })
}

export function createPack(pack: Pick<ExperiencePack, 'displayName' | 'description' | 'setupNotes'>) {
  return request<ExperiencePack>('/v1/packs', { method: 'POST', body: JSON.stringify(pack) })
}

export function deletePack(id: string) {
  return request<void>(`/v1/packs/${id}`, { method: 'DELETE' })
}

export function addPackAddon(id: string, addonName: string) {
  return request<ExperiencePack>(`/v1/packs/${id}:addAddon`, { method: 'POST', body: JSON.stringify({ addonName }) })
}

export function removePackAddon(id: string, addonID: string) {
  return request<void>(`/v1/packs/${id}/addons/${addonID}`, { method: 'DELETE' })
}

export function reorderPackAddons(id: string, addonNames: string[]) {
  return request<ExperiencePack>(`/v1/packs/${id}:reorderAddons`, { method: 'POST', body: JSON.stringify({ addonNames }) })
}
