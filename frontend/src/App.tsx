import { useState } from 'react'
import Assessment from './Assessment'
import Home from './Home'

function App() {
  const [page, setPage] = useState<'home' | 'assessment'>('home')
  const [name, setName] = useState('')
  const [solution, setSolution] = useState('')
  const [error, setError] = useState('')

  const handleStart = () => {
    const trimmedName = name.trim()

    if (!trimmedName) {
      setError('Enter name')
      return
    }

    setError('')
    console.log(trimmedName)
    setPage('assessment')
  }

  if (page === 'assessment') {
    return (
      <Assessment
        solution={solution}
        onSolutionChange={setSolution}
        onBack={() => setPage('home')}
      />
    )
  }

  return (
    <Home
      name={name}
      error={error}
      onNameChange={setName}
      onStart={handleStart}
    />
  )
}

export default App
