import './App.css'

type AssessmentProps = {
    solution: string
    onSolutionChange: (solution: string) => void
    onBack: () => void
}

function Assessment({ solution, onSolutionChange, onBack }: AssessmentProps) {
    return (
        <main className="app-container assessment-container">
            <section className="assessment-page">
                <h1>CodeSCE Assessment</h1>

                <div className="assessment-layout">
                    <section className="question-section">
                        <p>
                            Placeholder question text
                        </p>
                    </section>

                    <section className="right-section">
                        <section className="answer-section">
                            <textarea
                                value={solution}
                                onChange={(event) => onSolutionChange(event.target.value)}
                            />
                        </section>

                        <section className="results-section">
                            <h2>Compile Messages</h2>
                            <p>Compile Messages Output Placeholder</p>
                            <h2>Test Cases</h2>
                            <p>Test Cases Output Placeholder</p>
                        </section>
                    </section>
                </div>

                <button type="button" onClick={onBack}>
                    Back
                </button>
            </section>
        </main>
    )
}

export default Assessment
