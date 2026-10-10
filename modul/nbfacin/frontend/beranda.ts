// Kotak masuk Beranda modul NB FacIn (work owner 07-10-2026: login Team Leader tidak melihat case NB FacIn yang sudah
// di-submit ke "TEAM LEADER'S INBOX"). Pola modul NB Treaty In (keputusan work owner 06-10-2026: Beranda menunjukkan case
// yang menunggu akun per workbasket yang ia pegang; klik ID membuka layar case lewat `PropsRute.bukaKasus`).
//
// Kolom = gambar layar inbox "New Business" Pega dari work owner 07-10-2026 ("di tampilan inbox, keluarin kolom2 ini"):
// ID · Insured Name · Business name · Source Of Business · Ceding Company · Inception Date · %Share RNM · Marketing Name ·
// Status · Created date · Offering Date · Validity Date · Status Inbox · (tanpa judul: "Accept by <PIC>") · Time Since Last
// Update. Rule / report inbox itu TIDAK ada di korpus (label "Status Inbox", "Time Since Last Update" nihil) - label dari gambar
// layar, sama dengan cara modul NB Treaty In. Kolom tanpa isi di gambar (judul kosong tanpa data) tidak ditampilkan.

import type { DaftarBeranda } from '../../../inti/frontend/modul'
import { formatNumber } from '../../../inti/frontend/lib/format'
import { daftarMenungguNbFacIn, type BarisKotakMasukNbFacIn } from './api'
import { KOLOM_INBOX_NBFACIN, KOLOM_INBOX_UW_NBFACIN, WORKBASKET_UW_NBFACIN } from './labels'

/** %Share RNM tepat 2 desimal seperti gambar layar ("10,00"); `formatNumber` inti membuang nol di ekor. */
export function persenInbox(v: string): string {
  if (v.trim() === '') return ''
  const [kiri = '', kanan = ''] = formatNumber(v, 2).split(',')
  return `${kiri},${kanan.padEnd(2, '0')}`
}

/** Kabel `DD-MM-YYYY` -> `DD/MM/YYYY`; `pendek` = `DD/MM/YY` (Inception Date di gambar "31/05/26"). */
export function tanggalInbox(kabel: string, pendek = false): string {
  const m = /^(\d{2})-(\d{2})-(\d{4})/.exec(kabel.trim())
  if (!m) return kabel
  return `${m[1]}/${m[2]}/${pendek ? (m[3] ?? '').slice(2) : m[3]}`
}

/** Waktu kabel `YYYY-MM-DD HH:MM[:SS]` sebagai waktu lokal; selain itu null. */
function bacaWaktu(t: string): Date | null {
  const m = /^(\d{4})-(\d{2})-(\d{2})(?:[ T](\d{2}):(\d{2})(?::(\d{2}))?)?/.exec(t.trim())
  if (m === null) return null
  const n = (s: string | undefined) => Number.parseInt(s ?? '0', 10)
  return new Date(n(m[1]), n(m[2]) - 1, n(m[3]), n(m[4]), n(m[5]), n(m[6]))
}

/**
 * Lama sejak `tgl` sampai `kini` gaya kolom "Time Since Last Update" (gambar: "1m ago") - aturan sama dengan modul NB Treaty
 * In (`modul/nbtreatyin/frontend/beranda.ts` `sejak`; disalin karena modul tidak saling mengimpor).
 */
export function sejak(tgl: string, kini: Date): string {
  const w = bacaWaktu(tgl)
  if (w === null) return ''
  const detik = Math.max(0, Math.floor((kini.getTime() - w.getTime()) / 1000))
  if (detik < 60) return 'just now'
  if (detik < 3600) return `${Math.floor(detik / 60)}m ago`
  if (detik < 86400) return `${Math.floor(detik / 3600)}h ago`
  let bulan = (kini.getFullYear() - w.getFullYear()) * 12 + (kini.getMonth() - w.getMonth())
  if (kini.getDate() < w.getDate()) bulan--
  if (bulan < 1) return `${Math.floor(detik / 86400)}d ago`
  if (bulan < 12) return `${bulan}mo ago`
  const tahun = Math.floor(bulan / 12)
  const sisa = bulan % 12
  return sisa === 0 ? `${tahun}y ago` : `${tahun}y ${sisa}mo ago`
}

/**
 * Lama di Underwriting gaya gambar layar ("1d 1h ago"): hari + jam bila >= 1 hari, selain itu gaya `sejak`. Hanya bentuk
 * "1d 1h ago" yang terlihat di gambar; bentuk di bawah sehari mengikuti `sejak` (keputusan agent).
 */
export function lamaDiUw(tgl: string, kini: Date): string {
  const w = bacaWaktu(tgl)
  if (w === null) return ''
  const detik = Math.max(0, Math.floor((kini.getTime() - w.getTime()) / 1000))
  if (detik < 86400) return sejak(tgl, kini)
  const hari = Math.floor(detik / 86400)
  const jam = Math.floor((detik % 86400) / 3600)
  return jam === 0 ? `${hari}d ago` : `${hari}d ${jam}h ago`
}

/** Kolom "Time in Underwriting" tampil untuk workbasket UW; daftar semua workbasket (null) bila ada baris yang sudah ke UW. */
export function tampilKolomUw(workbasket: string | null, baris: readonly BarisKotakMasukNbFacIn[]): boolean {
  if (workbasket !== null) return (WORKBASKET_UW_NBFACIN as readonly string[]).includes(workbasket)
  return baris.some((b) => (b.submitToUw ?? '').trim() !== '')
}

/** Baris kotak masuk -> daftar Beranda (kolom `KOLOM_INBOX_NBFACIN`, urutan gambar layar; + "Time in Underwriting" untuk UW). */
export function keDaftarBerandaNbFacIn(
  baris: readonly BarisKotakMasukNbFacIn[],
  kini: Date,
  workbasket: string | null = null,
): DaftarBeranda {
  const kolom: { kunci: string; label: string }[] = KOLOM_INBOX_NBFACIN.map((k) => ({ kunci: k.kunci, label: k.label }))
  if (tampilKolomUw(workbasket, baris)) {
    const i = kolom.findIndex((k) => k.kunci === 'validitas')
    kolom.splice(i + 1, 0, { kunci: KOLOM_INBOX_UW_NBFACIN.kunci, label: KOLOM_INBOX_UW_NBFACIN.label })
  }
  return {
    kolom,
    baris: baris.map((b) => ({
      id: b.caseId,
      sel: {
        id: b.caseId,
        tertanggung: b.insuredName,
        bisnis: b.businessName,
        sob: b.sourceOfBusiness,
        ceding: b.cedingCompany,
        mulai: tanggalInbox(b.inceptionDate, true),
        share: persenInbox(b.shareRnm),
        marketing: b.marketing,
        status: b.status,
        dibuat: tanggalInbox(b.createdDate),
        penawaran: tanggalInbox(b.offeringDate),
        validitas: tanggalInbox(b.validityDate),
        diUw: lamaDiUw(b.submitToUw ?? '', kini),
        statusInbox: b.statusInbox,
        aksiTerakhir: b.aksiTerakhir,
        sejak: sejak(b.tglUpdate, kini),
      },
    })),
  }
}

/** Penyedia daftar kotak masuk Beranda (`MenuModul.daftarBeranda`). */
export async function daftarBeranda(workbasket: string | null): Promise<DaftarBeranda> {
  return keDaftarBerandaNbFacIn(await daftarMenungguNbFacIn(workbasket), new Date(), workbasket)
}
