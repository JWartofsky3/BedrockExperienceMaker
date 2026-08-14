import { Link, NavLink, useNavigate } from 'react-router-dom'
import type { ReactNode } from 'react'
import { useEffect, useState } from 'react'
import { getCurrentUser, logout, type User } from '../api'

export default function Layout({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User>()
  const navigate = useNavigate()

  useEffect(() => { void getCurrentUser().then(setUser).catch(() => setUser(undefined)) }, [])

  async function signOut() {
    await logout().catch(() => undefined)
    setUser(undefined)
    navigate('/experiences')
  }

  return (
    <main className="min-h-screen bg-slate-950 text-slate-100">
      <header className="border-b border-slate-800 bg-slate-900/70">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
          <Link to="/experiences" className="font-bold">Experience Pack Builder</Link>
          <nav className="flex items-center gap-4 text-sm font-medium text-slate-300">
            <NavLink className={({ isActive }) => isActive ? 'text-emerald-400' : 'hover:text-white'} to="/experiences">Experiences</NavLink>
            <NavLink className={({ isActive }) => isActive ? 'text-emerald-400' : 'hover:text-white'} to="/addons">Add-ons</NavLink>
            {user ? <><NavLink className={({ isActive }) => isActive ? 'text-emerald-400' : 'hover:text-white'} to="/packs">Manage packs</NavLink><span className="text-slate-400">{user.username}</span><button onClick={() => void signOut()} className="hover:text-white">Sign out</button></> : <NavLink className={({ isActive }) => isActive ? 'text-emerald-400' : 'hover:text-white'} to="/login">Sign in</NavLink>}
          </nav>
        </div>
      </header>
      {children}
    </main>
  )
}
