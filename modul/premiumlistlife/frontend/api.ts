// modul/premiumlistlife/frontend/api.ts - panggilan backend modul PremiumList Life, satu
// fungsi per endpoint. Klien HTTP-nya `inti/klien.ts` (refactor bentuk B,
// 30-09-2026: dipecah dari `services/api.ts` tanpa mengubah satu panggilan pun).

import { minta, mintaFormulir } from '../../../inti/frontend/klien'

// ---------------------------------------------------------------------------
// MODUL PREMIUMLIST LIFE — kotak masuk dan keputusan penawaran (tiket 01).
// ---------------------------------------------------------------------------

/** Satu baris kotak masuk PremiumList — `InboxPremiumList.xml`. */
export interface BarisInboxPolis {
  caseId: string
  /** `YYYY-MM-DD`; kosong berarti belum ada. */
  tglCreate: string
  createOpName: string
  statusWork: string
  /** `Offer` atau `Premium` — posisi LAYAR, bukan tahap. */
  position: string
  cedingCoName: string
  policyHolderName: string
  plNumber: string
  riSlipRnm: string
  type: string
  marketingName: string
  sobName: string
  dateReceived: string
}

/** Satu halaman kotak masuk PremiumList. */
export interface HalamanInboxPolis {
  baris: BarisInboxPolis[]
  /** Cacah SELURUH baris yang cocok, bukan yang di halaman ini. */
  total: number
  halaman: number
  ukuran: number
}

/**
 * Membaca kotak masuk PremiumList — `GET /api/polis-life`.
 *
 * `posisi` kosong berarti seluruh posisi.
 */
export async function ambilKotakMasukPolis(
  posisi = '',
  halaman = 1,
  ukuran = 20,
): Promise<HalamanInboxPolis> {
  return minta<HalamanInboxPolis>('/api/polis-life', {
    kueri: { posisi, halaman, ukuran },
  })
}

/**
 * `FlagPolicy` tombol portal — VERBATIM `Section/PremiumList.xml`:
 * `Input Offer` b3310/b3597, `Input Premium` b3958/b4233 (butir bn).
 */
export const FLAG_POLIS = {
  inputOffer: '0',
  inputPremium: '1',
} as const

export type FlagPolis = (typeof FLAG_POLIS)[keyof typeof FLAG_POLIS]

/** Kasus polis yang baru lahir — cukup untuk membukanya. */
export interface KasusPolisBaru {
  caseId: string
  /** Tahap pertamanya: `Input Offer Life` untuk KEDUA bendera. */
  statusWork: string
  position: string
  flag: FlagPolis
}

/**
 * Tombol portal `Input Offer` / `Input Premium` — `CreateInputLife`.
 *
 * ⛔ Kedua bendera mulai di tahap yang SAMA (`Input Offer Life`); benderanya
 * baru bekerja di `Decision3` sesudah `Confirm` (decision table
 * `IsFlagOnGoingPolicy`: "0" → Offer, "1" → Premium).
 */
export async function buatKasusPolis(flag: FlagPolis): Promise<KasusPolisBaru> {
  return minta<KasusPolisBaru>('/api/polis-life', { metode: 'POST', badan: { flag } })
}

/** Apa yang terjadi sesudah sebuah keputusan penawaran. */
export interface AkibatKeputusanPolis {
  /** Terisi bila kasus BERPINDAH tahap. */
  tahapTujuan: string
  /** Terisi bila kasus DITUTUP. */
  statusWork: string
  /**
   * PL Number yang terbit (Confirm di Input Premium Detail); kosong selain itu.
   * Layar menampilkannya sebelum kembali ke kotak masuk (03-10-2026).
   */
  plNumber?: string
  /** WPC `YYYY-MM-DD` yang dihitung bersama PL Number (03-10-2026). */
  wpc?: string
}

