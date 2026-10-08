// Penjaga layar Treaty In — ronde layar 1 (3 Oktober 2026).
//
// ⛔ Yang dijaga di sini bukan rupa layarnya, melainkan bahwa isinya tetap
// SALINAN. Tiap label, pesan galat, dan daftar tab datang dari ekspor Pega
// 2026-09; uji ini gagal pada hari seseorang memperhalus salah satunya.
//
// ⚠️ Uji ini membaca BERKAS SUMBER, bukan merender. Sebabnya: yang hendak
// dijaga adalah teks dan susunannya, dan merender menambah ketergantungan
// (jsdom, komponen inti) tanpa menambah satu pun hal yang dijaga.

import { DESIMAL_EGNPI } from './labelsEgnpi'
import { DESIMAL_RETENSI } from './labelsRetensi'

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import {
  DAFTAR_KONTRAK,
  FORM_KONTRAK,
  KOLOM_DAFTAR,
  REPORTING_PERIOD,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  KOLOM_LIMITS,
  KOLOM_SHARE,
  KOLOM_EVENT_LIMITS,
  KOLOM_RNM_SHARE,
  KOLOM_COIN_SCALE,
  LAMPIRAN,
  KOLOM_LAMPIRAN,
  KOLOM_HISTORY,
  KOLOM_IN2_SEMUA,
  KOLOM_IN2_TIDAK_DIPAKAI,
  KOLOM_EGNPI,
  KOLOM_RETENSI,
  KOLOM_ANGSURAN,
  JENIS_EGNPI,
  JENIS_RETENSI,
  JENIS_ANGSURAN,
  JENIS_LIMITS,
  JENIS_SHARE,
  JENIS_EVENT_LIMITS,
  JENIS_RNM_SHARE,
  JENIS_COIN_SCALE,
  golongan,
  NAMA_KATEGORI_LAMPIRAN_PROP,
} from './labels'
import { aksiUntuk, bolehRevisi, MEDAN_WARISAN } from './pages/DaftarKontrakTreatyIn'
import {
  NON_PROPORSIONAL,
  PROPORSIONAL,
  selAngka,
  tabUntuk,
  tataLetakKolom,
} from './pages/FormKontrakTreatyIn'

const AKAR = __dirname
const LABELS = readFileSync(join(AKAR, 'labels.ts'), 'utf8')
const MENU = readFileSync(join(AKAR, 'menu.ts'), 'utf8')
// ⭐ Layar ini DIPECAH 5 Oktober 2026: halaman + `components/`. `FORM`
// karena itu membaca KEDUANYA, supaya tiap pernyataan di bawah tetap
// menanyakan hal yang sama — "apakah kode layar modul ini memuat ini" —
// tanpa satu pun disunting. Itu bukti pemindahannya murni.
const SELA = String.fromCharCode(10)
// ⚠️ `HALAMAN` — teks halaman SAJA. Dipakai pernyataan yang menanyakan
// LETAK sesuatu DI DALAM halaman (irisan, `lastIndexOf`); `FORM` yang
// menggabungkan komponen akan menggeser letak itu dan membuat pernyataannya
// menanyakan hal yang berbeda.
const HALAMAN = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
const SARING = readFileSync(join(AKAR, 'saring.ts'), 'utf8')
const FORM =
  readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8') +
  readdirSync(join(AKAR, 'components'))
    .filter((f) => f.endsWith('.tsx') || f.endsWith('.ts'))
    .map((f) => readFileSync(join(AKAR, 'components', f), 'utf8'))
    .join(SELA)
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')

describe('layar daftar kontrak', () => {
  it('kesembilan kolom ada, dalam urutan layar lama', () => {
    expect(KOLOM_DAFTAR.map((k) => k.label)).toEqual([
      'ID',
      'Contract Name',
      'Reinsurance Type',
      'Source of Business',
      'Ceding',
      'Commencement',
      'Termination',
      'Position To',
      'Status Accept',
    ])
  })

  // ⛔ Pokok layar ini: tombolnya BERBEDA menurut keadaan baris, dan
  // meratakannya menghapus satu-satunya petunjuk bahwa sebuah kontrak sudah
  // tidak dapat disunting.
  it('baris Resolve Complete kehilangan Edit dan mendapat Copy + Revision', () => {
    expect(aksiUntuk('Resolve Complete')).toEqual(['View', 'Copy', 'Revision'])
  })

  it('baris yang masih dapat disunting mendapat Edit + View', () => {
    for (const keadaan of ['Open', 'Pending-Approval', 'DRAFT', '']) {
      expect(aksiUntuk(keadaan)).toEqual(['Edit', 'View'])
    }
  })

  it('kedua susunan tombol BERBEDA - penjaganya menggigit', () => {
    expect(aksiUntuk('Resolve Complete')).not.toEqual(aksiUntuk('Open'))
  })

  // ⚠️ RALAT 3 Oktober 2026. Uji ini dulu menuntut petunjuk kosong MENYEBUT
  // tiket 59 — benar ketika layar membaca `KONTRAK` (model baru, nol baris).
  // Sejak layar membaca `TREATY_IN` (1.854 baris) ia **tidak lagi benar**:
  // layar ini tidak menunggu tiket 59, dan petunjuk yang menunjuk tiket yang
  // salah membuat orang menagih pekerjaan yang tidak mengubah apa pun.
  // ⛔ DIBALIK 8 Oktober 2026 — kedua catatan DICABUT dari layar.
  it('⛔ kosongnya tidak lagi dijelaskan di layar', () => {
    expect(DAFTAR_KONTRAK.kosongPetunjuk).toBe('')
  })

  it('⛔ nomor tiket dan nama tabel tidak lagi tampil di layar', () => {
    expect(DAFTAR_KONTRAK.catatanSumber).toBe('')
    expect(DAFTAR_KONTRAK.keterangan).toBe('')
  })
})

describe('strip tab - DUA himpunan, dan keduanya tidak sama', () => {
  it('proporsional sebelas tab, non-proporsional dua belas', () => {
    expect(TAB_PROPORSIONAL).toHaveLength(11)
    expect(TAB_NON_PROPORSIONAL).toHaveLength(12)
  })

  // ⛔ Uji yang paling berguna di berkas ini. Perkiraan ronde ini: "keduanya
  // memuat kesepuluh namanya". Ekspor menyatakan sebaliknya, dan uji ini
  // membekukan temuan itu supaya tidak diratakan diam-diam nanti.
  it('hanya enam nama yang dimiliki KEDUA cabang', () => {
    const sama = TAB_PROPORSIONAL.filter((t) => (TAB_NON_PROPORSIONAL as readonly string[]).includes(t))
    expect(sama).toEqual([
      'Limits',
      'Share',
      'Retro',
      'Exclusions',
      'Special Conditions',
      'Information & Submit',
    ])
  })

  it('lima tab HANYA ada di proporsional, enam HANYA di non-proporsional', () => {
    const hanyaProp = TAB_PROPORSIONAL.filter((t) => !(TAB_NON_PROPORSIONAL as readonly string[]).includes(t))
    const hanyaNon = TAB_NON_PROPORSIONAL.filter((t) => !(TAB_PROPORSIONAL as readonly string[]).includes(t))
    expect(hanyaProp).toEqual([
      'Reporting Period',
      'Portfolio',
      'Co-Ins Scale',
      'Accumulation',
      'Achievement In IDR',
    ])
    expect(hanyaNon).toEqual([
      'Maximum Retention',
      'Event Limits',
      'EGNPI',
      'RNM Share',
      'Installment',
      'Value Difference',
    ])
  })

  it('Reporting Period hanya ada di cabang proporsional', () => {
    expect(TAB_PROPORSIONAL).toContain('Reporting Period')
    expect(TAB_NON_PROPORSIONAL).not.toContain('Reporting Period')
  })

  it('radio memilih himpunannya', () => {
    expect(tabUntuk(PROPORSIONAL)).toBe(TAB_PROPORSIONAL)
    expect(tabUntuk(NON_PROPORSIONAL)).toBe(TAB_NON_PROPORSIONAL)
  })
})

describe('teks yang disalin apa adanya', () => {
  // Pesan galat dirangkai Pega dari DUA rule: `pyCaption Start Date, Due`
  // dan `pyCaption . Must Not Be Empty`. Ejaannya tidak diperhalus.
  it('pesan galat Reporting Period persis seperti ekspor', () => {
    expect(REPORTING_PERIOD.galatKosong).toBe('Start Date, Due. Must Not Be Empty')
  })

  it('ejaan ekspor dipertahankan, termasuk yang tampak keliru', () => {
    // `Teritorial` satu R - begitu ekspornya, dan memperbaikinya membuat
    // pencarian atas layar lama tidak menemukan apa pun.
    expect(FORM_KONTRAK.lingkupWilayah).toBe('Teritorial Scope')
    // `TreatyIn.Bordeaux` adalah nama PROPERTINYA; labelnya `Bordereaux`.
    expect(FORM_KONTRAK.bordereaux).toBe('Bordereaux')
    expect(LABELS).toContain('TreatyIn.Bordeaux')
  })

  it('tiap medan membawa jejak posisi bitanya', () => {
    const jejak = LABELS.match(/@\d{4,}/g) ?? []
    expect(jejak.length).toBeGreaterThan(30)
  })
})

describe('blok bergaris mati tidak dibangun', () => {
  // Dua blok `pyContainerVisibleWhen = 1=2` ditemukan di
  // `TreatyInTabsProportional.xml`. Keduanya disebut di kepala berkas form
  // supaya yang membaca kode tahu apa yang SENGAJA tidak ada.
  it('keduanya disebut di kepala FormKontrakTreatyIn.tsx', () => {
    expect(FORM).toContain('1=2')
    expect(FORM).toContain('@1095716')
    expect(FORM).toContain('@1300558')
  })

  it('salinan HIDUP Account Reporting Period yang dipakai', () => {
    expect(FORM).toContain('@45828')
    expect(REPORTING_PERIOD.judul).toBe('Account Reporting Period')
  })
})

