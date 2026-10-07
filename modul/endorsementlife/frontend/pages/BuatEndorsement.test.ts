import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { kegagalanDari } from '../../../../inti/frontend/klien'
import { pesanPenolakanEDM } from '../api'
import { BUAT_EDM, OPSI_EDM_TYPE_BATAL, OPSI_EDM_TYPE_PERUBAHAN_DATA } from '../labels'
import { subJenisKasus } from '../tampilan'
import {
  alasanBelumSubmit,
  bolehSubmitEDM,
  EDM_TYPE_BATAL_BAWAAN,
  edmTypeSah,
  subJenisSah,
  subJenisSesudahGantiEdmType,
  type KeadaanSubmitEDM,
} from './BuatEndorsement'

// Tampilan Create Addendum dirapikan - keputusan work owner 07-10-2026 (logika tidak berubah).

const SUMBER = readFileSync(join(__dirname, 'BuatEndorsement.tsx'), 'utf8')
const CSS = readFileSync(join(__dirname, '..', 'endorsementlife.css'), 'utf8')
/** Sumber tanpa baris komentar - komentar yang MENERANGKAN aturan tidak dihitung sebagai kode. */
const KODE = SUMBER.split(/\r?\n/)
  .filter((b) => !/^\s*(\/\/|\*|\/\*\*)/.test(b))
  .join('\n')

/** Potongan sumber di antara dua penanda (keduanya harus ada, berurutan). */
function antara(a: string, b: string): string {
  const i = SUMBER.indexOf(a)
  const j = SUMBER.indexOf(b, i)
  expect(i).toBeGreaterThanOrEqual(0)
  expect(j).toBeGreaterThan(i)
  return SUMBER.slice(i, j)
}