/**
 * Menerapkan satu keputusan penawaran — `POST …/keputusan`.
 *
 * ⛔ `Reject` hanya sah dari tahap Input Premium Detail; di tahap penawaran
 * ia TIDAK punya konektor, dan backend menjawab **409**. Layar karena itu
 * tidak menawarkannya di sana — lihat `bolehRejectDiTahap`.
 *
 * ⭐ GILIRAN-14 butir bq: `Confirm` di tahap penawaran langsung menutup
 * (bendera `"0"`, Offer) atau memindahkan ke Input Premium Detail (`"1"`,
 * Premium) — `Decision3` tidak ditanyakan. Bendera di luar decision table
 * dijawab **409**.
 */
export async function putuskanPenawaran(
  polisID: string,
  keputusan: string,
): Promise<AkibatKeputusanPolis> {
  return minta<AkibatKeputusanPolis>(
    `/api/polis-life/${encodeURIComponent(polisID)}/keputusan`,
    { metode: 'POST', badan: { keputusan } },
  )
}

/** Nama tahap polis — VERBATIM `pyWorkStatus` tiap assignment. */
export const TAHAP_POLIS = {
  penawaran: 'Input Offer Life',
  detail: 'Input Premium Detail',
  summary: 'Input Premium Summary',
} as const

/**
 * Apakah kasus ber-`statusWork` ini dapat DIBUKA dari kotak masuk: hanya bila
 * statusnya salah satu tahap aktif (`TAHAP_POLIS`).
 *
 * ⛔ [keputusan work owner 02-10-2026] Kasus tertutup (`Resolved-Completed`,
 * `Resolved-Rejected`, dan status lain di luar tahap aktif) TIDAK dibuka:
 * layar keputusannya hanya menawarkan tombol yang server tolak ("kasus polis
 * sudah ditutup"). Daftar putih, bukan daftar hitam — status penutup baru
 * pun otomatis tidak dapat dibuka.
 */
export function kasusBisaDibuka(statusWork: string): boolean {
  return (Object.values(TAHAP_POLIS) as string[]).includes(statusWork.trim())
}

/**
 * Apakah `Reject` punya jalur dari tahap ini.
 *
 * ⛔ `Reject` muncul TEPAT SEKALI di seluruh flow — `Transition9` b2306 pada
 * `Decision2`, yaitu sesudah Input Premium Detail. Menawarkannya di tahap
 * penawaran berarti menjanjikan jalur yang tidak pernah ada.
 */
export function bolehRejectDiTahap(tahap: string): boolean {
  return tahap === TAHAP_POLIS.detail
}

/**
 * Periode produksi yang berlaku — `GET /api/polis-life/periode`.
 *
 * ⛔ Ditampilkan SEBELUM pemakai menyimpan (AC tiket 02), bukan tersimpan
 * diam-diam. Transaksi yang mendarat di bulan yang salah karena seseorang
 * menyimpannya lewat tengah malam adalah kekeliruan yang hanya dapat
 * dicegah dengan menunjukkannya lebih dulu.
 *
 * ⚠️ Menjawab **503** bila `POOLDATA.TANGGAL_CLOSING` kosong, dan pesannya
 * menyebut tabel itu. Pega diam-diam memakai `25`; kami menolak — fallback
 * diam membukukan ke periode yang salah tanpa meninggalkan jejak.
 */
export async function ambilPeriodeProduksi(): Promise<string> {
  const hasil = await minta<{ periode: string }>('/api/polis-life/periode')
  return hasil.periode
}

// ——— Tiket 03: Premium List Detail dan penomoran PL ———

/** Kepala polis di atas grid — `GET /api/polis-life/{id}`. */
export interface KepalaPolis {
  polisId: string
  /** `T_PREMIUM_LIST.TYPE` — QR / QP / TP / TR. */
  type: string
  businessCode: string
  /** `PL_NUMBER`; kosong berarti polis belum bernomor. */
  plNumber: string
  /**
   * Medan layar lama yang TIDAK kami punya kolomnya, beserta alasannya.
   *
   * ⛔ Datang dari server, bukan dikarang layar. Orang yang membandingkan
   * layar baru dengan layar lama akan menghitung kolomnya; yang menemukan
   * selisih tanpa penjelasan akan menyimpulkan datanya hilang.
   */
  medanTanpaKolom: { medan: string; alasan: string }[]
}

