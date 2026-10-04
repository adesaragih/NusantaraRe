// Kotak tanggal berformat dd/mm/yyyy - permintaan work owner 03-10-2026 ("Estimated Closing Date : ubah jadi
// format dd/mm/yyyy").
//
// `<input type="date">` bawaan menampilkan tanggal menurut BAHASA BROWSER (di mesin uji: mm/dd/yyyy) dan tidak
// dapat dipaksa ke dd/mm/yyyy. Karena itu yang terlihat dan diketik adalah kotak teks dd/mm/yyyy (garis miring
// disisipkan otomatis), dan tombol kalender membuka pemilih tanggal bawaan lewat `<input type="date">`
// tersembunyi. Nilai keluar = bentuk KABEL inti `DD-MM-YYYY` (`inti/frontend/lib/tanggalInput.ts`), kosong bila
// isian belum lengkap atau bukan tanggal yang ada.

import { useEffect, useRef, useState } from 'react'

import { dariInputTanggal, keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'

/** `DD-MM-YYYY` (kabel) -> `dd/mm/yyyy` (tampil). */
export function keTampil(kabel: string): string {
  return /^\d{2}-\d{2}-\d{4}$/.test(kabel) ? kabel.replace(/-/g, '/') : ''
}

/** Ketikan bebas -> angka saja, garis miring disisipkan sesudah hari dan bulan, paling panjang dd/mm/yyyy. */
export function rapikanKetikan(s: string): string {
  const d = s.replace(/\D/g, '').slice(0, 8)
  if (d.length <= 2) return d
  if (d.length <= 4) return `${d.slice(0, 2)}/${d.slice(2)}`
  return `${d.slice(0, 2)}/${d.slice(2, 4)}/${d.slice(4)}`
}

/** `dd/mm/yyyy` lengkap dan tanggalnya ada -> kabel `DD-MM-YYYY`; selain itu kosong. */
export function keKabel(tampil: string): string {
  if (!/^\d{2}\/\d{2}\/\d{4}$/.test(tampil)) return ''
  const kabel = tampil.replace(/\//g, '-')
  return keInputTanggal(kabel) === '' ? '' : kabel
}

export default function TanggalDMY({
  label,
  value,
  onChange,
  required,
  error,
  labelKalender,
  pesanFormat,
}: {
  label: string
  /** Bentuk kabel `DD-MM-YYYY`, atau kosong. */
  value: string
  onChange: (kabel: string) => void
  required?: boolean
  error?: string
  /** Nama aksesibel tombol kalender. */
  labelKalender: string
  /** Pesan bila isian sudah 10 karakter tetapi bukan tanggal yang ada. */
  pesanFormat: string
}) {
  const [teks, setTeks] = useState(keTampil(value))
  const pemilih = useRef<HTMLInputElement>(null)

  // Nilai dari luar berubah (mis. dipilih dari kalender) -> tampilan mengikuti.
  useEffect(() => {
    // Hanya `value` yang memicu: ketikan setengah jadi (teks belum lengkap, value kosong) tidak boleh
    // dihapus setiap kali pengguna mengetik satu angka.
    if (keKabel(teks) !== value) setTeks(keTampil(value))
  }, [value])

  const ketik = (s: string) => {
    const rapi = rapikanKetikan(s)
    setTeks(rapi)
    onChange(keKabel(rapi))
  }

  const salahFormat = teks.length === 10 && keKabel(teks) === ''
  const galat = error ?? (salahFormat ? pesanFormat : undefined)

  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <div className="nbf-tanggal">
        <input
          className={'field__input' + (galat ? ' field__input--error' : '')}
          type="text"
          inputMode="numeric"
          placeholder="dd/mm/yyyy"
          value={teks}
          onChange={(e) => ketik(e.target.value)}
        />
        <button
          type="button"
          className="btn btn--ghost btn--sm"
          aria-label={labelKalender}
          title={labelKalender}
          onClick={() => {
            const el = pemilih.current
            if (el && typeof el.showPicker === 'function') el.showPicker()
          }}
        >
          📅
        </button>
        <input
          ref={pemilih}
          className="nbf-tanggal__pemilih"
          type="date"
          tabIndex={-1}
          aria-hidden="true"
          value={keInputTanggal(value)}
          onChange={(e) => {
            const kabel = dariInputTanggal(e.target.value)
            setTeks(keTampil(kabel))
            onChange(kabel)
          }}
        />
      </div>
      {galat && <div className="field__error">{galat}</div>}
    </div>
  )
}
