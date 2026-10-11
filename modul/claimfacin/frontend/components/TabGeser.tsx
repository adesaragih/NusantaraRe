// Strip tab SEMUA tab Claim Fac In (work owner 11-10-2026 "UBAH TAMPILAN TAB NYA INI JADI LEBIH BAGUS DAN MODERN! DAN BUAT
// PERPINDAHAN TAB DAN HALAMAN NYA ITU SMOOTH" + "INTINYA SEMUA TAB DI MENU KLAIM FACIN"): kelas inti `tabs` / `tabs__item`
// seperti strip tab inti, ditambah kapsul penanda yang MELUNCUR ke tab terpilih. Letak kapsul diukur dari tombol aktif
// (label Pega panjang, lebar bebas, boleh terbungkus ke baris kedua) dan dianimasikan lewat left / top / width - tanpa
// transform (memerangkap Modal tanpa portal).

import { useLayoutEffect, useRef, useState, type KeyboardEvent, type ReactNode } from 'react'

interface Letak {
  left: number
  top: number
  width: number
  height: number
}

const samaLetak = (a: Letak | null, b: Letak | null) =>
  a === b || (!!a && !!b && a.left === b.left && a.top === b.top && a.width === b.width && a.height === b.height)

export function TabGeser<T extends string | number>({
  tab,
  aktif,
  onPilih,
  label,
}: {
  tab: readonly T[]
  aktif: T
  onPilih: (t: T) => void
  label?: (t: T) => ReactNode
}) {
  const jalur = useRef<HTMLDivElement>(null)
  const [letak, setLetak] = useState<Letak | null>(null)
  const i = tab.indexOf(aktif)

  useLayoutEffect(() => {
    const el = jalur.current
    if (!el) return
    const ukur = () => {
      const b = el.querySelector<HTMLElement>('.tabs__item--aktif')
      const baru = b ? { left: b.offsetLeft, top: b.offsetTop, width: b.offsetWidth, height: b.offsetHeight } : null
      setLetak((lama) => (samaLetak(lama, baru) ? lama : baru))
    }
    ukur()
    if (typeof ResizeObserver === 'undefined') return
    // jalur dan setiap tombol: huruf web yang selesai dimuat / lebar layar berubah menggeser tombol tanpa mengubah jalur
    const ro = new ResizeObserver(ukur)
    ro.observe(el)
    for (const b of el.querySelectorAll('.tabs__item')) ro.observe(b)
    return () => ro.disconnect()
  }, [i, tab.length])

  // panah kiri / kanan, Home, End berpindah tab (pola tablist WAI-ARIA)
  const tekan = (e: KeyboardEvent<HTMLDivElement>) => {
    const n = tab.length
    const ke: number | undefined = { ArrowRight: (i + 1) % n, ArrowLeft: (i - 1 + n) % n, Home: 0, End: n - 1 }[e.key]
    if (ke === undefined || n === 0) return
    e.preventDefault()
    onPilih(tab[ke]!)
    jalur.current?.querySelectorAll<HTMLElement>('.tabs__item')[ke]?.focus()
  }

  return (
    <div ref={jalur} className="tabs claimfacin__tabs" role="tablist" onKeyDown={tekan}>
      {letak && <span className="claimfacin__tabs-penanda" aria-hidden="true" style={letak} />}
      {tab.map((t, j) => (
        <button
          key={String(t)}
          type="button"
          role="tab"
          aria-selected={j === i}
          tabIndex={j === i || (i < 0 && j === 0) ? 0 : -1}
          className={'tabs__item' + (j === i ? ' tabs__item--aktif' : '')}
          onClick={() => onPilih(t)}
        >
          {label ? label(t) : t}
        </button>
      ))}
    </div>
  )
}
