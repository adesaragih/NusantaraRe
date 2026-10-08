// modul/treatyinadjustment/frontend/api.ts - panggilan backend modul Treaty In
// Adjustment (`backend/handlers/rute_treaty_in_adjustment.go`).
//
// ⛔ Seluruh rute modul ini BACA. Jalur tulis (Save, Submit, Actions, Decline
// offer) memakai rute modul Treaty In — lihat bagian TOMBOL TULIS di bawah.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_TREATYINADJUSTMENT = '/api/treaty-in-adjustment'

/** Kepala kontrak - lapisan BEKU yang seluruh versinya bagi (ADR-0040). */
export interface Kontrak {
  id: number
  nomorKontrakWarisan: string
  sifatProporsi: string
  tanggalMulai: string
  tanggalBerakhir: string
}

/**
 * Satu baris rantai versi - SAMA dengan `models.Versi`.
 *
 * `nomorUrutVersi` null = baris warisan yang belum dinomori ulang (tiket 10),
 * BUKAN nol. `idVersiDasar` null = versi PERTAMA kontrak itu, dan itu keadaan
 * yang benar, bukan data yang hilang.
 */
export interface Versi {
  id: number
  idKontrak: number
  nomorUrutVersi: number | null
  keadaanSiklusHidup: string
  jenisAddendum: string
  sifatMaterialAddendum: string
  tanggalBerlakuAddendum: string
  idVersiDasar: number | null
  namaKontrak: string
}

/** `GET /api/treaty-in-adjustment/kontrak` - kepala seluruh kontrak. */
export async function ambilKontrak(): Promise<Kontrak[]> {
  return minta<Kontrak[]>(`${PREFIX_TREATYINADJUSTMENT}/kontrak`)
}

/** `GET /api/treaty-in-adjustment/kontrak/{id}/versi` - rantai versi satu kontrak. */
export async function ambilRantaiVersi(idKontrak: number): Promise<Versi[]> {
  return minta<Versi[]>(`${PREFIX_TREATYINADJUSTMENT}/kontrak/${idKontrak}/versi`)
}

/**
 * Panel Attachment — dari `POOLDATA.M_ATTACHMENTTREATY_2`, tabel WARISAN
 * 43 baris. Nol tabel baru, nol migrasi.
 */
export interface BarisLampiranWarisan {
  id: string
  kodeKategori: string
  namaKategori: string
  namaBerkas: string
  jenisMime: string
  idSimpanan: string
  diunggah: string
  pengunggah: string
}

/**
 * Satu baris panel Attachment: Category + Count.
 *
 * ⭐ 8 Oktober 2026: `nama` dari `M_KATEGORIMASTERTREATY` (katalog RD Pega)
 * — kesebelas kode bernama, urut nama. `dipastikan` false hanya untuk kode
 * di data yang tidak ada di katalog; layar menampilkan kodenya.
 */
export interface BarisKategoriLampiran {
  kode: string
  nama: string
  cacah: number
  dipastikan: boolean
}

export interface LampiranKontrak {
  kategori: BarisKategoriLampiran[]
  berkas: BarisLampiranWarisan[]
}

/**
 * `GET /api/treaty-in-adjustment/kontrak-warisan/{id}/lampiran`.
 *
 * ⛔ Pengenalnya TEKS — `M_ATTACHMENTTREATY_2.TREATYID` adalah
 * `VARCHAR2(100)` milik sistem lama, bukan `idKontrak` model baru.
 *
 * `jenis` = `TreatyIn.ProportionType`: `NonProportional` menamai kode
 * `00007` dengan nama Non-Prop (`GetMasterTreatyCategory_Act` [2.2]).
 */
export async function ambilLampiran(masterID: string, jenis = ''): Promise<LampiranKontrak> {
  return minta<LampiranKontrak>(
    `${PREFIX_TREATYINADJUSTMENT}/kontrak-warisan/${encodeURIComponent(masterID)}/lampiran`,
    { kueri: { jenis: jenis !== '' ? jenis : undefined } },
  )
}

