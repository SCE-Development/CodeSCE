import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import './index.css'
import Home from './Home.tsx'
import Assessment from './Assessment.tsx'
import RecruiterLayout from './components/RecruiterLayout.tsx'
import AssessmentsList from './pages/AssessmentsList.tsx'
import AssessmentEditor from './pages/AssessmentEditor.tsx'

const queryClient = new QueryClient()

const router = createBrowserRouter([
  { path: '/', element: <Home /> },
  { path: '/assessment', element: <Assessment /> },
  {
    element: <RecruiterLayout />,
    children: [
      { path: '/recruiter', element: <div>Dashboard</div> },
      { path: '/recruiter/candidates', element: <div>Candidates</div> },
      { path: '/assessments', element: <AssessmentsList /> },
      { path: '/assessments/new', element: <AssessmentEditor /> },
      { path: '/assessments/:id', element: <AssessmentEditor /> },
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
