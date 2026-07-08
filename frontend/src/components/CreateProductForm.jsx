import { useState } from 'react'
import { useCreateProduct } from '../api/products'

export default function CreateProductForm({ onClose }) {
  const [form, setForm] = useState({
    sku: '', name: '', description: '', quantity: 0, price: '', low_stock_threshold: 10,
  })
  const createProduct = useCreateProduct()
  const [error, setError] = useState(null)

  async function handleSubmit(e) {
    e.preventDefault()
    setError(null)
    try {
      await createProduct.mutateAsync({
        ...form,
        quantity: parseInt(form.quantity, 10),
        price: parseFloat(form.price),
        low_stock_threshold: parseInt(form.low_stock_threshold, 10),
      })
      onClose()
    } catch (err) {
      setError(err.response?.data?.error ?? err.message)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <form onSubmit={handleSubmit} className="bg-white rounded-xl p-6 w-full max-w-md space-y-3">
        <h2 className="text-lg font-semibold">New product</h2>
        {error && <p className="text-xs text-red-600">{error}</p>}

        <input required placeholder="SKU" value={form.sku}
          onChange={e => setForm({ ...form, sku: e.target.value })}
          className="w-full border rounded px-3 py-2 text-sm" />
        <input required placeholder="Name" value={form.name}
          onChange={e => setForm({ ...form, name: e.target.value })}
          className="w-full border rounded px-3 py-2 text-sm" />
        <textarea placeholder="Description" value={form.description}
          onChange={e => setForm({ ...form, description: e.target.value })}
          className="w-full border rounded px-3 py-2 text-sm" rows={2} />
        <div className="grid grid-cols-3 gap-2">
          <input required type="number" placeholder="Qty" value={form.quantity}
            onChange={e => setForm({ ...form, quantity: e.target.value })}
            className="border rounded px-3 py-2 text-sm" />
          <input required type="number" step="0.01" placeholder="Price" value={form.price}
            onChange={e => setForm({ ...form, price: e.target.value })}
            className="border rounded px-3 py-2 text-sm" />
          <input required type="number" placeholder="Low stock at" value={form.low_stock_threshold}
            onChange={e => setForm({ ...form, low_stock_threshold: e.target.value })}
            className="border rounded px-3 py-2 text-sm" />
        </div>

        <div className="flex gap-2 pt-2">
          <button type="button" onClick={onClose} className="flex-1 border rounded py-2 text-sm">
            Cancel
          </button>
          <button type="submit" disabled={createProduct.isPending}
            className="flex-1 bg-blue-600 text-white rounded py-2 text-sm disabled:opacity-50">
            {createProduct.isPending ? 'Creating…' : 'Create'}
          </button>
        </div>
      </form>
    </div>
  )
}