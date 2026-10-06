// Penjaga ronde "layar disamakan dengan dokumen desain" (5 Oktober 2026).
//
// ⛔ Sumber setiap angka di berkas ini adalah salah satu dari 43 tangkapan
// layar `D:\NUSANTARA RE APP\design-treaty-in-gambar\`, dan nomor gambarnya
// ditulis di tiap uji. Uji yang tidak dapat menyebut gambarnya tidak punya
// tempat di sini.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import {
  KOLOM_AKUMULASI,
  KOLOM_EGNPI,
  KOLOM_LIHAT_BERKAS,
  KOLOM_POLIS_PRODUKSI,
  KOLOM_PORTOFOLIO,
  PORTOFOLIO,
  KOLOM_RETENSI,
  JENIS_EGNPI,
  JENIS_POLIS_PRODUKSI,
  JENIS_RETENSI,
  LAMPIRAN,
  NAMA_KATEGORI_LAMPIRAN_NON_PROP,
  NAMA_KATEGORI_LAMPIRAN_PROP,
  POLIS_PRODUKSI,
  SUB_TAB_SHARE,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  golongan,
  JENIS_LIMITS,
  JENIS_SHARE,
  JENIS_RNM_SHARE,
  desimalPadan,
  LIMITS_POHON,
  GRID_BERTOMBOL_TAMBAH,
  TOMBOL_TAMBAH_BELUM_BERGRID,
} from './labels'
import { padankanDesimal, selAngka } from './pages/FormKontrakTreatyIn'

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

describe('§2 panel Existing Policy for Master ID', () => {
  it('judul dan keempat kolomnya disalin dari gambar 01/26', () => {
    expect(POLIS_PRODUKSI.judul).toBe('Existing Policy for Master ID')
    expect([...KOLOM_POLIS_PRODUKSI]).toEqual(['Policy No', 'Pega ID', 'Quarter', 'Quarter Year'])
    expect(POLIS_PRODUKSI.tanpaBaris).toBe('No items')
  })

  it('keempat kolomnya TEKS — pengenal tidak pernah diformat', () => {
    expect(JENIS_POLIS_PRODUKSI).toHaveLength(KOLOM_POLIS_PRODUKSI.length)
    expect(JENIS_POLIS_PRODUKSI.every((j) => j === 'teks')).toBe(true)
  })

  it('panelnya dirender di kepala layar, kedua cabang', () => {
    // ⛔ Di LUAR cabang `tabTampil` — ia bukan isi tab mana pun.
    expect(FORM).toContain('<PanelPolisProduksi baris={warisan?.polisProduksi ?? []} />')
    const iPanel = FORM.indexOf('<PanelPolisProduksi baris=')
    const iTab = FORM.indexOf("tabTampil === 'Reporting Period'")
    expect(iPanel).toBeGreaterThan(0)
    expect(iPanel).toBeLessThan(iTab)
  })

  it('lebarnya dibatasi dan turun penuh di layar sempit', () => {
    expect(CSS).toContain('.treatyin .trin__polis')
    expect(CSS).toContain('max-width: 34rem')
    expect(CSS).toContain('@media (max-width: 52rem)')
  })
})

describe('§0 judul kolom disalin dari gambarnya, bukan dari kunci dokumen', () => {
  it('Portfolio — gambar 02', () => {
    expect([...KOLOM_PORTOFOLIO]).toEqual([
      'Portfolio Type',
      'Premium / Loss Type',
      'Description',
    ])
  })

  it('Accumulation — gambar 19', () => {
    expect([...KOLOM_AKUMULASI]).toEqual([
      'Period',
      'Reporting Date',
      'Submission Days',
      'Submission Due',
    ])
  })

  it('EGNPI — urut layar gambar 29, bukan abjad', () => {
    // ⭐ Enam kolom pertama persis urutan gambarnya.
    expect(KOLOM_EGNPI.slice(0, 6)).toEqual([
      'Treaty Group',
      'As Date',
      'Proportion %',
      'Currency',
      'Amount',
      'Amount in IDR',
    ])
    // ⛔ Dua kolom yang tidak terlihat di gambar TIDAK dihapus.
    expect(KOLOM_EGNPI).toContain('Class of Business')
    expect(KOLOM_EGNPI).toContain('Note')
    expect(JENIS_EGNPI).toHaveLength(KOLOM_EGNPI.length)
  })

  it('Maximum Retention — urut layar gambar 26', () => {
    expect(KOLOM_RETENSI.slice(0, 3)).toEqual(['Treaty Group', 'Currency', 'Amount'])
    expect(KOLOM_RETENSI).toContain('Class of Business')
    expect(KOLOM_RETENSI).toContain('Note')
    expect(JENIS_RETENSI).toHaveLength(KOLOM_RETENSI.length)
    // Kolom uangnya pindah posisi, dan golongannya ikut pindah.
    expect(golongan(JENIS_RETENSI[KOLOM_RETENSI.indexOf('Amount')]!)).toBe('uang')
  })
})

