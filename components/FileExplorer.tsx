'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { FileText, Image, Music, Mic, BookOpen } from 'lucide-react'
import { Button } from "@/components/ui/button"

const fileCategories = [
  { name: 'Documents', icon: FileText, color: 'bg-blue-500' },
  { name: 'Images', icon: Image, color: 'bg-green-500' },
  { name: 'Music', icon: Music, color: 'bg-purple-500' },
  { name: 'Voice Notes', icon: Mic, color: 'bg-red-500' },
  { name: 'Books', icon: BookOpen, color: 'bg-yellow-500' },
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

export default function FileExplorer() {
  const [activeCategory, setActiveCategory] = useState('Documents')
  const [expandedPile, setExpandedPile] = useState<string | null>(null)
  const [expandedFile, setExpandedFile] = useState<string | null>(null)

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

  return (
    <motion.div 
      className="w-full max-w-4xl bg-white rounded-3xl shadow-xl overflow-hidden flex flex-col"
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.5 }}
    >
      <div className="p-8 flex flex-col h-full">
        <h1 className="text-3xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500">
          Your Files
        </h1>
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
              className="rounded-full"
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
                      className={`absolute w-48 h-48 rounded-xl shadow-lg ${category?.color} flex items-center justify-center cursor-pointer overflow-hidden`}
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
                      <div className="absolute inset-0 bg-black bg-opacity-20" />
                      <div className="relative text-white text-center p-4">
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
  )
}