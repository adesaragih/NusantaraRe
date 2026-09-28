import { describe, expect, it } from 'vitest'

import {
  KONFIRMASI_BALIK,
  TOMBOL_AKSEPTASI,
  TOMBOL_MEDIS,
  TOMBOL_OS,
} from '../../assets/labels.claimlife'
import { TAHAP_JALUR } from '../../services/api'
import { tombolPindahTahap } from './PanelPindahTahap'

// Uji tombol perpindahan per layar — kelompok Medis dan Akseptasi.
//
// ⛔ Yang dijaga: tiap tahap menawarkan PERSIS perpindahan yang section-nya
// punya. Tombol yang disalin antarlayar membuka jalur yang di sistem lama
// tidak ada — dan backend akan menolaknya 409, sehingga yang lahir hanya
// tombol yang selalu gagal.

describe('tombolPindahTahap', () => {
  it('Outstanding menawarkan Register dan Medical Check', () => {
    // `InputOSClaimLife.xml` b21404, b21839.
    expect(tombolPindahTahap('Outstanding Claim')).toEqual([
      {
        label: TOMBOL_OS.kembaliKeRegister,
        tujuan: TAHAP_JALUR.inputRegister,
        konfirmasi: KONFIRMASI_BALIK.keAdmin,
      },
      { label: TOMBOL_OS.kirimKeMedis, tujuan: TAHAP_JALUR.medicalCheck },
    ])
  })

  it('Medical Check menawarkan Admin dan Claim Analyst', () => {
    // `MedicalCheckClaimLife.xml` b20256, b21151.
    expect(tombolPindahTahap('Medical Check')).toEqual([
      {
        label: TOMBOL_MEDIS.kembaliKeAdmin,
        tujuan: TAHAP_JALUR.outstanding,
        konfirmasi: KONFIRMASI_BALIK.keAdmin,
      },
      { label: TOMBOL_MEDIS.kirimKeAnalis, tujuan: TAHAP_JALUR.claimAnalis },
    ])
  })

  it('Claim Analis menawarkan Admin dan Medical', () => {
    // `InputAkseptasiClaimLife.xml` b20221, b20467.
    expect(tombolPindahTahap('Claim Analis')).toEqual([
      {
        label: TOMBOL_AKSEPTASI.kembaliKeAdmin,
        tujuan: TAHAP_JALUR.outstanding,
        konfirmasi: KONFIRMASI_BALIK.keAdmin,
      },
      {
        label: TOMBOL_AKSEPTASI.kembaliKeMedis,
        tujuan: TAHAP_JALUR.medicalCheck,
        konfirmasi: KONFIRMASI_BALIK.keMedis,
      },
    ])
  })

  it('Input Register dan tahap tak dikenal: KOSONG, gagal tertutup', () => {
    // ⛔ Satu-satunya perpindahan Input Register terjadi lewat `Submit`
    // pendaftaran, bukan tombol perpindahan. Dan tahap kosong — klaim yang
    // baris work-nya tidak ada — tidak menawarkan apa pun.
    expect(tombolPindahTahap('Input Register')).toEqual([])
    expect(tombolPindahTahap('')).toEqual([])
    expect(tombolPindahTahap('Selesai')).toEqual([])
  })

  it('TIDAK ADA tombol yang menunjuk tahapnya sendiri', () => {
    // Perpindahan ke tahap asal bukan perpindahan; backend menolaknya 409.
    const peta: Record<string, string> = {
      'Outstanding Claim': TAHAP_JALUR.outstanding,
      'Medical Check': TAHAP_JALUR.medicalCheck,
      'Claim Analis': TAHAP_JALUR.claimAnalis,
    }
    for (const [tahap, jalur] of Object.entries(peta)) {
      for (const t of tombolPindahTahap(tahap)) {
        expect(t.tujuan, `${tahap} menawarkan perpindahan ke dirinya sendiri`).not.toBe(
          jalur,
        )
      }
    }
  })

  it('tiap layar menawarkan TEPAT dua, dan labelnya tidak tertukar', () => {
    // ⛔ Cacahnya dikunci di kedua arah. `Send to Medical Check` milik
    // Outstanding; `Send Back to Medical` milik Claim Analis. Keduanya
    // menuju tahap yang sama dengan kalimat yang berbeda, dan menukarnya
    // tidak akan membuat satu pun uji perilaku merah.
    for (const tahap of ['Outstanding Claim', 'Medical Check', 'Claim Analis']) {
      expect(tombolPindahTahap(tahap)).toHaveLength(2)
    }
    expect(TOMBOL_OS.kirimKeMedis).toBe('Send to Medical Check')
    expect(TOMBOL_AKSEPTASI.kembaliKeMedis).toBe('Send Back to Medical')
    expect(TOMBOL_MEDIS.kirimKeAnalis).toBe('Send to Claim Analyst')
    expect(TOMBOL_MEDIS.kembaliKeAdmin).toBe('Send Back to Admin')
    expect(TOMBOL_AKSEPTASI.kembaliKeAdmin).toBe('Send Back to Admin')
  })
})

describe('konfirmasi jalur balik', () => {
  it('tombol maju TIDAK bertanya — activity langsung + finishAssignment', () => {
    // `Send to Medical Check` b21863/b21891, `Send to Claim Analyst` b21174/b21202.
    const maju = [
      ...tombolPindahTahap('Outstanding Claim'),
      ...tombolPindahTahap('Medical Check'),
    ].filter((t) => !t.label.startsWith('Send Back'))
    expect(maju).toHaveLength(2)
    for (const t of maju) expect(t.konfirmasi).toBeUndefined()
  })

  it('setiap tombol jalur balik bertanya', () => {
    const balik = ['Outstanding Claim', 'Medical Check', 'Claim Analis']
      .flatMap((tahap) => tombolPindahTahap(tahap))
      .filter((t) => t.label.startsWith('Send Back'))
    expect(balik).toHaveLength(4)
    for (const t of balik) expect(t.konfirmasi).toBeDefined()
  })
})
