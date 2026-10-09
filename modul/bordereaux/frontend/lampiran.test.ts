// Lampiran Bordereaux - `AttachmentsBdx` / `AttachmentDetailBdx` / `BordereauxAttach` (keputusan work owner
// 08-10-2026): label berbukti baris korpus, rute klien, penampil Office, dan syarat Upload File / Delete.

import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { PERAN } from '../../../inti/frontend/labels'
import {
  ambilIsiLampiran,
  ambilKategoriLampiran,
  ambilLampiran,
  hapusLampiran,
  PREFIX_BDX,
  tautanOffice,
  unggahLampiran,
  type Lampiran,
} from './api'
import { gabungBerkas } from './aturan'
import { LAMPIRAN_BDX } from './labels'
import { bisaViewOffice, jenisViewOnline, mimeViewOnline, PARAM_PENAMPIL, PENAMPIL_OFFICE } from './penampilOffice'

const KORPUS = 'D:\\XML\\RNM_BRD\\Bordereaux'
const adaKorpus = existsSync(KORPUS)

/** Satu baris (1-based) berkas korpus, tanpa indentasi. */
function baris(relatif: string, nomor: number): string {
  return (readFileSync(join(KORPUS, relatif), 'utf8').split('\n')[nomor - 1] ?? '').trim()
}

describe.skipIf(!adaKorpus)('label lampiran berbukti barisnya', () => {
  it('AttachmentsBdx, AttachmentDetailBdx, BordereauxAttach', () => {
    const A = 'Section\\AttachmentsBdx.xml'
    const D = 'Section\\AttachmentDetailBdx.xml'
    for (const [berkas, nomor, mau] of [
      [A, 2422, `<pyTitle>${LAMPIRAN_BDX.judul}</pyTitle>`],
      [A, 1800, `<pyLabel>${LAMPIRAN_BDX.refresh}</pyLabel>`],
      [A, 3039, `<pyValue>${LAMPIRAN_BDX.category}</pyValue>`],
      [A, 3192, `<pyValue>${LAMPIRAN_BDX.count}</pyValue>`],
      [A, 3342, `<pyValue>${LAMPIRAN_BDX.uploadFile}</pyValue>`],
      [A, 3490, `<pyValue>${LAMPIRAN_BDX.viewFile}</pyValue>`],
      [D, 1529, `<pyValue>${LAMPIRAN_BDX.fileName}</pyValue>`],
      [D, 2562, `<pyLabel>${LAMPIRAN_BDX.viewOffice}</pyLabel>`],
      [D, 1796, `<pyValue>${LAMPIRAN_BDX.type}</pyValue>`],
      [D, 3185, `<pyLabel>${LAMPIRAN_BDX.delete}</pyLabel>`],
      ['FlowAction\\BordereauxAttach.xml', 21, `<pySubmitLabel>${LAMPIRAN_BDX.submit}</pySubmitLabel>`],
      ['FlowAction\\BordereauxAttach.xml', 20, `<pyCancelLabel>${LAMPIRAN_BDX.cancel}</pyCancelLabel>`],
    ] as const) {
      expect(baris(berkas, nomor), `${berkas} b${nomor}`).toBe(mau)
    }
  })

  it('Download All mati, syarat Upload/Delete dan View Office Online, panel di InputBordereaux', () => {
    expect(baris('Section\\AttachmentsBdx.xml', 1566)).toBe('<pyLabel>Download All</pyLabel>')
    expect(baris('Section\\AttachmentsBdx.xml', 1656)).toBe('<pyCondition>never</pyCondition>')
    expect(Object.values(LAMPIRAN_BDX)).not.toContain('Download All')
    const syarat = "<pyCondition>BORDEREAUX.ViewStage !=1 || OperatorID.pyPosition='IT Developer'</pyCondition>"
    expect(baris('Section\\AttachmentsBdx.xml', 4459)).toBe(syarat)
    expect(baris('Section\\AttachmentDetailBdx.xml', 3376)).toBe(syarat)
    expect(baris('Section\\AttachmentDetailBdx.xml', 2876)).toBe(
      "<pyCondition>.HASIL5='xls' ||.HASIL5='xlsx' ||.HASIL5='doc' ||.HASIL5='docx'||.HASIL5='ppt' ||.HASIL5='pptx'</pyCondition>",
    )
    expect(baris('DownloadAttachmentBdx.xml', 388)).toBe(
      `<PropertiesValue>"${PENAMPIL_OFFICE}?${PARAM_PENAMPIL}="+@encodeURL(LinkDocument.url)</PropertiesValue>`,
    )
    expect(baris('Section\\InputBordereaux.xml', 27972)).toBe('<pyInclude>AttachmentsBdx</pyInclude>')
  })
})