/** Satu baris grid peserta. */
export interface BarisPesertaPolis {
  id: string
  /**
   * Nilai berkunci NAMA KOLOM, seluruhnya teks.
   *
   * ⛔ Uang tiba sebagai TEKS dan tetap teks (ADR-U-0003). `Number(...)`
   * atas premi delapan angka desimal membulatkannya diam-diam, dan
   * pembulatan di jalan pulang tidak kalah salah dari pembulatan saat
   * menyimpan.
   */
  nilai: Record<string, string>
}

/** Satu halaman grid peserta — `GET /api/polis-life/{id}/peserta`. */
export interface HalamanPesertaPolis {
  /**
   * Nama kolom, URUT seperti `PL_Detail_Sec` menampilkannya.
   *
   * ⛔ DATANG DARI SERVER, tidak diketik ulang di sini. Dua daftar kolom —
   * satu di Go, satu di TypeScript — akan berselisih, dan selisihnya muncul
   * sebagai angka di bawah judul kolom yang salah.
   */
  kolom: string[]
  /**
   * Kolom ber-jenis ANGKA (dari server) — diberi pemisah ribuan (03-10-2026).
   * Opsional: server lama tidak mengirimnya, dan tanpa itu tidak ada pemisah.
   */
  kolomAngka?: string[]
  baris: BarisPesertaPolis[]
  total: number
  halaman: number
  ukuran: number
}

/** Jawaban penerbitan atau pembacaan `PL_NUMBER`. */
export interface HasilNomorPL {
  nomor: string
  /** Membedakan nomor yang baru lahir dari nomor yang sudah ada. */
  baruTerbit: boolean
  periode: string
  barisPeserta: number
}

/** Membaca kepala polis — `GET /api/polis-life/{id}`. */
export async function ambilKepalaPolis(polisID: string): Promise<KepalaPolis> {
  return minta<KepalaPolis>(`/api/polis-life/${encodeURIComponent(polisID)}`)
}

/** Membaca satu halaman peserta — `GET /api/polis-life/{id}/peserta`. */
export async function ambilPesertaPolis(
  polisID: string,
  halaman = 1,
  ukuran = 50,
): Promise<HalamanPesertaPolis> {
  const q = `?halaman=${String(halaman)}&ukuran=${String(ukuran)}`
  return minta<HalamanPesertaPolis>(
    `/api/polis-life/${encodeURIComponent(polisID)}/peserta${q}`,
  )
}

/**
 * Menerbitkan `PL_NUMBER` — `POST /api/polis-life/{id}/nomor`.
 *
 * ⛔ TANPA BADAN PERMINTAAN, dan itu disengaja. Tidak satu pun bahan nomor
 * boleh datang dari klien: awalan, tipe, kode bisnis, periode, dan urut
 * seluruhnya dibaca server dari sumbernya. Nomor yang bahannya dapat disebut
 * pemanggil adalah nomor yang dapat dipilih pemanggil.
 *
 * ⛔ NOMOR LAHIR SEKALI. Permintaan kedua atas polis yang sudah bernomor
 * menjawab nomor yang SAMA dengan `baruTerbit: false`, dan penghitungnya
 * tidak bergerak.
 *
 * ⚠️ Menjawab **409** bila polis belum punya satu pun baris peserta —
 * nomornya belum punya tempat tersimpan. Unggah rincian peserta lebih dahulu.
 */
export async function terbitkanNomorPL(polisID: string): Promise<HasilNomorPL> {
  return minta<HasilNomorPL>(
    `/api/polis-life/${encodeURIComponent(polisID)}/nomor`,
    { metode: 'POST' },
  )
}

// ——— Tiket 04: unggahan CSV peserta ———

