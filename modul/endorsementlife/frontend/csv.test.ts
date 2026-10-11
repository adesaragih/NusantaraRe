import { describe, expect, it } from 'vitest'

import { KEPALA_UNDUH_CSV, barisUnggahan, namaBerkasUnduh, pemisahJudul, rakitCSV, uraiCSV } from './csv'

describe('CSV unggahan Endorsement Life (tampilan View Upload dan Generate Data Detail)', () => {
  it('mengurai kutip, koma di dalam kutip, CRLF, dan BOM', () => {
    const bom = String.fromCharCode(0xfeff)
    expect(uraiCSV(`${bom}PLAN,GROSS_PREMIUM\r\n"UJI, A","10,5"\r\n"UJI ""B""",7\n`)).toEqual([
      ['PLAN', 'GROSS_PREMIUM'],
      ['UJI, A', '10,5'],
      ['UJI "B"', '7'],
    ])
    expect(uraiCSV('')).toEqual([])
  })

  it('baris unggahan berkunci judul huruf besar; baris kosong dilewati', () => {
    const b = barisUnggahan('plan,Policy_Holder\nUJI-PLAN,UJI-PH\n,\nUJI-PLAN-2,\n')
    expect(b).toEqual([
      { PLAN: 'UJI-PLAN', POLICY_HOLDER: 'UJI-PH' },
      { PLAN: 'UJI-PLAN-2', POLICY_HOLDER: '' },
    ])
  })

  it('Generate Data Detail: kepala b276 berurutan, nilai berkutip bila perlu', () => {
    expect(KEPALA_UNDUH_CSV).toHaveLength(23)
    expect(KEPALA_UNDUH_CSV.slice(0, 3)).toEqual(['POLICY_NO', 'POLICY_HOLDER', 'CERTIFICATE_NO'])
    const teks = rakitCSV(['PLAN', 'GROSS_PREMIUM'], [{ PLAN: 'UJI, A', GROSS_PREMIUM: '10,5' }, { PLAN: 'UJI "B"' }])
    expect(teks).toBe('PLAN,GROSS_PREMIUM\r\n"UJI, A","10,5"\r\n"UJI ""B""",\r\n')
    expect(uraiCSV(teks)).toEqual([['PLAN', 'GROSS_PREMIUM'], ['UJI, A', '10,5'], ['UJI "B"', '']])
  })

  it('nama berkas DetailUpload + stempel waktu (b270, b272)', () => {
    expect(namaBerkasUnduh(new Date(2026, 9, 1, 9, 5, 7))).toBe('DetailUpload_20261001_090507.csv')
  })
})

describe('pemisah unggahan sama dengan server (keputusan work owner 10-10-2026, opsi a)', () => {
  const mau = [{ PLAN: 'UJI-PLAN', POLICY_HOLDER: 'UJI-PH', GROSS_PREMIUM: '14216802,6', NET_PREMIUM: '509,9775' }]

  it('titik koma (BOM diabaikan), desimal koma berkutip dan tanpa kutip', () => {
    expect(barisUnggahan('\uFEFFPLAN;POLICY_HOLDER;GROSS_PREMIUM;NET_PREMIUM\r\nUJI-PLAN;UJI-PH;"14216802,6";509,9775\r\n')).toEqual(mau)
    expect(pemisahJudul('PLAN;POLICY_HOLDER\n1;2\n')).toBe(';')
  })

  it('koma tetap koma; judul yang memuat keduanya = koma', () => {
    expect(barisUnggahan('PLAN,POLICY_HOLDER,GROSS_PREMIUM,NET_PREMIUM\nUJI-PLAN,UJI-PH,"14216802,6","509,9775"\n')).toEqual(mau)
    expect(pemisahJudul('A;B,C\n')).toBe(',')
  })
})
