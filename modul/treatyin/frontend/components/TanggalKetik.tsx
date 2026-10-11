// Kotak TANGGAL YANG BISA DIKETIK — permintaan pemakai 7 Oktober 2026:
// *"saya mau tiap inputan tanggal bisa diketik juga biar gampang"*.
//
// ⭐ BENTUKNYA BENTUK PEGA: kotak teks `DD/MM/YYYY` + ikon kalender di
// sampingnya (gambar Pega 38: `18/01/2024 📅`). `<input type="date">`
// bawaan peramban hanya dapat diisi per segmen atau lewat kalender.
//
// Yang diterima saat diketik (lalu dirapikan menjadi `DD/MM/YYYY` saat
// kotak ditinggalkan atau Enter):
//
//   18/01/2024 · 18-01-2024 · 18.01.2024 · 18 01 2024 · 1/8/2024
//   18012024 (DDMMYYYY) · 180124 (DDMMYY)
//   20240118 (YYYYMMDD, bentuk simpan) · 2024-01-18 (ISO)
//   tahun dua digit → 20YY — hasilnya langsung terlihat dirapikan.
//
// Teks yang bukan tanggal TIDAK diteruskan: kotaknya ditandai salah dan
// nilai sebelumnya tetap. Kosong = tanggal dikosongkan.
//
// ⛔ `FieldTanggal` inti (`dasar.tsx`) TIDAK disunting — berkas bersama.
// Konversi bentuk simpan/kabel memakai `tanggalInput.ts` inti yang sama,
// jadi keluarannya tetap bentuk KABEL `DD-MM-YYYY`, persis `FieldTanggal`:
// pemanggil lama tidak berubah.

import { useRef, useState } from 'react'

