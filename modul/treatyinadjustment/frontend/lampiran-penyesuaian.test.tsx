// Panel Attachment layar Adjustment — HIDUP dan sama dengan Treaty In
// (8 Oktober 2026). Render STATIS (`react-dom/server`): nol DOM, nol jaringan.
//
// ⚠️ Render statis tidak menjalankan efek dan tidak dapat menekan tombol:
// isi grid diuji lewat `PanelLampiranIsi` (data sebagai props), dan modal
// `View File` lewat `ModalLihatBerkas` — keadaan "terbuka"-nya dijaga dari
// sumber panel.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisKategoriLampiran, BarisLampiranWarisan } from './api'
import { berkasKategori, ModalLihatBerkas, PanelLampiranIsi } from './komponen/PanelLampiranPenyesuaian'
import { gabungBerkas, unggahBerurutan } from './komponen/unggahBerkas'
import { LAMPIRAN_KELOLA } from './labelsLampiran'

const AKAR = __dirname
const PANEL = readFileSync(join(AKAR, 'komponen', 'PanelLampiranPenyesuaian.tsx'), 'utf8')
const API = readFileSync(join(AKAR, 'api.ts'), 'utf8')
const HALAMAN = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
const CSS = readFileSync(join(AKAR, 'treatyinadjustment.css'), 'utf8')

const KATEGORI: BarisKategoriLampiran[] = [
  { kode: '00002', nama: 'Approval Email', cacah: 2, dipastikan: true },
  { kode: '00008', nama: 'Letter of Acknowledgment / LOA', cacah: 0, dipastikan: true },
]
const BERKAS: BarisLampiranWarisan[] = [
  {
    id: 'L1', kodeKategori: '00002', namaKategori: 'Approval Email', namaBerkas: 'setuju.pdf', jenisMime: 'pdf',
    idSimpanan: 'S1', diunggah: '08-10-2026', pengunggah: 'ADE',
  },
  {
    id: 'L2', kodeKategori: '00002', namaKategori: 'Approval Email', namaBerkas: 'hitung.xlsx', jenisMime: 'xlsx',
    idSimpanan: 'S2', diunggah: '08-10-2026', pengunggah: 'ADE',
  },
]
const ID = '1000080/R01'

/** Tag pembuka tombol ber-`title` tertentu (sel ikon Upload/View). */
const tombol = (html: string, judul: string) =>
  [...html.matchAll(new RegExp(`<button[^>]*title="${judul}"[^>]*>`, 'g'))].map((m) => m[0])

const panel = (p: Partial<Parameters<typeof PanelLampiranIsi>[0]> = {}) =>
  renderToStaticMarkup(<PanelLampiranIsi kategori={KATEGORI} berkas={BERKAS} idKontrak={ID} {...p} />)

