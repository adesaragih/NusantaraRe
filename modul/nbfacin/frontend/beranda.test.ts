// Kotak masuk Beranda NB FacIn (work owner 07-10-2026) - data uji sintetis.

import { describe, expect, it } from 'vitest'

import type { BarisKotakMasukNbFacIn } from './api'
import { keDaftarBerandaNbFacIn, lamaDiUw, persenInbox, sejak, tampilKolomUw, tanggalInbox } from './beranda'
import { KOLOM_INBOX_NBFACIN, WORKBASKET_UW_NBFACIN } from './labels'
import { PENDAFTARAN_MENU } from './menu'

describe('kotak masuk Beranda NB FacIn', () => {
  it('modul mendaftarkan penghitung dan daftar kotak masuk', () => {
    expect(PENDAFTARAN_MENU.antreanBeranda).toBeTypeOf('function')
    expect(PENDAFTARAN_MENU.daftarBeranda).toBeTypeOf('function')
  })
  it('kolom = gambar layar inbox "New Business" Pega, berurutan', () => {
    expect(KOLOM_INBOX_NBFACIN.map((k) => k.label)).toEqual([
      'ID', 'Insured Name', 'Business name', 'Source Of Business', 'Ceding Company', 'Inception Date', '%Share RNM', 'Marketing Name',
      'Status', 'Created date', 'Offering Date', 'Validity Date', 'Status Inbox', '', 'Time Since Last Update',
    ])
  })
  it('format: Inception DD/MM/YY, tanggal lain DD/MM/YYYY, %Share 2 desimal, waktu sejak', () => {
    expect(tanggalInbox('31-05-2026', true)).toBe('31/05/26')
    expect(tanggalInbox('06-10-2026')).toBe('06/10/2026')
    expect(persenInbox('10')).toBe('10,00')
    expect(persenInbox('12.5')).toBe('12,50')
    expect(persenInbox('')).toBe('')
    const kini = new Date(2026, 9, 7, 10, 0, 0)
    expect(sejak('2026-10-07 09:59:00', kini)).toBe('1m ago')
    expect(sejak('2026-10-07 10:00:00', kini)).toBe('just now')
    expect(sejak('', kini)).toBe('')
  })
  it('baris -> daftar Beranda: id = nomor case, sel sesuai kolom', () => {
    const b: BarisKotakMasukNbFacIn = {
      caseId: 'NB-UJI-1', insuredName: 'UJI-TERTANGGUNG', businessName: 'UJI-COB', sourceOfBusiness: 'UJI-SOB', cedingCompany: 'UJI-CEDING',
      inceptionDate: '31-05-2026', shareRnm: '10', marketing: 'UJI-MKT', status: 'UJI-STATUS', createdDate: '06-10-2026',
      offeringDate: '05-10-2026', validityDate: '06-10-2026', statusInbox: 'NB IS IN UJI\'S INBOX', aksiTerakhir: 'Accept by UJI',
      tglUpdate: '2026-10-07 09:59:00', submitToUw: '',
    }
    const d = keDaftarBerandaNbFacIn([b], new Date(2026, 9, 7, 10, 0, 0))
    expect(d.kolom.map((k) => k.kunci)).toEqual(KOLOM_INBOX_NBFACIN.map((k) => k.kunci))
    expect(d.baris[0]?.id).toBe('NB-UJI-1')
    expect(d.baris[0]?.sel).toMatchObject({ mulai: '31/05/26', share: '10,00', statusInbox: "NB IS IN UJI'S INBOX", aksiTerakhir: 'Accept by UJI', sejak: '1m ago' })
    for (const k of KOLOM_INBOX_NBFACIN) expect(d.baris[0]?.sel).toHaveProperty(k.kunci)
  })
})

describe('kolom "Time in Underwriting" (inbox workbasket UW, work owner 07-10-2026)', () => {
  const b: BarisKotakMasukNbFacIn = {
    caseId: 'NB-UJI-2', insuredName: 'UJI', businessName: 'UJI', sourceOfBusiness: 'UJI', cedingCompany: 'UJI', inceptionDate: '31-05-2026',
    shareRnm: '10', marketing: 'UJI', status: 'UJI', createdDate: '06-10-2026', offeringDate: '05-10-2026', validityDate: '06-10-2026',
    statusInbox: "NB IS IN UNDERWRITING'S INBOX", aksiTerakhir: 'Accept by UJI', tglUpdate: '2026-10-07 09:59:00', submitToUw: '2026-10-06 09:00:00',
  }
  const kini = new Date(2026, 9, 7, 10, 0, 0)

  it('lima workbasket UW sesuai daftar work owner', () => {
    expect([...WORKBASKET_UW_NBFACIN]).toEqual([
      'ReasFacInJuniorUnderwritingB', 'ReasFacInJuniorUnderwritingA', 'ReasFacInUnderwriting', 'ReasFacInSeniorUnderwriting', 'ReasFacInDepHeadUnderwriting',
    ])
  })
  it('format gambar layar "1d 1h ago"; di bawah sehari gaya Time Since Last Update; kosong = belum ke UW', () => {
    expect(lamaDiUw('2026-10-06 09:00:00', kini)).toBe('1d 1h ago')
    expect(lamaDiUw('2026-10-05 10:00:00', kini)).toBe('2d ago')
    expect(lamaDiUw('2026-10-07 07:30:00', kini)).toBe('2h ago')
    expect(lamaDiUw('', kini)).toBe('')
  })
  it('workbasket UW: kolom disisip sesudah Validity Date, sebelum Status Inbox', () => {
    for (const wb of WORKBASKET_UW_NBFACIN) {
      const d = keDaftarBerandaNbFacIn([b], kini, wb)
      const label = d.kolom.map((k) => k.label)
      expect(label.indexOf('Time in Underwriting')).toBe(label.indexOf('Validity Date') + 1)
      expect(label.indexOf('Status Inbox')).toBe(label.indexOf('Time in Underwriting') + 1)
      expect(d.baris[0]?.sel).toMatchObject({ diUw: '1d 1h ago' })
    }
  })
  it('workbasket non-UW (Team Leader, marketing): tanpa kolom; daftar semua workbasket: tampil bila ada baris sudah ke UW', () => {
    expect(tampilKolomUw('ReasFacInTeamLeader', [b])).toBe(false)
    expect(keDaftarBerandaNbFacIn([b], kini, 'ReasFacInTeamLeader').kolom.map((k) => k.label)).not.toContain('Time in Underwriting')
    expect(tampilKolomUw(null, [b])).toBe(true)
    expect(tampilKolomUw(null, [{ ...b, submitToUw: '' }])).toBe(false)
  })
})
