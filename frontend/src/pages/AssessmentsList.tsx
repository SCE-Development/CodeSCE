import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { listAssessments } from '../api/client'

type Assessment = {
  id: number
  title: string
  description?: string
  durationMinutes: number
  createdAt: string
  updatedAt: string
}

const dateFormatter = new Intl.DateTimeFormat(undefined, {
  year: 'numeric',
  month: 'short',
  day: 'numeric',
})

function AssessmentsList() {
  const { data, isLoading, isError, error } = useQuery<Assessment[]>({
    queryKey: ['assessments'],
    queryFn: () => listAssessments<Assessment[]>(),
  })

  return (
    <div className="max-w-4xl">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-semibold text-gray-900">Assessments</h1>
        <Link
          to="/assessments/new"
          className="px-4 py-2 bg-gray-900 text-white rounded-md text-sm font-medium hover:bg-gray-800"
        >
          Create Assessment
        </Link>
      </div>

      {isLoading && <p className="text-gray-600">Loading assessments…</p>}

      {isError && (
        <p className="text-red-700">
          {error instanceof Error ? error.message : 'Failed to load assessments'}
        </p>
      )}

      {!isLoading && !isError && data && data.length === 0 && (
        <div className="rounded-md border border-dashed border-gray-300 p-8 text-center">
          <p className="text-gray-700 font-medium">No assessments yet</p>
          <p className="text-gray-500 text-sm mt-1">
            Create your first assessment to get started.
          </p>
        </div>
      )}

      {!isLoading && !isError && data && data.length > 0 && (
        <ul className="divide-y divide-gray-200 rounded-md border border-gray-200 bg-white">
          {data.map((assessment) => (
            <li key={assessment.id} className="px-4 py-3 flex items-center justify-between">
              <div>
                <p className="font-medium text-gray-900">{assessment.title}</p>
                <p className="text-sm text-gray-500">
                  {assessment.durationMinutes} min · Created{' '}
                  {dateFormatter.format(new Date(assessment.createdAt))}
                </p>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default AssessmentsList
