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
  kosongPetunjuk: 'Pemindahan isinya dari sistem lama adalah tiket 44.',
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
  kosongPetunjuk:
    'Tabel warisan berisi 1.854 kontrak. Halaman ini kosong karena nomornya di luar jangkauan, atau karena penyaring tidak menemukan padanan.',
  keterangan:
    'Kesembilan kolom dan susunan tombolnya disalin dari Section/InputTreatyInOffer.xml ekspor 2026-09. Isinya dibaca dari tabel warisan POOLDATA.TREATY_IN — 1.854 kontrak, urut pengenal menurun.',
  /** Sumber datanya, dinyatakan di layar — bukan hanya di komentar kode. */
  catatanSumber:
    'Baris di atas dibaca dari POOLDATA.TREATY_IN, tabel sistem lama — BACA SAJA. Kontrak yang sudah dipindahkan ke model baru (KONTRAK/VERSI_KONTRAK) adalah daftar yang berbeda, dan tiket 59 yang memindahkannya.',
  /** Kolom yang model baru BELUM punya rumahnya — dinyatakan, bukan diisi. */
  // ⚠️ RALAT 3 Oktober 2026 — kalimat lamanya menyangkal kolom yang ADA.
  // "Position To" memang belum punya rumah di MODEL BARU, dan itu tetap
  // benar; tetapi layar ini membaca TABEL WARISAN, dan di sana
  // `POSITIONUSERNAME` ada. Dari 1.854 baris, 32 terisi dan 1.822 NULL —
  // jadi selnya kosong karena nilainya memang kosong, bukan karena kolomnya
  // tidak ada.
  /** Judul lipatan keterangan kaki — pendek, dan menyebut isinya. */
  kakiJudul: 'Tentang data di tabel ini',
  catatanPosisiKe:
    'Kolom "Position To" terisi pada 32 dari 1.854 baris warisan; sel kosong berarti nilainya memang kosong. Di model baru kolom ini belum ada — ia lahir bersama tiket 45 dan 49–53.',
  // ⚠️ RALAT 3 Oktober 2026 — di tabel warisan keduanya NAMA, bukan pengenal.
  // `CEDING` dan `LEADINGREINSSOURCE` memuat teksnya langsung. Penyangkalan
  // lama berlaku untuk model baru, yang memang hanya menyimpan pengenal.
  catatanNama:
    'Di tabel warisan "Ceding" dan "Source of Business" berupa nama. Di model baru keduanya pengenal — ERD.md §2.8 menempatkan kedua tabel acuan itu DI LUAR skema ini.',
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
  kursPetunjuk: 'Kurs per mata uang adalah tiket 20 (MATA_UANG_KONTRAK); barisnya tidak dikarang.',
  /**
   * Keterangan KAKI, bukan isi layar.
   *
   * ⛔ Dulu paragraf di tengah form, di antara medan dan grid kurs. Prosa
   * pengembang yang duduk di jalur baca pemakai membuat layar terbaca sebagai
   * catatan rilis; ia turun ke kaki, redup, satu baris.
   */
  catatanPilihLuar:
    'Tombol "Choose Ceding" dan "Choose Source of Business" menunggu modul pemilik tabel acuannya — ERD.md §2.8 menempatkan CEDANT dan ASAL_BISNIS di luar skema Treaty In.',

  /**
   * Keterangan medan yang KUNCINYA tidak ada di dokumen warisan.
   *
   * ⛔ Medan MATI dengan keterangan, bukan kotak kosong. Kotak kosong
   * terbaca "belum diisi"; medan mati terbaca "tidak ada di sistem lama".
   * Pola yang sama sudah dipakai tombol `Choose Ceding`.
   *
   * Sapuan 3 Oktober 2026 atas 1.854 dokumen: `ContractRefNo` ada di 742,
   * `TreatyLeader` di 659, `BordereauxNote` di 1.018. Jadi medan ini mati
   * pada SEBAGIAN kontrak, bukan pada semuanya.
   */
  takAdaDiWarisan: 'Tidak ada di dokumen sistem lama',
  /** Kolom tab Portfolio - dari kunci `Portfolio` di dokumen warisan. */
  petunjukPortofolio:
    'Portofolio dibaca dari dokumen warisan kontrak ini; terisi pada 152 dari 300 dokumen yang disapu. Kosong berarti kontrak ini memang tidak punya.',
  petunjukAkumulasi:
    'Akumulasi dibaca dari dokumen warisan kontrak ini; terisi pada 12 dari 300 dokumen yang disapu — jarang, dan kosong berarti kontrak ini memang tidak punya.',
  petunjukPeriode:
    'Periode pelaporan dibaca dari dokumen warisan kontrak ini; terisi pada 225 dari 300 dokumen yang disapu.',

  /** Petunjuk umum grid tab kosong — menyebut sebabnya, bukan cacahnya. */
  petunjukTabel:
    'Dibaca dari tabel pendaratan kontrak ini. Kosong berarti dokumen warisannya memang tidak punya baris untuk tab ini.',

  /**
   * ⛔ Petunjuk keempat tab `M_TREATY_IN2` BERBEDA, dan bedanya diukur.
   *
   * Tujuh tab lain dibaca dari tabel pendaratan yang dimuat dari dokumen
   * kontrak itu sendiri — kosong di sana berarti dokumennya memang tidak
   * punya. Keempat tab ini dibaca dari `M_TREATY_IN2`, tabel warisan yang
   * mencakup **1.340 dari 1.854 kontrak**. Terukur 3 Oktober 2026: 510
   * kontrak punya `Limits[]` berisi di dokumennya (1.210 elemen) tetapi nol
   * baris di tabel itu.
   *
   * Jadi kosong di sini TIDAK boleh berbunyi "kontrak ini memang tidak
   * punya" — pada 510 kontrak kalimat itu KELIRU, dan ia akan membuat orang
   * berhenti mencari data yang sebenarnya ada.
   */
  petunjukLayer:
    'Dibaca dari tabel layer sistem lama (M_TREATY_IN2), yang mencakup 1.340 dari 1.854 kontrak. Kosong di sini dapat berarti kontrak ini tidak punya layer, ATAU tabel itu tidak mencakupnya — 510 kontrak punya limit di dokumennya tanpa satu baris pun di sana.',
  /**
   * ⚠️ Event Limits dapat kosong karena sebab KETIGA: dari empat batas yang
   * ekspor Section sebut (`Earthquake`, `FloodJab`, `FloodNation`,
   * `RSMDLimit`), `M_TREATY_IN2` hanya punya `EARTHQUAKE` — dan kolom itu
   * terisi pada 759 dari 7.281 baris.
   */
  petunjukEventLimits:
    'Dibaca dari kolom EARTHQUAKE di M_TREATY_IN2, terisi pada 759 dari 7.281 baris layer. Tiga batas lain yang layar lama tampilkan (Flood Jabodetabek, Flood Nationwide, RSMD) tidak ada di tabel itu dan belum punya sumber tabel.',
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
  petunjukTeksPengecualian:
    'Dibaca dari kunci ExclusionsP (proporsional) atau Exclusions (non-proporsional) di dokumen warisan. Kosong berarti kunci yang sesuai cabang kontrak ini tidak ada — dan ejaan cabang seberang sengaja TIDAK dipakai sebagai pengganti, sebab isinya teks yang berbeda.',
  petunjukTeksSyarat:
    'Dibaca dari kunci SpecialConditionsP (proporsional) atau SpecialConditions (non-proporsional) di dokumen warisan. Ejaan ketiga SpecialConditionsp dipakai kedua cabang, jadi ia tidak pernah terpilih — bila berisi, ia disebut sebagai kunci lain. Kosong berarti kunci yang sesuai cabang tidak ada.',
  petunjukSkalaKoasuransi:
    'Skala koasuransi dibaca dari tabel pendaratan kontrak ini; terisi pada 186 dari 1.854 kontrak. Kosong berarti kontrak ini memang tidak punya skala.',

  belumDibangun: 'Tab ini belum dibangun.',
  /**
   * ⛔ Pesan ini pernah menyebut "hanya tab Reporting Period", dan menjadi
   * BASI diam-diam begitu tab kedua mendarat. Kalimat yang menyebut cacah
   * atau nama tab akan selalu basi; yang ini menyebut SEBABNYA, dan sebab
   * tidak berubah tiap ronde.
   */
  belumDibangunPetunjuk:
    'Sumber datanya belum ditemukan di ekspor sistem lama. Tabnya tetap tampil supaya susunan strip tidak berubah diam-diam ketika isinya menyusul.',
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
  'Retro', // @1031030
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
  'RNM Share', // @2093710 — panel "Summarry of RNM Share" @2324385 (ejaan ekspor), "Total All Layers RNM Share" @2486274
  'Retro', // @3291797
  'Installment', // @3465901
  'Value Difference', // @3848946
  'Exclusions', // @3907436
  'Special Conditions', // @3942540
  'Information & Submit', // @3977799
] as const