/** Satu baris panel History — dari `T_VIEW_COMMENT`, BACA SAJA. */
export interface BarisRiwayatWarisan {
  tanggal: string
  operator: string
  disetujui: string
  catatan: string
}

/** Satu baris panel `Existing Policy for Master ID` — `TREATYINPRODUCTION`. */
export interface BarisPolisMaster {
  nomorPolis: string
  pegaID: string
}

/** `GET /api/treaty-in-adjustment/kontrak-warisan/{id}/polis`. */
export async function ambilPolisMaster(idMaster: string): Promise<BarisPolisMaster[]> {
  return minta<BarisPolisMaster[]>(
    `${PREFIX_TREATYINADJUSTMENT}/kontrak-warisan/${encodeURIComponent(idMaster)}/polis`,
  )
}

/** `GET /api/treaty-in-adjustment/kontrak-warisan/{id}/riwayat`. */
export async function ambilRiwayat(masterID: string): Promise<BarisRiwayatWarisan[]> {
  return minta<BarisRiwayatWarisan[]>(
    `${PREFIX_TREATYINADJUSTMENT}/kontrak-warisan/${encodeURIComponent(masterID)}/riwayat`,
  )
}

// ===========================================================================
// LAYAR ADJUSTMENT — `TREATY_IN_EDM` + `M_TREATY_IN_EDM`, BACA SAJA.
// ===========================================================================

/**
 * Satu baris grid daftar — SAMA dengan `models.BarisPenyesuaian`.
 *
 * ⛔ Nilainya APA ADANYA dari kolom tabel warisan. `jenisPenyesuaian` dan
 * `jenisMaterial` adalah KODE; teks pilihannya tidak ada di ekspor.
 */
export interface BarisPenyesuaian {
  id: string
  idAsal: string
  jenisPenyesuaian: string
  jenisMaterial: string
  namaKontrak: string
  sifatProporsi: string
  asalBisnis: string
  cedant: string
  tanggalMulai: string
  tanggalBerakhir: string
  posisi: string
  statusAkseptasi: string
}

/**
 * Satu halaman di dalam dokumen — `TreatyIn` (New) atau `TreatyIn.OLDDATA`
 * (Old).
 *
 * ⛔ `medan` hanya memuat kunci yang ADA di dokumen. `kunci in medan` yang
 * salah berarti "tidak ada di sistem lama" — BUKAN "kosong".
 */
/** Nilai satu kunci baris: teks, atau larik anak (baris bersarang). */
export type NilaiBaris = string | BarisBersarang[]

/** Satu baris larik — kunci Pega; larik anak hanya pada larik yang punya anak. */
export interface BarisBersarang {
  [kunci: string]: NilaiBaris
}

export interface SisiPenyesuaian {
  medan: Record<string, string>
  larik: Record<string, BarisBersarang[]>
  /**
   * Larik akar yang punya larik anak (Limits, Share, FacultativeShareList,
   * Installment) sebagai simpul bersarang — dari tabel anak pendaratan.
   * Digabung ke `larik` saat panel lahir (`komponen/baris.ts`).
   */
  pohon?: Record<string, BarisBersarang[]>
}

export interface Penyesuaian {
  id: string
  idAsal: string
  baru: SisiPenyesuaian
  lama: SisiPenyesuaian
  /**
   * Isi sisi New ada di tabel pendaratan. `false` = kepalanya dilengkapi dari
   * `TREATY_IN_EDM`/master tetapi isi tabnya tidak ada — layar mengunci panel
   * New (mode lihat). Tidak ada = dianggap terdarat (bentuk lama).
   */
  terdarat?: boolean
}

/** `GET /api/treaty-in-adjustment/penyesuaian-warisan` — grid daftar. */
export async function ambilDaftarPenyesuaian(): Promise<BarisPenyesuaian[]> {
  return minta<BarisPenyesuaian[]>(`${PREFIX_TREATYINADJUSTMENT}/penyesuaian-warisan`)
}