describe('§6 RNM Share adalah SUB-TAB, dan nol tab dihapus', () => {
  it('strip sub-tab ada di dalam tab Share — gambar 16/17/34', () => {
    expect([...SUB_TAB_SHARE]).toEqual(['RNM Share'])
    expect(FORM).toContain('function SubTabShare(')
    expect(FORM).toContain('<StripTab tab={daftar} aktif={tampil} onPilih={setSub} />')
  })

  it('⛔ RNM Share TETAP di daftar tab — §6 melarang menghapusnya', () => {
    expect(TAB_NON_PROPORSIONAL).toContain('RNM Share')
    expect(TAB_NON_PROPORSIONAL).toHaveLength(12)
    expect(TAB_PROPORSIONAL).toHaveLength(11)
  })

  it('cabang tab RNM Share masih merender sesuatu', () => {
    // Tab yang dapat dipilih dan tidak merender apa pun adalah layar rusak.
    expect(FORM).toContain("tabTampil === 'RNM Share'")
    expect(FORM).toContain('<PanelRnmShare')
  })

  it('gridnya SATU bentuk, dipakai dua tempat', () => {
    // ⛔ Dua salinan grid yang sama berarti satu yang basi.
    expect(FORM.split('function PanelRnmShare(')).toHaveLength(2)
    expect(FORM.split('KOLOM_RNM_SHARE}').length - 1).toBe(1)
  })
})

describe('§4 permintaan perubahan — View File menjadi pop-up', () => {
  it('memakai Modal dari inti, bukan modal kedua', () => {
    expect(FORM).toContain('Modal,')
    expect(FORM).toContain('judul={LAMPIRAN.judulLihatBerkas}')
    expect(FORM).not.toContain('function ModalLihatBerkas(')
  })

  it('isinya DUA kolom, persis gambar 25', () => {
    expect([...KOLOM_LIHAT_BERKAS]).toEqual(['File Name', 'Type'])
  })

  it('tombol View hidup, tombol Upload mati', () => {
    expect(LAMPIRAN.lihatBerkas).toBe('View')
    expect(LAMPIRAN.unggah).toBe('Upload')
    const iUnggah = FORM.indexOf('{LAMPIRAN.unggah}')
    expect(FORM.slice(iUnggah - 160, iUnggah)).toContain('disabled')
    const iLihat = FORM.indexOf('{LAMPIRAN.lihatBerkas}')
    expect(FORM.slice(iLihat - 220, iLihat)).toContain('setBerkasDilihat')
  })

  it('penyaringnya memakai KODE kategori, bukan namanya', () => {
    // Nama kategori yang belum dipastikan tampil sebagai kodenya; menyaring
    // menurut yang tampil berarti menyaring menurut dua hal berbeda.
    expect(FORM).toContain('b.kodeKategori === kategori.kode')
  })
})

