// Isian angka - SATU komponen untuk medan uang layar, grid spreading / NonProp, dan popup survei (permintaan work
// owner 06-10-2026: "inputannya decimal, selain decimal tidak bisa terketik, format angka langsung kedetek decimal",
// "angka rata kanan", "sangat user friendly"):
//
//   mengetik      hanya angka, satu pemisah desimal, minus di depan (`ketikAngka.ts`); ribuan dipasang otomatis;
//                 titik ATAU koma yang diketik = desimal; teks berformat yang ditempel dikenali
//   tidak fokus   format sel (`sajian.ts`, K14 / bagian uang 4 desimal); kosong atau nol = placeholder "0"
//   tampilan      rata kanan, digit sama lebar (`.nbti__angka`); papan angka di HP (`inputMode="decimal"`)
//
// Nilai yang dikirim ke backend selalu MENTAH (`123456.678`), tidak pernah dibulatkan di sini (AC 24).

import { useLayoutEffect, useRef, useState } from 'react'

import { tampilKetik, ubahKetikan } from '../ketikAngka'
import { nilaiNol, sajikan, type Sajian } from '../sajian'

export default function InputAngka({
  value,
  sajian,
  label,
  onChange,
  onBlur,
  id,
}: {
  value: string
  sajian: Sajian
  label: string
  onChange: (v: string) => void
  onBlur: () => void
  id?: string
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
      className="field__input nbti__angka"
      inputMode="decimal"
      autoComplete="off"
      aria-label={label}
      placeholder="0"
      value={fokus ? tampil : nilaiNol(value) ? '' : sajian === 'tanggal' ? value : sajikan(value, sajian)}
      onFocus={() => {
        setTampil(nilaiNol(value) ? '' : tampilKetik(value))
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
