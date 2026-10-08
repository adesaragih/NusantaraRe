// Label layar modul Treaty In.
//
// ⛔ SETIAP teks di berkas ini DISALIN dari ekspor Pega 2026-09, dan tiap
// barisnya membawa jejaknya: nama rule `pyCaption`/`pyButtonLabel` dan posisi
// bitanya di dalam berkas Section. Posisi itu ada supaya pembaca berikutnya
// dapat memeriksanya sendiri, bukan mempercayai berkas ini.
//
// Cara memeriksa satu baris:
//
//   python - <<'PY'
//   import io,re
//   s=io.open(r'D:\XML_NURE\Treaty In\Section\InputTreatyInOffer.xml',encoding='utf8').read()
//   s=re.sub(r"<pyIncludedRuleXML>.*?</pyIncludedRuleXML>","",s,flags=re.S)
//   print(s[6518125:6518200])
//   PY
//
// ⛔ Salinan rule anak (`<pyIncludedRuleXML>`) DIBUANG lebih dulu; tanpa itu
// medan milik seksi anak terbaca sebagai milik induk. Posisi bita di bawah
// dihitung SESUDAH pembuangan — 8.194.691 → 7.295.384 bita.

import type { HimpunanAcuan } from './api'

export const MENU_TREATYIN = {
  acuan: 'Treaty In — Tabel Acuan',
  daftar: 'Treaty In — Daftar Kontrak',
} as const

export const LABEL_HIMPUNAN: Record<HimpunanAcuan, string> = {
  'jenis-potongan': 'Jenis Potongan',
  'kelas-bisnis': 'Kelas Bisnis',
  'kelompok-treaty': 'Kelompok Treaty',
  bahaya: 'Bahaya',
  'jenis-reasuransi': 'Jenis Reasuransi',
}

export const ACUAN_TREATYIN = {
  kolomKode: 'Kode',
  kolomNama: 'Nama',
  kolomAktif: 'Aktif',
  kolomInduk: 'Induk',
  // Dua kalimat, dua medan `Kosong` yang berbeda: `pesan` menyatakan KEADAAN
  // (tabelnya memang belum berisi), `petunjuk` menyatakan SIAPA yang akan
  // mengisinya. Digabung jadi satu paragraf, yang kedua terbaca sebagai alasan
  // kosongnya - padahal ia jadwal, bukan sebab.
  kosong: 'Tabel acuan ini belum berisi.',
  kosongPetunjuk: '',
} as const

// ===========================================================================
// LAYAR A — DAFTAR KONTRAK
// Sumber: `Section/InputTreatyInOffer.xml`
// ===========================================================================

/**
 * Kesembilan kolom, **persis urutan layar lama**.
 *
 * ⛔ Urutannya bagian dari salinan, bukan selera. Kolom yang ditukar membuat
 * orang yang hafal layar lama membaca angka di kolom yang salah.
 */
export const KOLOM_DAFTAR = [
  { kunci: 'id', label: 'ID' },
  { kunci: 'namaKontrak', label: 'Contract Name' }, // pyCaption @7257527
  { kunci: 'sifatProporsi', label: 'Reinsurance Type' }, // pyCaption @7236037
  { kunci: 'idAsalBisnis', label: 'Source of Business' }, // pyCaption @6540552
  { kunci: 'idCedant', label: 'Ceding' }, // pyCaption @6503393
  { kunci: 'tanggalMulai', label: 'Commencement' }, // pyCaption @6540017
  { kunci: 'tanggalBerakhir', label: 'Termination' }, // pyCaption @6523491
  { kunci: 'posisiKe', label: 'Position To' }, // pyCaption @7232843
  { kunci: 'keadaanSiklusHidup', label: 'Status Accept' }, // pyCaption @7249608
] as const

export const DAFTAR_KONTRAK = {
  judul: 'Treaty In — Daftar Kontrak',
  tambah: 'Add', // pyButtonLabel @1602683
  saring: 'Show/Hide filter', // pyButtonLabel @7276317
  aksi: 'Aksi',
  // Tombol aksi, dan keempatnya ada di ekspor sebagai `pyButtonLabel`.
  edit: 'Edit', // @7234432
  lihat: 'View', // @7258590
  salin: 'Copy', // @7283254
  revisi: 'Revision', // @7247465
  /**
   * Petunjuk arahkan-kursor tombol Revision. Event-nya `doubleclick` di
   * ekspor (`pyActionSets` cell 994; Copy/Edit/View `click`) dan ia langsung
   * MENYIMPAN — laporan pemakai 8 Oktober 2026: satu klik dikira tombolnya
   * tidak berjalan. Event tetap seperti ekspor; petunjuk ini tambahan.
   */
  petunjukRevisi: 'Klik dua kali untuk membuat revisi',
  /**
   * Keadaan yang membedakan susunan tombolnya.
   *
   * ⛔ Baris berstatus ini mendapat `View`+`Copy`+`Revision`; baris lain
   * mendapat `Edit`+`View`. Perbedaan itu KEADAAN SIKLUS HIDUP, bukan hiasan,
   * dan meratakannya menghapus satu-satunya petunjuk di layar bahwa sebuah
   * kontrak sudah tidak dapat disunting.
   */
  keadaanTerkunci: 'Resolve Complete',
  kosong: 'Tidak ada kontrak pada halaman ini.',
  // ⚠️ RALAT 3 Oktober 2026. Kalimat ini pernah berbunyi "Pemindahan kepala
  // kontrak warisan dari sistem lama adalah tiket 59; sampai ia jalan,
  // tabelnya memang kosong." Itu benar ketika layar membaca `KONTRAK`
  // (model baru, nol baris) — dan TIDAK LAGI BENAR sejak layar ini membaca
  // `TREATY_IN`, yang berisi 1.854 baris. Layar ini tidak menunggu tiket 59,
  // dan petunjuk yang menunjuk tiket yang salah membuat orang menagih
  // pekerjaan yang tidak akan mengubah apa pun di sini.
  kosongPetunjuk: '',
  keterangan: '',
  /** Sumber datanya, dinyatakan di layar — bukan hanya di komentar kode. */
  catatanSumber: '',
  /** Kolom yang model baru BELUM punya rumahnya — dinyatakan, bukan diisi. */
  // ⚠️ RALAT 3 Oktober 2026 — kalimat lamanya menyangkal kolom yang ADA.
  // "Position To" memang belum punya rumah di MODEL BARU, dan itu tetap
  // benar; tetapi layar ini membaca TABEL WARISAN, dan di sana
  // `POSITIONUSERNAME` ada. Dari 1.854 baris, 32 terisi dan 1.822 NULL —
  // jadi selnya kosong karena nilainya memang kosong, bukan karena kolomnya
  // tidak ada.
  /** Judul lipatan keterangan kaki — pendek, dan menyebut isinya. */
  kakiJudul: 'Tentang data di tabel ini',
  catatanPosisiKe: '',
  // ⚠️ RALAT 3 Oktober 2026 — di tabel warisan keduanya NAMA, bukan pengenal.
  // `CEDING` dan `LEADINGREINSSOURCE` memuat teksnya langsung. Penyangkalan
  // lama berlaku untuk model baru, yang memang hanya menyimpan pengenal.
  catatanNama: '',
} as const

// ===========================================================================
// LAYAR B — FORM KONTRAK
// Sumber: `Section/TreatyInNONProportional.xml` (tata letak dan ikatan medan),
// `Section/InputTreatyInOffer.xml` (kepala: ID dan Reinsurance Type)
// ===========================================================================

