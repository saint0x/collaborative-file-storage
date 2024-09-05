import { motion } from 'framer-motion'
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"
import { Badge } from "./ui/badge"
import { Button } from "./ui/button"
import { FileText, Image, Music, Mic, BookOpen, Users, Heart } from 'lucide-react'

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
  return (
    <motion.div
      className={`w-64 h-96 bg-white rounded-3xl shadow-xl overflow-hidden transform transition-all duration-500 ${isActive ? 'scale-105' : 'scale-95 opacity-70'}`}
      whileHover={{ scale: isActive ? 1.05 : 1 }}
    >
      <div className="p-4">
        <div className="flex flex-col items-center mb-4">
          <Avatar className="h-16 w-16 mb-2 ring-2 ring-violet-500 ring-offset-2">
            <AvatarImage src={friend.avatar} alt={friend.name} />
            <AvatarFallback>{friend.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
          </Avatar>
          <h2 className="text-lg font-semibold text-center">{friend.name}</h2>
          <div className="flex flex-wrap justify-center gap-1 mt-1">
            {friend.contexts.map((context) => (
              <Badge key={context} variant="secondary" className="text-xs rounded-full px-2 py-0.5">
                {context}
              </Badge>
            ))}
          </div>
        </div>
        <h3 className="font-semibold mb-2 text-sm">Shared Files</h3>
        <div className="relative h-36">
          {friend.sharedItems.map((item, itemIndex) => {
            const Icon = iconMap[item.icon as keyof typeof iconMap] || FileText
            return (
              <motion.div
                key={item.name}
                className={`absolute w-20 h-20 rounded-xl shadow-md ${item.color} flex items-center justify-center`}
                initial={{ rotate: 0, x: 0, y: 0 }}
                animate={{
                  rotate: itemIndex * 5 - 5,
                  x: itemIndex * 10,
                  y: itemIndex * 5,
                }}
                whileHover={{ scale: 1.1, rotate: 0, zIndex: 10 }}
                transition={{ type: 'spring', stiffness: 300, damping: 20 }}
              >
                <div className="text-white text-center">
                  <Icon className="h-6 w-6 mx-auto mb-1" />
                  <span className="text-xs font-medium">{item.name}</span>
                </div>
              </motion.div>
            )
          })}
        </div>
        <div className="flex justify-between items-center mt-4">
          <Button variant="outline" className="rounded-full text-xs px-2 py-1">
            <Users className="h-3 w-3 mr-1" />
            View All
          </Button>
          <Button variant="ghost" size="icon" className="rounded-full">
            <Heart className="h-4 w-4 text-red-500" />
          </Button>
        </div>
      </div>
    </motion.div>
  )
}