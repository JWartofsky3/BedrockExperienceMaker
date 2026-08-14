import { useEffect, useState } from 'react'

type Addon = {
  name: string
  displayName: string
  creatorName?: string
  description?: string
  iconPath?: string
  curseforgeUrl?: string
  mcpedlUrl?: string
  currentVersion?: string
  minecraftVersionNote?: string
}

function App() {
  const [addons, setAddons] = useState<Addon[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string>()

  useEffect(() => {
    const controller = new AbortController()
    async function loadAddons() {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL ?? ''}/v1/addons`, { signal: controller.signal })
        if (!response.ok) throw new Error(`The catalog service returned ${response.status}.`)
        const body = (await response.json()) as { addons?: Addon[] }
        setAddons(body.addons ?? [])
      } catch (loadError) {
        if (controller.signal.aborted) return
        setError(loadError instanceof Error ? loadError.message : 'The catalog service is unavailable.')
      } finally {
        if (!controller.signal.aborted) setIsLoading(false)
      }
    }
    void loadAddons()
    return () => controller.abort()
  }, [])

  return (
    <main className="min-h-screen bg-slate-950 text-slate-100">
      <header className="border-b border-slate-800 bg-slate-900/70">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-emerald-400">Minecraft Bedrock</p>
            <h1 className="mt-1 text-xl font-bold">Experience Pack Builder</h1>
          </div>
          <span className="rounded-full bg-slate-800 px-3 py-1 text-sm text-slate-300">Add-on catalog</span>
        </div>
      </header>
      <section className="mx-auto max-w-6xl px-6 py-12">
        <div className="max-w-2xl">
          <p className="text-sm font-medium text-emerald-400">Browse add-ons</p>
          <h2 className="mt-2 text-3xl font-bold tracking-tight sm:text-4xl">Build your next Bedrock experience.</h2>
          <p className="mt-4 text-lg leading-8 text-slate-300">Choose from the curated catalog, then open an official provider page when you are ready to install.</p>
        </div>
        {isLoading && <p className="mt-8 text-slate-400">Loading add-ons…</p>}
        {error && <p className="mt-8 rounded-lg border border-red-400/30 bg-red-400/10 px-4 py-3 text-sm text-red-100">Could not load the add-on catalog. {error}</p>}
        {!isLoading && !error && addons.length === 0 && <p className="mt-8 text-slate-400">No add-ons have been added to the catalog yet.</p>}
        <div className="mt-8 grid gap-5 md:grid-cols-2 lg:grid-cols-3">
          {addons.map((addon) => (
            <article key={addon.name} className="overflow-hidden rounded-xl border border-slate-800 bg-slate-900 shadow-sm">
              {addon.iconPath ? (
                <img className="h-44 w-full bg-slate-800 object-cover" src={addon.iconPath} alt={`${addon.displayName} preview`} />
              ) : (
                <div className="h-44 bg-slate-800" />
              )}
              <div className="flex min-h-64 flex-col p-6">
                <div className="flex items-start justify-between gap-4">
                  <div><h3 className="text-xl font-semibold text-white">{addon.displayName}</h3><p className="mt-1 text-sm text-slate-400">by {addon.creatorName || 'Unknown creator'}</p></div>
                  {addon.currentVersion && <span className="shrink-0 rounded bg-emerald-400/10 px-2 py-1 text-xs font-medium text-emerald-300">v{addon.currentVersion}</span>}
                </div>
                <p className="mt-5 text-sm leading-6 text-slate-300">{addon.description || 'No description has been added yet.'}</p>
                {addon.minecraftVersionNote && <p className="mt-4 text-xs leading-5 text-amber-200">{addon.minecraftVersionNote}</p>}
                <div className="mt-auto flex gap-3 pt-6 text-sm font-medium">
                  {addon.curseforgeUrl && <a className="text-emerald-400 hover:text-emerald-300" href={addon.curseforgeUrl} target="_blank" rel="noreferrer">CurseForge</a>}
                  {addon.mcpedlUrl && <a className="text-emerald-400 hover:text-emerald-300" href={addon.mcpedlUrl} target="_blank" rel="noreferrer">MCPEDL</a>}
                </div>
              </div>
            </article>
          ))}
        </div>
      </section>
    </main>
  )
}

export default App
