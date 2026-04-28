import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import './index.css'
import Home from './Home.tsx'
import Assessment from './Assessment.tsx'
import RecruiterLayout from './components/RecruiterLayout.tsx'

const queryClient = new QueryClient()

const router = createBrowserRouter([
  { path: '/', element: <Home /> },
  { path: '/assessment', element: <Assessment /> },
  {
    path: '/recruiter',
    element: <RecruiterLayout />,
    children: [
      { index: true, element: <div>Dashboard</div> },
      { path: 'candidates', element: <div>Candidates</div> },
      { path: 'assessments', element: <div>Assessments</div> },
    ],
  },
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
