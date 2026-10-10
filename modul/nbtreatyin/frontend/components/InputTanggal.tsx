// Isian tanggal NB Treaty In (work owner 10-10-2026: "supaya bisa di copy paste dan di ketik lancar") - pola
// `InputTanggal` Claim Prop, disalin (bukan diimpor, batas modul):
//
//   mengetik   kotak teks; cukup angka, pemisah disisipkan otomatis (`ketikTanggal.ts`), tampil dd-mm-yyyy
//   menempel   teks tempelan MENGGANTI isi kotak dan dibaca utuh (dd/mm/yyyy, yyyy-mm-dd, 1 Jul 2024, ...)
//   menyalin   isi kotak teks biasa - tersalin apa adanya (dd-mm-yyyy)
//   kalender   tombol membuka pemilih tanggal bawaan lewat isian bawaan tersembunyi (kelas inti `sr-only`)
//
// `onChange` hanya dengan tanggal lengkap dan sah, atau kosong (dikosongkan / ketikan setengah jadi) - sama dengan
// isian bawaan sebelumnya. Nilai halaman tetap `YYYY-MM-DD`.

import { useEffect, useRef, useState } from 'react'

import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import { CONTOH_TANGGAL, PANJANG_TANGGAL, halamanTanggal, rapikanTanggal, tampilTanggal } from '../ketikTanggal'

const awal = (v: string) => rapikanTanggal(tampilTanggal(v))

export default function InputTanggal({
  id,
  value,
  galat,
  onChange,
}: {
  id: string
  value: string
  /** Pesan galat medan dari server - kotak bertanda merah. */
  galat?: boolean
  onChange: (nilai: string) => void
}) {
  const [teks, setTeks] = useState(() => awal(value))
  const [fokus, setFokus] = useState(false)
  const pemilih = useRef<HTMLInputElement>(null)

  // Nilai dari luar berubah (server, kalender) -> tampilan mengikuti; ketikan setengah jadi (nilai halaman kosong)
  // tidak dihapus setiap kali pengguna mengetik satu angka.
  useEffect(() => {
    if (halamanTanggal(teks) !== keInputTanggal(value)) setTeks(awal(value))
  }, [value])

  const ketik = (s: string) => {
    const rapi = rapikanTanggal(s)
    setTeks(rapi)
    const nilai = halamanTanggal(rapi)
    if (nilai !== keInputTanggal(value)) onChange(nilai)
  }

  const salah = teks !== '' && halamanTanggal(teks) === '' && (teks.length === PANJANG_TANGGAL || !fokus)

  return (
    <span className="nbti__tanggal">
      <input
        id={id}
        className={'field__input' + (galat || salah ? ' field__input--error' : '')}
        type="text"
        inputMode="numeric"
        autoComplete="off"
        placeholder={CONTOH_TANGGAL}
        title={CONTOH_TANGGAL}
        value={teks}
        onFocus={() => setFokus(true)}
        onBlur={() => setFokus(false)}
        onChange={(e) => ketik(e.target.value)}
        onPaste={(e) => {
          // Tempelan menggantikan seluruh isi (bukan disisipkan di posisi kursor) supaya dibaca sebagai satu tanggal.
          e.preventDefault()
          ketik(e.clipboardData.getData('text'))
        }}
      />
      {/* pembungkus berposisi: pemilih tersembunyi (sr-only, absolut) berjangkar di tombol, kalender terbuka di sini */}
      <span className="nbti__tanggal-kalender">
        <button
          type="button"
          className="btn btn--ghost btn--sm"
          aria-label="Calendar"
          title="Calendar"
          onClick={() => {
            const el = pemilih.current
            if (el && typeof el.showPicker === 'function') el.showPicker()
          }}
        >
          <svg viewBox="0 0 24 24" width="15" height="15" aria-hidden="true" focusable="false">
            <path
              d="M4 6h16v14H4zM4 10h16M8 3v4M16 3v4"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </button>
        <input
          ref={pemilih}
          className="sr-only"
          type="date"
          tabIndex={-1}
          aria-hidden="true"
          value={keInputTanggal(value)}
          onChange={(e) => {
            setTeks(awal(e.target.value))
            onChange(e.target.value)
          }}
        />
      </span>
    </span>
  )
}
