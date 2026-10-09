// Lampiran "Reas" kasus - `AttachmentGridReas` / `ReasViewAttachment` / `AttachContentGIS` korpus NB FacIn, dipakai
// NB dan EDM Treaty In (keputusan work owner 08-10-2026): label berbukti baris korpus, rute klien, penampil Office,
// dan letak panel di layar kasus.

import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { BahasaUI } from '../components/ui/bahasaUI'
import { PERAN } from '../labels'
import { TEKS_UI } from '../lib/teksUI'
import {
  ambilDokumen,
  ambilGrid,
  ambilIsiDokumen,
  hapusDokumen,
  tautanOffice,
  unggahDokumen,
  type DokumenReas,
  type KategoriReas,
} from './api'
import { LAMPIRAN_REAS } from './labels'
import { GridKategori } from './PanelLampiranReas'
import {
  bisaViewOffice,
  gabungBerkas,
  jenisViewOnline,
  mimeViewOnline,
  PARAM_PENAMPIL,
  PENAMPIL_OFFICE,
} from './penampil'

const KORPUS = 'D:\\XML\\RNM_BRD\\NB FacIn'
const adaKorpus = existsSync(KORPUS)
const APLIKASI = join(__dirname, '..', '..', '..')

/** Satu baris (1-based) berkas korpus, tanpa indentasi. */
function baris(relatif: string, nomor: number): string {
  return (readFileSync(join(KORPUS, relatif), 'utf8').split('\n')[nomor - 1] ?? '').trim()
}

const panel = () => readFileSync(join(__dirname, 'PanelLampiranReas.tsx'), 'utf8')

describe.skipIf(!adaKorpus)('label lampiran Reas berbukti barisnya', () => {
  it('AttachmentGridReas, ReasViewAttachment, AttachContentGIS', () => {
    const G = 'Section\\AttachmentGridReas.xml'
    const V = 'Section\\ReasViewAttachment.xml'
    const F = 'FlowAction\\AttachContentGIS.xml'
    for (const [berkas, nomor, mau] of [
      [G, 702, `<pyTitle>${LAMPIRAN_REAS.judul}</pyTitle>`],
      [G, 1340, `<pyValue>${LAMPIRAN_REAS.category}</pyValue>`],
      [G, 1493, `<pyValue>${LAMPIRAN_REAS.count}</pyValue>`],
      [G, 1643, `<pyValue>${LAMPIRAN_REAS.uploadFile}</pyValue>`],
      [G, 1791, `<pyValue>${LAMPIRAN_REAS.viewFile}</pyValue>`],
      [V, 1363, `<pyValue>${LAMPIRAN_REAS.file}</pyValue>`],
      [V, 1618, `<pyValue>${LAMPIRAN_REAS.note}</pyValue>`],
      [V, 1762, `<pyValue>${LAMPIRAN_REAS.uploadDate}</pyValue>`],
      [V, 2670, `<pyLabel>${LAMPIRAN_REAS.viewOffice}</pyLabel>`],
      [V, 3434, `<pyLabel>${LAMPIRAN_REAS.delete}</pyLabel>`],
      [F, 22, `<pyCancelLabel>${LAMPIRAN_REAS.cancel}</pyCancelLabel>`],
      [F, 24, `<pySubmitLabel>${LAMPIRAN_REAS.attach}</pySubmitLabel>`],
    ] as const) {
      expect(baris(berkas, nomor), `${berkas} b${nomor}`).toBe(mau)
    }
  })

  it('syarat View Office Online (b2951) dan alamat penampil (DownloadDocumentPolis b540)', () => {
    expect(baris('Section\\ReasViewAttachment.xml', 2951)).toBe(
      "<pyCondition>(.MIME='xls' ||.MIME='xlsx' ||.MIME='doc' ||.MIME='docx'||.MIME='ppt' ||.MIME='pptx') &amp;&amp; .T_STORAGE_ID != ''</pyCondition>",
    )
    expect(baris('Activity\\DownloadDocumentPolis.xml', 540)).toBe(
      `<PropertiesValue>"${PENAMPIL_OFFICE}?${PARAM_PENAMPIL}="+@encodeURL(LinkDocument.url)</PropertiesValue>`,
    )
    expect(panel()).toContain('{d.adaObjek && bisaViewOffice(d.ekstensi) && (')
  })
})

