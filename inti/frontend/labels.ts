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

// Menu DATAR (keputusan work owner 30-09-2026): tombol modul berlabel
// `M_NAV_MENU.LABEL`; label butir navigasi (`LABEL_MENU_PREMIUMLIST`,
// `LABEL_MENU_KOMITE`, `MENU_TCO.treatyContractOut`) dibuang. `MENU` di atas
// TETAP: `MENU.inbox` dan `MENU.register` judul dan tombol layar Inbox Claim
// Life, dan Shell memakai `MENU.inbox` sebagai judul cadangan.

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
  ringkasan: 'Antrean Claim Life',
  catatanTahap: 'Claim Life',
  judulModul: 'Modul',
  kolomModul: 'Modul',
  kolomStatus: 'Status',
  kolomAntrean: 'Antrean',
  kolomAksi: 'Aksi',
  judulPeran: 'Peran Anda',
} as const

/**
 * Label bingkai aplikasi — sidebar, topbar, menu profil.
 *
 * ⚠️ `[tidak ada di korpus]` `[kerangka aplikasi, bukan menu Pega]`. Tak satu
 * pun teks ini berasal dari sistem lama: ia bingkai aplikasi baru (desain
 * `workpage-template.html`, 29-09-2026), bukan layar yang dimigrasi. Karena
 * itu ia berbahasa Indonesia dan tidak membawa nomor baris XML.
 */
export const KERANGKA = {
  lewati: 'Lewati ke konten utama',
  menuSamping: 'Menu samping',
  navUtama: 'Navigasi utama',
  bukaMenu: 'Buka menu',
  tutupMenu: 'Tutup menu',
  ciutkanMenu: 'Ciutkan menu',
  perluasMenu: 'Perluas menu',
  cariMenu: 'Cari menu',
  pintasCari: 'Ctrl K',
  modeGelap: 'Mode gelap',
  modeTerang: 'Mode terang',
  profil: 'Profil',
  keBeranda: 'ke Beranda',
  modeStub: 'mode stub',
  /** `[tidak ada di korpus]` — login M_LOGIN_GO (keputusan work owner 01-10-2026). */
  keluar: 'Keluar',
  /** `[tidak ada di korpus]` */
  gantiSandi: 'Ganti sandi',
  /** `[tidak ada di korpus]` — `GET /api/menu` sedang dibaca (menu dari M_NAV_MENU). */
  memuatMenu: 'Memuat menu…',
  /** `[tidak ada di korpus]` — `GET /api/menu` menjawab, tetapi tak satu baris pun dapat tampil. */
  menuKosong: 'Menu kosong — M_NAV_MENU tidak memuat baris aktif yang dapat dibuka di sini.',
} as const

/**
 * Label halaman login dan ganti sandi — `M_LOGIN_GO`, keputusan work owner
 * 01-10-2026.
 *
 * ⚠️ `[tidak ada di korpus]` `[kerangka aplikasi]`: tata letak `loginbaru.html`
 * (lampiran work owner 01-10-2026); teksnya BAHASA INGGRIS - permintaan work
 * owner 01-10-2026: "bahasa saat login ubah ke Inggris". Berlaku untuk layar
 * login, ganti sandi, dan periksa sesi. Pesan gagal dipilih dari STATUS
 * jawaban (dan kalimat galat ganti sandi yang dikenal), bukan teks backend.
 */
export const LOGIN = {
  judul: 'Hello Again!',
  sub: "Welcome back, it's great to see you again!",
  logo: 'Nusantara Re - Application login',
  akun: 'Username',
  isianAkun: 'Enter your username',
  sandi: 'Password',
  isianSandi: 'Password',
  tampilkanSandi: 'Show password',
  sembunyikanSandi: 'Hide password',
  masuk: 'Login',
  memproses: 'Checking…',
  memuatSesi: 'Checking your session…',
  bantuan: 'Forgot your password or need access? Contact IT.',
  salah: 'Incorrect username or password.',
  terkunci: 'Account temporarily locked after 5 incorrect passwords. Try again in 15 minutes.',
  tidakTersedia: 'Login is not available on this server yet. Contact IT.',
  gagal: 'Login failed. Try again, or contact IT if it keeps happening.',
  cobaLagi: 'Try again',
  judulGanti: 'Change password',
  wajibGanti: 'Your temporary password must be changed before you continue.',
  sandiLama: 'Current password',
  sandiBaru: 'New password',
  ulangiSandi: 'Repeat new password',
  simpanSandi: 'Save password',
  minimal: 'At least 10 characters.',
  terlaluPanjang: 'Password can be at most 72 bytes.',
  samaDenganLama: 'The new password must be different from the current one.',
  tidakSama: 'The two new passwords do not match.',
  lamaSalah: 'Current password is incorrect.',
  batal: 'Cancel',
} as const

/** Panjang sandi minimal — sama dengan `login.PanjangMinSandi` di backend. */
export const PANJANG_MIN_SANDI = 10

