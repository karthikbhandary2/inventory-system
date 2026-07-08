import { useReport } from '../api/products'

export default function ReportPage() {
  const { data: report, isLoading } = useReport()

  if (isLoading) return <div className="p-8">Loading report…</div>

  return (
    <div className="p-6 space-y-6">
      {/* Summary cards */}
      <div className="grid grid-cols-3 gap-4">
        <StatCard label="Total products" value={report.total_products} />
        <StatCard
          label="Total inventory value"
          value={`$${report.total_value.toFixed(2)}`}
        />
        <StatCard
          label="Low stock alerts"
          value={report.total_low_stock}
          alert={report.total_low_stock > 0}
        />
      </div>

      {/* Low stock table */}
      {report.low_stock_items?.length > 0 && (
        <div>
          <h2 className="text-base font-semibold mb-3 text-red-600">Low stock items</h2>
          <table className="w-full text-sm border-collapse">
            <thead>
              <tr className="border-b text-left text-gray-500">
                <th className="pb-2">SKU</th>
                <th className="pb-2">Name</th>
                <th className="pb-2 text-right">Qty</th>
                <th className="pb-2 text-right">Threshold</th>
              </tr>
            </thead>
            <tbody>
              {report.low_stock_items.map(p => (
                <tr key={p.id} className="border-b">
                  <td className="py-2 font-mono text-xs">{p.sku}</td>
                  <td className="py-2">{p.name}</td>
                  <td className="py-2 text-right text-red-600 font-semibold">{p.quantity}</td>
                  <td className="py-2 text-right text-gray-400">{p.low_stock_threshold}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

function StatCard({ label, value, alert }) {
  return (
    <div className={`rounded-xl border p-4 ${alert ? 'border-red-200 bg-red-50' : ''}`}>
      <p className="text-xs text-gray-500 mb-1">{label}</p>
      <p className={`text-2xl font-semibold ${alert ? 'text-red-600' : ''}`}>{value}</p>
    </div>
  )
}