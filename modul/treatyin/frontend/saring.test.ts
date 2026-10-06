import { describe, expect, it } from 'vitest'
import { saringTerdekat, peringkatKedekatan } from './saring'

const D = [
  { id: 'G0000006', nama: 'ASURANSI ADIRA DINAMIKA', kembar: true },
  { id: 'G0000010', nama: 'ASURANSI DIGITAL BERSAMA (D/H SARANA LINDUNG UPAYA)', kembar: false },
  { id: 'G0000038', nama: 'ASURANSI KRESNA MITRA TBK', kembar: true },
  { id: 'G0000041', nama: 'ASURANSI MEGA PRATAMA', kembar: false },
  { id: 'G0000050', nama: 'ASURANSI SIMAS INSURTECH', kembar: false },
]

describe('saringan terdekat', () => {
  it('"asuransi adira" hanya memunculkan yang itu', () => {
    const h = saringTerdekat(D, 'asuransi adira')
    expect(h.map((x) => x.nama)).toEqual(['ASURANSI ADIRA DINAMIKA'])
  })
  it('"adira" menemukan lewat awal KATA', () => {
    expect(saringTerdekat(D, 'adira').map((x) => x.nama)).toEqual(['ASURANSI ADIRA DINAMIKA'])
  })
  it('"asuransi" memang memunculkan semuanya — ketikannya memang seluas itu', () => {
    expect(saringTerdekat(D, 'asuransi')).toHaveLength(5)
  })
  it('kode pengenal tetap bisa dicari', () => {
    expect(saringTerdekat(D, 'G0000038').map((x) => x.id)).toEqual(['G0000038'])
  })
  it('nol cocok mengembalikan kosong', () => {
    expect(saringTerdekat(D, 'zzz')).toHaveLength(0)
  })
  it('peringkat: sama persis < awalan < awal kata < mengandung', () => {
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'asuransi adira dinamika')).toBe(0)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'asuransi')).toBe(1)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'adira')).toBe(2)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'dir')).toBe(3)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'zzz')).toBeNull()
  })
})