/**
 * Label Kelola User — CRUD akun `M_LOGIN_GO` beserta workbasket dan menunya
 * (keputusan work owner 01-10-2026). Tab Security dan "Change Password Next
 * Login": permintaan work owner 01-10-2026 (label tab dan centang VERBATIM).
 *
 * ⚠️ `[tidak ada di korpus]` `[kerangka aplikasi]`.
 */
export const KELOLA_USER = {
  judul: 'Kelola User',
  tambah: 'Tambah user',
  cari: 'Cari username, Contact ID, nama, atau email',
  kosong: 'Belum ada user.',
  tidakCocok: 'Tidak ada user yang cocok dengan pencarian.',
  kolomAkun: 'Username',
  kolomNama: 'Nama',
  kolomJenjang: 'Organisasi / Divisi / Unit',
  kolomStatus: 'Status',
  kolomLogin: 'Login terakhir',
  kolomAksi: 'Aksi',
  belumLogin: 'belum pernah',
  aktif: 'Aktif',
  nonaktif: 'Nonaktif',
  terkunci: 'Terkunci',
  wajibGanti: 'Wajib ganti sandi',
  anda: 'Anda',
  ubah: 'Ubah',
  bukaKunci: 'Buka kunci',
  nonaktifkan: 'Nonaktifkan',
  aktifkan: 'Aktifkan',
  hapus: 'Hapus',
  judulBaru: 'Tambah user',
  judulUbah: 'Ubah user',
  akun: 'Username',
  akunTetap: 'Username tidak dapat diubah sesudah dibuat.',
  nama: 'Nama',
  tabProfil: 'Profil',
  tabSecurity: 'Security',
  sandi: 'Password awal',
  ulangiSandi: 'Ulangi password awal',
  catatanSandi: 'Minimal 10 karakter.',
  sandiBaru: 'Password baru',
  ulangiSandiBaru: 'Ulangi password baru',
  catatanSandiUbah:
    'Kosongkan bila password tidak diganti. Password baru (minimal 10 karakter) mengakhiri semua sesi user itu dan membuka kuncinya.',
  catatanSandiSendiri: 'Ini akun Anda: sesi yang sedang Anda pakai tetap berjalan sesudah password diganti.',
  wajibGantiCentang: 'Change Password Next Login',
  catatanWajibGanti: 'Dicentang: user wajib mengganti password saat login berikutnya. Tidak dicentang: tidak perlu.',
  organisasi: 'Organisasi',
  divisi: 'Divisi',
  unit: 'Unit',
  tidakDiisi: '— tidak diisi —',
  workbasket: 'Workbasket',
  menu: 'Menu yang boleh dibuka',
  dipilih: (n: number, dari: number) => `${n} dari ${dari} dipilih`,
  pilihSemua: 'Pilih semua',
  kosongkan: 'Kosongkan',
  menuDiriSendiri: 'Kelola User tidak dapat dicabut dari akun Anda sendiri.',
  // Hak menu (migrasi 914, keputusan work owner 04-10-2026).
  hakMenu: (label: string) => `Akses ${label}`,
  hakPenuh: 'Full',
  hakLihat: 'View only',
  catatanHak: 'View only: menu tampil dan dapat dibaca, tetapi tombol tambah, ubah, dan hapus tidak ada.',
  hakPilihDulu: 'Centang menunya dulu untuk memilih Full atau View only.',
  simpan: 'Simpan',
  menyimpan: 'Menyimpan…',
  batal: 'Batal',
  judulHapus: 'Hapus user permanen',
  tanyaHapus: (akun: string) =>
    `Hapus user ${akun} beserta seluruh workbasket dan menunya? Tindakan ini tidak dapat dibatalkan.`,
  hapusPermanen: 'Hapus permanen',
  menghapus: 'Menghapus…',
  galatAkun: 'Username hanya huruf, angka, titik, garis bawah, @, atau tanda hubung (maks. 64 karakter).',
  galatNama: 'Nama wajib diisi (maks. 150 karakter).',
  galatSandi: 'Password minimal 10 karakter.',
  galatUlangi: 'Kedua password tidak sama.',
  // Kontak akun (permintaan work owner 03-10-2026: "tambahkan email, no hp, nik dan jabatan; buat dalam bahasa
  // inggris"). Pesan galat SAMA dengan `login.Err*TidakSah` backend.
  email: 'Email',
  telepon: 'Phone Number',
  nik: 'Employee ID (NIK)',
  jabatan: 'Position',
  contohEmail: 'name@company.com',
  contohTelepon: '+62 812-3456-7890',
  kolomJabatan: 'Position',
  galatEmail: 'Email is not valid; use a format like name@company.com (max. 254 characters)',
  galatTelepon: 'Phone Number must contain 8 to 15 digits, may start with +, and may use spaces or hyphens',
  galatNIK: 'Employee ID (NIK) may only contain letters, digits, dots, slashes, or hyphens (max. 30 characters)',
  galatJabatan: 'Position is too long (max. 150 characters)',
  // Identitas akun (migrasi 905, keputusan work owner 03-10-2026: "M_LOGIN_GO ID nya pake CON-xxx"; "tambahkan
  // proteksi ... email sudah terdaftar"; "proteksi username sudah ada"). Pesan ganda SAMA dengan jawaban 409
  // `kelola_rute.go`. Login lewat email dibatalkan work owner - layar login tetap username.
  contactId: 'Contact ID',
  catatanContactId: 'Generated automatically when the user is created; it never changes.',
  galatUsernameTerdaftar: 'Username is already registered',
  galatEmailTerdaftar: 'Email is already registered to another account',
  tersimpan: (akun: string) => `User ${akun} tersimpan.`,
  terhapus: (akun: string) => `User ${akun} dihapus permanen.`,
  dinonaktifkan: (akun: string) => `User ${akun} dinonaktifkan; sesinya berakhir pada permintaan berikutnya.`,
  diaktifkan: (akun: string) => `User ${akun} diaktifkan.`,
  dibukaKunci: (akun: string) => `Kunci user ${akun} dibuka.`,
  memuat: 'Memuat daftar user…',
  memuatPilihan: 'Memuat pilihan…',
} as const