describe('L-4 dicabut, dan pencabutannya tercatat', () => {
  it('menu.ts menyebut pencabutannya berikut cacah berkasnya', () => {
    expect(MENU).toContain('L-4')
    expect(MENU).toContain('DICABUT')
    expect(MENU).toMatch(/50 berkas `Section\/` dan 3 `Harness\/`/)
  })

  it('halaman awal modul kini daftar kontrak', () => {
    expect(MENU).toContain("HALAMAN_AWAL_TREATYIN: HalamanTreatyIn = 'treatyin-kontrak'")
  })
})

// ===========================================================================
// RONDE LAYAR 2 — susunan, bukan medan
// ===========================================================================

const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')
const DAFTAR = readFileSync(join(AKAR, 'pages', 'DaftarKontrakTreatyIn.tsx'), 'utf8')
const ACUAN = readFileSync(join(AKAR, 'pages', 'AcuanTreatyIn.tsx'), 'utf8')

describe('form dua kolom yang mengalir sendiri', () => {
  // ⛔ Pokoknya: kolom kanan LEBIH PANJANG. Grid pengisi baris tidak dapat
  // menghasilkan itu tanpa sel kosong yang disisipkan tangan, dan sel kosong
  // itu bergeser diam-diam pada medan berikutnya yang ditambahkan.
  it('kolom kanan lebih panjang daripada kiri', () => {
    expect(tataLetakKolom.kiri).toHaveLength(5)
    expect(tataLetakKolom.kanan).toHaveLength(7)
    expect(tataLetakKolom.kanan.length).toBeGreaterThan(tataLetakKolom.kiri.length)
  })

  it('urutannya urutan rujukan', () => {
    expect(tataLetakKolom.kiri[0]).toBe('Treaty Contract Name')
    expect(tataLetakKolom.kanan[0]).toBe('Commencement')
    expect(tataLetakKolom.kanan.slice(-3)).toEqual([
      'Ceding',
      'RNM as Treaty Leader',
      'Source of Business',
    ])
  })

  it('nol medan hilang atau kembar di antara kedua kolom', () => {
    const semua = [...tataLetakKolom.kiri, ...tataLetakKolom.kanan]
    expect(new Set(semua).size).toBe(semua.length)
    expect(semua).toHaveLength(12)
  })

  // Pengisi `<span />` yang dipakai sebelumnya DICABUT; uji ini menjaganya
  // tidak kembali diam-diam.
  //
  // ⚠️ Komentar dibuang lebih dulu. Kepala berkas itu MENYEBUT `<span />`
  // untuk menjelaskan kenapa ia dicabut, dan sapuan yang memindai teks
  // mentah menemukan penjelasannya sendiri — pelajaran yang sama yang sudah
  // membuat `tanpaKomentarGo` lahir di penjaga arsip.
  it('form tidak lagi memakai sel pengisi', () => {
    const kode = FORM.replace(/\{\/\*[\s\S]*?\*\/\}/g, '').replace(/^\s*\/\/.*$/gm, '')
    expect(kode).toContain('trin__dwikolom')
    expect(kode).not.toContain('<span />')
  })
})