describe('View Office Online', () => {
  it('hanya ekstensi xls/xlsx/doc/docx/ppt/pptx (b2876)', () => {
    expect(['xls', 'XLSX', 'doc', 'docx', 'ppt', 'pptx'].every(bisaViewOffice)).toBe(true)
    expect(['pdf', 'csv', 'png', ''].some(bisaViewOffice)).toBe(false)
  })

  it('alamat literal HANYA di penampilOffice.ts; dibuka lewat form GET ke bingkai, bukan jendela baru', () => {
    const berkas = (dir: string): string[] =>
      readdirSync(dir).flatMap((n) => {
        const p = join(dir, n)
        if (statSync(p).isDirectory()) return berkas(p)
        return /\.tsx?$/.test(n) && !n.includes('.test.') ? [p] : []
      })
    const beralamat = berkas(__dirname)
      .filter((p) => readFileSync(p, 'utf8').includes('://'))
      .map((p) => p.slice(__dirname.length + 1).replace(/\\/g, '/'))
    expect(beralamat).toEqual(['penampilOffice.ts'])
    const panel = readFileSync(join(__dirname, 'components', 'PanelLampiran.tsx'), 'utf8')
    expect(panel).toContain(
      '<form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>',
    )
    expect(panel).toContain('name={BINGKAI_PENAMPIL}')
  })
})

describe('panel di form berkas', () => {
  it('berkas tersimpan saja; Upload/Delete = hak.lampiran dan BUKAN mode View, siapa pun (keputusan work owner 08-10-2026)', () => {
    const form = readFileSync(join(__dirname, 'components', 'FormBordereaux.tsx'), 'utf8')
    expect(form).toContain('{!baru && <PanelLampiran bdxId={isian.bdxId} boleh={hak.lampiran && !lihat} />}')
    // Sesudah kaki Close/Save, sebelum History (urutan layout InputBordereaux).
    expect(form.indexOf('<PanelLampiran')).toBeGreaterThan(form.indexOf('bordereaux__kaki'))
    expect(form.indexOf('<PanelLampiran')).toBeLessThan(form.indexOf('{BDX.history}'))
  })
})

