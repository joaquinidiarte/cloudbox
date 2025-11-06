import api from './axios'

export const authAPI = {
  register: async (data) => {
    const response = await api.post('/auth/register', data)
    return response.data
  },

  login: async (data) => {
    const response = await api.post('/auth/login', data)
    return response.data
  },

  verify: async () => {
    const response = await api.post('/auth/verify')
    return response.data
  },
}
