import './App.css'

type HomeProps = {
    name: string
    error: string
    onNameChange: (name: string) => void
    onStart: () => void
}

function Home({ name, error, onNameChange, onStart }: HomeProps) {
    return (
        <main className="app-container">
            <section className="start-card">
                <h1>CodeSCE</h1>

                <div className="form-row">
                    <input
                        type="text"
                        value={name}
                        onChange={(event) => onNameChange(event.target.value)}
                        placeholder="Name"
                        aria-label="Name"
                    />
                    <button type="button" onClick={onStart}>
                        Start
                    </button>
                </div>

                {error && <p className="message error">{error}</p>}
            </section>
        </main>
    )
}

export default Home
