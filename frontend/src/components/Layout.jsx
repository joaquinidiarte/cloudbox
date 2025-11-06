import { LogOut, HardDrive } from 'lucide-react'
import { useAuthStore } from '../store/authStore'
import { useNavigate } from 'react-router-dom'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'

export default function Layout({ children }) {
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  const storagePercentage = user ? (user.storage_used / user.storage_limit) * 100 : 0
  
  const getStorageColor = () => {
    if (storagePercentage > 90) return '[&>div]:bg-destructive'
    if (storagePercentage > 70) return '[&>div]:bg-yellow-500'
    return '[&>div]:bg-primary'
  }

  const getUserInitials = () => {
    if (!user) return 'U'
    const firstInitial = user.first_name?.[0] || ''
    const lastInitial = user.last_name?.[0] || ''
    return `${firstInitial}${lastInitial}`.toUpperCase() || 'U'
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="w-full px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 items-center justify-between max-w-7xl mx-auto">
            <div className="flex items-center gap-3">
              <img src="/cloudbox-logo.png" alt="CloudBox Logo" width={60} height={60} />
            </div>

            <div className="flex items-center gap-4">
              {/* Storage Info - Desktop */}
              <div className="hidden md:flex items-center gap-2 text-sm text-muted-foreground">
                <HardDrive className="w-4 h-4" />
                <span>
                  {user && formatBytes(user.storage_used)} / {user && formatBytes(user.storage_limit)}
                </span>
              </div>

              {/* User Menu */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="relative h-10 gap-2">
                    <Avatar className="h-8 w-8">
                      <AvatarFallback className="text-xs">
                        {getUserInitials()}
                      </AvatarFallback>
                    </Avatar>
                    <span className="hidden sm:inline text-sm font-medium">
                      {user?.first_name} {user?.last_name}
                    </span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56" align="end" forceMount>
                  <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                      <p className="text-sm font-medium leading-none">
                        {user?.first_name} {user?.last_name}
                      </p>
                      <p className="text-xs leading-none text-muted-foreground">
                        {user?.email}
                      </p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  
                  {/* Storage info in mobile menu */}
                  <div className="md:hidden px-2 py-2">
                    <div className="flex items-center gap-2 text-xs text-muted-foreground mb-1">
                      <HardDrive className="w-3 h-3" />
                      <span>
                        {user && formatBytes(user.storage_used)} / {user && formatBytes(user.storage_limit)}
                      </span>
                    </div>
                  </div>
                  <DropdownMenuSeparator className="md:hidden" />
                  
                  <DropdownMenuItem onClick={handleLogout}>
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Cerrar sesión</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>
      </header>

      {/* Storage Bar */}
      {user && (
        <div className="border-b bg-muted/40">
          <div className="w-full px-4 sm:px-6 lg:px-8">
            <div className="max-w-7xl mx-auto py-3">
              <div className="flex items-center gap-3">
                <div className="flex-1">
                  <Progress 
                    value={Math.min(storagePercentage, 100)} 
                    className={`h-2 ${getStorageColor()}`}
                  />
                </div>
                <span className="text-xs text-muted-foreground whitespace-nowrap">
                  {storagePercentage.toFixed(1)}% usado
                </span>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Main Content */}
      <main className="w-full px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto py-8">
          {children}
        </div>
      </main>
    </div>
  )
}