import { useState } from 'react'
import { FolderPlus, AlertCircle } from 'lucide-react'
import { filesAPI } from '../api/files'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'

export default function CreateFolderModal({ open, onClose, onSuccess, parentId }) {
  const [folderName, setFolderName] = useState('')
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState('')

  const handleCreate = async () => {
    if (!folderName.trim()) return
    
    setCreating(true)
    setError('')
    
    try {
      await filesAPI.createFolder(folderName, parentId)
      setFolderName('') // Limpiar el input
      onSuccess()
    } catch (err) {
      setError(err.response?.data?.error || 'Error al crear la carpeta')
    } finally {
      setCreating(false)
    }
  }

  const handleOpenChange = (newOpen) => {
    if (!newOpen && !creating) {
      setFolderName('')
      setError('')
      onClose()
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <FolderPlus className="w-5 h-5" />
            Nueva Carpeta
          </DialogTitle>
          <DialogDescription>
            Crea una nueva carpeta para organizar tus archivos.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {error && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <div className="space-y-2">
            <Label htmlFor="folder-name">Nombre de la carpeta</Label>
            <Input
              id="folder-name"
              value={folderName}
              onChange={(e) => setFolderName(e.target.value)}
              placeholder="Mi carpeta"
              disabled={creating}
              autoFocus
              onKeyDown={(e) => {
                if (e.key === 'Enter' && folderName.trim() && !creating) {
                  handleCreate()
                }
              }}
            />
          </div>
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button
            type="button"
            variant="outline"
            onClick={onClose}
            disabled={creating}
          >
            Cancelar
          </Button>
          <Button
            type="button"
            onClick={handleCreate}
            disabled={!folderName.trim() || creating}
          >
            {creating ? 'Creando...' : 'Crear'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}