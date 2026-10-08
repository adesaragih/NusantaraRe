// Panel Attachment — tombol `Upload file` (`Section/WorkAttachments.xml`),
// HIDUP sejak 8 Oktober 2026: unggahan masuk `M_ATTACHMENTTREATY_2`.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import PanelLampiran from './components/PanelLampiran'
import type { BarisKategoriLampiran } from './api'
import { LAMPIRAN } from './labels'
import { gabungBerkas, unggahBerurutan } from './unggahBerkas'

const AKAR = __dirname
const PANEL = readFileSync(join(AKAR, 'components', 'PanelLampiran.tsx'), 'utf8')
const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')
const berkasUji = (nama: string) => new File(['isi'], nama)

const KATEGORI: BarisKategoriLampiran[] = [
  { kode: '00002', nama: 'Approval Email', cacah: 1, dipastikan: true },
  { kode: '00008', nama: 'Letter of Acknowledgment / LOA', cacah: 0, dipastikan: true },
]
const tombolUpload = (html: string) => [...html.matchAll(/<button[^>]*>Upload<\/button>/g)].map((m) => m[0])

describe('⭐ Upload file — syarat tampil ekspor', () => {
  it('ViewState != 1 / RevisionState = 1 (bisaUnggah) + kontrak ber-ID: tombol HIDUP per kategori', () => {
    const html = renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="1002305" bisaUnggah />)
    const t = tombolUpload(html)
    expect(t).toHaveLength(2)
    for (const b of t) expect(b).not.toContain('disabled')
  })

  it('mode lihat tanpa revisi: tombol TIDAK dirender', () => {
    expect(tombolUpload(renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="1002305" />))).toHaveLength(0)
  })

  it('kontrak belum tersimpan: tombol mati, dengan alasannya', () => {
    const html = renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="" bisaUnggah />)
    for (const b of tombolUpload(html)) expect(b).toContain('disabled')
    expect(html).toContain('Simpan kontrak lebih dulu untuk mengunggah lampiran.')
  })

  it('Refresh tampil; Download All (`pyCondition never`) tidak', () => {
    const html = renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="1002305" />)
    expect(html).toContain('>Refresh</button>')
    expect(html).not.toContain('Download All')
  })
})

// ⭐ 8 Oktober 2026 — permintaan pemakai: bentuk unggah DISAMAKAN dengan menu
// Master Product Name Life, sebab tim memakai bentuk itu.
describe('⭐ modal Attach — bentuk Master Product Name Life', () => {
  it('pilihan digabung, nama kembar (tanpa beda huruf) dilewati, urutan tetap', () => {
    const a = berkasUji('Bordero.xlsx')
    const hasil = gabungBerkas([a], [berkasUji('bordero.XLSX'), berkasUji('LOA.pdf')])
    expect(hasil.map((f) => f.name)).toEqual(['Bordero.xlsx', 'LOA.pdf'])
    expect(hasil[0]).toBe(a)
  })

  it('diunggah SATU PER SATU berurutan; satu gagal tidak menghentikan sisanya', async () => {
    const urutan: string[] = []
    const mulai: string[] = []
    const gagal = await unggahBerurutan(
      [berkasUji('a.pdf'), berkasUji('b.exe'), berkasUji('c.pdf')],
      async (f) => {
        urutan.push(f.name)
        if (f.name === 'b.exe') throw new Error('Jenis berkas .exe tidak dikenal')
        await Promise.resolve()
      },
      (f, i) => {
        mulai.push(`${String(i + 1)}/3 ${f.name}`)
      },
    )
    expect(urutan).toEqual(['a.pdf', 'b.exe', 'c.pdf'])
    expect(mulai).toEqual(['1/3 a.pdf', '2/3 b.exe', '3/3 c.pdf'])
    expect(gagal.map((g) => g.berkas.name)).toEqual(['b.exe'])
  })

  it('kotak seret-lepas, Remove per berkas, progres, daftar gagal', () => {
    expect(PANEL).toContain("className={'trin__unggah' + (seret ? ' trin__unggah--seret' : '')}")
    expect(PANEL).toContain('onDrop={(e) => {')
    expect(PANEL).toContain('{LAMPIRAN.buangPilihan}')
    expect(PANEL).toContain('{LAMPIRAN.mengunggah} {proses}')
    expect(PANEL).toContain('className="trin__unggah-gagal"')
    expect(LAMPIRAN.seretBerkas).toBe('Drag and drop files here, or click to choose files')
    expect(LAMPIRAN.buangPilihan).toBe('Remove')
    expect(LAMPIRAN.mengunggah).toBe('Uploading')
    for (const k of ['trin__unggah', 'trin__unggah--seret', 'trin__unggah-daftar', 'trin__unggah-gagal']) {
      expect(CSS).toMatch(new RegExp('\\.treatyin \\.' + k + '[ ,]'))
    }
  })

  it('SATU berkas per permintaan; berkas yang DITOLAK (berhasil false) tetap dihitung gagal', () => {
    expect(PANEL).toContain('await unggahLampiran(idKontrak, kode, [f])')
    expect(PANEL).toContain('if (b !== undefined && !b.berhasil) throw new Error(b.pesan)')
    // Yang gagal tinggal di pilihan; semua berhasil → modal tertutup.
    expect(PANEL).toContain('setDipilih(gagal.map((g) => g.berkas))')
    expect(PANEL).toContain('if (gagal.length === 0) setUnggahKe(null)')
  })

  // ⛔ Ralat 8 Oktober 2026: unduh lewat fetch beridentitas, View Office
  // Online di bingkai popup — nol `window.open` (`unduhdokumen.test.ts`).
  it('View File: unduh lewat fetch beridentitas, Office di bingkai popup', () => {
    expect(PANEL).toContain('unduhLampiran(idKontrak, b.id, b.namaBerkas)')
    expect(PANEL).toContain('bingkai.current.src = j.url')
    expect(PANEL).not.toMatch(/window\.open\(|location\.(href|assign|replace)\b/)
  })

  it('Attach tidak dimatikan saat pilihan kosong — kalimat penolakan datang dari backend', () => {
    expect(PANEL).toContain('disabled={sibuk} onClick={lampirkan}')
    expect(PANEL).toContain('await unggahLampiran(idKontrak, kode, [])')
  })
})
