import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { addTestCases } from '../api/client'

export type TestCase = {
  id: number
  questionId: number
  input: string
  expectedOutput: string
  isHidden: boolean
  score: number
}

type TestCasePayload = {
  input: string
  expectedOutput: string
  isHidden: boolean
  score: number
}

type TestCaseManagerProps = {
  assessmentId: string
  questionId: number
  testCases: TestCase[]
}

function TestCaseManager({ assessmentId, questionId, testCases }: TestCaseManagerProps) {
  const queryClient = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [input, setInput] = useState('')
  const [expectedOutput, setExpectedOutput] = useState('')
  const [isHidden, setIsHidden] = useState(false)
  const [score, setScore] = useState('1')
  const [formError, setFormError] = useState('')

  const addMutation = useMutation({
    mutationFn: (payload: TestCasePayload) =>
      addTestCases<{ testCases: TestCase[] }>(assessmentId, questionId, {
        testCases: [payload],
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assessment', assessmentId] })
      setInput('')
      setExpectedOutput('')
      setIsHidden(false)
      setScore('1')
      setShowForm(false)
    },
  })

  const submitError =
    addMutation.error instanceof Error ? addMutation.error.message : null

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    setFormError('')

    const numericScore = Number(score)
    if (!Number.isFinite(numericScore) || numericScore <= 0) {
      setFormError('Score must be greater than 0')
      return
    }

    addMutation.mutate({
      input,
      expectedOutput,
      isHidden,
      score: numericScore,
    })
  }

  return (
    <div className="mt-3 border-t border-gray-200 pt-3">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-sm font-semibold text-gray-700">Test cases</h3>
        {!showForm && (
          <button
            type="button"
            onClick={() => {
              addMutation.reset()
              setShowForm(true)
            }}
            className="text-sm font-medium text-gray-700 hover:text-gray-900"
          >
            + Add test case
          </button>
        )}
      </div>

      {testCases.length === 0 ? (
        <p className="text-xs text-gray-500">No test cases yet.</p>
      ) : (
        <ul className="space-y-2">
          {testCases.map((testCase) => (
            <li
              key={testCase.id}
              className="rounded-md border border-gray-200 p-2 text-xs"
            >
              <div className="flex items-center justify-between mb-1">
                <span className="font-medium text-gray-700">
                  {testCase.score} pts
                </span>
                {testCase.isHidden && (
                  <span className="px-1.5 py-0.5 rounded bg-gray-200 text-gray-700 uppercase tracking-wide">
                    Hidden
                  </span>
                )}
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <p className="text-gray-500 mb-0.5">Input</p>
                  <pre className="font-mono whitespace-pre-wrap break-words bg-gray-50 rounded p-1.5">
                    {testCase.input || <span className="text-gray-400">(empty)</span>}
                  </pre>
                </div>
                <div>
                  <p className="text-gray-500 mb-0.5">Expected output</p>
                  <pre className="font-mono whitespace-pre-wrap break-words bg-gray-50 rounded p-1.5">
                    {testCase.expectedOutput || (
                      <span className="text-gray-400">(empty)</span>
                    )}
                  </pre>
                </div>
              </div>
            </li>
          ))}
        </ul>
      )}

      {showForm && (
        <form
          onSubmit={handleSubmit}
          className="mt-3 space-y-3 rounded-md border border-gray-200 bg-gray-50 p-3"
        >
          <div>
            <label className="block text-xs font-medium text-gray-700 mb-1">
              Input
            </label>
            <textarea
              value={input}
              onChange={(event) => setInput(event.target.value)}
              rows={2}
              className="w-full font-mono text-sm border border-gray-300 rounded-md px-2 py-1.5"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700 mb-1">
              Expected output
            </label>
            <textarea
              value={expectedOutput}
              onChange={(event) => setExpectedOutput(event.target.value)}
              rows={2}
              className="w-full font-mono text-sm border border-gray-300 rounded-md px-2 py-1.5"
            />
          </div>
          <div className="flex items-center gap-4">
            <label className="flex items-center gap-2 text-sm text-gray-700">
              <input
                type="checkbox"
                checked={isHidden}
                onChange={(event) => setIsHidden(event.target.checked)}
              />
              Hidden
            </label>
            <div className="flex items-center gap-2">
              <label className="text-sm text-gray-700">Score</label>
              <input
                type="number"
                min={1}
                value={score}
                onChange={(event) => setScore(event.target.value)}
                className="w-20 border border-gray-300 rounded-md px-2 py-1"
              />
            </div>
          </div>

          {formError && <p className="text-red-700 text-xs">{formError}</p>}
          {submitError && <p className="text-red-700 text-xs">{submitError}</p>}

          <div className="flex gap-2">
            <button
              type="submit"
              disabled={addMutation.isPending}
              className="px-3 py-1.5 bg-gray-900 text-white rounded-md text-sm font-medium hover:bg-gray-800 disabled:opacity-60"
            >
              {addMutation.isPending ? 'Adding…' : 'Add test case'}
            </button>
            <button
              type="button"
              onClick={() => {
                addMutation.reset()
                setShowForm(false)
              }}
              className="px-3 py-1.5 border border-gray-300 rounded-md text-sm font-medium hover:bg-gray-100"
            >
              Cancel
            </button>
          </div>
        </form>
      )}
    </div>
  )
}

export default TestCaseManager