describe('§8 katalog kategori lampiran BERCABANG — gambar 24 lawan 42', () => {
  it('sebelas nama, dan tepat satu berbeda', () => {
    expect(NAMA_KATEGORI_LAMPIRAN_PROP).toHaveLength(11)
    expect(NAMA_KATEGORI_LAMPIRAN_NON_PROP).toHaveLength(11)
    const beda = NAMA_KATEGORI_LAMPIRAN_PROP.filter(
      (n) => !(NAMA_KATEGORI_LAMPIRAN_NON_PROP as readonly string[]).includes(n),
    )
    expect(beda).toEqual(['Pega Proportional Calculation /Perhitungan Pega Proportional'])
  })

  it('⛔ daftar nama TIDAK menutup pertanyaan kode↔nama', () => {
    // Urutan layarnya alfabetis, dan ronde 4 Oktober 2026 melarang keras
    // menyimpulkan kode dari urutan abjad.
    expect(LAMPIRAN.namaBelumBerumah).toContain('PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md')
  })
})

describe('§2 Treaty Year — tahun polos, TIDAK diformat', () => {
  it('nilainya lewat apa adanya, nol pemformat', () => {
    // Gambar 01: `Treaty Year 2025`. `selAngka` akan memberi `2.025`.
    const i = FORM.indexOf('const tahunTreaty =')
    expect(FORM.slice(i, i + 120)).toContain("warisan?.tahunTreaty ?? ''")
    expect(FORM.slice(i, i + 120)).not.toContain('selAngka')
  })

  it('dugaan 05/1 sebagai turunan dicabut di komentarnya', () => {
    expect(FORM).toContain('tahun polos')
    expect(FORM).not.toContain('PERTANYAAN TERBUKA: tangkapan layar Pega memperlihatkan `05/1`')
  })
})

describe('§24 desimal PER KOLOM, nol di ekor dipertahankan', () => {
  it('nol ditambahkan sampai presisi kolomnya', () => {
    // Gambar 01: `Value to IDR` berbunyi `1,00` dan `16.000,00`.
    expect(selAngka(['uang', 2], '1')).toBe('1,00')
    expect(selAngka(['uang', 2], '16000')).toBe('16.000,00')
    expect(selAngka(['uang', 2], '16000.5')).toBe('16.000,50')
  })

  it('⛔ NOL desimal berarti nol — bukan "tidak diketahui"', () => {
    // Gambar 29: `Amount in IDR` berbunyi `137.849.315.068`, tanpa koma.
    expect(selAngka(['uang', 0], '137849315068')).toBe('137.849.315.068')
    expect(selAngka(['uang', 0], '137849315068.00')).toBe('137.849.315.068')
  })

  it('⛔ yang SUDAH lebih panjang tidak dipotong — memotong itu membulatkan', () => {
    // ⚠️ DUA batas yang berbeda, dan keduanya berlaku berurutan:
    //
    //   `formatNumber` memotong di `DESIMAL_UANG` (4) — batas ATAS uang,
    //                  aturan lama yang §24 tidak sentuh;
    //   `padankanDesimal` menambah nol sampai presisi kolom (2) — dan
    //                  karena yang tiba sudah 4 desimal, ia tidak berbuat
    //                  apa-apa.
    //
    // Jadi hasilnya 4 desimal, bukan 2: pemadanan MENAMBAH, tidak pernah
    // memotong. Memotong berarti membulatkan, dan membulatkan menghilangkan
    // digit berarti.
    expect(selAngka(['uang', 2], '1500.12345')).toBe('1.500,1235')
    expect(padankanDesimal('1.500,12345', 2)).toBe('1.500,12345')
    expect(padankanDesimal('1.500,1235', 2)).toBe('1.500,1235')
  })

  it('tanpa presisi kolom, aturan ronde 69 tetap berlaku — nol di ekor dibuang', () => {
    // ⛔ Inilah yang menjaga §24 TIDAK bocor ke kolom yang tidak terbaca.
    expect(selAngka('uang', '1')).toBe('1')
    expect(selAngka('uang', '2484250')).toBe('2.484.250')
    expect(padankanDesimal('2.484.250', null)).toBe('2.484.250')
  })

  it('persen: nol masuk SEBELUM tanda %', () => {
    expect(selAngka(['persenShare', 2], '25')).toBe('25,00%')
    expect(padankanDesimal('100%', 4)).toBe('100,0000%')
  })

  it('⛔ teks bukan-angka lewat apa adanya — nol padanan', () => {
    for (const t of ['', '>=30% up to < 50%', 'n/a', '-', 'IDR']) {
      expect(padankanDesimal(t, 2)).toBe(t)
    }
  })

  it('kolom yang desimalnya terbaca membawa ANGKA, bukan hanya golongan', () => {
    // Gambar 29: dua kolom di baris yang SAMA, presisi berbeda.
    expect(desimalPadan(JENIS_EGNPI[KOLOM_EGNPI.indexOf('Amount')]!)).toBe(2)
    expect(desimalPadan(JENIS_EGNPI[KOLOM_EGNPI.indexOf('Amount in IDR')]!)).toBe(0)
    // Gambar 26: `Amount` di grid Maximum Retention nol desimal.
    expect(desimalPadan(JENIS_RETENSI[KOLOM_RETENSI.indexOf('Amount')]!)).toBe(0)
  })

  it('⚠️ kolom yang TIDAK terbaca di gambar tetap null — bukan ditebak', () => {
    // Nol dari 43 gambar memperlihatkan presisi kolom-kolom ini.
    for (const j of JENIS_LIMITS) expect(desimalPadan(j)).toBeNull()
    for (const j of JENIS_SHARE) expect(desimalPadan(j)).toBeNull()
    for (const j of JENIS_RNM_SHARE) expect(desimalPadan(j)).toBeNull()
  })

  it('⛔ format.ts dipanggil, tidak ditulis ulang — pemadanan di lapis modul', () => {
    expect(FORM).toContain("from '../../../../inti/frontend/lib/format'")
    expect(FORM).toContain('export function padankanDesimal(')
    // Pemadanan TIDAK menguraikan angka: nol pemformat kedua.
    expect(FORM).not.toMatch(/function formatNumber|function formatPersen/)
  })
})