describe('View Office Online dan View', () => {
  it('Office hanya xls/xlsx/doc/docx/ppt/pptx; View hanya pdf dan gambar raster', () => {
    expect(['xls', 'XLSX', 'doc', 'docx', 'ppt', 'pptx'].every(bisaViewOffice)).toBe(true)
    expect(['pdf', 'csv', 'png', ''].some(bisaViewOffice)).toBe(false)
    expect(['pdf', 'PNG', 'jpeg', 'svg', 'html', 'xlsx'].map(jenisViewOnline)).toEqual([
      'pdf',
      'gambar',
      'gambar',
      null,
      null,
      null,
    ])
    expect([mimeViewOnline('pdf'), mimeViewOnline('JPG'), mimeViewOnline('svg')]).toEqual([
      'application/pdf',
      'image/jpeg',
      '',
    ])
  })

  it('alamat literal HANYA di penampil.ts; dibuka lewat form GET ke bingkai, bukan jendela baru', () => {
    const berkas = (dir: string): string[] =>
      readdirSync(dir).flatMap((n) => {
        const p = join(dir, n)
        if (statSync(p).isDirectory()) return berkas(p)
        return /\.tsx?$/.test(n) && !n.includes('.test.') ? [p] : []
      })
    const beralamat = berkas(__dirname)
      .filter((p) => readFileSync(p, 'utf8').includes('://'))
      .map((p) => p.slice(__dirname.length + 1).replace(/\\/g, '/'))
    expect(beralamat).toEqual(['penampil.ts'])
    const isi = panel()
    expect(isi).toContain(
      '<form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>',
    )
    expect(isi).toContain('name={BINGKAI_PENAMPIL}')
  })

  it('popup penampil menggantikan daftar, ber-key berbeda (bug layar beku Bordereaux 08-10-2026)', () => {
    const isi = panel()
    expect(isi).toContain('if (penampil !== null) {')
    expect(isi).toMatch(
      /<Modal key="penampil" judul=\{penampil\.nama\} onTutup=\{tutupPenampil\} labelBatal=\{LAMPIRAN_REAS\.close\} penuh>/,
    )
    expect(isi).toMatch(/<Modal\s+key="daftar"/)
    expect(isi).toContain('URL.revokeObjectURL(objekAktif.current)')
  })
})

