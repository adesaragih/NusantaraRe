// Tampilan polis NB Fac In HANYA-BACA - tab Loss Record (FIRE) dan Correspondence, HANYA akun ber-JOB_POSITION "IT"
// (keputusan work owner 11-10-2026 "Ikut NB: hanya akun IT"; nbfacin tiket 55 / 56). Data dari rute sendiri
// `GET /kasus/{id}/polis-nb-it` (ditolak 403 bagi akun lain - tidak ikut halaman polis). Rupa = `modul/nbfacin/frontend/
// components/{TabLossRecord,TabCorrespondence}.tsx` baca-saja (pola, bukan impor); bar unggah CSV, Tambah / Hapus, ubah
// baris, editor email, Save / Send Email tidak dibawa. Label VERBATIM `modul/nbfacin/frontend/labels.ts`.
//
// Nilai TERSIMPAN (Uw Year, Loss Amount ( IDR ), Loss Ratio hasil Save NB terakhir); Total Loss per Class Of Business +
// Currency = jumlah baris tersimpan (`totalCobLossRecord` nbfacin, eksak).

import { minta } from '../../../inti/frontend/klien'
import { PREFIX_CFI } from './api'
import { NOL, eksak, tambahE, teksE, type Eksak } from './polisNB'

// ---------------------------------------------------------------- label (VERBATIM nbfacin/labels.ts)

/** `LOSS_RECORD` (tanpa bar unggah / tombol). */
export const LOSS_RECORD_NB = {
  premiPar: 'Premi PAR',
  premiMb: 'Premi MB',
  premiPl: 'Premi PL',
  premiOther: 'Premi Other COB',
  uwYear: 'Uw Year',
  classOfBusiness: 'Class Of Business',
  dateOfLoss: 'Date of Loss',
  causeOfLoss: 'Cause Of Loss',
  currency: 'Currency',
  lossAmount: 'Loss Amount',
  lossAmountIdr: 'Loss Amount ( IDR )',
  totalLoss: 'Total Loss Amount',
  totalLossIdr: 'Total Loss Amount ( IDR )',
  rasio: ['Loss Ratio 1 Year', 'Loss Ratio 3 Year', 'Loss Ratio 5 Year'],
  gross: 'Gross',
  net: 'Net',
} as const
/** `TEKS_LOSS_RECORD.kosong`. */
export const TEKS_LOSS_RECORD_NB = { kosong: 'No items' } as const
/** `KORESPONDENSI` (tanpa tombol). */
export const KORESPONDENSI_NB = {
  kolom: ['Date', 'Sender', 'Subject', 'Recipient', 'Status'],
  sender: 'Sender',
  senderMail: 'Sender Mail Address',
  penerima: 'Recipients:',
  subject: 'Subject',
  message: 'Message',
  belumTerkirim: 'Not Sent',
  terkirim: 'Sent',
} as const
/** `TEKS_KORESPONDENSI` (yang dipakai tampilan hanya-baca). */
export const TEKS_KORESPONDENSI_NB = { kosong: 'No items' } as const

// ---------------------------------------------------------------- data (`models.PolisNBIT` backend)

export interface BarisLossRecordNB {
  uwYear: string
  classOfBusiness: string
  dateOfLoss: string
  causeOfLoss: string
  currency: string
  lossAmount: string
  lossAmountIdr: string
}

export interface LossRecordNB {
  premiPar: string
  premiMb: string
  premiPl: string
  premiOther: string
  gross1: string
  gross3: string
  gross5: string
  net1: string
  net3: string
  net5: string
  baris: BarisLossRecordNB[]
}

export interface PenerimaNB {
  recipient: string
  recipientMail: string
}

export interface KorespondensiNB {
  tanggal: string
  sender: string
  senderMail: string
  subject: string
  message: string
  terkirim: boolean
  penerima: PenerimaNB[]
}

export interface PolisNBIT {
  lossRecord: LossRecordNB
  korespondensi: KorespondensiNB[]
}

/** Isi tab IT polis klaim; ditolak (bukan akun IT) / gagal = galat - layar menyembunyikan kedua tab (pola lunak nbfacin
 * `AksesTabIT`). */
export function ambilPolisNBIT(kasusId: string): Promise<PolisNBIT> {
  return minta<PolisNBIT>(`${PREFIX_CFI}/kasus/${encodeURIComponent(kasusId)}/polis-nb-it`)
}

// ---------------------------------------------------------------- hitungan tampilan

export interface TotalCobNB {
  classOfBusiness: string
  currency: string
  lossAmount: string
  /** '' bila seluruh Loss Amount ( IDR ) kelompok kosong. */
  lossAmountIdr: string
}

