// Pembungkus baca-saja layar case (work owner 07-10-2026: case yang sudah keluar dari tahap marketing dibuat baca-saja
// untuk marketing - "setuju, buatkan jangan halu"). Pega: sesudah Submit marketing tidak lagi memegang assignment-nya.
//
// Saat `aktif`, SEMUA isian (input / select / textarea / contenteditable) dan tombol di dalamnya dinonaktifkan, KECUALI
// navigasi yang dibutuhkan untuk MELIHAT data: tombol buka-baris (`aria-expanded`), strip tab (`role="tab"`), paginasi
// (`.pager`), dan elemen bertanda `data-navigasi`. MutationObserver menerapkan ulang aturan pada isian yang baru dirender
// (baris dibuka, sub-tab berganti) dan bila React memulihkan `disabled`. Ini lapis tampilan; backend menolak tulisan
// (409) secara mandiri. Keputusan agent K57-21: satu pembungkus generik, bukan prop baca-saja di tiap komponen tab.

import { useLayoutEffect, useRef, type ReactNode } from 'react'

/** Elemen navigasi yang tetap aktif saat baca-saja. */
export function navigasi(el: Pick<Element, 'closest' | 'getAttribute' | 'hasAttribute'>): boolean {
  return (
    el.hasAttribute('aria-expanded') ||
    el.getAttribute('role') === 'tab' ||
    el.closest('[data-navigasi]') !== null ||
    el.closest('.pager') !== null
  )
}

const PEMILIH = 'input, select, textarea, button, [contenteditable]'

/** Terapkan baca-saja pada satu akar; kembalikan jumlah elemen yang baru dinonaktifkan. */
export function terapkanBacaSaja(akar: ParentNode): number {
  let n = 0
  akar.querySelectorAll(PEMILIH).forEach((el) => {
    if (navigasi(el)) return
    if (el.hasAttribute('contenteditable')) {
      if (el.getAttribute('contenteditable') !== 'false') {
        el.setAttribute('contenteditable', 'false')
        n++
      }
      return
    }
    const kontrol = el as HTMLInputElement
    if (!kontrol.disabled) {
      kontrol.disabled = true
      n++
    }
  })
  return n
}

export default function BacaSajaKasus({ aktif, children }: { aktif: boolean; children: ReactNode }) {
  const ref = useRef<HTMLDivElement>(null)
  useLayoutEffect(() => {
    const akar = ref.current
    if (!akar || !aktif) return
    terapkanBacaSaja(akar)
    const pengamat = new MutationObserver(() => terapkanBacaSaja(akar))
    pengamat.observe(akar, { subtree: true, childList: true, attributes: true, attributeFilter: ['disabled', 'contenteditable'] })
    return () => pengamat.disconnect()
  }, [aktif])
  return (
    <div ref={ref} className={aktif ? 'nbf-bacasaja' : undefined}>
      {children}
    </div>
  )
}