/**
 * `GET /api/treaty-in-adjustment/penyesuaian-warisan/satu?id=…`.
 *
 * ⛔ Lewat parameter kueri: pengenalnya BERGARIS MIRING (`1000080/R02`).
 */
export async function ambilPenyesuaian(id: string): Promise<Penyesuaian> {
  return minta<Penyesuaian>(
    `${PREFIX_TREATYINADJUSTMENT}/penyesuaian-warisan/satu?id=${encodeURIComponent(id)}`,
  )
}

/** Picker tombol Add — `revisi` (Add Revision) atau `premi` (Add Adjustment Premium). */
export type JenisPicker = 'revisi' | 'premi'

/** Satu baris grid picker — `TempMasterList.pxResults` (`CARI1`…`CARI7`). */
export interface BarisMasterPilihan {
  id: string
  namaKontrak: string
  sifatProporsi: string
  asalBisnis: string
  cedant: string
  tanggalMulai: string
  tanggalBerakhir: string
}

/** `GET /api/treaty-in-adjustment/penyesuaian-warisan/master?jenis=…` — grid picker. */
export async function ambilDaftarMaster(jenis: JenisPicker): Promise<BarisMasterPilihan[]> {
  return minta<BarisMasterPilihan[]>(`${PREFIX_TREATYINADJUSTMENT}/penyesuaian-warisan/master?jenis=${jenis}`)
}

/** Parameter `TreatyInEDMSetValue` yang tombol `Choose` kirim. */
export interface MasukanDraf {
  id: string
  internalType: string
  materialType: string
}

/**
 * `POST /api/treaty-in-adjustment/penyesuaian-warisan/draf` — tombol `Choose`.
 *
 * ⛔ Menyusun penyesuaian baru TANPA menyimpannya: di Pega `Choose` juga
 * menyimpan, di sini simpanannya menunggu jalur Save (pemilik proses,
 * 7 Oktober 2026).
 */
export async function buatDrafPenyesuaian(m: MasukanDraf): Promise<Penyesuaian> {
  return minta<Penyesuaian>(`${PREFIX_TREATYINADJUSTMENT}/penyesuaian-warisan/draf`, { metode: 'POST', badan: m })
}

// ---------------------------------------------------------------------------
// Daftar pilihan dropdown/autocomplete panel New — rute modul Treaty In.
//
// ⭐ Panel New Adjustment ADALAH form Treaty In di Pega, dan kedua modul memakai
// RD yang sama (`BrowseCurrency_RD`, `BrowseTreatyGroup_RD`,
// `BrowseAgentNusaRe_RD`, `BrowseReinsuranceType_RD`). Rute Treaty In membaca
// tabel MASTER (CURRENCY, TREATYGROUP, AGENT, REINSURANCETYPE) dan katalog —
// nol JSON. Dipanggil lewat HTTP, bukan impor (lihat `komponen/rumus.ts`).
// ---------------------------------------------------------------------------

/** Prefix rute modul Treaty In. */
export const PREFIX_TREATYIN = '/api/treaty-in'

/** Satu baris master — `PilihanWarisan` modul Treaty In. */
export interface PilihanTreatyIn {
  id: string
  nama: string
  kembar?: boolean
  namaSoa?: string
}

/** Satu pilihan berlabel — `OpsiPilihan` modul Treaty In. */
export interface OpsiTreatyIn {
  value: string
  label: string
}

/** `GET /api/treaty-in/warisan/opsi-limits` — Treaty Type, Treaty Group, Currency. */
export interface OpsiLimitsTreatyIn {
  jenisTreaty: PilihanTreatyIn[]
  kelompokTreaty: PilihanTreatyIn[]
  mataUang: PilihanTreatyIn[]
  /**
   * ⭐ Dropdown `associated` rincian Layers Non-Prop — daftar yang SAMA
   * dengan tab Limits Non-Prop Treaty In (`Layer`/`Part Of` · `Cover` ·
   * `Currency Relation` · `Note` Reinstatement).
   */
  jenisLayer?: OpsiTreatyIn[]
  cover?: OpsiTreatyIn[]
  relasiMataUang?: OpsiTreatyIn[]
  catatanReinstatement?: OpsiTreatyIn[]
}
export async function ambilOpsiLimitsTreatyIn(): Promise<OpsiLimitsTreatyIn> {
  return minta<OpsiLimitsTreatyIn>(`${PREFIX_TREATYIN}/warisan/opsi-limits`)
}

