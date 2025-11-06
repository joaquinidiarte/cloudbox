import { useState, useEffect } from 'react'
import { Download, RotateCcw, Trash2, Clock, AlertCircle, Loader2 } from 'lucide-react'
import { filesAPI } from '../api/files'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Card, CardContent } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { useToast } from '@/hooks/use-toast'

export default function FileVersionsModal({ open, file, onClose, onVersionRestored }) {
  const [versions, setVersions] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [actionLoading, setActionLoading] = useState(null)
  const [alertDialog, setAlertDialog] = useState({ open: false, type: '', version: null })
  const { toast } = useToast()

  useEffect(() => {
    if (open && file?.id) {
      loadVersions()
    }
  }, [open, file?.id])

  const loadVersions = async () => {
    setLoading(true)
    setError('')
    try {
      const response = await filesAPI.getVersions(file.id)
      if (response.success) {
        setVersions(response.data || [])
      }
    } catch (err) {
      setError('Error al cargar versiones')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleDownloadVersion = async (version) => {
    setActionLoading(`download-${version}`)
    try {
      const response = await filesAPI.downloadVersion(file.id, version)
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `${file.original_name}_v${version}`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
      
      toast({
        title: 'Descarga iniciada',
        description: `Versión ${version} descargada correctamente`,
      })
    } catch (err) {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'No se pudo descargar la versión',
      })
      console.error(err)
    } finally {
      setActionLoading(null)
    }
  }

  const handleRestoreVersion = async () => {
    const version = alertDialog.version
    setAlertDialog({ open: false, type: '', version: null })
    setActionLoading(`restore-${version}`)
    
    try {
      await filesAPI.restoreVersion(file.id, version)
      toast({
        title: 'Versión restaurada',
        description: `La versión ${version} ahora es la actual`,
      })
      onVersionRestored?.()
      loadVersions()
    } catch (err) {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'No se pudo restaurar la versión',
      })
      console.error(err)
    } finally {
      setActionLoading(null)
    }
  }

  const handleDeleteVersion = async () => {
    const version = alertDialog.version
    setAlertDialog({ open: false, type: '', version: null })
    setActionLoading(`delete-${version}`)
    
    try {
      await filesAPI.deleteVersion(file.id, version)
      toast({
        title: 'Versión eliminada',
        description: `La versión ${version} ha sido eliminada`,
      })
      loadVersions()
    } catch (err) {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: err.response?.data?.error || 'No se pudo eliminar la versión',
      })
      console.error(err)
    } finally {
      setActionLoading(null)
    }
  }

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString('es-ES', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <>
      <Dialog open={open} onOpenChange={onClose}>
        <DialogContent className="max-w-3xl max-h-[85vh] flex flex-col p-0">
          <DialogHeader className="px-6 pt-6">
            <DialogTitle className="flex items-center gap-2">
              <Clock className="w-5 h-5" />
              Historial de Versiones
            </DialogTitle>
            <DialogDescription>{file?.original_name}</DialogDescription>
          </DialogHeader>

          <Separator />

          <ScrollArea className="flex-1 px-6">
            {loading ? (
              <div className="flex flex-col items-center justify-center py-12">
                <Loader2 className="w-8 h-8 animate-spin text-primary mb-2" />
                <p className="text-sm text-muted-foreground">Cargando versiones...</p>
              </div>
            ) : error ? (
              <Alert variant="destructive" className="my-4">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            ) : versions.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-12">
                <Clock className="w-16 h-16 text-muted-foreground/50 mb-4" />
                <p className="text-muted-foreground">No hay versiones disponibles</p>
              </div>
            ) : (
              <div className="space-y-3 py-4">
                {versions.map((version) => (
                  <Card
                    key={version.version}
                    className={version.is_current ? 'border-primary' : ''}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-start justify-between gap-4">
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-2">
                            <h3 className="font-semibold">Versión {version.version}</h3>
                            {version.is_current && (
                              <Badge>Actual</Badge>
                            )}
                          </div>
                          
                          <div className="space-y-1 text-sm text-muted-foreground">
                            <div className="flex items-center gap-1">
                              <Clock className="w-3.5 h-3.5" />
                              <span>{formatDate(version.uploaded_at)}</span>
                            </div>
                            <p>Tamaño: {formatBytes(version.size)}</p>
                            {version.comment && (
                              <p className="text-foreground italic mt-2">
                                "{version.comment}"
                              </p>
                            )}
                          </div>
                        </div>

                        <div className="flex items-center gap-1 flex-shrink-0">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDownloadVersion(version.version)}
                            disabled={actionLoading === `download-${version.version}`}
                            title="Descargar versión"
                          >
                            {actionLoading === `download-${version.version}` ? (
                              <Loader2 className="w-4 h-4 animate-spin" />
                            ) : (
                              <Download className="w-4 h-4" />
                            )}
                          </Button>

                          {!version.is_current && (
                            <>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() =>
                                  setAlertDialog({ 
                                    open: true, 
                                    type: 'restore', 
                                    version: version.version 
                                  })
                                }
                                disabled={actionLoading === `restore-${version.version}`}
                                title="Restaurar versión"
                                className="text-blue-600 hover:text-blue-600 hover:bg-blue-50"
                              >
                                {actionLoading === `restore-${version.version}` ? (
                                  <Loader2 className="w-4 h-4 animate-spin" />
                                ) : (
                                  <RotateCcw className="w-4 h-4" />
                                )}
                              </Button>

                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() =>
                                  setAlertDialog({ 
                                    open: true, 
                                    type: 'delete', 
                                    version: version.version 
                                  })
                                }
                                disabled={actionLoading === `delete-${version.version}`}
                                title="Eliminar versión"
                                className="text-destructive hover:text-destructive hover:bg-destructive/10"
                              >
                                {actionLoading === `delete-${version.version}` ? (
                                  <Loader2 className="w-4 h-4 animate-spin" />
                                ) : (
                                  <Trash2 className="w-4 h-4" />
                                )}
                              </Button>
                            </>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </ScrollArea>

          <Separator />

          <DialogFooter className="px-6 pb-6">
            <div className="flex items-center justify-between w-full">
              <p className="text-sm text-muted-foreground">
                <strong>{versions.length}</strong> versión{versions.length !== 1 ? 'es' : ''} en total
              </p>
              <Button variant="outline" onClick={onClose}>
                Cerrar
              </Button>
            </div>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* AlertDialog para confirmaciones */}
      <AlertDialog 
        open={alertDialog.open} 
        onOpenChange={(open) => !open && setAlertDialog({ open: false, type: '', version: null })}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {alertDialog.type === 'restore' 
                ? '¿Restaurar esta versión?' 
                : '¿Eliminar esta versión?'}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {alertDialog.type === 'restore' 
                ? `La versión ${alertDialog.version} se convertirá en la versión actual del archivo.`
                : `La versión ${alertDialog.version} será eliminada permanentemente. Esta acción no se puede deshacer.`}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              onClick={alertDialog.type === 'restore' ? handleRestoreVersion : handleDeleteVersion}
              className={alertDialog.type === 'delete' ? 'bg-destructive hover:bg-destructive/90' : ''}
            >
              {alertDialog.type === 'restore' ? 'Restaurar' : 'Eliminar'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}