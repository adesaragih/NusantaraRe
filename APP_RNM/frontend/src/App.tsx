import KlaimLife from './pages/KlaimLife'

// App = kerangka terluar tampilan. Setiap modul punya halamannya sendiri di
// src/pages/ dan dipasang di sini. Tiket 01 baru punya satu halaman: Klaim Life.
export default function App() {
  return (
    <main>
      <h1>Nusantara Re</h1>
      <KlaimLife />
    </main>
  )
}
