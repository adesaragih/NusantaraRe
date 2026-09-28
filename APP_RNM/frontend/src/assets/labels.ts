// PUSAT LABEL UI — F0.1, brief lanjutan 6 §2 aturan 4.
//
// # Aturan yang berkas ini tegakkan
//
// SEGALA yang soal TAMPILAN diubah di FRONTEND, dan sebisanya di BERKAS INI:
// "mau ganti sebutan" tidak pernah berarti membuka backend.
//
// ⛔ SETIAP TEKS MEMBAWA BUKTI XML-nya — path dan nomor baris. Teks tanpa
// bukti adalah teks yang dikarang, dan sistem lama akan berbeda sebutan dari
// sistem baru tanpa satu pun orang menyadarinya. Bila sebuah teks memang
// tidak ada di korpus, ia ditandai `[tidak ada di korpus]` dengan terang,
// bukan diam-diam dikarang.
//
// ⚠️ Bahasa Indonesia boleh berdiri DI SAMPING, bukan menggantikan. Judul
// tahap dan label tombol adalah kosakata yang dipakai pengguna sejak sistem
// lama; menerjemahkannya diam-diam membuat pelatihan, dokumen, dan percakapan
// sehari-hari tidak lagi cocok dengan layar.

/**
 * Menu sidebar — HANYA yang berbukti korpus (brief lanjutan 7 §1.2).
 *
 * ⛔ Aturan 7 brief 6 DICABUT. Sidebar tidak memuat sebelas modul lain, tidak
 * memuat butir `BelumTersedia`, dan tidak memuat *Dokumen / Komite / Detail &
 * Tutup / Cari Polis / Medical Check / Claim Analis* sebagai menu. Semua itu
 * dibuka DARI DALAM kasus lewat flow action dan popup — begitu Pega
 * melakukannya, dan menu yang tidak ada di sistem lama adalah menu yang
 * dikarang.
 *
 * ⛔ `Claim Life` sebagai nama kelompok: `[terverifikasi]`
 * `Flow/Register_Flow.xml:270` `<pyWorkTypeName>ClaimLife</pyWorkTypeName>`.
 *
 * ⚠️ `inbox` `[tidak ada di korpus]` sebagai label. Yang ada hanyalah
 * KOSAKATA-nya: berkas struktur pengekspor bernama `Struktur_InboxClaimLife`
 * dan report `InboxPremiumList`. Tidak ada harness portal Claim Life di
 * ekspor — modul yang punya menu portal diekspor bersama harness portalnya
 * (NB Treaty In `SFAPortalOpportunities`), dan Claim Life tidak.
 * `[terbuka — pemilik ekspor Pega]` apakah portalnya ada; bila jawabannya
 * datang, menu MENGIKUTI XML-nya.
 */
export const MENU = {
  /** Nama kelompok — `Register_Flow.xml:270` `pyWorkTypeName`. */
  kelompokClaimLife: 'Claim Life',
  /** `[tidak ada di korpus]` — kosakata `InboxPremiumList` / `Struktur_InboxClaimLife`. */
  inbox: 'Inbox Claim Life',
  /** VERBATIM `Flow/Register_Flow.xml:155` `<pyLabel>Register</pyLabel>`. */
  register: 'Register',
} as const

/**
 * Nama modul lain yang TIDAK BOLEH muncul di sidebar.
 *
 * ⛔ Daftar ini dipakai penjaga, bukan tampilan. Ia ada supaya "menambah satu
 * menu saja" untuk modul yang belum berbukti menjadi kegagalan uji, bukan
 * keputusan sepi yang tidak ada yang tinjau.
 */
export const MODUL_LAIN_TERLARANG = [
  'Treaty',
  'Endorsement',
  'Fac',
  'Prop',
  'Master',
  'Premium',
] as const

