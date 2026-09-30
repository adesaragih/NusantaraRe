// KATALOG DUA PULUH FOLDER KORPUS - lapisan aplikasi (perakit), bukan inti.
//
// Struktur tim satu folder per modul (30-09-2026): dulu `MODUL` di
// `inti/frontend/labels.ts`. Nama kelompok modul yang SUDAH dimigrasi kini
// tinggal di `menu.ts` modulnya sendiri (`KELOMPOK_<NAMA>`); katalog ini
// hanya kerangka yang menampilkan SETIAP folder korpus - termasuk yang belum
// dimigrasi - di Beranda. Isinya tidak berubah ketika modul bertambah (korpus
// tetap dua puluh folder), jadi tidak ada pengembang modul yang perlu
// menyuntingnya; `daftar.menuTabel.test.ts` menjaganya sama dengan isi awal
// M_NAV_MENU, dan menjaga setiap `KELOMPOK_<NAMA>` ada di sini.

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
export const FOLDER_KORPUS = {
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
  /**
   * Dua folder terakhir — ditambahkan brief menu `M_NAV_MENU` 30-09-2026:
   * isi awal tabel menu memuat DUA PULUH kelompok, satu per folder korpus,
   * dan LABEL-nya harus ada di sini (`frontend/daftar.menuTabel.test.ts`).
   */
  treatyIn: 'Treaty In',
  treatyInAdjustment: 'Treaty In Adjustment',
} as const
