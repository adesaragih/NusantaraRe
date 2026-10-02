// Isian form per activity Pega (`New*` = form baru, `Set*` = form dari baris) dan badan simpannya.

import { describe, expect, it } from 'vitest'

import type { Business, Kontrak, PratinjauSalin, Reinsurer, SecurityReinsurer, TahunTreaty } from './api'
import { opsiBusiness, opsiRate, opsiReinsurer } from './components/cariMaster'
import { penutup, rincianDampak, tanpaAnak } from './components/KonfirmasiHapus'
import { formBusinessBaru, formBusinessDari, keBusinessMasuk, viewRateTampil } from './components/PanelBusiness'
import { formKontrakBaru, formKontrakDari, keKontrakMasuk, opsiJenis } from './components/PanelKontrak'
import { formReinsurerBaru, formReinsurerDari, kelasTotal, keReinsurerMasuk } from './components/PanelReinsurer'
import { formSecurityBaru, formSecurityDari, keSecurityMasuk } from './components/PanelSecurity'
import { sasaranDariPratinjau } from './components/PratinjauSalin'
import { formTahunBaru, formTahunDari, keTahunMasuk } from './pages/MasterContractRetroLife'

const KINI = '2026-09-30 10:00:00'

const tahun: TahunTreaty = {
  id: 'UJI-T1', treatyYear: '2026', underwritingYear: '2025', startDate: '2026-01-01', endDate: '2026-12-31',
  userId: 'UJI-LAMA', tglUpdate: '2026-01-02 08:00:00',
}
const kontrak: Kontrak = {
  id: 'UJI-K1', idTreatyYear: 'UJI-T1', reinsTypeId: '10196', reinsTypeName: 'UJI QS', treatyStartDate: '2026-01-01',
  treatyEndDate: '2026-12-31', userId: 'UJI-LAMA', tglUpdate: '2026-01-02 08:00:00', idr: '1000000000', usd: '75000.5',
  bIdr: '0', bUsd: '0', idrSelisih: '1000000000', usdSelisih: '75000.5',
}
const reinsurer: Reinsurer = {
  id: 'UJI-R1', treatyYearId: 'UJI-T1', treatyContractId: 'UJI-K1', reinsTypeId: '10196', reinsTypeName: 'UJI QS',
  reinsurerId: 'UJI-L01', reinsurerName: 'UJI RE', pctShare: '40', komisi: '2.5', ovrComm: '1', userId: 'UJI-LAMA',
  tglUpdate: '2026-01-02 08:00:00',
}

describe('tahun treaty', () => {
  it('End Period/Add = NewInputTreatyYear_Life_Act: kosong, TGLUPDATE kini, Inputor operator', () => {
    expect(formTahunBaru('UJI-OP', KINI)).toEqual({
      id: '', treatyYear: '', underwritingYear: '', startDate: '', endDate: '', tglUpdate: KINI, userId: 'UJI-OP',
    })
  })

  it('Edit = SetTreatyYearLife_Act: salin lima medan, USERID operator, TGLUPDATE kini', () => {
    expect(formTahunDari(tahun, 'UJI-OP', KINI)).toEqual({
      id: 'UJI-T1', treatyYear: '2026', underwritingYear: '2025', startDate: '2026-01-01', endDate: '2026-12-31',
      tglUpdate: KINI, userId: 'UJI-OP',
    })
  })

  it('badan simpan: tanggal isian (DD-MM-YYYY) jadi YYYY-MM-DD, tanpa medan jejak', () => {
    const f = { ...formTahunBaru('UJI-OP', KINI), treatyYear: ' 2027 ', underwritingYear: '2027', startDate: '01-01-2027', endDate: '31-12-2027' }
    expect(keTahunMasuk(f)).toEqual({ id: '', treatyYear: '2027', underwritingYear: '2027', startDate: '2027-01-01', endDate: '2027-12-31' })
  })
})

