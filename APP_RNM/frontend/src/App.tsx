import KlaimLife from './pages/KlaimLife'
import RegisterKlaim from './pages/RegisterKlaim'

// App = kerangka terluar tampilan. Setiap modul punya halamannya sendiri di
// src/pages/ dan dipasang di sini. Tiket 01 punya halaman Klaim Life; tiket 02 menambah Register.
export default function App() {
  return (
    <main>
      <h1>Nusantara Re</h1>
      <RegisterKlaim />
      <KlaimLife />
    </main>
  )
}