export const FORM_KONTRAK = {
  judul: 'Input Treaty In', // pyValue "Input Treaty In"
  kembali: 'Kembali ke daftar',

  // --- kepala
  id: 'ID', // TreatyIn.ID
  jenisReasuransi: 'Reinsurance Type', // TreatyIn.ProportionType @7236037
  proporsional: 'Proportional',
  nonProporsional: 'Non Proportional',

  // --- kolom kiri
  namaKontrak: 'Treaty Contract Name', // TreatyIn.TreatyContractName · pyCaption @6518125
  nomorRujukan: 'Contract Ref No', // TreatyIn.ContractRefNo · @6510719
  lingkupWilayah: 'Teritorial Scope', // TreatyIn.TeritorialScope · @6507087
  bordereaux: 'Bordereaux', // TreatyIn.Bordeaux · @6527658
  bordereauxCatatan: 'Bordereaux Note', // TreatyIn.BordereauxNote · @6504996

  // --- kolom kanan
  mulai: 'Commencement', // TreatyIn.Commencement · @6540017
  berakhir: 'Termination', // TreatyIn.Termination · @6523491
  tahunTreaty: 'Treaty Year', // TreatyIn.TreatyYear · @6493381
  // ⚠️ SATU label untuk DUA properti — lihat `caraPembukuanNonPropNilai`.
  caraPembukuan: 'Accounting Mode', // TreatyIn.AccountingMode · @6535174
  cedant: 'Ceding', // TreatyIn.Ceding, tampil .ClientName · @6503393
  pilihCedant: 'Choose Ceding', // pyButtonLabel @6513401
  pemimpinTreaty: 'RNM as Treaty Leader', // TreatyIn.TreatyLeader · @6497589
  asalBisnis: 'Source of Business', // TreatyIn.LeadingReinsSource · @6540552
  pilihAsalBisnis: 'Choose Source of Business', // pyButtonLabel @6518676

  /**
   * Nilai pilihan yang ekspor perlihatkan — SATU per pilihan, dan hanya itu.
   *
   * ⚠️ Daftar penuhnya TIDAK ada di Section: Pega mengisinya saat jalan dari
   * sumber di luar ekspor. Yang tertulis di sini nilai yang sungguh terlihat;
   * menambah nilai kedua berarti mengarang.
   */
  /**
   * ⚠️ RALAT 3 Oktober 2026 — nilai NYATA, bukan yang layar lama tampilkan.
   *
   * Bentuk sebelumnya satu nilai per pilihan, disalin dari tangkapan layar
   * ekspor: `Reporting` dan `Accounting Year`. Sapuan atas 1.854 dokumen
   * `M_TREATY_IN.JSONDATA` menemukan DUA nilai per pilihan, dan keduanya
   * huruf kecil:
   *
   *   Bordeaux        `reporting` 1.164 · `nonreporting` 687 · tidak ada 3
   *   AccountingMode  `underwriting` 1.110 · `accounting` 741 · tidak ada 3
   *
   * ⛔ Pemetaan dari yang TERSIMPAN ke yang TAMPIL ("reporting" ->
   * "Reporting") tidak ada di ekspor mana pun, jadi ia TIDAK dikarang: yang
   * ditampilkan nilai apa adanya. Begitu pemetaannya diberikan, di sinilah
   * ia dipasang.
   */
  bordereauxNilai: ['reporting', 'nonreporting'] as readonly string[],
  caraPembukuanNilai: ['underwriting', 'accounting'] as readonly string[],

  /**
   * ⭐ `Accounting Mode` adalah DUA medan, bukan satu — dan sampai 5 Oktober
   * 2026 layar ini hanya punya satu, jadi cabang non-prop membaca properti
   * cabang seberang.
   *
   * Terukur di `Section/TreatyInNONProportional.xml` (bersih 6.904.904
   * bita), dan yang memutuskan BUKAN jaraknya melainkan `pyCondition` sel
   * masing-masing:
   *
   *   TreatyIn.AccountingMode         pyValue @108019 · pyCondition @111777
   *                                   `TreatyIn.ProportionType='Proportional'`
   *   TreatyIn.AccountingModeNonProp  pyValue @114500 · pyCondition @118211
   *                                   `TreatyIn.ProportionType='NonProportional'`
   *
   * Keduanya `pxDropdown`, keduanya berlabel "Accounting Mode", keduanya
   * ber-`pyDisabledWhen = TreatyIn.EDMMaterialType = 2`.
   *
   * ⛔ Nilainya himpunan yang BERBEDA. Sapuan 1.854 dokumen, 5 Oktober 2026:
   *
   *   AccountingMode         underwriting 1.110 · accounting 741 · tidak ada 3
   *   AccountingModeNonProp  loss 1.830 · risk 21 · tidak ada 3
   *
   * ⚠️ Keduanya ADA di KEDUA cabang di dalam dokumen — `AccountingModeNonProp`
   * bernilai `loss` pada SELURUH 1.079 kontrak proporsional. Yang memisahkan
   * keduanya layar, bukan data. Jadi "medan ini milik cabang itu" adalah
   * aturan TAMPIL, dan ia tidak dapat disimpulkan dari isi dokumen.
   */
  caraPembukuanNonPropNilai: ['loss', 'risk'] as readonly string[],

  // --- grid Rate of Exchange
  kurs: 'Rate of Exchange', // pyCaption @6502334
  kursMataUang: 'Currency', // .Currency
  kursKeIDR: 'Value to IDR', // pyCaption @6499731
  kursBerlakuDari: 'Valid From', // pyCaption @6501284
  kursBerlakuSampai: 'Valid Until', // pyCaption @6503943
  tambah: 'Add', // pyButtonLabel
  /** Teks grid kosong bawaan Pega — bukan rule `pyCaption`. */
  tanpaBaris: 'No items',
  /** Petunjuk di bawah grid kurs yang kosong — menyebut tiket pengisinya. */
  kursPetunjuk: '',
  /**
   * Keterangan KAKI, bukan isi layar.
   *
   * ⛔ Dulu paragraf di tengah form, di antara medan dan grid kurs. Prosa
   * pengembang yang duduk di jalur baca pemakai membuat layar terbaca sebagai
   * catatan rilis; ia turun ke kaki, redup, satu baris.
   */

  // =========================================================================
  // PEMILIH "Choose …" — keputusan pemilik proses 4 Oktober 2026
  // =========================================================================

  /**
   * ⛔ Nama KEMBAR ditandai, tidak disatukan dan tidak dibuang.
   *
   * Terukur: `REASURANSI MAIPARK INDONESIA` tercatat dengan 3 pengenal, dan
   * tujuh nama lain dengan 2 — di antaranya `ASURANSI ADIRA DINAMIKA`, yang
   * satu pengenalnya (`ASM-SFAGIS-WORK-ORG ORG-34`) adalah ID kerja Pega yang
   * menyelinap menjadi data. Pemilih yang menyatukannya harus DIAM-DIAM
   * memilih pengenal mana yang ditulis, dan kontrak nyangkut ke cedant yang
   * keliru tanpa seorang pun melihatnya sampai rekonsiliasi.
   */
  pilihKolomID: 'ID',
  pilihKolomNama: 'Name',
  pilihCari: 'Ketik untuk menyaring',
  pilihKembar: 'nama ini punya lebih dari satu pengenal',
  pilihKosong: 'Tidak ada yang cocok.',
  pilihKosongPetunjuk: 'Kosongkan penyaringnya untuk melihat seluruh daftar.',
  pilihMemuat: 'Memuat pilihan…',

  /** Kolom tab Portfolio - dari kunci `Portfolio` di dokumen warisan. */
  petunjukPortofolio: '',
  petunjukAkumulasi: '',
  petunjukPeriode: '',

  /** Petunjuk umum grid tab kosong — menyebut sebabnya, bukan cacahnya. */
  petunjukTabel: '',

  /**
   * ⭐ PETUNJUK LAMA DICABUT 5 Oktober 2026, dan sebabnya hilang bersamanya.
   *
   * Bentuk sebelumnya berbunyi: *"Dibaca dari tabel layer sistem lama
   * (M_TREATY_IN2), yang mencakup 1.340 dari 1.854 kontrak … 510 kontrak
   * punya limit di dokumennya tanpa satu baris pun di sana."*
   *
   * ⛔ Kalimat itu TIDAK BENAR LAGI. `M_TREATY_IN2` dicabut sebagai sumber;
   * keempat tab kini dibaca dari `M_TREATY_IN.JSONDATA`, dokumen yang sama
   * yang memuat kontraknya. Jangkauannya 1.340 → 1.850 kontrak, dan lubang
   * 510 itu TERTUTUP — dijaga `TestLubang510Tertutup`.
   *
   * ⭐ Karena itu kosong kembali punya SATU arti, dan petunjuknya menjadi
   * sama dengan tujuh tab lain: dokumennya memang tidak punya. Membiarkan
   * petunjuk lama berarti layar menjelaskan sebab yang sudah tidak ada.
   */
  petunjukLayer: '',
  /**
   * ⭐ Dan sebab KETIGA pun hilang.
   *
   * Bentuk sebelumnya menyatakan tiga dari empat batas — Flood Jabodetabek,
   * Flood Nationwide, RSMD — *"tidak ada di tabel itu dan belum punya sumber
   * tabel"*. Itu benar tentang `M_TREATY_IN2`, dan keliru tentang dokumen:
   * keempatnya ada di `Limits[].Detail[]` — `RSMDLimit` 493 · `Earthquake`
   * 569 · `FloodJab` 218 · `FloodNation` 550, masing-masing dengan mata
   * uangnya. Gambar `28` memperlihatkan keempat barisnya di layar.
   */
  petunjukEventLimits: '',
  /**
   * ⛔ Kedua tab teks punya pesan kosongnya SENDIRI. `tanpaBaris` berbunyi
   * tentang baris grid; tab ini tidak punya baris, ia punya satu medan.
   */
  tanpaTeks: 'Tidak ada teks',
  ejaanDipakai: 'Dibaca dari kunci:',
  /**
   * ⛔ Kalimat ini menyatakan ada teks yang pembacanya TIDAK lihat, dan ia
   * harus berbunyi begitu. Terukur: dari 303 dokumen yang punya lebih dari
   * satu ejaan `SpecialConditions*`, NOL yang isinya identik — jadi ejaan
   * lain BUKAN salinan, melainkan teks yang berbeda.
   */
  ejaanLainBerisi:
    '⚠️ Dokumen kontrak ini juga punya teks di bawah kunci lain, dan isinya BERBEDA — bukan salinan. Kunci itu:',
  petunjukTeksPengecualian: '',
  petunjukTeksSyarat: '',

  belumDibangun: 'Tab ini belum dibangun.',
  /**
   * ⛔ KEADAAN KEEMPAT, dan ia BUKAN "belum ada kode" maupun "tidak dipakai
   * lagi". Keputusan pemilik proses 4 Oktober 2026, `KEPUTUSAN §17`.
   *
   * "Belum ada kode" mengundang orang menagih pembangunannya. "Tidak dipakai
   * lagi" keliru, dan verifikasi membuktikannya: ketiga penjaga `1=2` di
   * `Section/ShareRetro.xml` membungkus sebuah tombol, satu blok BAR, dan
   * satu tombol kepala — NOL yang membungkus tabnya. `FlowAction` berbunyi
   * `pyRuleAvailable = Yes`. Retro HIDUP di Pega; ia jarang.
   *
   * Yang benar: fiturnya ada, dipakai lima kontrak, dan tidak dibangun
   * karena itu. Orang yang membacanya tahu persis apa yang ia lihat.
   */
  jarangDipakai: 'Tab ini tidak dibangun — fiturnya JARANG dipakai.',
  jarangDipakaiPetunjuk: '',
  /**
   * ⛔ Pesan ini pernah menyebut "hanya tab Reporting Period", dan menjadi
   * BASI diam-diam begitu tab kedua mendarat. Kalimat yang menyebut cacah
   * atau nama tab akan selalu basi; yang ini menyebut SEBABNYA, dan sebab
   * tidak berubah tiap ronde.
   */
  belumDibangunPetunjuk: '',
} as const

/**
 * Strip tab — DUA himpunan, dan keduanya TIDAK sama.
 *
 * ⛔ Temuan ronde ini, dan ia berlawanan dengan apa yang diperkirakan:
 * `TreatyInTabsProportional.xml` dan `TreatyInTabsNonProportional.xml`
 * **tidak memuat daftar tab yang sama**. Membangun satu strip yang dialihkan
 * radio akan memperlihatkan tab yang di sistem lama tidak pernah ada pada
 * cabang itu.
 *
 * Diurutkan menurut posisi bitanya di dalam masing-masing berkas, sesudah
 * `<pyIncludedRuleXML>` dibuang. Judul yang merupakan PANEL DI DALAM sebuah
 * tab tidak ikut — disebut di komentar supaya pemisahannya dapat diperiksa.
 */
export const TAB_PROPORSIONAL = [
  'Reporting Period', // @26514 — panel dalamnya "Account Reporting Period" @45828
  'Portfolio', // @417075
  'Limits', // @526981
  'Share', // @651816 — panel dalamnya "Total Share" @662148
  'Retro', // @1031030 — ⭐ BERSYARAT, lihat SYARAT_TAB
  'Co-Ins Scale', // @1077949
  'Accumulation', // @1230754 — panel dalamnya "Accumulation Control" @1250071
  'Exclusions', // @1455282
  'Special Conditions', // @1497175
  'Information & Submit', // @1539099
  'Achievement In IDR', // @1581141
] as const

