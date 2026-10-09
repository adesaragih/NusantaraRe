// Panel baris grid masterDetail UMUM Claim Fac In (objek -> item -> estimasi / adjustment): kunci panel `prefiks:jalur(n)`
// (`models.KunciPanel`), keputusan baris dapat dibuka, render rekursif panel di dalam panel, dan alamat aksi
// (`konteks` + `indeks`) unsur di dalam panel bersarang dan modal - `services.PermintaanAksi`.

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it, vi } from 'vitest'

import type { Halaman, Tata } from '../api'
import { barisTerbuka, kunciPanel, panelBaris } from './rincian'
import TataView, {
  gantiMedan,
  klikTombol,
  konteksDi,
  kunciGrid,
  PanelBaris,
  Tombol,
  type KonteksTata,
} from './TataView'

const OBJEK = 'ClaimData.ObjectList'
const ITEM = `${OBJEK}(1).ObjectItemList`
const PANEL_OBJEK = `est:${OBJEK}(1)`
const PANEL_ITEM = `estitem:${ITEM}(2)`

function konteks(ubah: Partial<KonteksTata> = {}): KonteksTata {
  return {
    h: { nilai: {}, daftar: {} },
    ubah: () => {},
    aksi: () => {},
    opsi: () => [],
    pesanMedan: {},
    sibuk: false,
    konteks: '',
    ...ubah,
  }
}

const tampil = (n: number) => Array.from({ length: n }, () => [{ tampil: true, hanyaBaca: true }])

describe('kunci panel baris (models.KunciPanel)', () => {
  it('prefiks:jalur(n) - juga jalur grid bersarang', () => {
    expect(kunciPanel('est', OBJEK, 1)).toBe(PANEL_OBJEK)
    expect(kunciPanel('estitem', ITEM, 2)).toBe(PANEL_ITEM)
    expect(kunciPanel('adjdtl', `${ITEM}(2).Adjustment`, 3)).toBe(
      'adjdtl:ClaimData.ObjectList(1).ObjectItemList(2).Adjustment(3)',
    )
  })

  it('baris dapat dibuka hanya bila grid ber-rincian DAN server mengirim panel baris itu', () => {
    const grid = { rincian: 'est', jalur: OBJEK }
    const panel = { [PANEL_OBJEK]: [] }
    expect(panelBaris(grid, 1, panel)).toBe(PANEL_OBJEK)
    expect(panelBaris(grid, 2, panel)).toBeNull()
    expect(panelBaris({ jalur: OBJEK }, 1, panel)).toBeNull()
    expect(panelBaris(grid, 0, panel)).toBeNull()
    expect(panelBaris(grid, 1, undefined)).toBeNull()
  })

  it('Add membuka baris baru; hapus menutup baris yang hilang; lainnya tetap', () => {
    expect(barisTerbuka(null, 1, 2)).toBe(2)
    expect(barisTerbuka(1, 1, 2)).toBe(2)
    expect(barisTerbuka(2, 2, 1)).toBeNull()
    expect(barisTerbuka(1, 3, 2)).toBe(1)
    expect(barisTerbuka(null, 2, 2)).toBeNull()
  })
})

describe('key grid', () => {
  // Grid objek Estimation (`est`) dan Adjustment (`adj`) berjalur sama di posisi yang sama: key berbeda supaya baris
  // terbuka satu tidak terwaris ke yang lain.
  it('posisi, jalur, dan prefiks panel ikut key', () => {
    expect(kunciGrid(0, { jalur: OBJEK, rincian: 'est' })).not.toBe(kunciGrid(0, { jalur: OBJEK, rincian: 'adj' }))
    expect(kunciGrid(0, { jalur: OBJEK })).not.toBe(kunciGrid(1, { jalur: OBJEK }))
  })
})