/** `totalCobLossRecord` nbfacin: per (Class Of Business, Currency) urut kemunculan; IDR kosong dilewati. */
export function totalCobNB(baris: readonly BarisLossRecordNB[]): TotalCobNB[] {
  const t: { classOfBusiness: string; currency: string; jumlah: Eksak; idr: Eksak | null }[] = []
  for (const b of baris) {
    let r = t.find((x) => x.classOfBusiness === b.classOfBusiness && x.currency === b.currency)
    if (!r) t.push((r = { classOfBusiness: b.classOfBusiness, currency: b.currency, jumlah: NOL, idr: null }))
    r.jumlah = tambahE(r.jumlah, eksak(b.lossAmount))
    if (b.lossAmountIdr.trim() !== '') r.idr = tambahE(r.idr ?? NOL, eksak(b.lossAmountIdr))
  }
  return t.map((r) => ({
    classOfBusiness: r.classOfBusiness,
    currency: r.currency,
    lossAmount: teksE(r.jumlah),
    lossAmountIdr: r.idr === null ? '' : teksE(r.idr),
  }))
}

/** Kolom Recipient grid (`daftarNamaPenerima` nbfacin): nama penerima (atau alamatnya), dipisah koma. */
export const daftarNamaPenerimaNB = (p: readonly PenerimaNB[]) =>
  p
    .map((x) => x.recipient.trim() || x.recipientMail.trim())
    .filter((x) => x !== '')
    .join(', ')

const JAM_WIB = 7

/**
 * Kolom Date korespondensi = `DD/MM/YY HH:MM` WIB (`tanggalKorespondensi` nbfacin atas kabel `waktuKorespondensiKeKabel`):
 * stempel Pega `YYYYMMDDTHHMMSS.mmm GMT` digeser ke WIB; kabel `DD-MM-YYYY HH:MM` dipendekkan; teks lain apa adanya.
 */
export function tanggalKorespondensiNB(v: string): string {
  const t = v.trim()
  const p = /^(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})\d{2}(?:\.\d+)? GMT$/.exec(t)
  if (p) {
    const d = new Date(Date.UTC(+p[1]!, +p[2]! - 1, +p[3]!, +p[4]! + JAM_WIB, +p[5]!))
    const dua = (n: number) => String(n).padStart(2, '0')
    return `${dua(d.getUTCDate())}/${dua(d.getUTCMonth() + 1)}/${dua(d.getUTCFullYear() % 100)} ${dua(d.getUTCHours())}:${dua(d.getUTCMinutes())}`
  }
  const k = /^(\d{2})-(\d{2})-\d{2}(\d{2}) (\d{2}:\d{2})$/.exec(t)
  return k ? `${k[1]}/${k[2]}/${k[3]} ${k[4]}` : t
}

// ---------------------------------------------------------------- pembersih HTML pesan (VERBATIM nbfacin
// `components/htmlKorespondensi.ts`): hanya tag polos yang diizinkan, tanpa atribut; isi script / style dibuang; teks
// di-escape. Hasilnya saja yang boleh ditaruh ke innerHTML.

export const TAG_HTML_KORESPONDENSI_NB = ['b', 'strong', 'i', 'em', 'u', 'p', 'br', 'div', 'ul', 'ol', 'li'] as const

const DIIZINKAN = new Set<string>(TAG_HTML_KORESPONDENSI_NB)
const ISI_DIBUANG = new Set(['script', 'style'])
const ENTITAS = /^&(?:[a-zA-Z][a-zA-Z0-9]{1,31}|#\d{1,7}|#[xX][0-9a-fA-F]{1,6});/

function escapeTeks(teks: string): string {
  let hasil = ''
  for (let i = 0; i < teks.length; i++) {
    const c = teks[i]
    if (c === '&') {
      const m = ENTITAS.exec(teks.slice(i))
      if (m) {
        hasil += m[0]
        i += m[0].length - 1
      } else hasil += '&amp;'
    } else if (c === '<') hasil += '&lt;'
    else if (c === '>') hasil += '&gt;'
    else hasil += c
  }
  return hasil
}

export function bersihkanHtmlNB(html: string): string {
  let hasil = ''
  let buang: string | null = null
  const pola = /<(\/?)([a-zA-Z][a-zA-Z0-9]*)\b[^>]*>|<!--[\s\S]*?-->/g
  let akhir = 0
  for (let m = pola.exec(html); m !== null; m = pola.exec(html)) {
    const teks = html.slice(akhir, m.index)
    akhir = pola.lastIndex
    if (buang === null) hasil += escapeTeks(teks)
    if (m[2] === undefined) continue // komentar
    const tutup = m[1] === '/'
    const nama = m[2].toLowerCase()
    if (buang !== null) {
      if (tutup && nama === buang) buang = null
      continue
    }
    if (ISI_DIBUANG.has(nama)) {
      if (!tutup) buang = nama
      continue
    }
    if (!DIIZINKAN.has(nama)) continue
    if (nama === 'br') hasil += tutup ? '' : '<br>'
    else hasil += tutup ? `</${nama}>` : `<${nama}>`
  }
  if (buang === null) hasil += escapeTeks(html.slice(akhir))
  return hasil
}
