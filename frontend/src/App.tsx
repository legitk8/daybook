import { useEffect, useState } from 'react'
import './App.css'

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

  useEffect(() => {
    fetchPages()
  }, [])

  const fetchPages = () => {
    fetch('http://localhost:8080/pages')
      .then(res => res.json())
      .then(data => setPages(data))
      .catch(err => console.error('Failed to fetch:', err))
  }

  // Create EMPTY page
  const handleAddPage = () => {
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

  const selectPage = (page: Page) => {
    setSelectedPage(page)
    setEditTitle(page.title)
    setEditContent(page.content)
  }

  const handleSave = () => {
    if (!selectedPage) return

    fetch(`http://localhost:8080/pages/${selectedPage.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: editTitle,
        content: editContent,
      }),
    })
      .then(res => res.json())
      .then(updatedPage => {
        // Update sidebar list
        setPages(prev =>
          prev.map(p => (p.id === updatedPage.id ? updatedPage : p))
        )

        // Refresh selected page
        setSelectedPage(updatedPage)
      })
      .catch(err => console.error('Save failed:', err))
  }


  const handleDelete = (id: number) => {
    fetch(`http://localhost:8080/pages/${id}`, {
      method: 'DELETE',
    })
      .then(() => {
        // Remove from UI immediately
        setPages(prev => prev.filter(p => p.id !== id))

        // If deleted page was open → reset view
        if (selectedPage?.id === id) {
          setSelectedPage(null)
        }
      })
      .catch(err => console.error('Delete failed:', err))
  }


  return (
    <div className='app'>
      <div className='sidebar'>
        <div className='sidebar-header'>
          <h2>Pages</h2>
          <button className='add-btn' onClick={handleAddPage}>+</button>
        </div>

        <ul>
          {pages.map(page => (
            <li
              key={page.id}
              className={selectedPage?.id === page.id ? 'active' : ''}
            >
              <span onClick={() => selectPage(page)}>
                {page.title}
              </span>

              <button
                className='delete-btn'
                onClick={(e) => {
                  e.stopPropagation()
                  handleDelete(page.id)
                }}
              >
                ✕
              </button>
            </li>
          ))}
        </ul>

      </div>

      <div className='content'>
        {selectedPage ? (
          <>
            <input
              className='title-input'
              value={editTitle}
              onChange={e => setEditTitle(e.target.value)}
            />

            <textarea
              className='content-input'
              value={editContent}
              onChange={e => setEditContent(e.target.value)}
              rows={10}
            />
            <button className='save-btn' onClick={handleSave}>
              Save
            </button>
          </>
        ) : (
          <>
            <h1>:)</h1>
            <p>Nothing to see here</p>
          </>
        )}
      </div>
    </div>
  )
}

export default App