describe('render panel rekursif', () => {
  const gridObjek: Tata = {
    jenis: 'grid',
    jalur: OBJEK,
    bernomor: true,
    rincian: 'est',
    kolom: [{ jenis: 'medan', jalur: 'ObjectName', label: 'Object Type', kendali: 'tampil' }],
    baris: tampil(2),
  }
  const h: Halaman = {
    nilai: { [`${ITEM}(2).TotalGrossEstimasi`]: '' },
    daftar: { [OBJEK]: [{ ObjectName: 'UJI-OBJEK-1' }, { ObjectName: 'UJI-OBJEK-2' }], [ITEM]: [{}, {}] },
  }

  it('grid ber-rincian: baris ber-panel dapat dibuka (tertutup sampai diklik), baris tanpa panel tidak', () => {
    const html = renderToStaticMarkup(
      createElement(TataView, { tata: [gridObjek], k: konteks({ h, panel: { [PANEL_OBJEK]: [] } }) }),
    )
    expect(html.match(/aria-expanded="false"/g)).toHaveLength(1)
    expect(html).not.toContain('aria-expanded="true"')
    expect(html).toContain('UJI-OBJEK-2')
    expect(html).toContain('▸')
  })

  it('grid tanpa nomor baris ber-panel mendapat kolom panah di depan', () => {
    const html = renderToStaticMarkup(
      createElement(TataView, {
        tata: [{ ...gridObjek, bernomor: false }],
        k: konteks({ h, panel: { [PANEL_OBJEK]: [] } }),
      }),
    )
    expect(html).toContain('<th></th><th>Object Type</th>')
  })

  it('isi panel objek dirender dengan grid item ber-rincian sendiri (panel di dalam panel)', () => {
    const gridItem: Tata = {
      jenis: 'grid',
      jalur: ITEM,
      rincian: 'estitem',
      kolom: [{ jenis: 'medan', jalur: 'ObjectItemName', label: 'Object Name', kendali: 'tampil' }],
      baris: tampil(2),
    }
    const panel: Record<string, Tata[]> = {
      [PANEL_OBJEK]: [{ jenis: 'label', label: 'UJI-ISI-PANEL-OBJEK' }, gridItem],
      // hanya item ke-2 yang punya panel
      [PANEL_ITEM]: [{ jenis: 'label', label: 'UJI-ISI-PANEL-ITEM' }],
    }
    const html = renderToStaticMarkup(createElement(PanelBaris, { kunci: PANEL_OBJEK, k: konteks({ h, panel }) }))
    expect(html).toContain('class="claimfacin__rinci"')
    expect(html).toContain('UJI-ISI-PANEL-OBJEK')
    expect(html.match(/aria-expanded="false"/g)).toHaveLength(1)
    // panel item tertutup sampai barisnya diklik
    expect(html).not.toContain('UJI-ISI-PANEL-ITEM')
    // panel item itu sendiri dapat dirender dari kuncinya
    expect(renderToStaticMarkup(createElement(PanelBaris, { kunci: PANEL_ITEM, k: konteks({ h, panel }) }))).toContain(
      'UJI-ISI-PANEL-ITEM',
    )
  })

  it('panel yang tidak dikirim server dirender kosong, bukan gagal', () => {
    const html = renderToStaticMarkup(createElement(PanelBaris, { kunci: 'est:UJI(9)', k: konteks() }))
    expect(html).toBe('<div class="claimfacin__rinci"><div class="claimfacin__tumpuk"></div></div>')
  })
})

