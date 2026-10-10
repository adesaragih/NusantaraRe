// Disalin dari `modul/claimprop/frontend/components/InputAngka.tsx` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Isian angka Claim Prop (work owner 08-10-2026: "perbaiki semua inputan format number, hanya bisa di isi format
// number, separator indonesia, 4 angka belakang koma") - pola isian angka nbtreatyin / edmtreatyin, disalin (bukan
// diimpor, batas modul):
//
//   mengetik      hanya angka, satu pemisah desimal (titik ATAU koma yang diketik), minus di depan, ribuan otomatis,
//                 paling banyak 4 desimal (`ketikAngka.ts`); teks berformat yang ditempel dikenali
//   tidak fokus   `tampilAngka` - titik ribuan, koma desimal, maks 4 desimal, nol ekor dibuang
//
// Nilai yang dikirim ke backend selalu MENTAH (`123456.678`).

import { useLayoutEffect, useRef, useState } from 'react'

import { tampilKetik, ubahKetikan } from '../ketikAngka'
import { tampilAngka } from '../nilai'

export default function InputAngka({
  id,
  className,
  value,
  disabled,
  onChange,
  onBlur,
}: {
  id?: string
  className: string
  value: string
  disabled?: boolean
  onChange: (mentah: string) => void
  onBlur: () => void
}) {
  const ref = useRef<HTMLInputElement>(null)
  const kursor = useRef<number | null>(null)
  const [fokus, setFokus] = useState(false)
  const [tampil, setTampil] = useState('')

  // kursor dikembalikan ke digit yang sama sesudah titik ribuan ditata ulang
  useLayoutEffect(() => {
    if (kursor.current !== null && ref.current !== null) {
      ref.current.setSelectionRange(kursor.current, kursor.current)
      kursor.current = null
    }
  })

  return (
    <input
      ref={ref}
      id={id}
      className={className}
      inputMode="decimal"
      autoComplete="off"
      disabled={disabled}
      value={fokus ? tampil : tampilAngka(value)}
      onFocus={() => {
        setTampil(tampilKetik(value))
        setFokus(true)
      }}
      onKeyDown={(e) => {
        // selain angka, pemisah, dan minus tidak terketik (pintasan Ctrl/Cmd tetap jalan)
        if (e.ctrlKey || e.metaKey || e.altKey) return
        if (e.key.length === 1 && !/[\d.,-]/.test(e.key)) e.preventDefault()
      }}
      onChange={(e) => {
        const r = ubahKetikan(tampil, e.target.value, e.target.selectionStart ?? e.target.value.length)
        setTampil(r.tampil)
        kursor.current = r.kursor
        const mentah = r.mentah === '-' ? '' : r.mentah
        if (mentah !== value) onChange(mentah)
      }}
      onBlur={() => {
        setFokus(false)
        onBlur()
      }}
    />
  )
}