describe('§1 tab Limits — pohon tiga tingkat', () => {
  it('puncaknya BERBEDA per cabang — gambar 03 lawan 30', () => {
    expect(LIMITS_POHON.puncakProp).toBe('Kind of Treaty')
    expect(LIMITS_POHON.puncakNonProp).toBe('Layers')
    expect(FORM).toContain('nonProp ? LIMITS_POHON.puncakNonProp : LIMITS_POHON.puncakProp')
  })

  it('⛔ medan pengelompokannya pun berbeda, bukan hanya judulnya', () => {
    // `jenisTreaty` terisi pada prop (1.360 elemen), nil pada 2.847 dari
    // 2.850 non-prop; `layer`/`layerType` sebaliknya. Satu medan untuk
    // keduanya memberi satu kelompok bernama kosong di cabang seberang.
    expect(FORM).toContain('const k = nonProp ? b.layer : b.jenisTreaty')
  })

  it('ketiga tingkatnya dirender, dan yang terdalam kelas bisnis', () => {
    expect(LIMITS_POHON.kelompokTreaty).toBe('Treaty Group')
    expect(LIMITS_POHON.kelasBisnis).toBe('Class of Business')
    expect(FORM).toContain('b.kelasBisnis.length === 0')
    expect(FORM).toContain('{LIMITS_POHON.kelasBisnis}')
  })

  it('dapat dilipat lewat <details>, bukan tiruan buatan sendiri', () => {
    // Papan ketik dan pembaca layar sudah menanganinya tanpa satu ARIA pun.
    expect(FORM).toContain('<details key={p.kunci}')
    expect(FORM).toContain('trin__pohon-tingkat')
    expect(CSS).toContain('.treatyin .trin__pohon-tingkat > summary')
  })

  it('⛔ keempat tombolnya MATI — nol jalur tulis', () => {
    expect(LIMITS_POHON.tambahLayer).toBe('add Layer') // huruf kecil, ejaan ekspor
    expect(LIMITS_POHON.tambahKelompok).toBe('Add Treaty Group')
    const i = FORM.indexOf('function PohonLimits(')
    const blok = FORM.slice(i, FORM.indexOf('function PerincianLayer('))
    expect((blok.match(/disabled/g) ?? []).length).toBe(2)
    expect(blok).not.toContain('onClick')
  })

  it('⚠️ Deductible ditandai 70,6% — bukan ditampilkan seolah pasti', () => {
    expect(LIMITS_POHON.belumPasti).toContain('70,6')
    expect(LIMITS_POHON.deductiblePetunjuk).toContain('1.340 kontrak')
    expect(FORM).toContain('belumPasti: true')
    // Dan HANYA medan itu yang ditandai.
    const i = FORM.indexOf('function PerincianLayer(')
    const blok = FORM.slice(i, i + 2000)
    expect((blok.match(/belumPasti: true/g) ?? []).length).toBe(1)
  })

  it('urutan kelompok = urutan dokumen, bukan abjad', () => {
    expect(FORM).not.toContain('puncak.sort(')
    expect(FORM).toContain('else puncak.push(')
  })
})

