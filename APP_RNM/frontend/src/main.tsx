import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import App from './App'

// Titik masuk frontend — berkas pertama yang dijalankan browser.
// index.html punya <div id="root">; React "menempel" ke elemen itu, lalu
// menggambar <App /> di dalamnya. StrictMode hanya menyalakan peringatan
// tambahan saat pengembangan; ia tidak mengubah tampilan.
const root = document.getElementById('root')
if (root === null) {
  throw new Error('index.html tidak punya elemen dengan id "root"')
}

createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
