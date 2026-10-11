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

/**
 * Pengecualian tahap Underwriting (tiket 60): kontrol Accumulation (Choose / Copy Accumulation, popup-nya, Save tab Coverage)
 * bertanda atribut ini tetap aktif bila pembungkus menerima `kecuali={PILIH_UBAH_AKUMULASI}`.
 */
export const PILIH_UBAH_AKUMULASI = '[data-ubah-akumulasi]'

/**
 * Terapkan baca-saja pada satu akar; kembalikan jumlah elemen yang baru dinonaktifkan. `kecuali` = selektor yang dibiarkan.
 * `sembunyi` (mode hanya lihat View Details / View Policy, work owner 09-10-2026 "button ubah,delete,save,submit dan lainnya
 * smua tidak ditampilan"): TOMBOL non-navigasi disembunyikan (`hidden`), bukan sekadar dinonaktifkan.
 */
export function terapkanBacaSaja(akar: ParentNode, kecuali?: string, sembunyi = false): number {
  let n = 0
  akar.querySelectorAll(PEMILIH).forEach((el) => {
    if (navigasi(el)) return
    if (kecuali && el.closest(kecuali) !== null) return
    if (sembunyi && el.tagName === 'BUTTON' && !(el as HTMLElement).hidden) {
      ;(el as HTMLElement).hidden = true
      n++
    }
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

export default function BacaSajaKasus({
  aktif,
  kecuali,
  sembunyi = false,
  children,
}: {
  aktif: boolean
  kecuali?: string
  /** Mode hanya lihat: tombol aksi disembunyikan (lihat `terapkanBacaSaja`). */
  sembunyi?: boolean
  children: ReactNode
}) {
  const ref = useRef<HTMLDivElement>(null)
  useLayoutEffect(() => {
    const akar = ref.current
    if (!akar || !aktif) return
    terapkanBacaSaja(akar, kecuali, sembunyi)
    const pengamat = new MutationObserver(() => terapkanBacaSaja(akar, kecuali, sembunyi))
    pengamat.observe(akar, { subtree: true, childList: true, attributes: true, attributeFilter: ['disabled', 'contenteditable', 'hidden'] })
    return () => pengamat.disconnect()
  }, [aktif, kecuali, sembunyi])
  return (
    <div ref={ref} className={aktif ? 'nbf-bacasaja' : undefined}>
      {children}
    </div>
  )
}