describe('Upload / Delete hanya selama kasus belum Resolve (keputusan work owner 08-10-2026)', () => {
  it('tombol Upload File dan Delete tergantung bolehUbah dari backend; View File selalu', () => {
    const isi = panel()
    expect(isi).toContain('const bolehUbah = grid?.bolehUbah ?? false')
    expect(isi.match(/\{bolehUbah && \(/g)).toHaveLength(2)
    expect(isi).toContain('onClick={() => onUnggah(k)}')
    expect(isi).toContain('await hapusDokumen(dasar, d)')
  })

  it('seret-lepas: berkas pilihan dan jatuhan digabung tanpa ganda', () => {
    const a = new File(['aa'], 'Slip.pdf')
    const b = new File(['bbb'], 'UW.xlsx')
    expect(gabungBerkas([a], [b, new File(['aa'], 'SLIP.PDF')]).map((f) => f.name)).toEqual(['Slip.pdf', 'UW.xlsx'])
    const isi = panel()
    expect(isi).toContain('tambah(Array.from(e.dataTransfer.files))')
    expect(isi).toMatch(/<input\s+type="file"\s+multiple/)
  })
})

describe('tata letak contoh work owner 08-10-2026', () => {
  const css = () => readFileSync(join(__dirname, '..', 'styles.css'), 'utf8')
  // 20 kategori; hanya yang terakhir berisi (UW ANALYSIS 3 di contoh).
  const daftar: KategoriReas[] = Array.from({ length: 20 }, (_, i) => ({
    nama: `K${String(i + 1).padStart(2, '0')}`,
    cacah: i === 19 ? 3 : 0,
  }))
  const render = (halaman: number, bolehUbah: boolean, isi = daftar) =>
    renderToStaticMarkup(
      createElement(
        BahasaUI.Provider,
        { value: 'en' },
        createElement(GridKategori, {
          daftar: isi,
          bolehUbah,
          halaman,
          onHalaman: () => {},
          onUnggah: () => {},
          onLihat: () => {},
        }),
      ),
    )
  const hitung = (html: string, kelas: string) => html.split(kelas).length - 1

  it('berjarak dari kaki Save / Submit; judul tanpa Refresh (tidak ada di XML)', () => {
    expect(panel()).toContain('<section className="panel lampiran-reas">')
    expect(css()).toMatch(/\n {2}margin-top: 24px;\n\}/)
    expect(panel()).toContain('<h3 className="lampiran-reas__judul">{LAMPIRAN_REAS.judul}</h3>')
    expect(Object.keys(LAMPIRAN_REAS)).not.toContain('refresh')
  })

  it('5 baris per halaman, pager di atas tabel; halaman 4 = baris 16-20', () => {
    const html = render(4, true)
    const en = TEKS_UI.en
    expect(html).toContain(en.menampilkan(16, 20, 20))
    expect(html).toContain(en.halamanDari(4, 4))
    expect(html).toContain(en.sebelumnya)
    expect(html).toContain(en.berikutnya)
    expect(html.indexOf('class="pager"')).toBeLessThan(html.indexOf('<table'))
    expect(['K15', 'K16', 'K20'].map((k) => html.includes(`<td>${k}</td>`))).toEqual([false, true, true])
    expect(hitung(html, '<tr>')).toBe(6) // kepala + 5 baris
    // Halaman di luar jangkauan dijepit ke halaman terakhir.
    expect(render(9, true)).toContain(en.menampilkan(16, 20, 20))
  })

  it('pager berbahasa Inggris: panel dibungkus BahasaUI en (permintaan work owner 08-10-2026)', () => {
    expect(panel()).toContain('<BahasaUI.Provider value="en">')
    expect(render(1, true)).toContain(TEKS_UI.en.menampilkan(1, 5, 20))
  })

  it('Count = lencana, hijau bila berisi; Upload / View File = tombol ikon', () => {
    const html = render(4, true)
    expect(html).toContain('<span class="lampiran-reas__cacah lampiran-reas__cacah--ada">3</span>')
    expect(hitung(html, '<span class="lampiran-reas__cacah">0</span>')).toBe(4)
    expect(hitung(html, 'lampiran-reas__ikon--unggah')).toBe(5)
    expect(hitung(html, 'lampiran-reas__ikon--lihat')).toBe(5)
    expect(html).toContain(`title="${LAMPIRAN_REAS.uploadFile}" aria-label="${LAMPIRAN_REAS.uploadFile} K20"`)
    expect(html).toContain(`title="${LAMPIRAN_REAS.viewFile}" aria-label="${LAMPIRAN_REAS.viewFile} K20"`)
    expect(css()).toMatch(/\.lampiran-reas__grid thead th \{[^}]*text-transform: uppercase;/)
  })

  it('kasus Resolve: tanpa tombol Upload File, View File tetap; daftar kosong = No items', () => {
    const html = render(1, false)
    expect(hitung(html, 'lampiran-reas__ikon--unggah')).toBe(0)
    expect(hitung(html, 'lampiran-reas__ikon--lihat')).toBe(5)
    expect(render(1, true, [])).toBe(`<p class="muted">${LAMPIRAN_REAS.kosong}</p>`)
  })
})

describe('panel di bawah layar kasus NB dan EDM Treaty In', () => {
  it.each([
    ['nbtreatyin', 'PREFIX_NBTREATYIN', 'nbti__kaki'],
    ['edmtreatyin', 'PREFIX_EDMTREATYIN', 'edmt__kaki'],
  ])('%s', (modul, prefix, kaki) => {
    const layar = readFileSync(join(APLIKASI, 'modul', modul, 'frontend', 'pages', 'LayarKasus.tsx'), 'utf8')
    const sisip = `<PanelLampiranReas dasar={\`\${${prefix}}/kasus/\${encodeURIComponent(id)}/lampiran\`} />`
    expect(layar).toContain(sisip)
    // Sesudah kaki Save / Submit.
    expect(layar.indexOf(sisip)).toBeGreaterThan(layar.indexOf(kaki))
  })
})

describe('klien lampiran Reas', () => {
  let tertangkap: { url: string; init: RequestInit }[] = []
  beforeEach(() => {
    tertangkap = []
    vi.stubEnv('VITE_AUTH_STUB', 'true')
    vi.stubEnv('VITE_STUB_PELAKU', 'UJI-MAKER')
    vi.stubEnv('VITE_STUB_PERAN', PERAN.admin)
    vi.stubGlobal('fetch', (url: string, init: RequestInit) => {
      tertangkap.push({ url, init })
      return Promise.resolve(
        new Response('{"daftar":[],"bolehUbah":true,"url":"U","ok":true}', {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
    })
  })
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.unstubAllEnvs()
  })

  it('jalur dan metode setiap aksi; kategori bergaris miring lewat kueri', async () => {
    const dasar = '/api/nb-treaty-in/kasus/NB-1/lampiran'
    const d: DokumenReas = {
      id: '20261008013000123',
      namaFile: 'a.xlsx',
      ekstensi: 'xlsx',
      kategori: 'R/I SLIP',
      tanggal: '08-10-2026 13:30',
      pengunggah: 'UJI',
      adaObjek: true,
    }
    await ambilGrid(dasar)
    await ambilDokumen(dasar, 'R/I SLIP')
    await unggahDokumen(dasar, 'R/I SLIP', new File(['isi'], 'a.xlsx'))
    await tautanOffice(dasar, d)
    await hapusDokumen(dasar, d)
    await ambilIsiDokumen(dasar, d)
    expect(tertangkap.map((t) => `${t.init.method ?? 'GET'} ${t.url}`)).toEqual([
      `GET ${dasar}`,
      `GET ${dasar}/dokumen?kategori=R%2FI+SLIP`,
      `POST ${dasar}/dokumen?kategori=R%2FI%20SLIP`,
      `GET ${dasar}/dokumen/20261008013000123/office`,
      `POST ${dasar}/dokumen/20261008013000123/hapus`,
      `GET ${dasar}/dokumen/20261008013000123/isi`,
    ])
    expect((tertangkap[2]!.init.body as FormData).get('berkas')).toBeInstanceOf(File)
    for (const t of tertangkap) expect((t.init.headers as Record<string, string>)['X-Pelaku']).toBe('UJI-MAKER')
  })
})