describe('§1 tombol Add — tiap grid yang gambarnya memperlihatkannya', () => {
  it('daftarnya berkunci NAMA TAB, supaya tab berikutnya tidak lupa', () => {
    // ⛔ RALAT 6 Oktober 2026 — `Portfolio` MASUK daftar.
    //
    // Daftar ini dibaca dari gambar, dan gambar `02` tangkapan layar
    // mode-BACA. Ekspor `Section/TreatyInTabsProportional.xml` menunjukkan
    // tombolnya ADA: sel `Add` 110 dan `Delete` 115, `pyVisible` = `OTHER`
    // dengan `pyCondition` = `TreatyIn.IsEditData!='1'`.
    expect([...GRID_BERTOMBOL_TAMBAH].sort()).toEqual(['Accumulation', 'EGNPI', 'Portfolio'])
  })

  it('⛔ MEKANISMENYA DIGANTI 6 Oktober 2026 — tombol mengikuti MODE', () => {
    // ⭐ Permintaan pemakai yang lebih baru mengalahkan daftar ini:
    //   "apabila di klik edit maka tiap function itu bisa digunakan semuanya
    //    seperti add pada tiap table, kalau view baru tidak bisa"
    //
    // Jadi `Add` tidak lagi dipasang per tab dari daftar gambar; ia muncul
    // pada SETIAP grid di mode `ubah` dan pada NOL grid di mode `lihat` —
    // dan ekspor mendukungnya: tombol-tombol itu dijaga
    // `TreatyIn.ViewState !='1'`.
    //
    // ⛔ Daftarnya TIDAK dihapus: ia tetap catatan grid mana yang gambar
    // mode-BACA perlihatkan bertombol, dan itu bukti terpisah dari mode.
    expect(FORM).toContain("mode === 'ubah'")
    expect(FORM).toContain('bisaTambah')
  })

  it('⛔ dua grid yang gambarnya TIDAK punya Add tetap tanpa tombol', () => {
    // Memasangnya di sana membuat layar berbeda dari layar lama.
    //
    // ⚠️ Dua, bukan tiga. `Portfolio` KELUAR dari daftar ini 6 Oktober
    // 2026 sebab ekspor membuktikan tombolnya ada — gambarnya yang
    // mode-BACA. Kedua yang tersisa BELUM diperiksa terhadap ekspor dengan
    // cara yang sama, dan itu dikatakan apa adanya: mereka tinggal di sini
    // atas dasar gambar saja, bukan atas dasar bukti bahwa tombolnya nihil.
    for (const tab of ['Co-Ins Scale', 'Maximum Retention']) {
      expect(GRID_BERTOMBOL_TAMBAH).not.toContain(tab)
    }
  })

  it('yang gambarnya punya tombol tetapi GRIDNYA belum ada: didaftarkan', () => {
    expect(TOMBOL_TAMBAH_BELUM_BERGRID).toHaveLength(3)
    expect(TOMBOL_TAMBAH_BELUM_BERGRID.join(' ')).toContain('Reinsurer')
    expect(TOMBOL_TAMBAH_BELUM_BERGRID.join(' ')).toContain('Update Value')
  })
})