export const TAB_NON_PROPORSIONAL = [
  'Maximum Retention', // @15962
  'Event Limits', // @223612
  'EGNPI', // @490872 — panel "Estimate Gross Net Premium Income" @509493
  'Limits', // @886054 — panel "Summary of Limit" @1061345, "Summary of MDP" @1204191, "Total All Layers" @1292282
  'Share', // @1695720
  // ⛔ `RNM Share` — PERTANYAAN TERBUKA, bukan tab menurut ekspor.
  //
  // Briefing 5 Oktober 2026 menyebutnya tab yang hilang dari tangkapan
  // layar. Pengukuran berkata lain, dan ia TIDAK dihapus sepihak — §0
  // ronde itu melarang menghapus tab mana pun. Yang terukur:
  //
  // `Section/TreatyInTabsNonProportional.xml` (bersih 4.203.858 bita) punya
  // TEPAT 11 wadah `pyHeaderType = TABBED`, dan `RNM Share` BUKAN salah
  // satunya. Kedua judul "RNM Share" (@2093710, @2133893) ber-`BAR`, dan
  // keduanya jatuh DI ANTARA tab `Share` (TABBED @1695720) dan tab `Retro`
  // (TABBED @3291797) — artinya ia panel DI DALAM tab Share, bersama
  // "Summarry of RNM Share" @2324385 dan "Total All Layers RNM Share"
  // @2486274.
  //
  // ⚠️ "BAR berarti bukan tab" TIDAK berlaku umum: di
  // `TreatyInTabsProportional.xml` nol wadah ber-TABBED dan kesebelas
  // tabnya justru `BAR`. Yang membedakan adalah pemakaian DI DALAM satu
  // berkas, dan di berkas non-prop TABBED-lah tabnya.
  //
  // Hitungannya pun cocok: 11 tab − `Value Difference` (syaratnya tidak
  // terpenuhi) = 10, persis yang tangkapan layar perlihatkan. Dengan
  // `RNM Share` ikut, daftar ini 12 — satu lebih banyak daripada ekspor.
  'RNM Share', // @2093710 — ⛔ lihat di atas
  'Retro', // @3291797
  'Installment', // @3465901
  'Value Difference', // @3848946 — ⭐ BERSYARAT, lihat SYARAT_TAB
  'Exclusions', // @3907436
  'Special Conditions', // @3942540
  'Information & Submit', // @3977799
] as const

/**
 * Sub-tab DI DALAM tab `Share` — kedua cabang.
 *
 * ⭐ Terbaca dari dokumen desain 5 Oktober 2026: gambar `16`/`17` (prop) dan
 * `34` (non-prop) memperlihatkan strip sub-tab berjudul `RNM Share` di bawah
 * panel atas tab `Share`.
 *
 * ⛔ Ini MENJAWAB pertanyaan §3 ronde sebelumnya, dan jawabannya sejalan
 * dengan pengukuran ekspor: `RNM Share` bukan salah satu dari 11 wadah
 * `TABBED`, dan kedua judulnya jatuh di dalam wilayah tab `Share`.
 *
 * ⚠️ `RNM Share` TETAP ADA di `TAB_NON_PROPORSIONAL` — nol tab dihapus.
 */
export const SUB_TAB_SHARE = ['RNM Share'] as const

// ===========================================================================
// Tab `Share` cabang PROPORSIONAL — bentuknya dari ekspor, bukan dari gambar
// ===========================================================================
//
// `Section/TreatyInShareProp.xml` (sesudah `pyIncludedRuleXML` bersarang
// dibuang dengan hitung kedalaman) menyebut seluruhnya apa adanya:
//
//   pyTitle   `RNM Share`
//   pyValue   `Kind of Treaty` · `.TreatyType` · `.Note`
//             `Total Share RNM Limit`   + `Value`  ← `.Currency` / `.Value`
//             `Total Value Spreading OR`  + `Value`
//             `Total Value Spreading R/I` + `Value`
//
// ⚠️ Cabang NON-PROPORSIONAL tidak memakai bentuk ini — ia punya gridnya
// sendiri per layer. Satu komponen untuk keduanya akan menampilkan kolom
// yang di cabang seberang tidak pernah ada.
export const TOTAL_SHARE = {
  judul: 'Total Share',
  segarkan: 'Refresh', // pyActionLabel `Refresh`
  persenRnmShare: '% RNM Share', // pyLabelFieldValue `% RNM Share` (Share.xml)
  persenBrokerage: '% Brokerage',
  opsi: 'Option',
  /** Petunjuk DI DALAM kotak — `pyValue` `%` di ekspor. */
  satuanPersen: '%',
  tanpaBaris: 'No items',
  /**
   * ⚠️ Petunjuk ketiga grid total yang KOSONG, dan ia menyatakan sebab yang
   * tepat: sumbernya belum punya tabel — BUKAN kontraknya yang kosong.
   * Keduanya terlihat sama di layar, dan hanya kalimat ini yang membedakan.
   */
  petunjukTotal: '',
} as const

/** Judul grid `Kind of Treaty` beserta kolom nilainya. */
export const KOLOM_KIND_OF_TREATY_SHARE = ['Kind of Treaty'] as const

/**
 * Ketiga grid total sub-tab `RNM Share`, berurut seperti di ekspor.
 *
 * ⛔ Kolom keduanya BERJUDUL `Value` pada ketiganya — itu bunyi ekspornya,
 * dan menamainya sendiri ("Jumlah", "Nilai") akan membuat layar berbeda dari
 * layar lama tanpa ada yang memintanya.
 */
export const GRID_TOTAL_RNM_SHARE = [
  'Total Share RNM Limit',
  'Total Value Spreading OR',
  'Total Value Spreading R/I',
] as const

export type TabProporsional = (typeof TAB_PROPORSIONAL)[number]
export type TabNonProporsional = (typeof TAB_NON_PROPORSIONAL)[number]

/**
 * Syarat TAMPIL per tab — dibaca dari `pyContainerVisibleWhen` di ekspor.
 *
 * ⛔ Tab yang syaratnya TIDAK terpenuhi tidak dirender sama sekali. Itu
 * keadaan yang berbeda dari `Kosong` (tabnya ada, isinya nol) dan dari
 * `.trin__belum` (tabnya ada, kodenya belum ditulis): tabnya memang bukan
 * bagian dari layar kontrak ini.
 *
 * ⚠️ Yang BUKAN di sini: tab tanpa syarat. Daftar ini sengaja memuat hanya
 * yang ekspornya memberi syarat — tab yang tidak disebut selalu tampil, dan
 * menuliskannya di sini sebagai `() => true` akan membuat daftar ini
 * terbaca seolah seluruh daftar tab ada dua kali.
 *
 * Terukur 5 Oktober 2026:
 *
 *   Retro (prop)              `TreatyIn.IsMultipleRetro`
 *     `Section/TreatyInTabsProportional.xml` pyTitle @1031030,
 *     pyContainerVisibleWhen di sel yang sama.
 *   Value Difference (nonprop) `TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1`
 *     `Section/TreatyInTabsNonProportional.xml` TABBED @3848982; syarat yang
 *     SAMA muncul di `InputTreatyInOffer.xml` @5479082 dan
 *     `TreatyInNONProportional.xml` @5840210 — tiga ekspor, satu kalimat.
 */
export interface SyaratTabKontrak {
  /** `TreatyIn.IsMultipleRetro` apa adanya — teks, bukan boolean. */
  retroBerganda: string
  /** `TreatyIn.EDMState` apa adanya. */
  edmState: string
  /** `TreatyIn.EDMMaterialType` apa adanya. */
  edmJenisMaterial: string
}

export type UjiSyaratTab = (k: SyaratTabKontrak) => boolean

/**
 * ⛔ BERCABANG, dan itu BUKAN kerapian — ia perbedaan yang terukur.
 *
 * `Retro` ada di KEDUA daftar tab, dan hanya yang PROPORSIONAL bersyarat:
 *
 *   TreatyInTabsProportional.xml     pyTitle `Retro` @1031030
 *                                    pyContainerVisibleWhen `TreatyIn.IsMultipleRetro`
 *   TreatyInTabsNonProportional.xml  TABBED  `Retro` @3291822
 *                                    NOL pyContainerVisibleWhen
 *
 * Peta bernama-tab-saja akan menyembunyikan `Retro` non-proporsional pada
 * 770 dari 775 kontrak yang berhak melihatnya — dan ekspornya tidak pernah
 * memintanya.
 */
export const SYARAT_TAB_PROPORSIONAL: Readonly<Record<string, UjiSyaratTab>> = {
  // `TreatyIn.IsMultipleRetro` — satu properti telanjang sebagai syarat
  // berarti "benar". Nilainya TEKS di dokumen: `"true"` pada 5 dari 1.854
  // kontrak, `"false"` pada 1.531, dan kuncinya tidak ada pada 318.
  //
  // ⛔ Kunci yang TIDAK ADA diperlakukan tidak-terpenuhi, dan itu bukan
  // sama dengan `false` — ia hanya kebetulan berakhir di tab yang sama.
  // Perbedaannya tercatat; yang dicatat dapat dibalik, yang dilebur tidak.
  Retro: (k) => k.retroBerganda === 'true',
}

export const SYARAT_TAB_NON_PROPORSIONAL: Readonly<Record<string, UjiSyaratTab>> = {
  // `TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1`.
  //
  // ⚠️ `!=` Pega terhadap nilai yang TIDAK ADA bernilai BENAR — properti
  // kosong bukan 3. Jadi yang mengikat praktis syarat keduanya.
  'Value Difference': (k) => k.edmState !== '3' && k.edmJenisMaterial === '1',
}

// ===========================================================================
// TAB — REPORTING PERIOD (Account Reporting Period)
// Sumber: `Section/TreatyInTabsProportional.xml`
// ===========================================================================

