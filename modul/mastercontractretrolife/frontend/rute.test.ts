import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

// Menu Master Contract Retro Life ditekan → layar kembali ke awal (03-10-2026).

describe('rute Master Contract Retro Life', () => {
  it('halaman dipasang ulang setiap kali menunya dipilih (key = ketukMenu)', () => {
    const rute = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')
    expect(rute).toContain('<MasterContractRetroLife key={ketukMenu ?? 0} />')
    const app = readFileSync(join(__dirname, '..', '..', '..', 'frontend', 'App.tsx'), 'utf8')
    expect(app).toContain('ketukMenu={ketukMenu}')
  })
})

describe('Inputor = akun login (04-10-2026)', () => {
  it('rute mencatat akun login; operatorKini membacanya lebih dahulu dari stub', async () => {
    const rute = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')
    expect(rute).toContain('catatOperator(masuk.akunID)')
    const { catatOperator, operatorKini } = await import('./tampilan')
    catatOperator('UJI-OPERATOR')
    expect(operatorKini()).toBe('UJI-OPERATOR')
    catatOperator('')
  })
})

describe('Inputor dan Modified Date di paling akhir setiap form (04-10-2026)', () => {
  it('di kelima form, Inputor adalah isian TERAKHIR form-grid', () => {
    const berkas = [
      'pages/MasterContractRetroLife.tsx',
      'components/PanelKontrak.tsx',
      'components/PanelBusiness.tsx',
      'components/PanelReinsurer.tsx',
      'components/PanelSecurity.tsx',
    ]
    for (const f of berkas) {
      const kode = readFileSync(join(__dirname, f), 'utf8')
      const awal = kode.indexOf('<div className="form-grid">')
      // Sampai baris tombol: grid dapat memuat div bersarang (mis. sel R/I RATE).
      const akhir = kode.indexOf('<div className="aksi-baris">', awal)
      const grid = kode.slice(awal, akhir)
      const iInputor = grid.search(/label=\{\w+_MCRL\.formInputor\}/)
      const sesudah = grid.slice(iInputor).replace(/^[^\n]*\n/, '')
      expect(iInputor, f).toBeGreaterThan(-1)
      expect(/<(Field|Pilih|PilihSaring|FieldTanggal)\b/.test(sesudah), f).toBe(false)
    }
  })
})

describe('satu REINS TYPE satu kontrak per tahun (04-10-2026)', () => {
  it('jenis milik kontrak lain disembunyikan; milik kontrak yang diubah tetap ada', async () => {
    const { jenisTersedia } = await import('./components/PanelKontrak')
    const jenis = [
      { id: '10196', note: 'QS' },
      { id: '10197', note: '2ND QS' },
      { id: '10200', note: 'OR' },
    ]
    const kontrak = [{ id: 'K1', reinsTypeId: '10196' }, { id: 'K2', reinsTypeId: '10200' }] as never[]
    expect(jenisTersedia(jenis, kontrak, '').map((j) => j.id)).toEqual(['10197'])
    expect(jenisTersedia(jenis, kontrak, 'K1').map((j) => j.id)).toEqual(['10196', '10197'])
  })
})

describe('kepala panel Reins Type (04-10-2026)', () => {
  it('ID Treaty Year paling depan, lalu Underwriting Year dan Transaction Year', () => {
    const kode = readFileSync(join(__dirname, 'components', 'PanelKontrak.tsx'), 'utf8')
    const uw = kode.indexOf('[TAHUN_MCRL.formUnderwritingYear, induk.underwritingYear]')
    const tr = kode.indexOf('[TAHUN_MCRL.formTransactionYear, induk.treatyYear]')
    const id = kode.indexOf('[KONTRAK_MCRL.idTreatyYear, induk.id]')
    expect(id).toBeGreaterThan(-1)
    expect(uw).toBeGreaterThan(id)
    expect(tr).toBeGreaterThan(uw)
  })
})

describe('tanpa paginasi (04-10-2026)', () => {
  it('kelima daftar menampilkan semua baris - nol Penomoran / potongHalaman', () => {
    for (const f of [
      'pages/MasterContractRetroLife.tsx',
      'components/PanelKontrak.tsx',
      'components/PanelBusiness.tsx',
      'components/PanelReinsurer.tsx',
      'components/PanelSecurity.tsx',
    ]) {
      const kode = readFileSync(join(__dirname, f), 'utf8')
      expect(kode, f).not.toContain('<Penomoran')
      expect(kode, f).not.toContain('potongHalaman(')
    }
  })
})

