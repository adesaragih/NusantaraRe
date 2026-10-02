// Isian tanggal modul PremiumList Life — TAMPIL dan DIKETIK dd/mm/yyyy, tetap
// dengan kalender (permintaan work owner 02-10-2026).
//
// ⛔ Kenapa bukan `<input type="date">` saja: bentuk tampil masukan bawaan itu
// ditentukan bahasa/wilayah PERAMBAN, bukan oleh kita — di peramban berbahasa
// Inggris (AS) ia tampil mm/dd/yyyy, dan 03/04/2026 terbaca dua tanggal
// berbeda. Masukan bawaan itu tetap dipakai, tersembunyi, hanya sebagai
// kalender (`showPicker`).
//
// ⛔ Nilai keluar-masuk tetap `YYYY-MM-DD` (atau kosong) — kontrak server tidak
// berubah. Teks yang belum lengkap / bukan tanggal mengirim `''` dan DITANDAI
// di bawah isian, supaya tanggal setengah diketik tidak tersimpan diam-diam
// sebagai tanggal lama.

import { useEffect, useRef, useState } from 'react'

import { isoDariTampil, tanggalTampil, topengTanggal } from '../tanggal'

/** Pesan isian yang bukan tanggal dd/mm/yyyy yang sah. */
export const PESAN_TANGGAL_SALAH = 'Use the dd/mm/yyyy format with a real date.'

export default function IsianTanggal({
  label,
  value,
  onChange,
  readOnly,
  required,
}: {
  label: string
  /** `YYYY-MM-DD` atau kosong. */
  value: string
  /** Menerima `YYYY-MM-DD`, atau `''` bila kosong / belum sah. */
  onChange: (iso: string) => void
  readOnly?: boolean
  required?: boolean
}) {
  const [teks, setTeks] = useState(() => tanggalTampil(value))
  const kalender = useRef<HTMLInputElement>(null)
  // Nilai terakhir yang DIKIRIM dari sini — supaya perubahan dari luar (muat
  // ulang, kalender) menimpa teks, sedangkan ketikan sendiri tidak.
  const terkirim = useRef(value)

  useEffect(() => {
    if (value !== terkirim.current) {
      terkirim.current = value
      setTeks(tanggalTampil(value))
    }
  }, [value])

  const kirim = (iso: string) => {
    terkirim.current = iso
    onChange(iso)
  }

  const salah = teks !== '' && isoDariTampil(teks) === null

  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <div className={readOnly ? 'pl-tanggal pl-tanggal--baca' : 'pl-tanggal'}>
        <input
          className={'field__input' + (salah ? ' field__input--error' : '') + (readOnly ? ' field__input--readonly' : '')}
          type="text"
          inputMode="numeric"
          placeholder="dd/mm/yyyy"
          aria-label={label}
          value={teks}
          readOnly={readOnly}
          onChange={(e) => {
            const t = topengTanggal(e.target.value)
            setTeks(t)
            kirim(isoDariTampil(t) ?? '')
          }}
        />
        {!readOnly && (
          <>
            <button
              type="button"
              className="btn btn--ghost pl-tanggal__tombol"
              aria-label={`Open calendar: ${label}`}
              title="Open calendar"
              onClick={() => {
                kalender.current?.showPicker()
              }}
            >
              <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true" focusable="false">
                <rect x="3" y="5" width="18" height="16" rx="2" fill="none" stroke="currentColor" strokeWidth="2" />
                <path d="M3 10h18M8 3v4M16 3v4" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
              </svg>
            </button>
            <input
              ref={kalender}
              className="pl-tanggal__asli"
              type="date"
              tabIndex={-1}
              aria-hidden="true"
              value={isoDariTampil(teks) ?? ''}
              onChange={(e) => {
                setTeks(tanggalTampil(e.target.value))
                kirim(e.target.value)
              }}
            />
          </>
        )}
      </div>
      {salah && <div className="field__error">{PESAN_TANGGAL_SALAH}</div>}
    </div>
  )
}
