import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function Home() {
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [error, setError] = useState('')

  const handleStart = () => {
    const trimmedName = name.trim()

    if (!trimmedName) {
      setError('Enter name')
      return
    }

    setError('')
    navigate('/assessment', { state: { name: trimmedName } })
  }

  return (
    <main className="min-h-screen flex items-center justify-center p-4">
      <section className="w-full max-w-sm text-center">
        <h1 className="text-2xl font-semibold mb-4">CodeSCE</h1>

        <div className="flex gap-2 justify-center">
          <input
            type="text"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="Name"
            aria-label="Name"
            className="flex-1 border border-gray-300 rounded-md px-2 py-1.5"
          />
          <button
            type="button"
            onClick={handleStart}
            className="px-3 py-1.5 bg-gray-900 text-white rounded-md hover:bg-gray-800"
          >
            Start
          </button>
        </div>

        {error && <p className="mt-2 text-red-700">{error}</p>}
      </section>
    </main>
  )
}

export default Home
