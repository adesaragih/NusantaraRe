// Peristiwa `change` Pega untuk isian layar Treaty In — 8 Oktober 2026.
//
// ⛔ MENGAPA. Activity/DataTransform Pega berjalan pada `change`: HANYA bila
// nilai isiannya berubah. Bentuk sebelumnya menjalankannya pada SETIAP `blur`
// — isian yang sekadar dilewati ikut menghitung ulang, dan itu menimpa data:
// melewati `Reinstatement` membangun ulang grid Reinstatement (persen yang
// sudah disunting hilang), melewati `Adjustment Rate`/`MDP %` menimpa Premium
// Earned/MDP yang diketik tangan, melewati `Deduction %` menghitung ulang
// `Deduction` dengan pembulatan baru.
//
// ⭐ RUMUSNYA TIDAK BERUBAH — hanya KAPAN ia dijalankan.
//
// Logikanya SAMA dengan `PemicuUbah` layar Adjustment
// (`treatyinadjustment/frontend/komponen/KerangkaTab.tsx`, sudah teruji):
// patokan = nilai saat isian DIMASUKI; saat ditinggalkan, aksi berjalan bila
// nilainya berbeda. Isian yang MENULIS nilainya di blur yang sama (tanggal
// ketik, angka yang dirapikan) diperiksa sekali lagi sesudah render
// berikutnya — hanya sesudah blur itu, supaya nilai yang diubah RUMUS tidak
// memicu `change`. Disalin, bukan diimpor: modul tidak saling impor.

import { useEffect, useRef, type ReactNode } from 'react'

/**
 * `masuk` / `keluar` untuk satu isian bernilai `nilai` (nilai MODEL, bukan
 * teks tampil). `aksi` berjalan HANYA bila nilainya berubah sejak dimasuki.
 */
export function usePemicuUbah(nilai: string, aksi: (() => void) | undefined) {
  const terakhir = useRef(nilai)
  const menyusul = useRef(false)
  const picu = () => {
    if (nilai === terakhir.current) return
    terakhir.current = nilai
    aksi?.()
  }
  useEffect(() => {
    if (!menyusul.current) return
    menyusul.current = false
    picu()
  })
  return {
    masuk: () => {
      terakhir.current = nilai
    },
    keluar: () => {
      picu()
      // Render dari blur ini (event diskret) selesai sebelum tenggat ini.
      menyusul.current = true
      setTimeout(() => {
        menyusul.current = false
      }, 0)
    },
  }
}

/**
 * Pembungkus isian yang menjalankan `aksi` pada peristiwa `change` Pega.
 * `aktif` false (mode lihat, terkunci) = tidak pernah memicu.
 */
export function PemicuUbah({
  nilai,
  aksi,
  aktif = true,
  className,
  children,
}: {
  nilai: string
  aksi: () => void
  aktif?: boolean
  className?: string
  children: ReactNode
}) {
  const p = usePemicuUbah(nilai, aktif ? aksi : undefined)
  return (
    <div className={className} onFocus={p.masuk} onBlur={p.keluar}>
      {children}
    </div>
  )
}
