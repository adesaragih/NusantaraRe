// Rute modul Claim Life - refactor bentuk B (30-09-2026): dipindah apa adanya
// dari `App.tsx`, beserta keadaan kasus yang sedang dibuka.

import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { HalamanClaimLife } from './menu'
import InboxClaimLife from './pages/InboxClaimLife'
import KlaimLife from './pages/KlaimLife'
import OutstandingClaimLife from './pages/OutstandingClaimLife'
import RegisterKlaim from './pages/RegisterKlaim'
import './claimlife.css'

export function RuteClaimLife({ halaman, masuk, onPindah, onBukaModul }: PropsRute<HalamanClaimLife>) {
  // Kasus yang sedang dibuka. Kosong berarti belum ada yang dipilih.
  const [kasus, setKasus] = useState('')
  // Case ID yang baru didaftarkan - pemberitahuan sukses di Inbox (10-10-2026). State modul, bukan penyimpanan peramban.
  const [terdaftar, setTerdaftar] = useState<string[]>([])
  // Pita sukses lain di Inbox (mis. sesudah Send to Medical Check, 11-10-2026).
  const [kabar, setKabar] = useState('')

  return (
    // Akar gaya modul: semua aturan `claimlife.css` diawali `.claimlife` (`display: contents`).
    // BahasaUI "en" (10-10-2026, seperti endorsementlife): teks bawaan komponen bersama (paginasi Previous / Page x of N /
    // Next, Memuat, Kosong) berbahasa Inggris, sejajar label korpus layar ini.
    <BahasaUI.Provider value="en">
    <div className="claimlife">
      {halaman === 'inbox' && (
        <InboxClaimLife
          peran={masuk.peran}
          onBukaModul={onBukaModul}
          onBuka={(workID) => {
            setKabar('')
            // Baris Inbox membuka layar kasus yang SAMA untuk setiap tahap (11-10-2026): OutstandingClaimLife menyesuaikan
            // isian dan tombol bawahnya dengan tahap kasus (Outstanding, Medical Check, Claim Analis).
            setKasus(workID)
            onPindah('outstanding')
          }}
          onRegister={() => {
            setTerdaftar([])
            onPindah('register')
          }}
          terdaftar={terdaftar}
          onTutupTerdaftar={() => setTerdaftar([])}
          kabar={kabar}
          onTutupKabar={() => setKabar('')}
        />
      )}
      {halaman === 'outstanding' && (
        <OutstandingClaimLife
          klaimID={kasus}
          peran={masuk.peran}
          onPindah={(k) => {
            setKabar(k ?? '')
            onPindah('inbox')
          }}
        />
      )}
      {halaman === 'register' && (
        <RegisterKlaim
          onTerdaftar={(ids) => {
            // Keputusan work owner 10-10-2026: sesudah Register berhasil langsung ke Inbox, tab paling depan (Inbox
            // dipasang ulang: tab bawaan tabTampil pertama, jumlah dan isi tab dimuat ulang).
            setTerdaftar(ids)
            onPindah('inbox')
          }}
        />
      )}
      {halaman === 'detail' && <KlaimLife />}
    </div>
    </BahasaUI.Provider>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanClaimLife> = RuteClaimLife
