import { useEffect, useState } from 'react'
import type { DragEvent } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { addPackAddon, deletePack, downloadPack, getCurrentUser, getPack, listAddons, removePackAddon, reorderPackAddons, type Addon, type ExperiencePack, type User } from '../api'
import Layout from '../components/Layout'

export default function PackDetailPage() {
  const { packId = '' } = useParams()
  const navigate = useNavigate()
  const [pack, setPack] = useState<ExperiencePack>()
  const [catalog, setCatalog] = useState<Addon[]>([])
  const [user, setUser] = useState<User>()
  const [error, setError] = useState<string>()

  async function load() {
    try {
      const [loadedPack, loadedCatalog] = await Promise.all([getPack(packId), listAddons()])
      setPack(loadedPack)
      setCatalog(loadedCatalog)
      getCurrentUser().then(setUser).catch(() => setUser(undefined))
    } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not load pack.') }
  }
  useEffect(() => { void load() }, [packId])

  async function add(addonName: string) {
    try { setPack(await addPackAddon(packId, addonName)); setError(undefined) } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not add add-on.') }
  }
  async function remove(addonName: string) {
    try { await removePackAddon(packId, addonName.replace('addons/', '')); await load(); setError(undefined) } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not remove add-on.') }
  }
  async function reorder(draggedName: string, targetName: string) {
    if (!pack?.addons) return
    const reordered = [...pack.addons]
    const draggedIndex = reordered.findIndex((item) => item.addon.name === draggedName)
    const targetIndex = reordered.findIndex((item) => item.addon.name === targetName)
    if (draggedIndex < 0 || targetIndex < 0 || draggedIndex === targetIndex) return
    const [dragged] = reordered.splice(draggedIndex, 1)
    reordered.splice(draggedIndex < targetIndex ? targetIndex - 1 : targetIndex, 0, dragged)
    try { setPack(await reorderPackAddons(packId, reordered.map((item) => item.addon.name))); setError(undefined) } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not reorder add-ons.') }
  }
  function dragStart(event: DragEvent<HTMLLIElement>, addonName: string) {
    event.dataTransfer.setData('text/plain', addonName)
    event.dataTransfer.effectAllowed = 'move'
  }
  function drop(event: DragEvent<HTMLLIElement>, addonName: string) {
    event.preventDefault()
    void reorder(event.dataTransfer.getData('text/plain'), addonName)
  }
  async function removePack() {
    if (!window.confirm(`Delete ${pack?.displayName}?`)) return
    try { await deletePack(packId); navigate('/packs') } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not delete pack.') }
  }

  const selected = new Set(pack?.addons?.map((item) => item.addon.name))
  const canManage = Boolean(user && pack?.creatorUserId === user.name.replace('users/', ''))
  return <Layout><section className="mx-auto max-w-4xl px-6 py-12"><Link className="text-sm text-emerald-400" to="/packs">← All packs</Link>
    {error && <p className="mt-5 rounded bg-red-400/10 p-4 text-red-100">{error}</p>}
    {pack && <><div className="mt-5 flex items-start justify-between gap-6"><div><p className="text-sm text-emerald-400">Created by {pack.creatorName}</p><h1 className="mt-2 text-3xl font-bold">{pack.displayName}</h1>{pack.description && <p className="mt-3 text-slate-300">{pack.description}</p>}</div><div className="flex gap-3"><button onClick={() => downloadPack(pack)} className="rounded border border-emerald-400/50 px-3 py-2 text-sm text-emerald-300">Download JSON</button>{canManage && <button onClick={() => void removePack()} className="rounded border border-red-400/50 px-3 py-2 text-sm text-red-300">Delete pack</button>}</div></div>
      {!canManage && <p className="mt-5 rounded border border-slate-800 bg-slate-900 p-4 text-sm text-slate-300">{user ? 'Only this pack’s creator can make changes.' : <><Link className="text-emerald-400" to="/login">Sign in</Link> to create or edit your own packs.</>}</p>}
      <div className="mt-10 grid gap-8 lg:grid-cols-[3fr_2fr]"><div><h2 className="text-xl font-semibold">Install order</h2><p className="mt-1 text-sm text-slate-400">{canManage ? 'Drag add-ons into the order they should be installed.' : 'Add-ons at the top are installed first.'}</p><ol className="mt-4 space-y-3">{pack.addons?.map((item) => <li key={item.addon.name} draggable={canManage} onDragStart={(event) => dragStart(event, item.addon.name)} onDragOver={(event) => canManage && event.preventDefault()} onDrop={(event) => canManage && drop(event, item.addon.name)} className={`flex items-center gap-3 rounded-lg border border-slate-800 bg-slate-900 p-3 ${canManage ? 'cursor-grab' : ''}`}><span className="text-sm text-slate-500">{item.installOrder}</span><div className="min-w-0 flex-1"><p className="font-medium">{item.addon.displayName}</p><p className="text-xs text-slate-400">{item.addon.creatorName}</p>{item.addon.dependencies && item.addon.dependencies.length > 0 && <p className="mt-1 text-xs text-amber-200">Requires: {item.addon.dependencies.map((dependency) => dependency.displayName).join(', ')}</p>}</div>{canManage && <button onClick={() => void remove(item.addon.name)} className="text-sm text-red-300">Remove</button>}</li>)}</ol>{pack.addons?.length === 0 && <p className="mt-4 text-slate-400">Add add-ons from the catalog to begin.</p>}</div>
        {canManage && <aside><h2 className="text-xl font-semibold">Add an add-on</h2><div className="mt-4 space-y-2">{catalog.filter((addon) => !selected.has(addon.name)).map((addon) => <button key={addon.name} onClick={() => void add(addon.name)} className="w-full rounded-lg border border-slate-800 bg-slate-900 p-3 text-left hover:border-emerald-400/50"><p className="font-medium">{addon.displayName}</p><p className="text-xs text-slate-400">{addon.creatorName}</p>{addon.dependencies && addon.dependencies.length > 0 && <p className="mt-1 text-xs text-amber-200">Requires: {addon.dependencies.map((dependency) => dependency.displayName).join(', ')}</p>}</button>)}</div></aside>}
      </div>
    </>}
  </section></Layout>
}