import { dariInputTanggal, keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'

/** Petunjuk ketik — di `title` dan pesan salah. */
export const PETUNJUK_TANGGAL = 'Type DD/MM/YYYY or pick from the calendar'
export const TANGGAL_SALAH = 'Date not recognised — type DD/MM/YYYY'
const LABEL_KALENDER = 'Pick from the calendar'

function wajar(th: number, bl: number, hr: number): boolean {
  if (th < 1900 || th > 9999 || bl < 1 || bl > 12 || hr < 1) return false
  const kabisat = (th % 4 === 0 && th % 100 !== 0) || th % 400 === 0
  const panjang = [31, kabisat ? 29 : 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]
  return hr <= (panjang[bl - 1] ?? 0)
}

const p2 = (n: number) => String(n).padStart(2, '0')
const tahunPenuh = (y: string) => (y.length === 2 ? 2000 + Number(y) : Number(y))

/**
 * Teks yang DIKETIK → bentuk kabel `DD-MM-YYYY`; `''` bila kosong; `null`
 * bila bukan tanggal. Manipulasi teks murni — nol `Date`, nol zona waktu
 * (alasan yang sama dengan `tanggalInput.ts` inti).
 */
export function uraiTanggalKetik(teks: string): string | null {
  const s = teks.trim()
  if (s === '') return ''
  const coba = (hr: number, bl: number, th: number) => (wajar(th, bl, hr) ? `${p2(hr)}-${p2(bl)}-${String(th)}` : null)
  let m = /^(\d{1,2})[/.\-\s](\d{1,2})[/.\-\s](\d{4}|\d{2})$/.exec(s)
  if (m) return coba(Number(m[1]), Number(m[2]), tahunPenuh(m[3] ?? ''))
  m = /^(\d{4})-(\d{1,2})-(\d{1,2})$/.exec(s)
  if (m) return coba(Number(m[3]), Number(m[2]), Number(m[1]))
  if (/^\d{8}$/.test(s)) {
    // DDMMYYYY lebih dulu — cara orang mengetik; bila bukan tanggal, bentuk
    // simpan Pega YYYYMMDD.
    return (
      coba(Number(s.slice(0, 2)), Number(s.slice(2, 4)), Number(s.slice(4))) ??
      coba(Number(s.slice(6)), Number(s.slice(4, 6)), Number(s.slice(0, 4)))
    )
  }
  if (/^\d{6}$/.test(s)) return coba(Number(s.slice(0, 2)), Number(s.slice(2, 4)), tahunPenuh(s.slice(4)))
  return null
}

/** Nilai apa pun (kabel, `YYYYMMDD`, ISO, stempel Pega) → `DD/MM/YYYY` untuk kotak. */
export function keTeksTanggal(nilai: string): string {
  const iso = keInputTanggal(nilai)
  return iso === '' ? '' : `${iso.slice(8, 10)}/${iso.slice(5, 7)}/${iso.slice(0, 4)}`
}

/** Kotak tanggal-ketik TANPA label (sel grid, atau di dalam `TanggalRedup`). */
export function KotakTanggalKetik({
  label,
  value,
  onChange,
  id,
}: {
  /** Nama untuk pembaca layar. */
  label: string
  /** Bentuk apa pun yang `keInputTanggal` inti baca. */
  value: string
  /** Menerima bentuk KABEL `DD-MM-YYYY`; kosong = dikosongkan. */
  onChange: (kabel: string) => void
  id?: string
}) {
  // `null` = tidak sedang diketik → tampilkan nilai yang dirapikan.
  const [draf, setDraf] = useState<string | null>(null)
  const [salah, setSalah] = useState(false)
  const kalender = useRef<HTMLInputElement>(null)
  const kabelKini = dariInputTanggal(keInputTanggal(value))
  const simpan = () => {
    if (draf === null) return
    const k = uraiTanggalKetik(draf)
    if (k === null) {
      setSalah(true)
      return
    }
    setSalah(false)
    setDraf(null)
    if (k !== kabelKini) onChange(k)
  }
  return (
    <span className="trin__tgl">
      <span className="trin__tgl-baris">
        <input
          className={'field__input' + (salah ? ' field__input--error' : '')}
          type="text"
          id={id}
          autoComplete="off"
          aria-label={label}
          aria-invalid={salah || undefined}
          title={PETUNJUK_TANGGAL}
          value={draf ?? keTeksTanggal(value)}
          onFocus={() => {
            if (draf === null) setDraf(keTeksTanggal(value))
          }}
          onChange={(e) => {
            // ⭐ Huruf TIDAK dapat diketik — hanya digit dan pemisah tanggal
            // `/ - .` serta spasi, maksimal 10 aksara (permintaan pemakai
            // 9 Oktober 2026). Salinan `saringTanggal` modul Adjustment.
            setDraf(e.target.value.replace(/[^\d/.\- ]/g, '').slice(0, 10))
            setSalah(false)
          }}
          onBlur={simpan}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault()
              simpan()
            }
          }}
        />
        <button
          type="button"
          className="trin__tgl-kalender"
          aria-label={`${LABEL_KALENDER} — ${label}`}
          title={LABEL_KALENDER}
          onClick={() => {
            const el = kalender.current
            if (el === null) return
            if (typeof el.showPicker === 'function') el.showPicker()
            else el.focus()
          }}
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
            <rect x="3" y="5" width="18" height="16" rx="2" />
            <path d="M3 10h18M8 3v4M16 3v4" />
          </svg>
        </button>
        {/* Kalender bawaan peramban, tak terlihat — hanya pemilihnya yang dibuka. */}
        <input
          ref={kalender}
          className="trin__tgl-asli"
          type="date"
          tabIndex={-1}
          aria-hidden="true"
          value={keInputTanggal(value)}
          onChange={(e) => {
            setDraf(null)
            setSalah(false)
            onChange(dariInputTanggal(e.target.value))
          }}
        />
      </span>
      {salah && (
        <span className="field__error" role="alert">
          {TANGGAL_SALAH}
        </span>
      )}
    </span>
  )
}