describe('klien lampiran', () => {
  let tertangkap: { url: string; init: RequestInit }[] = []
  beforeEach(() => {
    tertangkap = []
    vi.stubEnv('VITE_AUTH_STUB', 'true')
    vi.stubEnv('VITE_STUB_PELAKU', 'UJI-MAKER')
    vi.stubEnv('VITE_STUB_PERAN', PERAN.admin)
    vi.stubGlobal('fetch', (url: string, init: RequestInit) => {
      tertangkap.push({ url, init })
      return Promise.resolve(
        new Response('{"daftar":[],"url":"U","ok":true}', {
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

  it('jalur dan metode setiap aksi', async () => {
    const l: Lampiran = {
      id: '2026',
      kategoriId: '00 01',
      kategori: 'SOA',
      fileName: 'a.xlsx',
      ekstensi: 'xlsx',
      username: 'UJI',
    }
    await ambilKategoriLampiran('BDX-1')
    await ambilLampiran('BDX-1', '00 01')
    await unggahLampiran('BDX-1', '00 01', new File(['isi'], 'a.xlsx'))
    await tautanOffice('BDX-1', l)
    await hapusLampiran('BDX-1', l)
    const dasar = `${PREFIX_BDX}/berkas/BDX-1/lampiran`
    expect(tertangkap.map((t) => `${t.init.method ?? 'GET'} ${t.url}`)).toEqual([
      `GET ${dasar}`,
      `GET ${dasar}/00%2001`,
      `POST ${dasar}/00%2001`,
      `GET ${dasar}/00%2001/2026/office`,
      `POST ${dasar}/00%2001/2026/hapus`,
    ])
    const unggah = tertangkap[2]!.init.body as FormData
    expect((unggah.get('berkas') as File).name).toBe('a.xlsx')
    for (const t of tertangkap) expect((t.init.headers as Record<string, string>)['X-Pelaku']).toBe('UJI-MAKER')
  })
})

describe('Upload File seret-lepas (dragDropFileUpload)', () => {
  it('berkas pilihan dan jatuhan digabung tanpa ganda; urutan lama dipertahankan', () => {
    const a = new File(['aa'], 'Nota.pdf')
    const b = new File(['bbb'], 'SOA.xlsx')
    const aLagi = new File(['aa'], 'NOTA.PDF')
    const aLain = new File(['beda isi'], 'Nota.pdf')
    expect(gabungBerkas([a], [b, aLagi, aLain]).map((f) => f.name)).toEqual(['Nota.pdf', 'SOA.xlsx', 'Nota.pdf'])
    expect(gabungBerkas([], [])).toEqual([])
  })

  it('kotak menangkap dragover/drop, menggabung dataTransfer.files, dan tetap bisa diklik', () => {
    const panel = readFileSync(join(__dirname, 'components', 'PanelLampiran.tsx'), 'utf8')
    expect(panel).toContain('onDragOver={(e) => {')
    expect(panel).toContain('tambah(Array.from(e.dataTransfer.files))')
    expect(panel).toContain('setBerkas((b) => gabungBerkas(b, baru))')
    expect(panel).toMatch(/<input\s+type="file"\s+multiple/)
    expect(LAMPIRAN_BDX.seretBerkas).toBe('Drag and drop files here, or click to choose files')
  })
})

describe('View langsung seperti Product Name Life (permintaan work owner 08-10-2026)', () => {
  it('View hanya pdf dan gambar raster; tipe objek URL dari ekstensi', () => {
    expect(['pdf', 'PDF', 'png', 'jpg', 'jpeg', 'gif', 'bmp', 'webp'].map(jenisViewOnline)).toEqual([
      'pdf',
      'pdf',
      'gambar',
      'gambar',
      'gambar',
      'gambar',
      'gambar',
      'gambar',
    ])
    expect(['svg', 'html', 'xlsx', 'csv', ''].map(jenisViewOnline)).toEqual([null, null, null, null, null])
    expect([mimeViewOnline('pdf'), mimeViewOnline('JPG'), mimeViewOnline('svg')]).toEqual([
      'application/pdf',
      'image/jpeg',
      '',
    ])
  })

  it('View dan View Office Online membuka popup penampil layar penuh yang menggantikan daftar', () => {
    const panel = readFileSync(join(__dirname, 'components', 'PanelLampiran.tsx'), 'utf8')
    expect(panel).toContain('if (penampil !== null) {')
    // Bug 08-10-2026: tanpa key, popup daftar mewarisi status menutup penampil dan layar tidak bisa diklik.
    expect(panel).toContain('<Modal key="penampil"')
    expect(panel).toMatch(/<Modal\s+key="daftar"/)
    expect(panel).toMatch(
      /<Modal key="penampil" judul=\{penampil\.nama\} onTutup=\{tutupPenampil\} labelBatal=\{LAMPIRAN_BDX\.close\} penuh>/,
    )
    expect(panel).toContain('const isi = await ambilIsiLampiran(bdxId, l)')
    expect(panel).toContain('URL.revokeObjectURL(objekAktif.current)')
    expect(LAMPIRAN_BDX.view).toBe('View')
  })

  it('isi lampiran: GET rute unduh beridentitas, dijawab Blob', async () => {
    const tertangkap: { url: string; init: RequestInit }[] = []
    vi.stubEnv('VITE_AUTH_STUB', 'true')
    vi.stubEnv('VITE_STUB_PELAKU', 'UJI-CHK')
    vi.stubEnv('VITE_STUB_PERAN', PERAN.admin)
    vi.stubGlobal('fetch', (url: string, init: RequestInit) => {
      tertangkap.push({ url, init })
      return Promise.resolve(new Response('ISI PDF', { status: 200, headers: { 'Content-Type': 'application/pdf' } }))
    })
    try {
      const l: Lampiran = {
        id: '2026',
        kategoriId: '00000',
        kategori: 'Others',
        fileName: 'a.pdf',
        ekstensi: 'pdf',
        username: 'UJI',
      }
      const b = await ambilIsiLampiran('BDX-1', l)
      expect(await b.text()).toBe('ISI PDF')
      expect(tertangkap[0]!.url).toBe(`${PREFIX_BDX}/berkas/BDX-1/lampiran/00000/2026/isi`)
      expect(tertangkap[0]!.init.method).toBe('GET')
      expect((tertangkap[0]!.init.headers as Record<string, string>)['X-Pelaku']).toBe('UJI-CHK')
    } finally {
      vi.unstubAllGlobals()
      vi.unstubAllEnvs()
    }
  })
})
