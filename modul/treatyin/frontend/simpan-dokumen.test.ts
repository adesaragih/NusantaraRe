// Penyusun isi tombol tulis — Save/Submit/Actions/Decline offer.

import { describe, expect, it } from 'vitest'

import type { BarisKursWarisan, LimitsAkar, ShareNP } from './api'
import { bolehActions, keYYYYMMDD, kursBerubah, susunDokumen, type KepalaForm } from './simpanDokumen'

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

describe('kursBerubah — hanya baris baru dan yang berubah', () => {
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

  it('⛔ baris yang dihapus dari grid TIDAK dikirim untuk dihapus', () => {
    expect(kursBerubah([], asli, '2025')).toBeNull()
  })

  it('keYYYYMMDD menerima tiga bentuk', () => {
    expect(keYYYYMMDD('20250701')).toBe('20250701')
    expect(keYYYYMMDD('01-07-2025')).toBe('20250701')
    expect(keYYYYMMDD('2025-07-01')).toBe('20250701')
    expect(keYYYYMMDD('bukan')).toBe('')
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
