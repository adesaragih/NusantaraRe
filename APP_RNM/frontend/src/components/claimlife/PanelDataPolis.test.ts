// Uji panel data polis — A3 kelompok Register.

import { describe, expect, it } from 'vitest'

import { REGISTER } from '../../assets/labels.claimlife'
import { medanPolis } from './PanelDataPolis'

describe('himpunan medan PERSIS section', () => {
  const medan = medanPolis({})

  it('sebelas medan data, tidak kurang', () => {
    // ⛔ Medan yang DIHILANGKAN dari layar sama merusaknya dengan medan yang
    // dikarang - dan yang pertama tidak berbunyi sama sekali. Cacahnya
    // dikunci supaya menghilangkan satu menjadi kegagalan uji.
    expect(medan).toHaveLength(12)
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
    expect(belum).toHaveLength(11)
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
  medanTanpaSumber: ['TanggalRespon', 'TanggalKonfirmasi', 'TanggalRealisasi'],
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
    expect(isi.get(REGISTER.tanggalEmail)).toBe('2026-01-15')
    expect(isi.get(REGISTER.status)).toBe('UJI-STATUS')
    expect(isi.get(REGISTER.statusDiperbarui)).toBe('UJI-STATUS-BARU')
  })

  it('TEPAT tiga medan tetap tanpa sumber, dan mereka bernama', () => {
    // ⛔ Ketiganya TIDAK PUNYA KOLOM di migrasi 050–056 mana pun — bukan
    // "menunggu modul", sebab modulnya sudah ada.
    const medan = medanPolis({ polis: POLIS_UJI })
    const belum = medan.filter((m) => m.belumBersumber === true).map((m) => m.label)
    expect(belum).toEqual([
      REGISTER.tanggalRespon,
      REGISTER.tanggalKonfirmasi,
      REGISTER.tanggalRealisasi,
    ])
  })

  it('himpunan medannya TIDAK berubah — tetap sebelas', () => {
    // ⛔ Paritas dengan `InputRegisterClaimLife.xml` diukur dari HIMPUNAN
    // medannya, bukan dari berapa yang terisi.
    expect(medanPolis({ polis: POLIS_UJI })).toHaveLength(12)
    expect(medanPolis({})).toHaveLength(12)
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