/** Satu penolakan baris CSV. */
export interface PenolakanUnggah {
  /** Nomor BARIS DATA — 1 untuk baris pertama di bawah judul. */
  baris: number
  kolom: string
  /** Pesan VERBATIM dari sistem lama. */
  pesan: string
  /** Keterangan tepat: apa yang sebenarnya salah. */
  sebab: string
}

/** Jawaban tinjauan unggahan. */
export interface HasilTinjauUnggah {
  cacahBaris: number
  cacahDitolak: number
  ditolak: PenolakanUnggah[]
  lolos: boolean
}

/** Jawaban penyimpanan unggahan. */
export interface HasilSimpanUnggah {
  cacahBaris: number
  cacahDisimpan: number
  cacahDihapus: number
  ditolak: PenolakanUnggah[]
}

/**
 * Meninjau berkas CSV — `POST /api/polis-life/{id}/unggah/tinjau`.
 *
 * ⛔ TIDAK MENYIMPAN APA PUN. Ia mengurai, memvalidasi, dan mengembalikan
 * hasilnya — supaya pemakai dapat melihat apa yang lolos dan apa yang ditolak
 * SEBELUM apa pun tersimpan permanen (AC tiket 04).
 */
export async function tinjauUnggahPolis(
  polisID: string,
  berkas: File,
): Promise<HasilTinjauUnggah> {
  const isi = new FormData()
  isi.append('berkas', berkas)
  return mintaFormulir<HasilTinjauUnggah>(
    `/api/polis-life/${encodeURIComponent(polisID)}/unggah/tinjau`,
    isi,
  )
}

/**
 * Menyimpan berkas CSV — `POST /api/polis-life/{id}/unggah/simpan`.
 *
 * ⛔ BERKASNYA DIKIRIM ULANG, bukan sekadar "setujui yang tadi". Server
 * memvalidasi ulang: klien yang dapat melewatkan tinjauan adalah klien yang
 * dapat menyimpan apa saja.
 *
 * ⛔ MENGGANTI, bukan menumpuk — baris peserta lama polis ini dihapus di
 * dalam transaksi yang sama.
 *
 * ⚠️ Menjawab **409** bila masih ada penolakan, dan badan jawabannya MEMUAT
 * daftar penolakannya — jadi tidak perlu meninjau ulang untuk tahu apa yang
 * salah.
 */
export async function simpanUnggahPolis(
  polisID: string,
  berkas: File,
): Promise<HasilSimpanUnggah> {
  const isi = new FormData()
  isi.append('berkas', berkas)
  return mintaFormulir<HasilSimpanUnggah>(
    `/api/polis-life/${encodeURIComponent(polisID)}/unggah/simpan`,
    isi,
  )
}

// ——— Tiket 05a bagian 2: rekap premium list (ShowLifePremiumSummary) ———

/**
 * Satu baris rekap per mata uang.
 *
 * ⛔ UANG TETAP TEKS — `975.0000`, bukan `975`. Ekor nol itu bukti bahwa
 * pembulatan empat angka `@divide(…,1,4)` sudah terjadi di server.
 */
export interface RekapMataUangPolis {
  currency: string
  premium: string
  commission: string
  balance: string
  /** Ke-33 kolom yang dijumlah apa adanya, berkunci nama kolom. */
  jumlah: Record<string, string>
  cacahBaris: number
}

/** Jawaban `GET /api/polis-life/{id}/summary`. */
export interface HasilRekapPolis {
  tipe: string
  rekap: RekapMataUangPolis[]
}

/** Jawaban `POST /api/polis-life/{id}/summary`. */
export interface HasilSubmitRekapPolis {
  nomor: HasilNomorPL
  rekap: RekapMataUangPolis[]
  rekapDihapus: number
  pesertaWarisan: number
  /**
   * Tiket 06 — efek keluar SESUDAH commit. Ringkasan, bukan galat: premium
   * list sudah tersimpan walau Arasapas gagal.
   */
  efekKeluar: { dilewati: boolean; gagal: string[]; tidakTerantre: number }
}

