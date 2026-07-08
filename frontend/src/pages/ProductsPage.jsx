import { useState } from 'react'
import { useProducts } from '../api/products'
import ProductCard from '../components/ProductCard'
import StockModal from '../components/StockModal'

import CreateProductForm from '../components/CreateProductForm'

export default function ProductsPage() {
    const [search, setSearch] = useState('')
    const [lowStockOnly, setLowStockOnly] = useState(false)
    const [selectedProduct, setSelectedProduct] = useState(null)
    const [showCreate, setShowCreate] = useState(false)

    const { data: products, isLoading, error } = useProducts({
        search: search || undefined,
        low_stock: lowStockOnly || undefined,
    })

    if (isLoading) return <div className="p-8 text-center">Loading products…</div>
    if (error) return <div className="p-8 text-red-500">Error: {error.message}</div>

    return (
        <div className="p-6 space-y-4">
            {/* Filters */}
            <div className="flex justify-between items-center">
                <h1 className="text-base font-semibold">Products</h1>
                <button onClick={() => setShowCreate(true)} className="bg-blue-600 text-white text-sm px-3 py-1.5 rounded">
                    + New product
                </button>
            </div>

            {showCreate && <CreateProductForm onClose={() => setShowCreate(false)} />}
            <div className="flex gap-3">
                <input
                    type="text"
                    placeholder="Search by name or SKU…"
                    value={search}
                    onChange={e => setSearch(e.target.value)}
                    className="flex-1 border rounded px-3 py-2 text-sm"
                />
                <label className="flex items-center gap-2 text-sm">
                    <input
                        type="checkbox"
                        checked={lowStockOnly}
                        onChange={e => setLowStockOnly(e.target.checked)}
                    />
                    Low stock only
                </label>
            </div>

            {/* Product grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {products?.map(product => (
                    <ProductCard
                        key={product.id}
                        product={product}
                        onStockOp={() => setSelectedProduct(product)}
                    />
                ))}
            </div>

            {/* Stock operation modal */}
            {selectedProduct && (
                <StockModal
                    product={selectedProduct}
                    onClose={() => setSelectedProduct(null)}
                />
            )}
        </div>
    )
}