export default function ProductCard({ product, onStockOp }) {
  const isLowStock = product.quantity <= product.low_stock_threshold

  return (
    <div className={`bg-white border rounded-xl p-4 space-y-2 ${isLowStock ? 'border-red-300' : ''}`}>
      <div className="flex justify-between items-start">
        <div>
          <p className="text-xs text-gray-400 font-mono">{product.sku}</p>
          <h3 className="font-semibold text-sm">{product.name}</h3>
        </div>
        <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${
          isLowStock ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600'
        }`}>
          {isLowStock ? 'Low stock' : 'In stock'}
        </span>
      </div>

      {product.description && (
        <p className="text-xs text-gray-500 line-clamp-2">{product.description}</p>
      )}

      <div className="flex justify-between items-center pt-1">
        <div>
          <p className="text-lg font-semibold">{product.quantity}</p>
          <p className="text-xs text-gray-400">units · ${product.price.toFixed(2)} each</p>
        </div>
        <button
          onClick={onStockOp}
          className="text-sm border rounded-lg px-3 py-1.5 hover:bg-gray-50"
        >
          Stock op
        </button>
      </div>

      <div className="text-xs text-gray-400 pt-1 border-t">
        Low stock alert at {product.low_stock_threshold} units
      </div>
    </div>
  )
}