import { useState } from 'react'

export type QuestionType = 'mcq_single' | 'mcq_multi' | 'short_answer' | 'code'

export type QuestionPayload = {
  type: QuestionType
  prompt: string
  score: number
  options?: string[]
  correct_answer?: unknown
  language?: string
}

type QuestionFormProps = {
  onSubmit: (payload: QuestionPayload) => void
  onCancel: () => void
  isSubmitting?: boolean
  submitError?: string | null
}

const LANGUAGES = ['python', 'javascript', 'go', 'java', 'cpp']

function QuestionForm({ onSubmit, onCancel, isSubmitting, submitError }: QuestionFormProps) {
  const [type, setType] = useState<QuestionType>('mcq_single')
  const [prompt, setPrompt] = useState('')
  const [score, setScore] = useState('1')
  const [options, setOptions] = useState<string[]>(['', ''])
  const [singleCorrect, setSingleCorrect] = useState<number | null>(null)
  const [multiCorrect, setMultiCorrect] = useState<number[]>([])
  const [shortAnswer, setShortAnswer] = useState('')
  const [language, setLanguage] = useState<string>(LANGUAGES[0])
  const [formError, setFormError] = useState('')

  const isMcq = type === 'mcq_single' || type === 'mcq_multi'

  const updateOption = (index: number, value: string) => {
    setOptions((prev) => prev.map((option, i) => (i === index ? value : option)))
  }

  const addOption = () => setOptions((prev) => [...prev, ''])

  const removeOption = (index: number) => {
    setOptions((prev) => prev.filter((_, i) => i !== index))
    if (type === 'mcq_single') {
      setSingleCorrect((prev) => {
        if (prev === null) return prev
        if (prev === index) return null
        return prev > index ? prev - 1 : prev
      })
    } else {
      setMultiCorrect((prev) =>
        prev
          .filter((i) => i !== index)
          .map((i) => (i > index ? i - 1 : i)),
      )
    }
  }

  const toggleMultiCorrect = (index: number) => {
    setMultiCorrect((prev) =>
      prev.includes(index) ? prev.filter((i) => i !== index) : [...prev, index],
    )
  }

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    setFormError('')

    const trimmedPrompt = prompt.trim()
    if (!trimmedPrompt) {
      setFormError('Prompt is required')
      return
    }

    const numericScore = Number(score)
    if (!Number.isFinite(numericScore) || numericScore <= 0) {
      setFormError('Score must be greater than 0')
      return
    }

    const payload: QuestionPayload = {
      type,
      prompt: trimmedPrompt,
      score: numericScore,
    }

    if (isMcq) {
      const trimmedOptions = options.map((option) => option.trim())
      if (trimmedOptions.some((option) => !option) || trimmedOptions.length < 2) {
        setFormError('Provide at least 2 non-empty options')
        return
      }
      payload.options = trimmedOptions

      if (type === 'mcq_single') {
        if (singleCorrect === null) {
          setFormError('Select the correct answer')
          return
        }
        payload.correct_answer = singleCorrect
      } else {
        if (multiCorrect.length === 0) {
          setFormError('Select at least one correct answer')
          return
        }
        payload.correct_answer = [...multiCorrect].sort((a, b) => a - b)
      }
    } else if (type === 'short_answer') {
      const trimmedAnswer = shortAnswer.trim()
      if (trimmedAnswer) {
        payload.correct_answer = trimmedAnswer
      }
    } else if (type === 'code') {
      if (!language) {
        setFormError('Language is required for code questions')
        return
      }
      payload.language = language
    }

    onSubmit(payload)
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="space-y-4 rounded-md border border-gray-200 bg-white p-4"
    >
      <div>
        <label htmlFor="question-type" className="block text-sm font-medium text-gray-700 mb-1">
          Type
        </label>
        <select
          id="question-type"
          value={type}
          onChange={(event) => setType(event.target.value as QuestionType)}
          className="w-full border border-gray-300 rounded-md px-3 py-2"
        >
          <option value="mcq_single">Multiple choice (single)</option>
          <option value="mcq_multi">Multiple choice (multi)</option>
          <option value="short_answer">Short answer</option>
          <option value="code">Code</option>
        </select>
      </div>

      <div>
        <label htmlFor="question-prompt" className="block text-sm font-medium text-gray-700 mb-1">
          Prompt
        </label>
        <textarea
          id="question-prompt"
          value={prompt}
          onChange={(event) => setPrompt(event.target.value)}
          rows={3}
          className="w-full border border-gray-300 rounded-md px-3 py-2"
        />
      </div>

      <div>
        <label htmlFor="question-score" className="block text-sm font-medium text-gray-700 mb-1">
          Score
        </label>
        <input
          id="question-score"
          type="number"
          min={1}
          value={score}
          onChange={(event) => setScore(event.target.value)}
          className="w-full border border-gray-300 rounded-md px-3 py-2"
        />
      </div>

      {isMcq && (
        <div>
          <div className="flex items-center justify-between mb-2">
            <label className="block text-sm font-medium text-gray-700">Options</label>
            <button
              type="button"
              onClick={addOption}
              className="text-sm font-medium text-gray-700 hover:text-gray-900"
            >
              + Add option
            </button>
          </div>
          <ul className="space-y-2">
            {options.map((option, index) => (
              <li key={index} className="flex items-center gap-2">
                {type === 'mcq_single' ? (
                  <input
                    type="radio"
                    name="mcq-correct"
                    checked={singleCorrect === index}
                    onChange={() => setSingleCorrect(index)}
                    aria-label={`Mark option ${index + 1} correct`}
                  />
                ) : (
                  <input
                    type="checkbox"
                    checked={multiCorrect.includes(index)}
                    onChange={() => toggleMultiCorrect(index)}
                    aria-label={`Mark option ${index + 1} correct`}
                  />
                )}
                <input
                  type="text"
                  value={option}
                  onChange={(event) => updateOption(index, event.target.value)}
                  placeholder={`Option ${index + 1}`}
                  className="flex-1 border border-gray-300 rounded-md px-3 py-2"
                />
                <button
                  type="button"
                  onClick={() => removeOption(index)}
                  disabled={options.length <= 2}
                  className="text-sm text-red-700 hover:text-red-900 disabled:opacity-40"
                >
                  Remove
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}

      {type === 'short_answer' && (
        <div>
          <label
            htmlFor="short-answer"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Expected answer (optional)
          </label>
          <input
            id="short-answer"
            type="text"
            value={shortAnswer}
            onChange={(event) => setShortAnswer(event.target.value)}
            className="w-full border border-gray-300 rounded-md px-3 py-2"
          />
        </div>
      )}

      {type === 'code' && (
        <div>
          <label htmlFor="code-language" className="block text-sm font-medium text-gray-700 mb-1">
            Language
          </label>
          <select
            id="code-language"
            value={language}
            onChange={(event) => setLanguage(event.target.value)}
            className="w-full border border-gray-300 rounded-md px-3 py-2"
          >
            {LANGUAGES.map((lang) => (
              <option key={lang} value={lang}>
                {lang}
              </option>
            ))}
          </select>
        </div>
      )}

      {formError && <p className="text-red-700 text-sm">{formError}</p>}
      {submitError && <p className="text-red-700 text-sm">{submitError}</p>}

      <div className="flex gap-2">
        <button
          type="submit"
          disabled={isSubmitting}
          className="px-4 py-2 bg-gray-900 text-white rounded-md text-sm font-medium hover:bg-gray-800 disabled:opacity-60"
        >
          {isSubmitting ? 'Adding…' : 'Add question'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium hover:bg-gray-100"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}

export default QuestionForm
