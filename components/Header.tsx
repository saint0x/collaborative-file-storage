import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Plus, Search } from "lucide-react"

export default function Header() {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-6 py-4 bg-gradient-to-r from-gray-100 to-gray-200 bg-opacity-80 backdrop-blur-md rounded-b-3xl shadow-lg">
      <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-gray-600 to-gray-800">
        Your Circle
      </h1>
      <div className="flex items-center space-x-4">
        <div className="relative">
          <Input 
            className="w-64 rounded-full border-2 border-gray-300 focus:border-gray-500 transition-colors pl-10" 
            placeholder="Search connections..." 
          />
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
        </div>
        <Button 
          size="icon" 
          className="rounded-full bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500 transition-all duration-300 shadow-md hover:shadow-lg"
        >
          <Plus className="h-5 w-5" />
        </Button>
        <Avatar className="h-12 w-12 ring-2 ring-gray-400 ring-offset-2">
          <AvatarImage src="/placeholder.svg?height=48&width=48" alt="Your Avatar" />
          <AvatarFallback>YA</AvatarFallback>
        </Avatar>
      </div>
    </header>
  )
}