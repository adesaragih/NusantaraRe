// Uji layar tahun treaty — tiket 03 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { TAHUN_TCO } from '../labels'
import type { TahunTreaty } from '../api'
import { formDari, formKosong, isiDariMulai, keMasuk, namaGrup, selTahun, tahunDari } from './InboxTreatyContract'

const SUMBER = readFileSync(join(__dirname, 'InboxTreatyContract.tsx'), 'utf8')
const KODE = SUMBER.split('\n')
  .filter((b) => !b.trimStart().startsWith('//'))
  .join('\n')

const BARIS: TahunTreaty = {
  id: '1000001', treatyYear: '2026', underwritingYear: '2026', treatyGroupId: '10001',
  treatyGroupName: 'UJI GRUP', proportion: '10003', startDate: '2026-01-01', endDate: '2026-12-31',
  userId: 'UJI-OP', tglUpdate: '2026-09-28 10:00:00',
}

describe('form tahun treaty', () => {
  it('baris baru: ID kosong — identitas dibuat server (AC 5)', () => {
    const f = formKosong()
    expect(f.id).toBe('')
    expect(keMasuk(f).id).toBe('')
  })
  it('Edit menyalin seluruh medan baris (SetTreatyYear_Act b1110–b1353)', () => {
    const f = formDari(BARIS)
    expect(f).toMatchObject({ id: '1000001', treatyGroupName: 'UJI GRUP', proportion: '10003', startDate: '2026-01-01' })
  })
  it('tanggal dikirim YYYY-MM-DD dari bentuk kabel DD-MM-YYYY maupun ISO', () => {
    const f = { ...formDari(BARIS), startDate: '01-01-2026', endDate: '2026-12-31' }
    const m = keMasuk(f)
    expect(m.startDate).toBe('2026-01-01')
    expect(m.endDate).toBe('2026-12-31')
  })
  it('teks yang bukan tanggal dikirim APA ADANYA — server yang menolak dengan pesan', () => {
    expect(keMasuk({ ...formKosong(), startDate: 'kapan-kapan' }).startDate).toBe('kapan-kapan')
    expect(keMasuk({ ...formKosong(), startDate: '  ' }).startDate).toBe('')
  })
  it('tahun dikirim sebagai TEKS, di-trim, tidak lewat Number()', () => {
    expect(keMasuk({ ...formKosong(), treatyYear: ' 2026 ' }).treatyYear).toBe('2026')
    expect(KODE).not.toMatch(/Number\(|parseInt|parseFloat/)
  })
  it('sel kosong ditandai', () => {
    expect(selTahun('')).toBe('—')
    expect(selTahun('2026')).toBe('2026')
  })
  it('nama grup dari master, kosong bila tidak dikenal', () => {
    const grup = [{ id: '10001', treatyGroupName: 'UJI GRUP' }]
    expect(namaGrup(grup, '10001')).toBe('UJI GRUP')
    expect(namaGrup(grup, '99')).toBe('')
  })
})

describe('paritas layar', () => {
  it('enam kepala kolom grid dan empat tombol dari label, tidak diketik ulang', () => {
    for (const k of ['kolomUnderwritingYear', 'kolomTransactionYear', 'kolomStartDate', 'kolomEndDate',
      'kolomTreatyGroup', 'kolomReinsuranceType', 'add', 'edit', 'reinsType', 'listDescription',
      'formId', 'formTreatyGroup', 'formReinsuranceType', 'formStartDate', 'formEndDate',
      'formUnderwritingYear', 'formTransactionYear', 'formModifiedDate', 'formUsername', 'save', 'cancel'] as const) {
      expect(KODE, k).toContain(`TAHUN_TCO.${k}`)
      expect(KODE).not.toContain(`'${TAHUN_TCO[k]}'`)
    }
  })
  it('kolom Reinsurance Type grid menampilkan proportion (b19724)', () => {
    expect(KODE).toContain('selTahun(labelProporsi(b.proportion))')
  })
  it('nol tombol Copy dan nol form salin (AC 72)', () => {
    expect(KODE).not.toMatch(/Copy|salin|BrowseCopyData/)
  })
  it('ReinsType membuka editor kontrak (tiket 04); List Description membuka layar klausul (tiket 08)', () => {
    expect(KODE).toContain('TAHUN_TCO.reinsType')
    expect(KODE).toContain("bukaRinci({ jenis: 'kontrak', tahun: b })")
    expect(KODE).toContain('<PanelKontrakTahun')
    expect(KODE).toContain('TAHUN_TCO.listDescription')
    expect(KODE).toContain("bukaRinci({ jenis: 'klausul', tahun: b })")
    expect(KODE).toContain('<PanelKlausulTahun')
    expect(KODE).not.toMatch(/`\$\{TAHUN_TCO.menungguTiket\} 0[48]`/)
  })
  it('ID, Modified Date, Username hanya dibaca', () => {
    expect((KODE.match(/readOnly/g) ?? []).length).toBe(3)
  })
  it('layar tidak menyaring daftar; Reinsurance Type dari daftar Proportional/NonProportional (OQ-TCO-04)', () => {
    expect(KODE).not.toContain('.filter(')
    expect(KODE).toContain("from '../proporsi'")
  })
  it('tiket 12: panel lampiran hanya untuk tahun yang sudah ber-ID, bukan syarat simpan', () => {
    expect(KODE).toContain('<PanelLampiranTahun tahunID={form.id} />')
    expect(KODE).toMatch(/form\.id !== '' &&/)
    // [keputusan work owner 30-09-2026] catatan "simpan dulu" tidak tampil.
    expect(KODE).not.toContain('simpanDulu')
    // Simpan tahun tidak menunggu lampiran: fungsi simpan tidak menyebut lampiran.
    const awal = KODE.indexOf('async function simpan')
    const simpan = KODE.slice(awal, KODE.indexOf('\n  }\n', awal))
    expect(simpan).toContain('simpanTahunTreaty(')
    expect(simpan).not.toMatch(/[Ll]ampiran/)
  })

  it('panel rinci SATU saja, dan selama terbuka tabel utama tidak dirender (keputusan work owner 30-09-2026)', () => {
    expect(KODE).toContain('const [rinci, setRinci] = useState<RinciTahun | null>(null)')
    expect(KODE).not.toMatch(/setTahunKontrak|setTahunKlausul/)
    // Kembalian dini SEBELUM tabel: tabel, form, dan penomoran tidak ikut.
    const dini = KODE.indexOf('if (rinci !== null) {')
    expect(dini).toBeGreaterThan(0)
    expect(dini).toBeLessThan(KODE.indexOf('<table className="inbox__tabel">'))
    expect(KODE.slice(dini, KODE.indexOf('\n  }\n', dini))).not.toMatch(/<table|<Halaman|form !== null/)
  })
  it('nol catatan pengembang di layar (keputusan work owner 30-09-2026)', () => {
    expect(KODE).not.toMatch(/catatanLabelBersilang|polis__catatan/)
  })
})

describe('Add tahun treaty: Start Date mengisi tahun dan End Date (keputusan work owner 30-09-2026)', () => {
  it('tahun diambil dari Start Date', () => {
    expect(tahunDari('2026-06-01')).toBe('2026')
    expect(tahunDari('01-06-2026')).toBe('2026')
    expect(tahunDari('bukan tanggal')).toBe('')
  })
  it('form BARU: Underwriting Year dan Transaction Year = tahun Start Date', () => {
    const f = isiDariMulai(formKosong(), '2026-06-01')
    expect(f.startDate).toBe('2026-06-01')
    // `treatyYear` berlabel Underwriting Year, `underwritingYear` berlabel Transaction Year (OQ-TCO-05).
    expect(f.treatyYear).toBe('2026')
    expect(f.underwritingYear).toBe('2026')
  })
  it('tetap dapat diubah: tahun yang diketik sesudahnya tidak ditimpa sampai Start Date berubah lagi', () => {
    const f = { ...isiDariMulai(formKosong(), '2026-06-01'), treatyYear: '2027' }
    expect(f.treatyYear).toBe('2027')
    expect(isiDariMulai(f, '2028-01-01').treatyYear).toBe('2028')
  })
  it('form yang sudah ber-ID tidak diisi ulang', () => {
    const lama = { ...formKosong(), id: '1000001', treatyYear: '2020', underwritingYear: '2019' }
    const f = isiDariMulai(lama, '2026-06-01')
    expect(f.treatyYear).toBe('2020')
    expect(f.underwritingYear).toBe('2019')
  })
  it('End Date dari server - aturan yang sama dengan kontrak; hanya form baru', () => {
    expect(KODE).toContain('onChange={ubahMulai}')
    expect(KODE).toContain('ambilAkhirBawaanTahun(iso)')
    expect(KODE).toMatch(/form\.id !== '' \|\| iso === ''\) return/)
  })
})

