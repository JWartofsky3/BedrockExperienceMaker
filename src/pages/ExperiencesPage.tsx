import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { listPacks, type ExperiencePack } from '../api'
import Layout from '../components/Layout'

export default function ExperiencesPage() {
  const [packs, setPacks] = useState<ExperiencePack[]>([])
  const [error, setError] = useState<string>()

  useEffect(() => { void listPacks().then(setPacks).catch((reason: unknown) => setError(reason instanceof Error ? reason.message : 'Could not load experiences.')) }, [])

  return <Layout><section className="mx-auto max-w-6xl px-6 py-12">
    <p className="text-sm font-medium text-emerald-400">Minecraft Bedrock</p>
    <h1 className="mt-2 text-3xl font-bold">Experience packs</h1>
    <p className="mt-3 max-w-2xl text-slate-300">Browse curated add-on collections and follow their recommended installation order.</p>
    {error && <p className="mt-8 rounded-lg bg-red-400/10 p-4 text-red-100">{error}</p>}
    <div className="mt-8 grid gap-5 md:grid-cols-2 lg:grid-cols-3">
      {packs.map((pack) => <Link key={pack.name} to={`/experiences/${pack.name.replace('packs/', '')}`} className="rounded-xl border border-slate-800 bg-slate-900 p-6 hover:border-emerald-400/50">
        <p className="text-sm text-emerald-400">Created by {pack.creatorName}</p>
        <h2 className="mt-2 text-xl font-semibold">{pack.displayName}</h2>
        <p className="mt-3 text-sm text-slate-300">{pack.description || 'A curated Bedrock add-on collection.'}</p>
      </Link>)}
    </div>
    {!error && packs.length === 0 && <p className="mt-8 text-slate-400">No experience packs have been published yet.</p>}
  </section></Layout>
}
