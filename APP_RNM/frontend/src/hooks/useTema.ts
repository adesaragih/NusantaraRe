/**
 * Tema terang/gelap untuk seluruh aplikasi — lihat `lib/tema.ts`.
 *
 * Hook ini hanya menyambungkan logika murni itu ke peramban: membaca dan
 * menyimpan pilihan, mengikuti perubahan tema sistem selama pemakai belum
 * memilih, dan memasang `data-theme` pada <html>.
 */

import { useCallback, useEffect, useState } from 'react'

import type { GudangMini } from '../lib/lipatMenu'
import { atributTema, bacaTema, simpanTema, temaBerlaku, type Tema } from '../lib/tema'

const KUERI_GELAP = '(prefers-color-scheme: dark)'

// localStorage disentuh lewat fungsi supaya kegagalannya terkurung.
function gudang(): GudangMini | null {
  try {
    return window.localStorage
  } catch {
    return null
  }
}

function sistemGelapKini(): boolean {
  return typeof window.matchMedia === 'function' && window.matchMedia(KUERI_GELAP).matches
}

export function useTema(): { tema: Tema; balik: () => void } {
  const [tersimpan, setTersimpan] = useState<Tema | null>(() => bacaTema(gudang()))
  const [sistemGelap, setSistemGelap] = useState(sistemGelapKini)

  // Tema sistem diganti saat aplikasi terbuka: ikut, SELAMA pemakai belum
  // memilih sendiri (pilihan tersimpan tetap menang, lihat temaBerlaku).
  useEffect(() => {
    if (typeof window.matchMedia !== 'function') return
    const mq = window.matchMedia(KUERI_GELAP)
    const ubah = (): void => {
      setSistemGelap(mq.matches)
    }
    mq.addEventListener('change', ubah)
    return () => {
      mq.removeEventListener('change', ubah)
    }
  }, [])

  const tema = temaBerlaku(tersimpan, sistemGelap)

  useEffect(() => {
    document.documentElement.dataset.theme = atributTema(tema)
  }, [tema])

  const balik = useCallback(() => {
    const baru: Tema = tema === 'gelap' ? 'terang' : 'gelap'
    simpanTema(gudang(), baru)
    setTersimpan(baru)
  }, [tema])

  return { tema, balik }
}
