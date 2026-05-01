import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function Assessment() {
  const navigate = useNavigate()
  const [solution, setSolution] = useState('')

  return (
    <main className="min-h-screen">
      <section className="w-full">
        <h1 className="text-2xl font-semibold p-3">CodeSCE Assessment</h1>

        <div className="grid grid-cols-1 md:grid-cols-[1fr_2fr] min-h-[calc(100vh-6rem)]">
          <section className="p-3 border-b md:border-b-0 md:border-r border-gray-300">
            <p>Placeholder question text</p>
          </section>

          <section className="flex flex-col">
            <section className="p-3">
              <textarea
                value={solution}
                onChange={(event) => setSolution(event.target.value)}
                className="w-full min-h-[220px] resize-y border border-gray-300 rounded-md p-2"
              />
            </section>

            <section className="p-3 border-t border-gray-300">
              <h2 className="text-lg font-semibold">Compile Messages</h2>
              <p>Compile Messages Output Placeholder</p>
              <h2 className="text-lg font-semibold mt-2">Test Cases</h2>
              <p>Test Cases Output Placeholder</p>
            </section>
          </section>
        </div>

        <div className="p-3">
          <button
            type="button"
            onClick={() => navigate('/')}
            className="px-3 py-1.5 bg-gray-900 text-white rounded-md hover:bg-gray-800"
          >
            Back
          </button>
        </div>
      </section>
    </main>
  )
}

export default Assessment
