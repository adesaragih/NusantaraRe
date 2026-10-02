// Penjaga bukti label dan tombol Master Product Name Life - paket 10.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka korpus XML pada baris yang tiap
// label sebut dan memastikan teksnya memang ada di sana, lalu memastikan:
//   - setiap kunci label berbukti XML ATAU terdaftar `[tidak ada di korpus]` - tidak keduanya, tidak nol;
//   - setiap tombol/tautan yang layar render bertekskan label berbukti `pyLabel` / `pySubmitLabel` (atau teks
//     bukan-korpus yang beralasan di `labels.ts`);
//   - ke-25 tombol hidup `InboxProductName` + tautan `View Office Online` + tombol dialog SEMUANYA dirender - tidak
//     ada aksi XML yang dilewati; ketujuh tombol mati beralasan di PARITAS; tombol `Choose*` dan tombol popup
//     pemilihnya sengaja TIDAK dirender (diganti dropdown, keputusan work owner 02-10-2026).
// Korpus READ-ONLY - hanya dibaca. Bila korpus tidak terjangkau, uji korpus DILEWATI, bukan gagal.

import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import * as LABEL from './labels'

const KORPUS = 'D:\\XML\\RNM_BRD\\Master Product Name Life'
const adaKorpus = existsSync(KORPUS)

const S_INBOXPRODUCTNAME = 'Section\\InboxProductName.xml'
const S_SAVEPRODUCTNAME_CONFIRM = 'Section\\SaveProductName_Confirm.xml'
const S_EDITPRODUCTNAME_CONFIRM = 'Section\\EditProductName_Confirm.xml'
const S_CEDING = 'Section\\Ceding_Section.xml'
const S_RIRATE = 'Section\\RIRate_Section.xml'
const S_VIEWRATE = 'Section\\ViewRate.xml'
const F_SAVEPRODUCTNAME_CONFIRM = 'FlowAction\\SaveProductName_Confirm.xml'
const F_EDITPRODUCTNAME_CONFIRM = 'FlowAction\\EditProductName_Confirm.xml'
const F_CHOOSECEDING = 'FlowAction\\ChooseCeding.xml'
const F_VIEWRATE = 'FlowAction\\ViewRate.xml'
const F_PRODUCTNAMEATTACHCONTENT = 'FlowAction\\ProductNameAttachContent.xml'
const A_GENERATEUPLOAD = 'Activity\\GenerateUpload_Act.xml'
const D_COPYPRODUCT = 'DataTransform\\CopyProduct.xml'

type Bukti = readonly [kunci: string, berkas: string, baris: number, tag: string]

