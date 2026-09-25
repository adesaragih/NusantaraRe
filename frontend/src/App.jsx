import KlaimLife from './pages/KlaimLife.jsx'

// Tiket 01: satu halaman yang menampilkan klaim Life dan baris-barisnya.
// Layar modul lain lahir bersama tiketnya sendiri, di src/pages/.
export default function App() {
  return (
    <main>
      <h1>Nusantara Re</h1>
      <KlaimLife />
    </main>
  )
}