describe('§1.2 pemilih bercari — papan tik, pola ditiru bukan diimpor', () => {
  const DD = readFileSync(join(AKAR, 'components', 'DropdownWarisan.tsx'), 'utf8')

  it('⛔ pola DITIRU, kode TIDAK diimpor dari modul lain', () => {
    // ⚠️ Namanya MEMANG disebut — di komentar, sebagai asal polanya. Yang
    // dilarang IMPORNYA, dan itu yang diuji: nol baris `import` menyebut
    // modul mana pun selain `inti` dan modul ini sendiri.
    expect(DD).toContain('DropdownCari.tsx')
    for (const ln of DD.split(String.fromCharCode(10))) {
      if (ln.trimStart().startsWith('import ')) {
        // Yang DIIZINKAN: paket (`react`), `inti`, dan berkas modul ini
        // sendiri. Selain itu impor lintas modul, dan itu dilarang.
        const m = /from '([^']+)'/.exec(ln)
        if (m !== null) {
          const jalur = m[1]!
          const sah =
            !jalur.startsWith('.') ||
            jalur.includes('/inti/frontend/') ||
            /^\.\.?\/[^.]/.test(jalur)
          expect(`${jalur} sah=${String(sah)}`).toContain('sah=true')
        }
      }
    }
  })

  it('kelima tombol papan tik ditangani', () => {
    for (const k of ['ArrowDown', 'ArrowUp', 'PageDown', 'PageUp', 'Enter', 'Escape', 'Tab']) {
      expect(DD).toContain(`'${k}'`)
    }
  })

  it('⛔ geseran DIJEPIT, tidak melingkar', () => {
    // Melingkar membuat yang menekan terus tidak pernah tahu ia di ujung.
    expect(DD).toContain('Math.min(cocok.length - 1, Math.max(0, a + langkah))')
  })

  it('baris aktif disebut kepada pembaca layar dan digulirkan', () => {
    expect(DD).toContain('aria-activedescendant')
    expect(DD).toContain('scrollIntoView')
    expect(DD).toContain('aria-controls')
  })

  it('Escape menutup TANPA memilih', () => {
    const i = DD.indexOf("e.key === 'Escape'")
    const blok = DD.slice(i, DD.indexOf("e.key === 'Tab'", i))
    expect(blok).toContain('setBuka(false)')
    expect(blok).toContain('setAktif(-1)')
    // ⛔ Nol pemilihan di cabang itu.
    expect(blok).not.toContain('pilih(')
  })

  it('⚠️ ketikan baru mengosongkan baris aktif', () => {
    // Baris ke-3 daftar LAMA bukan baris ke-3 daftar baru.
    const i = DD.indexOf('setKetik(e.target.value)')
    expect(DD.slice(i, i + 400)).toContain('setAktif(-1)')
  })

  it('⭐ nama KEMBAR tetap dibedakan pengenalnya — beda yang milik kita', () => {
    expect(DD).toContain('b.kembar ? `${b.nama} — ${b.id}` : b.nama')
  })

  it('⚠️ beda yang DISENGAJA dari modul contoh dicatat: dapat diketik', () => {
    expect(DD).toContain('MENOLAK ketikan bebas')
    expect(DD).toContain('HARUS')
  })

  it('baris aktif punya tanda yang BUKAN hover', () => {
    expect(DD).toContain('trin__ketikpilih-baris--aktif')
    expect(CSS).toContain('.treatyin .trin__ketikpilih-baris--aktif')
  })
})