export const REPORTING_PERIOD = {
  judul: 'Account Reporting Period', // pyTitle @45828
  mulai: 'Start Date', // pyCaption @1586974 · TreatyIn.ReportingStart
  akhir: 'End Date', // pyCaption @1673706 · TreatyIn.ReportingEnd
  periode: 'Period', // pyCaption @1610195 · TreatyIn.ReportingPeriod
  periodeNilai: 'Quarter Year', // pyCaption @7251718
  penyerahan: 'Submission', // pyCaption @1671554 · TreatyIn.ReportingSubmission
  konfirmasi: 'Confirmation', // pyCaption @1624384 · TreatyIn.ReportingConfirmation
  pelunasan: 'Settlement', // pyCaption @1605349 · TreatyIn.ReportingSettlement
  hari: 'Days', // pyValue "Days"
  terapkan: 'Apply', // pyButtonLabel @1662987

  // grid
  kolomPeriode: 'Period', // .Period
  kolomOtomatis: 'Auto Calculate', // pyCaption @1611297 · .AutoCalculate
  kolomTanggalAwal: 'Initial Date', // pyCaption @1606483 · .InitialDate
  kolomPenyerahan: 'Submission Due', // pyCaption @1587534 · .SubmissionDue
  kolomKonfirmasi: 'Confirmation Due', // pyCaption @1609623 · .ConfirmationDue
  kolomPelunasan: 'Settlement Due', // pyCaption @1639462 · .SettlementDue
  tanpaBaris: 'No items',

  /**
   * Pesan galat — DISALIN APA ADANYA, dan ia dirangkai dari DUA rule.
   *
   * `pyCaption Start Date, Due` @1612960 + `pyCaption . Must Not Be Empty`
   * @1609072. Pega merangkainya saat jalan; di sini ia ditulis utuh supaya
   * teks yang sampai ke layar tidak bergantung pada perangkaian kita.
   *
   * ⚠️ Ada SAUDARANYA di ekspor yang sama: `pyCaption , and Interval`
   * @1629752, yang membentuk *"Start Date, Due, and Interval. Must Not Be
   * Empty"*. Varian itu BELUM dipakai di sini — ia milik jalur yang juga
   * memeriksa Interval, dan jalur itu belum dibangun.
   */
  galatKosong: 'Start Date, Due. Must Not Be Empty',
  /** Varian bila Period `other` — `, and Interval` @276424 tampil (`ReportingPeriod='other'`). */
  galatKosongInterval: 'Start Date, Due, and Interval. Must Not Be Empty',
  interval: 'Interval', // @144723 · TreatyIn.ReportingInterval — tampil bila Period `other`
} as const

/**
 * Kolom tab Portfolio — judul dari gambar `02`, PADANAN KUNCI dari ekspor.
 *
 * ---------------------------------------------------------------------
 * ⛔ RALAT 6 Oktober 2026 — KEDUA KOLOM PERTAMA TERTUKAR SEJAK 5 Oktober.
 * ---------------------------------------------------------------------
 * Bunyi sebelumnya: *"kunci `Type` menjadi kolom `Portfolio Type`"*. Itu
 * TERBALIK, dan layar memperlihatkan `Premium`/`Loss` di kolom yang
 * berjudul `Portfolio Type`.
 *
 * Yang membuktikannya `Section/TreatyInTabsProportional.xml` — grid
 * `TreatyIn.Portfolio` (`pyPageListPropertyClass` =
 * `ASM-FW-GISFW-Data-TreatyInPortfolio`) punya DUA baris sel yang lurus
 * berpasangan, dan pasangannya tidak dapat ditafsirkan dua cara:
 *
 *   kolom  judul (baris 1)        properti (baris 2)
 *     1    `Portfolio Type`  107   `.TypePortfolio`  112
 *     2    `Premium / Loss Type` 108   `.Type`       113
 *     3    `Description`     109   `.Description`    114  (pxTextArea)
 *     4    tombol `Add`      110   tombol `Delete`   115
 *
 * ⭐ Dan dua saksi LAIN mengatakan hal yang sama, tanpa saling menyalin:
 *
 *   1. Judul kolomnya sendiri. `Premium / Loss Type` MENYEBUT kedua
 *      nilainya, dan yang bernilai `Premium`/`Loss` adalah kunci `Type` —
 *      terukur, lihat `PORTOFOLIO.opsiJenis` di bawah.
 *   2. Ekspor modul Treaty In Adjustment, dua seksi terpisah
 *      (`TreatyInTabsProportional`, `TreatyInTabsProportionalOldData`),
 *      berpasangan persis sama.
 *
 * ⚠️ Ejaan kuncinya `TypePortfolio`, TANPA spasi. Catatan lama menulis
 * `Type Portfolio`; nol dokumen memakai ejaan itu — 1.925 elemen yang
 * disapu seluruhnya memakai `TypePortfolio`.
 */
export const KOLOM_PORTOFOLIO = [
  'Portfolio Type', // sel 107 ↔ 112 · kunci dokumen `TypePortfolio`
  'Premium / Loss Type', // sel 108 ↔ 113 · kunci dokumen `Type`
  'Description', // sel 109 ↔ 114 · kunci dokumen `Description`
] as const

/**
 * Tab Portfolio — isian, tombol, dan PILIHAN kedua dropdown-nya.
 *
 * ---------------------------------------------------------------------
 * DARI MANA PILIHANNYA DATANG — dan apa yang ekspor TIDAK katakan
 * ---------------------------------------------------------------------
 * Kedua sel bernilai `pyControlDisplayTitle` = *"Control inherited from
 * property"*, dan `pyListDataSource`-nya KOSONG. Artinya daftar pilihannya
 * hidup di `Rule-Obj-Property` kelas `ASM-FW-GISFW-Data-TreatyInPortfolio`
 * — dan aturan properti TIDAK ikut diekspor ke `D:\XML_NURE` (korpusnya
 * hanya Activity · ConnectREST · DataTransform · DecisionTable ·
 * FlowAction · Harness · RDBList · ReportDefinition · Section ·
 * SystemSettings · When).
 *
 * ⛔ Jadi daftarnya TIDAK dikarang, dan juga TIDAK diambil dari ekspor —
 * ia DIUKUR dari datanya sendiri. Sapuan 6 Oktober 2026 atas SELURUH 1.855
 * dokumen `POOLDATA.M_TREATY_IN` (nol `ROWNUM`, nol gagal urai): 1.925
 * elemen `Portfolio[]` di 845 dokumen.
 *
 *   `TypePortfolio`   Withdrawal 1.578 · Assumption 347   — 2 nilai, 0 kosong
 *   `Type`            Premium    1.032 · Loss       893   — 2 nilai, 0 kosong
 *
 * Keempat kombinasinya terpakai (845 · 733 · 187 · 160), sejalan dengan
 * `UQ_PORTOFOLIO (ID_VERSI_KONTRAK, ARAH_PORTOFOLIO, JENIS_PORTOFOLIO)`
 * pada migrasi `406` dan dengan `TestTiket24PortofolioEmpatKombinasiDiterima`.
 *
 * ⚠️ YANG MASIH TERBUKA, dan ia ditulis bukan ditebak: himpunan terukur
 * adalah yang TERPAKAI, belum tentu seluruh yang aturan propertinya
 * izinkan. Nilai yang sah tetapi belum pernah dipakai akan hilang dari
 * daftar ini. Pertanyaannya di
 * `docs/PERTANYAAN-TERBUKA-PILIHAN-PORTFOLIO.md`.
 *
 * ⚠️ Pilihan KOSONG tetap ada, dan itu bukan kelalaian: tombol `Add`
 * menjalankan `Activity/TreatyInPropAdd.xml`, yang hanya menetapkan
 * `TreatyIn.Portfolio(<APPEND>).Description = ""` — baris barunya lahir
 * dengan kedua pilihan BELUM terisi.
 */
export const PORTOFOLIO = {
  judul: 'Portfolio',
  tambah: 'Add', // sel 110 · Activity `TreatyInPropAdd`, param Type="portfolio"
  hapus: 'Delete', // sel 115 · `Embed-SelectedContextAPI-DeleteRow`
  belumDipilih: '', // teks pilihan kosong — layar Pega menampilkannya kosong
  /** Kolom 1 `Portfolio Type` ← `TypePortfolio`. Terukur, 1.925 elemen. */
  opsiArah: ['Withdrawal', 'Assumption'] as const,
  /** Kolom 2 `Premium / Loss Type` ← `Type`. Terukur, 1.925 elemen. */
  opsiJenis: ['Premium', 'Loss'] as const,
} as const

/**
 * Kolom tab Accumulation — RALAT 5 Oktober 2026, dari gambar `19`.
 *
 * ⛔ Ketiga judul singkat sebelumnya (`Report Date`, `Sub Days`,
 * `Sub Due Date`) adalah nama kunci dokumen, bukan judul layar.
 */
export const KOLOM_AKUMULASI = [
  'Period', // gambar 19
  'Reporting Date', // gambar 19 · kunci dokumen `Report Date`
  'Submission Days', // gambar 19 · kunci dokumen `Sub Days`
  'Submission Due', // gambar 19 · kunci dokumen `Sub Due Date`
] as const

/** Kolom tab EGNPI — nama kunci di dokumen warisan apa adanya. */
/**
 * ⛔ URUTAN DAN JUDUL dari gambar `29`, bukan urutan abjad kunci dokumen.
 *
 * Bentuk sebelumnya delapan kunci dokumen berurut abjad. Layar Pega
 * memperlihatkan ENAM kolom dalam urutan yang berbeda:
 *
 *   Treaty Group · As Date · Proportion % · Currency · Amount · Amount in IDR
 *
 * ⚠️ `Class of Business` dan `Note` TIDAK dihapus. Keduanya ada di dokumen,
 * dan `Note` terlihat di RINCIAN baris gambar 29 (teks dua baris di bawah
 * `Proportion %`). Yang tidak terlihat di gambar belum tentu tidak ada —
 * keduanya pindah ke ujung, bukan hilang.
 */
export const KOLOM_EGNPI = [
  'Treaty Group', // gambar 29
  'As Date', // gambar 29
  'Proportion %', // gambar 29 · kunci dokumen `Proportion`
  'Currency', // gambar 29
  'Amount', // gambar 29
  'Amount in IDR', // gambar 29 · kunci dokumen `Amount IDR`
  'Class of Business', // ⚠️ tidak terlihat di gambar 29 — lihat di atas
  'Note', // ⚠️ di RINCIAN baris, bukan di gridnya
] as const

/**
 * ⚠️ `Proportion` adalah persen SHARE, bukan persen biasa — dan itu diukur,
 * bukan ditebak. Dijumlahkan per kontrak ia menghasilkan **tepat 100** (6 dari
 * 6 kontrak bersampel pada 4 Oktober 2026), dan nilainya menyimpan 20 desimal
 * (`41.63535247256876597000`). Keduanya persis alasan §13 memilih 8 desimal:
 * sifat jumlah-tepat-100 harus terjaga, dan 20 melampaui `NUMBER(38,8)`.
 */
/**
 * ⭐ DESIMAL PER KOLOM — §24, dibaca dari gambar `29`.
 *
 * ⛔ Dua kolom di BARIS YANG SAMA berbeda, dan itu justru buktinya:
 *
 *   `Amount`         `137.849.315.068,00`   2 desimal
 *   `Amount in IDR`  `137.849.315.068`      0 desimal
 *
 * `Proportion %` tampil `100,00` di grid dan `100,0000000000` di rinciannya.
 * Yang dipasang di sini **2** — ini gridnya. Angka 10 milik rincian baris,
 * dan layar kita belum punya baris yang dapat dibuka.
 */
export const JENIS_EGNPI: readonly JenisAngka[] = [
  'teks', // Treaty Group
  'teks', // As Date
  ['persenShare', 2], // Proportion % — gambar 29, grid
  'teks', // Currency
  ['uang', 2], // Amount — gambar 29
  ['uang', 0], // Amount in IDR — gambar 29, baris yang SAMA
  'teks', // Class of Business
  'teks', // Note
]

