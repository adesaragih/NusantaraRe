import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { BUAT_EDM, KASUS_EDM, POLIS_LAMA_EDM, SIMPAN_EDM, UMUM_EDM, UNGGAH_EDM } from '../labels'
import { bagianKepalaKasus, labelMetodePremi } from '../tampilan'

// Layar "Endorsement Life Detail" dirapikan - keputusan work owner 07-10-2026 (tampilan saja).

const SUMBER = readFileSync(join(__dirname, 'InputEDMLife.tsx'), 'utf8')
const CSS = readFileSync(join(__dirname, '..', 'endorsementlife.css'), 'utf8')

const KEPALA = {
  PRODUCT_NAME: 'UJI-PRODUK DENGAN NAMA YANG SANGAT PANJANG',
  PRODUCT_NAME_ID: 'UJI-P1',
  TYPE: 'QR',
  TYPE_CEDING: '1',
  BUSINESS_NAME: 'UJI-COB',
  PRO_RATE_TYPE: '2',
  CEDING_CO_NAME: 'UJI-CEDING',
  SOB_NAME: 'UJI-SOB',
  POLICY_HOLDER_NAME: 'UJI-PH',
  MARKETING_NAME: 'UJI-MO',
}

function kasus(edmType: string) {
  return {
    policyNo: 'UJI-PL-1', edmType, edmTypePerubahanData: 'Perubahan data (DOB)', edmTypeBatal: 'Batal Eksternal',
    description: 'UJI-CATATAN', edmDate: '2026-10-07', kepala: KEPALA,
  }
}

