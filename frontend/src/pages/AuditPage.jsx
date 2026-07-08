import { useQuery } from '@tanstack/react-query'
import api from '../api/client'

export default function AuditPage() {
  const { data: logs, isLoading } = useQuery({
    queryKey: ['audit'],
    queryFn: () => api.get('/audit').then(r => r.data),
  })

  if (isLoading) return <div className="p-8">Loading audit log…</div>

  return (
    <div className="p-6">
      <h1 className="text-base font-semibold mb-4">Audit log</h1>
      <div className="space-y-2">
        {logs?.map(log => (
          <div key={log.id} className="border rounded-lg p-3 text-sm">
            <div className="flex justify-between text-xs text-gray-500 mb-1">
              <span className="capitalize font-medium text-gray-700">
                {log.action} · {log.entity_type}
              </span>
              <span>{new Date(log.created_at).toLocaleString()}</span>
            </div>
            <p className="text-xs text-gray-400">by {log.performed_by}</p>
            {log.action === 'update' && (
              <div className="mt-2 text-xs grid grid-cols-2 gap-2">
                <pre className="bg-red-50 rounded p-2 overflow-x-auto">
                  {JSON.stringify(log.old_values, null, 2)}
                </pre>
                <pre className="bg-green-50 rounded p-2 overflow-x-auto">
                  {JSON.stringify(log.new_values, null, 2)}
                </pre>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}