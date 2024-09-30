"use client"

import { useState, useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Plus, Search, Settings, LogOut, User, HelpCircle } from "lucide-react"

export default function Header() {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const toggleDropdown = () => {
    setIsDropdownOpen(!isDropdownOpen)
  }

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [])

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
        <div className="relative" ref={dropdownRef}>
          <Avatar 
            className="h-12 w-12 ring-2 ring-gray-400 ring-offset-2 cursor-pointer transition-all duration-300 hover:ring-gray-500"
            onClick={toggleDropdown}
          >
            <AvatarImage src="/placeholder.svg?height=48&width=48" alt="Your Avatar" />
            <AvatarFallback>YA</AvatarFallback>
          </Avatar>
          <AnimatePresence>
            {isDropdownOpen && (
              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: 0.2 }}
                className="absolute right-0 mt-2 w-56 rounded-xl shadow-lg bg-white ring-1 ring-black ring-opacity-5 overflow-hidden"
              >
                <div className="py-1" role="menu" aria-orientation="vertical" aria-labelledby="options-menu">
                  {[
                    { icon: User, label: 'Your Profile' },
                    { icon: Settings, label: 'Settings' },
                    { icon: HelpCircle, label: 'Help & Support' },
                    { icon: LogOut, label: 'Sign out' },
                  ].map((item, index) => (
                    <a
                      key={index}
                      href="#"
                      className="block px-4 py-2 text-sm text-gray-700 hover:bg-gradient-to-r hover:from-gray-100 hover:to-gray-200 hover:text-gray-900 transition-colors duration-200 flex items-center"
                      role="menuitem"
                    >
                      <item.icon className="mr-3 h-4 w-4" />
                      {item.label}
                    </a>
                  ))}
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </header>
  )
}