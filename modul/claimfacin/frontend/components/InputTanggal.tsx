// Disalin dari `modul/claimnonprop/frontend/components/InputTanggal.tsx` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Isian tanggal Claim Prop (work owner 08-10-2026: "format tanggalnya kenapa susah banget di ketik ya? perbaiki") - pola
// `TanggalDMY` nbfacin, disalin (bukan diimpor, batas modul):
//
//   mengetik   kotak teks; cukup angka, pemisah disisipkan otomatis (`ketikTanggal.ts`), tampil dd-mm-yyyy (+ hh:mm)
//   kalender   tombol membuka pemilih tanggal bawaan lewat isian bawaan tersembunyi (kelas inti `sr-only`)
//   aksi       `onChange(nilai, true)` hanya saat isian lengkap dan sah, atau dikosongkan - sama dengan isian bawaan
//              sebelumnya; ketikan setengah jadi mengosongkan nilai halaman tanpa aksi server
//
// Nilai halaman tetap `YYYY-MM-DD` / `YYYY-MM-DD HH:MM:SS`.

import { useEffect, useRef, useState } from 'react'

import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import { CONTOH, PANJANG, halamanTanggal, rapikanTanggal, tampilTanggal, type JenisTanggal } from '../ketikTanggal'
import { dariInputWaktu, keInputWaktu } from '../nilai'

const awal = (v: string, jenis: JenisTanggal) => rapikanTanggal(tampilTanggal(v, jenis), jenis)

export default function InputTanggal({
  id,
  className,
  jenis,
  value,
  disabled,
  onChange,
}: {
  id?: string
  className: string
  jenis: JenisTanggal
  value: string
  disabled?: boolean
  /** `final` = isian lengkap dan sah, atau dikosongkan: aksi server boleh berjalan. */
  onChange: (nilai: string, final: boolean) => void
}) {
  const [teks, setTeks] = useState(() => awal(value, jenis))
  const [fokus, setFokus] = useState(false)
  const pemilih = useRef<HTMLInputElement>(null)

  // Nilai dari luar berubah (server, kalender) -> tampilan mengikuti; ketikan setengah jadi (nilai halaman kosong)
  // tidak dihapus setiap kali pengguna mengetik satu angka.
  useEffect(() => {
    if (halamanTanggal(teks, jenis) !== value) setTeks(awal(value, jenis))
  }, [value])

  const ketik = (s: string) => {
    const rapi = rapikanTanggal(s, jenis)
    setTeks(rapi)
    const nilai = halamanTanggal(rapi, jenis)
    if (nilai !== '') {
      if (nilai !== value) onChange(nilai, true)
    } else if (value !== '') {
      onChange('', rapi === '')
    }
  }

  const salah = teks !== '' && halamanTanggal(teks, jenis) === '' && (teks.length === PANJANG[jenis] || !fokus)

  // Satu baris: kotak + tombol kalender tidak pernah terpisah, juga di sel grid sempit (Estimation Date; work owner
  // 08-10-2026 "perbaiki tampilan ini"). Tata letaknya di claimfacin.css (`claimfacin__tanggal`).
  return (
    <span className="claimfacin__tanggal">
      <input
        id={id}
        className={className + (salah ? ' field__input--error' : '')}
        type="text"
        inputMode="numeric"
        autoComplete="off"
        placeholder={CONTOH[jenis]}
        title={CONTOH[jenis]}
        value={teks}
        disabled={disabled}
        onFocus={() => setFokus(true)}
        onBlur={() => setFokus(false)}
        onChange={(e) => ketik(e.target.value)}
      />
      {/* pembungkus berposisi: pemilih tersembunyi (sr-only, absolut) berjangkar di tombol, kalender terbuka di sini */}
      <span className="claimfacin__tanggal-kalender">
        <button
          type="button"
          className="btn btn--ghost btn--sm claimfacin__ikon"
          disabled={disabled}
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
          type={jenis === 'tanggal' ? 'date' : 'datetime-local'}
          tabIndex={-1}
          aria-hidden="true"
          value={jenis === 'tanggal' ? keInputTanggal(value) : keInputWaktu(value)}
          onChange={(e) => {
            const nilai = jenis === 'tanggal' ? e.target.value : dariInputWaktu(e.target.value)
            setTeks(awal(nilai, jenis))
            onChange(nilai, true)
          }}
        />
      </span>
    </span>
  )
}
