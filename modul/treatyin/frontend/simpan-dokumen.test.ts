// Penyusun isi tombol tulis — Save/Submit/Actions/Decline offer.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { BarisKursWarisan, LimitsAkar, ShareNP } from './api'
import { bolehActions, bolehSave, keYYYYMMDD, kursBerubah, susunDokumen, type KepalaForm } from './simpanDokumen'

const KEPALA: KepalaForm = {
  nonProporsional: true,
  nama: 'KONTRAK UJI',
  rujukan: 'REF-1',
  wilayah: 'Indonesia',
  bordereaux: 'reporting',
  bordereauxNote: 'catatan',
  mulai: '01-07-2025',
  berakhir: '30-06-2026',
  tahunTreaty: '2025',
  pembukuan: '',
  pembukuanNonProp: 'loss',
  cedant: 'ASURANSI A',
  idCedant: '10001',
  asalBisnis: 'DIRECT',
  idAsalBisnis: '20001',
  pemimpin: true,
}

describe('susunDokumen — isi layar → properti TreatyIn ejaan Pega', () => {
  it('kepala form: nilai TERSIMPAN (NonProportional, YYYYMMDD, "true")', () => {
    const d = susunDokumen(KEPALA, {}, null, null)
    expect(d).toMatchObject({
      ProportionType: 'NonProportional',
      TreatyContractName: 'KONTRAK UJI',
      Commencement: '20250701',
      Termination: '20260630',
      Ceding: 'ASURANSI A',
      CedingID: '10001',
      LeadingReinsSource: 'DIRECT',
      LeadingReinsSourceID: '20001',
      TreatyLeader: 'true',
      AccountingModeNonProp: 'loss',
    })
  })

  it('penampung halaman ikut, tetapi kepala form MENANG atas penampung', () => {
    const d = susunDokumen(KEPALA, { Portfolio: [{ Description: 'x' }], TreatyContractName: 'lama' }, null, null)
    expect(d.Portfolio).toEqual([{ Description: 'x' }])
    expect(d.TreatyContractName).toBe('KONTRAK UJI')
  })

  it('⛔ status, posisi, dan riwayat TIDAK disusun layar', () => {
    const d = susunDokumen(KEPALA, {}, null, null)
    for (const k of ['StatusAkseptasi', 'Position', 'PositionUsername', 'CommentList']) expect(d).not.toHaveProperty(k)
  })

  it('Limits / Share Non-Prop diratakan ke kunci akar — total per nama properti', () => {
    const akar: LimitsAkar = {
      LimitSummaryList: [{ Layer: '1' }],
      Total: { TotalLimitIOONP: [{ Currency: 'IDR', Value: '5' }] },
      TotalLimitsROL: '12',
    }
    const share = { RNMShare: '10', Share: [], Total: { TotalShareNP: [] } } as unknown as ShareNP
    const d = susunDokumen(KEPALA, {}, { layers: [{ Layer: '1' }], akar }, share)
    expect(d.Limits).toEqual([{ Layer: '1' }])
    expect(d.TotalLimitIOONP).toEqual([{ Currency: 'IDR', Value: '5' }])
    expect(d.TotalLimitsROL).toBe('12')
    expect(d.RNMShare).toBe('10')
    expect(d).toHaveProperty('TotalShareNP')
    expect(d).not.toHaveProperty('Total')
  })
})

