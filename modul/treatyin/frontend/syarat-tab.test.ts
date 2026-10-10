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
import { totalTerkunci } from './labelsRetensi'
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
const RETENSI_SRC = readFileSync(join(AKAR, 'components', 'TabRetensi.tsx'), 'utf8')
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

  // ⛔ DIBALIK 8 Oktober 2026 — tab itu DIKELUARKAN dari proyek atas
  // keputusan pemilik proses. Kedua uji lama menjaga kapan ia tampil; yang
  // ini menjaga bahwa ia TIDAK PERNAH tampil, pada nilai apa pun.
  it('⛔ tab itu nol di kedua cabang, pada keadaan dokumen apa pun', () => {
    for (const k of [KOSONG, { ...KOSONG, retroBerganda: 'false' }, { ...KOSONG, retroBerganda: 'true' }]) {
      expect(tabUntuk(PROPORSIONAL, k)).not.toContain('Retro')
      expect(tabUntuk(NON_PROPORSIONAL, k)).not.toContain('Retro')
    }
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
    // tetap seluruh butirnya; tab yang tidak dirender berkurang satu.
    //
    // ⚠️ Sejak 8 Oktober 2026 hanya cabang NON-PROPORSIONAL yang masih
    // punya tab bersyarat (`Value Difference`). Cabang proporsional nol —
    // satu-satunya syaratnya dikeluarkan dari proyek, jadi daftarnya utuh.
    expect(tabUntuk(PROPORSIONAL, KOSONG)).toHaveLength(TAB_PROPORSIONAL.length)
    expect(tabUntuk(NON_PROPORSIONAL, KOSONG)).toHaveLength(TAB_NON_PROPORSIONAL.length - 1)
  })

  // ⛔ RALAT 8 Oktober 2026 — satu-satunya tab yang dihapus: `RNM Share`
  // Non-Prop, atas keputusan pemilik proses (*"di non prop tab RNM SHARE itu
  // tidak ada"*). Pertanyaan terbukanya TERTUTUP; alasannya tercatat.
  it('hanya RNM Share Non-Prop yang dihapus, dan keputusannya tercatat', () => {
    expect(TAB_PROPORSIONAL).toHaveLength(10)
    expect(TAB_NON_PROPORSIONAL).toHaveLength(10)
    expect(TAB_NON_PROPORSIONAL as readonly string[]).not.toContain('RNM Share')
    expect(LABELS).toContain('di non prop tab RNM SHARE itu tidak')
  })

  it('SATU tab bersyarat tersisa, dan sumbernya disebut', () => {
    // ⛔ Peta cabang proporsional kini KOSONG: satu-satunya entrinya sudah
    // dikeluarkan dari proyek. Bentuknya tetap ada sebab `tabUntuk`
    // menerima peta untuk kedua cabang.
    expect(Object.keys(SYARAT_TAB_PROPORSIONAL)).toEqual([])
    expect(Object.keys(SYARAT_TAB_NON_PROPORSIONAL)).toEqual(['Value Difference'])
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

// ⭐ 7 Oktober 2026 — permintaan pemakai: "tiap inputan tanggal bisa diketik
// juga biar gampang". `<input type="date">` (yang menulis `dd/mm/yyyy`
// sendiri saat kosong, dan dulu ditambal `trin__tanggal--kosong`) diganti
// kotak TEKS `DD/MM/YYYY` + ikon kalender — bentuk Pega gambar 38. Contoh
// ketikan dijaga `tanggal-ketik.test.tsx`.
describe('§4 medan tanggal — dapat DIKETIK', () => {
  it('kaitnya di modul, dan dasar.tsx tidak disentuh', () => {
    expect(FORM).toContain('function TanggalRedup(')
    expect(FORM).toContain('<KotakTanggalKetik')
    expect(FORM).not.toContain('<FieldTanggal')
  })

  it('tambalan `dd/mm/yyyy` dicabut — kotak teks yang kosong memang kosong', () => {
    expect(CSS).not.toContain(':not(:focus)::-webkit-datetime-edit')
    expect(CSS).not.toContain('.treatyin .trin__tanggal--kosong')
  })

  it('kalender bawaan tetap ada — tak terlihat, bukan `display: none`', () => {
    expect(CSS).toContain('.treatyin .trin__tgl-asli')
    expect(CSS).toContain('opacity: 0;')
  })
})

describe('§5 Update Total dan Total Retention Amount', () => {
  it('panelnya dua kolom, dan judul kolom pertama adalah nama panelnya', () => {
    expect(TOTAL_RETENSI.judul).toBe('Total Retention Amount')
    expect(TOTAL_RETENSI.kolomNilai).toBe('Value')
    expect(TOTAL_RETENSI.tanpaBaris).toBe('No items')
  })

  // ⛔ TOMBOLNYA HIDUP 6 Oktober 2026, dan pernyataan ini ikut berubah.
  //
  // Bentuk lama menuntutnya MATI, dan itu benar selama rumusnya belum ada:
  // tombol hidup yang tidak menghitung apa pun lebih buruk daripada tombol
  // mati yang jujur. Rumusnya kini ada (`hitung_retensi.go`), jadi yang
  // dijaga berpindah — bukan "mati", melainkan "hidup, dan satu-satunya
  // yang mematikannya syarat yang TERUKUR di ekspor".
  it('tombolnya HIDUP, dan hanya `pyDisabledWhen` yang mematikannya', () => {
    expect(TOTAL_RETENSI.perbarui).toBe('Update Total')
    expect(RETENSI_SRC).toContain('{TOTAL_RETENSI.perbarui}')
    const i = RETENSI_SRC.indexOf('{TOTAL_RETENSI.perbarui}')
    const sebelum = RETENSI_SRC.slice(i - 400, i)
    // `pyDisabledWhen` = `TreatyIn.EDMMaterialType = 2` — @202558.
    expect(sebelum).toContain('disabled={totalTerkunci(edmJenisMaterial)}')
    expect(totalTerkunci('2')).toBe(true)
    expect(totalTerkunci('1')).toBe(false)
    // ⚠️ Tombol KEDUA tepat di atasnya di ekspor ber-`pyCondition 1=2`.
    // Ia MATI dan TIDAK dibangun — nol tombol kedua di berkas ini.
    expect((RETENSI_SRC.match(/<button/g) ?? []).length).toBe(1)
  })

  it('keempat tab lain yang juga punya Update Total TERCATAT', () => {
    // ⚠️ Yang dibangun ronde ini satu; yang terukur lima.
    expect([...TOTAL_RETENSI.tabLain].sort()).toEqual(['EGNPI', 'Installment', 'Limits', 'Share'])
  })

  it('rumusnya tidak dihitung ulang di layar', () => {
    // ⛔ Layar MEMANGGIL rumusnya, nol menghitungnya. Total tersimpan
    // dokumen tetap dibaca sebagai nilai AWAL — Pega memperlihatkan total
    // yang sudah ada di clipboard, bukan panel kosong.
    expect(FORM).toContain('warisan?.totalRetensi')
    expect(RETENSI_SRC).toContain('hitungRetensi(')
    expect(FORM).not.toContain('reduce((t, b) => t + Number(b.jumlah)')
    // Nol aritmetika di berkas tabnya pula.
    const kode = RETENSI_SRC.split(String.fromCharCode(10))
      .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
      .join(String.fromCharCode(10))
    expect(kode).not.toContain('reduce(')
    expect(kode).not.toContain('+ Number(')
  })
})

// ⛔ NOL JEJAK DI STRIP TAB — keputusan pemilik proses 8 Oktober 2026.
//
// ⚠️ Mekanisme `TAB_DISEMBUNYIKAN` ikut DICABUT, bukan sekadar dikosongkan:
// saringan yang nol menyaring apa pun adalah undangan bagi ronde berikutnya
// untuk memasukkan sesuatu ke dalamnya diam-diam.
describe('⛔ tab yang dikeluarkan nol jejak di strip', () => {
  it('daftar tab dan peta syaratnya bersih', () => {
    expect([...TAB_PROPORSIONAL, ...TAB_NON_PROPORSIONAL]).not.toContain('Retro')
    expect(Object.keys(SYARAT_TAB_PROPORSIONAL)).not.toContain('Retro')
    expect(Object.keys(SYARAT_TAB_NON_PROPORSIONAL)).not.toContain('Retro')
  })

  it('nol saringan tab tersisa di halaman', () => {
    const form = readFileSync(join(__dirname, 'pages/FormKontrakTreatyIn.tsx'), 'utf8')
    expect(form).not.toContain('TAB_DISEMBUNYIKAN')
    expect(form).toContain('const tab = tabUntuk(jenis, syaratTab)')
  })
})
