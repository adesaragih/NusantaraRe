// Uji editor kontrak treaty — tiket 04 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { alihRinci, formKontrakDari, formKontrakKosong, keMasukKontrak } from './PanelKontrakTahun'

const KODE = readFileSync(join(__dirname, 'PanelKontrakTahun.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
  .join('\n')

describe('form kontrak', () => {
  it('baris baru: NewInputTreatyContract_Act mengosongkan ID, jenis, dan kedua tanggal', () => {
    expect(formKontrakKosong()).toEqual({
      id: '', reinsTypeId: '', treatyStartDate: '', treatyEndDate: '', tglUpdate: '', userId: '',
    })
  })
  it('Edit: SetUbahTreatyContract menyalin ID, jenis, dan kedua tanggal dari baris', () => {
    const f = formKontrakDari({
      id: '1000003', idTreatyYear: '1000001', reinsTypeId: '10003', reinsTypeName: 'UJI QS',
      treatyStartDate: '2026-01-01', treatyEndDate: '2027-01-01', userId: 'UJI-ADMIN', tglUpdate: '2026-09-29 10:00:00',
    })
    expect(f).toEqual({
      id: '1000003', reinsTypeId: '10003', treatyStartDate: '2026-01-01', treatyEndDate: '2027-01-01',
      tglUpdate: '2026-09-29 10:00:00', userId: 'UJI-ADMIN',
    })
  })
  it('tanggal dikirim YYYY-MM-DD; teks bukan tanggal dikirim apa adanya untuk ditolak server', () => {
    const m = keMasukKontrak({ ...formKontrakKosong(), reinsTypeId: '10003', treatyStartDate: '01-02-2026',
      treatyEndDate: 'kapan-kapan' })
    expect(m).toEqual({ id: '', reinsTypeId: '10003', treatyStartDate: '2026-02-01', treatyEndDate: 'kapan-kapan' })
  })
})

describe('paritas layar kontrak', () => {
  it('seluruh label dari KONTRAK_TCO, tidak diketik ulang', () => {
    for (const k of ['judul', 'headerUnderwritingYear', 'headerReinsType', 'formStartDate', 'formEndDate', 'save',
      'formModifiedDate', 'formUsername', 'undo', 'information', 'add', 'kolomReinsType', 'kolomTreatyStart',
      'kolomTreatyEnd', 'edit', 'businessList', 'reinsurerList', 'delete'] as const) {
      expect(KODE).toContain(`KONTRAK_TCO.${k}`)
    }
  })
  it('jenis reasuransi dari pemilih tersaring tiket 02, bukan teks bebas', () => {
    expect(KODE).toContain('<PilihJenisReasuransi')
    expect(KODE).not.toMatch(/reinsTypeName:\s*e\.target/)
  })
  it('tanggal mulai mengisi tanggal akhir bawaan dari server; tanggal akhir hanya dirinya (IsEndDate)', () => {
    const mulai = KODE.slice(KODE.indexOf('function ubahMulai'), KODE.indexOf('function ubahAkhir'))
    expect(mulai).toContain('ambilAkhirBawaanKontrak(')
    const akhir = KODE.slice(KODE.indexOf('function ubahAkhir'))
    expect(akhir.slice(0, akhir.indexOf('\n  }\n'))).not.toContain('ambilAkhirBawaanKontrak(')
  })
  it('Undo mengembalikan isian terakhir yang dimuat', () => {
    expect(KODE).toMatch(/setForm\(asal\)/)
  })
  it('Reinsurer List (05) dan Business List (07) membuka panelnya; Delete menunggu 10', () => {
    expect(KODE).toContain("setRinci((r) => alihRinci(r, 'reinsurer', k.id))")
    expect(KODE).toContain('<PanelReinsurerKombinasi')
    expect(KODE).toContain("setRinci((r) => alihRinci(r, 'business', k.id))")
    expect(KODE).toContain('<PanelBusinessKombinasi')
    expect(KODE).not.toContain('`${TAHUN_TCO.menungguTiket} 05`')
    expect(KODE).not.toContain('`${TAHUN_TCO.menungguTiket} 07`')
    expect(KODE).not.toContain('`${TAHUN_TCO.menungguTiket} 10`')
    expect(KODE).toContain('mintaHapus(k)')
  })
  it('ID, Modified Date, Username hanya dibaca', () => {
    expect(KODE).toMatch(/label=\{KONTRAK_TCO\.formId\}[^/]*readOnly/)
    expect(KODE).toMatch(/label=\{KONTRAK_TCO\.formModifiedDate\}[^/]*readOnly/)
    expect(KODE).toMatch(/label=\{KONTRAK_TCO\.formUsername\}[^/]*readOnly/)
  })
})

describe('Business List / Reinsurer List — satu panel rinci (keputusan work owner 30-09-2026)', () => {
  it('membuka daftar kontrak lain MENUTUP yang terbuka', () => {
    const b1 = alihRinci(null, 'business', 'K1')
    expect(b1).toEqual({ daftar: 'business', kontrakID: 'K1' })
    // Business List kontrak ke-2 dibuka: Business List pertama tertutup.
    expect(alihRinci(b1, 'business', 'K2')).toEqual({ daftar: 'business', kontrakID: 'K2' })
    // Begitu pula Reinsurer List: satu saja, termasuk lintas jenis daftar.
    const r1 = alihRinci(null, 'reinsurer', 'K1')
    expect(alihRinci(r1, 'reinsurer', 'K2')).toEqual({ daftar: 'reinsurer', kontrakID: 'K2' })
    expect(alihRinci(b1, 'reinsurer', 'K1')).toEqual({ daftar: 'reinsurer', kontrakID: 'K1' })
  })
  it('menekan tombol yang sama menutupnya', () => {
    expect(alihRinci({ daftar: 'business', kontrakID: 'K1' }, 'business', 'K1')).toBeNull()
  })
  it('SATU keadaan, dan panelnya dirender di bawah baris kontraknya', () => {
    expect(KODE).toContain('const [rinci, setRinci] = useState<RinciKontrak | null>(null)')
    expect(KODE).not.toMatch(/setKontrakBusiness|setKontrakReinsurer/)
    expect(KODE).toContain('{rinci?.kontrakID === k.id && (')
    expect(KODE).toContain('<tr className="inbox__rinci">')
    expect(KODE).toContain('<td colSpan={4}>')
  })
  it('nol catatan pengembang di layar', () => {
    expect(KODE).not.toMatch(/catatanLabelBersilang|polis__catatan/)
  })
  it('baris rinci TIDAK mewarisi nowrap sel tabel (temuan /code-review 30-09-2026)', () => {
    // `.inbox__tabel td { white-space: nowrap }` diwarisi isi sel: tanpa aturan
    // ini seluruh kalimat panel Business/Reinsurer/Security (dan popup hapus di
    // dalamnya) tidak membungkus.
    const css = readFileSync(join(__dirname, '..', '..', '..', 'inti', 'styles.css'), 'utf8')
    expect(css).toMatch(/\.inbox__tabel tr\.inbox__rinci > td \{\s*white-space: normal;/)
  })
})
