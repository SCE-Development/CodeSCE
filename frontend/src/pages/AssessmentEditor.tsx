import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '../api/client'
import QuestionForm, { type QuestionPayload } from '../components/QuestionForm'
import TestCaseManager, { type TestCase } from '../components/TestCaseManager'

type Question = {
  id: number
  type: string
  prompt: string
  score: number
  language?: string
  testCases?: TestCase[]
}

type AssessmentDetail = {
  id: number
  title: string
  description?: string
  durationMinutes: number
  createdAt: string
  updatedAt: string
  questions: Question[]
}

type AssessmentPayload = {
  title: string
  description?: string
  duration_minutes: number
}

function AssessmentEditor() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const params = useParams<{ id: string }>()
  const isEdit = Boolean(params.id)
  const assessmentId = params.id

  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [durationMinutes, setDurationMinutes] = useState('30')
  const [formError, setFormError] = useState('')
  const [showQuestionForm, setShowQuestionForm] = useState(false)

  const detailQuery = useQuery<AssessmentDetail>({
    queryKey: ['assessment', assessmentId],
    queryFn: () => apiClient.get<AssessmentDetail>(`/assessments/${assessmentId}`),
    enabled: isEdit,
  })

  useEffect(() => {
    const data = detailQuery.data
    if (!data) return
    setTitle(data.title)
    setDescription(data.description ?? '')
    setDurationMinutes(String(data.durationMinutes))
  }, [detailQuery.data])

  const createMutation = useMutation({
    mutationFn: (payload: AssessmentPayload) =>
      apiClient.post<AssessmentDetail>('/assessments', payload),
    onSuccess: (created) => {
      queryClient.invalidateQueries({ queryKey: ['assessments'] })
      navigate(`/assessments/${created.id}`)
    },
  })

  const updateMutation = useMutation({
    mutationFn: (payload: AssessmentPayload) =>
      apiClient.put<AssessmentDetail>(`/assessments/${assessmentId}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assessments'] })
      queryClient.invalidateQueries({ queryKey: ['assessment', assessmentId] })
    },
  })

  const addQuestionMutation = useMutation({
    mutationFn: (payload: QuestionPayload) =>
      apiClient.post<Question>(`/assessments/${assessmentId}/questions`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assessment', assessmentId] })
      setShowQuestionForm(false)
    },
  })

  const questionSubmitError =
    addQuestionMutation.error instanceof Error ? addQuestionMutation.error.message : null

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    setFormError('')

    const trimmedTitle = title.trim()
    if (!trimmedTitle) {
      setFormError('Title is required')
      return
    }

    const duration = Number(durationMinutes)
    if (!Number.isFinite(duration) || duration <= 0) {
      setFormError('Duration must be greater than 0')
      return
    }

    const payload: AssessmentPayload = {
      title: trimmedTitle,
      description: description.trim() || undefined,
      duration_minutes: duration,
    }

    if (isEdit) {
      updateMutation.mutate(payload)
    } else {
      createMutation.mutate(payload)
    }
  }

  const activeMutation = isEdit ? updateMutation : createMutation
  const isSubmitting = activeMutation.isPending
  const submitError =
    activeMutation.error instanceof Error ? activeMutation.error.message : null

  if (isEdit && detailQuery.isLoading) {
    return <p className="text-gray-600">Loading assessment…</p>
  }

  if (isEdit && detailQuery.isError) {
    return (
      <p className="text-red-700">
        {detailQuery.error instanceof Error
          ? detailQuery.error.message
          : 'Failed to load assessment'}
      </p>
    )
  }

  return (
    <div className="max-w-3xl">
      <h1 className="text-2xl font-semibold text-gray-900 mb-6">
        {isEdit ? 'Edit Assessment' : 'Create Assessment'}
      </h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
            Title
          </label>
          <input
            id="title"
            type="text"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            className="w-full border border-gray-300 rounded-md px-3 py-2"
          />
        </div>

        <div>
          <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
            Description
          </label>
          <textarea
            id="description"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            rows={4}
            className="w-full border border-gray-300 rounded-md px-3 py-2"
          />
        </div>

        <div>
          <label
            htmlFor="duration"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Duration (minutes)
          </label>
          <input
            id="duration"
            type="number"
            min={1}
            value={durationMinutes}
            onChange={(event) => setDurationMinutes(event.target.value)}
            className="w-full border border-gray-300 rounded-md px-3 py-2"
          />
        </div>

        {formError && <p className="text-red-700 text-sm">{formError}</p>}
        {submitError && <p className="text-red-700 text-sm">{submitError}</p>}

        <div className="flex gap-2">
          <button
            type="submit"
            disabled={isSubmitting}
            className="px-4 py-2 bg-gray-900 text-white rounded-md text-sm font-medium hover:bg-gray-800 disabled:opacity-60"
          >
            {isSubmitting
              ? 'Saving…'
              : isEdit
                ? 'Save changes'
                : 'Create assessment'}
          </button>
          <button
            type="button"
            onClick={() => navigate('/assessments')}
            className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium hover:bg-gray-100"
          >
            Cancel
          </button>
        </div>
      </form>

      {isEdit && detailQuery.data && (
        <section className="mt-10">
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-lg font-semibold text-gray-900">Questions</h2>
            {!showQuestionForm && (
              <button
                type="button"
                onClick={() => {
                  addQuestionMutation.reset()
                  setShowQuestionForm(true)
                }}
                className="px-3 py-1.5 bg-gray-900 text-white rounded-md text-sm font-medium hover:bg-gray-800"
              >
                Add Question
              </button>
            )}
          </div>

          {showQuestionForm && (
            <div className="mb-4">
              <QuestionForm
                onSubmit={(payload) => addQuestionMutation.mutate(payload)}
                onCancel={() => {
                  addQuestionMutation.reset()
                  setShowQuestionForm(false)
                }}
                isSubmitting={addQuestionMutation.isPending}
                submitError={questionSubmitError}
              />
            </div>
          )}

          {detailQuery.data.questions.length === 0 ? (
            <p className="text-gray-500 text-sm">No questions added yet.</p>
          ) : (
            <ul className="divide-y divide-gray-200 rounded-md border border-gray-200 bg-white">
              {detailQuery.data.questions.map((question, index) => (
                <li key={question.id} className="px-4 py-3">
                  <div className="flex items-center justify-between">
                    <p className="font-medium text-gray-900">
                      {index + 1}. {question.prompt}
                    </p>
                    <span className="text-xs text-gray-500 uppercase tracking-wide">
                      {question.type}
                    </span>
                  </div>
                  <p className="text-sm text-gray-500 mt-1">
                    {question.score} pts
                    {question.language ? ` · ${question.language}` : ''}
                  </p>
                  {question.type === 'code' && assessmentId && (
                    <TestCaseManager
                      assessmentId={assessmentId}
                      questionId={question.id}
                      testCases={question.testCases ?? []}
                    />
                  )}
                </li>
              ))}
            </ul>
          )}
        </section>
      )}
    </div>
  )
}

export default AssessmentEditor
