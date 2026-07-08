import { useState } from 'react'
import { useStockOperation } from '../api/products'

export default function StockModal({ product, onClose }) {
  const [operation, setOperation] = useState('in')
  const [quantity, setQuantity] = useState('')
  const [notes, setNotes] = useState('')

  const stockOp = useStockOperation()

  async function handleSubmit(e) {
    e.preventDefault()
    try {
      await stockOp.mutateAsync({
        productId: product.id,
        operation,
        quantity: parseInt(quantity, 10),
        notes,
      })
      onClose()
    } catch (err) {
      // Error message comes from the Go backend
      alert(err.response?.data?.error ?? err.message)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
        <h2 className="text-lg font-semibold mb-1">{product.name}</h2>
        <p className="text-sm text-gray-500 mb-4">Current stock: {product.quantity}</p>

        <form onSubmit={handleSubmit} className="space-y-3">
          <div className="flex gap-2">
            {['in', 'out', 'adjustment'].map(op => (
              <button
                key={op}
                type="button"
                onClick={() => setOperation(op)}
                className={`px-3 py-1 rounded text-sm border capitalize
                  ${operation === op ? 'bg-blue-600 text-white border-blue-600' : 'border-gray-300'}`}
              >
                {op}
              </button>
            ))}
          </div>

          <input
            type="number"
            min="1"
            placeholder="Quantity"
            value={quantity}
            onChange={e => setQuantity(e.target.value)}
            required
            className="w-full border rounded px-3 py-2 text-sm"
          />

          <textarea
            placeholder="Notes (optional)"
            value={notes}
            onChange={e => setNotes(e.target.value)}
            className="w-full border rounded px-3 py-2 text-sm"
            rows={2}
          />

          <div className="flex gap-2 pt-2">
            <button type="button" onClick={onClose} className="flex-1 border rounded py-2 text-sm">
              Cancel
            </button>
            <button
              type="submit"
              disabled={stockOp.isPending}
              className="flex-1 bg-blue-600 text-white rounded py-2 text-sm disabled:opacity-50"
            >
              {stockOp.isPending ? 'Saving…' : 'Confirm'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}