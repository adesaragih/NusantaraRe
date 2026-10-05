import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, ambilHak, ambilLama, cariInduk, periksaNama, salinLama, tambah, ubah, type Detail, type OrgLama, type Pilihan } from './api'
import {
  buangDi,
  catatanLama,
  gantiDi,
  isianDari,
  isianKosong,
  jabatanAkun,
  jumlahHalaman,
  KODE_LAIN,
  kodeAreaLain,
  kataNama,
  kodeNegara,
  labelKode,
  opsiAkun,
  opsiDari,
  opsiKodeArea,
  opsiNegara,
  opsiTitle,
  periksa,
  picKosong,
  potong,
  ringkasSalin,
  saringLama,
  titleDiNama,
} from './aturan'
import { CD } from './labels'

const detail: Detail = {
  id: 'ASM-SFAGIS-WORK-ORG ORG-115', idView: 'ORG-115', nama: 'UJI Anak Usaha', title: 'PT.', npwp: '', country: '001',
  countryName: 'UJI NEGARA SATU', businessField: '23', parentId: '', parentName: 'UJI Induk Grup', note: '',
  createdBy: '', createdAt: '', updatedBy: '', updatedAt: '',
  pic: [{ userIdentifier: 'PIC-2', nama: 'UJI Kontak', position: 'Staff', gender: '2', email: '', dateOfBirth: '15-01-1990', phone: '' }],
  alamat: [{ asal: 'UJI Jalan', type: '2', address: 'UJI Jalan', telfax: [{ type: '2', code: '021', no: 'UJI-FAX' }] }],
}

const telfax: Pilihan[] = [
  { kode: '2', label: 'FAX', aktif: false },
  { kode: '3', label: 'MOBILE PHONE', aktif: true },
  { kode: '5', label: 'OFFICE PHONE', aktif: true },
]