describe('struktur Create Addendum', () => {
  it('baris 1: Policy No, EDM Type, EDM Date; baris 2: Description sendiri', () => {
    const baris1 = antara('<div className="edm-buat__baris">', '</div>')
    expect(baris1).toContain('label={BUAT_EDM.policyNo}')
    expect(baris1).toContain('label={BUAT_EDM.edmType}')
    expect(baris1).toContain('label={BUAT_EDM.edmDate}')
    expect(baris1).not.toContain('BUAT_EDM.description')
    expect(baris1.indexOf('policyNo')).toBeLessThan(baris1.indexOf('edmType'))
    expect(baris1.indexOf('edmType')).toBeLessThan(baris1.indexOf('edmDate'))
    const baris2 = antara('<div className="edm-buat__lebar">', '</div>')
    expect(baris2).toContain('label={BUAT_EDM.description}')
    expect(baris2).toContain('baris={3}')
    expect(CSS).toMatch(/\.edm-buat__baris \{[^}]*grid-template-columns: repeat\(3, minmax\(0, 1fr\)\)/)
    // Selebar penuh area konten (07-10-2026): tanpa batas lebar.
    expect(CSS).toMatch(/\.edm-buat \{[^}]*width: 100%/)
    expect(CSS).not.toMatch(/\.edm-buat \{[^}]*max-width/)
    expect(CSS).toMatch(/@media \(max-width: 720px\) \{\s*\.edm-buat__baris \{\s*grid-template-columns: minmax\(0, 1fr\)/)
  })

  it('pesan memeriksa dan pesan gerbang di antara isian dan footer', () => {
    const iIsian = SUMBER.indexOf('<div className="edm-buat__lebar">')
    const iStatus = SUMBER.indexOf('UMUM_EDM.memeriksa')
    const iPesan = SUMBER.indexOf('pesanGerbang.map(')
    const iKaki = SUMBER.indexOf('<footer className="edm-buat__kaki">')
    expect(iIsian).toBeLessThan(iStatus)
    expect(iStatus).toBeLessThan(iPesan)
    expect(iPesan).toBeLessThan(iKaki)
    expect(SUMBER).toContain('className="alert alert--info edm-buat__status" role="status"')
    expect(SUMBER).toContain('className="alert alert--error edm-buat__pesan" role="alert"')
  })

  it('footer berisi Cancel (onTutup) lalu Submit', () => {
    const kaki = antara('<footer className="edm-buat__kaki">', '</footer>')
    expect(kaki).toContain('<button type="button" className="btn btn--ghost" onClick={onTutup}>')
    expect(kaki).toContain('{BUAT_EDM.cancel}')
    expect(kaki.indexOf('BUAT_EDM.cancel')).toBeLessThan(kaki.indexOf('BUAT_EDM.submit'))
    expect(BUAT_EDM.cancel).toBe('Cancel')
    expect(CSS).toMatch(/\.edm-buat__kaki \{[^}]*justify-content: flex-end/)
  })
})

describe('Submit selalu tampil, nonaktif selama gerbang belum lolos', () => {
  it('disabled={!bolehSubmit}; tidak lagi dibungkus {bolehSubmit && ...}; kirim() yang sama', () => {
    const kaki = antara('<footer className="edm-buat__kaki">', '</footer>')
    expect(kaki).toContain('disabled={!bolehSubmit}')
    expect(kaki).toContain('onClick={() => void kirim()}')
    expect(SUMBER).not.toContain('{bolehSubmit && (')
    expect(SUMBER).toContain('const bolehSubmit = bolehSubmitEDM(keadaan)')
    expect(SUMBER).toContain('keputusan work owner 07-10-2026')
    expect(SUMBER).toContain('Pega b4226 VIS')
  })

  // Gerbang kini dipanggil TANPA EDM Type, jadi `boleh` server selalu false - Submit tidak bergantung padanya.
  const dasar: KeadaanSubmitEDM = {
    policyNo: 'UJI-PL-1', edmType: '1', kelayakan: { pesan: [], boleh: false }, penolakan: [], memeriksa: false, mengirim: false,
    perubahanData: OPSI_EDM_TYPE_PERUBAHAN_DATA[0] ?? '', batal: '',
  }

  it('bolehSubmitEDM: Policy No, hasil gerbang tanpa pesan, EDM Type 1/3, tidak sibuk, tanpa penolakan', () => {
    expect(bolehSubmitEDM(dasar)).toBe(true)
    expect(bolehSubmitEDM({ ...dasar, edmType: '3', perubahanData: '', batal: EDM_TYPE_BATAL_BAWAAN })).toBe(true)
    expect(bolehSubmitEDM({ ...dasar, policyNo: ' ' })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, kelayakan: null })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, kelayakan: { pesan: ['UJI-PESAN'], boleh: false } })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, edmType: '' })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, edmType: '2' })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, edmType: 'toString' })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, memeriksa: true })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, mengirim: true })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, penolakan: ['UJI-TOLAK'] })).toBe(false)
    expect(edmTypeSah('1') && edmTypeSah('3') && !edmTypeSah('2')).toBe(true)
    expect(KODE).not.toContain('.boleh')
  })

  it('teks bantuan menyebut apa yang masih kurang', () => {
    expect(alasanBelumSubmit({ ...dasar, policyNo: '  ' })).toBe(BUAT_EDM.bantuan.policyNo)
    expect(alasanBelumSubmit({ ...dasar, edmType: '' })).toBe(BUAT_EDM.bantuan.edmType)
    expect(alasanBelumSubmit({ ...dasar, edmType: '2' })).toBe(BUAT_EDM.bantuan.edmType)
    expect(alasanBelumSubmit({ ...dasar, memeriksa: true })).toBe(BUAT_EDM.bantuan.memeriksa)
    expect(alasanBelumSubmit({ ...dasar, kelayakan: null })).toBe(BUAT_EDM.bantuan.menunggu)
    expect(alasanBelumSubmit({ ...dasar, kelayakan: { pesan: ['UJI-PESAN'], boleh: false } })).toBe(BUAT_EDM.bantuan.pesan)
    expect(alasanBelumSubmit({ ...dasar, penolakan: ['UJI-TOLAK'] })).toBe(BUAT_EDM.bantuan.penolakan)
    expect(alasanBelumSubmit({ ...dasar, mengirim: true })).toBe(BUAT_EDM.bantuan.mengirim)
  })

  it('gerbang lolos → tanpa teks bantuan (Submit aktif)', () => {
    expect(alasanBelumSubmit(dasar)).toBe('')
    expect(alasanBelumSubmit({ ...dasar, edmType: '3', perubahanData: '', batal: EDM_TYPE_BATAL_BAWAAN })).toBe('')
    expect(SUMBER).toContain("const alasan = bolehSubmit ? '' : alasanBelumSubmit(keadaan)")
  })
})