/** Template Manager (keputusan work owner 04-10-2026) - pengelola berkas templat unduhan semua menu. */
export const TEMPLATE_MANAGER = {
  judul: 'Template Manager',
  cari: 'Search menu or template',
  semuaMenu: 'All menus',
  memuat: 'Loading templates…',
  kosong: 'No template is registered yet.',
  tidakCocok: 'No template matches the search.',
  jumlah: (n: number) => `${n} template${n === 1 ? '' : 's'}`,
  kolomNama: 'Template',
  kolomBerkas: 'Active file',
  kolomKolom: 'Columns',
  kolomVersi: 'Version',
  kolomUnggah: 'Uploaded',
  kolomAksi: 'Action',
  bawaan: 'built-in',
  unduh: 'Download',
  unggah: 'Upload',
  riwayat: 'History',
  judulUnggah: 'Upload new version',
  untuk: 'For',
  dipakaiDi: 'Used in',
  syarat: 'Requirement',
  syaratTeks: (ekstensi: string, pemisah: string, kolom: number) =>
    [ekstensi, pemisah === '' ? '' : `separator ${pemisah}`, kolom > 0 ? `${kolom} columns` : ''].filter((s) => s !== '').join(' · '),
  pilihBerkas: 'Choose file',
  memeriksa: 'Checking file…',
  hasilCek: 'Check result',
  lolos: (kolom: number) => (kolom > 0 ? `${kolom} columns, matches the upload of this menu` : 'File accepted'),
  perbedaan: (n: number) => `${n} column title${n === 1 ? '' : 's'} differ from the active version:`,
  kolomKe: (n: number) => `column ${n}`,
  catatan: 'Note',
  catatanPetunjuk: 'Why is this version uploaded?',
  simpan: (versi: number) => `Save as v${versi}`,
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (versi: number) => `Saved as version ${versi}. Users now download this file.`,
  judulRiwayat: 'Version history',
  memuatRiwayat: 'Loading history…',
  kolomCatatan: 'Note',
  aktif: 'ACTIVE',
  aktifkan: 'Activate',
  diaktifkan: (versi: number) => (versi === 0 ? 'The built-in file is active again.' : `Version ${versi} is active again.`),
  aplikasi: '(application)',
} as const

// ---------------------------------------------------------------------------
// Peran dan nama produk - refactor bentuk B (30-09-2026): pindah apa adanya
// dari `labels.claimlife.ts`. Dipakai LEBIH DARI SATU modul (identitas sesi,
// Shell, Komite), jadi tempatnya berkas bersama.
// ---------------------------------------------------------------------------

/**
 * Peran — `[terverifikasi]` `Register_Flow.xml` `<pyPosition>` tiap
 * Assignment, dan ADR-U-0002.
 *
 * ⛔ Nilainya adalah PENGENAL yang dikirim ke backend sebagai `X-Peran`;
 * mengubahnya mengubah wewenang, bukan tampilan.
 */
export const PERAN = {
  admin: 'ReasLifeAdmin',
  medis: 'ReasLifeMedicalAdvisor',
  spv: 'ReasLifeSPV',
} as const

export type KodePeran = (typeof PERAN)[keyof typeof PERAN]

/** Sebutan peran di layar — pendamping, bukan pengganti. */
export const PERAN_ID: Record<KodePeran, string> = {
  [PERAN.admin]: 'Admin Klaim Jiwa',
  [PERAN.medis]: 'Penasihat Medis',
  [PERAN.spv]: 'Supervisor Klaim',
}

/**
 * Nama produk di topbar dan judul dokumen.
 *
 * ⛔ `[tidak ada di korpus]` — aset dan judul **e-Treaty** milik produk lain
 * dan TIDAK dipakai (brief §2 aturan 1). Teks netral dipakai sampai aset
 * resmi diberikan.
 */
export const PRODUK = {
  nama: 'Nusantara Re',
  sub: 'Reasuransi',
} as const
