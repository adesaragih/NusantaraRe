import { useState } from 'react'

import { PagarGalat } from './PagarGalat'
import KlaimLife from './pages/KlaimLife'
import Masuk from './pages/Masuk'
import RegisterKlaim from './pages/RegisterKlaim'
import { sesi, type Sesi } from './store/sesi'

// App = kerangka terluar tampilan.
//
// ⚠️ SEMENTARA. F0.3 menggantikan isi cabang "sudah masuk" dengan Shell
// referensi (sidebar, topbar, palet), dan F0.5 memindahkan kedua halaman ini
// ke dalamnya. Yang lahir di F0.2 hanyalah GERBANGnya: tanpa sesi, tidak ada
// halaman yang dirender - sebab setiap permintaannya akan ditolak backend
// tanpa identitas, dan layar penuh galat 401 tidak menjelaskan apa pun.
export default function App() {
  // Sesi dibaca SEKALI saat menyala: memuat ulang (F5) di tab yang sama tetap
  // masuk, sebab `sessionStorage` bertahan selama tab hidup.
  const [masuk, setMasuk] = useState<Sesi | null>(() => sesi.baca())

  if (masuk === null) {
    return <Masuk onMasuk={setMasuk} />
  }

  return (
    // ⛔ PagarGalat membungkus isinya: satu galat render di satu halaman tidak
    // boleh memutihkan seluruh layar tanpa pesan.
    <PagarGalat>
      <main>
        <header className="bilah-sesi">
          <h1>Nusantara Re</h1>
          <span className="bilah-sesi__pelaku">
            {masuk.akunID}
            <code className="bilah-sesi__peran">{masuk.peran.join(', ')}</code>
          </span>
          <button
            type="button"
            className="bilah-sesi__keluar"
            onClick={() => {
              sesi.hapus()
              setMasuk(null)
            }}
          >
            Keluar
          </button>
        </header>
        <RegisterKlaim />
        <KlaimLife />
      </main>
    </PagarGalat>
  )
}
