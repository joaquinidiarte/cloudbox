import { File, Folder, Download, Trash2, Clock, MoreVertical } from 'lucide-react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

export default function FileList({ files, onFolderClick, onDelete, onDownload, onViewVersions }) {
  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const getFileIcon = (file) => {
    if (file.is_folder) {
      return <Folder className="w-8 h-8 text-yellow-500" />
    }
    return <File className="w-8 h-8 text-blue-500" />
  }

  if (files.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center border rounded-lg bg-muted/50">
        <Folder className="w-16 h-16 text-muted-foreground/50 mb-4" />
        <p className="text-muted-foreground">No hay archivos en esta carpeta</p>
      </div>
    )
  }

  return (
    <div className="border rounded-lg">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Nombre</TableHead>
            <TableHead className="hidden sm:table-cell">Tamaño</TableHead>
            <TableHead className="hidden md:table-cell">Modificado</TableHead>
            <TableHead className="text-right">Acciones</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {files.map((file) => (
            <TableRow key={file.id} className="group">
              <TableCell>
                <div
                  className={`flex items-center gap-3 ${
                    file.is_folder ? 'cursor-pointer hover:underline' : ''
                  }`}
                  onClick={() => file.is_folder && onFolderClick(file)}
                >
                  {getFileIcon(file)}
                  <div className="min-w-0 flex-1">
                    <p className="font-medium truncate">
                      {file.original_name || file.name}
                    </p>
                    {file.is_folder && (
                      <p className="text-xs text-muted-foreground">Carpeta</p>
                    )}
                  </div>
                </div>
              </TableCell>
              
              <TableCell className="hidden sm:table-cell text-muted-foreground">
                {file.is_folder ? '-' : formatBytes(file.size)}
              </TableCell>
              
              <TableCell className="hidden md:table-cell text-muted-foreground">
                {formatDate(file.updated_at)}
              </TableCell>
              
              <TableCell className="text-right">
                <div className="flex items-center justify-end gap-1">
                  {!file.is_folder && (
                    <>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => onDownload(file)}
                            >
                              <Download className="w-4 h-4" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Descargar</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>

                      {file.version_count > 1 && (
                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => onViewVersions(file)}
                                className="relative"
                              >
                                <Clock className="w-4 h-4" />
                                <Badge 
                                  variant="default" 
                                  className="absolute -top-1 -right-1 h-5 w-5 flex items-center justify-center p-0 text-xs"
                                >
                                  {file.version_count}
                                </Badge>
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                              <p>Ver {file.version_count} versiones</p>
                            </TooltipContent>
                          </Tooltip>
                        </TooltipProvider>
                      )}
                    </>
                  )}

                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => onDelete(file.id)}
                          className="text-destructive hover:text-destructive hover:bg-destructive/10"
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>Eliminar</p>
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}