/**
 * `GET /api/treaty-in/warisan/kelas-bisnis` — autocomplete `Class of
 * Business` (`BrowseTreatyBusinessWOType_RD`, parameter `pTreatyGroupId`):
 * `id` = `.BizCode`, `nama` = `.BIZNAME`. Sama dengan Treaty In.
 */
export async function ambilKelasBisnisTreatyIn(treatyGroupId: string): Promise<PilihanTreatyIn[]> {
  return minta<PilihanTreatyIn[]>(`${PREFIX_TREATYIN}/warisan/kelas-bisnis`, { kueri: { treatyGroupId } })
}

/** Satu susunan induk spreading — `SusunanSpreading` modul Treaty In. */
export interface SusunanSpreadingTreatyIn {
  reinsTypeId: string
  reinsTypeName: string
}

/**
 * `GET /api/treaty-in/warisan/spreading-induk` — dropdown `Spreading Type`
 * (`BrowseTreatyArrangement_ParentReinsMasterTrt`, parameter `TreatyGroupID`
 * dan `StartDate = TreatyIn.Commencement`). Sama dengan Treaty In.
 */
export async function ambilIndukSpreadingTreatyIn(treatyGroupId: string, mulai: string): Promise<SusunanSpreadingTreatyIn[]> {
  return minta<SusunanSpreadingTreatyIn[]>(`${PREFIX_TREATYIN}/warisan/spreading-induk`, { kueri: { treatyGroupId, mulai } })
}

/** `GET /api/treaty-in/warisan/opsi-kepala` — Bordereaux, Accounting Mode, Reporting Period. */
export interface OpsiKepalaTreatyIn {
  bordereaux: OpsiTreatyIn[]
  caraPembukuan: OpsiTreatyIn[]
  caraPembukuanNonProp: OpsiTreatyIn[]
  periodePelaporan: OpsiTreatyIn[]
}
export async function ambilOpsiKepalaTreatyIn(): Promise<OpsiKepalaTreatyIn> {
  return minta<OpsiKepalaTreatyIn>(`${PREFIX_TREATYIN}/warisan/opsi-kepala`)
}

/** `GET /api/treaty-in/warisan/cedant` — agen aktif, saringan `BrowseAgentNusaRe_RD`. */
export async function ambilAgenTreatyIn(): Promise<PilihanTreatyIn[]> {
  return minta<PilihanTreatyIn[]>(`${PREFIX_TREATYIN}/warisan/cedant`)
}


// ---------------------------------------------------------------------------
// ⭐ TOMBOL TULIS — Save, Submit, Actions, Decline offer (7 Oktober 2026).
//
// Penulisnya di modul Treaty In: dokumen Adjustment mendarat di tabel
// `T_TREATY_*` yang SAMA (kedua sisi, `MASTERID` dan `#LAMA`), dan kepalanya
// di `TREATY_IN_EDM`. Dipanggil lewat HTTP, bukan impor — pola rute
// `/hitung/*`. Isian masuk basis data HANYA dari tombol-tombol ini.
// ---------------------------------------------------------------------------

/** Satu halaman layar apa adanya — medan skalar dan larik (pohon tergabung). */
export interface SisiKiriman {
  medan: Record<string, string>
  larik: Record<string, BarisBersarang[]>
}

