import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { addPackAddon, deletePack, getPack, listAddons, removePackAddon, reorderPackAddons, type Addon, type ExperiencePack } from '../api'
import Layout from '../components/Layout'

export default function PackDetailPage() {
  const { packId = '' } = useParams()
  const navigate = useNavigate()
  const [pack, setPack] = useState<ExperiencePack>()
  const [catalog, setCatalog] = useState<Addon[]>([])
  const [error, setError] = useState<string>()

  async function load() {
    try { setPack(await getPack(packId)); setCatalog(await listAddons()) } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not load pack.') }
  }
  useEffect(() => { void load() }, [packId])

  async function add(addonName: string) {
    try { setPack(await addPackAddon(packId, addonName)); setError(undefined) } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not add add-on.') }
  }
  async function remove(addonName: string) {
    try { await removePackAddon(packId, addonName.replace('addons/', '')); await load(); setError(undefined) } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not remove add-on.') }
  }
  async function move(index: number, direction: -1 | 1) {
    if (!pack?.addons) return
    const reordered = [...pack.addons]
    const next = index + direction
    if (next < 0 || next >= reordered.length) return
    ;[reordered[index], reordered[next]] = [reordered[next], reordered[index]]
    try { setPack(await reorderPackAddons(packId, reordered.map((item) => item.addon.name))); setError(undefined) } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not reorder add-ons.') }
  }
  async function removePack() {
    if (!window.confirm(`Delete ${pack?.displayName}?`)) return
    try { await deletePack(packId); navigate('/packs') } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not delete pack.') }
  }

  const selected = new Set(pack?.addons?.map((item) => item.addon.name))
  return <Layout><section className="mx-auto max-w-4xl px-6 py-12"><Link className="text-sm text-emerald-400" to="/packs">← All packs</Link>
    {error && <p className="mt-5 rounded bg-red-400/10 p-4 text-red-100">{error}</p>}
    {pack && <><div className="mt-5 flex items-start justify-between gap-6"><div><p className="text-sm text-emerald-400">Created by {pack.creatorName}</p><h1 className="mt-2 text-3xl font-bold">{pack.displayName}</h1>{pack.description && <p className="mt-3 text-slate-300">{pack.description}</p>}</div><button onClick={() => void removePack()} className="rounded border border-red-400/50 px-3 py-2 text-sm text-red-300">Delete pack</button></div>
      <div className="mt-10 grid gap-8 lg:grid-cols-[3fr_2fr]"><div><h2 className="text-xl font-semibold">Install order</h2><p className="mt-1 text-sm text-slate-400">Add-ons at the top are installed first.</p><ol className="mt-4 space-y-3">{pack.addons?.map((item, index) => <li key={item.addon.name} className="flex items-center gap-3 rounded-lg border border-slate-800 bg-slate-900 p-3"><span className="text-sm text-slate-500">{item.installOrder}</span><div className="min-w-0 flex-1"><p className="font-medium">{item.addon.displayName}</p><p className="text-xs text-slate-400">{item.addon.creatorName}</p></div><button onClick={() => void move(index, -1)} disabled={index === 0} className="text-sm text-emerald-400 disabled:text-slate-600">Up</button><button onClick={() => void move(index, 1)} disabled={index === (pack.addons?.length ?? 0) - 1} className="text-sm text-emerald-400 disabled:text-slate-600">Down</button><button onClick={() => void remove(item.addon.name)} className="text-sm text-red-300">Remove</button></li>)}</ol>{pack.addons?.length === 0 && <p className="mt-4 text-slate-400">Add add-ons from the catalog to begin.</p>}</div>
        <aside><h2 className="text-xl font-semibold">Add an add-on</h2><div className="mt-4 space-y-2">{catalog.filter((addon) => !selected.has(addon.name)).map((addon) => <button key={addon.name} onClick={() => void add(addon.name)} className="w-full rounded-lg border border-slate-800 bg-slate-900 p-3 text-left hover:border-emerald-400/50"><p className="font-medium">{addon.displayName}</p><p className="text-xs text-slate-400">{addon.creatorName}</p></button>)}</div></aside>
      </div>
    </>}
  </section></Layout>
}
