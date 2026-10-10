// Blok Attachment File (tiket 67).

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import { LAMPIRAN as L } from '../labels'

const GRID = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\AttachmentGridReas.xml'
const VIEW = 'D:\\migrasi\\RNM\\NB FacIn\\Harness\\ViewPictureList.xml'
const baca = (...p: string[]) => readFileSync(join(__dirname, ...p), 'utf8').replace(/\r\n/g, '\n')
const SUMBER = baca('BlokLampiran.tsx')
const API = baca('..', 'api.ts')

describe('Attachment File - label korpus', () => {
  it.skipIf(!existsSync(GRID))('grid AttachmentGridReas: judul + Category / Count / Upload File / View File', () => {
    const xml = readFileSync(GRID, 'utf8')
    expect(xml).toContain(`<pyTitle>${L.judul}</pyTitle>`)
    for (const t of [L.kategori, L.jumlah, L.unggah, L.lihat]) expect(xml).toContain(`<pyValue>${t}</pyValue>`)
  })
  it.skipIf(!existsSync(VIEW))('popup ViewPictureList: File / Note / Upload Date / Delete', () => {
    const xml = readFileSync(VIEW, 'utf8')
    for (const t of [L.berkas, L.catatan, L.tanggal]) expect(xml).toContain(`<pyValue>${t}</pyValue>`)
    expect(xml).toContain(`<pyLabel>${L.hapus}</pyLabel>`)
  })
})

describe('BlokLampiran', () => {
  it('unggah banyak berkas: satu permintaan multipart per berkas (pola Treaty Contract Out)', () => {
    expect(SUMBER).toContain('type="file" multiple hidden')
    expect(SUMBER).toMatch(/for \(const f of daftar\) \{\s*setMengunggah\(kategori\)\s*await unggahLampiran\(caseId, kategori, f\)/)
    expect(API).toContain("isi.append('berkas', berkas)")
    expect(API).toContain("isi.append('kategori', kategori)")
    expect(API).toContain('return mintaFormulir<{ id: string }>(`${jalurLampiran(caseId)}/berkas`, isi)')
  })
  it('kategori lewat kueri (teks memuat "/"); unduh beridentitas; hapus hanya bila bolehHapus + konfirmasi', () => {
    expect(API).toContain('{ kueri: { kategori } }')
    expect(SUMBER).toContain('unduhBerkasBeridentitas(jalurIsiLampiran(caseId, b.id), b.namaFile)')
    expect(SUMBER).toContain('{v.bolehHapus && (')
    expect(SUMBER).toContain('<p>{T.yakinHapus(hapus.namaFile)}</p>')
    expect(SUMBER).toContain('disabled={!v.bolehUnggah || mengunggah !== \'\'}')
  })
  it('status unggah asinkron: menunggu / gagal + Ulangi (ADR-0008)', () => {
    expect(API).toContain("status: '' | 'menunggu' | 'gagal'")
    expect(API).toContain("/ulangi`, { metode: 'POST' })")
    expect(SUMBER).toContain("{b.status === 'gagal' && (")
    expect(SUMBER).toContain('onClick={() => void ulangi(b)}')
  })
  it('5 kategori per halaman, pager di atas tabel; tombol ikon ber-aria-label', () => {
    expect(SUMBER).toContain('export const UKURAN_KATEGORI_LAMPIRAN = 5')
    expect(SUMBER).toMatch(/<HalamanPopup halaman=\{kini\} ukuran=\{UKURAN_KATEGORI_LAMPIRAN\}[\s\S]{0,200}<table className="nbf-tabel nbf-lampiran__tabel">/)
    expect(SUMBER).toContain('{halamanIni.map((k) => (')
    expect(SUMBER).toContain('aria-label={`${L.unggah} - ${k.kategori}`}')
    expect(SUMBER).toContain('aria-label={`${L.lihat} - ${k.kategori}`}')
  })
  it('halaman case memasang blok di bawah riwayat akseptasi', () => {
    const HALAMAN = baca('..', 'pages', 'InwardFacultative.tsx')
    expect(HALAMAN).toMatch(/\{bolehSubmit && <BlokEmailType s=\{sub\} hanyaRiwayat=\{hanyaLihat\} \/>\}[\s\S]{0,400}<BlokLampiran caseId=\{kasus\.caseId\} \/>/)
  })
})
