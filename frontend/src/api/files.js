import api from './axios'

export const filesAPI = {
  list: async (parentId = null) => {
    const params = parentId ? { parent_id: parentId } : {}
    const response = await api.get('/files/', { params })
    return response.data
  },

  upload: async (file, parentId = null) => {
    const formData = new FormData()
    formData.append('file', file)
    if (parentId) {
      formData.append('parent_id', parentId)
    }
    
    const response = await api.post('/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  download: async (fileId) => {
    const response = await api.get(`/files/${fileId}/download`, {
      responseType: 'blob',
    })
    return response
  },

  delete: async (fileId) => {
    const response = await api.delete(`/files/${fileId}`)
    return response.data
  },

  update: async (fileId, data) => {
    const response = await api.put(`/files/${fileId}`, data)
    return response.data
  },

  createFolder: async (name, parentId = null) => {
    const response = await api.post('/files/folders', {
      name,
      parent_id: parentId,
    })
    return response.data
  },

  getFolderContents: async (folderId) => {
    const response = await api.get(`/files/folders/${folderId}`)
    return response.data
  },

  // Version management
  getVersions: async (fileId) => {
    const response = await api.get(`/files/${fileId}/versions`)
    return response.data
  },

  downloadVersion: async (fileId, version) => {
    const response = await api.get(`/files/${fileId}/versions/${version}/download`, {
      responseType: 'blob',
    })
    return response
  },

  restoreVersion: async (fileId, version) => {
    const response = await api.post(`/files/${fileId}/versions/${version}/restore`)
    return response.data
  },

  deleteVersion: async (fileId, version) => {
    const response = await api.delete(`/files/${fileId}/versions/${version}`)
    return response.data
  },
}