/** Menghitung rekap tanpa menyimpan — `GET /api/polis-life/{id}/summary`. */
export async function ambilRekapPolis(polisID: string): Promise<HasilRekapPolis> {
  return minta<HasilRekapPolis>(`/api/polis-life/${encodeURIComponent(polisID)}/summary`)
}

/**
 * Submit rekap — `POST /api/polis-life/{id}/summary`.
 *
 * ⛔ TANPA badan: tidak satu pun angka rekap boleh datang dari layar. Server
 * menomori, merekap, dan menyalin peserta warisan dalam SATU transaksi.
 */
export async function submitRekapPolis(polisID: string): Promise<HasilSubmitRekapPolis> {
  return minta<HasilSubmitRekapPolis>(
    `/api/polis-life/${encodeURIComponent(polisID)}/summary`,
    { metode: 'POST' },
  )
}

// ——— Tiket 01 bagian 3: form penawaran (layar Input Offer) ———

/** Satu pilihan tertutup — kode yang disimpan, nama yang tampil. */
export interface PilihanKode {
  kode: string
  nama: string
}

/** Satu baris riwayat penawaran — `T_VIEW_SUGGEST` (`AddHistorySuggest`). */
export interface BarisRiwayatPenawaran {
  no: number
  /** ISO dari server. */
  dateSuggest: string
  picSuggest: string
  isCedingConfirm: string
  commentSuggest: string
  /** `Offer` atau `Bind`. */
  initialSuggest: string
}

/** Isi layar Input Offer — `GET /api/polis-life/{id}/penawaran`. */
export interface PenawaranPolis {
  caseId: string
  /** `FlagOnGoingPolicy` — menentukan radio Status mana yang tampil. */
  flag: string
  tahap: string
  /** False di luar tahap `Input Offer Life`: layar hanya menampilkan. */
  bolehDisimpan: boolean
  noOffer: string
  cedingCo: string
  cedingCoName: string
  policyHolder: string
  policyHolderName: string
  typeCeding: string
  typeCedingName: string
  /** "Reinsurance Type" — turunan System Reinsurance (`SetReinsuranceType`), dihitung server. */
  jenisAsuransi: string
  businessCode: string
  businessName: string
  /** ISO dari server; `null` bila belum diisi. */
  dateReceived: string | null
  description: string
  /** Sel penawaran migrasi 059 (+ SUM_INSURED, STATUS_UPDATE). */
  batasUsiaPeserta: number | null
  periodePertanggungan: string
  /** Uang sebagai TEKS — jangan `Number(...)` (ADR-U-0003). */
  sumInsured: string
  tanggalPenawaran: string | null
  tanggalRespon: string | null
  tanggalKonfirmasi: string | null
  tbc: number | null
  /** "Max TBC" — dihitung server (`SetMaxTBCLife_Act`). */
  tanggalTbc: string | null
  statusUpdate: string
  keteranganMarketing: string
  /** Sel migrasi 060. */
  qqName: string
  jenisUsaha: string
  ketentuanUnderwriting: string
  tanggalKonfirmasiBalik: string | null
  tanggalRealisasi: string | null
  tanggalBind: string | null
  statusFinal: string
  /** Status TERAKHIR disimpan (`STATUS_PENAWARAN` baris utama, migrasi 062). */
  status: string
  riwayat: BarisRiwayatPenawaran[]
  /**
   * ⛔ Pilihan tertutup DATANG DARI SERVER — satu daftar di Go, nol salinan
   * di sini. Dua daftar akan berselisih, dan selisihnya menyimpan kode yang
   * server tolak.
   */
  pilihan: {
    typeCeding: PilihanKode[]
    classOfBusiness: PilihanKode[]
    status: PilihanKode[]
  }
}

