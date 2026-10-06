// Rute modul NB Treaty In - portal (daftar + Create) dan layar satu kasus.

import { useEffect, useRef, useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanNBTreatyIn } from './menu'
import LayarKasus from './pages/LayarKasus'
import PortalNBTreatyIn from './pages/PortalNBTreatyIn'
import './nbtreatyin.css'

export function RuteNBTreatyIn({ halaman, ketukMenu }: PropsRute<HalamanNBTreatyIn>) {
  const [kasus, setKasus] = useState('')
  const [pesan, setPesan] = useState('')

  // Memilih menu NB Treaty In (lagi) kembali ke portal.
  const ketukLalu = useRef(ketukMenu)
  useEffect(() => {
    if (ketukMenu !== undefined && ketukMenu !== ketukLalu.current) setKasus('')
    ketukLalu.current = ketukMenu
  }, [ketukMenu])

  return (
    // Akar gaya modul: semua aturan `nbtreatyin.css` diawali `.nbtreatyin` (`display: contents`).
    <div className="nbtreatyin">
      {halaman === 'nbtreatyin-portal' &&
        (kasus === '' ? (
          <PortalNBTreatyIn
            pesan={pesan}
            onBuka={(id) => {
              setPesan('')
              setKasus(id)
            }}
          />
        ) : (
          <LayarKasus
            id={kasus}
            onKembali={(p) => {
              setPesan(p ?? '')
              setKasus('')
            }}
          />
        ))}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanNBTreatyIn> = RuteNBTreatyIn