describe('kursBerubah — seluruh grid; yang tidak berubah bertanda tetap', () => {
  const asli: BarisKursWarisan[] = [
    {
      id: '13358',
      mataUang: 'USD',
      mataUangID: '10001',
      nilaiKeIDR: '16500',
      berlakuDari: '01/07/25',
      berlakuSampai: '30/06/26',
      berlakuDariAsli: '20250701',
      berlakuSampaiAsli: '20260630',
      tetap: true, // tercatat MILIK kontrak (`T_TREATY_KURS`)
    },
  ]

  it('tidak ada perubahan → null (nol tulisan ke tabel bersama)', () => {
    expect(kursBerubah(asli, asli, '2025')).toBeNull()
  })

  it('nilai berubah, tanggal dari kotak (DD-MM-YYYY) dikirim YYYYMMDD; baris baru ikut', () => {
    const kini: BarisKursWarisan[] = [
      { ...asli[0]!, nilaiKeIDR: '17000', berlakuDariAsli: '02-07-2025' },
      {
        mataUang: 'SGD',
        mataUangID: '10002',
        nilaiKeIDR: '12000',
        berlakuDari: '',
        berlakuSampai: '',
        berlakuDariAsli: '01-07-2025',
        berlakuSampaiAsli: '',
      },
      { mataUang: '', mataUangID: '', nilaiKeIDR: '', berlakuDari: '', berlakuSampai: '', berlakuDariAsli: '', berlakuSampaiAsli: '' },
    ]
    const k = kursBerubah(kini, asli, '2025')
    expect(k?.tahun).toBe('2025')
    expect(k?.baris).toHaveLength(2)
    expect(k?.baris[0]).toMatchObject({ id: '13358', nilaiKeIDR: '17000', berlakuDariAsli: '20250702', berlakuSampaiAsli: '20260630' })
    expect(k?.baris[1]).toMatchObject({ mataUang: 'SGD', berlakuDariAsli: '20250701' })
  })

  it('⭐ baris yang tidak diubah IKUT dikirim bertanda tetap — tidak lepas dari kontrak', () => {
    const kini: BarisKursWarisan[] = [
      asli[0]!,
      { mataUang: 'IDR', mataUangID: '10026', nilaiKeIDR: '1', berlakuDari: '', berlakuSampai: '', berlakuDariAsli: '01-07-2025', berlakuSampaiAsli: '30-06-2026' },
    ]
    const k = kursBerubah(kini, asli, '2025')
    expect(k?.baris).toHaveLength(2)
    expect(k?.baris[0]).toMatchObject({ id: '13358', tetap: true })
    expect(k?.baris[1]?.tetap).toBeUndefined()
  })

  it('⛔ baris yang di-Delete hanya LEPAS: grid tanpa baris itu dikirim, nol hapus di tabel', () => {
    expect(kursBerubah([], asli, '2025')).toEqual({ tahun: '2025', baris: [] })
  })

  // ⛔ Laporan 10 Oktober 2026: kontrak tanpa catatan membaca kurs SELURUH
  // tahun treaty; Save pertama mencatat semuanya sebagai miliknya.
  describe('grid cadangan seluruh-tahun (baris dimuat tanpa tanda tetap)', () => {
    const baris = (id: string, cur: string): BarisKursWarisan => ({
      id,
      mataUang: cur,
      mataUangID: '1',
      nilaiKeIDR: '1',
      berlakuDari: '',
      berlakuSampai: '',
      berlakuDariAsli: '20250701',
      berlakuSampaiAsli: '20260630',
    })
    const cadangan = [baris('1', 'BND'), baris('2', 'CAD'), baris('3', 'CHF')]
    const baru: BarisKursWarisan = { mataUang: 'IDR', mataUangID: '10026', nilaiKeIDR: '1', berlakuDari: '', berlakuSampai: '', berlakuDariAsli: '01-01-2026', berlakuSampaiAsli: '31-12-2026' }

    it('tanpa perubahan → null: kurs kontrak lama tidak disentuh', () => {
      expect(kursBerubah(cadangan, cadangan, '2025')).toBeNull()
    })

    it('⭐ hanya baris BARU yang dikirim — cadangan yang tak disentuh tidak ikut tercatat', () => {
      const k = kursBerubah([...cadangan, baru], cadangan, '2025')
      expect(k?.baris).toHaveLength(1)
      expect(k?.baris[0]).toMatchObject({ mataUang: 'IDR', berlakuDariAsli: '20260101' })
      expect(k?.baris[0]?.tetap).toBeUndefined()
    })

    it('baris cadangan yang DIUBAH ikut dikirim tanpa tanda tetap (supaya benar-benar ditulis)', () => {
      const k = kursBerubah([{ ...cadangan[0]!, nilaiKeIDR: '5' }, cadangan[1]!, cadangan[2]!], cadangan, '2025')
      expect(k?.baris).toHaveLength(1)
      expect(k?.baris[0]).toMatchObject({ id: '1', nilaiKeIDR: '5' })
      expect(k?.baris[0]?.tetap).toBeUndefined()
    })

    it('⭐ pemakai menghapus baris cadangan → sisanya dicatat bertanda tetap', () => {
      const k = kursBerubah([cadangan[0]!], cadangan, '2025')
      expect(k?.baris).toHaveLength(1)
      expect(k?.baris[0]).toMatchObject({ id: '1', tetap: true })
    })
  })

  it('baris milik kontrak yang diubah dikirim TANPA tanda tetap (kalau tidak, ubahannya tidak ditulis)', () => {
    const k = kursBerubah([{ ...asli[0]!, nilaiKeIDR: '17000' }], asli, '2025')
    expect(k?.baris[0]).toMatchObject({ id: '13358', nilaiKeIDR: '17000' })
    expect(k?.baris[0]?.tetap).toBeUndefined()
  })

  it('keYYYYMMDD menerima tiga bentuk', () => {
    expect(keYYYYMMDD('20250701')).toBe('20250701')
    expect(keYYYYMMDD('01-07-2025')).toBe('20250701')
    expect(keYYYYMMDD('2025-07-01')).toBe('20250701')
    expect(keYYYYMMDD('bukan')).toBe('')
  })
})