describe('kontrak', () => {
  it('Add = NewInputTreatyLimit_Life: batas dan jenis kosong, Inputor operator', () => {
    expect(formKontrakBaru('UJI-OP', KINI)).toEqual({ id: '', reinsTypeId: '', bIdr: '', idr: '', bUsd: '', usd: '', userId: 'UJI-OP', tglUpdate: KINI })
  })

  it('Edit = SetRetroListLife_Act: batas dan jenis disalin, USERID operator', () => {
    const f = formKontrakDari(kontrak, 'UJI-OP', KINI)
    expect(f).toEqual({ id: 'UJI-K1', reinsTypeId: '10196', bIdr: '0', idr: '1000000000', bUsd: '0', usd: '75000.5', userId: 'UJI-OP', tglUpdate: KINI })
  })

  it('badan simpan tanpa tanggal treaty (K4: disalin server dari tahun) dan tanpa selisih (K5)', () => {
    const m = keKontrakMasuk({ ...formKontrakDari(kontrak, 'UJI-OP', KINI), idr: ' 5 ' })
    expect(m).toEqual({ id: 'UJI-K1', reinsTypeId: '10196', idr: '5', usd: '75000.5', bIdr: '0', bUsd: '0' })
    expect(Object.keys(m)).not.toContain('treatyStartDate')
  })

  it('dropdown REINS TYPE: nilai .ID, tampil .Note', () => {
    expect(opsiJenis([{ id: '10196', note: 'UJI QS' }, { id: '10197', note: ' ' }])).toEqual([
      { value: '10196', label: 'UJI QS' },
      { value: '10197', label: '10197' },
    ])
  })
})

describe('reinsurer', () => {
  it('Add = NewInputSecurityLife_Act (Page-New): semua kosong', () => {
    expect(formReinsurerBaru('UJI-OP')).toEqual({ id: '', reinsurerId: '', reinsurerName: '', pctShare: '', komisi: '', ovrComm: '', userId: 'UJI-OP' })
  })

  it('Edit = SetSecurityLife_Act: share, discount (COMMISION), ovr comm disalin', () => {
    expect(formReinsurerDari(reinsurer, 'UJI-OP')).toEqual({
      id: 'UJI-R1', reinsurerId: 'UJI-L01', reinsurerName: 'UJI RE', pctShare: '40', komisi: '2.5', ovrComm: '1', userId: 'UJI-OP',
    })
  })

  it('badan simpan: angka teks di-trim, bukan Number', () => {
    const m = keReinsurerMasuk({ ...formReinsurerDari(reinsurer, 'UJI-OP'), pctShare: ' 33,333333 ' })
    expect(m.pctShare).toBe('33,333333')
    expect(typeof m.komisi).toBe('string')
  })

  it('total ≠ 100 mencolok (tiket 05 AC 19)', () => {
    expect(kelasTotal(true)).toContain('mcrl-total--bukan100')
    expect(kelasTotal(false)).not.toContain('mcrl-total--bukan100')
  })
})

describe('security reinsurer', () => {
  const s: SecurityReinsurer = {
    id: 'UJI-S1', treatyYearId: 'UJI-T1', treatyContractId: 'UJI-K1', treatyReinsurerId: 'UJI-R1', reinsurerId: 'UJI-L02',
    reinsurerName: 'UJI SEC', pctShare: '10', userId: 'UJI-LAMA', tglUpdate: '2026-01-02 08:00:00',
  }

  it('Add: form KOSONG (OQ-MCRL-12)', () => {
    expect(formSecurityBaru('UJI-OP')).toEqual({ id: '', reinsurerId: '', reinsurerName: '', pctShare: '', userId: 'UJI-OP' })
  })

  it('Edit = SetSecurityReinsurerLife_Act: REINSURERNAME, REINSURERID, PCTSHARE, ID', () => {
    const f = formSecurityDari(s)
    expect(keSecurityMasuk(f)).toEqual({ id: 'UJI-S1', reinsurerName: 'UJI SEC', reinsurerId: 'UJI-L02', pctShare: '10' })
  })
})

