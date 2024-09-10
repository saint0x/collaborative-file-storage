'use client'

import { useState } from 'react'
import Header from '@/components/Header'
import FileExplorer from '@/components/FileExplorer'

export default function FilesPage() {
  const [isFileExplorerOpen, setIsFileExplorerOpen] = useState(true)

  const handleCloseFileExplorer = () => {
    setIsFileExplorerOpen(false)
  }

  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-br from-gray-100 to-gray-200">
      <Header />
      <main className="flex-1 flex flex-col items-center justify-center p-6 mt-20">
        <FileExplorer 
          isOpen={isFileExplorerOpen} 
          onClose={handleCloseFileExplorer}
        />
      </main>
    </div>
  )
}