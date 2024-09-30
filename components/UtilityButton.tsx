"use client"

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Upload, UserPlus, Link, Filter, X } from 'lucide-react'

const UtilityButton = () => {
  const [isExpanded, setIsExpanded] = useState(false)

  const toggleExpand = () => setIsExpanded(!isExpanded)

  const buttonVariants = {
    collapsed: { scale: 1, rotate: 0 },
    expanded: { scale: 1, rotate: 0 },
  }

  const optionVariants = {
    collapsed: {
      x: 0,
      y: 0,
      scale: 0,
      opacity: 0,
      transition: { duration: 0.3 },
    },
    expanded: (index: number) => {
      const totalButtons = 4
      const spacing = 12 // Reduced spacing between buttons
      const buttonSize = 48 // Size of each button
      const radius = (buttonSize + spacing) * (totalButtons - 1) / 2 + buttonSize / 2

      const angle = (Math.PI / 2) * (index / (totalButtons - 1))
      return {
        x: -Math.sin(angle) * radius - buttonSize / 2,
        y: -Math.cos(angle) * radius - buttonSize / 2,
        scale: 1,
        opacity: 1,
        transition: { duration: 0.3, delay: index * 0.05 },
      }
    },
  }

  const options = [
    { icon: <Upload size={20} />, label: 'Upload File' },
    { icon: <Filter size={20} />, label: 'Filter' },
    { icon: <Link size={20} />, label: 'Add Integration' },
    { icon: <UserPlus size={20} />, label: 'Add Connection' },
  ]

  return (
    <div className="fixed bottom-6 right-6 z-50">
      <motion.div
        className="relative"
        animate={isExpanded ? 'expanded' : 'collapsed'}
      >
        <motion.button
          className="w-14 h-14 rounded-full bg-gradient-to-r from-blue-400 to-blue-600 text-white shadow-lg flex items-center justify-center"
          onClick={toggleExpand}
          variants={buttonVariants}
        >
          <X size={24} />
        </motion.button>
        <AnimatePresence>
          {isExpanded && (
            <>
              {options.map((option, index) => (
                <motion.button
                  key={index}
                  className="absolute w-12 h-12 rounded-full bg-white text-blue-600 shadow-md flex items-center justify-center"
                  custom={index}
                  variants={optionVariants}
                  initial="collapsed"
                  animate="expanded"
                  exit="collapsed"
                  whileHover={{ scale: 1.1 }}
                  whileTap={{ scale: 0.9 }}
                >
                  {option.icon}
                  <span className="sr-only">{option.label}</span>
                </motion.button>
              ))}
            </>
          )}
        </AnimatePresence>
      </motion.div>
    </div>
  )
}

export default UtilityButton