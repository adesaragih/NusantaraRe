// Unsur bersama layar Adjustment: pemformat sel dan medan baca-saja.
//
// ⛔ SATU pemformat untuk kedua panel. Panel Old dan New menampilkan angka
// yang dibandingkan mata ke mata; dua jalur pemformatan berarti beda
// pembulatan yang terbaca sebagai beda data.

import { formatDate, formatNumber } from '../../../../inti/frontend/lib/format'
import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import { DESIMAL_PERSEN, DESIMAL_PERSEN_SHARE, DESIMAL_UANG, PENYESUAIAN, type JenisAngka } from '../labelsPenyesuaian'

/** Desimal CADANGAN per golongan — hanya bila ekspor tidak menyatakannya. */
const CADANGAN: Readonly<Record<'uang' | 'persen' | 'persenShare', number>> = {
  uang: DESIMAL_UANG,
  persen: DESIMAL_PERSEN,
  persenShare: DESIMAL_PERSEN_SHARE,
}

/**
 * Padankan pecahan hasil `formatNumber` ke presisi kolom.
 *
 * ⛔ Ini MEMADANKAN, bukan memformat ulang: `formatNumber` (inti, tidak
 * disunting) sudah membulatkan ke `desimal` digit lalu membuang nol di ekor;
 * di sini nol itu dikembalikan sampai presisi yang ekspor nyatakan, supaya
 * `1,00` tetap `1,00`. Teks yang bukan bentuk angka tidak disentuh.
 */
export function padankan(teks: string, desimal: number): string {
  if (desimal <= 0 || !/^-?\d{1,3}(?:\.\d{3})*(?:,\d*)?$/.test(teks)) return teks
  const [bulat, pecahan = ''] = teks.split(',')
  return `${bulat ?? ''},${pecahan.padEnd(desimal, '0')}`
}

/** Hanya memutuskan APAKAH pemformat angka dipanggil — bukan pemformat kedua. */
export function angkaMurni(nilai: string): boolean {
  return /^[+-]?\d*(?:\.\d*)?$/.test(nilai.trim()) && /\d/.test(nilai)
}

/**
 * Satu nilai, diformat menurut GOLONGANNYA dan DESIMAL kolomnya.
 *
 * `desimal` = `pyDecimalPlaces` ekspor; `null`/tak diberikan = tidak
 * dinyatakan → desimal cadangan golongan, nol di ekor dibuang (aturan lama).
 *
 * ⛔ Nol `%` ditempel: di sistem lama tanda persen berdiri di LABEL kolom
 * (`Proportion %`, `% Share`) atau sel terpisah (`%` @337597), bukan di
 * nilainya. Teks bukan angka dikembalikan APA ADANYA.
 */
export function selNilai(jenis: JenisAngka, nilai: string, desimal: number | null = null): string {
  switch (jenis) {
    case 'uang':
    case 'persen':
    case 'persenShare': {
      if (!angkaMurni(nilai)) return nilai
      const t = formatNumber(nilai, desimal ?? CADANGAN[jenis])
      return desimal === null ? t : padankan(t, desimal)
    }
    case 'tanggal': {
      // `formatDate` mengembalikan "" untuk yang tak terbaca — teks tak
      // terbaca yang TIDAK kosong tampil apa adanya, bukan dihapus.
      const t = formatDate(nilai)
      return t === '' ? nilai : t
    }
    default:
      return nilai
  }
}

/** Medan yang KUNCINYA tidak ada di dokumen — ditandai, bukan kotak kosong. */
export function MedanTakAda({ label }: { label: string }) {
  return (
    <div className="field">
      <label className="field__label">{label}</label>
      <input className="field__input field__input--readonly" type="text" value="" readOnly disabled />
      <span className="tria__tak-ada">{PENYESUAIAN.takAdaDiWarisan}</span>
    </div>
  )
}

/**
 * Tanggal BACA-SAJA — kontrol `type="date"` lewat `keInputTanggal` inti.
 *
 * ⛔ `nilai` WAJIB nilai TERSIMPAN (`20180101`), bukan nilai tampilan.
 * `keInputTanggal` mengenali `DD-MM-YYYY`, `YYYY-MM-DD`, `YYYYMMDD` — dan
 * nol bentuk bergaris miring; mengisi kotak dengan terjemahan tampil
 * (`01/01/18`) membuatnya kosong. Pola itu rusak di modul tetangga; uji
 * `penyesuaian.test.tsx` menjaga layar ini tidak mengulangnya.
 */
export function TanggalBacaSaja({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <label className="field__label">{label}</label>
      <input className="field__input field__input--readonly" type="date" value={keInputTanggal(nilai)} readOnly />
    </div>
  )
}

/** Teks panjang BACA-SAJA — `Area` inti tidak punya `readOnly`. */
export function AreaBacaSaja({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field field--lebar">
      <label className="field__label">{label}</label>
      <textarea className="field__input field__input--readonly" rows={4} value={nilai} readOnly />
    </div>
  )
}

export function Centang({
  label,
  nilai,
  bacaSaja,
  onUbah,
}: {
  label: string
  nilai: string
  bacaSaja: boolean
  onUbah: (v: string) => void
}) {
  return (
    <label className="tria__centang">
      <input
        type="checkbox"
        checked={nilai === 'true'}
        // ⛔ `disabled`, bukan `readOnly`: kotak centang ber-`readOnly`
        // TETAP dapat dicentang di peramban.
        disabled={bacaSaja}
        onChange={(e) => {
          onUbah(e.target.checked ? 'true' : 'false')
        }}
      />
      {label}
    </label>
  )
}

/** Kotak "belum ada KODE" — SENGAJA bukan `Kosong` ("belum ada DATA"). */
export function BelumDibangun({ judul, petunjuk = PENYESUAIAN.belumDibangunPetunjuk }: { judul: string; petunjuk?: string }) {
  return (
    <div className="tria__belum" role="note">
      <span className="tria__belum-judul">{judul}</span>
      <span className="tria__belum-petunjuk">{petunjuk}</span>
    </div>
  )
}
