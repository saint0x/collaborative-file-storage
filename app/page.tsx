'use client'

import { useState, useRef, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from 'lucide-react'
import Header from '@/components/Header'
import FriendCard from '@/components/FriendCard'

const friends = [
  {
    id: 1,
    name: 'Alice Johnson',
    avatar: '/placeholder.svg?height=128&width=128',
    contexts: ['Work', 'Book Club'],
    sharedItems: [
      { type: 'file', name: 'Project Proposal', icon: 'FileText', color: 'bg-blue-500' },
      { type: 'book', name: 'Current Read: 1984', icon: 'BookOpen', color: 'bg-green-500' },
      { type: 'image', name: 'Team Photo', icon: 'Image', color: 'bg-yellow-500' },
    ]
  },
  {
    id: 2,
    name: 'Bob Smith',
    avatar: '/placeholder.svg?height=128&width=128',
    contexts: ['Music', 'Photography'],
    sharedItems: [
      { type: 'playlist', name: 'Road Trip Mix', icon: 'Music', color: 'bg-purple-500' },
      { type: 'album', name: 'Summer Vacation', icon: 'Image', color: 'bg-pink-500' },
      { type: 'file', name: 'Photo Editing Tips', icon: 'FileText', color: 'bg-blue-500' },
    ]
  },
  {
    id: 3,
    name: 'Charlie Brown',
    avatar: '/placeholder.svg?height=128&width=128',
    contexts: ['Work', 'Fitness'],
    sharedItems: [
      { type: 'file', name: 'Meeting Notes', icon: 'FileText', color: 'bg-blue-500' },
      { type: 'voice', name: 'Workout Plan', icon: 'Mic', color: 'bg-red-500' },
      { type: 'image', name: 'Gym Progress', icon: 'Image', color: 'bg-green-500' },
    ]
  },
]

export default function Home() {
  const [activeIndex, setActiveIndex] = useState(0)
  const containerRef = useRef<HTMLDivElement>(null)

  const nextCard = () => {
    setActiveIndex((prevIndex) => (prevIndex + 1) % friends.length)
  }

  const prevCard = () => {
    setActiveIndex((prevIndex) => (prevIndex - 1 + friends.length) % friends.length)
  }

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'ArrowRight') {
        nextCard()
      } else if (event.key === 'ArrowLeft') {
        prevCard()
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [])

  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-br from-pink-100 to-violet-200">
      <Header />
      <main className="flex-1 flex items-center justify-center p-6 mt-20">
        <div className="w-full max-w-4xl h-[500px] relative">
          <div ref={containerRef} className="w-full h-full perspective-1000">
            <AnimatePresence>
              {friends.map((friend, index) => {
                const rotationAngle = ((index - activeIndex) / friends.length) * -360
                const zIndex = friends.length - Math.abs(index - activeIndex)
                const scale = 1 - Math.abs(index - activeIndex) * 0.1

                return (
                  <motion.div
                    key={friend.id}
                    className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"
                    style={{
                      zIndex,
                      rotateY: `${rotationAngle}deg`,
                      translateZ: '300px',
                    }}
                    initial={{ opacity: 0, rotateY: `${rotationAngle + 60}deg` }}
                    animate={{ opacity: 1, rotateY: `${rotationAngle}deg`, scale }}
                    exit={{ opacity: 0, rotateY: `${rotationAngle - 60}deg` }}
                    transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                  >
                    <FriendCard friend={friend} isActive={index === activeIndex} />
                  </motion.div>
                )
              })}
            </AnimatePresence>
          </div>

          {/* Navigation Buttons */}
          <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex justify-center space-x-4">
            <Button onClick={prevCard} size="icon" className="rounded-full bg-white bg-opacity-50 hover:bg-opacity-75 transition-all duration-300">
              <ChevronLeft className="h-6 w-6 text-gray-800" />
            </Button>
            <Button onClick={nextCard} size="icon" className="rounded-full bg-white bg-opacity-50 hover:bg-opacity-75 transition-all duration-300">
              <ChevronRight className="h-6 w-6 text-gray-800" />
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}