/** Kolom tab Maximum Retention. */
/**
 * ⛔ URUTAN dari gambar `26`/`27`, bukan abjad.
 *
 * Grid layar lama berkolom TIGA — `Treaty Group` · `Currency` · `Amount` —
 * dan barisnya dapat dibuka menjadi rincian `Treaty Group` · `Amount`
 * (mata uang + nilai) · `Note` (gambar 27).
 *
 * ⚠️ `Class of Business` dan `Note` TIDAK dihapus, alasan yang sama dengan
 * `KOLOM_EGNPI`: `Note` terlihat di rinciannya, dan yang tidak terlihat
 * belum tentu tidak ada.
 */
export const KOLOM_RETENSI = [
  'Treaty Group', // gambar 26
  'Currency', // gambar 26
  'Amount', // gambar 26
  'Class of Business', // ⚠️ tidak terlihat di gambar 26/27
  'Note', // ⚠️ di RINCIAN baris (gambar 27), bukan di gridnya
] as const

/**
 * ⚠️ `Amount` di GRID gambar 26 berbunyi `3.500.000.000` — NOL desimal —
 * sementara panel `Total Retention Amount` tepat di bawahnya berbunyi
 * `3.500.000.000,00`. Nilai yang sama, dua presisi, dua tempat.
 */
export const JENIS_RETENSI: readonly JenisAngka[] = [
  'teks', // Treaty Group
  'teks', // Currency
  ['uang', 0], // Amount — gambar 26, grid
  'teks', // Class of Business
  'teks', // Note
]

/**
 * Panel `Existing Policy for Master ID` — kanan atas layar, KEDUA cabang.
 *
 * ⛔ Nol jejaknya di kode kita sampai 5 Oktober 2026, dan sumbernya
 * DITELUSURI sampai SQL-nya — bukan dikarang:
 *
 *   Section/InputTreatyInOffer.xml  pyTitle @89.949
 *                                   pyPageListProperty `PolisList.pxResults` @105.310
 *   Activity/FetchTreatyExistingProduction.xml  mengisinya
 *   RDBList/FetchTreatyInProductionUsingNooffer.xml  SQL-nya:
 *
 *     SELECT DISTINCT NOPOLIS, IDPEGA, QUARTER, QUARTER_YEAR
 *       FROM POOLDATA.TREATYINPRODUCTION
 *      WHERE SUBSTR(NOOFFER,1,7) = SUBSTR({InputData.CARI1},1,7)
 *
 * ⭐ Diadu dengan gambarnya dan cocok: kontrak `1001841` (gambar 26) satu
 * baris `RNM-QR.T02.05.2025.11987` / `NB-147044`; kontrak `1001846`
 * (gambar 01) nol baris, dan gambarnya berbunyi `No items`.
 */
export const POLIS_PRODUKSI = {
  judul: 'Existing Policy for Master ID', // pyTitle @89949
  tanpaBaris: 'No items',
} as const

/** Keempat kolomnya, urut layar — `pyValue` di Section-nya. */
export const KOLOM_POLIS_PRODUKSI = [
  'Policy No', // @110151 · NOPOLIS
  'Pega ID', // @114008 · IDPEGA dipotong 18 aksara
  'Quarter', // @117127 · QUARTER
  'Quarter Year', // @121484 · QUARTER_YEAR
] as const

/**
 * ⛔ KEEMPATNYA TEKS, dan itu diukur bukan diasumsikan. `QUARTER` dan
 * `QUARTER_YEAR` bertipe `VARCHAR2` di `TREATYINPRODUCTION`, dan `Policy No`
 * maupun `Pega ID` pengenal — pengenal tidak pernah diformat.
 */
export const JENIS_POLIS_PRODUKSI: readonly JenisAngka[] = ['teks', 'teks', 'teks', 'teks']

/**
 * Panel `Total Retention Amount` dan tombol `Update Total` — keduanya DI
 * BAWAH grid Maximum Retention, dan keduanya hilang dari layar ini sampai
 * 5 Oktober 2026.
 *
 * Terukur di `Section/TreatyInTabsNonProportional.xml` (bersih 4.203.858):
 *
 *   pyPageListProperty  TreatyIn.TotalRetentionAmountNP   @147731
 *   judul kolom 1       `Total Retention Amount`          @152591
 *   judul kolom 2       `Value`                           @156684
 *   isi kolom 1         `.Currency`                       @161105
 *   isi kolom 2         `.Value`                          @166072
 *   tombol `Update Total` pyLabel                         @202657
 *
 * ⛔ Judul panelnya adalah JUDUL KOLOM PERTAMA, bukan tajuk di atas grid —
 * dan kolom pertama itu berisi MATA UANG. Briefing menyebut panel ini
 * berkolom `Value` saja; ekspornya dua kolom.
 */
export const TOTAL_RETENSI = {
  /** Judul kolom pertama = nama panel. Isinya `.Currency`. */
  judul: 'Total Retention Amount', // @152591
  kolomNilai: 'Value', // @156684
  tanpaBaris: 'No items',

  /**
   * ⭐ Tombol `Update Total` — HIDUP di layar lama, MATI di sini.
   *
   * Yang terukur tentang tombol itu, @200987..@207693:
   *
   *   pyFormat        pxButton
   *   pyLabel         `Update Total`        @202657
   *   pyAction        `refresh`             @203023 (param `retention` @203951)
   *   pyDisabledWhen  `TreatyIn.EDMMaterialType = 2`  @202558
   *   pyVisible/OTHER `TreatyIn.IsEditData !='1'`     @207592/@207693
   *
   * ⚠️ Ada tombol KEDUA tepat di atasnya, @194468, dan ia MATI: penjaganya
   * `pyCondition 1=2` @198452. Ia tak berlabel (`.pyTemplateButton`) —
   * kemungkinan pendahulu tombol ini. Yang dibangun yang hidup.
   *
   * ⛔ Di sini ia DINONAKTIFKAN, bukan dihilangkan: layar ini baca-saja,
   * dan totalnya sudah dihitung di services setiap kali kontrak dibaca.
   * Tombol hidup yang tidak mengubah apa pun berbohong; tombol hilang
   * menyembunyikan bahwa layar lama punya langkah ini.
   */
  perbarui: 'Update Total', // pyLabel @202657
  perbaruiPetunjuk: '',

  /**
   * ⚠️ `Update Total` TIDAK hanya milik Maximum Retention.
   *
   * Sapuan berkas yang sama menemukannya pada LIMA tab non-prop —
   * Maximum Retention @202657 · EGNPI @843732 · Limits @1641962 ·
   * Share @3188309 · Installment @3828165. Yang dibangun ronde ini satu;
   * keempat sisanya tercatat di `KEPUTUSAN-PENYELARASAN-REPO.md` §21.
   */
  tabLain: ['EGNPI', 'Limits', 'Share', 'Installment'] as readonly string[],
} as const

/** Kolom JADWAL tab Installment. */
export const KOLOM_ANGSURAN = [
  'Installment', 'Currency', 'Amount', 'Pct', 'Due Date', 'Payment Date', 'WPC',
] as const

/**
 * ⚠️ `Pct` juga persen SHARE: dijumlahkan per kontrak ia 100, atau 200 bila
 * kontraknya berjadwal dua mata uang. Datanya hari ini hanya 2 desimal
 * (`25.00`), jadi 8 maupun 2 sama-sama menampilkan `25%` — penggolongannya
 * tetap `persenShare` supaya ALASANNYA yang tercatat, bukan kebetulan
 * datanya.
 */
/**
 * ⭐ Gambar `38`: `% Installment` berbunyi `25,00` (2) sementara `% Total` di
 * bawahnya `100,0000` (4), dan `Amount` `161.168.704,90` (2).
 *
 * ⚠️ Yang dipasang di sini kolom GRID-nya. `% Total` bukan kolom grid; ia
 * medan panel yang layar kita belum punya.
 */
export const JENIS_ANGSURAN: readonly JenisAngka[] = [
  'teks', // Installment
  'teks', // Currency
  ['uang', 2], // Amount — gambar 38
  ['persenShare', 2], // Pct — gambar 38 (`% Installment`)
  'teks', // Due Date
  'teks', // Payment Date
  'teks', // WPC
]

/** Kolom tab Information & Submit. */
export const KOLOM_CATATAN = ['Date', 'Operator', 'Approved', 'Suggest'] as const

// ===========================================================================
// Tab `Information & Submit` — FORM, bukan grid riwayat
// ===========================================================================
//
// ⛔ CACAT YANG DIPERBAIKI 6 Oktober 2026: tab ini menampilkan grid riwayat
// (`KOLOM_CATATAN` di atas), padahal di Pega ia FORM. Riwayatnya sudah punya
// panelnya sendiri di kaki layar — jadi yang ditampilkan bukan sekadar
// salah, ia SALINAN dari yang sudah ada di layar yang sama.
//
// Bentuknya dari `Section/TreatyInfoSubmit.xml`, sesudah `pyIncludedRuleXML`
// bersarang dibuang dengan hitung kedalaman:
//
//   TreatyIn.Information  label `Additional Information`  Text area
//   TreatyIn.Comment      label `Comment`                 Text area
//
//   tombol  `Submit`         pyActivity TreatyInSubmit     gaya Strong
//           pyCondition      TreatyIn.ViewState !='1'
//                            && TreatyIn.StatusAkseptasi != 'Resolve Complete'
//   tombol  `Decline offer`                                gaya Simple
//           pyCondition      TreatyIn.ViewState != '1'
//
//   pyDisabledWhen  TreatyIn.ID = ''   ·   TreatyIn.EDMEffective = ''
export const INFO_SUBMIT = {
  judul: 'Information & Submit',
  infoTambahan: 'Additional Information', // pyLabelFieldValue
  komentar: 'Comment', // pyLabelFieldValue
  kirim: 'Submit', // pyLabel, gaya Strong
  tolak: 'Decline offer', // pyLabel, gaya Simple
  /**
   * ⚠️ Kedua tombol MATI, dan sebabnya bukan kelalaian: `Submit` memanggil
   * `TreatyInSubmit`, yang bermuara ke prosedur yang menulis `M_TREATY_IN`
   * dan `TREATY_IN` — keduanya dilarang keras dipakai aplikasi. Sasaran
   * tulisnya belum diputuskan pemilik proses.
   */
  petunjukTombol: '',
} as const

