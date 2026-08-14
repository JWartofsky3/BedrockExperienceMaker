import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { downloadPack, getPack, type ExperiencePack } from '../api'
import Layout from '../components/Layout'

export default function ExperienceDetailPage() {
  const { packId = '' } = useParams()
  const [pack, setPack] = useState<ExperiencePack>()
  const [error, setError] = useState<string>()

  useEffect(() => { void getPack(packId).then(setPack).catch((reason: unknown) => setError(reason instanceof Error ? reason.message : 'Could not load this experience.')) }, [packId])

  return <Layout><section className="mx-auto max-w-5xl px-6 py-12">
    <Link className="text-sm text-emerald-400" to="/experiences">← All experiences</Link>
    {error && <p className="mt-5 rounded-lg bg-red-400/10 p-4 text-red-100">{error}</p>}
    {pack && <><div className="mt-6 flex items-start justify-between gap-6"><div><p className="text-sm text-emerald-400">Created by {pack.creatorName}</p><h1 className="mt-2 text-3xl font-bold">{pack.displayName}</h1></div><button onClick={() => downloadPack(pack)} className="rounded border border-emerald-400/50 px-3 py-2 text-sm text-emerald-300 hover:bg-emerald-400/10">Download JSON</button></div>
      {pack.description && <p className="mt-4 max-w-3xl text-slate-300">{pack.description}</p>}
      {pack.setupNotes && <section className="mt-8 rounded-xl border border-amber-400/30 bg-amber-400/10 p-5"><h2 className="font-semibold text-amber-100">Setup notes</h2><p className="mt-2 text-sm leading-6 text-amber-50">{pack.setupNotes}</p></section>}
      <section className="mt-10"><h2 className="text-2xl font-semibold">Installed add-ons</h2><p className="mt-2 text-sm text-slate-400">Install them in this order.</p>
        <ol className="mt-5 space-y-4">{pack.addons?.map((item) => <li key={item.addon.name} className="flex gap-4 rounded-xl border border-slate-800 bg-slate-900 p-4"><span className="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-emerald-400/10 text-sm font-semibold text-emerald-300">{item.installOrder}</span>{item.addon.iconPath && <img className="h-16 w-16 rounded-lg bg-slate-800 object-cover" src={item.addon.iconPath} alt="" />}<div className="min-w-0 flex-1"><h3 className="font-semibold">{item.addon.displayName}</h3><p className="text-sm text-slate-400">by {item.addon.creatorName}</p><p className="mt-2 text-sm text-slate-300">{item.addon.description}</p><div className="mt-3 flex gap-3 text-sm text-emerald-400">{item.addon.curseforgeUrl && <a href={item.addon.curseforgeUrl} target="_blank" rel="noreferrer">CurseForge</a>}{item.addon.mcpedlUrl && <a href={item.addon.mcpedlUrl} target="_blank" rel="noreferrer">MCPEDL</a>}</div></div></li>)}</ol>
      </section>
    </>}
  </section></Layout>
}