describe('kepala kasus tiga bagian berjudul', () => {
  it('Policy Data, Parties, Endorsement dengan medan sesuai daftar', () => {
    const b = bagianKepalaKasus(kasus('1'))
    expect(b.map((x) => x.judul)).toEqual([KASUS_EDM.bagianPolis, KASUS_EDM.bagianPihak, KASUS_EDM.bagianEndorsement])
    expect(b[0]?.medan.map(([l]) => l)).toEqual([
      BUAT_EDM.policyNo, KASUS_EDM.productName, KASUS_EDM.productNameId, KASUS_EDM.type,
      KASUS_EDM.reinsuranceSystem, KASUS_EDM.classOfBusiness, KASUS_EDM.premiumMethod,
    ])
    expect(b[1]?.medan.map(([l]) => l)).toEqual([KASUS_EDM.ceding, KASUS_EDM.sob, KASUS_EDM.policyHolder, KASUS_EDM.marketingOfficer])
    expect(b[2]?.medan.map(([l]) => l)).toEqual([
      KASUS_EDM.edmType, BUAT_EDM.edmTypePerubahanData, BUAT_EDM.edmDate, KASUS_EDM.description,
    ])
    // Nilai tampil: Premium Method bernama, tanggal dd/mm/yyyy, Description & Product Name selebar bagian.
    expect(b[0]?.medan.find(([l]) => l === KASUS_EDM.premiumMethod)?.[1]).toBe('Annually')
    expect(b[2]?.medan.find(([l]) => l === BUAT_EDM.edmDate)?.[1]).toBe('07/10/2026')
    expect(b[2]?.medan.find(([l]) => l === KASUS_EDM.description)?.[2]).toBe(true)
    expect(b[0]?.medan.find(([l]) => l === KASUS_EDM.productName)?.[2]).toBe(true)
  })

  it('sub-jenis hanya yang sesuai EDM Type', () => {
    const batal = bagianKepalaKasus(kasus('3'))[2]?.medan.map(([l, n]) => [l, n])
    expect(batal).toContainEqual([BUAT_EDM.edmTypeBatal, 'Batal Eksternal'])
    expect(batal?.map(([l]) => l)).not.toContain(BUAT_EDM.edmTypePerubahanData)
    const ubah = bagianKepalaKasus(kasus('1'))[2]?.medan.map(([l]) => l)
    expect(ubah).not.toContain(BUAT_EDM.edmTypeBatal)
  })

  it('struktur layar: alert Locked di bawah judul, bagian berjudul, kolom tiga bertumpuk di layar sempit', () => {
    const iJudul = SUMBER.indexOf('{KASUS_EDM.judul}')
    const iAlert = SUMBER.indexOf('<p className="alert alert--info edm-kasus__info"')
    const iBagian = SUMBER.indexOf('<div className="edm-kasus__tigabagian">')
    expect(iJudul).toBeGreaterThan(0)
    expect(iAlert).toBeGreaterThan(iJudul)
    expect(iBagian).toBeGreaterThan(iAlert)
    expect(SUMBER).toContain('{UMUM_EDM.terkunci}')
    expect(SUMBER).toContain('<h3 className="edm-subjudul">{b.judul}</h3>')
    expect(CSS).toMatch(/\.edm-kasus__tigabagian \{[^}]*grid-template-columns: repeat\(3, minmax\(0, 1fr\)\)/)
    expect(CSS).toMatch(/@media \(max-width: 1024px\) \{\s*\.edm-kasus__tigabagian \{\s*grid-template-columns: minmax\(0, 1fr\)/)
    expect(CSS).toMatch(/\.edm-kasus__grid dd \{[^}]*overflow-wrap: anywhere/)
    expect(UMUM_EDM.terkunci.startsWith('Locked')).toBe(true)
  })

  it('Participants berjudul dengan toolbar tiga tombol CSV; kondisi tombol tidak berubah', () => {
    const kepala = SUMBER.slice(SUMBER.indexOf('<div className="edm-kasus__kepalapeserta">'), SUMBER.indexOf('{periksa !== null'))
    expect(kepala).toContain('{KASUS_EDM.peserta}')
    for (const t of ['UNGGAH_EDM.uploadCsv', 'UNGGAH_EDM.viewUpload', 'UNGGAH_EDM.addCsvData']) {
      expect(kepala).toContain(t)
    }
    expect(kepala).toContain('{denganCentang && terbukaKasus && (')
    expect(kepala).toContain('disabled={kasus.csvTerkunci} onClick={() => setJendelaCSV(\'unggah\')}')
    expect(kepala).toContain('disabled={kasus.csvTerkunci || berkasCSV === null || menambah}')
    expect(KASUS_EDM.peserta).toBe('Participants')
    expect(UNGGAH_EDM.uploadCsv).toBe('Upload CSV')
  })

  it('tombol kanan atas tetap: View Old Policy, Save (utama), Back', () => {
    const kepala = SUMBER.slice(SUMBER.indexOf('<header className="inbox__kepala edm-kasus__kepala">'), SUMBER.indexOf('</header>'))
    expect(kepala.indexOf('POLIS_LAMA_EDM.viewOldPolicy')).toBeLessThan(kepala.indexOf('SIMPAN_EDM.save'))
    expect(kepala.indexOf('SIMPAN_EDM.save')).toBeLessThan(kepala.indexOf('UMUM_EDM.kembali'))
    expect(kepala).toContain('className="btn btn--primary btn--sm" disabled={kasus.sudahSimpan || menyimpan}')
    expect(POLIS_LAMA_EDM.viewOldPolicy).toBe('View Old Policy')
    expect(SIMPAN_EDM.save).toBeTruthy()
  })
})

describe('Premium Method bernama', () => {
  it('1/2/3 → Single/Annually/Others; nilai lain apa adanya; kosong → tanda kosong', () => {
    expect(labelMetodePremi('1')).toBe('Single')
    expect(labelMetodePremi('2')).toBe('Annually')
    expect(labelMetodePremi('3')).toBe('Others')
    expect(labelMetodePremi(' 2 ')).toBe('Annually')
    expect(labelMetodePremi('9')).toBe('9')
    expect(labelMetodePremi('')).toBe('—')
  })
})