// ===========================================================================
// EMPAT TAB DARI `M_TREATY_IN2` — Limits · Share · Event Limits · RNM Share
//
// ⛔ SATU TABEL, EMPAT TAB. Keempatnya proyeksi atas baris yang sama, satu
// baris per layer. Nol tabel baru dibuat untuk mereka.
//
// ⛔ Nama kolom Oracle BUKAN nama properti Pega: dari 41 kolom, hanya 8 yang
// cocok harfiah dengan `pyValue` di ekspor Section. Pemetaannya diturunkan
// dari NILAI — bukti kolom demi kolom di `docs/PEMETAAN-M-TREATY-IN2.md`.
//
// ⚠️ Judul kolom di bawah memakai NAMA KOLOM apa adanya, bukan label layar
// lama. Sebabnya diukur: ekspor Section menamai medannya menurut properti
// Pega (`.Limit`, `.MDP`, `.Deductible`), dan properti itu TIDAK berpadanan
// satu-satu dengan kolom tabel ini. Memberi label layar lama berarti
// mengaku tahu padanan yang justru belum terbukti.
// ===========================================================================

/** Golongan angka satu kolom — menentukan pemformat mana yang dipanggil. */
export type GolonganAngka = 'uang' | 'persen' | 'persenShare' | 'teks'

/**
 * Golongan satu kolom — DAN, bila terbaca di gambar, jumlah desimalnya.
 *
 * ⭐ KEPUTUSAN §24, 5 Oktober 2026: desimal di layar ini **per kolom**, dan
 * nol di ekor **DIPERTAHANKAN** sampai presisi kolomnya. Gambar desain
 * memperlihatkan `1,00` dan `15.000,00`; ronde 69 memperlihatkan `2484250`
 * tanpa ekor. **Keduanya benar — pada layar yang berbeda.**
 *
 * ⛔ §24 MEMPERSEMPIT ronde 69, tidak membatalkannya. Ronde 69 tetap berlaku
 * penuh pada ketiga layar yang ia sebut sendiri:
 *
 *   1. Laporan Realisasi (view existing)
 *   2. Limit Treaty In
 *   3. grid XOL
 *
 * Nol di antaranya form Treaty In, dan form inilah yang §24 atur.
 *
 * ⚠️ Bentuk BERPASANGAN dipakai HANYA untuk kolom yang desimalnya terbaca di
 * gambar. Kolom yang tidak terbaca tetap memakai bentuk teks biasa — yaitu
 * aturan lama, nol di ekor dibuang — dan didaftarkan di
 * `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §13. Menebaknya akan mengubah setiap
 * angka di setiap grid sekaligus.
 *
 * ⛔ `format.ts` TIDAK disentuh. `formatNumber` tetap membuang nol di ekor;
 * yang memadankan kembali `selAngka` di lapis modul, dan hanya untuk layar
 * ini.
 */
export type JenisAngka = GolonganAngka | readonly [GolonganAngka, number]

/** Golongan satu kolom, apa pun bentuknya. */
export function golongan(j: JenisAngka): GolonganAngka {
  return typeof j === 'string' ? j : j[0]
}

/**
 * Jumlah desimal yang DIPADANKAN, atau `null` bila kolomnya tidak terbaca di
 * gambar mana pun.
 *
 * ⚠️ `null` BUKAN nol. Nol berarti "padankan sampai nol desimal" (kolom
 * `Amount in IDR` gambar 29); `null` berarti "tidak diketahui, pakai aturan
 * lama".
 */
export function desimalPadan(j: JenisAngka): number | null {
  return typeof j === 'string' ? null : j[1]
}

/**
 * ⛔ `persenShare` BUKAN `persen`, dan bedanya SUDAH diputuskan —
 * **8 desimal**, keputusan pemilik proses 4 Oktober 2026
 * (`KEPUTUSAN-PENYELARASAN-REPO.md` §13).
 *
 * Kedua ujungnya ditolak, masing-masing dengan sebabnya:
 *
 *   2 desimal — `PctTotal` `99.999999999999900` akan tampil `100%`, menutupi
 *     persis selisih yang orang cari ketika memeriksa. Dan `SD-05`/`BR-01`
 *     benar: tiga share 33,333 berjumlah TEPAT 100, tiga share 33,33 tidak.
 *   tanpa batas — `2,825601535925207120348922139444%` (30 desimal, nyata di
 *     satu kontrak) tidak terbaca di dalam sel grid, dan penyimpanannya
 *     hanya 8 desimal (`NUMBER(38,8)`, dijaga `TestNolNumberTanpaPresisi`).
 *     Menampilkan 30 berarti mengaku lebih teliti daripada yang disimpan.
 *
 * Delapan menyelesaikan keduanya: sama persis dengan batas penyimpanan, dan
 * untuk share seberapa pun realistis sifat jumlah-tepat-100 tetap terjaga.
 *
 * ⚠️ `DESIMAL_PERSEN` (2) tetap berlaku untuk persen yang BUKAN share —
 * `MDP_RATIO`, `ADJ_RATE`, `ROL`. Ketiganya tidak dijumlahkan menjadi 100,
 * dan dua di antaranya memang melampaui 100.
 */
export const DESIMAL_UANG = 4
export const DESIMAL_PERSEN = 2
export const DESIMAL_PERSEN_SHARE = 8

/** Kolom tab Limits — dari `M_TREATY_IN2`, satu baris per layer. */
export const KOLOM_LIMITS = [
  'LAYER', 'LAYERTYPE', 'BASIS_COVER', 'TREATYTYPE', 'CURRENCY',
  'LIMIT_100', 'CEDANT_RETENTION', 'MDP', 'MDP_RATIO', 'ROL',
  'ADJ_RATE', 'EARN_PREMIUM', 'CURRENCYRELATION', 'CURRENCYLIMIT',
] as const

export const JENIS_LIMITS: readonly JenisAngka[] = [
  'teks', 'teks', 'teks', 'teks', 'teks',
  'uang', 'uang', 'uang', 'persen', 'persen',
  'persen', 'uang', 'teks', 'teks',
]

/**
 * Tab Limits kini dua komponen yang membaca ekspor langsung:
 * `labelsLimitsProp.ts` (Prop: Kind of Treaty → Treaty Type → Treaty Group
 * → DetailLimits) dan `labelsLimitsNP.ts` (Non-Prop: grid layer → Layers →
 * Summary of Limit → Total All Layers). Pohon lama `LIMITS_POHON` dicabut
 * 6 Oktober 2026 bersama `PohonLimits`.
 */
/**
 * Tombol `Add` di baris kepala grid — dan daftar tab yang BERHAK punya.
 *
 * ⛔ Keluhan pemilik proses 5 Oktober 2026: seluruh layar hanya punya SATU
 * tombol `Add` (Rate of Exchange), sementara gambar desain memperlihatkan
 * tiap tab ber-grid punya miliknya sendiri.
 *
 * ⚠️ Daftarnya DARI SAPUAN, bukan dari asumsi "tiap grid punya". Sapuan ke-43
 * gambar menemukan tiga grid yang JUSTRU TIDAK punya tombol `Add`, dan
 * memasangkannya di sana akan membuat layar berbeda dari layar lama:
 *
 *   Portfolio          gambar 02 — nol tombol, hanya "No items"
 *   Co-Ins Scale       gambar 18 — nol tombol
 *   Maximum Retention  gambar 26/27 — nol `Add`; yang ada `Update Total`
 *
 * ⛔ Tombolnya MATI. Modul ini nol jalur tulis.
 */
export const GRID_TAMBAH = {
  tambah: 'Add', // pyButtonLabel, gambar 19/29
  hapus: 'Delete', // per baris, gambar 19/29
  petunjuk: '',
} as const

/**
 * Tab yang gambarnya MEMPERLIHATKAN tombol `Add` di kepala gridnya.
 *
 * ⭐ Berkunci nama tab supaya tab berikutnya yang dibangun tidak lupa — uji
 * `desain-pega.test.ts` menuntut tiap nama di sini punya tombolnya di layar.
 *
 * ⚠️ `Limits` TIDAK di sini walau punya tombol: pohonnya memasang tombolnya
 * sendiri (`LIMITS_NP.tambahLayer` / `.tambahGrup`, `LIMITS_PROP.tambah`), dua tingkat dan
 * berbeda nama per cabang. Satu tanda boolean tidak dapat menyatakannya.
 *
 * ⛔ `Portfolio` DITAMBAHKAN 6 Oktober 2026, dan daftar ini pernah
 * MENYATAKAN SEBALIKNYA. Sebabnya daftar ini dibaca dari GAMBAR, dan
 * gambar `02` adalah tangkapan layar mode-BACA — di sana tombolnya memang
 * tidak terlihat. Ekspor menerangkan mengapa: sel `Add` (110) dan `Delete`
 * (115) ber-`pyVisible` = `OTHER` dengan
 * `pyCondition` = `TreatyIn.IsEditData!='1'`. Keduanya ADA; yang
 * menyembunyikannya keadaan layarnya, bukan ketiadaannya.
 */
export const GRID_BERTOMBOL_TAMBAH = ['EGNPI', 'Accumulation', 'Portfolio'] as const

/**
 * ⚠️ Yang gambarnya memperlihatkan tombol tetapi GRIDNYA belum dibangun —
 * didaftarkan supaya tidak terbaca sebagai sudah selesai.
 *
 *   Share non-prop   gambar 34 — grid `Reinsurer Name · Layer · % Share`
 *                    ber-`Add`/`Delete`. Grid itu BUKAN proyeksi baris layer
 *                    yang kita render; ia daftar reasuradur tersendiri, dan
 *                    sumbernya belum terpetakan.
 *   Installment      gambar 38 — tombol `Update Value` di samping medan
 *                    `Installment`. Ia memicu penghitungan ulang jadwal,
 *                    dan rumusnya belum dibaca dari Activity.
 *   11 sub-tab Limits gambar 05–15 — tiap daftar mata uang punya
 *                    `Add`/`Remove`; sub-tabnya sendiri belum dibangun.
 */
export const TOMBOL_TAMBAH_BELUM_BERGRID = [
  'Share non-prop · Reinsurer',
  'Installment · Update Value',
  'Limits · 11 sub-tab',
] as const

/** Kolom tab Share — dari `M_TREATY_IN2`. */
export const KOLOM_SHARE = [
  'LAYER', 'CESSIONPCT', 'SPREADINGTYPE', 'BROKERAGEPERCENTP',
  'CESSION_TO_RI', 'QSOR', 'QSRI', 'LIABILITYQSOR', 'LIABILITYQSRI',
  'EPI100', 'RIOGR',
] as const

export const JENIS_SHARE: readonly JenisAngka[] = [
  'teks', 'persenShare', 'teks', 'persenShare',
  'uang', 'persenShare', 'persenShare', 'uang', 'uang',
  'uang', 'uang',
]

