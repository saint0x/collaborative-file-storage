import { motion } from 'framer-motion'
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"
import { Badge } from "./ui/badge"
import { Button } from "./ui/button"
import { FileText, Image, Music, Mic, BookOpen, Users, Heart } from 'lucide-react'
import { useState } from 'react'
import FileExplorer from './FileExplorer'

const iconMap = {
  FileText,
  Image,
  Music,
  Mic,
  BookOpen,
}

interface SharedItem {
  type: string
  name: string
  icon: string
  color: string
}

interface Friend {
  id: number
  name: string
  avatar: string
  contexts: string[]
  sharedItems: SharedItem[]
}

interface FriendCardProps {
  friend: Friend
  isActive: boolean
}

export default function FriendCard({ friend, isActive }: FriendCardProps) {
  const [showFileExplorer, setShowFileExplorer] = useState(false)
  const [selectedFile, setSelectedFile] = useState<SharedItem | null>(null)

  const handleFileClick = (item: SharedItem) => {
    setSelectedFile(item)
    setShowFileExplorer(true)
  }

  return (
    <>
      <div
        className={`w-72 h-[450px] bg-gradient-to-br from-gray-100 to-gray-300 rounded-3xl shadow-xl overflow-hidden transform transition-all duration-500`}
        style={{
          filter: isActive ? 'none' : 'blur(1px)',
        }}
      >
        <div className="p-6 flex flex-col h-full">
          <div className="flex items-center mb-4">
            <Avatar className="h-16 w-16 mr-4 ring-2 ring-gray-400 ring-offset-2">
              <AvatarImage src={friend.avatar} alt={friend.name} />
              <AvatarFallback>{friend.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
            </Avatar>
            <div>
              <h2 className="text-xl font-semibold text-gray-800">{friend.name}</h2>
              <div className="flex flex-wrap gap-1 mt-1">
                {friend.contexts.map((context) => (
                  <Badge key={context} variant="secondary" className="text-xs rounded-full px-2 py-0.5 bg-gray-200 text-gray-700">
                    {context}
                  </Badge>
                ))}
              </div>
            </div>
          </div>
          <h3 className="font-semibold mb-3 text-sm text-gray-600">Shared Files</h3>
          <div className="flex-grow relative">
            {friend.sharedItems.map((item, itemIndex) => {
              const Icon = iconMap[item.icon as keyof typeof iconMap] || FileText
              return (
                <motion.div
                  key={item.name}
                  className={`absolute w-24 h-24 rounded-xl shadow-md bg-gradient-to-br from-gray-200 to-gray-400 flex items-center justify-center cursor-pointer`}
                  initial={{ rotate: 0, x: 0, y: 0 }}
                  animate={{
                    rotate: itemIndex * 5 - 5,
                    x: itemIndex * 15,
                    y: itemIndex * 10,
                  }}
                  whileHover={{ scale: 1.1, rotate: 0, zIndex: 10 }}
                  transition={{ type: 'spring', stiffness: 300, damping: 20 }}
                  onClick={() => handleFileClick(item)}
                >
                  <div className="text-gray-700 text-center">
                    <Icon className="h-8 w-8 mx-auto mb-2" />
                    <span className="text-xs font-medium">{item.name}</span>
                  </div>
                </motion.div>
              )
            })}
          </div>
          <div className="flex justify-between items-center mt-4">
            <Button variant="outline" className="rounded-full text-sm px-4 py-2 bg-gradient-to-r from-gray-100 to-gray-300 text-gray-700 border-gray-400 hover:bg-gradient-to-r hover:from-gray-200 hover:to-gray-400">
              <Users className="h-4 w-4 mr-2" />
              View All
            </Button>
            <Button variant="ghost" size="icon" className="rounded-full bg-gradient-to-r from-gray-100 to-gray-300 text-gray-700 hover:bg-gradient-to-r hover:from-gray-200 hover:to-gray-400">
              <Heart className="h-5 w-5 text-red-400" />
            </Button>
          </div>
        </div>
      </div>
      {showFileExplorer && selectedFile && (
        <FileExplorer
          isOpen={showFileExplorer}
          onClose={() => setShowFileExplorer(false)}
          initialCategory={selectedFile.type}
          initialFile={selectedFile.name}
        />
      )}
    </>
  )
}