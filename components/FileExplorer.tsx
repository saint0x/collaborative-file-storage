'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { FileText, Image, Music, Mic, BookOpen, X } from 'lucide-react'
import { Button } from "@/components/ui/button"

const fileCategories = [
  { name: 'Documents', icon: FileText, color: 'from-gray-200 to-gray-400' },
  { name: 'Images', icon: Image, color: 'from-gray-200 to-gray-400' },
  { name: 'Music', icon: Music, color: 'from-gray-200 to-gray-400' },
  { name: 'Voice Notes', icon: Mic, color: 'from-gray-200 to-gray-400' },
  { name: 'Books', icon: BookOpen, color: 'from-gray-200 to-gray-400' },
]

const files = {
  Documents: [
    { name: 'Project Proposal.docx', date: '2023-05-15' },
    { name: 'Budget Spreadsheet.xlsx', date: '2023-06-05' },
    { name: 'Meeting Minutes.pdf', date: '2023-06-20' },
  ],
  Images: [
    { name: 'Vacation Photo.jpg', date: '2023-06-01' },
    { name: 'Family Portrait.png', date: '2023-05-25' },
    { name: 'Product Mockup.psd', date: '2023-06-15' },
  ],
  Music: [
    { name: 'Summer Playlist.mp3', date: '2023-06-10' },
    { name: 'Workout Mix.mp3', date: '2023-06-12' },
    { name: 'Relaxation Sounds.wav', date: '2023-06-18' },
  ],
  'Voice Notes': [
    { name: 'Meeting Notes.mp3', date: '2023-06-15' },
    { name: 'Lecture Recording.mp3', date: '2023-06-18' },
    { name: 'Idea Brainstorm.m4a', date: '2023-06-22' },
  ],
  Books: [
    { name: '1984.epub', date: '2023-05-20' },
    { name: 'The Great Gatsby.pdf', date: '2023-05-30' },
    { name: 'To Kill a Mockingbird.mobi', date: '2023-06-08' },
  ],
}

interface FileExplorerProps {
  isOpen: boolean
  onClose: () => void
  initialCategory?: string
  initialFile?: string
}

export default function FileExplorer({ isOpen, onClose, initialCategory = 'Documents', initialFile }: FileExplorerProps) {
  const [activeCategory, setActiveCategory] = useState(initialCategory)
  const [expandedPile, setExpandedPile] = useState<string | null>(null)
  const [expandedFile, setExpandedFile] = useState<string | null>(initialFile || null)

  useEffect(() => {
    if (initialFile) {
      setExpandedPile(activeCategory)
    }
  }, [initialFile, activeCategory])

  const handlePileClick = (pileName: string) => {
    if (expandedPile === pileName) {
      setExpandedPile(null)
      setExpandedFile(null)
    } else {
      setExpandedPile(pileName)
      setExpandedFile(null)
    }
  }

  const handleFileClick = (fileName: string) => {
    setExpandedFile(expandedFile === fileName ? null : fileName)
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <motion.div 
        className="w-full max-w-4xl bg-gradient-to-br from-gray-100 to-gray-300 rounded-3xl shadow-xl overflow-hidden flex flex-col"
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.9 }}
        transition={{ duration: 0.3 }}
      >
        <div className="p-8 flex flex-col h-full">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-800">
              Your Files
            </h1>
            <Button variant="ghost" size="icon" onClick={onClose} className="rounded-full bg-gradient-to-r from-gray-200 to-gray-400 text-gray-700 hover:bg-gradient-to-r hover:from-gray-300 hover:to-gray-500">
              <X className="h-6 w-6" />
            </Button>
          </div>
          <div className="flex flex-wrap gap-2 mb-6">
            {fileCategories.map((category) => (
              <Button
                key={category.name}
                onClick={() => {
                  setActiveCategory(category.name)
                  setExpandedPile(null)
                  setExpandedFile(null)
                }}
                variant={activeCategory === category.name ? 'default' : 'outline'}
                className={`rounded-full bg-gradient-to-r ${activeCategory === category.name ? 'from-gray-300 to-gray-500' : 'from-gray-100 to-gray-300'} text-gray-700 hover:bg-gradient-to-r hover:from-gray-200 hover:to-gray-400`}
              >
                <category.icon className="mr-2 h-4 w-4" />
                {category.name}
              </Button>
            ))}
          </div>
          <div className="flex-grow flex items-center justify-center">
            <div className="relative w-72 h-72">
              <AnimatePresence mode="wait">
                <motion.div
                  key={activeCategory}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.3 }}
                  className="absolute inset-0"
                >
                  {files[activeCategory as keyof typeof files].map((file, index) => {
                    const category = fileCategories.find(c => c.name === activeCategory)
                    const isExpanded = expandedPile === activeCategory
                    const isFileExpanded = expandedFile === file.name
                    const Icon = category ? category.icon : FileText
                    return (
                      <motion.div
                        key={file.name}
                        className={`absolute w-48 h-48 rounded-xl shadow-lg bg-gradient-to-br ${category?.color} flex items-center justify-center cursor-pointer overflow-hidden`}
                        initial={false}
                        animate={{
                          rotate: isExpanded ? 0 : Math.min(index * 5 - 5, 10),
                          x: isExpanded ? (isFileExpanded ? index * 260 - 260 : index * 60) : Math.min(index * 10, 20),
                          y: isExpanded ? (isFileExpanded ? 0 : index * 20) : Math.min(index * 10, 20),
                          zIndex: isFileExpanded ? 10 : files[activeCategory as keyof typeof files].length - index,
                          opacity: expandedFile && !isFileExpanded ? 0.3 : 1,
                          scale: isFileExpanded ? 1.1 : 1,
                        }}
                        transition={{ type: 'spring', stiffness: 300, damping: 20 }}
                        onClick={() => isExpanded ? handleFileClick(file.name) : handlePileClick(activeCategory)}
                        whileHover={{ scale: isExpanded ? 1.05 : 1.1 }}
                      >
                        <div className="absolute inset-0 bg-black bg-opacity-10" />
                        <div className="relative text-gray-700 text-center p-4">
                          <Icon className="h-12 w-12 mx-auto mb-2" />
                          <span className="text-sm font-medium break-words">{file.name}</span>
                          {isFileExpanded && (
                            <motion.div
                              initial={{ opacity: 0 }}
                              animate={{ opacity: 1 }}
                              className="mt-2 text-xs"
                            >
                              Date: {file.date}
                            </motion.div>
                          )}
                        </div>
                      </motion.div>
                    )
                  })}
                </motion.div>
              </AnimatePresence>
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  )
}