/** Kunci `OBJEK.medan` → baris tag korpus yang ISINYA persis teks label (dibangkitkan dari korpus, PARITAS §0). */
const BUKTI: readonly Bukti[] = [
  ['GRID_MPNL.labelSelAdd', S_INBOXPRODUCTNAME, 71783, 'pyLabelFieldValue'],
  ['GRID_MPNL.add', S_INBOXPRODUCTNAME, 71865, 'pyLabel'],
  ['GRID_MPNL.tooltipAdd', S_INBOXPRODUCTNAME, 71863, 'pyTooltip'],
  ['GRID_MPNL.kolomId', S_INBOXPRODUCTNAME, 72403, 'pyValue'],
  ['GRID_MPNL.kolomCeding', S_INBOXPRODUCTNAME, 72551, 'pyValue'],
  ['GRID_MPNL.kolomTreatyNumber', S_INBOXPRODUCTNAME, 72707, 'pyValue'],
  ['GRID_MPNL.kolomTreatyName', S_INBOXPRODUCTNAME, 72866, 'pyValue'],
  ['GRID_MPNL.kolomCreateOp', S_INBOXPRODUCTNAME, 73025, 'pyValue'],
  ['GRID_MPNL.kolomUpdateOp', S_INBOXPRODUCTNAME, 73184, 'pyValue'],
  ['GRID_MPNL.view', S_INBOXPRODUCTNAME, 74798, 'pyLabel'],
  ['UMUM_MPNL.productName', S_INBOXPRODUCTNAME, 3620, 'pyLabelFieldValue'],
  ['UMUM_MPNL.productCode', S_INBOXPRODUCTNAME, 3894, 'pyLabelFieldValue'],
  ['UMUM_MPNL.ceding', S_INBOXPRODUCTNAME, 4075, 'pyLabelFieldValue'],
  ['UMUM_MPNL.chooseCeding', S_INBOXPRODUCTNAME, 5222, 'pyLabel'],
  ['UMUM_MPNL.sob', S_INBOXPRODUCTNAME, 4463, 'pyLabelFieldValue'],
  ['UMUM_MPNL.chooseSob', S_INBOXPRODUCTNAME, 5563, 'pyLabel'],
  ['UMUM_MPNL.deduction', S_INBOXPRODUCTNAME, 7104, 'pyLabelFieldValue'],
  ['UMUM_MPNL.riRisk', S_INBOXPRODUCTNAME, 7398, 'pyLabelFieldValue'],
  ['UMUM_MPNL.chooseRiRisk', S_INBOXPRODUCTNAME, 8205, 'pyLabel'],
  ['UMUM_MPNL.treatyName', S_INBOXPRODUCTNAME, 8719, 'pyLabelFieldValue'],
  ['UMUM_MPNL.treatyNumber', S_INBOXPRODUCTNAME, 8900, 'pyLabelFieldValue'],
  ['UMUM_MPNL.causeOfLoss', S_INBOXPRODUCTNAME, 10661, 'pyLabelFieldValue'],
  ['UMUM_MPNL.chooseCause', S_INBOXPRODUCTNAME, 11360, 'pyLabel'],
  ['UMUM_MPNL.onRetention', S_INBOXPRODUCTNAME, 47312, 'pyCheckboxCaption'],
  ['INWARD_MPNL.policyHolder', S_INBOXPRODUCTNAME, 17097, 'pyLabelFieldValue'],
  ['INWARD_MPNL.choosePolicyHolder', S_INBOXPRODUCTNAME, 17827, 'pyLabel'],
  ['INWARD_MPNL.insured', S_INBOXPRODUCTNAME, 18341, 'pyLabelFieldValue'],
  ['INWARD_MPNL.addendumNo', S_INBOXPRODUCTNAME, 18834, 'pyLabelFieldValue'],
  ['INWARD_MPNL.addendum', S_INBOXPRODUCTNAME, 19432, 'pyLabelFieldValue'],
  ['INWARD_MPNL.amandementNo', S_INBOXPRODUCTNAME, 20155, 'pyLabelFieldValue'],
  ['INWARD_MPNL.amandement', S_INBOXPRODUCTNAME, 20760, 'pyLabelFieldValue'],
  ['INWARD_MPNL.maxExpiredClaim', S_INBOXPRODUCTNAME, 21781, 'pyLabelFieldValue'],
  ['INWARD_MPNL.begin', S_INBOXPRODUCTNAME, 21969, 'pyLabelFieldValue'],
  ['INWARD_MPNL.stnc', S_INBOXPRODUCTNAME, 22304, 'pyLabelFieldValue'],
  ['INWARD_MPNL.cedingRetention', S_INBOXPRODUCTNAME, 22494, 'pyLabelFieldValue'],
  ['INWARD_MPNL.cedingLimit', S_INBOXPRODUCTNAME, 22657, 'pyLabelFieldValue'],
  ['INWARD_MPNL.brokerage', S_INBOXPRODUCTNAME, 22915, 'pyLabelFieldValue'],
  ['INWARD_MPNL.minAge', S_INBOXPRODUCTNAME, 23078, 'pyLabelFieldValue'],
  ['INWARD_MPNL.maxAge', S_INBOXPRODUCTNAME, 23265, 'pyLabelFieldValue'],
  ['INWARD_MPNL.expiryAge', S_INBOXPRODUCTNAME, 23452, 'pyLabelFieldValue'],
  ['INWARD_MPNL.extraPremi', S_INBOXPRODUCTNAME, 23635, 'pyLabelFieldValue'],
  ['INWARD_MPNL.minSumInsured', S_INBOXPRODUCTNAME, 23796, 'pyLabelFieldValue'],
  ['INWARD_MPNL.maxSumInsured', S_INBOXPRODUCTNAME, 23956, 'pyLabelFieldValue'],
  ['INWARD_MPNL.maxSumReasured', S_INBOXPRODUCTNAME, 24214, 'pyLabelFieldValue'],
  ['INWARD_MPNL.rnmShare', S_INBOXPRODUCTNAME, 24674, 'pyLabelFieldValue'],
  ['INWARD_MPNL.ofSumReasured', S_INBOXPRODUCTNAME, 24880, 'pyValue'],
  ['INWARD_MPNL.rnmLimit', S_INBOXPRODUCTNAME, 25236, 'pyLabelFieldValue'],
  ['INWARD_MPNL.premiumFactor', S_INBOXPRODUCTNAME, 25398, 'pyLabelFieldValue'],
  ['INWARD_MPNL.payment', S_INBOXPRODUCTNAME, 25611, 'pyLabelFieldValue'],
  ['INWARD_MPNL.subjectTo', S_INBOXPRODUCTNAME, 25924, 'pyLabelFieldValue'],
  ['INWARD_MPNL.annuityInterest', S_INBOXPRODUCTNAME, 26092, 'pyLabelFieldValue'],
  ['INWARD_MPNL.premiumRefundFactor', S_INBOXPRODUCTNAME, 26306, 'pyLabelFieldValue'],
  ['INWARD_MPNL.maxDataReceive', S_INBOXPRODUCTNAME, 27062, 'pyLabelFieldValue'],
  ['INWARD_MPNL.mature', S_INBOXPRODUCTNAME, 27250, 'pyLabelFieldValue'],
  ['INWARD_MPNL.birthday', S_INBOXPRODUCTNAME, 27960, 'pyLabelFieldValue'],
  ['INWARD_MPNL.currency', S_INBOXPRODUCTNAME, 28140, 'pyLabelFieldValue'],
  ['INWARD_MPNL.chooseCurrency', S_INBOXPRODUCTNAME, 28742, 'pyLabel'],
  ['INWARD_MPNL.extraMortality', S_INBOXPRODUCTNAME, 29256, 'pyLabelFieldValue'],
  ['INWARD_MPNL.maxContract', S_INBOXPRODUCTNAME, 29442, 'pyLabelFieldValue'],
  ['INWARD_MPNL.proportionalTable', S_INBOXPRODUCTNAME, 29626, 'pyLabelFieldValue'],
  ['LIEN_MPNL.judul', S_INBOXPRODUCTNAME, 12201, 'pyValue'],
  ['LIEN_MPNL.usia', S_INBOXPRODUCTNAME, 12741, 'pyValue'],
  ['LIEN_MPNL.manfaat', S_INBOXPRODUCTNAME, 12890, 'pyValue'],
  ['DOKUMEN_MPNL.judul', S_INBOXPRODUCTNAME, 14601, 'pyValue'],
  ['DOKUMEN_MPNL.documentList', S_INBOXPRODUCTNAME, 15125, 'pyValue'],
  ['PLAN_MPNL.judul', S_INBOXPRODUCTNAME, 31557, 'pyValue'],
  ['PLAN_MPNL.planName', S_INBOXPRODUCTNAME, 31845, 'pyValue'],
  ['PLAN_MPNL.bussines', S_INBOXPRODUCTNAME, 31994, 'pyValue'],
  ['PLAN_MPNL.benefit', S_INBOXPRODUCTNAME, 32143, 'pyValue'],
  ['PLAN_MPNL.riRate', S_INBOXPRODUCTNAME, 32296, 'pyValue'],
  ['PLAN_MPNL.add', S_INBOXPRODUCTNAME, 32823, 'pyLabel'],
  ['PLAN_MPNL.viewRate', S_INBOXPRODUCTNAME, 34113, 'pyLabel'],
  ['PLAN_MPNL.chooseRiRate', S_INBOXPRODUCTNAME, 34589, 'pyLabel'],
  ['PLAN_MPNL.delete', S_INBOXPRODUCTNAME, 35075, 'pyLabel'],
  ['FINUW_MPNL.judul', S_INBOXPRODUCTNAME, 37148, 'pyValue'],
  ['FINUW_MPNL.minInsured', S_INBOXPRODUCTNAME, 37670, 'pyValue'],
  ['FINUW_MPNL.maxInsured', S_INBOXPRODUCTNAME, 37818, 'pyValue'],
  ['FINUW_MPNL.employee', S_INBOXPRODUCTNAME, 37966, 'pyValue'],
  ['FINUW_MPNL.nonEmployee', S_INBOXPRODUCTNAME, 38115, 'pyValue'],
  ['FINUW_MPNL.add', S_INBOXPRODUCTNAME, 38425, 'pyLabel'],
  ['FINUW_MPNL.delete', S_INBOXPRODUCTNAME, 40024, 'pyLabel'],
  ['UWLIMIT_MPNL.judul', S_INBOXPRODUCTNAME, 42075, 'pyValue'],
  ['UWLIMIT_MPNL.minInsured', S_INBOXPRODUCTNAME, 42594, 'pyValue'],
  ['UWLIMIT_MPNL.maxInsured', S_INBOXPRODUCTNAME, 42742, 'pyValue'],
  ['UWLIMIT_MPNL.minAge', S_INBOXPRODUCTNAME, 42890, 'pyValue'],
  ['UWLIMIT_MPNL.maxAge', S_INBOXPRODUCTNAME, 43038, 'pyValue'],
  ['UWLIMIT_MPNL.medical', S_INBOXPRODUCTNAME, 43187, 'pyValue'],
  ['UWLIMIT_MPNL.description', S_INBOXPRODUCTNAME, 43336, 'pyValue'],
  ['UWLIMIT_MPNL.add', S_INBOXPRODUCTNAME, 43643, 'pyLabel'],
  ['UWLIMIT_MPNL.delete', S_INBOXPRODUCTNAME, 45682, 'pyLabel'],
  ['KOMENTAR_MPNL.date', S_INBOXPRODUCTNAME, 62095, 'pyValue'],
  ['KOMENTAR_MPNL.pic', S_INBOXPRODUCTNAME, 62246, 'pyValue'],
  ['KOMENTAR_MPNL.comment', S_INBOXPRODUCTNAME, 62399, 'pyValue'],
  ['TOMBOL_MPNL.close', S_INBOXPRODUCTNAME, 58770, 'pyLabel'],
  ['TOMBOL_MPNL.save', S_INBOXPRODUCTNAME, 59041, 'pyLabel'],
  ['TOMBOL_MPNL.edit', S_INBOXPRODUCTNAME, 59489, 'pyLabel'],
  ['TOMBOL_MPNL.copy', S_INBOXPRODUCTNAME, 59854, 'pyLabel'],
  ['TOMBOL_MPNL.generate', S_INBOXPRODUCTNAME, 60122, 'pyLabel'],
  ['SIMPAN_MPNL.tanya', S_SAVEPRODUCTNAME_CONFIRM, 519, 'pyValue'],
  ['SIMPAN_MPNL.comment', S_SAVEPRODUCTNAME_CONFIRM, 1025, 'pyLabelFieldValue'],
  ['SIMPAN_MPNL.save', F_SAVEPRODUCTNAME_CONFIRM, 34, 'pySubmitLabel'],
  ['SIMPAN_MPNL.cancel', F_SAVEPRODUCTNAME_CONFIRM, 33, 'pyCancelLabel'],
  ['EDIT_MPNL.tanya', S_EDITPRODUCTNAME_CONFIRM, 496, 'pyValue'],
  ['EDIT_MPNL.edit', F_EDITPRODUCTNAME_CONFIRM, 19, 'pySubmitLabel'],
  ['EDIT_MPNL.cancel', F_EDITPRODUCTNAME_CONFIRM, 18, 'pyCancelLabel'],
  ['PEMILIH_MPNL.search', S_CEDING, 513, 'pyLabelFieldValue'],
  ['PEMILIH_MPNL.kolomId', S_CEDING, 1502, 'pyValue'],
  ['PEMILIH_MPNL.kolomName', S_CEDING, 1643, 'pyValue'],
  ['PEMILIH_MPNL.choose', S_CEDING, 2265, 'pyLabel'],
  ['PEMILIH_MPNL.submit', F_CHOOSECEDING, 32, 'pySubmitLabel'],
  ['PEMILIH_MPNL.cancel', F_CHOOSECEDING, 31, 'pyCancelLabel'],
  ['PEMILIH_MPNL.kolomRiRateName', S_RIRATE, 1739, 'pyValue'],
  ['RATE_MPNL.judul', S_VIEWRATE, 843, 'pyValue'],
  ['RATE_MPNL.kolomId', S_VIEWRATE, 1147, 'pyValue'],
  ['RATE_MPNL.usedby', S_VIEWRATE, 1293, 'pyValue'],
  ['RATE_MPNL.gender', S_VIEWRATE, 1439, 'pyValue'],
  ['RATE_MPNL.contract', S_VIEWRATE, 1585, 'pyValue'],
  ['RATE_MPNL.age', S_VIEWRATE, 1731, 'pyValue'],
  ['RATE_MPNL.rate', S_VIEWRATE, 1877, 'pyValue'],
  ['RATE_MPNL.submit', F_VIEWRATE, 18, 'pySubmitLabel'],
  ['RATE_MPNL.cancel', F_VIEWRATE, 20, 'pyCancelLabel'],
  ['LAMPIRAN_MPNL.add', S_INBOXPRODUCTNAME, 64747, 'pyLabel'],
  ['LAMPIRAN_MPNL.refresh', S_INBOXPRODUCTNAME, 65270, 'pyLabel'],
  ['LAMPIRAN_MPNL.peringatanNama', S_INBOXPRODUCTNAME, 66071, 'pyValue'],
  ['LAMPIRAN_MPNL.peringatanGanti', S_INBOXPRODUCTNAME, 66230, 'pyValue'],
  ['LAMPIRAN_MPNL.downloadAll', S_INBOXPRODUCTNAME, 67657, 'pyLabel'],
  ['LAMPIRAN_MPNL.fileName', S_INBOXPRODUCTNAME, 68426, 'pyValue'],
  ['LAMPIRAN_MPNL.viewOffice', S_INBOXPRODUCTNAME, 69291, 'pyLabel'],
  ['LAMPIRAN_MPNL.delete', S_INBOXPRODUCTNAME, 69714, 'pyLabel'],
  ['LAMPIRAN_MPNL.attach', F_PRODUCTNAMEATTACHCONTENT, 24, 'pySubmitLabel'],
  ['LAMPIRAN_MPNL.cancel', F_PRODUCTNAMEATTACHCONTENT, 22, 'pyCancelLabel'],
]

