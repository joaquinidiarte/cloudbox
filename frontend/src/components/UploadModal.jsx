import { useState, useRef } from 'react'
import { Upload, AlertCircle, File, X, Loader2 } from 'lucide-react'
import { filesAPI } from '../api/files'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { useToast } from '@/hooks/use-toast'

export default function UploadModal({ open, onClose, onSuccess, parentId }) {
  const [selectedFile, setSelectedFile] = useState(null)
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [error, setError] = useState('')
  const fileInputRef = useRef(null)
  const { toast } = useToast()

  const handleFileSelect = (e) => {
    const file = e.target.files?.[0]
    if (file) {
      // Validar tamaño (100MB)
      if (file.size > 100 * 1024 * 1024) {
        setError('El archivo es demasiado grande (máximo 100MB)')
        setSelectedFile(null)
        return
      }
      setSelectedFile(file)
      setError('')
    }
  }

  const handleUpload = async () => {
    if (!selectedFile) return

    setUploading(true)
    setError('')
    setUploadProgress(0)

    try {
      // Simulación de progreso si tu API no lo soporta
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => {
          if (prev >= 90) {
            clearInterval(progressInterval)
            return prev
          }
          return prev + 10
        })
      }, 200)

      await filesAPI.upload(selectedFile, parentId)
      
      clearInterval(progressInterval)
      setUploadProgress(100)

      toast({
        title: 'Archivo subido',
        description: `${selectedFile.name} se subió correctamente`,
      })

      // Resetear y cerrar
      setTimeout(() => {
        setSelectedFile(null)
        setUploadProgress(0)
        onSuccess()
      }, 500)
    } catch (err) {
      setError(err.response?.data?.error || 'Error al subir el archivo')
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'No se pudo subir el archivo',
      })
    } finally {
      setUploading(false)
    }
  }

  const handleOpenChange = (newOpen) => {
    if (!newOpen && !uploading) {
      setSelectedFile(null)
      setError('')
      setUploadProgress(0)
      onClose()
    }
  }

  const handleDragOver = (e) => {
    e.preventDefault()
    e.stopPropagation()
  }

  const handleDrop = (e) => {
    e.preventDefault()
    e.stopPropagation()

    const file = e.dataTransfer.files?.[0]
    if (file) {
      if (file.size > 100 * 1024 * 1024) {
        setError('El archivo es demasiado grande (máximo 100MB)')
        return
      }
      setSelectedFile(file)
      setError('')
    }
  }

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Upload className="w-5 h-5" />
            Subir Archivo
          </DialogTitle>
          <DialogDescription>
            Selecciona un archivo para subir a tu almacenamiento (máximo 100MB)
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {error && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {/* Drop Zone / File Input */}
          <div
            onDragOver={handleDragOver}
            onDrop={handleDrop}
            onClick={() => !uploading && fileInputRef.current?.click()}
            className={`
              border-2 border-dashed rounded-lg p-8 text-center cursor-pointer
              transition-colors
              ${selectedFile 
                ? 'border-primary bg-primary/5' 
                : 'border-muted-foreground/25 hover:border-primary hover:bg-accent'
              }
              ${uploading ? 'pointer-events-none opacity-60' : ''}
            `}
          >
            <input
              ref={fileInputRef}
              type="file"
              onChange={handleFileSelect}
              className="hidden"
              disabled={uploading}
            />

            {selectedFile ? (
              <div className="space-y-2">
                <File className="w-12 h-12 mx-auto text-primary" />
                <div>
                  <p className="font-medium text-sm">{selectedFile.name}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {formatBytes(selectedFile.size)}
                  </p>
                </div>
                {!uploading && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={(e) => {
                      e.stopPropagation()
                      setSelectedFile(null)
                      setError('')
                    }}
                    className="mt-2"
                  >
                    <X className="w-4 h-4 mr-2" />
                    Cambiar archivo
                  </Button>
                )}
              </div>
            ) : (
              <div className="space-y-2">
                <Upload className="w-12 h-12 mx-auto text-muted-foreground" />
                <div>
                  <p className="font-medium">
                    Haz clic o arrastra un archivo aquí
                  </p>
                  <p className="text-sm text-muted-foreground mt-1">
                    Máximo 100MB
                  </p>
                </div>
              </div>
            )}
          </div>

          {/* Progress Bar */}
          {uploading && (
            <div className="space-y-2">
              <Progress value={uploadProgress} className="h-2" />
              <p className="text-sm text-center text-muted-foreground">
                Subiendo... {uploadProgress}%
              </p>
            </div>
          )}
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button
            type="button"
            variant="outline"
            onClick={onClose}
            disabled={uploading}
          >
            Cancelar
          </Button>
          <Button
            type="button"
            onClick={handleUpload}
            disabled={!selectedFile || uploading}
          >
            {uploading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Subiendo...
              </>
            ) : (
              <>
                <Upload className="mr-2 h-4 w-4" />
                Subir
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}