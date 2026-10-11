// Uji panel data polis — A3 kelompok Register.

import { describe, expect, it } from 'vitest'

import { REGISTER } from '../labels'
import { medanPolis } from './PanelDataPolis'

describe('himpunan medan PERSIS section', () => {
  const medan = medanPolis({})

  // ⛔ 10-10-2026 (keputusan work owner): System Reinsurance dan Premium Method TIDAK ditampilkan - cacah turun 20 -> 18.
  it('delapan belas medan bagi non-treaty (tanpa dua medan retro, tanpa System Reinsurance / Premium Method)', () => {
    // ⛔ Medan yang DIHILANGKAN dari layar sama merusaknya dengan medan yang
    // dikarang - dan yang pertama tidak berbunyi sama sekali. Cacahnya
    // dikunci supaya menghilangkan satu menjadi kegagalan uji.
    expect(medan).toHaveLength(18)
  })

  it('setiap label datang dari REGISTER, bukan diketik ulang', () => {
    const sah = new Set<string>(Object.values(REGISTER))
    for (const m of medan) {
      expect(sah.has(m.label), `label ${m.label}`).toBe(true)
    }
  })

  it('label tidak ganda', () => {
    const lihat = new Set(medan.map((m) => m.label))
    expect(lihat.size).toBe(medan.length)
  })

  it('yang MENUNGGU MODUL dinyatakan, bukan dihilangkan', () => {
    const belum = medan.filter((m) => m.belumBersumber === true)
    // ⛔ SEBELAS, bukan sepuluh: `Type` ikut menunggu sejak ralat
    // pembacaan 27-09-2026 - ia `pyReadOnly` true di Pega dan terikat
    // `.PolicyDataLife.Type`, bukan isian bebas.
    expect(belum).toHaveLength(17) // 19 - System Reinsurance - Premium Method (10-10-2026)
    for (const m of belum) {
      expect(m.nilai).toBe('—')
    }
  })

  it('yang PUNYA sumber terisi dari nilainya', () => {
    const isi = medanPolis({ namaTertanggung: 'UJI-TERTANGGUNG' })
    const nama = isi.find((m) => m.label === REGISTER.namaTertanggung)
    expect(nama?.nilai).toBe('UJI-TERTANGGUNG')
    expect(nama?.belumBersumber).toBeUndefined()
  })

  it('Type TETAP menunggu walau nilainya diberikan', () => {
    // ⛔ Ia read-only di Pega; nilai yang layar punya hari ini datang
    // dari isian sementara kita, BUKAN dari PolicyDataLife.
    const isi = medanPolis({ type: 'QP' })
    const type = isi.find((m) => m.label === REGISTER.type)
    expect(type?.belumBersumber).toBe(true)
  })

  it('nilai kosong tampil sebagai tanda pisah, bukan string kosong', () => {
    // Sel kosong tidak dapat dibedakan dari sel yang gagal dimuat.
    const isi = medanPolis({ namaTertanggung: '' })
    expect(isi.find((m) => m.label === REGISTER.namaTertanggung)?.nilai).toBe('—')
  })
})

// ——— Butir av: PolicyDataLife dari PremiumList Life, 28-09-2026 ———

const POLIS_UJI = {
  nomorPolis: 'UJI-POL-1',
  type: 'QR',
  marketingName: 'UJI MARKETING',
  cedingCoName: 'UJI CEDING',
  policyHolderName: 'UJI HOLDER',
  businessName: 'UJI COB',
  dateReceived: '2026-01-15',
  status: 'UJI-STATUS',
  statusUpdate: 'UJI-STATUS-BARU',
  productNameId: 'UJI-PROD-1',
  productName: 'UJI PRODUK',
  prodKe: 3,
  medanTanpaSumber: [
    'TanggalRespon',
    'TanggalKonfirmasi',
    'TanggalRealisasi',
    'TanggalKonfirmasiBalik',
    'KetentuanUnderwriting',
  ],
  typeCeding: '1',
  typeCedingName: 'UJI QS',
  proRateType: 'UJI-PRORATA',
  wpc: '2026-02-01',
  retroName: 'UJI RETRO',
  securityReinsurer: 'UJI SECURITY',
  sobName: 'UJI SOB',
}