/** Badan `PUT /api/polis-life/{id}/penawaran` — kode, bukan nama. */
export interface IsiPenawaranPolis {
  cedingCo: string
  cedingCoName: string
  policyHolder: string
  policyHolderName: string
  typeCeding: string
  businessCode: string
  /** `YYYY-MM-DD` atau kosong. */
  dateReceived: string
  description: string
  status: string
  /** Angka dikirim sebagai TEKS; kosong berarti tidak diisi. */
  batasUsiaPeserta: string
  periodePertanggungan: string
  sumInsured: string
  tanggalPenawaran: string
  tanggalRespon: string
  tanggalKonfirmasi: string
  tbc: string
  statusUpdate: string
  keteranganMarketing: string
  qqName: string
  jenisUsaha: string
  ketentuanUnderwriting: string
  tanggalKonfirmasiBalik: string
  tanggalRealisasi: string
  tanggalBind: string
  statusFinal: string
}

/** Satu baris popup pilihan master. */
export interface BarisRujukanPolis {
  id: string
  nama: string
  keterangan?: string
}

/** Membaca isi layar Input Offer. */
export async function ambilPenawaranPolis(polisID: string): Promise<PenawaranPolis> {
  return minta<PenawaranPolis>(`/api/polis-life/${encodeURIComponent(polisID)}/penawaran`)
}

/**
 * `Save Offer` — menyimpan isian dan menambah satu baris riwayat.
 *
 * Menjawab isi layar TERBARU (nama turunan dan riwayat dihitung server).
 * ⚠️ **400** bila isian wajib kosong — pesannya kalimat Pega, satu per baris.
 */
export async function simpanPenawaranPolis(
  polisID: string,
  isi: IsiPenawaranPolis,
): Promise<PenawaranPolis> {
  return minta<PenawaranPolis>(`/api/polis-life/${encodeURIComponent(polisID)}/penawaran`, {
    metode: 'PUT',
    badan: isi,
  })
}

/** Popup `Choose Ceding Name` — `BrowseCedingCoLife_RD`. */
export async function cariCedingPolis(cari: string): Promise<BarisRujukanPolis[]> {
  return minta<BarisRujukanPolis[]>('/api/polis-life/cari-ceding', { kueri: { cari } })
}

/** Popup Policy Holder — `BrowseClientNusaRe_RD`. */
export async function cariPemegangPolis(cari: string): Promise<BarisRujukanPolis[]> {
  return minta<BarisRujukanPolis[]>('/api/polis-life/cari-pemegang-polis', { kueri: { cari } })
}

/** ISO dari server → `YYYY-MM-DD` untuk `<input type="date">`. */
export function tanggalMasukan(iso: string | null): string {
  if (iso === null || iso === '') return ''
  return iso.slice(0, 10)
}

/** Isi layar → badan simpan, dengan isian yang sedang diketik. */
export function isiDariPenawaran(p: PenawaranPolis): IsiPenawaranPolis {
  return {
    cedingCo: p.cedingCo,
    cedingCoName: p.cedingCoName,
    policyHolder: p.policyHolder,
    policyHolderName: p.policyHolderName,
    typeCeding: p.typeCeding,
    businessCode: p.businessCode,
    dateReceived: tanggalMasukan(p.dateReceived),
    description: p.description,
    status: p.status,
    batasUsiaPeserta: p.batasUsiaPeserta === null ? '' : String(p.batasUsiaPeserta),
    periodePertanggungan: p.periodePertanggungan,
    sumInsured: p.sumInsured,
    tanggalPenawaran: tanggalMasukan(p.tanggalPenawaran),
    tanggalRespon: tanggalMasukan(p.tanggalRespon),
    tanggalKonfirmasi: tanggalMasukan(p.tanggalKonfirmasi),
    tbc: p.tbc === null ? '' : String(p.tbc),
    statusUpdate: p.statusUpdate,
    keteranganMarketing: p.keteranganMarketing,
    qqName: p.qqName,
    jenisUsaha: p.jenisUsaha,
    ketentuanUnderwriting: p.ketentuanUnderwriting,
    tanggalKonfirmasiBalik: tanggalMasukan(p.tanggalKonfirmasiBalik),
    tanggalRealisasi: tanggalMasukan(p.tanggalRealisasi),
    tanggalBind: tanggalMasukan(p.tanggalBind),
    statusFinal: p.statusFinal,
  }
}