// ---------------------------------------------------------------------
// §25 tab Portfolio — bentuknya dari EKSPOR, pilihannya dari PENGUKURAN
// ---------------------------------------------------------------------
//
// ⛔ Berkas ini biasanya menuntut tiap uji menyebut GAMBARNYA. Blok ini
// pengecualian yang dinyatakan: gambar `02` tangkapan layar mode-BACA,
// dan mode-baca tidak memperlihatkan kendali apa pun — nol dropdown, nol
// tombol. Jadi sumbernya di sini `Section/TreatyInTabsProportional.xml`,
// disebut sel demi sel, dan gambar `02` tetap menjadi sumber JUDULNYA.
describe('§25 tab Portfolio — dua daftar pilihan, satu area teks', () => {
  const TP = readFileSync(join(AKAR, 'components', 'TabPortofolio.tsx'), 'utf8')

  it('⛔ kolom 1 memakai TypePortfolio, kolom 2 memakai Type — TIDAK tertukar', () => {
    // Pasangan sel di ekspor: 107↔112 (`.TypePortfolio`), 108↔113 (`.Type`).
    // Saksi kedua: judul kolom 2 MENYEBUT kedua nilainya, dan yang bernilai
    // `Premium`/`Loss` adalah kunci `Type`.
    expect([...KOLOM_PORTOFOLIO]).toEqual(['Portfolio Type', 'Premium / Loss Type', 'Description'])
    const k1 = TP.indexOf('arah: b.jenisPortfolio')
    const k2 = TP.indexOf('jenis: b.jenis')
    expect(k1).toBeGreaterThan(-1)
    expect(k2).toBeGreaterThan(k1)
    // ⛔ Dan `arah` BENAR-BENAR dirender di kolom pertama.
    expect(TP.indexOf('nilai={b.arah}')).toBeLessThan(TP.indexOf('nilai={b.jenis}'))
  })

  it('pilihannya DIUKUR, bukan dikarang — dan hanya itu nilainya', () => {
    // Sapuan 6 Oktober 2026 atas seluruh 1.855 dokumen `M_TREATY_IN`:
    // 1.925 elemen `Portfolio[]` di 845 dokumen, nol gagal urai, nol kosong.
    //   TypePortfolio  Withdrawal 1.578 · Assumption 347
    //   Type           Premium    1.032 · Loss       893
    expect([...PORTOFOLIO.opsiArah]).toEqual(['Withdrawal', 'Assumption'])
    expect([...PORTOFOLIO.opsiJenis]).toEqual(['Premium', 'Loss'])
  })

  it('⛔ keduanya <select>, bukan kotak teks', () => {
    expect(TP).toContain('<select')
    expect(TP).toContain('opsi={PORTOFOLIO.opsiArah}')
    expect(TP).toContain('opsi={PORTOFOLIO.opsiJenis}')
    // ⛔ Nol `<input type="text">` tersisa di tab ini — itulah cacat yang
    // ronde ini tutup.
    expect(TP).not.toContain("type=\"text\"")
  })

  it('Description area teks — sel 114 ber-`pyFormat` pxTextArea', () => {
    expect(TP).toContain('<textarea')
  })

  it('⚠️ nilai di luar daftar TIDAK dijatuhkan diam-diam', () => {
    // Daftarnya terukur dari yang TERPAKAI, bukan dibaca dari aturan
    // properti (yang tidak ikut diekspor). Menjatuhkan nilai asing akan
    // mengubah kontrak yang sekadar dibuka.
    expect(TP).toContain('const asing =')
  })

  it('Add dan Delete mengikuti satu sakelar: IsEditData', () => {
    // Sel 110 dan 115: `pyCondition` = `TreatyIn.IsEditData!='1'`.
    // Sel 112·113·114: `pyReadOnlyCondition` = `TreatyIn.IsEditData= 1`.
    expect(TP).toContain('IsEditData')
    expect(TP).toContain('const bisaUbah = mode === ')
    expect(TP).toContain('{bisaUbah && (')
  })

  it('baris baru lahir KOSONG bertiga — seperti Activity-nya', () => {
    // `Activity/TreatyInPropAdd.xml`, langkah 1 (`param.Type=="portfolio"`):
    // satu-satunya penetapan adalah
    // `TreatyIn.Portfolio(<APPEND>).Description = ""`.
    expect(TP).toContain("{ arah: '', jenis: '', keterangan: '' }")
  })

  it('⛔ EDMMaterialType DICATAT, tidak dikarang', () => {
    // `pyDisabledWhen` = `TreatyIn.EDMMaterialType = 2` ada di keempat sel,
    // tetapi medannya belum terbaca layar ini. Menyalakannya berarti
    // menebak nilainya.
    expect(TP).toContain('EDMMaterialType')
  })
})
