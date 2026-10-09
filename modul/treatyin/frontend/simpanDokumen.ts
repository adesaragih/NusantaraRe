// Penyusun isi tombol tulis form Treaty In — Save, Submit, Actions, Decline
// offer. Isi layar → properti `TreatyIn` EJAAN PEGA, satu lawan satu dengan
// peta tabel pendaratan (`repository.PetaPendaratan`) dan kepala `TREATY_IN`.
//
// ⭐ Yang dikirim hanya yang LAYAR pegang. Server menimpakannya ke dokumen
// tersimpan, jadi tab yang tidak pernah dibuka tidak terhapus — padanan
// clipboard Pega.
//
// ⛔ Status akseptasi, posisi, dan riwayat (`CommentList`) TIDAK dikirim:
// hanya tombolnya (di server) yang mengubahnya.

import type { BarisKursWarisan, KursSimpan, LimitsAkar, ShareNP, SimpulLimit } from './api'
import { keSimpan } from './components/tanggalIso'
import type { Halaman } from './halaman'

/** Medan kepala form, apa adanya dari keadaan `FormKontrakTreatyIn`. */
export interface KepalaForm {
  nonProporsional: boolean
  nama: string
  rujukan: string
  wilayah: string
  bordereaux: string
  bordereauxNote: string
  /** Bentuk kabel `DD-MM-YYYY` (atau `YYYYMMDD`). */
  mulai: string
  berakhir: string
  tahunTreaty: string
  pembukuan: string
  pembukuanNonProp: string
  cedant: string
  idCedant: string
  asalBisnis: string
  idAsalBisnis: string
  pemimpin: boolean
}

/** `DD-MM-YYYY`, `YYYY-MM-DD`, atau `YYYYMMDD` → `YYYYMMDD`; lainnya kosong. */
export function keYYYYMMDD(nilai: string): string {
  const s = nilai.trim()
  if (/^\d{8}$/.test(s)) return s
  return keSimpan(s)
}

/**
 * Dokumen `TreatyIn` dari isi layar.
 *
 *   penampung halaman   seluruh properti tab yang pernah dibuka (ejaan Pega)
 *   kepala form         19 medan kepala — menang atas penampung
 *   Limits Non-Prop     `Limits` + `LimitSummaryList` + `TotalLimitsROL` +
 *                       keempat total (`Total*NP`), bila tab itu disentuh
 *   Share Non-Prop      skalar + larik Share/Retro/Fakultatif + total, bila
 *                       tab itu disentuh
 */
export function susunDokumen(
  k: KepalaForm,
  halaman: Halaman,
  limitsNP: { layers: readonly SimpulLimit[]; akar: LimitsAkar } | null,
  shareNP: ShareNP | null,
): Record<string, unknown> {
  const doc: Record<string, unknown> = { ...halaman }
  Object.assign(doc, {
    // ⛔ Nilai TERSIMPAN, bukan label tampil (`Non Proportional`).
    ProportionType: k.nonProporsional ? 'NonProportional' : 'Proportional',
    TreatyContractName: k.nama,
    ContractRefNo: k.rujukan,
    TeritorialScope: k.wilayah,
    Bordeaux: k.bordereaux,
    BordereauxNote: k.bordereauxNote,
    Commencement: keYYYYMMDD(k.mulai),
    Termination: keYYYYMMDD(k.berakhir),
    TreatyYear: k.tahunTreaty,
    AccountingMode: k.pembukuan,
    AccountingModeNonProp: k.pembukuanNonProp,
    Ceding: k.cedant,
    CedingID: k.idCedant,
    LeadingReinsSource: k.asalBisnis,
    LeadingReinsSourceID: k.idAsalBisnis,
    TreatyLeader: k.pemimpin ? 'true' : 'false',
  })
  if (limitsNP !== null) {
    doc.Limits = limitsNP.layers
    doc.LimitSummaryList = limitsNP.akar.LimitSummaryList
    doc.TotalLimitsROL = limitsNP.akar.TotalLimitsROL
    for (const [nama, larik] of Object.entries(limitsNP.akar.Total)) doc[nama] = larik
  }
  if (shareNP !== null) {
    const { Total, ...lain } = shareNP
    Object.assign(doc, lain)
    for (const [nama, larik] of Object.entries(Total)) doc[nama] = larik
  }
  return doc
}

/**
 * Grid Rate of Exchange → `KursSimpan`: hanya baris BARU (ber-mata uang) dan
 * baris yang BERUBAH dari yang dimuat. `null` = tidak ada yang berubah.
 *
 * ⛔ Baris yang dihapus dari grid TIDAK dikirim untuk dihapus — keputusan
 * pemilik proses: `TREATYEXCHANGEYEARLY` berkunci TAHUN dan dipakai bersama.
 * Tanggal dikirim `YYYYMMDD` di `…Asli`; server mengubahnya menjadi stempel.
 */
export function kursBerubah(kini: readonly BarisKursWarisan[], asli: readonly BarisKursWarisan[], tahun: string): KursSimpan | null {
  const dimuat = new Map<string, BarisKursWarisan>()
  for (const a of asli) if (a.id !== undefined && a.id !== '') dimuat.set(a.id, a)
  const baris: BarisKursWarisan[] = []
  for (const b of kini) {
    const id = b.id ?? ''
    if (id === '') {
      if (b.mataUangID.trim() === '' && b.mataUang.trim() === '') continue
    } else {
      const a = dimuat.get(id)
      const sama =
        a !== undefined &&
        a.mataUang === b.mataUang &&
        a.mataUangID === b.mataUangID &&
        a.nilaiKeIDR === b.nilaiKeIDR &&
        keYYYYMMDD(a.berlakuDariAsli) === keYYYYMMDD(b.berlakuDariAsli) &&
        keYYYYMMDD(a.berlakuSampaiAsli) === keYYYYMMDD(b.berlakuSampaiAsli)
      if (sama) continue
    }
    baris.push({ ...b, berlakuDariAsli: keYYYYMMDD(b.berlakuDariAsli), berlakuSampaiAsli: keYYYYMMDD(b.berlakuSampaiAsli) })
  }
  return baris.length === 0 ? null : { tahun, baris }
}

/**
 * Syarat tampil tombol `Actions` — `TreatyInActionButtons`: pemakai memegang
 * workbasket `TreatyIn.Position`, dan posisi itu anak tangga PENYETUJU
 * (SecHead / DeptHead / Director). Admin memakai `Submit` di tab
 * Information & Submit.
 */
export const POSISI_PENYETUJU = ['ReasTreatyInSecHead', 'ReasTreatyInDeptHead', 'ReasTreatyInDirector'] as const

export function bolehActions(posisi: string, workbasket: readonly string[], status: string): boolean {
  if (status === 'Resolve Complete') return false
  return (POSISI_PENYETUJU as readonly string[]).includes(posisi) && workbasket.includes(posisi)
}
