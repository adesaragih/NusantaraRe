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
 * Nama kedua puluh kelompok sidebar — butir **bg** (dulu tujuh belas; ralat di bawah).
 *
 * ⛔ `[nama folder korpus]`, bukan kosakata kami. Setiap nama di bawah
 * adalah nama folder di `D:\XML\RNM_BRD\` apa adanya, termasuk ejaannya
 * yang janggal: `Endorsment Fac In` *(tanpa `e`)* dan `Komite Claim FacIn`
 * *(tanpa spasi)*. Merapikannya berarti sidebar menyebut modul dengan nama
 * yang tidak cocok dengan korpusnya, dan orang yang mencari foldernya tidak
 * akan menemukannya.
 *
 * ⛔ Kelompok yang BELUM dimigrasi tetap berdiri, terlipat dan tanpa butir.
 * Menyembunyikannya membuat aplikasi tampak lengkap padahal enam belas
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
   * dua lainnya menunggu keputusan asisten/work owner - keduanya masuk
   * 30-09-2026 (di bawah).
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

/**
 * Nama TAMPILAN modul yang BUKAN nama folder korpus - keputusan work owner 03-10-2026 ("ganti nama modul ... hapus kata
 * Master nya", nama tampilan saja; kode modul, folder, rute, dan MODUL_AKTIF tetap). Dipakai `M_NAV_MENU.LABEL` sesudah
 * slot menu modulnya (959 / 961), `kelompok` menu modulnya, dan kartu Beranda. Pasangan Go: `labelTampilDisetujui`
 * (`inti/backend/penjaga/menu_test.go`) - keduanya dikunci uji.
 */
export const LABEL_TAMPIL: Readonly<Partial<Record<keyof typeof FOLDER_KORPUS, string>>> = {
  masterContractRetroLife: 'Contract Retro Life',
  masterProductNameLife: 'Product Name Life',
}

/**
 * Modul DI LUAR dua puluh folder korpus (`PANDUAN-TIM-PER-MODUL.md` bab 5) — nama menunya = `M_NAV_MENU.LABEL` barisnya,
 * yang dibuat langkah inti tersendiri (bukan isi awal 900). Pasangan Go: `modulLuarKorpus`
 * (`inti/backend/penjaga/menu_test.go`) - keduanya dikunci uji.
 *
 * Keputusan work owner 03-10-2026: modul `marketingofficer` (insert/update `MARKETINGOFFICER`), label "Marketing
 * Officer", kelompok MASTER, baris menu migrasi inti 906. Keputusan work owner 04-10-2026: modul `companydetail`
 * (layar Pega SFAGIS Company Detail, tabel datar CLIENT), label "Company Detail", kelompok MASTER, migrasi inti 907.
 * Keputusan work owner 04-10-2026: modul `accounts` (layar Pega SFAGIS Account, tabel T_M_ACCOUNT), label "Accounts",
 * kelompok MASTER, migrasi inti 908. Keputusan work owner 04-10-2026: "Modul baru 'masterdata'" (menu master insert /
 * update / aktif / nonaktif), label "Master Data", migrasi inti 911 - lalu DIPECAH (keputusan work owner 04-10-2026:
 * "bukan di satuin begini, di pisah per sub modul", "8 modul terpisah"): delapan modul `master<nama>` satu menu per
 * master di kelompok MASTER, baris menu migrasi inti 912-919; baris 'masterdata' dihapus 920.
 */
export const MODUL_LUAR_KORPUS = {
  marketingOfficer: 'Marketing Officer',
  companyDetail: 'Company Detail',
  accounts: 'Accounts',
  masterNation: 'Nation',
  masterProvince: 'Province',
  masterCity: 'City',
  masterDistrict: 'District',
  masterCzone: 'CZone',
  masterAccumulatedType: 'Accumulated Type',
  masterAccumulation: 'Accumulation',
  masterObjectItemType: 'Object Item Type',
} as const

/**
 * Nama menu setiap modul: folder korpus dengan nama tampilan bila diputuskan (`LABEL_TAMPIL`), selain itu nama folder
 * VERBATIM; lalu modul di luar korpus (`MODUL_LUAR_KORPUS`).
 */
export const LABEL_MENU = {
  ...Object.fromEntries(
    Object.entries(FOLDER_KORPUS).map(([k, v]) => [k, LABEL_TAMPIL[k as keyof typeof FOLDER_KORPUS] ?? v]),
  ),
  ...MODUL_LUAR_KORPUS,
} as Readonly<Record<keyof typeof FOLDER_KORPUS | keyof typeof MODUL_LUAR_KORPUS, string>>
