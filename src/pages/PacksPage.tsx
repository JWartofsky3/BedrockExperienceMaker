import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { createPack, getCurrentUser, listPacks, type ExperiencePack, type User } from '../api'
import Layout from '../components/Layout'

export default function PacksPage() {
  const [packs, setPacks] = useState<ExperiencePack[]>([])
  const [error, setError] = useState<string>()
  const [saving, setSaving] = useState(false)
  const [user, setUser] = useState<User>()
  const navigate = useNavigate()

  useEffect(() => {
    void listPacks().then(setPacks).catch((reason: unknown) => setError(reason instanceof Error ? reason.message : 'Could not load packs.'))
    void getCurrentUser().then(setUser).catch(() => setUser(undefined))
  }, [])

  async function create(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    setSaving(true)
    setError(undefined)
    try {
      const pack = await createPack({ displayName: String(form.get('displayName')), description: String(form.get('description')), setupNotes: '' })
      navigate(`/packs/${pack.name.replace('packs/', '')}`)
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : 'Could not create pack.')
    } finally {
      setSaving(false)
    }
  }

  return <Layout><section className="mx-auto grid max-w-6xl gap-10 px-6 py-12 lg:grid-cols-[2fr_1fr]">
    <div><p className="text-sm font-medium text-emerald-400">Browse packs</p><h1 className="mt-2 text-3xl font-bold">Experience packs</h1>
      <div className="mt-8 space-y-3">{packs.map((pack) => <Link key={pack.name} to={`/packs/${pack.name.replace('packs/', '')}`} className="block rounded-xl border border-slate-800 bg-slate-900 p-5 hover:border-emerald-400/50"><h2 className="font-semibold">{pack.displayName}</h2><p className="mt-1 text-sm text-slate-400">Created by {pack.creatorName}</p>{pack.description && <p className="mt-3 text-sm text-slate-300">{pack.description}</p>}</Link>)}{!error && packs.length === 0 && <p className="text-slate-400">No packs have been created yet.</p>}</div>
    </div>
    {user ? <form onSubmit={create} className="h-fit rounded-xl border border-slate-800 bg-slate-900 p-6"><h2 className="text-lg font-semibold">Create a pack</h2><p className="mt-1 text-sm text-slate-400">Creating as {user.username}</p>{error && <p className="mt-3 text-sm text-red-300">{error}</p>}
      <label className="mt-5 block text-sm">Pack name<input required name="displayName" className="mt-1 w-full rounded border border-slate-700 bg-slate-950 p-2" /></label>
      <label className="mt-4 block text-sm">Description<textarea name="description" className="mt-1 w-full rounded border border-slate-700 bg-slate-950 p-2" rows={3} /></label>
      <button disabled={saving} className="mt-5 rounded bg-emerald-500 px-4 py-2 font-medium text-slate-950 disabled:opacity-50">{saving ? 'Creating…' : 'Create pack'}</button>
    </form> : <aside className="h-fit rounded-xl border border-slate-800 bg-slate-900 p-6"><h2 className="text-lg font-semibold">Create a pack</h2><p className="mt-2 text-sm text-slate-400">Sign in to create and manage your experience packs.</p><Link to="/login" className="mt-5 inline-block rounded bg-emerald-500 px-4 py-2 font-medium text-slate-950">Sign in</Link></aside>}
  </section></Layout>
}
