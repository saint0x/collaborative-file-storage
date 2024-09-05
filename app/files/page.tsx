import Header from '@/components/Header'
import FileExplorer from '@/components/FileExplorer'

export default function FilesPage() {
  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-br from-pink-100 to-violet-200">
      <Header />
      <main className="flex-1 flex flex-col items-center justify-center p-6 mt-20">
        <FileExplorer />
      </main>
    </div>
  )
}