/**
 * Nama ketujuh belas kelompok sidebar — butir **bg**.
 *
 * ⛔ `[nama folder korpus]`, bukan kosakata kami. Setiap nama di bawah
 * adalah nama folder di `D:\XML\RNM_BRD\` apa adanya, termasuk ejaannya
 * yang janggal: `Endorsment Fac In` *(tanpa `e`)* dan `Komite Claim FacIn`
 * *(tanpa spasi)*. Merapikannya berarti sidebar menyebut modul dengan nama
 * yang tidak cocok dengan korpusnya, dan orang yang mencari foldernya tidak
 * akan menemukannya.
 *
 * ⛔ Kelompok yang BELUM dimigrasi tetap berdiri, terlipat dan tanpa butir.
 * Menyembunyikannya membuat aplikasi tampak lengkap padahal empat belas
 * modul belum ada — dan layar yang tampak lengkap padahal tidak adalah
 * layar yang tidak akan dicari lagi (pelajaran butir av).
 *
 * ⚠️ Kelompok `Master`, `Offer`, `Realization`, `Citrix`, dan `Borderaux`
 * milik REFERENSI_UI TIDAK dibawa: itu aplikasi Treaty, bukan korpus ini.
 */
export const MODUL = {
  claimFacIn: 'Claim Fac In',
  claimLife: 'Claim Life',
  claimNonProp: 'Claim Non Prop',
  claimProp: 'Claim Prop',
  edmTreatyIn: 'EDM Treaty In',
  endorsementLife: 'Endorsement Life',
  /** ⚠️ Ejaan folder apa adanya — `Endorsment`, tanpa huruf `e`. */
  endorsmentFacIn: 'Endorsment Fac In',
  /** ⚠️ Ejaan folder apa adanya — `FacIn`, tanpa spasi. */
  komiteClaimFacIn: 'Komite Claim FacIn',
  komiteClaimLife: 'Komite Claim Life',
  komiteClaimNonProp: 'Komite Claim Non Prop',
  komiteClaimProp: 'Komite Claim Prop',
  masterContractRetroLife: 'Master Contract Retro Life',
  masterProductNameLife: 'Master Product Name Life',
  nbFacIn: 'NB FacIn',
  nbTreatyIn: 'NB Treaty In',
  premiumListLife: 'PremiumList Life',
  rnwFacIn: 'RNW Fac In',
  /**
   * ⚠️ RALAT 28-09-2026 (sesi Treaty Contract Out): korpus memuat DUA PULUH
   * folder, bukan tujuh belas — `Treaty Contract Out`, `Treaty In`, dan
   * `Treaty In Adjustment` terlewat (`PROMPT-EKSEKUSI-HULU-HILIR.md` §8).
   * Yang ditambahkan di sini HANYA kelompok modul yang sesi itu bangun;
   * dua lainnya menunggu keputusan asisten/work owner.
   */
  treatyContractOut: 'Treaty Contract Out',
} as const

/**
 * Butir menu modul yang SUDAH berbukti XML — butir **bg**.
 *
 * ⛔ Empat butir, dan tidak lebih. Tiga modul punya bukti korpusnya:
 * Claim Life *(dua butir)*, PremiumList Life, Komite Claim Life. Empat
 * belas kelompok lain berdiri TANPA butir.
 */
export const MENU_MODUL = {
  /** Harness portal `PremiumLife_harness` → section `PremiumList`. */
  premiumList: 'PremiumList',
  /** `[tidak ada di korpus]` — worklist `KomiteRouter`, kosakata kami. */
  inboxKomite: 'Inbox Komite',
} as const

/** Kata untuk kelompok yang modulnya belum dipindahkan. */
export const KETERANGAN_BELUM_DIMIGRASI = 'belum dimigrasi'

/**
 * Label Beranda — butir **bg**.
 *
 * ⚠️ `[kerangka aplikasi, bukan menu Pega]`. Beranda pengganti layar awal
 * portal, sebagaimana `PremiumLife_harness` *(kelas `Data-Portal`)* menjadi
 * layar awal PremiumList. Claim Life **tidak punya** harness portal yang
 * terekspor — OQ ke pemilik ekspor tetap terbuka — jadi bentuk Beranda ini
 * keputusan kami, dan ditandai begitu.
 */
export const BERANDA = {
  judul: 'Beranda',
  salam: 'Selamat datang',
  aktif: 'aktif',
  antrean: 'antrean',
  tanpaAntrean: 'belum ada kotak masuk',
} as const