describe('form Contract tanpa TREATY START/END (04-10-2026)', () => {
  it('kedua isian tidak dirender di form', () => {
    const kode = readFileSync(join(__dirname, 'components', 'PanelKontrak.tsx'), 'utf8')
    expect(kode).not.toContain('label={KONTRAK_MCRL.formTreatyStart}')
    expect(kode).not.toContain('label={KONTRAK_MCRL.formTreatyEnd}')
  })
})

describe('Business List tanpa ID Reins Type (04-10-2026)', () => {
  it('kepala panel Business DAN Reinsurer tidak menampilkan ID Reins Type', () => {
    const reas = readFileSync(join(__dirname, 'components', 'PanelReinsurer.tsx'), 'utf8')
    expect(reas).not.toContain('[REINSURER_MCRL.idReinsType, induk.id]')
    const kode = readFileSync(join(__dirname, 'components', 'PanelBusiness.tsx'), 'utf8')
    expect(kode).not.toContain('[BUSINESS_MCRL.idReinsType, induk.id]')
    expect(kode).toContain('[BUSINESS_MCRL.reinsType, induk.reinsTypeName]')
  })
})


describe('View Rate di pojok kanan baris tombol form Business (04-10-2026)', () => {
  it('tombol berada di baris aksi form, sesudah Cancel, ber-kelas mcrl-aksi-kanan', () => {
    const kode = readFileSync(join(__dirname, 'components', 'PanelBusiness.tsx'), 'utf8')
    const aksi = kode.indexOf('<div className="aksi-baris">')
    const batal = kode.indexOf('{BUSINESS_MCRL.cancel}', aksi)
    const tombol = kode.indexOf('{BUSINESS_MCRL.viewRateForm}')
    expect(aksi).toBeGreaterThan(-1)
    expect(tombol).toBeGreaterThan(batal)
    expect(kode).toContain('className="btn btn--ghost mcrl-aksi-kanan"')
    expect(kode).not.toContain('mcrl-sel-rate')
  })
})

describe('satu Business Name sekali per kontrak (04-10-2026)', () => {
  it('nama yang sudah dipakai business lain disembunyikan; milik yang diubah tetap ada', async () => {
    const { bisnisTersedia } = await import('./components/PanelBusiness')
    const pilihan = [
      { value: '10172', label: 'GROUP TERM LIFE' },
      { value: '10180', label: 'CREDIT LIFE' },
      { value: '10190', label: 'HEALTH' },
    ]
    const daftar = [{ id: 'B1', bizCode: '10172' }, { id: 'B2', bizCode: '10180' }] as never[]
    expect(bisnisTersedia(pilihan, daftar, '').map((o) => o.value)).toEqual(['10190'])
    expect(bisnisTersedia(pilihan, daftar, 'B1').map((o) => o.value)).toEqual(['10172', '10190'])
  })
})

describe('satu Reinsurer Name sekali per kontrak (04-10-2026)', () => {
  it('reinsurer yang sudah dipakai disembunyikan; milik baris yang diubah tetap ada', async () => {
    const { reinsurerTersedia } = await import('./components/PanelReinsurer')
    const pilihan = [
      { value: 'L01', label: 'AON' },
      { value: 'L02', label: 'GALLAGHER' },
      { value: 'L03', label: 'MARSH' },
    ]
    const daftar = [{ id: 'R1', reinsurerId: 'L01' }, { id: 'R2', reinsurerId: 'L02' }] as never[]
    expect(reinsurerTersedia(pilihan, daftar, '').map((o) => o.value)).toEqual(['L03'])
    expect(reinsurerTersedia(pilihan, daftar, 'R1').map((o) => o.value)).toEqual(['L01', 'L03'])
  })
})

describe('Security Reinsurer memakai aturan Reinsurer List (04-10-2026)', () => {
  it('dropdown tanpa nama yang sudah dipakai, dan baris Total Share tampil', () => {
    const kode = readFileSync(join(__dirname, 'components', 'PanelSecurity.tsx'), 'utf8')
    expect(kode).toContain('opsi={reinsurerTersedia(master.pilihan, daftar, form.id)}')
    expect(kode).toContain('kelasTotal(jawab.totalBukan100 === true)')
    expect(kode).toContain("from './aturanDaftar'")
    expect(kode).not.toContain("from './PanelReinsurer'")
  })
})
