import { Link, NavLink } from 'react-router-dom'
import type { ReactNode } from 'react'

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <main className="min-h-screen bg-slate-950 text-slate-100">
      <header className="border-b border-slate-800 bg-slate-900/70">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
          <Link to="/experiences" className="font-bold">Experience Pack Builder</Link>
          <nav className="flex gap-4 text-sm font-medium text-slate-300">
            <NavLink className={({ isActive }) => isActive ? 'text-emerald-400' : 'hover:text-white'} to="/experiences">Experiences</NavLink>
            <NavLink className={({ isActive }) => isActive ? 'text-emerald-400' : 'hover:text-white'} to="/addons">Add-ons</NavLink>
            <NavLink className={({ isActive }) => isActive ? 'text-emerald-400' : 'hover:text-white'} to="/packs">Manage packs</NavLink>
          </nav>
        </div>
      </header>
      {children}
    </main>
  )
}