describe('gerbang hanya saat Policy No berubah (07-10-2026)', () => {
  const efek = (() => {
    const i = SUMBER.indexOf('  useEffect(() => {')
    const j = SUMBER.indexOf('  }, [policyNo])', i)
    expect(i).toBeGreaterThanOrEqual(0)
    expect(j).toBeGreaterThan(i)
    return SUMBER.slice(i, j + '  }, [policyNo])'.length)
  })()

  it('EDM Type / EDM Date / Description tidak memicu cekKelayakan', () => {
    expect(efek.endsWith('}, [policyNo])')).toBe(true)
    expect(SUMBER).not.toContain('}, [policyNo, edmType')
    expect(KODE.split('cekKelayakan(').length - 1).toBe(1)
    expect(KODE.split('useEffect(').length - 1).toBe(1)
  })

  it('Policy No kosong → tanpa panggilan server, tanpa pesan', () => {
    const awal = efek.indexOf("if (np === '') {")
    const kosong = efek.slice(awal, efek.indexOf('return', awal) + 'return'.length)
    expect(kosong).toContain('setKelayakan(null)')
    expect(kosong).not.toContain('cekKelayakan')
    expect(efek.indexOf("if (np === '') {")).toBeLessThan(efek.indexOf('cekKelayakan('))
  })

  it('Policy No terisi → cekKelayakan dengan EDM Type KOSONG, debounce 300ms', () => {
    expect(efek).toContain('const np = policyNo.trim()')
    expect(efek).toContain("cekKelayakan(np, '')")
    expect(efek).toContain('}, 300)')
  })
})

describe('penolakan Submit berpesan tampil sebagai daftar pesan gerbang', () => {
  it('422 {galat} → satu pesan per baris; selain itu null', () => {
    const tolak = kegagalanDari(422, JSON.stringify({ galat: 'UJI-PESAN-A\nUJI-PESAN-B\n' }))
    expect(pesanPenolakanEDM(tolak)).toEqual(['UJI-PESAN-A', 'UJI-PESAN-B'])
    expect(pesanPenolakanEDM(kegagalanDari(409, JSON.stringify({ galat: 'UJI-BENTROK' })))).toBeNull()
    expect(pesanPenolakanEDM(kegagalanDari(422, '{}'))).toBeNull()
    expect(pesanPenolakanEDM(new Error('UJI-JARINGAN'))).toBeNull()
  })

  it('kirim(): penolakan → daftar pesan (alert merah yang sama), bukan galat umum; isian berubah menghapusnya', () => {
    expect(SUMBER).toContain('const pesan = pesanPenolakanEDM(e)')
    expect(SUMBER).toContain('if (pesan !== null) setPenolakan(pesan)')
    expect(SUMBER).toContain('else setGalat(e)')
    expect(SUMBER).toContain('const pesanGerbang = [...(kelayakan?.pesan ?? []), ...penolakan]')
    expect(SUMBER).toContain('const h = await buatKasus({')
    const ubah = SUMBER.slice(SUMBER.indexOf('const ubah = '), SUMBER.indexOf('const ubah = ') + 120)
    expect(ubah).toContain('setPenolakan([])')
  })
})

describe('penanda wajib sesuai gerbang', () => {
  it('hanya Policy No (masukan gerbang) dan EDM Type (syarat Submit 1/3); EDM Date dan Description tidak', () => {
    expect(SUMBER).toContain('<Field label={BUAT_EDM.policyNo} value={policyNo} onChange={ubah(setPolicyNo)} autoFocus required />')
    expect(SUMBER).toContain('<Pilih label={BUAT_EDM.edmType} value={edmType} onChange={gantiEdmType} opsi={OPSI} required />')
    expect(SUMBER).toContain('<FieldTanggal label={BUAT_EDM.edmDate} value={edmDate} onChange={ubah(setEdmDate)} />')
    expect(SUMBER).toContain('<Area label={BUAT_EDM.description} value={deskripsi} onChange={ubah(setDeskripsi)} baris={3} />')
    expect(SUMBER.match(/ required \/>/g)?.length).toBe(2)
  })
})

