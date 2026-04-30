const baseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:6767/api'

type RequestOptions = Omit<RequestInit, 'body'> & {
  body?: unknown
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { body, headers, ...rest } = options
  const init: RequestInit = {
    ...rest,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
  }

  if (body !== undefined) {
    init.body = typeof body === 'string' ? body : JSON.stringify(body)
  }

  const response = await fetch(`${baseUrl}${path}`, init)

  if (!response.ok) {
    const text = await response.text()
    throw new Error(text || `Request failed: ${response.status}`)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export type CreateAssessmentBody = {
  title: string
  description?: string
  duration_minutes: number
}

export function listAssessments<T = unknown>(): Promise<T> {
  return request<T>('/assessments', { method: 'GET' })
}

export function createAssessment<T = unknown>(body: CreateAssessmentBody): Promise<T> {
  return request<T>('/assessments', { method: 'POST', body })
}

export function getAssessment<T = unknown>(assessmentId: string): Promise<T> {
  return request<T>(`/assessments/${encodeURIComponent(assessmentId)}`, { method: 'GET' })
}

export function updateAssessment<T = unknown>(
  assessmentId: string,
  body: CreateAssessmentBody,
): Promise<T> {
  return request<T>(`/assessments/${encodeURIComponent(assessmentId)}`, {
    method: 'PUT',
    body,
  })
}

export function addQuestion<T = unknown>(
  assessmentId: string,
  body: Record<string, unknown>,
): Promise<T> {
  return request<T>(`/assessments/${encodeURIComponent(assessmentId)}/questions`, {
    method: 'POST',
    body,
  })
}

export function addTestCases<T = unknown>(
  assessmentId: string,
  questionId: number,
  body: { testCases: unknown[] },
): Promise<T> {
  return request<T>(
    `/assessments/${encodeURIComponent(assessmentId)}/questions/${questionId}/test-cases`,
    { method: 'POST', body },
  )
}

export function createInvite<T = unknown>(
  assessmentId: string,
  body: Record<string, unknown>,
): Promise<T> {
  return request<T>(`/assessments/${encodeURIComponent(assessmentId)}/invites`, {
    method: 'POST',
    body,
  })
}

export function listInvites<T = unknown>(): Promise<T> {
  return request<T>('/invites', { method: 'GET' })
}

export function validateInvite<T = unknown>(token: string): Promise<T> {
  return request<T>(`/invites/${encodeURIComponent(token)}`, { method: 'GET' })
}

export function startAttempt<T = unknown>(token: string, body?: Record<string, unknown>): Promise<T> {
  return request<T>(`/invites/${encodeURIComponent(token)}/start`, {
    method: 'POST',
    body: body ?? {},
  })
}

export function listAttempts<T = unknown>(params?: { assessmentId?: string }): Promise<T> {
  const search =
    params?.assessmentId !== undefined && params.assessmentId !== ''
      ? `?assessmentId=${encodeURIComponent(params.assessmentId)}`
      : ''
  return request<T>(`/attempts${search}`, { method: 'GET' })
}

export function getAttempt<T = unknown>(attemptId: string): Promise<T> {
  return request<T>(`/attempts/${encodeURIComponent(attemptId)}`, { method: 'GET' })
}

export function getAttemptQuestions<T = unknown>(attemptId: string): Promise<T> {
  return request<T>(`/attempts/${encodeURIComponent(attemptId)}/questions`, { method: 'GET' })
}

export function saveAnswers<T = unknown>(
  attemptId: string,
  body: Record<string, unknown>,
): Promise<T> {
  return request<T>(`/attempts/${encodeURIComponent(attemptId)}/answers`, {
    method: 'POST',
    body,
  })
}

export function submitAttempt<T = unknown>(attemptId: string): Promise<T> {
  return request<T>(`/attempts/${encodeURIComponent(attemptId)}/submit`, { method: 'POST' })
}

export function createSubmission<T = unknown>(
  attemptId: string,
  questionId: number,
  body: Record<string, unknown>,
): Promise<T> {
  return request<T>(
    `/attempts/${encodeURIComponent(attemptId)}/questions/${questionId}/submissions`,
    { method: 'POST', body },
  )
}

export function getSubmission<T = unknown>(submissionId: string): Promise<T> {
  return request<T>(`/submissions/${encodeURIComponent(submissionId)}`, { method: 'GET' })
}

export function runCode<T = unknown>(
  attemptId: string,
  questionId: number,
  body?: Record<string, unknown>,
): Promise<T> {
  return request<T>(
    `/attempts/${encodeURIComponent(attemptId)}/questions/${questionId}/run`,
    { method: 'POST', body: body ?? {} },
  )
}
