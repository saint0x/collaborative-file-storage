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
  {
    id: 4,
    name: 'Diana Prince',
    avatar: '/placeholder.svg?height=128&width=128',
    contexts: ['Art', 'Travel'],
    sharedItems: [
      { type: 'image', name: 'Paris Sketches', icon: 'Image', color: 'bg-yellow-500' },
      { type: 'playlist', name: 'Travel Tunes', icon: 'Music', color: 'bg-purple-500' },
      { type: 'file', name: 'Itinerary', icon: 'FileText', color: 'bg-blue-500' },
    ]
  },
  {
    id: 5,
    name: 'Ethan Hunt',
    avatar: '/placeholder.svg?height=128&width=128',
    contexts: ['Sports', 'Movies'],
    sharedItems: [
      { type: 'file', name: 'Game Strategy', icon: 'FileText', color: 'bg-blue-500' },
      { type: 'playlist', name: 'Workout Mix', icon: 'Music', color: 'bg-green-500' },
      { type: 'image', name: 'Team Logo', icon: 'Image', color: 'bg-red-500' },
    ]
  },
]

export default function Home() {
  const [activeIndex, setActiveIndex] = useState(0)
  const [direction, setDirection] = useState(0)
  const containerRef = useRef<HTMLDivElement>(null)

  const nextCard = () => {
    setDirection(1)
    setActiveIndex((prevIndex) => (prevIndex + 1) % friends.length)
  }

  const prevCard = () => {
    setDirection(-1)
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
      <main className="flex-1 flex flex-col items-center justify-center p-6 mt-20 overflow-hidden">
        <div className="w-full max-w-7xl h-[500px] relative mb-16">
          <div ref={containerRef} className="w-full h-full perspective-1000 flex items-center justify-center">
            <AnimatePresence initial={false} custom={direction}>
              {friends.map((friend, index) => {
                const isActive = index === activeIndex
                const offset = (index - activeIndex + friends.length) % friends.length
                const factor = offset > friends.length / 2 ? offset - friends.length : offset

                return (
                  <motion.div
                    key={friend.id}
                    custom={direction}
                    variants={{
                      enter: (direction: number) => ({
                        x: direction > 0 ? '100%' : '-100%',
                        opacity: 0,
                        scale: 0.8,
                        zIndex: 0,
                      }),
                      center: {
                        x: isActive ? 0 : `${factor * 110}%`,
                        opacity: isActive ? 1 : Math.max(0.3, 1 - Math.abs(factor) * 0.3),
                        scale: isActive ? 1 : Math.max(0.8, 1 - Math.abs(factor) * 0.1),
                        zIndex: isActive ? friends.length : friends.length - Math.abs(factor),
                      },
                      exit: (direction: number) => ({
                        x: direction < 0 ? '100%' : '-100%',
                        opacity: 0,
                        scale: 0.8,
                        zIndex: 0,
                      }),
                    }}
                    initial="enter"
                    animate="center"
                    exit="exit"
                    transition={{
                      x: { type: "spring", stiffness: 300, damping: 30 },
                      opacity: { duration: 0.5 },
                      scale: { duration: 0.5 },
                    }}
                    style={{
                      position: 'absolute',
                      left: '50%',
                      top: '50%',
                      translateX: '-50%',
                      translateY: '-50%',
                    }}
                  >
                    <FriendCard friend={friend} isActive={isActive} />
                  </motion.div>
                )
              })}
            </AnimatePresence>
          </div>
        </div>

        {/* Navigation Buttons */}
        <div className="flex justify-center space-x-8">
          <Button
            onClick={prevCard}
            variant="outline"
            size="icon"
            className="rounded-full w-16 h-16 bg-gradient-to-br from-red-500 to-pink-500 hover:from-red-600 hover:to-pink-600 border-none text-white transition-all duration-300 shadow-lg hover:shadow-xl overflow-hidden group"
          >
            <div className="absolute inset-0 bg-gradient-to-t from-transparent to-white opacity-20 group-hover:opacity-30 transition-opacity duration-300"></div>
            <ChevronLeft className="h-8 w-8 relative z-10" />
          </Button>
          <Button
            onClick={nextCard}
            variant="outline"
            size="icon"
            className="rounded-full w-16 h-16 bg-gradient-to-br from-red-500 to-pink-500 hover:from-red-600 hover:to-pink-600 border-none text-white transition-all duration-300 shadow-lg hover:shadow-xl overflow-hidden group"
          >
            <div className="absolute inset-0 bg-gradient-to-t from-transparent to-white opacity-20 group-hover:opacity-30 transition-opacity duration-300"></div>
            <ChevronRight className="h-8 w-8 relative z-10" />
          </Button>
        </div>
      </main>
    </div>
  )
}