/**
 * Kolom tab Event Limits — ⚠️ SATU kolom, dan itu temuan.
 *
 * Ekspor `TreatyInTabsNonProportional.xml` wilayah tab Event Limits
 * (@223612–@490872) menyebut EMPAT batas: `TreatyIn.Earthquake`,
 * `TreatyIn.FloodJab`, `TreatyIn.FloodNation`, `TreatyIn.RSMDLimit`,
 * masing-masing dengan mata uangnya. Dari keempatnya, `M_TREATY_IN2` hanya
 * punya `EARTHQUAKE`. Ketiga sisanya ada di dokumen (`FloodJab`,
 * `FloodNation`, `RSMDLimit` di akar `JSONDATA`) dan TIDAK di tabel ini.
 */
/**
 * ⭐ KEEMPAT BATAS, dari gambar `28` — bukan satu.
 *
 * Bentuk sebelumnya tiga kolom (`LAYER`, `CURRENCY`, `EARTHQUAKE`), sebab
 * `M_TREATY_IN2` hanya punya `EARTHQUAKE`. Dokumen punya keempatnya, dan
 * layar lama menampilkan keempatnya — masing-masing dengan mata uangnya
 * sendiri, bukan satu mata uang untuk semua.
 */
/**
 * Tab **Event Limits** — label keempat barisnya, disalin dari gambar `28`.
 *
 * ⛔ Urutannya URUTAN GAMBAR, bukan abjad: RSMD · Earthquake · Flood
 * (Jabodetabek) · Flood (Nationwide). Mengurutkannya abjad membuat layar ini
 * berbeda dari layar yang orang hafal.
 */
export const EVENT_LIMITS = {
  judul: 'Event Limits', // gambar 28
  layer: 'Layer',
  // ⛔ RALAT 6 Oktober 2026: tab Event Limits NON-PROP mengikat properti
  // AKAR (`TreatyIn.RSMDLimit` …), bukan `Detail[].…`. Label sel ekspor
  // `pyLabelFieldValue`, sama persis dengan gambar 28.
  rsmd: 'RSMD Limit', // gambar 28 · TreatyIn.CurrencyRSMD / RSMDLimit
  gempa: 'Earthquake Limit', // gambar 28 · TreatyIn.CurrencyEarthquake / Earthquake
  banjirJab: 'Flood Limit (Jabodetabek)', // gambar 28 · TreatyIn.CurrencyFloodJab / FloodJab
  banjirNas: 'Flood Limit (Nationwide)', // gambar 28 · TreatyIn.CurrencyFloodNat / FloodNation
  /**
   * ⛔ BUKAN "tidak ada": nilai akar ini TERISI di 49 dari 772 kontrak
   * Non-Prop (sapuan seluruh korpus, 6 Oktober 2026).
   *
   * ⭐ 7 Oktober 2026: rumahnya kini DITETAPKAN — `T_TREATY_HAZARD_LIMIT`
   * (diagram v2 `TreatyIn [BATAS_BAHAYA]`, migrasi `446`). Yang tersisa dua:
   * tabelnya belum terpasang di basis data, dan nilai kontrak LAMA hanya ada
   * di dokumen JSON — `TREATYINDETAIL` tidak menulis nilai akar Non-Prop —
   * sedangkan JSON dilarang dibaca. Kontrak lama karena itu tetap kosong.
   */
  belumTerjangkau:
    'Nilai tersimpan belum dapat dibaca: tabelnya (T_TREATY_HAZARD_LIMIT) belum terpasang, dan nilai kontrak lama hanya ada di dokumen JSON yang tidak boleh dibaca.',
} as const

/**
 * ⚠️ DIPERTAHANKAN walau tabnya tidak lagi memakai grid.
 *
 * Tab Event Limits kini berbentuk empat baris berlabel (gambar 28), tetapi
 * daftar kolom ini tetap berdiri sebagai catatan kesembilan medan yang
 * pencabutan `M_TREATY_IN2` kembalikan — dan uji yang menjaganya masih
 * menanyakannya.
 */
export const KOLOM_EVENT_LIMITS = [
  'Layer', // Limits[].Layer
  'RSMD Limit', // gambar 28 · Detail[].RSMDLimit
  'RSMD Currency', // Detail[].CurrencyRSMD
  'Earthquake Limit', // gambar 28 · Detail[].Earthquake
  'Earthquake Currency', // Detail[].CurrencyEarthquake
  'Flood Limit (Jabodetabek)', // gambar 28 · Detail[].FloodJab
  'Flood (Jabodetabek) Currency', // Detail[].CurrencyFloodJab
  'Flood Limit (Nationwide)', // gambar 28 · Detail[].FloodNation
  'Flood (Nationwide) Currency', // Detail[].CurrencyFloodNat
] as const

export const JENIS_EVENT_LIMITS: readonly JenisAngka[] = [
  'teks',
  'uang', 'teks',
  'uang', 'teks',
  'uang', 'teks',
  'uang', 'teks',
]

/** Kolom tab RNM Share — dari `M_TREATY_IN2`. */
export const KOLOM_RNM_SHARE = [
  'LAYER', 'RNMSHARE', 'LIABILITY_RNM', 'MDP_RNM_100',
  'EPIRNMQS100', 'RNM_RETAINED_PREMI', 'RNM_QS_PREMI',
] as const

export const JENIS_RNM_SHARE: readonly JenisAngka[] = [
  'teks', 'persenShare', 'uang', 'uang',
  'uang', 'uang', 'uang',
]

/**
 * Kolom tab Co-Ins Scale — tabel pendaratan kesembilan, migrasi 432.
 *
 * ⚠️ `Co-Insurance Share` BUKAN angka melainkan PITA: `>=30% up to < 50%`.
 * Golongannya `teks`, dan `formatNumber` memang mengembalikan teks
 * bukan-angka apa adanya — tetapi menggolongkannya `teks` membuat niat itu
 * terbaca alih-alih bergantung pada kebetulan.
 */
export const KOLOM_COIN_SCALE = [
  'Co-Insurance Share', // sel 264 · .CoInShare · pxTextInput
  '% Treaty Limit', // sel 265 · .PctLimit · pxNumber
] as const
/**
 * ⛔ DUA kolom, bukan empat — 6 Oktober 2026. Bentuk sebelumnya menambahkan
 * `Disusun oleh` dan `Disusun pada` (jejak audit Pega pada elemennya). Grid
 * `TreatyIn.CoInScale` di ekspor hanya punya `.CoInShare` dan `.PctLimit`,
 * dan gambar `18` — tangkapan layar Pega yang berjalan — memperlihatkan dua
 * kolom itu saja. Datanya TIDAK dibuang: ia tetap di model
 * (`skalaKoasuransi[].penyusun`/`.disusunPada`), hanya tidak dirender.
 */
export const JENIS_COIN_SCALE: readonly JenisAngka[] = ['teks', 'persenShare']

/**
 * Tab Co-Ins Scale — grid DAN dua medan, 6 Oktober 2026.
 *
 * ---------------------------------------------------------------------
 * ⛔⛔ RALAT DI RONDE YANG SAMA: GRIDNYA HIDUP, BUKAN MATI
 * ---------------------------------------------------------------------
 * Bunyi sebelumnya — dan laporan pencocokan §5 temuan 2 — menyatakan grid
 * `TreatyIn.CoInScale` MATI, sebab wadahnya ber-`pyContainerVisibleWhen` =
 * `1=2`. ITU SALAH BACA. Wadah yang sama juga ber-`pyIsVisibilityOption` =
 * `ALWAYS`, dan pilihan itu MENIMPA syaratnya — persis seperti `pyVisible`
 * = `ALWAYS` menimpa `pyCondition` pada sel. Saksi penentunya gambar `18`
 * (tangkapan layar Pega yang berjalan): grid `Co-Insurance Share · % Treaty
 * Limit` TAMPIL, dengan kedua medan di bawahnya.
 *
 * ⚠️ Jadi `1=2` di XML belum berarti mati. Yang menentukan pasangan
 * pilihan-visibilitas + syaratnya: `OTHER` + `1=2` mati; `ALWAYS` + `1=2`
 * hidup, dan `1=2`-nya hanya sisa.
 *
 * Yang SUNGGUH berlaku, dari `Section/TreatyInTabsProportional.xml`:
 *
 *   grid `TreatyIn.CoInScale`    TAMPIL (wadah `pyIsVisibilityOption` ALWAYS)
 *     sel 262 tambah (`addRow`)    di sel KEPALA ketiga, `pyCondition` = `TreatyIn.IsEditData!='1'`
 *     sel 266 hapus (`deleteRow`)  per baris, `pyCondition` = `TreatyIn.IsEditData!='1'`
 *     sel 264/265 baca-saja bila `TreatyIn.IsEditData='1'`
 *   sel 277 `.MaxCoNonGroup`     TAMPIL, `pxNumber`, baca-saja bila `ViewState ='1'`
 *   sel 278 `.MaxCoGroup`        sama
 *
 * ⛔ Kedua tombol grid ber-`pyLabel` KOSONG — hanya ikon (`IconAdd.png`/
 * `IconTrash.png`), dengan `pyControlDisplayTitle` = `Button`. Itulah sebab
 * sapuan berdasarkan LABEL melewatkannya; sapuan harus berdasarkan API aksi
 * (`addRow`/`deleteRow`). Teks yang tampil di layar datang dari Pega yang
 * berjalan, bukan dari XML — lihat `tambah`/`hapus` di bawah.
 *
 * ⚠️ Label medan diambil dari `pyLabelFieldValue` sel itu sendiri, LENGKAP
 * dengan kata `Panel` — ringkasan "Max Co-Insurance Non Group" menghilangkannya.
 */
export const CO_INS_SCALE = {
  nonGroup: 'Max Co-Insurance Panel (Non Group)', // sel 277 · .MaxCoNonGroup
  group: 'Max Co-Insurance Panel (Group)', // sel 278 · .MaxCoGroup
  /**
   * Sel 262. `pyLabel` dan `pyCaption` KOSONG di XML — bandingkan Rate of
   * Exchange, yang tombolnya ber-`pyLabel` = `Add`. Teks `Tambah` diambil
   * dari tangkapan layar Pega pemakai 6 Oktober 2026 (mode ubah): Pega
   * mengisi tombol ikon tanpa teks dengan teks bawaannya sendiri.
   */
  tambah: 'Tambah',
  /**
   * Sel 266 — `pyLabel` kosong, ikon `IconTrash.png`. ⚠️ Belum ada tangkapan
   * layar Treaty In yang memperlihatkan barisnya. `Hapus` disamakan dengan
   * Pega yang SAMA di NB FacIn (`GRID_OBJEK`, tangkapan layar work owner
   * 3 Oktober 2026): pasangan `IconAdd`/`IconTrash` tanpa teks tampil sebagai
   * `Tambah`/`Hapus`.
   */
  hapus: 'Hapus',
  /**
   * `pyFieldValueForNoRows` = `GridNoResultsOnLoad` — field value bawaan
   * Pega, teksnya terlihat di gambar 18 dan di tangkapan layar pemakai.
   */
  kosong: 'No items',
  /** Sel 265 — `pySymbol` constant `%`, `pySymbolPosition` right. */
  simbolPersen: '%',
} as const

