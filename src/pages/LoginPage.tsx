import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../api'
import Layout from '../components/Layout'

export default function LoginPage() {
  const [error, setError] = useState<string>()
  const [saving, setSaving] = useState(false)
  const navigate = useNavigate()

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    setSaving(true)
    setError(undefined)
    try {
      await login(String(form.get('username')), String(form.get('password')))
      navigate('/packs')
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : 'Could not sign in.')
    } finally {
      setSaving(false)
    }
  }

  return <Layout><section className="mx-auto max-w-md px-6 py-16"><form onSubmit={submit} className="rounded-xl border border-slate-800 bg-slate-900 p-6"><p className="text-sm font-medium text-emerald-400">Creator access</p><h1 className="mt-2 text-2xl font-bold">Sign in</h1>{error && <p className="mt-4 rounded bg-red-400/10 p-3 text-sm text-red-100">{error}</p>}
    <label className="mt-6 block text-sm">Username<input required name="username" autoComplete="username" className="mt-1 w-full rounded border border-slate-700 bg-slate-950 p-2" /></label>
    <label className="mt-4 block text-sm">Password<input required name="password" type="password" autoComplete="current-password" className="mt-1 w-full rounded border border-slate-700 bg-slate-950 p-2" /></label>
    <button disabled={saving} className="mt-6 rounded bg-emerald-500 px-4 py-2 font-medium text-slate-950 disabled:opacity-50">{saving ? 'Signing in…' : 'Sign in'}</button>
  </form></section></Layout>
}