describe('sub-jenis EDM Type (OQ-EDM-005, 07-10-2026)', () => {
  const dasar: KeadaanSubmitEDM = {
    policyNo: 'UJI-PL-1', edmType: '1', kelayakan: { pesan: [], boleh: false }, penolakan: [], memeriksa: false, mengirim: false,
    perubahanData: '', batal: '',
  }

  it('opsi = layar Pega; Batal bawaan Batal Internal tanpa Choose', () => {
    expect(OPSI_EDM_TYPE_PERUBAHAN_DATA).toEqual([
      'Perubahan data (No Polis)', 'Perubahan data (No Sertifikat)', 'Perubahan data (Nama Peserta)',
      'Perubahan data (DOB)', 'Perubahan data (Begin & Expired Date)', 'Perubahan data (UP & Premi)',
    ])
    expect(OPSI_EDM_TYPE_BATAL).toEqual(['Batal Internal', 'Batal Eksternal'])
    expect(EDM_TYPE_BATAL_BAWAAN).toBe('Batal Internal')
    expect(BUAT_EDM.pilih).toBe('Choose')
    const batal = antara('{edmType === \'3\' && (', '</select>')
    expect(batal).toContain('<select')
    expect(batal).not.toContain('<option value="">')
    expect(batal).toContain('<span className="field__req">*</span>')
  })

  it('dropdown muncul sesuai EDM Type, di bawah EDM Type', () => {
    const ubah = antara("{edmType === '1' && (", "{edmType === '3' && (")
    expect(ubah).toContain('label={BUAT_EDM.edmTypePerubahanData}')
    expect(ubah).toContain('kosong={BUAT_EDM.pilih}')
    expect(ubah).toContain('required')
    expect(ubah).not.toContain('edmTypeBatal')
    expect(SUMBER).toContain('label={BUAT_EDM.edmTypePerubahanData}')
    expect(SUMBER.indexOf('label={BUAT_EDM.edmType} ')).toBeLessThan(SUMBER.indexOf("{edmType === '1' && ("))
    expect(SUMBER.indexOf("{edmType === '3' && (")).toBeLessThan(SUMBER.indexOf('<div className="edm-buat__lebar">'))
    expect(CSS).toMatch(/\.edm-buat__sub \{[^}]*grid-column: 2/)
  })

  it('ganti EDM Type → sub-jenis yang tidak relevan dikosongkan; Batal ke bawaan', () => {
    expect(subJenisSesudahGantiEdmType('1')).toEqual({ perubahanData: '', batal: '' })
    expect(subJenisSesudahGantiEdmType('3')).toEqual({ perubahanData: '', batal: 'Batal Internal' })
    expect(subJenisSesudahGantiEdmType('')).toEqual({ perubahanData: '', batal: '' })
    const ganti = SUMBER.slice(SUMBER.indexOf('const gantiEdmType = '), SUMBER.indexOf('const gantiEdmType = ') + 260)
    expect(ganti).toContain('subJenisSesudahGantiEdmType(v)')
    expect(ganti).toContain('setPerubahanData(sub.perubahanData)')
    expect(ganti).toContain('setBatal(sub.batal)')
  })

  it('Submit nonaktif tanpa sub-jenis yang sah; teks bantuannya', () => {
    expect(bolehSubmitEDM(dasar)).toBe(false)
    expect(alasanBelumSubmit(dasar)).toBe(BUAT_EDM.bantuan.perubahanData)
    expect(bolehSubmitEDM({ ...dasar, perubahanData: 'Perubahan data (DOB)' })).toBe(true)
    expect(bolehSubmitEDM({ ...dasar, perubahanData: 'UJI-BUKAN-OPSI' })).toBe(false)
    expect(bolehSubmitEDM({ ...dasar, edmType: '3' })).toBe(false)
    expect(alasanBelumSubmit({ ...dasar, edmType: '3' })).toBe(BUAT_EDM.bantuan.batal)
    expect(bolehSubmitEDM({ ...dasar, edmType: '3', batal: 'Batal Eksternal' })).toBe(true)
    expect(subJenisSah('3', 'Perubahan data (DOB)', 'Batal Internal')).toBe(false)
    expect(subJenisSah('1', 'Perubahan data (DOB)', 'Batal Internal')).toBe(false)
  })

  it('sub-jenis tidak memanggil gerbang; payload buatKasus memuatnya sesuai EDM Type', () => {
    expect(KODE.split('cekKelayakan(').length - 1).toBe(1)
    expect(SUMBER).toContain('}, [policyNo])')
    expect(SUMBER).toContain("edmTypePerubahanData: edmType === '1' ? perubahanData : '',")
    expect(SUMBER).toContain("edmTypeBatal: edmType === '3' ? batal : '',")
  })

  it('layar kasus: sub-jenis baca-saja sesuai EDM Type', () => {
    expect(subJenisKasus({ edmType: '1', edmTypePerubahanData: 'Perubahan data (DOB)', edmTypeBatal: '' })).toEqual([
      [BUAT_EDM.edmTypePerubahanData, 'Perubahan data (DOB)'],
    ])
    expect(subJenisKasus({ edmType: '3', edmTypePerubahanData: '', edmTypeBatal: 'Batal Eksternal' })).toEqual([
      [BUAT_EDM.edmTypeBatal, 'Batal Eksternal'],
    ])
    expect(subJenisKasus({ edmType: '2' })).toEqual([])
    // Kolom inbox EDM Type Batal dihapus 07-10-2026; nilainya tetap tampil di layar kasus.
    const inbox = readFileSync(join(__dirname, 'InboxEndorsementLife.tsx'), 'utf8')
    expect(inbox).not.toContain('edmTypeBatal')
  })
})
