// DAFTAR PILIHAN `<select>` YANG SELALU TERBUKA KE BAWAH — permintaan pemakai
// 8 Oktober 2026 (tangkapan layar dropdown Treaty Group): *"masih keatas drop
// nya"*. ⛔ DITIRU dari `treatyin/frontend/components/DaftarPilihBawah.tsx`,
// bukan diimpor — modul tidak saling impor.
//
// Daftar bawaan browser untuk `<select>` TIDAK dapat diatur arahnya lewat CSS:
// Chrome membukanya ke ATAS bila ruang di bawah kotak kurang dari panjang
// daftarnya. Komponen ini dipasang SEKALI di akar modul (`rute.tsx`) dan
// menggantikan daftar bawaan itu untuk setiap `<select>` di dalam `.treatyinadjustment`:
//
//   1. `mousedown` pada `<select>` dicegah → daftar bawaan tidak muncul;
//   2. bila ruang di bawah sempit, halaman digulir dulu supaya kotaknya ke
//      tengah layar;
//   3. daftar buatan sendiri digambar TEPAT DI BAWAH kotak;
//   4. memilih = mengisi `<select>` ASLI lalu meniupkan peristiwa `change` —
//      `onChange` React milik setiap pemanggil berjalan seperti biasa.
//
// ⭐ `<select>` aslinya TETAP: nilai, `disabled`, papan ketik (panah atas/bawah
// mengganti pilihan), dan semua uji markup tidak berubah. Nol pemanggil yang
// perlu disunting.

import { useEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'

interface Butir {
  nilai: string
  teks: string
  mati: boolean
}

interface Terbuka {
  el: HTMLSelectElement
  butir: Butir[]
  atas: number
  kiri: number
  lebar: number
  tinggiMaks: number
}

/** Ruang minimum di bawah kotak sebelum halaman digulir lebih dulu. */
const RUANG_MIN = 240
const TINGGI_MAKS = 320

function ukur(el: HTMLSelectElement): Omit<Terbuka, 'el' | 'butir'> {
  const r = el.getBoundingClientRect()
  return {
    atas: r.bottom + 2,
    kiri: r.left,
    lebar: r.width,
    tinggiMaks: Math.max(120, Math.min(TINGGI_MAKS, window.innerHeight - r.bottom - 12)),
  }
}

/** Isi `<select>` asli lalu tiupkan `change` — React membaca peristiwa ini. */
function pilihNilai(el: HTMLSelectElement, nilai: string) {
  const setel = Object.getOwnPropertyDescriptor(HTMLSelectElement.prototype, 'value')?.set
  setel?.call(el, nilai)
  el.dispatchEvent(new Event('change', { bubbles: true }))
}

export default function DaftarPilihBawah({ akar }: { akar: string }) {
  const [buka, setBuka] = useState<Terbuka | null>(null)
  const [aktif, setAktif] = useState(-1)
  const daftar = useRef<HTMLUListElement>(null)

  useEffect(() => {
    const tekan = (e: MouseEvent) => {
      const t = e.target
      if (daftar.current?.contains(t as Node)) return
      if (!(t instanceof HTMLSelectElement) || t.multiple || t.disabled || t.closest(akar) === null) {
        setBuka(null)
        return
      }
      if (e.button !== 0) return
      e.preventDefault()
      t.focus()
      if (buka?.el === t) {
        setBuka(null)
        return
      }
      if (window.innerHeight - t.getBoundingClientRect().bottom < RUANG_MIN) {
        t.scrollIntoView({ block: 'center' })
      }
      const butir = [...t.options].map((o) => ({ nilai: o.value, teks: o.text, mati: o.disabled }))
      setAktif(Math.max(0, t.selectedIndex))
      setBuka({ el: t, butir, ...ukur(t) })
    }
    document.addEventListener('mousedown', tekan, true)
    return () => {
      document.removeEventListener('mousedown', tekan, true)
    }
  }, [akar, buka])

  useEffect(() => {
    if (buka === null) return
    const tutup = () => {
      setBuka(null)
    }
    // Halaman digulir (termasuk gulir halus sesudah `scrollIntoView`) — daftar
    // MENGIKUTI kotaknya, tidak ditutup.
    const geser = (e: Event) => {
      if (daftar.current?.contains(e.target as Node)) return
      setBuka((b) => (b === null ? null : { ...b, ...ukur(b.el) }))
    }
    const kunci = (e: KeyboardEvent) => {
      if (e.key === 'Escape' || e.key === 'Tab') {
        setBuka(null)
        return
      }
      if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
        e.preventDefault()
        const arah = e.key === 'ArrowDown' ? 1 : -1
        setAktif((a) => {
          let i = a
          for (let n = 0; n < buka.butir.length; n++) {
            i = (i + arah + buka.butir.length) % buka.butir.length
            if (buka.butir[i]?.mati !== true) return i
          }
          return a
        })
        return
      }
      if (e.key === 'Enter') {
        e.preventDefault()
        const b = buka.butir[aktif]
        if (b !== undefined && !b.mati) pilihNilai(buka.el, b.nilai)
        setBuka(null)
      }
    }
    window.addEventListener('resize', tutup)
    window.addEventListener('scroll', geser, true)
    document.addEventListener('keydown', kunci, true)
    return () => {
      window.removeEventListener('resize', tutup)
      window.removeEventListener('scroll', geser, true)
      document.removeEventListener('keydown', kunci, true)
    }
  }, [buka, aktif])

  // Butir aktif selalu terlihat di dalam daftar yang menggulir.
  useEffect(() => {
    daftar.current?.querySelector('[data-aktif="1"]')?.scrollIntoView({ block: 'nearest' })
  }, [aktif, buka])

  if (buka === null) return null
  return createPortal(
    // Dibungkus akar modul supaya gaya `.treatyinadjustment …` berlaku di `document.body`.
    <div className={akar.replace(/^\./, '')}>
      <ul
        ref={daftar}
        className="tria__pilih-bawah"
        role="listbox"
        style={{ top: buka.atas, left: buka.kiri, width: buka.lebar, maxHeight: buka.tinggiMaks }}
      >
        {buka.butir.map((b, i) => (
          <li
            key={`${String(i)}-${b.nilai}`}
            role="option"
            aria-selected={buka.el.value === b.nilai}
            aria-disabled={b.mati || undefined}
            data-aktif={i === aktif ? '1' : undefined}
            className={
              'tria__pilih-bawah-butir' +
              (i === aktif ? ' tria__pilih-bawah-butir--aktif' : '') +
              (buka.el.value === b.nilai ? ' tria__pilih-bawah-butir--terpilih' : '') +
              (b.mati ? ' tria__pilih-bawah-butir--mati' : '')
            }
            onMouseEnter={() => {
              if (!b.mati) setAktif(i)
            }}
            onMouseDown={(e) => {
              e.preventDefault()
              if (b.mati) return
              pilihNilai(buka.el, b.nilai)
              setBuka(null)
            }}
          >
            {b.teks}
          </li>
        ))}
      </ul>
    </div>,
    document.body,
  )
}