// ——— Tiket 03 bagian 2: data polis layar Input Premium Detail ———

/** Data polis — `GET /api/polis-life/{id}/data-polis`. */
export interface DataPolis {
  caseId: string
  tahap: string
  /** False di luar tahap Input Premium Detail: layar hanya menampilkan. */
  bolehDisimpan: boolean
  type: string
  productNameId: string
  productName: string
  sourceOfBusiness: string
  sobName: string
  cedingCo: string
  cedingCoName: string
  policyHolder: string
  policyHolderName: string
  riSlipRnm: string
  proRateType: string
  moId: string
  marketingCode: string
  marketingName: string
  /** Desimal sebagai TEKS (ADR-U-0003); kosong bila belum diisi. */
  annuityInterest: string
  premiumRefundFactor: string
  retroId: string
  retroName: string
  securityReinsurerId: string
  securityReinsurer: string
  /** Email Received Date — ISO; DAPAT DIISI dan WAJIB di layar ini (02-10-2026). */
  dateReceived: string | null
  /** Dibaca saja (WPCLife_Act belum dibawa). */
  wpc: string | null
  /** Pesan SavePremiumList_Act langkah 9 sesudah Save Data (data tetap tersimpan). */
  peringatan: string[] | null
  pilihan: { type: PilihanKode[]; proRateType: PilihanKode[] }
}

/** Badan `PUT .../data-polis`. */
export type IsiDataPolis = Omit<
  DataPolis,
  'caseId' | 'tahap' | 'bolehDisimpan' | 'wpc' | 'peringatan' | 'pilihan' | 'dateReceived'
> & {
  /** `YYYY-MM-DD` (masukan `type="date"`); kosong = belum diisi. */
  dateReceived: string
}

/** Satu baris popup Choose Product Name. */
export interface BarisProdukPolis {
  id: string
  inwardName: string
  cedingId: string
  ceding: string
  sobId: string
  sobName: string
  policyHolder: string
  policyHolderName: string
}

/** Satu pilihan Marketing Officer. */
export interface BarisMarketingPolis {
  id: string
  kode: string
  nama: string
}

/** Membaca data polis. */
export async function ambilDataPolis(polisID: string): Promise<DataPolis> {
  return minta<DataPolis>(`/api/polis-life/${encodeURIComponent(polisID)}/data-polis`)
}

/** Menyimpan data polis (tanpa perhitungan premi). Menjawab isi terbaru. */
export async function simpanDataPolis(polisID: string, isi: IsiDataPolis): Promise<DataPolis> {
  return minta<DataPolis>(`/api/polis-life/${encodeURIComponent(polisID)}/data-polis`, {
    metode: 'PUT',
    badan: isi,
  })
}

/** Popup Choose Product Name — produk milik Ceding kasus. */
export async function cariProdukPolis(polisID: string, cari: string): Promise<BarisProdukPolis[]> {
  return minta<BarisProdukPolis[]>(`/api/polis-life/${encodeURIComponent(polisID)}/cari-produk`, {
    kueri: { cari },
  })
}

/** Autocomplete Marketing Officer. */
export async function cariMarketingPolis(cari: string): Promise<BarisMarketingPolis[]> {
  return minta<BarisMarketingPolis[]>('/api/polis-life/cari-marketing', { kueri: { cari } })
}

/** Autocomplete R/I SLIP RNM No. — nomor polis RNML-Q / RNML-F. */
export async function cariRISlipPolis(cari: string): Promise<BarisRujukanPolis[]> {
  return minta<BarisRujukanPolis[]>('/api/polis-life/cari-rislip', { kueri: { cari } })
}

/** Data polis → badan simpan. */
export function isiDariDataPolis(d: DataPolis): IsiDataPolis {
  const { caseId: _c, tahap: _t, bolehDisimpan: _b, wpc: _w, peringatan: _r, pilihan: _p, ...isi } = d
  return { ...isi, dateReceived: tanggalMasukan(d.dateReceived) }
}
