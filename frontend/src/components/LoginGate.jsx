import { useState } from 'react'

export default function LoginGate({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('token'))
  const [input, setInput] = useState('')

  if (token) return children

  function handleSubmit(e) {
    e.preventDefault()
    localStorage.setItem('token', input.trim())
    setToken(input.trim())
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <form onSubmit={handleSubmit} className="bg-white border rounded-xl p-6 w-full max-w-sm space-y-3">
        <h1 className="text-base font-semibold">Paste your test JWT</h1>
        <p className="text-xs text-gray-500">
          Temporary — replace with a real login form once you build /auth/login on the backend.
        </p>
        <textarea
          value={input}
          onChange={e => setInput(e.target.value)}
          rows={4}
          className="w-full border rounded px-3 py-2 text-xs font-mono"
          placeholder="eyJhbGciOi..."
        />
        <button type="submit" className="w-full bg-blue-600 text-white rounded py-2 text-sm">
          Continue
        </button>
      </form>
    </div>
  )
}