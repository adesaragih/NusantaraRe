// Penjaga bukti label Endorsement Life.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka korpus XML pada baris
// yang tiap label sebut dan memastikan teksnya memang ada di sana, lalu memastikan setiap kunci label
// berbukti XML ATAU terdaftar `BUKAN_KORPUS` - tidak keduanya, tidak nol.
// Korpus READ-ONLY - hanya dibaca. Bila korpus tidak terjangkau, uji korpus DILEWATI, bukan gagal.

import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KEPALA_UNDUH_CSV } from './csv'
import * as KOLOM from './kolomKorpus'
import * as LABEL from './labels'

const KORPUS = 'D:\\XML\\RNM_BRD\\Endorsement Life'
const adaKorpus = existsSync(KORPUS)

const INBOX = 'Section\\InboxEndorsementLife.xml'
const INPUT = 'Section\\InputEDMLife.xml'
const BUAT = 'Section\\EndorsmentLife_Section.xml'
const RINGKAS = 'Section\\ShowLifePremiumSummary_EDM.xml'
const MODAL_CSV = 'FlowAction\\UploadCSV_LifeEndorsement.xml'
const HASIL_CSV = 'Section\\ViewCSVResult_LifeEDM.xml'
const UNDUH_CSV = 'Activity\\GenerateDataDtlLife_act.xml'
const PUTUSAN = 'Section\\ConfirmSection.xml'
const RIWAYAT = 'Activity\\AddHistorySuggest.xml'

type Bukti = readonly [kunci: string, berkas: string, baris: number, tag: string]

/** Kunci `OBJEK.medan` → baris tag korpus yang memuat teksnya. */
export const BUKTI: readonly Bukti[] = [
  ['INBOX_EDM.judul', INBOX, 5505, 'pyValue'],
  ['INBOX_EDM.createAddendum', INBOX, 6620, 'pyLabel'],
  ['INBOX_EDM.kolomCaseId', INBOX, 8418, 'pyValue'],
  ['INBOX_EDM.kolomEndorsementNo', INBOX, 8556, 'pyValue'],
  ['INBOX_EDM.kolomType', INBOX, 8692, 'pyValue'],
  ['INBOX_EDM.kolomEdmType', INBOX, 8828, 'pyValue'],
  ['INBOX_EDM.kolomCeding', INBOX, 9236, 'pyValue'],
  ['INBOX_EDM.kolomCreateDate', INBOX, 9644, 'pyValue'],
  ['INBOX_EDM.kolomCreateOperator', INBOX, 9780, 'pyValue'],
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
  ['BUAT_EDM.edmTypePerubahanData', BUAT, 1645, 'pyLabelFieldValue'],
  ['BUAT_EDM.edmTypeBatal', BUAT, 1937, 'pyLabelFieldValue'],
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
  ['UNGGAH_EDM.uploadCsv', INPUT, 8973, 'pyLabel'],
  ['UNGGAH_EDM.viewUpload', INPUT, 9340, 'pyLabel'],
  ['UNGGAH_EDM.addCsvData', INPUT, 10405, 'pyLabel'],
  ['UNGGAH_EDM.judulModal', MODAL_CSV, 191, 'pyLabel'],
  ['UNGGAH_EDM.generateDataDetail', HASIL_CSV, 14322, 'pyLabel'],
  ['PUTUSAN_EDM.status', PUTUSAN, 496, 'pyLabelFieldValue'],
  ['PUTUSAN_EDM.comment', PUTUSAN, 829, 'pyLabelFieldValue'],
  ['PUTUSAN_EDM.kolomDate', PUTUSAN, 1984, 'pyValue'],
  ['PUTUSAN_EDM.kolomPic', PUTUSAN, 2125, 'pyValue'],
  ['PUTUSAN_EDM.kolomStatus', PUTUSAN, 2263, 'pyValue'],
  ['PUTUSAN_EDM.kolomComment', PUTUSAN, 2399, 'pyValue'],
]