/** Teks di dalam ekspresi atau nilai berkutip - baris MEMUATNYA, bukan sama dengannya. */
const BUKTI_SEBAGIAN: ReadonlyArray<readonly [kunci: string, berkas: string, baris: number, potongan: (v: string) => string]> = [
  ['PEMBAYARAN_MPNL.annual', A_GENERATEUPLOAD, 1141, (v) => `PAYMENT==1,"${v}"`],
  ['PEMBAYARAN_MPNL.semiAnnual', A_GENERATEUPLOAD, 1141, (v) => `PAYMENT==2,"${v}"`],
  ['PEMBAYARAN_MPNL.quarterly', A_GENERATEUPLOAD, 1141, (v) => `PAYMENT==3,"${v}"`],
  ['PEMBAYARAN_MPNL.monthly', A_GENERATEUPLOAD, 1141, (v) => `PAYMENT==4,"${v}"`],
  ['PESAN_MPNL.copy', D_COPYPRODUCT, 236, (v) => `<pyPropertiesValue>"${v}"</pyPropertiesValue>`],
]

/** Kunci `[tidak ada di korpus]` - masing-masing beralasan di `labels.ts`. */
const BUKAN_KORPUS: readonly string[] = [
  'MENU_MPNL.kelompok', // nama FOLDER korpus, dibuktikan terpisah di bawah
  'LAIN_MPNL.kosong',
  'LAIN_MPNL.terpotong',
  'LAIN_MPNL.tambahBaris',
  'LAIN_MPNL.hapusBaris',
  'LAIN_MPNL.salinBaris',
  'LAIN_MPNL.terunggah',
  'LAIN_MPNL.gagal',
  'LAIN_MPNL.belum',
  'LAIN_MPNL.ulangi',
  'LAIN_MPNL.dropdownTerpotong',
]