describe('panel Attachment Adjustment — bentuk Treaty In', () => {
  it('Download All lalu Refresh tampil, hidup, di deret kanan atas grid', () => {
    const html = panel()
    const i = html.indexOf('tria__lampiran-aksi')
    expect(i).toBeGreaterThan(0)
    const blok = html.slice(i, html.indexOf('</div>', i))
    expect(blok.indexOf(LAMPIRAN_KELOLA.unduhSemua)).toBeGreaterThan(0)
    expect(blok.indexOf(LAMPIRAN_KELOLA.segarkan)).toBeGreaterThan(blok.indexOf(LAMPIRAN_KELOLA.unduhSemua))
    expect(blok).not.toContain('disabled')
    // Spanduk aturan nama berkas sebelum deret tombol, grid sesudahnya.
    expect(html.indexOf('tria__spanduk')).toBeLessThan(i)
    expect(html.indexOf('tria__lampiran-grid')).toBeGreaterThan(i)
  })

  it('keempat kolom, tiga terakhir rata tengah, lebar 64/12/12/12', () => {
    const html = panel()
    for (const k of ['Category', 'Count', 'Upload file', 'View File']) expect(html).toContain(`>${k}</th>`)
    expect(html.match(/class="tria__lampiran-tengah"/g)).toHaveLength(3)
    expect(html).toContain('width:64%')
    expect(html.match(/width:12%/g)).toHaveLength(3)
  })

  it('⭐ Upload HIDUP per kategori bila bisaUnggah (ViewState != 1 / RevisionState = 1) dan ber-ID', () => {
    const t = tombol(panel({ bisaUnggah: true }), LAMPIRAN_KELOLA.unggah)
    expect(t).toHaveLength(KATEGORI.length)
    for (const x of t) expect(x).not.toContain('disabled')
  })

  it('⛔ mode lihat tanpa revisi: Upload TIDAK dirender; View tetap ada', () => {
    const html = panel({ bisaUnggah: false })
    expect(tombol(html, LAMPIRAN_KELOLA.unggah)).toHaveLength(0)
    expect(tombol(html, LAMPIRAN_KELOLA.lihatBerkas)).toHaveLength(KATEGORI.length)
  })

  it('⛔ draf belum tersimpan: Upload MATI dengan kalimat "simpan dulu" Treaty In', () => {
    const html = panel({ bisaUnggah: true, draf: true, idKontrak: '1000080' })
    const t = tombol(html, LAMPIRAN_KELOLA.unggah)
    expect(t).toHaveLength(KATEGORI.length)
    for (const x of t) expect(x).toContain('disabled')
    expect(html).toContain(LAMPIRAN_KELOLA.simpanDulu)
    expect(LAMPIRAN_KELOLA.simpanDulu).toBe('Save the contract first to upload attachments.')
    // Download All tetap hidup — lampiran asal boleh diunduh.
    expect(html).not.toMatch(/<button[^>]*disabled[^>]*>Download All</)
  })

  it('pembacaan pertama: tombol panel mati sampai data datang', () => {
    const html = panel({ kategori: [], berkas: [], memuat: true })
    expect(html).toMatch(/<button[^>]*disabled[^>]*>Download All</)
    expect(html).toMatch(/<button[^>]*disabled[^>]*>Refresh</)
  })

  it('Count TIDAK diformat dan nama kategori tampil apa adanya', () => {
    const html = panel()
    expect(html).toContain('<td>Approval Email</td><td>2</td>')
  })
})