describe('bolehSave — Save hanya bagi pemegang workbasket posisi berkas', () => {
  it('posisi kosong = Admin; workbasket lain tidak melihat Save', () => {
    expect(bolehSave('', ['ReasTreatyInAdmin'], '', '')).toBe(true)
    expect(bolehSave('ReasTreatyInAdmin', ['ReasTreatyInAdmin'], 'Reject', '')).toBe(true)
    expect(bolehSave('', ['ReasTreatyInSecHead'], '', '')).toBe(false)
    expect(bolehSave('ReasTreatyInSecHead', ['ReasTreatyInAdmin'], 'Accept', '')).toBe(false)
    expect(bolehSave('ReasTreatyInDeptHead', [], 'Accept', '')).toBe(false)
  })
  it('divisi IT (Force Edit) dikecualikan; Resolve Complete tetap tanpa Save', () => {
    expect(bolehSave('ReasTreatyInDeptHead', [], 'Accept', 'IT')).toBe(true)
    expect(bolehSave('', ['ReasTreatyInAdmin'], 'Resolve Complete', 'IT')).toBe(false)
  })
})

describe('bolehActions — syarat tampil tombol Actions', () => {
  it('pemegang workbasket posisi penyetuju', () => {
    expect(bolehActions('ReasTreatyInSecHead', ['ReasTreatyInSecHead'], 'Accept')).toBe(true)
    expect(bolehActions('ReasTreatyInDirector', ['ReasTreatyInSecHead'], 'Accept')).toBe(false)
  })
  it('Admin memakai Submit, bukan Actions; tuntas = tidak tampil', () => {
    expect(bolehActions('ReasTreatyInAdmin', ['ReasTreatyInAdmin'], '')).toBe(false)
    expect(bolehActions('ReasTreatyInSecHead', ['ReasTreatyInSecHead'], 'Resolve Complete')).toBe(false)
  })
})

// ⛔ TOMBOL `Actions` NOL MEMBACA `Position To`.
//
// Permintaan pemilik proses 8 Oktober 2026: satu nama di tampilan, tetapi
// *"bisa diakses semua yg dapat Workbasket itu"*. Sejak `PositionUsername`
// menyusut jadi SATU nama, membandingkannya dengan pemakai yang sedang masuk
// akan mengunci layar bagi pemegang workbasket yang namanya tidak terpilih.
//
// ⚠️ Pasangannya di backend: `TestSemuaPemegangTetapBolehBertindak`.
describe('Actions dinilai dari PERAN, bukan dari nama Position To', () => {
  it('⛔ dua pemegang workbasket yang sama sama-sama boleh, nama siapa pun yang tampil', () => {
    for (const akun of ['SEC1', 'SEC2', 'SIAPA SAJA']) {
      // `akun` sengaja TIDAK ikut dinilai — itulah intinya.
      expect(akun.length).toBeGreaterThan(0)
      expect(bolehActions('ReasTreatyInSecHead', ['ReasTreatyInSecHead'], 'Accept')).toBe(true)
    }
  })

  it('⛔ sumbernya tidak menyebut `pemegangPosisi` sama sekali', () => {
    const src = readFileSync(join(__dirname, 'simpanDokumen.ts'), 'utf8')
    const i = src.indexOf('export function bolehActions')
    expect(i).toBeGreaterThan(0)
    expect(src.slice(i, i + 400)).not.toMatch(/pemegangPosisi|PositionUsername/)
  })
})