/** Nilai label untuk kunci `OBJEK.medan`. */
function nilai(kunci: string): string {
  const [objek, medan] = kunci.split('.') as [string, string]
  const o = (LABEL as unknown as Record<string, Record<string, string>>)[objek]
  const v = o?.[medan]
  if (v === undefined) throw new Error(`label ${kunci} tidak ada`)
  return v
}

/** Setiap kunci label yang diekspor `labels.ts`. */
function semuaKunci(): string[] {
  return Object.entries(LABEL as unknown as Record<string, Record<string, string>>).flatMap(([objek, isi]) =>
    Object.keys(isi).map((m) => `${objek}.${m}`),
  )
}

const xml = (t: string): string => t.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

const PEMILIH = ['Ceding', 'SOB', 'PolicyHolder', 'Currency', 'RIRISK', 'RIRate', 'CauseOfLoss'] as const

describe.skipIf(!adaKorpus)('label Master Product Name Life berbukti barisnya', () => {
  const baris = new Map<string, string[]>()
  const baca = (berkas: string, nomor: number): string => {
    if (!baris.has(berkas)) baris.set(berkas, readFileSync(join(KORPUS, berkas), 'utf8').split('\n'))
    return (baris.get(berkas)?.[nomor - 1] ?? '').trim()
  }

  it.each(BUKTI.map((b) => [...b]))('%s = baris %s:%i <%s>', (kunci, berkas, nomor, tag) => {
    expect(baca(berkas as string, nomor as number)).toBe(`<${tag}>${xml(nilai(kunci as string))}</${tag}>`)
  })

  it.each(BUKTI_SEBAGIAN.map((b) => [b[0], b[1], b[2]] as const))('%s dimuat baris %s:%i', (kunci, berkas, nomor) => {
    const potongan = BUKTI_SEBAGIAN.find((b) => b[0] === kunci)?.[3] ?? ((v: string) => v)
    expect(baca(berkas, nomor)).toContain(potongan(nilai(kunci)))
  })

  it('nama menu = nama folder korpus', () => {
    expect(readdirSync(join(KORPUS, '..'))).toContain(LABEL.MENU_MPNL.kelompok)
  })

  it('ketujuh section pemilih memuat Search, ID, Name / RIRate Name, Choose; ketujuh FlowAction Submit / Cancel', () => {
    for (const p of PEMILIH) {
      const s = readFileSync(join(KORPUS, 'Section', `${p}_Section.xml`), 'utf8')
      for (const t of [LABEL.PEMILIH_MPNL.search, LABEL.PEMILIH_MPNL.choose, LABEL.PEMILIH_MPNL.kolomId]) expect(s, p).toContain(`>${t}<`)
      expect(s, p).toContain(`>${p === 'RIRate' ? LABEL.PEMILIH_MPNL.kolomRiRateName : LABEL.PEMILIH_MPNL.kolomName}<`)
      const f = readFileSync(join(KORPUS, 'FlowAction', `Choose${p === 'RIRISK' ? 'RIRisk' : p}.xml`), 'utf8')
      expect(f, p).toContain(`<pySubmitLabel>${LABEL.PEMILIH_MPNL.submit}</pySubmitLabel>`)
      expect(f, p).toContain(`<pyCancelLabel>${LABEL.PEMILIH_MPNL.cancel}</pyCancelLabel>`)
    }
  })

  it('sensus tombol: 33 sel pxButton di InboxProductName, 1 per section pemilih (PARITAS §0)', () => {
    const hitung = (f: string): number =>
      readFileSync(join(KORPUS, 'Section', f), 'utf8')
        .split('\n')
        .filter((l) => l.trim() === '<pyFormat>pxButton</pyFormat>').length
    expect(hitung('InboxProductName.xml')).toBe(33)
    for (const p of PEMILIH) expect(hitung(`${p}_Section.xml`), p).toBe(1)
  })
})

