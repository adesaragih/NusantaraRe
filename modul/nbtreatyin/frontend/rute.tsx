// Rute modul NB Treaty In - portal (daftar + Create) dan layar satu kasus.

import { useEffect, useRef, useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { tujuanKembali, type AsalKasus } from './asalKasus'
import { STATUS_PORTAL, type StatusPortal } from './labels'
import type { HalamanNBTreatyIn } from './menu'
import LayarKasus from './pages/LayarKasus'
import PortalNBTreatyIn from './pages/PortalNBTreatyIn'
import './nbtreatyin.css'

export function RuteNBTreatyIn({ halaman, ketukMenu, bukaKasus, onBeranda }: PropsRute<HalamanNBTreatyIn>) {
  const [kasus, setKasus] = useState('')
  const [pesan, setPesan] = useState('')
  // Tempat berkas yang terbuka dibuka - tombol Back kembali ke sana (permintaan work owner 06-10-2026).
  const [asal, setAsal] = useState<AsalKasus>('portal')
  // Switch portal In Progress / Resolved - bertahan sesudah Back; memilih menu lagi kembali ke bawaan In Progress.
  const [status, setStatus] = useState<StatusPortal>(STATUS_PORTAL[0])

  // Memilih menu NB Treaty In (lagi) kembali ke portal.
  const ketukLalu = useRef(ketukMenu)
  useEffect(() => {
    if (ketukMenu !== undefined && ketukMenu !== ketukLalu.current) {
      setKasus('')
      setAsal('portal')
      setStatus(STATUS_PORTAL[0])
    }
    ketukLalu.current = ketukMenu
  }, [ketukMenu])

  // Berkas yang dibuka langsung dari daftar kotak masuk Beranda (keputusan work owner 06-10-2026).
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
    // Akar gaya modul: semua aturan `nbtreatyin.css` diawali `.nbtreatyin` (`display: contents`).
    <div className="nbtreatyin">
      {halaman === 'nbtreatyin-portal' &&
        (kasus === '' ? (
          <PortalNBTreatyIn
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
export const RUTE_MODUL: RuteModul<HalamanNBTreatyIn> = RuteNBTreatyIn
