import { useState, useEffect } from 'react'
import { Upload, FolderPlus, RefreshCw, Loader2 } from 'lucide-react'
import { filesAPI } from '../api/files'
import FileList from '../components/FileList'
import UploadModal from '../components/UploadModal'
import CreateFolderModal from '../components/CreateFolderModal'
import FileVersionsModal from '../components/FileVersionsModal'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
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
import { useToast } from '@/hooks/use-toast'

export default function Dashboard() {
  const [files, setFiles] = useState([])
  const [loading, setLoading] = useState(true)
  const [currentFolder, setCurrentFolder] = useState(null)
  const [showUploadModal, setShowUploadModal] = useState(false)
  const [showFolderModal, setShowFolderModal] = useState(false)
  const [showVersionsModal, setShowVersionsModal] = useState(false)
  const [selectedFile, setSelectedFile] = useState(null)
  const [deleteDialog, setDeleteDialog] = useState({ open: false, fileId: null, fileName: '' })
  const [breadcrumbs, setBreadcrumbs] = useState([{ id: null, name: 'Mis Archivos' }])
  const [refreshing, setRefreshing] = useState(false)
  const { toast } = useToast()

  const loadFiles = async (folderId = null) => {
    setLoading(true)
    try {
      const response = await filesAPI.list(folderId)
      if (response.success) {
        setFiles(response.data || [])
      }
    } catch (error) {
      console.error('Error loading files:', error)
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'No se pudieron cargar los archivos',
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadFiles(currentFolder)
  }, [currentFolder])

  const handleFolderClick = (folder) => {
    setCurrentFolder(folder.id)
    setBreadcrumbs([...breadcrumbs, { id: folder.id, name: folder.name }])
  }

  const handleBreadcrumbClick = (index) => {
    const newBreadcrumbs = breadcrumbs.slice(0, index + 1)
    setBreadcrumbs(newBreadcrumbs)
    setCurrentFolder(newBreadcrumbs[newBreadcrumbs.length - 1].id)
  }

  const handleUploadSuccess = () => {
    setShowUploadModal(false)
    loadFiles(currentFolder)
  }

  const handleFolderCreated = () => {
    setShowFolderModal(false)
    loadFiles(currentFolder)
  }

  const handleDeleteClick = (fileId, fileName) => {
    setDeleteDialog({ open: true, fileId, fileName })
  }

  const handleDeleteConfirm = async () => {
    const fileId = deleteDialog.fileId
    setDeleteDialog({ open: false, fileId: null, fileName: '' })

    try {
      await filesAPI.delete(fileId)
      toast({
        title: 'Archivo eliminado',
        description: 'El archivo se eliminó correctamente',
      })
      loadFiles(currentFolder)
    } catch (error) {
      console.error('Error deleting file:', error)
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'No se pudo eliminar el archivo',
      })
    }
  }

  const handleDownload = async (file) => {
    try {
      const response = await filesAPI.download(file.id)
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', file.original_name || file.name)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
      
      toast({
        title: 'Descarga iniciada',
        description: `${file.original_name} se está descargando`,
      })
    } catch (error) {
      console.error('Error downloading file:', error)
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'No se pudo descargar el archivo',
      })
    }
  }

  const handleViewVersions = (file) => {
    setSelectedFile(file)
    setShowVersionsModal(true)
  }

  const handleVersionRestored = () => {
    loadFiles(currentFolder)
  }

  const handleRefresh = async () => {
    setRefreshing(true)
    await loadFiles(currentFolder)
    setRefreshing(false)
    toast({
      title: 'Lista actualizada',
      description: 'Los archivos se han actualizado',
    })
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumbs */}
      <Breadcrumb>
        <BreadcrumbList>
          {breadcrumbs.map((crumb, index) => (
            <div key={crumb.id || 'root'} className="flex items-center">
              {index > 0 && <BreadcrumbSeparator />}
              <BreadcrumbItem>
                {index === breadcrumbs.length - 1 ? (
                  <BreadcrumbPage>{crumb.name}</BreadcrumbPage>
                ) : (
                  <BreadcrumbLink
                    onClick={() => handleBreadcrumbClick(index)}
                    className="cursor-pointer"
                  >
                    {crumb.name}
                  </BreadcrumbLink>
                )}
              </BreadcrumbItem>
            </div>
          ))}
        </BreadcrumbList>
      </Breadcrumb>

      {/* Actions */}
      <div className="flex flex-wrap gap-3">
        <Button onClick={() => setShowUploadModal(true)}>
          <Upload className="mr-2 h-4 w-4" />
          Subir Archivo
        </Button>

        <Button variant="outline" onClick={() => setShowFolderModal(true)}>
          <FolderPlus className="mr-2 h-4 w-4" />
          Nueva Carpeta
        </Button>

        <Button
          variant="outline"
          onClick={handleRefresh}
          disabled={refreshing}
        >
          <RefreshCw className={`mr-2 h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
          Actualizar
        </Button>
      </div>

      {/* Files List */}
      {loading ? (
        <div className="flex flex-col items-center justify-center py-12 space-y-4">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Cargando archivos...</p>
        </div>
      ) : (
        <FileList
          files={files}
          onFolderClick={handleFolderClick}
          onDelete={(fileId) => {
            const file = files.find(f => f.id === fileId)
            handleDeleteClick(fileId, file?.original_name || file?.name || 'este archivo')
          }}
          onDownload={handleDownload}
          onViewVersions={handleViewVersions}
        />
      )}

      {/* Modals */}
      <UploadModal
        open={showUploadModal}
        onClose={() => setShowUploadModal(false)}
        onSuccess={handleUploadSuccess}
        parentId={currentFolder}
      />

      <CreateFolderModal
        open={showFolderModal}
        onClose={() => setShowFolderModal(false)}
        onSuccess={handleFolderCreated}
        parentId={currentFolder}
      />

      {selectedFile && (
        <FileVersionsModal
          open={showVersionsModal}
          file={selectedFile}
          onClose={() => {
            setShowVersionsModal(false)
            setSelectedFile(null)
          }}
          onVersionRestored={handleVersionRestored}
        />
      )}

      {/* Delete Confirmation Dialog */}
      <AlertDialog 
        open={deleteDialog.open} 
        onOpenChange={(open) => !open && setDeleteDialog({ open: false, fileId: null, fileName: '' })}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>¿Estás seguro?</AlertDialogTitle>
            <AlertDialogDescription>
              Esta acción eliminará permanentemente <strong>{deleteDialog.fileName}</strong>.
              Esta acción no se puede deshacer.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteConfirm}
              className="bg-destructive hover:bg-destructive/90"
            >
              Eliminar
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}