/** Tombol Save — `SaveTreatyIn_EDM_Act` dengan pra-DT `TreatyInAddNew`. */
export interface MasukanSimpanPenyesuaian {
  id: string
  /** Penyesuaian dari `Choose` yang belum pernah disimpan. */
  draf: boolean
  baru: SisiKiriman
  /** Panel Old — hanya untuk draf (penyesuaian tersimpan: baca saja). */
  lama?: SisiKiriman
}

/** Submit (`TreatyInSubmitEDM`) atau Actions (`TreatyInAkseptasiEDM_Act`). */
export interface MasukanKirimPenyesuaian extends MasukanSimpanPenyesuaian {
  aksi: 'submit' | 'akseptasi'
  /** `ChooseStatusAkseptasi` — Accept/Reject/Decline. */
  pilihan?: string
}

/** Jawaban tombol tulis — `services.HasilSimpan` modul Treaty In. */
/** Nasib satu berkas salinan lampiran master. */
export interface BerkasSalinanLampiran {
  nama: string
  berhasil: boolean
  pesan: string
}

/**
 * Salinan lampiran master (`TreatyRevisionCopyAttachment`) — hanya pada
 * tulisan PERTAMA draf penyesuaian.
 */
export interface SalinanLampiran {
  sumber: string
  tersalin: number
  berkas: BerkasSalinanLampiran[]
  pesan?: string
}

export interface HasilSimpanPenyesuaian {
  id: string
  pesan: string
  posisi: string
  status: string
  pemegangPosisi: string
  /** Properti terkirim yang TIDAK tersimpan — dilaporkan, tidak ditelan. */
  kunciTakTersimpan: string[]
  /** Draf yang PERTAMA kali tersimpan: salinan lampiran master. */
  salinanLampiran?: SalinanLampiran
}

/** `POST /api/treaty-in/penyesuaian/simpan` — tombol Save. */
export async function simpanPenyesuaian(m: MasukanSimpanPenyesuaian): Promise<HasilSimpanPenyesuaian> {
  return minta<HasilSimpanPenyesuaian>(`${PREFIX_TREATYIN}/penyesuaian/simpan`, { metode: 'POST', badan: m })
}

/** `POST /api/treaty-in/penyesuaian/kirim` — Submit atau Actions. */
export async function kirimPenyesuaian(m: MasukanKirimPenyesuaian): Promise<HasilSimpanPenyesuaian> {
  return minta<HasilSimpanPenyesuaian>(`${PREFIX_TREATYIN}/penyesuaian/kirim`, { metode: 'POST', badan: m })
}

/**
 * `POST /api/treaty-in/penyesuaian/hapus` — Decline offer
 * (`TreatyInDeclineConfirmation_postactEDM`): baris EDM DIHAPUS fisik.
 */
export async function hapusPenyesuaian(id: string): Promise<HasilSimpanPenyesuaian> {
  return minta<HasilSimpanPenyesuaian>(`${PREFIX_TREATYIN}/penyesuaian/hapus`, { metode: 'POST', badan: { id } })
}

/** Satu baris `AchievementLists` untuk log — ejaan properti Pega. */
export type BarisLogAchievement = Record<
  'Quarter' | 'QUARTERYEAR' | 'CurrencyID' | 'Currency' | 'PREMIUM' | 'RICOMM' | 'BROKERAGE' | 'NETPREMIUM' |
    'PaidClaim' | 'CASHCALL' | 'OutstandingClaim' | 'IncuredClaim' | 'Total' | 'LossRatio',
  string
>

/**
 * `POST /api/treaty-in/achievement/log` — tombol `Submit` sub-tab Achievement
 * (`InsertToLogAchievement` → `LOG_ACHIEVEMENT`, keputusan pemakai 8 Oktober
 * 2026). Baris ber-Quarter kosong dilewati server.
 */
export async function catatLogAchievement(m: { idKontrak: string; baris: BarisLogAchievement[] }): Promise<{ disisipkan: number; dilewati: number }> {
  return minta<{ disisipkan: number; dilewati: number }>(`${PREFIX_TREATYIN}/achievement/log`, { metode: 'POST', badan: m })
}
