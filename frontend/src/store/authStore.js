import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export const useAuthStore = create(
  persist(
    (set) => ({
      token: null,
      user: null,

      setAuth: (token, user) => set({ token, user }),

      logout: () => set({ token: null, user: null }),

      updateUser: (userData) => set((state) => ({
        user: { ...state.user, ...userData }
      })),

      refreshUser: async () => {
        try {
          const { usersAPI } = await import('../api/users')
          const response = await usersAPI.getMe()
          if (response.success) {
            set({ user: response.data })
          }
        } catch (error) {
          console.error('Failed to refresh user:', error)
        }
      },
    }),
    {
      name: 'cloudbox-auth',
    }
  )
)
