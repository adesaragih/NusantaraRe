// Penjaga ronde "layar Treaty In disamakan dengan Pega" (5 Oktober 2026).
//
// ⛔ Yang dijaga di sini KEADAAN KELIMA: tab yang syaratnya tidak terpenuhi
// TIDAK DIRENDER SAMA SEKALI. Itu berbeda dari `Kosong` (tabnya ada, isinya
// nol) dan dari `.trin__belum` (tabnya ada, kodenya belum ditulis), dan
// perbedaan itu hanya bertahan kalau ada yang menjaganya.
//
// ⚠️ Mengikuti `layar.test.ts`: membaca BERKAS SUMBER untuk hal yang berupa
// teks, dan memanggil fungsinya untuk hal yang berupa aturan.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import {
  SYARAT_TAB_NON_PROPORSIONAL,
  SYARAT_TAB_PROPORSIONAL,
  type SyaratTabKontrak,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  TOTAL_RETENSI,
} from './labels'
import { NON_PROPORSIONAL, PROPORSIONAL, tabUntuk } from './pages/FormKontrakTreatyIn'

const AKAR = join(__dirname)
// ⭐ Layar ini DIPECAH 5 Oktober 2026: halaman + `components/`. `FORM`
// karena itu membaca KEDUANYA, supaya tiap pernyataan di bawah tetap
// menanyakan hal yang sama — "apakah kode layar modul ini memuat ini" —
// tanpa satu pun disunting. Itu bukti pemindahannya murni.
const SELA = String.fromCharCode(10)
const FORM =
  readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8') +
  readdirSync(join(AKAR, 'components'))
    .filter((f) => f.endsWith('.tsx') || f.endsWith('.ts'))
    .map((f) => readFileSync(join(AKAR, 'components', f), 'utf8'))
    .join(SELA)
const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')
const LABELS = readFileSync(join(AKAR, 'labels.ts'), 'utf8')

/** Dokumen yang nol kuncinya — keadaan 318 kontrak di POOLDATA. */
const KOSONG: SyaratTabKontrak = {
  retroBerganda: '',
  edmState: '',
  edmJenisMaterial: '',
}

describe('§0 syarat tampil tab', () => {
  it('kontrak BARU memperlihatkan seluruh tab, kedua cabang', () => {
    // ⛔ Nol dokumen berarti pertanyaannya belum dapat dijawab, dan tab yang
    // disembunyikan karena itu tidak akan pernah diisi.
    expect(tabUntuk(PROPORSIONAL)).toEqual(TAB_PROPORSIONAL)
    expect(tabUntuk(NON_PROPORSIONAL)).toEqual(TAB_NON_PROPORSIONAL)
  })

  it('Retro HILANG dari cabang prop saat IsMultipleRetro bukan "true"', () => {
    // Keadaan 1.531 kontrak "false" + 318 tanpa kunci.
    expect(tabUntuk(PROPORSIONAL, KOSONG)).not.toContain('Retro')
    expect(tabUntuk(PROPORSIONAL, { ...KOSONG, retroBerganda: 'false' })).not.toContain('Retro')
  })

  it('Retro MUNCUL pada kelima kontrak yang IsMultipleRetro-nya "true"', () => {
    const tab = tabUntuk(PROPORSIONAL, { ...KOSONG, retroBerganda: 'true' })
    expect(tab).toContain('Retro')
    // ⭐ Dan HANYA Retro yang berubah — sisanya utuh.
    expect(tab).toEqual(TAB_PROPORSIONAL)
  })

  it('Value Difference HILANG dari cabang non-prop pada SELURUH 1.854 kontrak', () => {
    // ⚠️ Terukur 5 Oktober 2026: `EDMMaterialType` nol di SELURUH dokumen,
    // jadi syarat `== 1` tidak pernah terpenuhi lewat JSONDATA.
    expect(tabUntuk(NON_PROPORSIONAL, KOSONG)).not.toContain('Value Difference')
  })

  it('Value Difference MUNCUL hanya bila KEDUA syaratnya terpenuhi', () => {
    const ya: SyaratTabKontrak = { ...KOSONG, edmState: '1', edmJenisMaterial: '1' }
    expect(tabUntuk(NON_PROPORSIONAL, ya)).toContain('Value Difference')

    // State 3 menutupnya walau jenis materialnya benar.
    expect(
      tabUntuk(NON_PROPORSIONAL, { ...KOSONG, edmState: '3', edmJenisMaterial: '1' }),
    ).not.toContain('Value Difference')
    // Jenis material 2 (Non Material) menutupnya walau statenya benar.
    expect(
      tabUntuk(NON_PROPORSIONAL, { ...KOSONG, edmState: '1', edmJenisMaterial: '2' }),
    ).not.toContain('Value Difference')
  })

  it('EDMState kosong dihitung BUKAN 3 — != Pega bernilai benar', () => {
    const uji = SYARAT_TAB_NON_PROPORSIONAL['Value Difference']
    expect(uji).toBeDefined()
    expect(uji?.({ ...KOSONG, edmJenisMaterial: '1' })).toBe(true)
  })

  it('tab bersyarat HILANG dari daftar, bukan dirender kosong', () => {
    // ⛔ Pembedanya: panjang daftarnya BERUBAH. Tab yang dirender kosong
    // tetap sebelas butir; tab yang tidak dirender menjadi sepuluh.
    expect(tabUntuk(PROPORSIONAL, KOSONG)).toHaveLength(TAB_PROPORSIONAL.length - 1)
    expect(tabUntuk(NON_PROPORSIONAL, KOSONG)).toHaveLength(TAB_NON_PROPORSIONAL.length - 1)
  })

  it('Retro TETAP ADA di cabang non-prop — syaratnya milik cabang prop saja', () => {
    // ⛔ Jebakan yang sungguh terjadi saat ronde ini ditulis: peta syarat
    // berkunci NAMA TAB saja menyembunyikan `Retro` non-proporsional pada
    // 770 dari 775 kontrak. Ekspor non-prop (TABBED @3291822) NOL syarat.
    expect(tabUntuk(NON_PROPORSIONAL, KOSONG)).toContain('Retro')
    expect(Object.keys(SYARAT_TAB_NON_PROPORSIONAL)).not.toContain('Retro')
  })

  it('NOL tab dihapus dari kedua daftar — §0 melarangnya', () => {
    expect(TAB_PROPORSIONAL).toHaveLength(11)
    // ⚠️ DUA BELAS, dan butir ke-12 (RNM Share) adalah pertanyaan terbuka —
    // lihat komentarnya di `labels.ts`. Ia TIDAK dihapus.
    expect(TAB_NON_PROPORSIONAL).toHaveLength(12)
    expect(TAB_NON_PROPORSIONAL).toContain('RNM Share')
    expect(LABELS).toContain('PERTANYAAN TERBUKA, bukan tab menurut ekspor')
  })

  it('hanya DUA tab yang bersyarat, dan keduanya disebut sumbernya', () => {
    expect(Object.keys(SYARAT_TAB_PROPORSIONAL)).toEqual(['Retro'])
    expect(Object.keys(SYARAT_TAB_NON_PROPORSIONAL)).toEqual(['Value Difference'])
    expect(LABELS).toContain('TreatyIn.IsMultipleRetro')
    expect(LABELS).toContain('TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1')
  })
})