/** Kunci tanpa baris korpus - beserta alasannya di komentar. */
export const BUKAN_KORPUS: readonly string[] = [
  'UNGGAH_EDM.ringkasMasalah', // modal Upload CSV disederhanakan (keputusan work owner 10-10-2026)
  'UNGGAH_EDM.ringkasLolos', // modal Upload CSV disederhanakan (keputusan work owner 10-10-2026)
  'UNGGAH_EDM.kolomAsing', // modal Upload CSV disederhanakan (keputusan work owner 10-10-2026)
  'UNGGAH_EDM.catatanBaris', // modal Upload CSV disederhanakan (keputusan work owner 10-10-2026)
  'UNGGAH_EDM.kolomMasalah', // modal Upload CSV disederhanakan (keputusan work owner 10-10-2026)
  'UNGGAH_EDM.petunjukTemplat', // modal Upload CSV disederhanakan (keputusan work owner 10-10-2026)
  'UNGGAH_EDM.downloadTemplate', // templat CSV (keputusan work owner 10-10-2026)
  'POLIS_LAMA_EDM.endorsementNo', // kolom Endorsement No rekap View Old Policy (keputusan work owner 10-10-2026)
  'POLIS_LAMA_EDM.newBusiness', // baris versi new business di rekap View Old Policy (keputusan work owner 10-10-2026)
  'BUAT_EDM.infoAkhiran', // Policy No bernomor endorsement (keputusan work owner 10-10-2026)
  'BUAT_EDM.versiTerakhir', // Policy No bernomor endorsement (keputusan work owner 10-10-2026)
  'BUAT_EDM.versiNB', // Policy No bernomor endorsement (keputusan work owner 10-10-2026)
  'INBOX_EDM.tutup', // kasus tertutup di kotak masuk tidak dapat dibuka (keputusan work owner 09-10-2026)
  'HASIL_EDM.judul', // popup hasil Confirm menggantikan ConfirmSubmitEDM (permintaan work owner 09-10-2026)
  'HASIL_EDM.noEndorsement', // popup hasil Confirm menggantikan ConfirmSubmitEDM (permintaan work owner 09-10-2026)
  'HASIL_EDM.edmType', // popup hasil Confirm menggantikan ConfirmSubmitEDM (permintaan work owner 09-10-2026)
  'HASIL_EDM.subType', // popup hasil Confirm menggantikan ConfirmSubmitEDM (permintaan work owner 09-10-2026)
  'HASIL_EDM.caseId', // popup hasil Confirm menggantikan ConfirmSubmitEDM (permintaan work owner 09-10-2026)
  'HASIL_EDM.ok', // popup hasil Confirm menggantikan ConfirmSubmitEDM (permintaan work owner 09-10-2026)
  'KONFIRMASI_PUTUSAN_EDM.confirm', // Submit b37494/b38109 diganti nama Confirm + konfirmasi (keputusan work owner 08-10-2026)
  'KONFIRMASI_PUTUSAN_EDM.judulAccept', // Submit b37494/b38109 diganti nama Confirm + konfirmasi (keputusan work owner 08-10-2026)
  'KONFIRMASI_PUTUSAN_EDM.isiAccept', // Submit b37494/b38109 diganti nama Confirm + konfirmasi (keputusan work owner 08-10-2026)
  'KONFIRMASI_PUTUSAN_EDM.judulDecline', // Submit b37494/b38109 diganti nama Confirm + konfirmasi (keputusan work owner 08-10-2026)
  'KONFIRMASI_PUTUSAN_EDM.isiDecline', // Submit b37494/b38109 diganti nama Confirm + konfirmasi (keputusan work owner 08-10-2026)
  'KONFIRMASI_PUTUSAN_EDM.batal', // Submit b37494/b38109 diganti nama Confirm + konfirmasi (keputusan work owner 08-10-2026)
  'KONFIRMASI_PUTUSAN_EDM.ya', // Submit b37494/b38109 diganti nama Confirm + konfirmasi (keputusan work owner 08-10-2026)
  'PERBANDINGAN_EDM.judul', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.judulRekap', // Summary Comparison dirapikan (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.judulKeputusan', // Summary Comparison dirapikan (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.kolomBaris', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.kolomPeserta', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.before', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.deleted', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.added', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.after', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.difference', // Summary Comparison (permintaan work owner 08-10-2026)
  'PERBANDINGAN_EDM.kosong', // Summary Comparison (permintaan work owner 08-10-2026)
  'PROSES_EDM.sudahDiproses', // Process Endorsement langkah final, CSV terkunci (keputusan work owner 08-10-2026)
  'PROSES_EDM.tanpaPerubahan', // proteksi Process Endorsement tanpa perubahan (keputusan work owner 08-10-2026)
  'PROSES_EDM.processEndorsement', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.memproses', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.judulKonfirmasi', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.hapusN', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.hapusNol', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.baruN', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.batalN', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.batalSemua', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.rekap', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.batal', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'PROSES_EDM.lanjut', // tombol Save b37202 diganti nama + konfirmasi (keputusan work owner 08-10-2026)
  'MENU_EDM.kelompok', // nama FOLDER korpus, dibuktikan terpisah di bawah
  'UMUM_EDM.kosong', // teks bawaan grid kosong
  'UMUM_EDM.tutup', // ikon tutup harness `pxIconCancel` tanpa teks
  'INBOX_EDM.kolomPolicyNo', // korpus b8964 "Policy No"; kini "PL_NUMBER" seperti kotak masuk PremiumList Life (keputusan work owner 07-10-2026) - isinya tetap OLD_POLICY_NO
  'INBOX_EDM.statusTerbuka', // sel Work Status kasus terbuka = nama tahap InputEDMLife (keputusan work owner 08-10-2026) - tampilan, STATUSS tetap kosong
  'INBOX_EDM.kolomStatus', // korpus b10052 "Status"; kini "Work Status" seperti kotak masuk PremiumList Life (keputusan work owner 07-10-2026) - isinya tetap STATUSS
  'KASUS_EDM.bagianPolis', // judul bagian layar kasus - tampilan (keputusan work owner 07-10-2026)
  'KASUS_EDM.bagianPihak', // judul bagian layar kasus - tampilan (keputusan work owner 07-10-2026)
  'KASUS_EDM.bagianEndorsement', // judul bagian layar kasus - tampilan (keputusan work owner 07-10-2026)
  'KASUS_EDM.peserta', // judul bagian peserta layar kasus - tampilan (keputusan work owner 07-10-2026)
  'RINCI_EDM.bagianPeserta', // judul bagian popup Details - tampilan (keputusan work owner 08-10-2026)
  'RINCI_EDM.bagianPeriode', // judul bagian popup Details - tampilan (keputusan work owner 08-10-2026)
  'RINCI_EDM.bagianPertanggungan', // judul bagian popup Details - tampilan (keputusan work owner 08-10-2026)
  'RINCI_EDM.bagianPremi', // judul bagian popup Details - tampilan (keputusan work owner 08-10-2026)
  'RINCI_EDM.bagianSpreading', // judul bagian popup Details - tampilan (keputusan work owner 08-10-2026)
  'UMUM_EDM.nomor', // kolom nomor baris popup View Old Policy - tampilan saja (keputusan work owner 07-10-2026)
  'BUAT_EDM.pilih', // teks pilihan kosong dropdown sub-jenis (layar Pega, diberikan work owner 07-10-2026); opsi properti tidak diekspor
  'BUAT_EDM.cancel', // footer Create Addendum - keputusan work owner 07-10-2026 (harness hanya ikon tutup)
  'BUAT_EDM.bantuan', // teks bantuan Submit nonaktif - keputusan work owner 07-10-2026 (Pega menyembunyikan Submit)
  'UMUM_EDM.memuat', // teks bawaan muat
  'UMUM_EDM.memeriksa', // gerbang sedang berjalan
  'UMUM_EDM.rinci', // tombol buka rincian baris (`pyEditingMode` `expandPane`, tanpa teks di korpus)
  'UMUM_EDM.terkunci', // pesan penjelas kunci field - spec §9 `[keputusan work owner + desain]`, AC 34
  'UMUM_EDM.menyimpan', // `Save` sedang berjalan
  'UMUM_EDM.tandaiHapus', // nama aksesibel kotak centang `.EdmBatal` b15753 (`pyLabelPreview` kosong)
  'UMUM_EDM.mengunggah', // unggahan sedang berjalan
  'UMUM_EDM.pilihBerkas', // isian berkas wizard `pxUploadCSVResults` (teks bawaan Pega, tidak diekspor)
  'UMUM_EDM.barisDisimpan', // hasil `Add CSV Data`
  'UMUM_EDM.pesanTerpotong', // server mengembalikan 200 penolakan pertama
  'UMUM_EDM.kolomBaris', // tabel penolakan (AC 37: nomor baris dan kolom)
  'UMUM_EDM.memutuskan', // keputusan sedang dikirim
  'UMUM_EDM.efekDilewati', // gerbang `IsPEGAPROD` b7230 - efek keluar hanya di produksi (tiket 10)
  'UMUM_EDM.efekGagal', // alarm di layar menggantikan `SendEmailNotification` 15 `//` (E5)
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

  it('kepala Generate Data Detail = CSVPropHeaders b276, nama berkas b270 + stempel b272', () => {
    expect(baca(UNDUH_CSV, 276)).toBe(`<CSVPropHeaders>${KEPALA_UNDUH_CSV.join(',')}</CSVPropHeaders>`)
    expect(baca(UNDUH_CSV, 270)).toBe('<FileName>DetailUpload</FileName>')
    expect(baca(UNDUH_CSV, 272)).toBe('<AppendTimeStampToFileName>true</AppendTimeStampToFileName>')
  })

  it('kedua Submit berlabel sama (b37494 Confirm, b38109 Decline); opsi radio dari AddHistorySuggest b353', () => {
    expect(baca(INPUT, 38109)).toBe('<pyLabel>Submit</pyLabel>')
    const b353 = baca(RIWAYAT, 353)
    expect(b353).toContain(`.EmailTypePL=1,"${LABEL.OPSI_KEPUTUSAN['1']}","${LABEL.OPSI_KEPUTUSAN['2']}"`)
    expect(Object.keys(LABEL.OPSI_KEPUTUSAN)).toEqual(['1', '2'])
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
