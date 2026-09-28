import { useState } from 'react'

import Beranda from './pages/Beranda'
import InboxPremiumList from './pages/premiumlist/InboxPremiumList'
import InputOffer from './pages/premiumlist/InputOffer'
import InboxClaimLife from './pages/InboxClaimLife'
import KlaimLife from './pages/KlaimLife'
import OutstandingClaimLife from './pages/OutstandingClaimLife'
import RegisterKlaim from './pages/RegisterKlaim'
import { Shell, type Halaman } from './components/Shell'
import { BelumTersedia } from './components/ui/dasar'
import { pelakuStub } from './store/sesi'

// App = identitas + Shell.
//
// ⛔ TIDAK ADA FORM LOGIN `[perintah work owner 27-09-2026]`. Aplikasi membuka
// Inbox langsung; identitasnya datang dari env saat menyala (`store/sesi`).
//
// ⛔ Bila stub-nya mati, layar MENYATAKANNYA - bukan pecah, dan bukan pula
// diam-diam menampilkan daftar kosong yang terbaca "tidak ada pekerjaan".
export default function App() {
  const masuk = pelakuStub()
  // ⛔ Beranda layar AWAL sejak butir bg: ia yang menyebut modul mana
  // yang sudah ada dan mana yang belum. Membuka langsung ke Inbox membuat
  // aplikasi tampak hanya punya satu modul.
  const [halaman, setHalaman] = useState<Halaman>('beranda')
  // Kasus yang sedang dibuka. Kosong berarti belum ada yang dipilih.
  const [kasus, setKasus] = useState('')
  // Polis yang sedang dibuka, beserta tahapnya - tiket 01 PremiumList.
  const [polis, setPolis] = useState({ id: '', tahap: '' })

  if (masuk === null) {
    return (
      // ⚠️ `BelumTersedia` hanya menerima `apa`: parameter `sebab`-nya
      // sengaja dibuang di REFERENSI_UI ronde 70 karena diterima lalu
      // diabaikan. Keterangan teknisnya karena itu berdiri di sini.
      <section>
        <BelumTersedia apa="Identitas pelaku" />
        <p className="polis__catatan" role="note">
          Selama IAM belum terpasang, identitas dibentuk dari env Vite dan
          menuntut <code>VITE_AUTH_STUB=true</code>. Tanpa itu setiap
          permintaan ditolak backend, dan layar ini tidak dapat
          menampilkan pekerjaan siapa pun.
        </p>
      </section>
    )
  }

  return (
    <Shell masuk={masuk} halaman={halaman} onPindah={setHalaman}>
      {halaman === 'beranda' && <Beranda masuk={masuk} onBuka={setHalaman} />}
      {halaman === 'premiumlist' && polis.id === '' && (
        <InboxPremiumList
          onBuka={(caseID, tahap) => {
            setPolis({ id: caseID, tahap })
          }}
        />
      )}
      {halaman === 'premiumlist' && polis.id !== '' && (
        <InputOffer
          polisID={polis.id}
          tahap={polis.tahap}
          onSelesai={() => {
            setPolis({ id: '', tahap: '' })
          }}
        />
      )}
      {halaman === 'inbox' && (
        <InboxClaimLife
          peran={masuk.peran}
          onBuka={(workID) => {
            // ⚠️ Baris Inbox membuka layar TAHAPnya. Tab Outstanding
            // membuka `OSClaimLife`; tahap lain menyusul bersama
            // kelompok A3 masing-masing.
            setKasus(workID)
            setHalaman('outstanding')
          }}
          onRegister={() => {
            setHalaman('register')
          }}
        />
      )}
      {halaman === 'outstanding' && (
        <OutstandingClaimLife
          klaimID={kasus}
          onPindah={() => {
            setHalaman('inbox')
          }}
          onDetail={() => {
            setHalaman('detail')
          }}
        />
      )}
      {halaman === 'register' && <RegisterKlaim />}
      {halaman === 'detail' && <KlaimLife />}
    </Shell>
  )
}
