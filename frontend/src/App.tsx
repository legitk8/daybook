import { useEffect, useState, useRef } from 'react'

type Page = {
  id: number
  title: string
  content: string
}

function App() {
  const [pages, setPages] = useState<Page[]>([])
  const [selectedPage, setSelectedPage] = useState<Page | null>(null)

  const [editTitle, setEditTitle] = useState('')
  const [editContent, setEditContent] = useState('')

  const [isDirty, setIsDirty] = useState(false)
  const [saveStatus, setSaveStatus] = useState<
    'idle' | 'unsaved' | 'saving' | 'saved'
  >('idle')

  const debounceRef = useRef<number | null>(null)
  const isSavingRef = useRef(false)

  useEffect(() => {
    fetchPages()
  }, [])

  const fetchPages = () => {
    fetch('http://localhost:8080/pages')
      .then(res => res.json())
      .then(data => setPages(data))
      .catch(err => console.error('Failed to fetch:', err))
  }

  // Dirty tracking
  useEffect(() => {
    if (!selectedPage) return

    const changed =
      editTitle !== selectedPage.title ||
      editContent !== selectedPage.content

    if (changed) {
      setIsDirty(true)
      setSaveStatus('unsaved')
    }
  }, [editTitle, editContent])

  // Debounced autosave
  useEffect(() => {
    if (!selectedPage || !isDirty) return

    if (debounceRef.current) {
      clearTimeout(debounceRef.current)
    }

    debounceRef.current = window.setTimeout(() => {
      handleSave()
    }, 800)

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current)
      }
    }
  }, [editTitle, editContent, isDirty, selectedPage])

  // Snapshot-safe save
  const handleSave = async () => {
    if (!selectedPage || isSavingRef.current) return

    // If nothing changed → still show Saved feedback
    if (!isDirty) {
      setSaveStatus('saved')
      return
    }

    const pageId = selectedPage.id
    const titleSnapshot = editTitle
    const contentSnapshot = editContent

    isSavingRef.current = true
    setSaveStatus('saving')

    try {
      const res = await fetch(`http://localhost:8080/pages/${pageId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: titleSnapshot,
          content: contentSnapshot,
        }),
      })

      const updatedPage = await res.json()

      // Ignore stale response
      if (selectedPage?.id !== pageId) return

      setPages(prev =>
        prev.map(p => (p.id === updatedPage.id ? updatedPage : p))
      )

      setSelectedPage(updatedPage)
      setIsDirty(false)
      setSaveStatus('saved')
    } catch (err) {
      console.error('Save failed:', err)

      if (selectedPage?.id === pageId) {
        setSaveStatus('unsaved')
      }
    } finally {
      isSavingRef.current = false
    }
  }

  // Save-before-add
  const handleAddPage = async () => {
    if (isDirty) {
      await handleSave()
    }

    fetch('http://localhost:8080/pages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: 'Untitled',
        content: '',
      }),
    })
      .then(res => res.json())
      .then(newPage => {
        setPages(prev => [...prev, newPage])
        selectPage(newPage)
      })
      .catch(err => console.error('Create failed:', err))
  }

  const handleDelete = async (id: number) => {
    // Cancel pending autosave
    if (debounceRef.current) {
      clearTimeout(debounceRef.current)
    }

    fetch(`http://localhost:8080/pages/${id}`, {
      method: 'DELETE',
    })
      .then(() => {
        setPages(prev => {
          const updated = prev.filter(p => p.id !== id)

          if (selectedPage?.id === id) {
            if (updated.length > 0) {
              const deletedIndex = prev.findIndex(p => p.id === id)
              const nextPage =
                updated[deletedIndex] || updated[deletedIndex - 1]

              if (nextPage) {
                selectPage(nextPage)
              }
            } else {
              setSelectedPage(null)
            }
          }

          return updated
        })
      })
      .catch(err => console.error('Delete failed:', err))
  }

  // Save-before-switch
  const selectPage = async (page: Page) => {
    if (isDirty) {
      await handleSave()
    }

    if (debounceRef.current) {
      clearTimeout(debounceRef.current)
    }

    setSelectedPage(page)
    setEditTitle(page.title)
    setEditContent(page.content)

    setIsDirty(false)
    setSaveStatus('idle')
  }

  return (
    <div className="h-screen flex bg-gradient-to-br from-slate-100 to-slate-200 text-slate-800">

      {/* Sidebar */}
      <div className="w-[20%] backdrop-blur-md bg-white/70 border-r border-white/40 p-5 flex flex-col">

        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold tracking-tight">Pages</h2>

          <button
            onClick={handleAddPage}
            className="w-9 h-9 rounded-xl bg-white shadow-sm hover:shadow-md hover:scale-105 transition"
          >
            +
          </button>
        </div>

        <ul className="space-y-1 overflow-y-auto">
          {pages.map(page => (
            <li
              key={page.id}
              onClick={() => selectPage(page)}
              className={`group flex items-center justify-between px-3 py-2 rounded-xl cursor-pointer transition-all duration-200
                ${selectedPage?.id === page.id
                  ? 'bg-indigo-500 text-white shadow-sm'
                  : 'hover:bg-slate-200/70'
                }`}
            >
              <span className="truncate text-sm font-medium">
                {page.title || 'Untitled'}
              </span>

              <button
                onClick={(e) => {
                  e.stopPropagation()
                  handleDelete(page.id)
                }}
                className="opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-xs text-slate-400 hover:text-red-500"
              >
                ✕
              </button>
            </li>
          ))}
        </ul>
      </div>

      {/* Editor */}
      <div className="w-[80%] p-6 flex justify-center">
        {selectedPage ? (
          <div className="w-full max-w-5xl bg-white rounded-3xl shadow-xl p-8 flex flex-col">

            <input
              value={editTitle}
              onChange={e => setEditTitle(e.target.value)}
              placeholder="Untitled"
              className="w-full text-3xl font-semibold mb-6 outline-none"
            />

            <textarea
              value={editContent}
              onChange={e => setEditContent(e.target.value)}
              placeholder="Start writing..."
              className="w-full flex-1 min-h-[300px] resize-none outline-none text-slate-600 leading-relaxed"
            />

            <button
              onClick={handleSave}
              className="mt-6 px-6 py-2.5 rounded-xl bg-indigo-600 text-white font-medium shadow-sm hover:bg-indigo-500 hover:shadow-md transition-all"
            >
              Save
            </button>

            {/* Save Status Indicator */}
            <div className="mt-2 text-sm text-slate-400">
              {saveStatus === 'unsaved' && 'Unsaved changes'}
              {saveStatus === 'saving' && 'Saving...'}
              {saveStatus === 'saved' && 'Saved'}
            </div>

          </div>
        ) : (
          <div className="text-slate-400 text-sm">
            Nothing to see here ;)
          </div>
        )}
      </div>
    </div>
  )
}

export default App
