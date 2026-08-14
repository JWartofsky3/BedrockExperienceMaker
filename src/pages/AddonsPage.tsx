import { useEffect, useState } from 'react'
import { listAddons, type Addon } from '../api'
import Layout from '../components/Layout'

export default function AddonsPage() {
  const [addons, setAddons] = useState<Addon[]>([])
  const [error, setError] = useState<string>()

  useEffect(() => { void listAddons().then(setAddons).catch((reason: unknown) => setError(reason instanceof Error ? reason.message : 'Could not load add-ons.')) }, [])

  return <Layout><section className="mx-auto max-w-6xl px-6 py-12">
    <p className="text-sm font-medium text-emerald-400">Browse add-ons</p>
    <h1 className="mt-2 text-3xl font-bold">Add-on catalog</h1>
    {error && <p className="mt-8 rounded-lg bg-red-400/10 p-4 text-red-100">{error}</p>}
    <div className="mt-8 grid gap-5 md:grid-cols-2 lg:grid-cols-3">
      {addons.map((addon) => <article key={addon.name} className="overflow-hidden rounded-xl border border-slate-800 bg-slate-900">
        {addon.iconPath ? <img className="h-44 w-full bg-slate-800 object-cover" src={addon.iconPath} alt={`${addon.displayName} preview`} /> : <div className="h-44 bg-slate-800" />}
        <div className="flex min-h-56 flex-col p-6">
          <div className="flex justify-between gap-3"><div><h2 className="text-xl font-semibold">{addon.displayName}</h2><p className="text-sm text-slate-400">by {addon.creatorName || 'Unknown creator'}</p></div>{addon.currentVersion && <span className="h-fit rounded bg-emerald-400/10 px-2 py-1 text-xs text-emerald-300">v{addon.currentVersion}</span>}</div>
          <p className="mt-4 text-sm text-slate-300">{addon.description}</p>
          <div className="mt-auto flex gap-3 pt-5 text-sm text-emerald-400">{addon.curseforgeUrl && <a href={addon.curseforgeUrl} target="_blank" rel="noreferrer">CurseForge</a>}{addon.mcpedlUrl && <a href={addon.mcpedlUrl} target="_blank" rel="noreferrer">MCPEDL</a>}</div>
        </div>
      </article>)}
    </div>
  </section></Layout>
}