describe('alamat aksi: konteks + indeks', () => {
  it('tombol sel grid di panel item (di dalam panel objek): konteks = kunci panel item, indeks = baris grid estimasi', () => {
    const aksi = vi.fn()
    const k = konteksDi(konteksDi(konteks({ aksi }), PANEL_OBJEK), PANEL_ITEM)
    klikTombol({ aksi: 'DeleteValueEstimation' }, k, 3)
    expect(aksi).toHaveBeenCalledWith('DeleteValueEstimation', 3, {}, PANEL_ITEM)
  })

  it('tombol Add kepala grid: indeks 0 di konteks panelnya', () => {
    const aksi = vi.fn()
    klikTombol({ aksi: 'TambahItem' }, konteksDi(konteks({ aksi }), PANEL_OBJEK), 0)
    expect(aksi).toHaveBeenCalledWith('TambahItem', 0, {}, PANEL_OBJEK)
  })

  it('tombol layar utama (baris grid objek): konteks kosong, indeks = baris objek', () => {
    const aksi = vi.fn()
    klikTombol({ aksi: 'BukaOutstanding' }, konteks({ aksi }), 2)
    expect(aksi).toHaveBeenCalledWith('BukaOutstanding', 2, {}, '')
  })

  it('sel autocomplete di panel objek: nilai + param pilihan, indeks = baris item', () => {
    const aksi = vi.fn()
    const k = konteksDi(konteks({ aksi }), PANEL_OBJEK)
    gantiMedan({ aksi: 'PilihOkupasi' }, k, `${ITEM}(2).OccupationName`, 2, 'UJI-OKP', true, 'UJI-OKP')
    expect(aksi).toHaveBeenCalledWith(
      'PilihOkupasi',
      2,
      { [`${ITEM}(2).OccupationName`]: 'UJI-OKP' },
      PANEL_OBJEK,
      'UJI-OKP',
    )
  })

  it('medan tanpa aksi (postValue) hanya mengubah nilai lokal', () => {
    const aksi = vi.fn()
    const ubah = vi.fn()
    gantiMedan({}, konteksDi(konteks({ aksi, ubah }), PANEL_ITEM), `${ITEM}(2).DeductibleType`, 0, 'true', true)
    expect(aksi).not.toHaveBeenCalled()
    expect(ubah).toHaveBeenCalledWith(`${ITEM}(2).DeductibleType`, 'true')
  })

  it('medan beraksi di panel di luar grid: indeks 0', () => {
    const aksi = vi.fn()
    gantiMedan({ aksi: 'CountTSI' }, konteksDi(konteks({ aksi }), PANEL_ITEM), `${ITEM}(2).Amount`, 0, '5', true)
    expect(aksi).toHaveBeenCalledWith('CountTSI', 0, { [`${ITEM}(2).Amount`]: '5' }, PANEL_ITEM, '')
  })
})

describe('kunci modal = konteks aksi di dalamnya', () => {
  it('Submit Print PLA beralamat modal pla:<objek>', () => {
    const aksi = vi.fn()
    klikTombol({ aksi: 'GeneratePLA' }, konteksDi(konteks({ aksi }), 'pla:1'), 0)
    expect(aksi).toHaveBeenCalledWith('GeneratePLA', 0, {}, 'pla:1')
  })

  it('Choose pop-up Cedant: aksi berparameter dikirim apa adanya, indeks = baris CedingCedantList', () => {
    const aksi = vi.fn()
    const kunci = 'cedant:ClaimData.ObjectList(1).ObjectItemList(1).Adjustment(2)'
    klikTombol({ aksi: 'SetCedant:1' }, konteksDi(konteks({ aksi }), kunci), 3)
    expect(aksi).toHaveBeenCalledWith('SetCedant:1', 3, {}, kunci)
  })

  it('tombol tanpa aksi (Cancel) menutup modal / panelnya tanpa permintaan server', () => {
    const aksi = vi.fn()
    const tutup = vi.fn()
    klikTombol({}, konteksDi(konteks({ aksi }), 'pilihPolis', tutup), 0)
    expect(tutup).toHaveBeenCalledTimes(1)
    expect(aksi).not.toHaveBeenCalled()
  })
})

describe('tombol nonaktif-OQ dari server', () => {
  it('ikon "+" Consultant / Adjuster: nonaktif, keterangan OQ di title (tombolOQ server)', () => {
    const oq = 'OQ-CFI-32: UJI keterangan'
    const html = renderToStaticMarkup(
      createElement(Tombol, {
        t: { jenis: 'tombol', id: 'TambahAdjuster', ikon: 'tambah', nonaktif: true, catatan: oq },
        k: konteks(),
      }),
    )
    expect(html).toContain('disabled=""')
    expect(html).toContain(`title="${oq}"`)
  })
})
