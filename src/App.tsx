import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import AddonsPage from './pages/AddonsPage'
import ExperienceDetailPage from './pages/ExperienceDetailPage'
import ExperiencesPage from './pages/ExperiencesPage'
import PackDetailPage from './pages/PackDetailPage'
import PacksPage from './pages/PacksPage'

export default function App() {
  return <BrowserRouter><Routes>
    <Route path="/experiences" element={<ExperiencesPage />} />
    <Route path="/experiences/:packId" element={<ExperienceDetailPage />} />
    <Route path="/addons" element={<AddonsPage />} />
    <Route path="/packs" element={<PacksPage />} />
    <Route path="/packs/:packId" element={<PackDetailPage />} />
    <Route path="*" element={<Navigate to="/experiences" replace />} />
  </Routes></BrowserRouter>
}
