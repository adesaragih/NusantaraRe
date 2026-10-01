// Penjaga bukti label Endorsement Life.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka korpus XML pada baris
// yang tiap label sebut dan memastikan teksnya memang ada di sana, lalu memastikan setiap kunci label
// berbukti XML ATAU terdaftar `BUKAN_KORPUS` - tidak keduanya, tidak nol.
// Korpus READ-ONLY - hanya dibaca. Bila korpus tidak terjangkau, uji korpus DILEWATI, bukan gagal.

import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import * as KOLOM from './kolomKorpus'
import * as LABEL from './labels'

const KORPUS = 'D:\\XML\\RNM_BRD\\Endorsement Life'
const adaKorpus = existsSync(KORPUS)

const INBOX = 'Section\\InboxEndorsementLife.xml'
const INPUT = 'Section\\InputEDMLife.xml'
const BUAT = 'Section\\EndorsmentLife_Section.xml'
const RINGKAS = 'Section\\ShowLifePremiumSummary_EDM.xml'

type Bukti = readonly [kunci: string, berkas: string, baris: number, tag: string]

/** Kunci `OBJEK.medan` → baris tag korpus yang memuat teksnya. */
export const BUKTI: readonly Bukti[] = [
  ['INBOX_EDM.judul', INBOX, 5505, 'pyValue'],
  ['INBOX_EDM.createAddendum', INBOX, 6620, 'pyLabel'],
  ['INBOX_EDM.kolomCaseId', INBOX, 8418, 'pyValue'],
  ['INBOX_EDM.kolomEndorsementNo', INBOX, 8556, 'pyValue'],
  ['INBOX_EDM.kolomType', INBOX, 8692, 'pyValue'],
  ['INBOX_EDM.kolomEdmType', INBOX, 8828, 'pyValue'],
  ['INBOX_EDM.kolomPolicyNo', INBOX, 8964, 'pyValue'],
  ['INBOX_EDM.kolomSob', INBOX, 9100, 'pyValue'],
  ['INBOX_EDM.kolomCeding', INBOX, 9236, 'pyValue'],
  ['INBOX_EDM.kolomPolicyHolder', INBOX, 9372, 'pyValue'],
  ['INBOX_EDM.kolomMarketingName', INBOX, 9508, 'pyValue'],
  ['INBOX_EDM.kolomCreateDate', INBOX, 9644, 'pyValue'],
  ['INBOX_EDM.kolomCreateOperator', INBOX, 9780, 'pyValue'],
  ['INBOX_EDM.kolomEdmTypeBatal', INBOX, 9916, 'pyValue'],
  ['INBOX_EDM.kolomStatus', INBOX, 10052, 'pyValue'],
  ['KASUS_EDM.judul', INPUT, 817, 'pyValue'],
  ['KASUS_EDM.productName', INPUT, 3399, 'pyLabelFieldValue'],
  ['KASUS_EDM.productNameId', INPUT, 3589, 'pyLabelFieldValue'],
  ['KASUS_EDM.type', INPUT, 3812, 'pyLabelFieldValue'],
  ['KASUS_EDM.reinsuranceSystem', INPUT, 4143, 'pyLabelFieldValue'],
  ['KASUS_EDM.classOfBusiness', INPUT, 4473, 'pyLabelFieldValue'],
  ['KASUS_EDM.sob', INPUT, 4658, 'pyLabelFieldValue'],
  ['KASUS_EDM.policyHolder', INPUT, 4897, 'pyLabelFieldValue'],
  ['KASUS_EDM.premiumMethod', INPUT, 5976, 'pyLabelFieldValue'],
  ['KASUS_EDM.ceding', INPUT, 6339, 'pyLabelFieldValue'],
  ['KASUS_EDM.marketingOfficer', INPUT, 6578, 'pyLabelFieldValue'],
  ['KASUS_EDM.edmType', INPUT, 6910, 'pyLabelFieldValue'],
  ['KASUS_EDM.description', INPUT, 7697, 'pyLabelFieldValue'],
  ['BUAT_EDM.policyNo', BUAT, 1075, 'pyLabelFieldValue'],
  ['BUAT_EDM.edmType', BUAT, 1347, 'pyLabelFieldValue'],
  ['BUAT_EDM.description', BUAT, 2140, 'pyLabelFieldValue'],
  ['BUAT_EDM.edmDate', BUAT, 3160, 'pyLabelFieldValue'],
  ['BUAT_EDM.submit', BUAT, 4226, 'pyLabel'],
  ['POLIS_LAMA_EDM.viewOldPolicy', RINGKAS, 64965, 'pyLabel'],
  ['GRID_EDM.policyNo', INPUT, 12035, 'pyValue'],
  ['GRID_EDM.policyHolder', INPUT, 12177, 'pyValue'],
  ['GRID_EDM.certificateNo', INPUT, 12320, 'pyValue'],
  ['GRID_EDM.nameOfInsured', INPUT, 12463, 'pyValue'],
  ['GRID_EDM.sex', INPUT, 12602, 'pyValue'],
  ['GRID_EDM.dateOfBirth', INPUT, 12738, 'pyValue'],
  ['GRID_EDM.entryAge', INPUT, 12874, 'pyValue'],
  ['GRID_EDM.plan', INPUT, 13010, 'pyValue'],
  ['GRID_EDM.beginDate', INPUT, 13146, 'pyValue'],
  ['GRID_EDM.effectiveDate', INPUT, 13282, 'pyValue'],
  ['GRID_EDM.expiredDate', INPUT, 13417, 'pyValue'],
  ['SIMPAN_EDM.deleteAll', INPUT, 13607, 'pyLabel'],
  ['SIMPAN_EDM.save', INPUT, 37202, 'pyLabel'],
]