describe('§1 Accounting Mode adalah DUA properti', () => {
  it('nilai kedua himpunan berbeda dan tidak bercampur', () => {
    expect(LABELS).toContain("caraPembukuanNonPropNilai: ['loss', 'risk']")
    expect(LABELS).toContain("caraPembukuanNilai: ['underwriting', 'accounting']")
  })

  it('layar memilih menurut cabang, bukan memakai satu medan', () => {
    expect(FORM).toContain('pembukuanNonProp')
    // ⚠️ RALAT 6 Oktober 2026: daftar nilainya tidak lagi tertanam di layar —
    // ia datang dari sumber lewat `opsi`. Yang dijaga tetap yang pokok:
    // cabang non-prop membaca PROPERTI KEDUA, bukan properti cabang seberang.
    expect(FORM).toContain('opsi?.caraPembukuanNonProp')
    // Keadaannya DUA, supaya berpindah cabang tidak menimpa yang ditinggal.
    expect(FORM).toContain("const [pembukuanNonProp, setPembukuanNonProp] = useState('')")
  })
})

describe('§2 Bordereaux', () => {
  it('ejaannya Bordeaux di dokumen, Bordereaux di layar', () => {
    expect(LABELS).toContain("bordereaux: 'Bordereaux'")
    expect(LABELS).toContain('TreatyIn.Bordeaux')
  })

  it('dropdownnya PROPORSIONAL SAJA', () => {
    expect(FORM).toContain('{jenis !== NON_PROPORSIONAL && (')
    expect(FORM).toContain('@60649')
  })
})

describe('§4 medan tanggal kosong', () => {
  it('kaitnya di modul, dan dasar.tsx tidak disentuh', () => {
    expect(FORM).toContain('trin__tanggal--kosong')
    expect(FORM).toContain('function TanggalRedup(')
  })

  it('aturannya hanya berlaku saat KOSONG dan TIDAK difokus', () => {
    expect(CSS).toContain('.trin__tanggal--kosong')
    expect(CSS).toContain(':not(:focus)::-webkit-datetime-edit')
    // ⭐ Yang transparan TEKSNYA saja — medannya tetap dapat diklik dan
    // ikon pemilih tanggalnya tetap terlihat.
    expect(CSS).toContain('color: transparent;')
  })

  it('batasnya DINYATAKAN, tidak disembunyikan', () => {
    expect(CSS).toContain('Firefox')
  })
})

describe('§5 Update Total dan Total Retention Amount', () => {
  it('panelnya dua kolom, dan judul kolom pertama adalah nama panelnya', () => {
    expect(TOTAL_RETENSI.judul).toBe('Total Retention Amount')
    expect(TOTAL_RETENSI.kolomNilai).toBe('Value')
    expect(TOTAL_RETENSI.tanpaBaris).toBe('No items')
  })

  it('tombolnya ADA dan MATI — bukan dihilangkan', () => {
    expect(TOTAL_RETENSI.perbarui).toBe('Update Total')
    expect(FORM).toContain('<PanelTotalRetensi baris={warisan?.totalRetensi ?? []} />')
    expect(FORM).toContain('{TOTAL_RETENSI.perbarui}')
    const i = FORM.indexOf('{TOTAL_RETENSI.perbarui}')
    expect(FORM.slice(i - 200, i)).toContain('disabled')
  })

  it('keempat tab lain yang juga punya Update Total TERCATAT', () => {
    // ⚠️ Yang dibangun ronde ini satu; yang terukur lima.
    expect([...TOTAL_RETENSI.tabLain].sort()).toEqual(['EGNPI', 'Installment', 'Limits', 'Share'])
  })

  it('rumusnya tidak dihitung ulang di layar', () => {
    // ⛔ Layar membaca `totalRetensi`; nol penjumlahan di berkas layar.
    expect(FORM).toContain('warisan?.totalRetensi')
    expect(FORM).not.toContain('reduce((t, b) => t + Number(b.jumlah)')
  })
})
