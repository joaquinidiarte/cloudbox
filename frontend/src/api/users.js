import api from './axios'

export const usersAPI = {
  getMe: async () => {
    const response = await api.get('/users/me')
    return response.data
  },

  updateMe: async (data) => {
    const response = await api.put('/users/me', data)
    return response.data
  },
}