/** Kunci tanpa baris korpus - beserta alasannya di komentar. */
export const BUKAN_KORPUS: readonly string[] = [
  'MENU_EDM.kelompok', // nama FOLDER korpus, dibuktikan terpisah di bawah
  'UMUM_EDM.kosong', // teks bawaan grid kosong
  'UMUM_EDM.tutup', // ikon tutup harness `pxIconCancel` tanpa teks
  'UMUM_EDM.kembali', // kembali dari layar kasus ke kotak masuk (`finishAssignment`)
  'UMUM_EDM.memuat', // teks bawaan muat
  'UMUM_EDM.memeriksa', // gerbang sedang berjalan
  'UMUM_EDM.rinci', // tombol buka rincian baris (`pyEditingMode` `expandPane`, tanpa teks di korpus)
  'UMUM_EDM.terkunci', // pesan penjelas kunci field - spec §9 `[keputusan work owner + desain]`, AC 34
  'UMUM_EDM.menyimpan', // `Save` sedang berjalan
  'UMUM_EDM.tandaiHapus', // nama aksesibel kotak centang `.EdmBatal` b15753 (`pyLabelPreview` kosong)
]

/** Nilai label untuk kunci `OBJEK.medan`. */
function nilai(kunci: string): string {
  const [objek, medan] = kunci.split('.') as [string, string]
  const o = (LABEL as unknown as Record<string, Record<string, string>>)[objek]
  const v = o?.[medan]
  if (v === undefined) throw new Error(`label ${kunci} tidak ada`)
  return v
}

/** Setiap kunci label ber-objek `*_EDM` (peta opsi `OPSI_*` diuji di `tampilan.test.ts`). */
function semuaKunci(): string[] {
  return Object.entries(LABEL as unknown as Record<string, Record<string, string>>)
    .filter(([objek]) => objek.endsWith('_EDM'))
    .flatMap(([objek, isi]) => Object.keys(isi).map((m) => `${objek}.${m}`))
}

const xml = (t: string): string => t.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

describe.skipIf(!adaKorpus)('label Endorsement Life berbukti barisnya', () => {
  const baris = new Map<string, string[]>()
  const baca = (berkas: string, nomor: number): string => {
    if (!baris.has(berkas)) baris.set(berkas, readFileSync(join(KORPUS, berkas), 'utf8').split('\n'))
    return (baris.get(berkas)?.[nomor - 1] ?? '').trim()
  }

  it.each(BUKTI.map((b) => [...b]))('%s = baris %s:%i <%s>', (kunci, berkas, nomor, tag) => {
    expect(baca(berkas as string, nomor as number)).toBe(`<${tag}>${xml(nilai(kunci as string))}</${tag}>`)
  })

  it('nama menu = nama folder korpus', () => {
    expect(readdirSync(join(KORPUS, '..'))).toContain(LABEL.MENU_EDM.kelompok)
  })

  it('keempat tombol View Old Policy berlabel sama (b64965, b65522, b66083, b66640)', () => {
    for (const n of [64965, 65522, 66083, 66640]) expect(baca(RINGKAS, n)).toBe('<pyLabel>View Old Policy</pyLabel>')
  })

  // Setiap kolom terbangkitkan (`kolomKorpus.ts`) berbukti barisnya: label VERBATIM pada tag yang disebut.
  const semuaKolom = Object.entries(KOLOM).flatMap(([nama, isi]) => {
    if (Array.isArray(isi)) return (isi as KOLOM.KolomKorpus[]).map((k) => [nama, k] as const)
    if (typeof isi === 'object' && isi !== null) {
      return Object.entries(isi as unknown as Record<string, readonly KOLOM.KolomKorpus[]>).flatMap(([tipe, d]) => d.map((k) => [`${nama}.${tipe}`, k] as const))
    }
    return []
  })

  it('kolom terbangkitkan cukup banyak (pembacanya tidak rusak)', () => {
    expect(semuaKolom.length).toBeGreaterThan(250)
  })

  it.each(semuaKolom.map(([nama, k]) => [nama, k.label, k.berkas, k.baris, k.tag]))('%s %s = %s:%i <%s>', (_, label, berkas, nomor, tag) => {
    expect(baca(berkas as string, nomor as number)).toBe(`<${tag}>${xml(label as string)}</${tag}>`)
  })

  it('halaman awal = satu-satunya harness portal (kelas Data-Portal)', () => {
    const harness = readdirSync(join(KORPUS, 'Harness')).filter((f) => f.endsWith('.xml'))
    const portal = harness.filter((f) => readFileSync(join(KORPUS, 'Harness', f), 'utf8').includes('<pyClassName>Data-Portal</pyClassName>'))
    expect(portal).toEqual(['InboxEndorsementLife.xml'])
    expect(baca('Harness\\InboxEndorsementLife.xml', 843)).toBe('<pyInclude>InboxEndorsementLife</pyInclude>')
  })
})

describe('setiap label berbukti XML atau dinyatakan bukan korpus', () => {
  it('dua daftar menutup seluruh kunci, tanpa tumpang tindih', () => {
    const berbukti = new Set(BUKTI.map((b) => b[0]))
    const bukan = new Set(BUKAN_KORPUS)
    for (const k of berbukti) expect(bukan.has(k), k).toBe(false)
    expect(new Set([...berbukti, ...bukan])).toEqual(new Set(semuaKunci()))
    expect(berbukti.size).toBe(BUKTI.length)
  })
})
