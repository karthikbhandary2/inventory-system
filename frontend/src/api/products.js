import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from './client'

// Fetch + cache products list with optional search/filter
export function useProducts(params = {}) {
  return useQuery({
    queryKey: ['products', params],
    queryFn: () => api.get('/products', { params }).then(r => r.data),
    staleTime: 30_000, // cache for 30s
  })
}

export function useProduct(id) {
  return useQuery({
    queryKey: ['products', id],
    queryFn: () => api.get(`/products/${id}`).then(r => r.data),
    enabled: !!id,
  })
}

export function useCreateProduct() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: data => api.post('/products', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['products'] }),
  })
}

export function useStockOperation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ productId, ...data }) =>
      api.post(`/products/${productId}/stock`, data).then(r => r.data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['products'] })
      qc.invalidateQueries({ queryKey: ['report'] })
    },
  })
}

export function useReport() {
  return useQuery({
    queryKey: ['report'],
    queryFn: () => api.get('/reports/inventory').then(r => r.data),
  })
}