/**
 * ⛔ Kolom `M_TREATY_IN2` yang TIDAK dipakai satu tab pun — disebut namanya,
 * bukan dibuang diam-diam.
 *
 * Kesepuluhnya KEPALA KONTRAK yang berulang pada setiap baris layer, dan
 * form sudah menampilkan kesepuluhnya di bagian atas dari `TREATY_IN`:
 * `MASTERID` `PROPORTIONTYPE` `CEDINGID` `CEDING` `SOBID` `SOB`
 * `TREATYGROUP` `TREATYCONTRACTNAME` `COMMENCEMENT` `TERMINATION`.
 *
 * Menampilkannya lagi di dalam grid layer akan mengulang nilai yang sama
 * pada setiap baris tanpa menambah satu pun keterangan.
 */
/**
 * Ke-41 kolom `M_TREATY_IN2` apa adanya, dalam urutan katalog Oracle.
 *
 * ⛔ Ada di sini SUPAYA DAPAT DIJUMLAHKAN. Tanpa sebuah angka yang harus
 * pas, kolom ke-41 dapat menguap tanpa suara ketika seseorang menyusun
 * ulang sebuah tab — dan kolom yang hilang dari layar tidak menimbulkan
 * satu pun galat.
 */
export const KOLOM_IN2_SEMUA = [
  'MASTERID', 'PROPORTIONTYPE', 'CEDINGID', 'CEDING', 'SOBID', 'SOB',
  'TREATYGROUP', 'TREATYCONTRACTNAME', 'COMMENCEMENT', 'TERMINATION',
  'TREATYTYPE', 'CESSIONPCT', 'CEDANT_RETENTION', 'BASIS_COVER',
  'LAYERTYPE', 'LAYER', 'SPREADINGTYPE', 'CURRENCY', 'LIMIT_100',
  'ADJ_RATE', 'EARN_PREMIUM', 'MDP_RATIO', 'MDP', 'ROL',
  'CURRENCYRELATION', 'CESSION_TO_RI', 'EPI100', 'RIOGR',
  'BROKERAGEPERCENTP', 'EARTHQUAKE', 'RNMSHARE', 'CURRENCYLIMIT',
  'LIABILITY_RNM', 'MDP_RNM_100', 'QSOR', 'QSRI', 'LIABILITYQSRI',
  'LIABILITYQSOR', 'EPIRNMQS100', 'RNM_RETAINED_PREMI', 'RNM_QS_PREMI',
] as const

export const KOLOM_IN2_TIDAK_DIPAKAI = [
  'MASTERID', 'PROPORTIONTYPE', 'CEDINGID', 'CEDING', 'SOBID', 'SOB',
  'TREATYGROUP', 'TREATYCONTRACTNAME', 'COMMENCEMENT', 'TERMINATION',
] as const

// ===========================================================================
// PANEL ATTACHMENT + HISTORY
// Sumber: `Section/WorkAttachments.xml` (Category · Count · Upload file ·
// View File · Download All · Refresh · spanduk biru) dan
// `Section/ShowAttachmentTreaty.xml` (File Name · Type · Button).
//
// ⛔ Datanya dari `POOLDATA.M_ATTACHMENTTREATY_2` — tabel WARISAN, 43 baris.
// Nol tabel baru, nol migrasi.
// ===========================================================================

/** Kolom panel Attachment — `pyCaption` di `WorkAttachments.xml` apa adanya. */
export const KOLOM_LAMPIRAN = ['Category', 'Count', 'Upload file', 'View File'] as const

/**
 * ⛔ `Count` adalah CACAH, dan cacah TIDAK diformat.
 *
 * Pemisah ribuan pada cacah berkas membuat `1.234` terbaca seribu dua ratus
 * padahal ia seribu dua ratus — kebetulan benar di sini, tetapi aturannya
 * sama dengan `SubDays`: yang dihitung butir, bukan uang.
 */
export const JENIS_LAMPIRAN: readonly JenisAngka[] = ['teks', 'teks', 'teks', 'teks']

/** Kolom daftar berkas di dalam satu kategori. */
export const KOLOM_BERKAS_LAMPIRAN = ['File Name', 'Type', 'Uploaded', 'By'] as const
export const JENIS_BERKAS_LAMPIRAN: readonly JenisAngka[] = ['teks', 'teks', 'teks', 'teks']

/**
 * Kolom panel History — Date · PIC · Approval · Comment.
 *
 * ⛔ Sumbernya `T_VIEW_COMMENT`, tabel pendaratan yang SUDAH ada dan
 * SUDAH terisi 11.365 baris — keempat kolomnya berpadanan satu-satu dengan
 * `TANGGAL`, `OPERATORNAME`, `ISAPPROVED`, `SUGGEST`. Nol pembacaan baru,
 * nol tabel baru; `repository.BacaCatatan` dipakai ulang apa adanya.
 */
export const KOLOM_HISTORY = ['Date', 'PIC', 'Approval', 'Comment'] as const
export const JENIS_HISTORY: readonly JenisAngka[] = ['teks', 'teks', 'teks', 'teks']

export const LAMPIRAN = {
  judul: 'Attachment',
  judulHistory: 'History',
  /** `pyCaption Recommended safe substitute should be . or _` — disalin apa adanya. */
  spanduk: 'Recommended safe substitute should be . or _',
  unduhSemua: 'Download All', // pyButtonLabel DOWNLOAD ALL
  segarkan: 'Refresh', // pyButtonLabel REFRESH
  unggah: 'Upload', // kolom `Upload file`, gambar 24/42
  lihatBerkas: 'View', // kolom `View File`, gambar 24/42
  tutup: 'Close', // tombol merah kaki layar, gambar 24/42/43
  /**
   * ⭐ PERMINTAAN PERUBAHAN — judul modal `View File`.
   *
   * Pemilik proses, keterangan gambar `25`: *"dan view upload saran
   * dibuatkan pop up dan ini diubah menjadi design nya bagus"*. Di Pega ia
   * jendela Chrome terpisah berjudul `ShowAttachmentTreaty`; judul teknis
   * itu TIDAK dibawa ke modal — yang dibawa nama panelnya.
   */
  judulLihatBerkas: 'View File',
  tanpaBaris: 'No items',
  /** `Kosong` panel History — teks layar lama apa adanya. */
  tanpaRiwayat: 'No items',
  tanpaLampiran: 'No items',
  petunjukLampiran: '',
  petunjukHistory: '',
  /**
   * ⛔ Penanda kategori yang pasangan kode↔namanya BELUM dipastikan.
   *
   * Empat kode — 00003, 00004, 00008, 00009 — punya nol baris, dan keempat
   * namanya tidak ada di korpus kedua modul. Menebak pasangannya menaruh
   * berkas di kategori yang salah, dan itu baru ketahuan bertahun kemudian.
   */
  kategoriBelumPasti: 'nama kategori belum dipastikan',
  /** FlowAction `TreatyAttachContent` — `pyCaption ASM Attach Content`. */
  judulUnggah: 'ASM Attach Content',
  /** Tombol FlowAction — `pyButtonLabel Attach` / `Cancel`. */
  lampirkan: 'Attach',
  batal: 'Cancel',
  /** `pyAttachmentScreen` — pemilih berkas. */
  pilihBerkas: 'Select file(s)',
  /**
   * Kotak unggah — teks DISALIN dari Master Product Name Life (`LAIN_MPNL`),
   * permintaan pemakai 8 Oktober 2026: bentuk unggah disamakan.
   */
  seretBerkas: 'Drag and drop files here, or click to choose files',
  mengunggah: 'Uploading',
  buangPilihan: 'Remove',
  gagal: 'Failed',
  /** Kontrak baru belum ber-ID — lampiran menempel pada `TREATYID`. */
  simpanDulu: 'Simpan kontrak lebih dulu untuk mengunggah lampiran.',
  /** `ShowAttachmentTreaty` — `pyLabel` apa adanya. */
  viewOffice: 'View Office Online',
  hapus: 'Delete',
  gantiKategori: 'Change Category',
  simpanKategori: 'Save',
  /**
   * ⛔ KEDUA kolom modal `View File` — gambar `25` memperlihatkan tepat dua.
   *
   * Panel induknya punya empat (`File Name` · `Type` · `Uploaded` · `By`);
   * jendela `ShowAttachmentTreaty` punya DUA. Menyamakan keduanya membuat
   * modal ini bukan modal yang gambar 25 perlihatkan.
   */
  namaBelumBerumah:
    'Empat nama kategori ada di layar lama tetapi kodenya belum dipastikan, jadi keempatnya belum ditampilkan sebagai nama: Binding, signed share Email · Claim Data · Info Pack · Letter of Acknowledgment / LOA. Lihat docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md.',
} as const

/** Kolom modal `View File` — DUA, gambar 25. */
export const KOLOM_LIHAT_BERKAS = ['File Name', 'Type'] as const

/**
 * Kesebelas nama kategori lampiran, dibaca dari gambar `24` (prop) dan `42`
 * (non-prop).
 *
 * ⛔ Daftar ini TIDAK menutup `PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`.
 * Yang terbuka adalah pasangan KODE↔NAMA, dan gambar hanya memberi namanya;
 * urutannya di layar alfabetis, dan ronde 4 Oktober 2026 melarang keras
 * menyimpulkan kode dari urutan abjad.
 *
 * ⭐ Yang BARU dari gambar: daftarnya BERCABANG. Sepuluh nama sama persis;
 * butir kesepuluh berbeda — `Pega Proportional Calculation` di cabang
 * proporsional, `Pega Non Proportional Calculation` di non-proporsional.
 */
export const NAMA_KATEGORI_LAMPIRAN_PROP = [
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
] as const

export const NAMA_KATEGORI_LAMPIRAN_NON_PROP = [
  'Analysed Email',
  'Approval Email',
  'Assessment Inward Treaty Form / Format Analisa Treaty',
  'Binding, signed share Email',
  'Claim Data',
  'Info Pack',
  'Letter of Acknowledgment / LOA',
  'Offer Email',
  'Others',
  'Pega Non Proportional Calculation /Perhitungan Pega Non Proportional',
  'Summary Treaty Leader',
] as const