// ⛔ Laporan 10 Oktober 2026: Save menampilkan "No column in the table yet, so
// NOT saved: DeductionTotalList, GrossPremiumList, GrossPremiumMinList,
// NetPremiumList". Daftar turunan tingkat spreading tidak dikirim lagi.
describe('susunDokumen — daftar turunan spreading tidak ikut dikirim', () => {
  it('Gross / Gross Min / Deduction / Net / RnmLimit tingkat spreading dibuang; isian tetap', () => {
    const nilai = [{ Currency: 'IDR', CurrencyID: '', Value: '1' }]
    const share = {
      Share: [
        {
          Layer: '1',
          GrossPremiumList: nilai,
          SpreadingListXOL: [
            {
              ReinsTypeID: '10007',
              ReinsTypeName: 'ORS',
              Pct: '22.5',
              GrossPremiumList: nilai,
              GrossPremiumMinList: nilai,
              DeductionTotalList: nilai,
              NetPremiumList: nilai,
              RnmLimitList: nilai,
              // ⭐ 11 Oktober 2026 — rincian anak susunan + OR / R/I tingkat spreading
              Currency: 'IDR',
              Value: '1',
              BreakDownSprdListXOL: [{ ReinsTypeName: 'ORS', ReinsTypeID: '10007', Currency: 'IDR', Amount: '1', SharePct: '100' }],
              RNMSpreadedListXOL: nilai,
              RNMSpreadedListRIXOL: nilai,
              RNMSpreadedListNetXOL: nilai,
            },
          ],
        },
      ],
    }
    const d = susunDokumen(KEPALA, share, null, null)
    const sp = (d.Share as { SpreadingListXOL: Record<string, unknown>[]; GrossPremiumList: unknown }[])[0]!
    expect(sp.SpreadingListXOL[0]).toEqual({ ReinsTypeID: '10007', ReinsTypeName: 'ORS', Pct: '22.5' })
    // daftar tingkat BARIS Share (punya kolom) tetap terkirim
    expect(sp.GrossPremiumList).toEqual(nilai)
  })
})

// ⛔ Laporan 10 Oktober 2026: Copy lalu Save mengubah Rate of Exchange. Kontrak
// baru hasil Copy tidak punya catatan, jadi grid milik sumber harus ikut dikirim.
describe('kursBerubah — Copy (kontrak baru)', () => {
  const milik: BarisKursWarisan = {
    id: '10146',
    mataUang: 'IDR',
    mataUangID: '10026',
    nilaiKeIDR: '1',
    berlakuDari: '',
    berlakuSampai: '',
    berlakuDariAsli: '20260101',
    berlakuSampaiAsli: '20261231',
    tetap: true,
  }
  it('grid milik sumber dikirim bertanda tetap walau tak diubah', () => {
    const k = kursBerubah([milik], [milik], '2026', true)
    expect(k?.baris).toHaveLength(1)
    expect(k?.baris[0]).toMatchObject({ id: '10146', tetap: true })
  })
  it('bukan Copy: tak diubah → null (perilaku lama)', () => {
    expect(kursBerubah([milik], [milik], '2026', false)).toBeNull()
  })
  it('Copy dari sumber tanpa catatan (cadangan seluruh-tahun) → tidak dikirim', () => {
    const { tetap: _t, ...cadangan } = milik
    expect(kursBerubah([cadangan], [cadangan], '2026', true)).toBeNull()
  })
})