describe('business', () => {
  const b: Business = {
    id: 'UJI-B1', treatyYearId: 'UJI-T1', treatyYear: '2026', treatyContractId: 'UJI-K1', reinsTypeId: '10196',
    reinsTypeName: 'UJI QS', bizCode: 'UJI-BZ', bizName: 'UJI TERM', riRateId: 'UJI-RT', riRate: 'UJI_RATE_TABLE',
    userId: 'UJI-LAMA', tglUpdate: '2026-01-02 08:00:00',
  }

  it('Add = NewInputBusinessLife_Act langkah 2 (langkah 1 ber-// tidak jalan): lima medan kosong', () => {
    expect(formBusinessBaru('UJI-OP', KINI)).toEqual({ id: '', bizCode: '', bizName: '', riRateId: '', riRate: '', userId: 'UJI-OP', tglUpdate: KINI })
  })

  it('Edit = SetBusinessListLife_Act: ID, BIZCODE, BIZNAME, RIRATEID, RIRATE; USERID operator', () => {
    const f = formBusinessDari(b, 'UJI-OP')
    expect(f.userId).toBe('UJI-OP')
    expect(keBusinessMasuk(f)).toEqual({ id: 'UJI-B1', bizCode: 'UJI-BZ', bizName: 'UJI TERM', riRateId: 'UJI-RT', riRate: 'UJI_RATE_TABLE' })
  })

  it("View Rate form tampil bila RIRATEID != ''", () => {
    expect(viewRateTampil(formBusinessDari(b, 'UJI-OP'))).toBe(true)
    expect(viewRateTampil(formBusinessBaru('UJI-OP', KINI))).toBe(false)
  })

  it('Copy to all Reinstype: Yes mengirim PERSIS sasaran pratinjau', () => {
    const p: PratinjauSalin = { business: b, reinsTypeId: '10196', sasaran: [{ ...kontrak, id: 'UJI-K2' }, { ...kontrak, id: 'UJI-K3' }] }
    expect(sasaranDariPratinjau(p)).toEqual(['UJI-K2', 'UJI-K3'])
  })
})

describe('popup hapus', () => {
  it('kontrak menyebut ketiga anak, reinsurer menyebut security, daun nol baris', () => {
    const d = { reinsurer: 2, security: 3, business: 4 }
    expect(rincianDampak('kontrak', d)).toEqual(['2 reinsurer row(s)', '3 security reinsurer row(s)', '4 business row(s)'])
    expect(rincianDampak('reinsurer', d)).toEqual(['3 security reinsurer row(s)'])
    expect(rincianDampak('security', d)).toEqual([])
    expect(rincianDampak('business', d)).toEqual([])
  })

  it('Cancel tidak berlaku selama Yes berjalan (AC 30/35)', () => {
    let ditutup = 0
    const tutup = () => {
      ditutup++
    }
    penutup(true, tutup)()
    expect(ditutup).toBe(0)
    penutup(false, tutup)()
    expect(ditutup).toBe(1)
  })

  it('tanpa anak dinyatakan (AC 36)', () => {
    expect(tanpaAnak({ reinsurer: 0, security: 0, business: 0 })).toBe(true)
    expect(tanpaAnak({ reinsurer: 0, security: 1, business: 0 })).toBe(false)
  })
})

describe('butir autocomplete', () => {
  it('REINSURER NAME: tampil ClientName, isi ID', () => {
    expect(opsiReinsurer({ id: 'UJI-L01', clientName: 'UJI RE' })).toEqual({ value: 'UJI-L01', label: 'UJI RE', keterangan: 'UJI-L01' })
  })

  it('BUSINESS NAME: tampil Note, isi ID, keterangan OLDID', () => {
    expect(opsiBusiness({ id: 'UJI-BZ', note: 'UJI TERM', oldId: 'L01' })).toEqual({ value: 'UJI-BZ', label: 'UJI TERM', keterangan: 'L01' })
  })

  it('R/I RATE: tampil USEDBY, isi ID', () => {
    expect(opsiRate({ id: 'UJI-RT', usedBy: 'UJI_RATE_TABLE' })).toEqual({ value: 'UJI-RT', label: 'UJI_RATE_TABLE', keterangan: 'UJI-RT' })
  })
})
