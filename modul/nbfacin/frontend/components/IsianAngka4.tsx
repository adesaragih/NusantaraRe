// Kotak isian angka yang nilainya juga datang dari backend (sampai 8 / 20 desimal) - TAMPIL 4 desimal (work owner
// 06-10-2026: "ini %installment msh blm 4 digit belakang koma"):
// - tidak difokus: `angka4` ("100,0000");
// - difokus: 4 desimal bertitik ("100.0000") untuk diketik; koma yang diketik dibaca sebagai titik desimal;
// - TIDAK diketik -> nilai lengkap asli tidak berubah (mis. % cicilan 33.33333333333333333333 tetap dikirim utuh).
// Selesai (lepas fokus / Enter) memanggil `onSelesai` = event change / enter Pega.

import { useState } from 'react'

import { angka4 } from '../angka4'

/** Bentuk 4 desimal bertitik untuk diketik: "100,0000" -> "100.0000" (ribuan dibuang). */
export function untukDiketik(nilai: string): string {
  const tampil = angka4(nilai)
  if (tampil === '' || tampil === nilai.trim()) return nilai
  return tampil.replace(/\./g, '').replace(',', '.')
}

/** Ketikan -> teks kabel: koma desimal menjadi titik, spasi dibuang. */
export const dariKetikan = (teks: string) => teks.replace(/\s/g, '').replace(',', '.')

export default function IsianAngka4({
  value,
  onChange,
  onSelesai,
  ariaLabel,
  id,
}: {
  value: string
  onChange: (kabel: string) => void
  onSelesai?: () => void
  ariaLabel?: string
  id?: string
}) {
  const [fokus, setFokus] = useState(false)
  // Teks yang sedang diketik; null = belum diketik sejak fokus (nilai asli dipertahankan).
  const [ketik, setKetik] = useState<string | null>(null)
  const tampil = !fokus ? angka4(value) : ketik ?? untukDiketik(value)
  const selesai = () => {
    setKetik(null)
    onSelesai?.()
  }
  return (
    <input
      id={id}
      className="nbf-angka"
      aria-label={ariaLabel}
      inputMode="decimal"
      value={tampil}
      onFocus={() => {
        setFokus(true)
        setKetik(null)
      }}
      onChange={(e) => {
        setKetik(e.target.value)
        onChange(dariKetikan(e.target.value))
      }}
      onBlur={() => {
        setFokus(false)
        selesai()
      }}
      onKeyDown={(e) => {
        if (e.key === 'Enter') selesai()
      }}
    />
  )
}