describe('aturan Company Detail', () => {
  it('isian ubah menyalin organisasi, PIC, dan alamat tanpa berbagi objek', () => {
    const isi = isianDari(detail)
    expect(isi).toMatchObject({ nama: 'UJI Anak Usaha', country: '001', businessField: '23', parentId: '' })
    isi.alamat[0]!.telfax[0]!.no = 'UJI-LAIN'
    expect(detail.alamat[0]!.telfax[0]!.no).toBe('UJI-FAX')
  })

  it('pemeriksaan awal: wajib isi utama, PIC, alamat ganda, dan nomor', () => {
    expect(periksa(isianKosong())).toBe(CD.galatNama)
    const isi = isianDari(detail)
    expect(periksa(isi)).toBeNull()
    expect(periksa({ ...isi, country: ' ' })).toBe(CD.galatCountry)
    expect(periksa({ ...isi, businessField: '' })).toBe(CD.galatBusinessField)
    expect(periksa({ ...isi, pic: [...isi.pic, picKosong()] })).toBe(CD.galatPICNama(2))
    expect(periksa({ ...isi, pic: [{ ...picKosong(), nama: 'UJI' }] })).toBe(CD.galatPICPosisi(1))
    expect(periksa({ ...isi, alamat: [...isi.alamat, { asal: '', type: '1', address: ' UJI Jalan ', telfax: [] }] })).toBe(
      CD.galatAlamatGanda(2, 1),
    )
    expect(periksa({ ...isi, alamat: [{ ...isi.alamat[0]!, telfax: [{ type: '3', code: '', no: '' }] }] })).toBe(CD.galatNomor(1, 1))
    expect(periksa({ ...isi, alamat: [{ ...isi.alamat[0]!, type: '' }] })).toBe(CD.galatAlamatType(1))
  })

  it('dropdown hanya menawarkan pilihan aktif; nilai lama yang terpasang tetap terbaca', () => {
    expect(opsiDari(telfax, '').map((o) => o.value)).toEqual(['3', '5'])
    expect(opsiDari(telfax, '2')).toContainEqual({ value: '2', label: CD.nilaiLama('FAX') })
    expect(opsiDari([{ kode: '0542', label: '', aktif: true }], '', true)).toEqual([{ value: '0542', label: '0542' }])
    expect(opsiTitle([{ kode: '04', label: 'PT.', aktif: true }, { kode: '01', label: 'TN.', aktif: false }], 'TN.')).toEqual([
      { value: 'PT.', label: 'PT.' },
      { value: 'TN.', label: CD.nilaiLama('TN.') },
    ])
    expect(labelKode(telfax, '5')).toBe('OFFICE PHONE')
    expect(labelKode(telfax, '9')).toBe('9')
  })

  it('PIC Name: akun login aktif, nama kembar sekali; nama PIC lama yang bukan akun tetap terbaca', () => {
    const akun = [
      { loginId: 'UJI-AKUN-1', nama: 'UJI Kontak A', jabatan: 'Direktur' },
      { loginId: 'UJI-AKUN-2', nama: 'UJI Kontak B', jabatan: 'Staff' },
      { loginId: 'UJI-AKUN-3', nama: 'UJI Kontak A', jabatan: 'Manager' },
    ]
    expect(opsiAkun(akun, '')).toEqual([
      { value: 'UJI Kontak A', label: 'UJI Kontak A' },
      { value: 'UJI Kontak B', label: 'UJI Kontak B' },
    ])
    expect(opsiAkun(akun, 'UJI Kontak B')).toHaveLength(2)
    expect(opsiAkun(akun, 'UJI Kontak Lama')).toContainEqual({ value: 'UJI Kontak Lama', label: CD.nilaiLama('UJI Kontak Lama') })
  })

  it('PIC Position = JOB_POSITION akun yang dipilih; nama bukan akun = tidak diketahui', () => {
    const akun = [
      { loginId: 'UJI-AKUN-1', nama: 'UJI Kontak A', jabatan: 'Direktur' },
      { loginId: 'UJI-AKUN-2', nama: 'UJI Kontak B', jabatan: '' },
    ]
    expect(jabatanAkun(akun, 'UJI Kontak A')).toBe('Direktur')
    expect(jabatanAkun(akun, 'UJI Kontak B')).toBe('')
    expect(jabatanAkun(akun, 'UJI Kontak Lama')).toBeUndefined()
  })

  it('kode area: daftar kodehp aktif + Others; kode di luar daftar = Others diisi sendiri', () => {
    const kode: Pilihan[] = [
      { kode: '021', label: 'UJI AREA', aktif: true },
      { kode: '0542', label: '', aktif: false },
    ]
    expect(opsiKodeArea(kode, '')).toEqual([
      { value: '021', label: '021 - UJI AREA' },
      { value: KODE_LAIN, label: CD.lainnya },
    ])
    expect(opsiKodeArea(kode, '0542')).toContainEqual({ value: '0542', label: CD.nilaiLama('0542') })
    expect(kodeAreaLain(kode, '')).toBe(false)
    expect(kodeAreaLain(kode, '021')).toBe(false)
    expect(kodeAreaLain(kode, '0542')).toBe(false)
    expect(kodeAreaLain(kode, '0999')).toBe(true)
  })

  it('COUNTRY: kode = OLDID, atau ID NATION bila tanpa OLDID', () => {
    const negara = [
      { id: '100901', oldId: '001', nama: 'UJI NEGARA SATU', nationInitial: '' },
      { id: '100902', oldId: '', nama: 'UJI NEGARA BARU', nationInitial: '' },
    ]
    expect(negara.map(kodeNegara)).toEqual(['001', '100902'])
    expect(opsiNegara(negara, '777', 'UJI LAMA').at(-1)).toEqual({ value: '777', label: CD.nilaiLama('UJI LAMA') })
  })

  it('title di dalam nama: kata utuh, titik dan tanda baca tidak menghalangi', () => {
    const title: Pilihan[] = [
      { kode: '04', label: 'PT.', aktif: true },
      { kode: '05', label: 'CV.', aktif: true },
    ]
    expect(kataNama('P.T. UJI-Maju, Tbk')).toEqual(['PT', 'UJI', 'MAJU', 'TBK'])
    expect(titleDiNama('PT UJI Maju', title)).toBe('PT.')
    expect(titleDiNama('UJI Maju, CV', title)).toBe('CV.')
    expect(titleDiNama('P.T. UJI Maju', title)).toBe('PT.')
    expect(titleDiNama('UJI PTERA Maju', title)).toBe('')
    expect(titleDiNama('', title)).toBe('')
    expect(titleDiNama('UJI NY Trading', [...title, { kode: '02', label: 'NY.', aktif: false }])).toBe('')
  })

  it('title di nama ditolak; Edit nama lama yang tidak diubah tetap boleh', () => {
    const title: Pilihan[] = [{ kode: '04', label: 'PT.', aktif: true }]
    const isi = { ...isianDari(detail), nama: 'PT UJI BARU' }
    expect(periksa(isi, title)).toBe(CD.galatTitleDiNama('PT.'))
    expect(periksa(isi, title, 'pt uji baru')).toBeNull()
    expect(periksa({ ...isi, nama: 'UJI BARU' }, title)).toBeNull()
  })

  it('Copy Old: saring, catatan, dan ringkasan hasil', () => {
    const lama: OrgLama[] = [
      { id: 'A', idView: 'ORG-118', nama: 'UJI SATU', baru: true, pic: 2, nomor: 0, isi: [] },
      { id: 'B', idView: 'ORG-119', nama: 'UJI DUA', baru: false, pic: 0, nomor: 3, isi: ['PARENT_ID', 'NOTE'] },
    ]
    expect(saringLama(lama, 'dua').map((o) => o.id)).toEqual(['B'])
    expect(saringLama(lama, 'org-118').map((o) => o.id)).toEqual(['A'])
    expect(catatanLama(lama[0]!)).toBe('new organization; 2 PIC')
    expect(catatanLama(lama[1]!)).toBe('3 phone/fax; fill PARENT_ID, NOTE')
    expect(ringkasSalin([{ id: 'A', status: 'disalin', pesan: [] }, { id: 'B', status: 'gagal', pesan: [] }])).toEqual({
      disalin: 1, sudahAda: 0, ditolak: 0, gagal: 1,
    })
  })

  it('Process Copy dipecah per kelompok berurutan', () => {
    expect(potong([1, 2, 3, 4, 5], 2)).toEqual([[1, 2], [3, 4], [5]])
    expect(potong([], 20)).toEqual([])
  })

  it('halaman dan larik', () => {
    expect(jumlahHalaman(0, 50)).toBe(1)
    expect(jumlahHalaman(101, 50)).toBe(3)
    expect(gantiDi([1, 2, 3], 1, 9)).toEqual([1, 9, 3])
    expect(buangDi([1, 2, 3], 0)).toEqual([2, 3])
  })
})