export type TabProporsional = (typeof TAB_PROPORSIONAL)[number]
export type TabNonProporsional = (typeof TAB_NON_PROPORSIONAL)[number]

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
} as const

/** Kolom tab Portfolio — nama kuncinya di dokumen warisan apa adanya. */
export const KOLOM_PORTOFOLIO = ['Type', 'Type Portfolio', 'Description'] as const

/** Kolom tab Accumulation — idem. */
export const KOLOM_AKUMULASI = ['Period', 'Report Date', 'Sub Days', 'Sub Due Date'] as const

/** Kolom tab EGNPI — nama kunci di dokumen warisan apa adanya. */
export const KOLOM_EGNPI = [
  'Amount', 'Amount IDR', 'As Date', 'Class of Business',
  'Currency', 'Note', 'Proportion', 'Treaty Group',
] as const

/** Kolom tab Maximum Retention. */
export const KOLOM_RETENSI = [
  'Amount', 'Class of Business', 'Currency', 'Note', 'Treaty Group',
] as const

/** Kolom JADWAL tab Installment. */
export const KOLOM_ANGSURAN = [
  'Installment', 'Currency', 'Amount', 'Pct', 'Due Date', 'Payment Date', 'WPC',
] as const

/** Kolom tab Information & Submit. */
export const KOLOM_CATATAN = ['Date', 'Operator', 'Approved', 'Suggest'] as const

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
export type JenisAngka = 'uang' | 'persen' | 'persenShare' | 'teks'

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
export const KOLOM_EVENT_LIMITS = ['LAYER', 'CURRENCY', 'EARTHQUAKE'] as const
export const JENIS_EVENT_LIMITS: readonly JenisAngka[] = ['teks', 'teks', 'uang']

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
  'Co-Insurance Share', '% Treaty Limit', 'Disusun oleh', 'Disusun pada',
] as const
export const JENIS_COIN_SCALE: readonly JenisAngka[] = ['teks', 'persenShare', 'teks', 'teks']

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
