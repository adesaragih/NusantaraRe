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
  buangTurunanSpreading(doc)
  return doc
}

/**
 * Daftar TURUNAN tingkat spreading — dihitung ulang backend tiap kontrak
 * dibuka (`sebarManual` / `sebarLama`), dan `T_TREATY_SHARE_SPREADING` tidak
 * punya kolomnya. Dikirim, penyimpan hanya melaporkannya sebagai "No column in
 * the table yet, so NOT saved" (laporan pemakai 10 Oktober 2026). Isian
 * pemakai (jenis dan Pct) tetap terkirim.
 */
export const TURUNAN_SPREADING = [
  'GrossPremiumList',
  'GrossPremiumMinList',
  'DeductionTotalList',
  'NetPremiumList',
  'RnmLimitList',
  // ⭐ 11 Oktober 2026 — rincian anak susunan + bagian OR / R/I tingkat
  // spreading (`SetSpreadingXOL` [5.3.8]–[5.9]), juga turunan tanpa kolom.
  'Currency',
  'Value',
  'BreakDownSprdListXOL',
  'RNMSpreadedListXOL',
  'RNMSpreadedListRIXOL',
  'RNMSpreadedListGrossXOL',
  'RNMSpreadedListGrossRIXOL',
  'RNMSpreadedListGrossMinXOL',
  'RNMSpreadedListGrossRIMinXOL',
  'RNMSpreadedListDeductXOL',
  'RNMSpreadedListDeductRIXOL',
  'RNMSpreadedListNetXOL',
  'RNMSpreadedListNetRIXOL',
] as const

function buangTurunanSpreading(doc: Record<string, unknown>): void {
  for (const kunci of ['Share', 'FacultativeShareList']) {
    const larik = doc[kunci]
    if (!Array.isArray(larik)) continue
    doc[kunci] = larik.map((baris: unknown) => {
      if (typeof baris !== 'object' || baris === null) return baris
      const b = baris as Record<string, unknown>
      const sebar = b.SpreadingListXOL
      if (!Array.isArray(sebar)) return baris
      return {
        ...b,
        SpreadingListXOL: sebar.map((sp: unknown) => {
          if (typeof sp !== 'object' || sp === null) return sp
          const salin = { ...(sp as Record<string, unknown>) }
          for (const k of TURUNAN_SPREADING) delete salin[k]
          return salin
        }),
      }
    })
  }
}

/**
 * Grid Rate of Exchange → `KursSimpan`: SELURUH baris grid (yang kosong
 * dibuang) — yang tidak berubah bertanda `tetap` dan tidak ditulis, yang
 * baru/berubah ditulis. `null` = tidak ada yang berubah.
 *
 * ⛔ Baris yang dihapus dari grid tidak dihapus dari `TREATYEXCHANGEYEARLY`
 * (keputusan pemilik proses: tabel bersama) — ia hanya lepas dari kontrak.
 * Tanggal dikirim `YYYYMMDD` di `…Asli`; server mengubahnya menjadi stempel.
 */