describe('klien Company Detail', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sesuai backend; ID beralamat spasi di-encode', async () => {
    const panggil: [string, string][] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push([String(init.method), url])
        return Promise.resolve(new Response('{"daftar":[]}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar('uji', 2)
    await cariInduk('uji', detail.id)
    await periksaNama('PT UJI', '')
    await ambilHak()
    await ambilLama()
    await salinLama(['A'])
    await ambil(detail.id)
    await tambah(isianKosong())
    await ubah(detail.id, isianKosong())
    expect(panggil.map(([m, u]) => `${m} ${u.replace(/^https?:\/\/[^/]+/, '')}`)).toEqual([
      'GET /api/company-detail?q=uji&halaman=2',
      'GET /api/company-detail/induk?q=uji&kecuali=ASM-SFAGIS-WORK-ORG+ORG-115',
      'GET /api/company-detail/periksa-nama?nama=PT+UJI',
      'GET /api/company-detail/hak',
      'GET /api/company-detail/lama',
      'POST /api/company-detail/lama/salin',
      'GET /api/company-detail/ASM-SFAGIS-WORK-ORG%20ORG-115',
      'POST /api/company-detail',
      'PUT /api/company-detail/ASM-SFAGIS-WORK-ORG%20ORG-115',
    ])
  })
})
