// Rute modul EDM Treaty In - portal (`SFAPortal_Endorsement_Treaty` + popup Create) dan layar satu kasus. Asal pola:
// `modul/nbtreatyin/frontend/rute.tsx` (06-10-2026: kembali ke portal saat menu dipilih lagi, buka berkas dari
// Beranda, tombol Back ke tempat asal). Tanpa switch status (portal EDM tanpa tab).

import { useEffect, useRef, useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { tujuanKembali, type AsalKasus } from './asalKasus'
import { STATUS_PORTAL, type StatusPortal } from './labels'
import type { HalamanEDMTreatyIn } from './menu'
import LayarKasus from './pages/LayarKasus'
import PortalEDMTreatyIn from './pages/PortalEDMTreatyIn'
import './edmtreatyin.css'

export function RuteEDMTreatyIn({ halaman, ketukMenu, bukaKasus, onBeranda }: PropsRute<HalamanEDMTreatyIn>) {
  const [kasus, setKasus] = useState('')
  const [pesan, setPesan] = useState('')
  // Tempat berkas yang terbuka dibuka - tombol Back kembali ke sana (pola NB).
  const [asal, setAsal] = useState<AsalKasus>('portal')
  // Switch portal In Progress / Resolved (aturan portal NB, WO 07-10-2026): bertahan sesudah Back, klik menu = In
  // Progress.
  const [status, setStatus] = useState<StatusPortal>(STATUS_PORTAL[0])

  // Memilih menu EDM Treaty In (lagi) kembali ke portal.
  const ketukLalu = useRef(ketukMenu)
  useEffect(() => {
    if (ketukMenu !== undefined && ketukMenu !== ketukLalu.current) {
      setKasus('')
      setAsal('portal')
      setStatus(STATUS_PORTAL[0])
    }
    ketukLalu.current = ketukMenu
  }, [ketukMenu])

  // Berkas yang dibuka langsung dari daftar kotak masuk Beranda.
  const ketukBuka = bukaKasus?.ketuk
  const idBuka = bukaKasus?.id
  useEffect(() => {
    if (ketukBuka !== undefined && idBuka !== undefined) {
      setPesan('')
      setKasus(idBuka)
      setAsal('beranda')
    }
  }, [ketukBuka, idBuka])

  return (
    // Akar gaya modul: semua aturan `edmtreatyin.css` diawali `.edmtreatyin` (`display: contents`).
    <div className="edmtreatyin">
      {halaman === 'edmtreatyin-portal' &&
        (kasus === '' ? (
          <PortalEDMTreatyIn
            pesan={pesan}
            status={status}
            onStatus={setStatus}
            onBuka={(id) => {
              setPesan('')
              setKasus(id)
              setAsal('portal')
            }}
          />
        ) : (
          <LayarKasus
            id={kasus}
            onKembali={(p) => {
              const keBeranda = tujuanKembali(asal, p, onBeranda !== undefined) === 'beranda'
              setPesan(p ?? '')
              setKasus('')
              setAsal('portal')
              if (keBeranda) onBeranda?.()
            }}
          />
        ))}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanEDMTreatyIn> = RuteEDMTreatyIn