describe('lebar yang disetel membawa asal-usulnya', () => {
  // Perbandingan 193 : 349 : 80 adalah satu-satunya lebar yang ekspor
  // sungguh bedakan di layar ini.
  it('grid Rate of Exchange memakai perbandingan dari ekspor', () => {
    expect(FORM).toContain('trin__kol-mata-uang')
    // ⛔ RALAT 7 Oktober 2026 — rasio ekspor (27:48:12,5:12,5) DIGANTI
    // atas permintaan pemilik proses: kotak `Value to IDR` selebar setengah
    // tabel untuk angka sependek `286,07`, dan kolom tanggal terdesak.
    //
    // ⭐ Keempatnya kini berjumlah 84%, bukan 100% — sisanya untuk kolom
    // tombol Remove, yang di `colgroup` memang tidak punya `<col>`.
    expect(CSS).toMatch(/\.trin__kol-mata-uang\s*\{\s*width:\s*14%/)
    expect(CSS).toMatch(/\.trin__kol-nilai\s*\{\s*width:\s*28%/)
    expect(CSS).toMatch(/\.trin__kol-tanggal\s*\{\s*width:\s*21%/)
    expect(CSS).toContain('193')
    expect(CSS).toContain('349')
  })

  // ⛔ RALAT yang dijaga: lebar kolom daftar BUKAN dari ekspor. Kesembilan
  // sel baris saring `pyWidth=130`, seragam. Komentar yang mengaku
  // sebaliknya sudah dicabut, dan uji ini menahannya tetap tercabut.
  it('lebar kolom daftar TIDAK mengaku berasal dari ekspor', () => {
    expect(CSS).not.toMatch(/tidak seragam - mengikuti ekspor/)
    expect(CSS).toContain('pyWidth=130')
    expect(CSS).toContain('keputusan rupa KITA')
  })

  /**
   * ⛔ RALAT 3 Oktober 2026. Uji ini dulu menuntut `white-space: nowrap` atas
   * SELURUH sel `.trin__tabel`, membekukan pembacaan `pyWrap=false` yang benar
   * tentang Pega dan salah tentang kita.
   *
   * Yang sungguh terjadi di layar: dengan nama kontrak sepanjang "WHOLE ACCOUNT
   * RISK AND CATASTROPHE EXCESS OF LOSS REINSURANCE 2025", tabelnya melar jauh
   * melewati kartunya dan HANYA DUA dari sembilan kolom terlihat — tujuh
   * sisanya di balik gulir mendatar. Uji lama akan menolak setiap perbaikannya.
   *
   * Pega merender pada lebar tetap; tabel ini tidak. `pyWrap=false` karena itu
   * dihormati di tempat ia berarti — pengenal, jenis, dan tanggal tidak boleh
   * terpotong — dan tidak di tiga kolom teks panjang.
   */
  it('nowrap HANYA pada kolom pendek, tiga kolom teks membungkus', () => {
    // `sifatProporsi` DICABUT dari daftar ini 3 Oktober 2026: dua kata yang
    // meluber ke kolom sebelahnya, bukan satu kata yang rusak bila dipenggal.
    for (const k of ['id', 'tanggalMulai', 'tanggalBerakhir']) {
      expect(CSS).toMatch(new RegExp(`trin__k--${k}[^}]*white-space: nowrap`, 's'))
    }
    // Tiga kolom teks panjang TIDAK boleh nowrap.
    for (const k of ['namaKontrak', 'idAsalBisnis', 'idCedant', 'sifatProporsi']) {
      expect(CSS).not.toMatch(new RegExp(`trin__k--${k} \\{[^}]*white-space: nowrap`, 's'))
    }
    // Dan nowrap menyeluruh tidak boleh kembali.
    expect(CSS).not.toMatch(/\.trin__tabel td,\s*\n\s*\.treatyin \.trin__tabel th \{\s*\n\s*white-space: nowrap/)
  })

  /**
   * Kesembilan kolom harus MUAT. `table-layout: fixed` + lebar persen adalah
   * yang menjamin itu; tanpa keduanya kolom melebar mengikuti isi terpanjang.
   */
  it('kesembilan kolom punya lebar persen dan tabelnya fixed', () => {
    expect(CSS).toMatch(/\.treatyin \.trin__tabel \{[^}]*table-layout: fixed/s)
    const kunci = ['id', 'namaKontrak', 'sifatProporsi', 'idAsalBisnis', 'idCedant',
      'tanggalMulai', 'tanggalBerakhir', 'posisiKe', 'keadaanSiklusHidup']
    for (const k of kunci) {
      expect(CSS).toMatch(new RegExp(`trin__k--${k}[^}]*width: \\d+%`, 's'))
    }
  })

  /**
   * ⛔ Nama kelas lebar HARUS sama dengan `kunci` kolomnya. Enam di antaranya
   * pernah salah eja (`jenisReasuransi` lawan `sifatProporsi`, `mulai` lawan
   * `tanggalMulai`, …), sehingga aturannya mati tanpa suara: CSS yang tidak
   * cocok tidak melempar galat, ia hanya tidak berlaku.
   */
  it('tiap kunci kolom punya aturan lebarnya, dan ejaannya cocok', () => {
    const kunci = [...LABELS.matchAll(/kunci: '([a-zA-Z]+)'/g)].map((m) => m[1])
    expect(kunci.length).toBe(9)
    for (const k of kunci) {
      expect(CSS).toContain(`trin__k--${k}`)
    }
  })
})

describe('daftar: saring per kolom dan penomoran halaman', () => {
  it('kesembilan kolom punya ikon saringnya', () => {
    expect(DAFTAR).toContain('trin__ikon-saring')
    expect(DAFTAR).toContain('aria-pressed')
  })

  it('baris saring berisi kotak, bukan sel kosong', () => {
    expect(DAFTAR).toContain('trin__baris-saring')
    expect(DAFTAR).toMatch(/aria-label=\{`Saring \$\{k\.label\}`\}/)
  })

  // ⚠️ RALAT 3 Oktober 2026 — syaratnya berubah bersama sumber datanya.
  // Dengan 25 baris per halaman dari 1.854, `tampil.length` SELALU 25 dan
  // penomoran akan mengira hanya ada satu halaman. Yang menentukan `total`
  // dari server.
  it('penomoran halaman memakai total dari SERVER, bukan panjang halaman', () => {
    expect(DAFTAR).toContain('trin__kepala-kanan')
    expect(DAFTAR).toContain('<Halaman')
    expect(DAFTAR).toMatch(/total=\{total\}/)
    expect(DAFTAR).not.toMatch(/total=\{tampil\.length\}/)
  })

  // ⛔ Server sudah memotong halamannya. Memotong lagi di peramban membuat
  // halaman kedua menampilkan halaman kosong.
  it('nol pemotongan halaman kedua kali di peramban', () => {
    expect(DAFTAR).not.toContain('.slice((halaman - 1) * UKURAN_HALAMAN')
  })

  // Ukuran halaman keputusan kita, bukan salinan — dan ia menyatakannya.
  it('ukuran halaman dinyatakan sebagai keputusan kita', () => {
    expect(DAFTAR).toContain('UKURAN_HALAMAN')
    expect(DAFTAR).toContain('KEPUTUSAN KITA')
  })
})

describe('tab belum dibangun dibedakan dari data kosong', () => {
  // ⛔ Dua pernyataan yang berbeda: "belum ada KODE" versus "belum ada DATA".
  // Yang salah membacanya menagih pemuatan data padahal yang kurang layarnya.
  it('memakai .trin__belum, bukan komponen Kosong', () => {
    expect(FORM).toContain('trin__belum')
    const potong = HALAMAN.slice(HALAMAN.indexOf('tabTampil === '))
    expect(potong).not.toContain('<Kosong')
  })

  it('bentuknya berbeda di CSS - tepi putus-putus, bukan ikon', () => {
    expect(CSS).toMatch(/\.trin__belum \{[\s\S]*?border: 1px dashed/)
  })
})

describe('AcuanTreatyIn ikut bahasa visual yang sama', () => {
  it('memakai trin__tabel dan merender kepala kolom walau kosong', () => {
    expect(ACUAN).toContain('className="trin__tabel"')
    expect(ACUAN).toContain('colSpan={bersusun ? 4 : 3}')
  })
})

describe('pin .form-grid tidak merusak lebar telepon', () => {
  // ⛔ RALAT yang dijaga. Pin dua-kolom berkekhususan lebih tinggi daripada
  // aturan inti `@media (max-width: 860px) { .form-grid { 1fr } }`, jadi
  // tanpa media query sendiri ia menimpanya — dan Reporting Period bertahan
  // dua kolom di lebar telepon, persis cacat yang aturan inti itu cegah.
  it('pin hanya berlaku dari 861px ke atas', () => {
    expect(CSS).toMatch(/@media \(min-width: 861px\) \{\s*\n\s*\.treatyin \.form-grid \{/)
  })

  it('dwikolom runtuh satu kolom di bawah 900px', () => {
    expect(CSS).toMatch(/@media \(max-width: 899px\) \{\s*\n\s*\.treatyin \.trin__dwikolom \{\s*\n\s*grid-template-columns: minmax\(0, 1fr\)/)
  })
})

// ===========================================================================
// DATA NYATA — `POOLDATA.TREATY_IN`
// ===========================================================================

describe('layar daftar membaca tabel warisan', () => {
  it('memanggil jalur warisan, bukan jalur model baru', () => {
    expect(DAFTAR).toContain('ambilDaftarWarisan')
    expect(DAFTAR).not.toContain('ambilDaftarKontrak(')
  })

  // ⛔ Satu halaman per permintaan: 1.854 baris tidak ditarik sekaligus.
  it('memuat ulang setiap kali halamannya berganti', () => {
    expect(DAFTAR).toMatch(/\}, \[halaman\]\)/)
  })

  // Pemetaan kolom layar -> medan warisan berdiri sendiri supaya dapat diuji.
  it('kesembilan kolom terpetakan ke medan warisan', () => {
    const kunci = KOLOM_DAFTAR.map((k) => k.kunci)
    expect(Object.keys(MEDAN_WARISAN).sort()).toEqual([...kunci].sort())
    expect(MEDAN_WARISAN.idCedant).toBe('cedant')
    expect(MEDAN_WARISAN.idAsalBisnis).toBe('asalBisnis')
    expect(MEDAN_WARISAN.keadaanSiklusHidup).toBe('statusAkseptasi')
    expect(MEDAN_WARISAN.posisiKe).toBe('posisiKe')
  })

  // ⭐ Aturan tombol TIDAK ditulis ulang — ia aturan yang sama, kini diberi
  // makan oleh `STATUSAKSEPTASI` yang nyata. Sapuan menemukan EMPAT nilai,
  // bukan satu: `Resolve Complete` 1.820 · `Accept` 12 · `Decline` 11 ·
  // NULL 11. Ketiga yang terakhir 34 baris nyata, bukan kasus teoretis.
  // ⭐ `Revision` MENULIS sejak 7 Oktober 2026 — syarat lengkap cell 994:
  // workbasket Admin, Position kosong, Resolve Complete.
  it('Revision hanya bagi Admin atas kontrak tuntas yang tidak berposisi', () => {
    const tuntas = { statusAkseptasi: 'Resolve Complete', posisi: '' }
    expect(bolehRevisi(tuntas, ['ReasTreatyInAdmin'])).toBe(true)
    expect(bolehRevisi(tuntas, ['ReasTreatyInSecHead'])).toBe(false)
    expect(bolehRevisi({ ...tuntas, posisi: 'ReasTreatyInSecHead' }, ['ReasTreatyInAdmin'])).toBe(false)
    expect(bolehRevisi({ ...tuntas, statusAkseptasi: 'Accept' }, ['ReasTreatyInAdmin'])).toBe(false)
    expect(DAFTAR).toContain('bolehRevisi(b, workbasket)')
  })

  it('keempat nilai STATUSAKSEPTASI nyata memilih tombolnya', () => {
    expect(aksiUntuk('Resolve Complete')).toEqual(['View', 'Copy', 'Revision'])
    for (const lain of ['Accept', 'Decline', '']) {
      expect(aksiUntuk(lain)).toEqual(['Edit', 'View'])
    }
  })

  it('baris memakai statusAkseptasi, bukan medan model baru', () => {
    expect(DAFTAR).toContain('aksiUntuk(b.statusAkseptasi)')
  })

  // ⛔ Pengenal TEKS. `TREATY_IN.ID` adalah VARCHAR2(100).
  it('pengenal kontrak bertipe teks di seluruh jalurnya', () => {
    expect(DAFTAR).toContain('onBuka: (id: string, mode: ModeForm) => void')
    expect(FORM).toContain('idKontrak: string')
    expect(RUTE).toContain('useState<string | null>(null)')
  })
})

// ===========================================================================
// FORM MEMBACA KONTRAK WARISAN NYATA
// ===========================================================================

describe('form membaca kontrak warisan', () => {
  it('memanggil jalur warisan, bukan jalur model baru', () => {
    expect(FORM).toContain('ambilKontrakWarisan')
    expect(FORM).not.toContain('ambilKontrak(')
  })

  it('memuat ulang ketika pengenalnya berganti, dan melewati kontrak baru', () => {
    expect(FORM).toMatch(/\}, \[idKontrak\]\)/)
    expect(FORM).toContain("if (idKontrak === '')")
  })

  // ⭐ Radio disambungkan ke NILAI NYATA. Kontrak non-proporsional yang
  // membuka strip proporsional memperlihatkan sebelas tab yang tidak satu
  // pun miliknya.
  it('radio jenis diisi dari kontraknya, bukan dibiarkan bawaan', () => {
    expect(FORM).toMatch(/setJenis\(k\.sifatProporsi === NON_PROPORSIONAL/)
  })

  // ⭐ RALAT 7 Oktober 2026 — ketiga medan SELALU tampil, persis ekspor
  // `Section/TreatyInNONProportional.xml` (`ContractRefNo` /
  // `BordereauxNote` `pyVisible ALWAYS`, `TreatyLeader` tanpa syarat
  // tampil). Dahulu `MedanTakAda` "Tidak ada di dokumen sistem lama"
  // menggantikannya bila kuncinya tidak ada — syarat yang tidak dikenal Pega.
  it('Contract Ref No · Bordereaux Note · RNM as Treaty Leader selalu dirender — nol MedanTakAda', () => {
    expect(FORM).not.toContain('MedanTakAda')
    expect(FORM).not.toContain('adaKunci(')
    expect(FORM_KONTRAK).not.toHaveProperty('takAdaDiWarisan')
    // `pxTextInput` — baca-saja bila `TreatyIn.EDMMaterialType = 1`.
    expect(FORM).toMatch(/<Field\s+label=\{FORM_KONTRAK\.nomorRujukan\}[^>]*readOnly=\{edmMateri\}/)
    // `pxTextArea` — idem, berlabel lewat `htmlFor`.
    expect(FORM).toMatch(/htmlFor=\{idBordereauxNote\}/)
    expect(FORM).toMatch(/<textarea[^>]*value=\{bordereauxNote\}[^>]*readOnly=\{edmMateri\}/)
    // `pxCheckbox` dengan `pyCheckboxCaption`, `pyIncludeLabel=false`.
    expect(FORM).toMatch(/<label className="trin__centang">\s*<input\s+type="checkbox"\s+checked=\{pemimpin\}/)
    expect(FORM_KONTRAK.pemimpinTreaty).toBe('RNM as Treaty Leader')
  })

  // ⭐ `TREATYYEAR` adalah kolom — dibaca, bukan dipotong dari tanggal.
  it('Treaty Year dibaca dari kolomnya SAAT MEMBUKA kontrak', () => {
    // ⛔ Bukan dihitung dari Commencement saat membuka: 5 dari 1.854 kontrak
    // sengaja bertahun berbeda dari tahun mulainya, dan menghitung ulang akan
    // menimpa kelimanya tanpa suara.
    expect(FORM).toContain('setTahunTreaty(k.tahunTreaty)')
  })

  it('⛔ medan tanggal diisi dari bentuk NILAI, bukan bentuk BACA', () => {
    // `<input type="date">` hanya menerima `yyyy-mm-dd`. Disuapi `01/01/2025`
    // ia menolaknya TANPA BERSUARA dan menampilkan placeholder — medan yang
    // menolak nilainya terlihat persis seperti medan yang memang kosong.
    expect(FORM).toContain('setMulai(keKabel(k.tanggalMulaiAsli))')
    expect(FORM).toContain('setBerakhir(keKabel(k.tanggalBerakhirAsli))')
    expect(FORM).not.toContain('setMulai(k.tanggalMulai)')
    // ⛔ SATU bentuk kabel saja. Bentuk ISO di sini membuat kedua turunan
    // (`Treaty Year`, `Termination`) KOSONG begitu pemakai memilih tanggal,
    // sebab `FieldTanggal` mengembalikan `DD-MM-YYYY` — keluhan pemilik
    // proses 6 Oktober 2026.
    expect(FORM).not.toContain('keIso(')
  })

  // ⚠️ RALAT: nilai pilihan adalah yang TERSIMPAN, dan ada DUA per pilihan.
  it('pilihan membawa kedua nilai nyata, huruf kecil seperti tersimpan', () => {
    expect(FORM_KONTRAK.bordereauxNilai).toEqual(['reporting', 'nonreporting'])
    expect(FORM_KONTRAK.caraPembukuanNilai).toEqual(['underwriting', 'accounting'])
  })

  // Terjemahan TIDAK ditulis dua kali — ia dipakai ulang dari services.
  it('nol penerjemah kedua di frontend', () => {
    expect(FORM).not.toMatch(/NonProportional'\s*\?\s*'Non Proportional'/)
    expect(FORM).not.toContain('/19')
  })
})

describe('empat tab dari M_TREATY_IN2 — satu tabel, empat proyeksi', () => {
  it('ketiganya membaca `warisan.layer` — Event Limits KELUAR, dengan bukti', () => {
    for (const tab of ['Limits', 'Share', 'Event Limits', 'RNM Share']) {
      expect(FORM).toContain(`tabTampil === '${tab}'`)
    }
    // ⛔ Satu seam. Banyak pembacaan berarti banyak kueri untuk satu baca,
    // dan kesempatan agar yang satu melihat layer yang lain tidak.
    //
    // ⛔ RALAT 6 Oktober 2026: EMPAT menjadi TIGA, dan itu bukan pelonggaran.
    // Tab Event Limits Non-Prop di ekspor mengikat properti AKAR
    // (`TreatyIn.RSMDLimit`, `TreatyIn.Earthquake`, …) — satu set per
    // kontrak — bukan nilai per layer. Diukur atas seluruh 1.855 dokumen:
    // Non-Prop berisi di AKAR pada 49 kontrak, di `Detail[]` pada 1, tumpang
    // tindih 0. Membacanya dari `warisan.layer` membuat tab itu kosong pada
    // 771 dari 772 kontrak Non-Prop.
    //
    // ⛔ RALAT 7 Oktober 2026: TIGA menjadi DUA. Tab Share NON-PROP kini
    // membaca `warisan.shareNP` (pendaratan `T_TREATY_SHARE*`) — bentuk
    // ekspornya per baris `TreatyIn.Share`, bukan kolom `Detail[]` Prop yang
    // dulu ditampilkannya. Yang tersisa: Share Prop dan RNM Share.
    //
    // ⛔ RALAT 7 Oktober 2026: DUA menjadi SATU. Share PROP kini membaca
    // `TreatyIn.Limits` (penampung halaman, nilai awal `warisan.limitsPohon`)
    // — grid `Kind of Treaty` di ekspor ADALAH page list `TreatyIn.Limits`.
    // Yang tersisa hanya RNM Share.
    const pakai = FORM.match(/warisan\?\.layer \?\? \[\]/g) ?? []
    expect(pakai.length).toBe(1)
    expect(FORM).toContain('pohon={warisan?.limitsPohon ?? []}')
    expect(FORM).toContain('share={shareNP ?? warisan?.shareNP}')
    expect(FORM).toContain('<TabEventLimits mode={mode} />')
  })

  it('nol tabel baru dibuat untuk keempatnya', () => {
    for (const salah of [
      'M_TREATYIN_LIMITS',
      'M_TREATYIN_SHARE',
      'M_TREATYIN_EVENTLIMITS',
      'M_TREATYIN_RNMSHARE',
    ]) {
      expect(FORM).not.toContain(salah)
      expect(LABELS).not.toContain(salah)
    }
  })

  it('tiap kolom punya golongan angkanya, dan panjangnya sejajar', () => {
    expect(JENIS_LIMITS.length).toBe(KOLOM_LIMITS.length)
    expect(JENIS_SHARE.length).toBe(KOLOM_SHARE.length)
    expect(JENIS_EVENT_LIMITS.length).toBe(KOLOM_EVENT_LIMITS.length)
    expect(JENIS_RNM_SHARE.length).toBe(KOLOM_RNM_SHARE.length)
    expect(JENIS_COIN_SCALE.length).toBe(KOLOM_COIN_SCALE.length)
  })

  it('kolom layer tidak pernah diformat — ia pengenal, dan nol di depannya bermakna', () => {
    // ⚠️ Judulnya BERUBAH dari `LAYER` menjadi `Layer` di Event Limits,
    // sebab keempat tab tidak lagi memakai nama kolom `M_TREATY_IN2`
    // sebagai judul. Yang dijaga tetap sama: ia TEKS, bukan angka.
    for (const [kolom, jenis] of [
      [KOLOM_LIMITS, JENIS_LIMITS],
      [KOLOM_SHARE, JENIS_SHARE],
      [KOLOM_EVENT_LIMITS, JENIS_EVENT_LIMITS],
      [KOLOM_RNM_SHARE, JENIS_RNM_SHARE],
    ] as const) {
      const i = kolom.findIndex((k) => k.toUpperCase() === 'LAYER')
      expect(i).toBeGreaterThanOrEqual(0)
      expect(jenis[i]).toBe('teks')
    }
  })

  it('kesepuluh kolom kepala disebut namanya, bukan dibuang diam-diam', () => {
    expect(KOLOM_IN2_TIDAK_DIPAKAI.length).toBe(10)
    // ⛔ Dan tidak satu pun dari kesepuluhnya menyelinap ke grid: ia akan
    // mengulang nilai yang sama di setiap baris tanpa menambah keterangan.
    const semuaKolom = [
      ...KOLOM_LIMITS,
      ...KOLOM_SHARE,
      ...KOLOM_EVENT_LIMITS,
      ...KOLOM_RNM_SHARE,
    ]
    for (const k of KOLOM_IN2_TIDAK_DIPAKAI) {
      expect(semuaKolom).not.toContain(k)
    }
  })

  it('⛔ penghitungan ke-41 kolom DICABUT — sumbernya bukan tabel itu lagi', () => {
    // ⭐ Uji ini DIUBAH ARTINYA, bukan dihapus.
    //
    // Bentuk sebelumnya menjumlahkan 31 kolom terpakai + 10 kolom kepala =
    // 41 kolom `M_TREATY_IN2`, supaya kolom yang hilang dari layar TERLIHAT.
    // Penjumlahan itu kehilangan maknanya pada 5 Oktober 2026: `M_TREATY_IN2`
    // dicabut sebagai sumber, dan keempat tab dibaca dari
    // `M_TREATY_IN.JSONDATA` — dokumen yang tidak punya 41 kolom.
    //
    // ⛔ Kedua daftarnya TIDAK dihapus: keduanya mencatat bentuk tabel lama,
    // dan catatan itu masih dipakai `PEMETAAN-M-TREATY-IN2.md`. Yang dicabut
    // kesimpulannya, bukan datanya.
    expect(KOLOM_IN2_SEMUA.length).toBe(41)
    expect(KOLOM_IN2_TIDAK_DIPAKAI.length).toBe(10)

    // ⭐ Dan yang kini dijaga: Event Limits memperoleh TIGA kolom yang tabel
    // itu tidak pernah bisa berikan.
    for (const k of [
      'RSMD Limit',
      'Flood Limit (Jabodetabek)',
      'Flood Limit (Nationwide)',
    ]) {
      expect(KOLOM_EVENT_LIMITS).toContain(k)
    }
    expect(JENIS_EVENT_LIMITS).toHaveLength(KOLOM_EVENT_LIMITS.length)
  })
})

describe('angka: titik ribuan, koma desimal, dan batas desimalnya', () => {
  it('uang memakai 4 desimal dan titik ribuan', () => {
    expect(selAngka('uang', '5700000000')).toBe('5.700.000.000')
    expect(selAngka('uang', '1234.56789')).toBe('1.234,5679')
  })

  it('persen biasa dipangkas 2 desimal dan diberi tanda', () => {
    expect(selAngka('persen', '12.1000')).toBe('12,1%')
    expect(selAngka('persen', '33.333')).toBe('33,33%')
  })

  it('⛔ persen SHARE dipotong pada DELAPAN desimal — keputusan §13', () => {
    // Tiga share 33,333 berjumlah TEPAT 100; tiga share 33,33 tidak. Delapan
    // desimal menjaga sifat itu tanpa mengaku lebih teliti daripada
    // penyimpanannya, yang juga 8 (`NUMBER(38,8)`).
    expect(selAngka('persenShare', '33.333')).toBe('33,333%')

    // ⚠️ `99.999999999999900` TETAP tampil `100%` pada 8 desimal — dan itu
    // DIUKUR, bukan diandaikan. Pembulatannya setengah-ke-atas dan
    // limpahannya merambat, jadi selisihnya baru terlihat pada 13 desimal
    // ke atas. Salah satu dari dua alasan pemilik proses menolak 2 desimal
    // karena itu TIDAK terpenuhi oleh 8; alasan keduanya (sama dengan batas
    // penyimpanan) terpenuhi. Dicatat di KEPUTUSAN §13.1 dan ditagih balik.
    expect(selAngka('persenShare', '99.999999999999900')).toBe('100%')
    expect(selAngka('persenShare', '99.9999999')).toBe('99,9999999%')

    // ⛔ Tiga puluh desimal DIPOTONG pada delapan, tidak dibiarkan utuh.
    expect(selAngka('persenShare', '2.825601535925207120348922139444')).toBe('2,82560154%')
  })

  it('persen BUKAN-share tetap 2 desimal — ia tidak dijumlahkan menjadi 100', () => {
    expect(selAngka('persen', '33.333')).toBe('33,33%')
    expect(selAngka('persen', '12.1000')).toBe('12,1%')
  })

  it('⛔ pita tidak diberi tanda % kedua walau salah golong', () => {
    // `formatPersen` menempelkan `%` pada teks bukan-angka apa adanya,
    // menghasilkan `>=30% up to < 50%%`. Satu huruf salah di JENIS_* sudah
    // cukup untuk memicunya.
    expect(selAngka('persenShare', '>=30% up to < 50%')).toBe('>=30% up to < 50%')
    expect(selAngka('persen', '1/04/2023')).toBe('1/04/2023')
  })

  it('nol di ekor tetap dibuang, dan uang tetap 4 desimal', () => {
    expect(selAngka('persenShare', '12.10000000')).toBe('12,1%')
    expect(selAngka('persenShare', '50')).toBe('50%')
    expect(selAngka('uang', '1234.56789')).toBe('1.234,5679')
  })



  it('teks bukan-angka dikembalikan APA ADANYA, tidak berpura-pura nol', () => {
    // `CoInShare` adalah PITA, bukan bilangan — 702 barisnya berbentuk ini.
    expect(selAngka('teks', '>=30% up to < 50%')).toBe('>=30% up to < 50%')
    expect(selAngka('uang', '1/04/2023')).toBe('1/04/2023')
  })

  it('pengenal tidak kehilangan nol di depannya', () => {
    expect(selAngka('teks', '007')).toBe('007')
    expect(selAngka('teks', '1A')).toBe('1A')
  })

  it('kosong tetap kosong — bukan 0', () => {
    for (const j of ['uang', 'persen', 'persenShare', 'teks'] as const) {
      expect(selAngka(j, '')).toBe('')
    }
  })

  it('format.ts dipanggil, tidak ditulis ulang', () => {
    expect(FORM).toContain("from '../../../../inti/frontend/lib/format'")
    // ⛔ Nol pemformat kedua di dalam modul ini.
    expect(FORM).not.toMatch(/function formatNumber|function formatPersen/)
    expect(LABELS).not.toMatch(/function formatNumber|function formatPersen/)
  })
})

describe('tab Co-Ins Scale dan petunjuk kosong yang membedakan sebabnya', () => {
  it('Co-Ins Scale membaca tabel pendaratan kesembilan — DAN dua medan hidupnya', () => {
    // ⛔ RALAT DI RONDE YANG SAMA, 6 Oktober 2026. Uji ini sempat DIBALIK
    // untuk menuntut grid `TreatyIn.CoInScale` TIDAK dirender, atas dasar
    // `pyContainerVisibleWhen` = `1=2` di wadahnya. Bacaan itu SALAH: wadah
    // yang sama ber-`pyIsVisibilityOption` = `ALWAYS`, yang menimpa
    // syaratnya, dan gambar `18` (tangkapan layar Pega) memperlihatkan
    // gridnya TAMPIL. Bunyi aslinya — grid dirender — dipulihkan.
    //
    // Yang memang KURANG, dan kini ditambahkan: sel 277 `.MaxCoNonGroup`
    // dan sel 278 `.MaxCoGroup` di bawah grid.
    expect(FORM).toContain("tabTampil === 'Co-Ins Scale'")
    expect(FORM).toContain('warisan?.skalaKoasuransi')
    expect(FORM).toContain('<TabCoInsScale')
  })

  it('⭐ petunjuk kosong keempat tab layer DIBALIK — lubang 510 tertutup', () => {
    // ⛔ Uji ini DIBALIK, bukan dihapus.
    //
    // Bentuk sebelumnya MENUNTUT petunjuknya menyebut "1.340 dari 1.854" dan
    // "510", sebab kosong punya DUA arti dan layar tidak boleh mengaku tahu
    // yang mana. Sesudah `M_TREATY_IN2` dicabut, kosong punya SATU arti —
    // dan petunjuk lama menjelaskan sebab yang sudah tidak ada.
    // ⛔ 8 Oktober 2026 petunjuknya DICABUT seluruhnya — bukan hanya
    // dibersihkan dari `M_TREATY_IN2`. Lihat penjaga menyeluruh di akhir
    // berkas ini.
    expect(FORM_KONTRAK.petunjukLayer).toBe('')
  })

  it('⭐ petunjuk Event Limits DIBALIK — ketiga batas kini punya sumbernya', () => {
    // Bentuk sebelumnya MENUNTUT ketiganya disebut sebagai yang TIDAK punya
    // sumber. Dokumen punya keempatnya di `Limits[].Detail[]`, dan ketiganya
    // kini kolom di layar — jadi petunjuknya tidak perlu menyebutnya lagi.
    expect(FORM_KONTRAK.petunjukEventLimits).not.toContain('M_TREATY_IN2')
    expect(FORM_KONTRAK.petunjukEventLimits).not.toContain('belum punya sumber')
    for (const k of ['RSMD Limit', 'Flood Limit (Jabodetabek)', 'Flood Limit (Nationwide)']) {
      expect(KOLOM_EVENT_LIMITS).toContain(k)
    }
  })

  it('keempat tab layer memakai petunjuknya sendiri, bukan petunjuk umum', () => {
    for (const tab of ['Limits', 'Share', 'RNM Share']) {
      const i = FORM.indexOf(`tabTampil === '${tab}'`)
      expect(i).toBeGreaterThan(0)
      // ⛔ Dipotong sampai CABANG BERIKUTNYA, bukan 1.400 aksara tetap —
      // 6 Oktober 2026. Jendela tetap pecah begitu komentar cabangnya
      // memanjang (cabang Share kini merender dua komponen Prop/Non-Prop
      // dengan komentar bukti), padahal pemakaian petunjuknya tidak berubah.
      const j = FORM.indexOf('tabTampil ===', i + 1)
      const blok = FORM.slice(i, j === -1 ? undefined : j)
      expect(blok, tab).toContain('petunjukLayer')
    }
  })
})

describe('dua tab TEKS — Jalan B, keputusan §15', () => {
  it('keduanya membaca dokumen, bukan tabel pendaratan', () => {
    expect(FORM).toContain("tabTampil === 'Exclusions'")
    expect(FORM).toContain("tabTampil === 'Special Conditions'")
    expect(FORM).toContain('warisan?.pengecualian')
    expect(FORM).toContain('warisan?.syaratKhusus')
    // ⛔ Nol tabel pendaratan dibuat untuk keduanya — itulah Jalan B.
    for (const salah of ['M_TREATYIN_TEKS', 'M_TREATYIN_EXCLUSIONS', 'M_TREATYIN_SPECIALCONDITIONS']) {
      expect(FORM).not.toContain(salah)
      expect(LABELS).not.toContain(salah)
    }
  })

  it('teksnya DAPAT DIGULIR, bukan satu baris terpotong', () => {
    expect(FORM).toContain('trin__teks')
    // Isinya mencapai 23.453 aksara; sel yang memotongnya menyembunyikan
    // teks tanpa memberi tahu ada yang disembunyikan.
    expect(CSS).toMatch(/\.trin__teks \{[\s\S]*?overflow-y: auto/)
    expect(CSS).toMatch(/\.trin__teks \{[\s\S]*?white-space: pre-wrap/)
  })

  it('⛔ kosong memakai Kosong (belum ada DATA), bukan .trin__belum', () => {
    const i = FORM.indexOf('function TabTeksPanjang')
    expect(i).toBeGreaterThan(0)
    // Seluruh badan fungsi — bukan jendela beraksara tetap, yang bergeser
    // setiap kali fungsi ini tumbuh (`TreatyInCopyConditions`, 7 Oktober 2026).
    const blok = FORM.slice(i, FORM.indexOf('\n}\n', i))
    expect(blok).toContain('<Kosong')
    expect(blok).not.toContain('trin__belum')
  })

  it('peringatan ejaan lain menyatakan isinya BERBEDA, bukan salinan', () => {
    expect(FORM_KONTRAK.ejaanLainBerisi).toContain('BERBEDA')
    expect(FORM_KONTRAK.ejaanLainBerisi).toContain('bukan salinan')
    expect(FORM).toContain('ejaanLainBerisi')
  })

  // ⛔ DIBALIK 8 Oktober 2026 — ejaan kunci dokumen adalah urusan
  // pengembang; ia tidak lagi tampil di layar.
  it('⛔ petunjuk kosong kedua tab teks dicabut', () => {
    expect(FORM_KONTRAK.petunjukTeksPengecualian).toBe('')
    expect(FORM_KONTRAK.petunjukTeksSyarat).toBe('')
  })

  it('⛔ Value Difference TETAP belum dibangun — ia objek, bukan teks', () => {
    expect(FORM).not.toContain("tabTampil === 'Value Difference'")
  })
})

describe('empat grid yang sampai 4 Oktober 2026 LOLOS dari pemformat', () => {
  // ⛔ Lima grid layer memakai `selAngka`, empat grid pendaratan tidak —
  // mereka menyusun barisnya sendiri dan melewati `barisLayer`. Akibatnya
  // `Amount` 11080000000 terbaca apa adanya di layar. Penjaga di bawah
  // memakai nilai yang BENAR-BENAR ada di Oracle, bukan contoh karangan.

  it('⛔ keempat grid kini melewati barisAngka', () => {
    // ⛔ Installment KELUAR 7 Oktober 2026 — bukan lagi grid umum: tabnya
    // sendiri (`TabAngsuran.tsx`) memformat Amount/Total lewat `selAngka`
    // (`['uang', 2]`, % lewat `persenShare`), seperti gambar 38.
    const ANG = readFileSync(join(AKAR, 'components', 'TabAngsuran.tsx'), 'utf8')
    expect(ANG).toContain("selAngka(['uang', 2], v)")
    expect(ANG).toContain("selAngka(['persenShare', 2], r.InstallmentPct)")
    // Rate of Exchange tidak punya JENIS_*; kursnya diformat langsung.
    expect(FORM).toContain("selAngka(['uang', 2], b.nilaiKeIDR)")
  })

  // ⛔ EGNPI KELUAR dari daftar di atas 6 Oktober 2026 — bukan karena
  // pemformatannya dicabut, melainkan karena ia bukan lagi grid umum.
  // Tabnya sendiri (`TabEgnpi.tsx`) memformat lewat `formatLimit` +
  // `DESIMAL_EGNPI`, jadi yang dijaga tetap sama: nol angka mentah di layar.
  //
  // ⚠️ Dua desimal yang BERBEDA di baris yang sama adalah intinya —
  // `Amount` 2, `Amount in IDR` 0 (gambar 29). Menyeragamkannya menghapus
  // bukti yang membuat angka ini terbaca seperti di Pega.
  // ⛔ Maximum Retention keluar dari daftar di atas sebab yang sama
  // dengan EGNPI — bukan pemformatan yang dicabut, melainkan grid umumnya.
  it('⛔ Maximum Retention memformat lewat tabnya sendiri', () => {
    const RT = readFileSync(join(__dirname, 'components', 'TabRetensi.tsx'), 'utf8')
    expect(RT).toContain('formatLimit')
    expect(RT).toContain('DESIMAL_RETENSI.jumlah')
    expect(RT).toContain('DESIMAL_RETENSI.nilaiTotal')
    // ⚠️ Grid 0 desimal, panel total 2 — gambar 26. Nilai yang sama,
    // dua presisi, dua tempat.
    expect(DESIMAL_RETENSI.jumlah).not.toBe(DESIMAL_RETENSI.nilaiTotal)
  })

  it('⛔ EGNPI memformat lewat tabnya sendiri', () => {
    const EG = readFileSync(join(__dirname, 'components', 'TabEgnpi.tsx'), 'utf8')
    expect(EG).toContain('formatLimit')
    expect(EG).toContain('DESIMAL_EGNPI.jumlah')
    expect(EG).toContain('DESIMAL_EGNPI.jumlahIDR')
    expect(DESIMAL_EGNPI.jumlah).not.toBe(DESIMAL_EGNPI.jumlahIDR)
  })

  it('tiap JENIS_* sepanjang KOLOM_* pasangannya', () => {
    expect(JENIS_EGNPI.length).toBe(KOLOM_EGNPI.length)
    expect(JENIS_RETENSI.length).toBe(KOLOM_RETENSI.length)
    expect(JENIS_ANGSURAN.length).toBe(KOLOM_ANGSURAN.length)
  })

  it('nilai Oracle NYATA tampil bertitik ribuan', () => {
    // POOLDATA.T_TREATY_EGNPI.AMOUNT / .AMOUNTIDR
    expect(selAngka('uang', '11080000000')).toBe('11.080.000.000')
    expect(selAngka('uang', '11080000000.0')).toBe('11.080.000.000')
    // POOLDATA.T_TREATY_RETENTION.AMOUNT
    expect(selAngka('uang', '5000000000')).toBe('5.000.000.000')
    // POOLDATA.T_TREATY_INSTALLMENT_ITEM.AMOUNT
    expect(selAngka('uang', '9267363.876000000000')).toBe('9.267.363,876')
  })

  it('⛔ Proportion EGNPI adalah persen SHARE — ia berjumlah TEPAT 100', () => {
    // Diukur 4 Oktober 2026: 6 dari 6 kontrak bersampel menjumlahkan
    // `PROPORTION` menjadi tepat 100. Nilainya menyimpan 20 desimal, jadi
    // batas 8 benar-benar menggigit di sini.
    // ⚠️ `golongan()` — sejak §24 larik `JENIS_*` dapat membawa PASANGAN
    // [golongan, desimal]. Yang diuji di sini golongannya; desimalnya diuji
    // terpisah di `desain-pega.test.ts`.
    expect(golongan(JENIS_EGNPI[KOLOM_EGNPI.indexOf('Proportion %')]!)).toBe('persenShare')
    expect(selAngka('persenShare', '41.63535247256876597000')).toBe('41,63535247%')
  })

  it('Pct Installment juga persen share, dan 25.00 tetap terbaca 25%', () => {
    expect(golongan(JENIS_ANGSURAN[KOLOM_ANGSURAN.indexOf('Pct')]!)).toBe('persenShare')
    expect(selAngka('persenShare', '25.00')).toBe('25%')
  })

  it('contoh pemilik proses, 4 Oktober 2026', () => {
    // "1500000000 seharusnya 1.500.000.000" bila tanpa koma — nol di ekor
    // TIDAK dikarang, sesuai kalimat "kalau tanpa koma langsung masuk saja".
    expect(selAngka('uang', '1500000000')).toBe('1.500.000.000')
    // Uang berdesimal dipangkas pada 4.
    expect(selAngka('uang', '1500000000.43133134')).toBe('1.500.000.000,4313')
    // Persen biasa pada 2.
    expect(selAngka('persen', '15.21')).toBe('15,21%')
    // ⚠️ Delapan desimal hanya untuk persen SHARE, dan hanya di sana —
    // penyimpanannya `NUMBER(38,8)`.
    expect(selAngka('persenShare', '1500000000.43133134')).toBe('1.500.000.000,43133134%')
  })

  it('⛔ hitungan hari BUKAN uang — SubDays tidak digolongkan angka', () => {
    // `T_TREATY_ACCUMULATION.SUBDAYS` bernilai 30. Ia cacah hari, bukan
    // nominal; memberinya aturan uang berarti mengarang satuan yang pemilik
    // proses tidak sebut. Grid Accumulation karena itu TIDAK punya JENIS_*.
    expect(FORM).not.toContain('JENIS_AKUMULASI')
  })
})

describe('panel Attachment + History — nol tabel baru', () => {
  it('Attachment dibaca dari tabel WARISAN, bukan tabel pendaratan baru', () => {
    expect(FORM).toContain('PanelLampiran')
    expect(FORM).toContain('warisan?.kategoriLampiran')
    expect(FORM).toContain('warisan?.lampiran')
    // ⛔ Nol tabel pendaratan dibuat untuk lampiran — ia sudah relasional.
    for (const salah of ['M_TREATYIN_LAMPIRAN', 'M_TREATYIN_ATTACHMENT']) {
      expect(FORM).not.toContain(salah)
      expect(LABELS).not.toContain(salah)
    }
  })

  it('keempat kolomnya disalin dari WorkAttachments.xml apa adanya', () => {
    expect([...KOLOM_LAMPIRAN]).toEqual(['Category', 'Count', 'Upload file', 'View File'])
    expect([...KOLOM_HISTORY]).toEqual(['Date', 'PIC', 'Approval', 'Comment'])
  })

  it('spanduk biru disalin apa adanya — ia ATURAN, bukan hiasan', () => {
    expect(LAMPIRAN.spanduk).toBe('Recommended safe substitute should be . or _')
    expect(FORM).toContain('trin__spanduk')
    expect(CSS).toMatch(/\.trin__spanduk \{/)
  })

  // ⭐ RALAT 6 Oktober 2026 — panel menampilkan NAMA, kesebelasnya.
  //
  // Bentuk sebelumnya merender empat baris sebagai KODE-nya (`00003 nama
  // kategori belum dipastikan`), sebab daftarnya datang dari katalog basis
  // data yang hanya memuat tujuh kategori yang pernah dipakai. Tangkapan
  // layar pemilik proses memperlihatkan kesebelas NAMA beserta `Count 0`.
  //
  // ⛔ Yang TIDAK berubah: pasangan kode↔nama tetap tidak ditebak. Penjaga
  // itu pindah ke `services.SusunKategoriLampiran`
  // (`TestEmpatNamaTanpaKodeTIDAKDitebakKodenya`) — tempat pasangannya
  // sungguh dibentuk, bukan layar yang hanya menampilkannya.
  it('⭐ panel menampilkan NAMA kategori, bukan kodenya', () => {
    const i = FORM.indexOf('function PanelLampiran')
    expect(i).toBeGreaterThan(0)
    const blok = FORM.slice(i, FORM.indexOf('function berkasKategori', i))
    expect(blok).toContain('<td>{k.nama}</td>')
    expect(blok).not.toContain('{k.kode} <span')
  })

  it('⛔ berkas disaring menurut kode bila diketahui, NAMA bila tidak', () => {
    // Keempat kategori tanpa kode akan selalu berbunyi `No items` kalau
    // penyaringnya hanya kode — dan kosong yang salah terbaca persis seperti
    // kosong yang benar.
    expect(FORM).toContain('b.kodeKategori === kategori.kode')
    expect(FORM).toContain('b.namaKategori === kategori.nama')
  })

  it('kesebelas nama kategori sama dengan tangkapan layar pemilik proses', () => {
    expect([...NAMA_KATEGORI_LAMPIRAN_PROP]).toEqual([
      'Analysed Email',
      'Approval Email',
      'Assessment Inward Treaty Form / Format Analisa Treaty',
      'Binding, signed share Email',
      'Claim Data',
      'Info Pack',
      'Letter of Acknowledgment / LOA',
      'Offer Email',
      'Others',
      'Pega Proportional Calculation /Perhitungan Pega Proportional',
      'Summary Treaty Leader',
    ])
  })

  it('⛔ Count TIDAK diformat — ia cacah butir, bukan uang', () => {
    const i = FORM.indexOf('function PanelLampiran')
    const blok = FORM.slice(i, FORM.indexOf('function berkasKategori', i))
    expect(blok).toContain('String(k.cacah)')
    // Nol pemformat angka menyentuhnya — cacah butir tidak diberi pemisah ribuan.
    expect(blok).not.toMatch(/selAngka\([^)]*cacah/)
  })

  it('History memakai tabel pendaratan yang SUDAH ada, nol pembacaan baru', () => {
    expect(FORM).toContain('PanelHistory')
    expect(FORM).toContain('warisan?.catatan')
    expect(LAMPIRAN.tanpaRiwayat).toBe('No items')
  })

  it('kosong memakai Kosong (belum ada DATA), bukan .trin__belum', () => {
    for (const nama of ['function PanelLampiran', 'function PanelHistory']) {
      const i = FORM.indexOf(nama)
      expect(i).toBeGreaterThan(0)
      // Panel Attachment kini memuat jalur unggah — dipotong sampai ujung fungsinya.
      const ujung = nama === 'function PanelLampiran' ? FORM.indexOf('function berkasKategori', i) : i + 4200
      const blok = FORM.slice(i, ujung)
      expect(blok).toContain('<Kosong')
      expect(blok).not.toContain('trin__belum')
    }
  })

  it('⭐ Save dan Actions HIDUP menurut syarat tampil ekspor, Close selalu', () => {
    // 7 Oktober 2026 — keputusan pemilik proses: jalur tulis ke tabel
    // masing-masing. Syarat tampil disalin dari `TreatyInActionButtons`.
    const i = HALAMAN.indexOf('className="trin__aksi trin__aksi--kaki"')
    expect(i).toBeGreaterThan(0)
    const blok = HALAMAN.slice(i, i + 1200)
    expect(blok).toMatch(/\{saveTampil && \(\s*<button[^>]*onClick=\{tekanSave\}/)
    expect(blok).toContain('onClick={onKembali}')
    expect(blok).toMatch(/\{actionsTampil && \(/)
    expect(HALAMAN).toContain("const saveTampil = bisaUbah && statusKini !== 'Resolve Complete'")
    expect(HALAMAN).toContain("bolehActions(warisan?.posisi ?? '', workbasket, statusKini)")
  })

  it('⭐ properti penampung yang TERSIMPAN (migrasi 448) disemai sesudah kosongkan', () => {
    // Tanpa ini RNMShareP/OptionLimit/InstallmentNo yang di-Save kembali
    // kosong begitu kontrak dibaca ulang — tabnya menyemai bawaan ''.
    const i = HALAMAN.indexOf('Object.entries(k.penampung ?? {})')
    expect(i).toBeGreaterThan(0)
    expect(HALAMAN.lastIndexOf('penampung.kosongkan()', i)).toBeGreaterThan(HALAMAN.lastIndexOf('setWarisan(k)', i))
    expect(HALAMAN.slice(i, i + 120)).toContain('penampung.ubah(kunci, () => v)')
  })

  it('deret Save · Close · Actions berjarak dari kartu History di bawahnya', () => {
    // Tangkapan layar pemakai 7 Oktober 2026: History menempel ke tombol.
    // Jarak bawahnya sama dengan jarak antar kartu (`.treatyin .panel`).
    const i = HALAMAN.lastIndexOf('trin__aksi')
    expect(HALAMAN.slice(i - 40, i + 40)).toContain('className="trin__aksi trin__aksi--kaki"')
    expect(HALAMAN.indexOf('<PanelHistory', i)).toBeGreaterThan(i)
    expect(CSS).toMatch(/\.treatyin \.trin__aksi--kaki \{\s*margin-bottom: 18px;/)
    expect(CSS).toMatch(/\.treatyin \.panel \{\s*margin-bottom: 18px;/)
  })
})

describe('pemilih "Choose …" adalah DROPDOWN sejak 5 Oktober 2026', () => {
  // ⛔ Bentuk sebelumnya tombol yang membuka modal. Pemilik proses
  // menggantinya menjadi daftar pilihan: satu ketukan, bukan dua.

  it('nol tombol "Choose" tersisa — keduanya daftar pilihan', () => {
    expect(FORM).not.toContain('setPemilih')
    expect(FORM).not.toContain('PemilihWarisan')
    expect(FORM).toContain('DropdownWarisan')
  })

  it('keduanya memanggil katalognya masing-masing', () => {
    expect(FORM).toContain('ambil={ambilDaftarCedant}')
    expect(FORM).toContain('ambil={ambilDaftarAsalBisnis}')
  })

  it('⛔ nama KEMBAR tetap dapat dibedakan — dropdown memperburuk soalnya', () => {
    // Modal masih memperlihatkan kolom pengenal; dropdown tidak. Dua baris
    // bernama sama di dalam daftar pilihan karena itu MUSTAHIL dibedakan
    // kecuali pengenalnya dibubuhkan ke labelnya. Terukur: 36 nama cedant
    // dipakai lebih dari satu pengenal.
    expect(FORM).toContain('b.kembar ?')
    expect(FORM).toContain('${b.nama} — ${b.id}')
  })

  it('⛔ tiap baris berkunci PENGENAL, dan barisnya diserahkan UTUH', () => {
    // Nama sebagai kunci membuat dua baris kembar saling menimpa, dan yang
    // memilih tidak akan pernah tahu baris mana yang terpilih.
    expect(FORM).toContain('key={b.id}')
    expect(FORM).toContain('onPilih(b)')
  })

  it('kotaknya DAPAT DIKETIK dan daftarnya menyaring', () => {
    // 131 cedant dan 96 asal bisnis — menggulir lebih lambat daripada
    // mengetik tiga huruf.
    expect(FORM).toContain("role=\"combobox\"")
    expect(FORM).toContain('setKetik')
    // Saringan membaca NAMA dan PENGENAL: yang hafal kodenya mengetik kode.
    //
    // ⚠️ RALAT 5 Oktober 2026: saringannya PINDAH ke `saring.ts`, jadi
    // pernyataan ini menanyakannya di sana. Yang dijaga tidak berubah —
    // bahwa saringan membaca nama DAN pengenal.
    expect(SARING).toContain('nama.toLowerCase()')
    expect(SARING).toContain('id.toLowerCase()')

  })

  it('⛔ pilihan dibaca sebelum kotaknya kehilangan fokus', () => {
    // `onBlur` menyala mendahului `onClick`; memakai `onClick` membuat
    // pilihan tidak pernah sampai.
    expect(FORM).toContain('onMouseDown')
  })

  it('⛔ keterangan kaki tentang pemilih DIBUANG dari layar', () => {
    // Keputusan pemilik proses 5 Oktober 2026: prosa penjelas tidak duduk di
    // jalur baca pemakai.
    expect(FORM).not.toContain('catatanPilihLuar')
  })
})

describe('Retro — keadaan KEEMPAT: jarang, bukan belum-ada-kode', () => {
  it('tab Retro memakai teks JARANG, bukan belumDibangun', () => {
    expect(FORM).toContain("tabTampil === 'Retro' ? FORM_KONTRAK.jarangDipakai")
    expect(FORM_KONTRAK.jarangDipakai).toContain('JARANG')
    expect(FORM_KONTRAK.jarangDipakai).not.toContain('belum dibangun.')
  })

  // ⛔ DIBALIK 8 Oktober 2026 — petunjuknya DICABUT dari layar.
  //
  // ⚠️ Yang dijaga uji lama tetap BENAR dan tetap perlu diketahui: Retro
  // HIDUP di sistem lama (`pyRuleAvailable = Yes`, nol penjaga `1=2`), hanya
  // jarang dipakai. Keterangan itu kini hidup di komentar kode dan dokumen
  // modul, bukan di layar pemakai.
  it('⛔ petunjuk Retro dicabut dari layar', () => {
    expect(FORM_KONTRAK.jarangDipakaiPetunjuk).toBe('')
    expect(FORM_KONTRAK.belumDibangunPetunjuk).toBe('')
    // Pesan tab-nya sendiri TETAP — ia kalimat polos, bukan catatan internal.
    expect(FORM_KONTRAK.belumDibangun).toBe('Tab ini belum dibangun.')
  })

  it('tab lain TETAP memakai belumDibangun — keempat keadaan tidak tertukar', () => {
    expect(FORM).toContain('FORM_KONTRAK.belumDibangun}')
    expect(FORM_KONTRAK.belumDibangun).toBe('Tab ini belum dibangun.')
  })
})

// ===========================================================================
// Galat `Cannot read properties of null (reading 'length')` — 6 Oktober 2026
// ===========================================================================
//
// Dilaporkan pemakai sesudah menekan `Edit` lalu `View`: SELURUH halaman
// berhenti, bukan satu grid. Sebabnya bukan di layar melainkan di kawatnya —
// backend mengirim `"kurs": null` untuk kontrak yang dokumennya tidak punya
// `CurrencyList`, sementara `api.ts` menyatakan `kurs: BarisKursWarisan[]`.
//
// ⛔ TypeScript TIDAK DAPAT menangkap ini, dan itu pokoknya: tiap tipe di
// `api.ts` adalah KLAIM tentang kawat yang tidak ada yang membuktikan. Yang
// menegakkan janjinya ada di Go (`services/warisan_nol.go` beserta ujinya);
// yang dijaga di sini adalah lapis keduanya.
describe('larik dari API tidak boleh menghentikan halaman', () => {
  it('⛔ `kurs` dibaca lewat `?? []` — satu-satunya larik yang lewat useState', () => {
    expect(FORM).toContain('setKurs(k.kurs ?? [])')
    // Dan ia TIDAK boleh kembali ke bentuk telanjang.
    expect(FORM).not.toMatch(/setKurs\(k\.kurs\)/)
  })

  it('⛔ setiap medan larik `warisan` dibaca lewat `?? []`', () => {
    // Tiap `warisan?.<medan>` yang diteruskan sebagai larik harus berakhir
    // `?? []`. Yang lolos dari pola ini akan jatuh persis seperti `kurs`.
    // ⛔ Tangkap juga `??`-nya, jangan pakai lookahead negatif: `\w+` akan
    // mundur satu aksara untuk memuaskannya, dan daftar yang keluar berisi
    // nama yang terpotong — `polisProduks`, `laye`.
    const telanjang = [...FORM.matchAll(/warisan\?\.(\w+)(\s*\?\?)?/g)]
      .filter((m) => m[2] === undefined)
      .map((m) => m[1])
      .filter((n): n is string => n !== undefined)
      // Medan BUKAN larik — dibaca apa adanya, dan `?? []` di atasnya keliru.
      // `limitsAkar` OBJEK, bukan larik — `TabLimitsNonProp` memakai bentuk
      // kosongnya sendiri bila tak ada, dan larik di dalamnya dijamin
      // services (`nolkanLarik`).
      // `shareNP` OBJEK pula — `TabShareNonProp` memakai `SHARE_NP_KOSONG`
      // bila tak ada, dan lariknya dijamin services (`lengkapiShare`).
      .filter((n) => !['tahunTreaty', 'pengecualian', 'syaratKhusus', 'limitsAkar', 'shareNP'].includes(n))
    expect(telanjang).toEqual([])
  })
})

// ===========================================================================
// Grid Rate of Exchange — bentuk angka dan tanggalnya, 6 Oktober 2026
// ===========================================================================
//
// Cacat yang melahirkannya: mode Edit merender nilai MENTAH (`10584.39`)
// sementara mode View merender `10.584,39`. Satu layar, dua bentuk, dan yang
// menyuntingnya mengira angkanya memang berbeda.
// ⭐ 7 Oktober 2026 — sel grid Rate of Exchange mengikuti
// `Section/TreatyInNONProportional.xml` (@439485–@466983). Dahulu keempatnya
// kotak teks bebas di mode Edit.
describe('grid Rate of Exchange — kontrol sel sesuai ekspor', () => {
  const HAL = readFileSync(join(__dirname, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
  const i = HAL.indexOf('{bisaUbah &&\n                kurs.map')
  const blokUbah = HAL.slice(i, HAL.indexOf('{!bisaUbah && kurs.map', i))

  it('Currency `.CurrencyID` = pxDropdown `BrowseCurrency_RD` — bukan teks bebas', () => {
    expect(blokUbah).toMatch(/<DropdownDaftar[\s\S]*?nilai=\{b\.mataUang\}[\s\S]*?pilihan=\{mataUangKurs\}/)
    // `SetCurrNameMasterTreaty_Act`: memilih pengenal mengisi namanya juga.
    expect(blokUbah).toContain('ubahKurs(i, { mataUang: nama, mataUangID: id })')
    // Opsinya rute yang sama dengan dropdown mata uang tab Limits.
    expect(HAL).toContain('if (!dibuang) setMataUangKurs(o.mataUang)')
    // ⛔ Kotak teks yang mengulang keempat kunci sudah tidak ada.
    expect(HAL).not.toContain("(['mataUang', 'nilaiKeIDR', 'berlakuDari', 'berlakuSampai'] as const)")
  })

  it('Valid From / Valid Until = pxDateTime — kotak tanggal, dan mode baca memformat tanggalnya', () => {
    expect((blokUbah.match(/<KotakTanggalKetik/g) ?? []).length).toBe(2)
    expect(HAL).toContain('<td>{formatDate(b.berlakuDari)}</td>')
    expect(HAL).toContain('<td>{formatDate(b.berlakuSampai)}</td>')
  })

  it('Add = `TreatyInAddCurrency` — baris baru ber-CurrencyID kosong', () => {
    // ⛔ Baris baru kini membawa DUA pasang medan tanggal: bentuk TAMPIL
    // (`berlakuDari`) dan bentuk TERSIMPAN (`…Asli`). Yang kedua lahir
    // 7 Oktober 2026 sebab kotak tanggal tidak dapat membaca `dd/mm/yy`.
    expect(HAL).toContain(
      "{ mataUang: '', mataUangID: '', nilaiKeIDR: '', berlakuDari: '', berlakuSampai: '', berlakuDariAsli: '', berlakuSampaiAsli: '' }",
    )
  })
})

describe('grid Rate of Exchange — satu bentuk angka di kedua mode', () => {
  it('⛔ sel isian pun diformat, bukan hanya sel baca', () => {
    // Sel isian: diformat kecuali sedang diketik.
    expect(FORM).toContain("selDiketik === `${String(i)}:nilaiKeIDR` ? b.nilaiKeIDR : selAngka(['uang', 2], b.nilaiKeIDR)")
    // Mode baca memakai pemformat yang sama; keduanya harus sama.
    expect(FORM).toContain("<td>{selAngka(['uang', 2], b.nilaiKeIDR)}</td>")
  })

  it('⭐ pemformat berhenti saat selnya sedang diketik', () => {
    // Memformat di tiap ketukan membuat koma desimal mustahil diketik:
    // `10,` berubah menjadi `10` sebelum angka berikutnya sempat masuk.
    expect(FORM).toContain('selDiketik')
    expect(FORM).toContain('onFocus')
    expect(FORM).toContain('onBlur')
  })

  it('pemisah ribuan TITIK, desimal KOMA — terbukti atas nilai nyata', () => {
    expect(selAngka(['uang', 2], '10584.39')).toBe('10.584,39')
    expect(selAngka(['uang', 2], '2294.18')).toBe('2.294,18')
    // Nol di ekor DIPERTAHANKAN — gambar 01 berbunyi `1,00`.
    expect(selAngka(['uang', 2], '1')).toBe('1,00')
  })
})

// ⛔ Ketikan tidak boleh terhapus oleh peristiwa fokus susulan.
//
// Keluhan pemilik proses 6 Oktober 2026: diketik `zurich`, daftar kembali
// memperlihatkan seluruh 131 cedant sementara kotaknya masih berbunyi
// `zurich`. Saringannya BENAR — `saring.test.ts` membuktikannya atas nama
// yang persis sama — yang keliru pengosongan `ketik` pada SETIAP fokus.
describe('pemilih ketik-pilih menjaga ketikannya', () => {
  const DD = readFileSync(join(__dirname, 'components', 'DropdownWarisan.tsx'), 'utf8')

  it('⛔ fokus hanya mengosongkan ketikan bila daftarnya TERTUTUP', () => {
    const i = DD.indexOf('onFocus={() => {')
    expect(i).toBeGreaterThan(0)
    const blok = DD.slice(i, i + 1600)
    expect(blok).toContain('if (!buka) {')
    // Pengosongan tanpa syarat adalah bentuk yang melahirkan keluhannya.
    expect(blok).not.toMatch(/onFocus=\{\(\) => \{\s*setKetik\(''\)/)
  })

  it('daftar yang dirender adalah hasil SARINGAN, bukan daftar penuh', () => {
    expect(DD).toContain('const cocok = saringTerdekat(daftar, q)')
    expect(DD).toContain('{cocok.map((b, i) => (')
    expect(DD).not.toContain('{daftar.map(')
  })
})


// ⛔ KEPALA KOLOM TETAP TERLIHAT SAAT TABEL DIGULIR.
//
// Permintaan pemilik proses 7 Oktober 2026: *"label per table kalau ada
// scroll nya labelnya itu netap tidak ikut tenggelam"*.
//
// ⚠️ Ketiga bagiannya DIJAGA TERPISAH, sebab tiap satunya punya cara
// gagal sendiri yang TIDAK terlihat di layar sampai seseorang menggulir.
describe('kepala kolom tidak ikut tenggelam', () => {
  const ATURAN = CSS.slice(CSS.indexOf('.treatyin .table-wrap .trin__tabel thead th'))

  it('⛔ sticky DIBATASI `.table-wrap` — di luar itu ia menempel di belakang topbar', () => {
    // Tanpa penggulir sendiri, penggulir terdekat adalah HALAMAN, dan kepala
    // kolom berhenti di tepi atas layar di belakang `.topbar` (z-index 20).
    expect(CSS).toMatch(/\.treatyin \.table-wrap \.trin__tabel thead th \{/)
    expect(CSS).not.toMatch(/\.treatyin \.trin__tabel thead th \{\s*position: sticky/)
  })

  it('⚠️ berlatar TIDAK tembus pandang, atau baris yang lewat menimpa labelnya', () => {
    expect(ATURAN.slice(0, 260)).toMatch(/background: var\(--bg\)/)
  })

  it('⚠️ garis bawahnya `box-shadow`, bukan `border-bottom` yang ikut tergulir pergi', () => {
    expect(ATURAN.slice(0, 260)).toMatch(/box-shadow: inset 0 -1px 0 var\(--border\)/)
  })

  it('⛔ tabel rincian BERSARANG dikecualikan — kepalanya milik satu baris', () => {
    // `TabAngsuran`/`TabShareProp` menaruh tabel di dalam sel baris induk;
    // penggulir terdekatnya `.table-wrap` MILIK INDUK, jadi kepalanya akan
    // melesat ke puncak dan memberi label pada baris yang bukan miliknya.
    expect(CSS).toMatch(/\.treatyin \.trin__rincian \.trin__tabel thead th \{\s*position: static/)
    for (const f of ['components/TabAngsuran.tsx', 'components/TabShareProp.tsx']) {
      expect(readFileSync(join(__dirname, f), 'utf8')).toContain('trin__rincian')
    }
  })
})