describe('setiap label berbukti XML atau dinyatakan bukan korpus', () => {
  it('tiga daftar menutup seluruh kunci, tanpa tumpang tindih', () => {
    const berbukti = new Set([...BUKTI.map((b) => b[0]), ...BUKTI_SEBAGIAN.map((b) => b[0])])
    const bukan = new Set(BUKAN_KORPUS)
    for (const k of berbukti) expect(bukan.has(k), k).toBe(false)
    expect(new Set([...berbukti, ...bukan])).toEqual(new Set(semuaKunci()))
    expect(berbukti.size).toBe(BUKTI.length + BUKTI_SEBAGIAN.length)
  })
})

// ---------------------------------------------------------------------------
// Tombol: yang dirender ⊆ berbukti, dan yang berbukti ⊆ yang dirender.
// ---------------------------------------------------------------------------

/** Isi berkas TSX layar modul. */
function sumberLayar(): { berkas: string; isi: string }[] {
  const akar = __dirname
  return ['pages', 'components'].flatMap((d) =>
    readdirSync(join(akar, d))
      .filter((f) => f.endsWith('.tsx'))
      .map((f) => ({ berkas: join(d, f), isi: readFileSync(join(akar, d, f), 'utf8') })),
  )
}

/** Setiap `<button …>…</button>` dan `<a …>…</a>` (tautan `pxLink`). */
function tombolDirender(): { berkas: string; isi: string }[] {
  return sumberLayar().flatMap((s) =>
    [...s.isi.matchAll(/<button\b[\s\S]*?<\/button>|<a\b[\s\S]*?<\/a>/g)].map((m) => ({ berkas: s.berkas, isi: m[0] })),
  )
}