describe('modal View File — ShowAttachmentTreaty', () => {
  const modal = (p: Partial<Parameters<typeof ModalLihatBerkas>[0]> = {}) =>
    renderToStaticMarkup(
      <ModalLihatBerkas
        kategori={KATEGORI}
        berkas={berkasKategori(BERKAS, KATEGORI[0] ?? null)}
        idKontrak={ID}
        bolehTulis
        bisaUnggah
        bolehGantiKategori
        gantiKategori={false}
        kategoriBaru={{}}
        sibuk={false}
        galatPanel=""
        onTutup={() => undefined}
        onGantiKategori={() => undefined}
        onSimpanKategori={() => undefined}
        onPilihKategori={() => undefined}
        onUnduh={() => undefined}
        onOffice={() => undefined}
        onHapus={() => undefined}
        {...p}
      />,
    )

  it('⭐ tombol View membuka modal ini (keadaan panel), isinya DUA kolom', () => {
    expect(PANEL).toContain('setBerkasDilihat(k.kode)')
    expect(PANEL).toMatch(/berkasDilihat !== null && \(\s*<ModalLihatBerkas/)
    const html = modal()
    expect(html).toContain('class="modal')
    expect(html).toContain(LAMPIRAN_KELOLA.judulLihatBerkas)
    expect(html).toContain('>File Name</th>')
    expect(html).toContain('>Type</th>')
    expect(html).not.toContain('>Uploaded</th>')
  })

  it('nama berkas = tautan unduh; View Office Online hanya untuk berkas kantor', () => {
    const html = modal()
    expect(html).toMatch(/class="tria__tautan"[^>]*>setuju\.pdf</)
    expect(html).toMatch(/class="tria__tautan"[^>]*>hitung\.xlsx</)
    expect(html.match(/View Office Online/g)).toHaveLength(1)
  })

  it('Delete + Change Category bila boleh; hilang di mode lihat / draf / status akhir', () => {
    const hidup = modal()
    expect(hidup.match(/>Delete</g)).toHaveLength(2)
    expect(hidup).toContain('>Change Category<')
    expect(modal({ bisaUnggah: false })).not.toContain('>Delete<')
    const draf = modal({ bolehTulis: false })
    expect(draf).not.toContain('>Delete<')
    expect(draf).not.toContain('>Change Category<')
    expect(modal({ bolehGantiKategori: false })).not.toContain('>Change Category<')
  })

  it('CARI30 = 1: tombol berganti Save dan kategori menjadi dropdown', () => {
    const html = modal({ gantiKategori: true })
    expect(html).toContain('>Save<')
    expect(html).not.toContain('>Change Category</button>')
    expect(html.match(/<select/g)).toHaveLength(2)
  })

  it('syarat Change Category = StatusAkseptasi bukan Resolve Complete / Decline', () => {
    expect(PANEL).toContain("statusAkseptasi !== 'Resolve Complete' && statusAkseptasi !== 'Decline'")
  })
})

describe('kabel ke layar Adjustment dan ke rute Treaty In', () => {
  it('⭐ layar memasang panel di wadah tersendiri dengan syarat ekspor', () => {
    expect(HALAMAN).toContain('<div className="tria__inbox-lampiran">')
    expect(HALAMAN).toContain("bisaUnggah={mode !== '1' || p.baru.medan.RevisionState === '1'}")
    expect(HALAMAN).toContain("statusAkseptasi={p.baru.medan.StatusAkseptasi ?? ''}")
    expect(HALAMAN).toContain('draf={draf !== undefined}')
    // Draf menampilkan lampiran asalnya; sesudah Save, pengenal baru.
    expect(HALAMAN).toContain('idKontrak={draf !== undefined ? p.idAsal : p.id}')
    // Attachment → deret tombol → History, ketiganya di wadah yang sama.
    const i = HALAMAN.indexOf('<div className="tria__inbox-lampiran">')
    expect(HALAMAN.indexOf('<DeretTombol', i)).toBeGreaterThan(i)
    expect(HALAMAN.indexOf('<PanelRiwayat', i)).toBeGreaterThan(HALAMAN.indexOf('<DeretTombol', i))
  })

  it('pengenal bergaris miring dikodekan; nol impor modul Treaty In', () => {
    expect(API).toContain('${PREFIX_TREATYIN}/kontrak/${encodeURIComponent(idKontrak)}/lampiran')
    for (const fn of ['unggahLampiran', 'hapusLampiran', 'ubahKategoriLampiran', 'ambilTautanLampiran', 'unduhLampiran', 'unduhSemuaLampiran']) {
      expect(API).toContain(`export async function ${fn}(`)
    }
    expect(PANEL).not.toMatch(/from '.*modul\/treatyin\//)
    expect(PANEL).not.toMatch(/from '.*\.\.\/\.\.\/treatyin\//)
  })

  it('CSS: wadah, spanduk persegi, aksi kanan, ikon berwarna token', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.tria__inbox-lampiran \{[^}]*margin-top: 16px;/)
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.tria__lampiran > \.tria__spanduk \{[^}]*border-radius: 0;/)
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.tria__lampiran > \.tria__lampiran-aksi \{\s*justify-content: flex-end;/)
    expect(tanpaKomentar).toMatch(/\.tria__lampiran-ikon--unggah \{\s*color: var\(--success-text\);/)
    expect(tanpaKomentar).toMatch(/\.tria__lampiran-ikon--lihat \{\s*color: var\(--accent\);/)
  })
})

describe('unggah berurutan — salinan bentuk Treaty In', () => {
  const f = (n: string) => new File(['isi'], n)

  it('pilihan digabung, nama kembar (tanpa beda huruf) dilewati', () => {
    expect(gabungBerkas([f('a.pdf')], [f('A.PDF'), f('b.pdf')]).map((x) => x.name)).toEqual(['a.pdf', 'b.pdf'])
  })

  it('satu gagal tidak menghentikan sisanya', async () => {
    const urut: string[] = []
    const gagal = await unggahBerurutan([f('a'), f('b'), f('c')], (x) => {
      urut.push(x.name)
      return x.name === 'b' ? Promise.reject(new Error('tolak')) : Promise.resolve()
    })
    expect(urut).toEqual(['a', 'b', 'c'])
    expect(gagal.map((g) => g.berkas.name)).toEqual(['b'])
  })
})
