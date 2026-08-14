import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import AddonsPage from './pages/AddonsPage'
import ExperienceDetailPage from './pages/ExperienceDetailPage'
import ExperiencesPage from './pages/ExperiencesPage'
import PackDetailPage from './pages/PackDetailPage'
import PacksPage from './pages/PacksPage'
import LoginPage from './pages/LoginPage'
import RequireUser from './components/RequireUser'

export default function App() {
  return <BrowserRouter><Routes>
    <Route path="/experiences" element={<ExperiencesPage />} />
    <Route path="/experiences/:packId" element={<ExperienceDetailPage />} />
    <Route path="/addons" element={<AddonsPage />} />
    <Route path="/packs" element={<RequireUser><PacksPage /></RequireUser>} />
    <Route path="/packs/:packId" element={<RequireUser><PackDetailPage /></RequireUser>} />
    <Route path="/login" element={<LoginPage />} />
    <Route path="*" element={<Navigate to="/experiences" replace />} />
  </Routes></BrowserRouter>
}
