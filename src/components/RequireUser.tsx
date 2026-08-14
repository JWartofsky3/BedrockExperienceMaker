import { useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { getCurrentUser } from '../api'
import Layout from './Layout'

export default function RequireUser({ children }: { children: ReactNode }) {
  const [checking, setChecking] = useState(true)
  const [allowed, setAllowed] = useState(false)

  useEffect(() => {
    void getCurrentUser()
      .then(() => setAllowed(true))
      .catch(() => setAllowed(false))
      .finally(() => setChecking(false))
  }, [])

  if (checking) return <Layout><section className="mx-auto max-w-6xl px-6 py-12 text-slate-400">Checking access…</section></Layout>
  if (allowed) return <>{children}</>
  return <Layout><section className="mx-auto max-w-6xl px-6 py-12"><p className="text-sm font-medium text-red-300">403 Forbidden</p><h1 className="mt-2 text-3xl font-bold">This page is unavailable.</h1></section></Layout>
}