describe('butir av — delapan medan bersumber, tiga tidak', () => {
  it('delapan medan terisi dari polis', () => {
    const medan = medanPolis({ polis: POLIS_UJI, namaTertanggung: 'UJI PESERTA' })
    const isi = new Map(medan.map((m) => [m.label, m.nilai]))
    expect(isi.get(REGISTER.type)).toBe('QR')
    expect(isi.get(REGISTER.marketing)).toBe('UJI MARKETING')
    expect(isi.get(REGISTER.ceding)).toBe('UJI CEDING')
    expect(isi.get(REGISTER.pemegangPolis)).toBe('UJI HOLDER')
    expect(isi.get(REGISTER.kelasBisnis)).toBe('UJI COB')
    expect(isi.get(REGISTER.tanggalEmail)).toBe('15/01/2026') // dd/mm/yyyy sejak 10-10-2026
    expect(isi.get(REGISTER.status)).toBe('UJI-STATUS')
    expect(isi.get(REGISTER.statusDiperbarui)).toBe('UJI-STATUS-BARU')
  })

  it('TEPAT lima medan tetap tanpa sumber, dan mereka bernama', () => {
    // ⛔ Kelimanya TIDAK PUNYA KOLOM di migrasi 050–056 mana pun — bukan
    // "menunggu modul", sebab modulnya sudah ada.
    const medan = medanPolis({ polis: POLIS_UJI })
    const belum = medan.filter((m) => m.belumBersumber === true).map((m) => m.label)
    expect(belum).toEqual([
      REGISTER.tanggalRespon,
      REGISTER.tanggalKonfirmasi,
      REGISTER.konfirmasiBalik,
      REGISTER.tanggalRealisasi,
      REGISTER.catatanUnderwriter,
    ])
  })

  it('himpunan medannya TIDAK bergantung isi — delapan belas bagi non-treaty (10-10-2026)', () => {
    // ⛔ Paritas dengan `InputRegisterClaimLife.xml` diukur dari HIMPUNAN
    // medannya, bukan dari berapa yang terisi.
    expect(medanPolis({ polis: POLIS_UJI })).toHaveLength(18)
    expect(medanPolis({})).toHaveLength(18)
  })

  it('av-2 — medan terisi dari polis; System Reinsurance / Premium Method tidak ditampilkan (10-10-2026)', () => {
    const isi = new Map(medanPolis({ polis: POLIS_UJI }).map((m) => [m.label, m.nilai]))
    expect(isi.has(REGISTER.sistemReasuransi)).toBe(false)
    expect(isi.has(REGISTER.metodePremi)).toBe(false)
    expect(isi.get(REGISTER.wpc)).toBe('01/02/2026') // dd/mm/yyyy
    expect(isi.get(REGISTER.sob)).toBe('UJI SOB')
    expect(isi.get(REGISTER.idProduk)).toBe('UJI-PROD-1')
    expect(isi.get(REGISTER.namaProduk)).toBe('UJI PRODUK')
  })

  it('av-2 — medan retro HANYA bagi TP/TR (pyCondition b10541/b10824)', () => {
    const treaty = medanPolis({ polis: { ...POLIS_UJI, type: 'TP' } })
    const isi = new Map(treaty.map((m) => [m.label, m.nilai]))
    expect(treaty).toHaveLength(20) // 22 - System Reinsurance - Premium Method (10-10-2026)
    expect(isi.get(REGISTER.retroName)).toBe('UJI RETRO')
    expect(isi.get(REGISTER.securityReinsurer)).toBe('UJI SECURITY')
    const labelQR = medanPolis({ polis: POLIS_UJI }).map((m) => m.label)
    expect(labelQR).not.toContain(REGISTER.retroName)
    expect(labelQR).not.toContain(REGISTER.securityReinsurer)
  })

  it('tanggal tak terurai / kosong tampil tanda pisah (dd/mm/yyyy, 10-10-2026)', () => {
    const isi = new Map(medanPolis({ polis: { ...POLIS_UJI, wpc: '', dateReceived: 'bukan-tanggal' } }).map((m) => [m.label, m.nilai]))
    expect(isi.get(REGISTER.wpc)).toBe('—')
    expect(isi.get(REGISTER.tanggalEmail)).toBe('—')
  })

  it('tanpa polis, Type jatuh kembali ke pilihan pemakai', () => {
    // ⚠️ Klaim dapat didaftarkan atas polis yang belum ada di PremiumList
    // Life; layar tetap berguna.
    const medan = medanPolis({ type: 'TP' })
    const isi = new Map(medan.map((m) => [m.label, m.nilai]))
    expect(isi.get(REGISTER.type)).toBe('TP')
  })

  it('kolom polis yang KOSONG ditandai, bukan dibiarkan kosong', () => {
    const medan = medanPolis({ polis: { ...POLIS_UJI, cedingCoName: '' } })
    const isi = new Map(medan.map((m) => [m.label, m.nilai]))
    expect(isi.get(REGISTER.ceding)).toBe('—')
  })
})