/** Kunci label teks tombol - anak `{X_MPNL.k}` tepat sebelum penutupnya; atribut (`title`, `aria-label`) tidak dihitung. */
function kunciTeks(isi: string): string[] {
  const m = /\{([A-Z]+_MPNL)\.(\w+)\}\s*<\/(?:button|a)>$/.exec(isi)
  return m ? [`${m[1]}.${m[2]}`] : []
}

/** Tombol `Cancel` kaki FlowAction dirender `Modal` bersama lewat `labelBatal={X_MPNL.cancel}`. */
function kunciBatal(): string[] {
  return sumberLayar().flatMap((s) => [...s.isi.matchAll(/labelBatal=\{([A-Z]+_MPNL)\.(\w+)\}/g)].map((m) => `${m[1]}.${m[2]}`))
}

describe('tombol layar = tombol korpus', () => {
  const tombol = tombolDirender()
  const berbuktiTombol = new Set(
    BUKTI.filter((b) => b[3] === 'pyLabel' || b[3] === 'pySubmitLabel' || b[3] === 'pyCancelLabel').map((b) => b[0]),
  )
  // Ikon grid bawaan / salin baris / kirim ulang - beralasan di `labels.ts` (`LAIN_MPNL`).
  const tambahanSah = new Set(['LAIN_MPNL.tambahBaris', 'LAIN_MPNL.hapusBaris', 'LAIN_MPNL.salinBaris', 'LAIN_MPNL.ulangi'])
  // Keputusan work owner 02-10-2026 ("perubahan pada tampilan untuk semua Choose ubah jadi dropdown saja"): ketujuh
  // tombol `Choose*` dan tombol popup FlowAction `Choose*` (`Choose` baris, `Submit`, `Cancel`) diganti dropdown
  // master (`DropdownMaster`). Labelnya tetap berbukti korpus; tombolnya sengaja TIDAK dirender.
  const digantiDropdown = new Set([
    'UMUM_MPNL.chooseCeding',
    'UMUM_MPNL.chooseSob',
    'UMUM_MPNL.chooseRiRisk',
    'UMUM_MPNL.chooseCause',
    'INWARD_MPNL.choosePolicyHolder',
    'INWARD_MPNL.chooseCurrency',
    'PLAN_MPNL.chooseRiRate',
    'PEMILIH_MPNL.choose',
    'PEMILIH_MPNL.submit',
    'PEMILIH_MPNL.cancel',
  ])

  it('terbaca: ada tombol di layar', () => {
    expect(tombol.length).toBeGreaterThan(25)
  })

  it('setiap tombol yang dirender bertekskan label berbukti (atau teks bukan-korpus beralasan)', () => {
    const langgar: string[] = []
    for (const t of tombol) {
      const kunci = kunciTeks(t.isi)
      if (kunci.length === 0) {
        // Tautan nama berkas b68903 menampilkan `.pyFileName` - nilai data, bukan label.
        if (!t.isi.includes('{l.fileName}')) langgar.push(`${t.berkas}: tombol tanpa teks label: ${t.isi.slice(0, 80)}`)
        continue
      }
      for (const k of kunci) if (!berbuktiTombol.has(k) && !tambahanSah.has(k)) langgar.push(`${t.berkas}: ${k}`)
    }
    expect(langgar).toEqual([])
  })

  it('setiap tombol korpus hidup dirender - tidak ada aksi XML yang dilewati', () => {
    const dirender = new Set([...tombol.flatMap((t) => kunciTeks(t.isi)), ...kunciBatal()])
    const hilang = [...berbuktiTombol].filter((k) => !dirender.has(k) && !digantiDropdown.has(k))
    expect(hilang).toEqual([])
    // 25 pxButton hidup InboxProductName + pxLink `View Office Online` + `Choose` pemilih + 10 tombol FlowAction
    // (`Save`/`Cancel`, `Edit`/`Cancel`, `Submit`/`Cancel` pemilih, `Submit`/`Cancel` View Rate, `Attach`/`Cancel`).
    expect(berbuktiTombol.size).toBe(25 + 1 + 1 + 10)
  })

  it('tombol Choose* dan tombol popup pemilih tidak dirender - diganti dropdown (keputusan work owner 02-10-2026)', () => {
    const dirender = new Set([...tombol.flatMap((t) => kunciTeks(t.isi)), ...kunciBatal()])
    expect([...digantiDropdown].filter((k) => !berbuktiTombol.has(k))).toEqual([])
    expect([...digantiDropdown].filter((k) => dirender.has(k))).toEqual([])
  })
})
