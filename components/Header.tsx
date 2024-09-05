import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Plus, Search } from "lucide-react"

export default function Header() {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-6 py-4 bg-white bg-opacity-80 backdrop-blur-md rounded-b-3xl shadow-lg">
      <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500">
        Your Circle
      </h1>
      <div className="flex items-center space-x-4">
        <Input className="w-64 rounded-full border-2 border-gray-300 focus:border-violet-500 transition-colors" placeholder="Search connections..." />
        <Button size="icon" className="rounded-full bg-gradient-to-r from-pink-500 to-violet-500 text-white hover:from-pink-600 hover:to-violet-600 transition-all duration-300">
          <Plus className="h-5 w-5" />
        </Button>
        <Avatar className="h-12 w-12 ring-2 ring-violet-500 ring-offset-2">
          <AvatarImage src="/placeholder.svg?height=48&width=48" alt="Your Avatar" />
          <AvatarFallback>YA</AvatarFallback>
        </Avatar>
      </div>
    </header>
  )
}