export function kursBerubah(
  kini: readonly BarisKursWarisan[],
  asli: readonly BarisKursWarisan[],
  tahun: string,
  /**
   * Kontrak BARU hasil `Copy`: ia belum memiliki catatan apa pun, jadi seluruh
   * grid milik sumbernya HARUS dikirim walau tak ada yang diubah — kalau tidak,
   * kontrak baru jatuh ke cadangan seluruh kurs tahun treaty (laporan pemakai
   * 10 Oktober 2026: "Copy lalu Save, Rate of Exchange jadi berubah").
   */
  kontrakBaru = false,
): KursSimpan | null {
  const dimuat = new Map<string, BarisKursWarisan>()
  for (const a of asli) if (a.id !== undefined && a.id !== '') dimuat.set(a.id, a)
  // ⭐ `tetap` pada baris DIMUAT = baris itu tercatat MILIK kontrak (backend,
  // `T_TREATY_KURS`). Baris tanpa tanda itu hanya cadangan seluruh-tahun untuk
  // kontrak yang belum punya catatan — BUKAN input kontrak ini. Laporan
  // pemakai 10 Oktober 2026: Save pertama mencatat seluruh kurs tahun itu
  // sebagai milik kontrak, sehingga grid jadi puluhan baris.
  const adaCadangan = asli.some((a) => a.id !== undefined && a.id !== '' && a.tetap !== true)
  const idKini = new Set(kini.map((b) => b.id ?? '').filter((id) => id !== ''))
  const adaDihapus = [...dimuat.keys()].some((id) => !idKini.has(id))
  // Grid cadangan yang baris-barisnya di-Delete pemakai = pemakai MENYUSUN
  // isi kontrak: yang tersisa ikut dicatat.
  const kurasi = adaCadangan && adaDihapus
  const baris: BarisKursWarisan[] = []
  let berubah = false
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
      if (sama) {
        // Milik kontrak (atau disusun pemakai): dikirim bertanda `tetap` — tidak
        // ditulis, hanya dicatat; server mencatat ulang dari SELURUH grid.
        // Cadangan yang tak disentuh: TIDAK dikirim.
        if (a.tetap === true || kurasi) baris.push({ ...b, tetap: true })
        continue
      }
    }
    berubah = true
    const { tetap: _tetap, ...sisa } = b
    baris.push({ ...sisa, berlakuDariAsli: keYYYYMMDD(b.berlakuDariAsli), berlakuSampaiAsli: keYYYYMMDD(b.berlakuSampaiAsli) })
  }
  // Baris milik kontrak yang di-Delete (atau urutan yang berganti) juga perubahan.
  const idAsli = asli.filter((a) => a.tetap === true).map((a) => a.id ?? '').filter((id) => id !== '')
  const idKirim = baris.map((b) => b.id ?? '').filter((id) => id !== '')
  if (idAsli.join('|') !== idKirim.join('|')) berubah = true
  // Salinan: grid milik sumber ikut dikirim bertanda `tetap` supaya tercatat
  // sebagai milik kontrak baru. Grid cadangan (tanpa catatan) tidak dikirim.
  const salinanPunyaGrid = kontrakBaru && asli.some((x) => x.tetap === true) && baris.length > 0
  return berubah || kurasi || salinanPunyaGrid ? { tahun, baris } : null
}

/**
 * Syarat tampil tombol `Actions` — `TreatyInActionButtons`: pemakai memegang
 * workbasket `TreatyIn.Position`, dan posisi itu anak tangga PENYETUJU
 * (SecHead / DeptHead / Director). Admin memakai `Submit` di tab
 * Information & Submit.
 */
export const POSISI_PENYETUJU = ['ReasTreatyInSecHead', 'ReasTreatyInDeptHead', 'ReasTreatyInDirector'] as const

/**
 * Syarat tampil tombol `Save` di form — keputusan pemakai 9 Oktober 2026:
 * *"tombol save muncul apabila workbasketnya sesuai dengan yg di usernya"*.
 * Workbasket berkas = `Position` (kosong = `ReasTreatyInAdmin`, kontrak
 * baru/draf); Save tampil hanya bila akun memegang workbasket itu. Akun
 * divisi IT (Force Edit) dikecualikan. Tetap tidak untuk `Resolve Complete`.
 */
export function bolehSave(posisi: string, workbasket: readonly string[], status: string, divisi: string): boolean {
  if (status === 'Resolve Complete') return false
  const wb = posisi === '' ? 'ReasTreatyInAdmin' : posisi
  return workbasket.includes(wb) || divisi === DIVISI_IT
}

/** Divisi akun (`M_LOGIN_GO.DIVISION_CODE`) pemegang Force Edit. */
export const DIVISI_IT = 'IT'

export function bolehActions(posisi: string, workbasket: readonly string[], status: string): boolean {
  if (status === 'Resolve Complete') return false
  return (POSISI_PENYETUJU as readonly string[]).includes(posisi) && workbasket.includes(posisi)
}
