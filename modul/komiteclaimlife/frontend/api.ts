// modul/komiteclaimlife/frontend/api.ts - panggilan backend modul Komite Claim Life, satu fungsi
// per endpoint. Klien HTTP-nya `inti/klien.ts` (refactor bentuk B, 30-09-2026:
// dipecah dari `services/api.ts` tanpa mengubah satu panggilan pun).

import { BATAS_WAKTU_MS, kegagalanDari, minta, rakitURL } from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'

// ——— Komite Claim Life tiket 01: Inbox Komite ———

/** Satu baris Inbox Komite. Uang TEKS; status baris KATA. */
export interface BarisInboxKomite {
  kasusId: string
  tglUpdate: string
  statusWork: string
  klaimId: string
  nomorKlaim: string
  tingkatBerjalan: number
  komiteLoop: number
  nilaiKlaim: string
  mataUang: string
  statusBaris: string
}

export interface AnggotaKasusKomite {
  urut: number
  jabatan: string
  approval: string
  kataApproval: string
  komentar: string
  tglApprove: string
  saya: boolean
}

/** Satu efek keluar kasus komite — KATA, bukan kode (tiket 08). */
export interface EfekKomite {
  jenis: string
  keadaan: string
  percobaan: number
  sejak: string
}

export interface KasusKomite {
  kasus: BarisInboxKomite
  adjustmentId: string
  tangga: AnggotaKasusKomite[]
  giliranSaya: boolean
  /** Tiket 08 — keadaan ringkas ('' bila belum ada keputusan) dan tiap efek. */
  efek: { keadaan: string; efek: EfekKomite[] }
}

/** Keputusan asli satu tingkat yang tertimpa langkah 5.1 (OQ-K-05). */
export interface KeputusanAsliKomite {
  status: string
  comment: string
}

/** Riwayat tangga — `GET /api/komite/{id}/riwayat` (tiket 09; siapa pun). */
export interface RiwayatKomite {
  kasusId: string
  adjustmentId: string
  tangga: {
    urut: number
    committee: string
    anggota: string
    /** KATA — Setuju / Tolak / Menunggu / Dilewati (eskalasi). */
    status: string
    dateApprove: string
    comment: string
    /**
     * Keputusan tingkat ini SEBELUM langkah 5.1 `KomitePostAdjustment`
     * menimpanya (tingkat akhir menolak → seluruh tingkat `2`, komentar
     * kosong). Dibaca dari jejak — OQ-K-05, GILIRAN-17. Tidak ada bila tidak
     * tertimpa.
     */
    asli?: KeputusanAsliKomite
  }[]
  eskalasi: { dariTingkat: number; keTingkat: number; oleh: string; waktu: string }[]
  /** Date / PIC / Remarks pop-up Send Claim to Committee (Claim Life, 12-10-2026) — jejak klaim induk; baca saja. */
  penyerahan?: { tanggal: string; pic: string; remarks: string }
}

export async function ambilRiwayatKomite(kasusID: string): Promise<RiwayatKomite> {
  return minta<RiwayatKomite>(`/api/komite/${encodeURIComponent(kasusID)}/riwayat`)
}

/** Jawaban `POST /api/komite/{id}/keputusan`. */
export interface HasilKeputusanKomite {
  tingkatDiputus: number
  keputusan: string
  kataKeputusan: string
  berlanjut: boolean
  tingkatBerikut: number
  akseptasiAkhir: boolean
  tolakAkhir: boolean
  /** Tiket 04a — terisi hanya pada Setuju di tingkat akhir. */
  nomorAkseptasi: string
  /** Tiket 06 — efek keluar yang DIANTRE: tersimpan, belum tuntas. */
  efekTertunda: string[]
}

/**
 * Keputusan satu tingkat — `POST /api/komite/{id}/keputusan`.
 *
 * ⚠️ Menjawab **501** di tingkat akhir sampai akseptasi (tiket 04a/04b) dan
 * penolakan ke baris klaim (tiket 05) dibangun — pesan server tampil apa adanya.
 */
export async function putuskanKomite(
  kasusID: string,
  keputusan: string,
  komentar: string,
): Promise<HasilKeputusanKomite> {
  return minta<HasilKeputusanKomite>(`/api/komite/${encodeURIComponent(kasusID)}/keputusan`, {
    metode: 'POST',
    badan: { keputusan, komentar },
  })
}

/** Eskalasi naik satu tingkat — `POST /api/komite/{id}/eskalasi` (admin). */
export async function eskalasiKomite(
  kasusID: string,
): Promise<{ dariTingkat: number; keTingkat: number }> {
  return minta<{ dariTingkat: number; keTingkat: number }>(
    `/api/komite/${encodeURIComponent(kasusID)}/eskalasi`,
    { metode: 'POST' },
  )
}

/** Satu kasus komite — `GET /api/komite/{id}`. 403 bila bukan anggota tangga. */
export async function ambilKasusKomite(kasusID: string): Promise<KasusKomite> {
  return minta<KasusKomite>(`/api/komite/${encodeURIComponent(kasusID)}`)
}

// ——— Form of Internal Claim Requisition (work owner 12-10-2026, "perbaiki tampilannya lagi") ———
// Layar kasus komite membaca klaim induk BACA SAJA dari rute Claim Life yang sudah ada, dipinjam di `cmd/api/rakit.go`
// (`ruteDipinjam`): modul ini tanpa menu, layarnya hanya terbuka bagi pemegang menu Claim Life - yang memang sudah
// dilayani rute-rute itu. Tipe di bawah SALINAN BAGIAN dari kontrak JSON Claim Life (`modul/claimlife/frontend/api.ts`;
// disalin, bukan diimpor - isolasi modul): hanya medan yang ditampilkan `Section/ShowTransfer.xml`. Nol penulisan.

/** Uang Claim Life: `amount` TEKS desimal, tidak pernah Number. */
export interface UangKlaim {
  amount: string
  currency: string
}

/** Satu baris adjustment klaim induk (medan dialog AdjustmentDetail). */
export interface BarisAdjustmentKlaim {
  id: string
  jumlahKlaim: UangKlaim
  shareNusantaraRe?: UangKlaim
  cedingRetention?: UangKlaim
  sumInsured?: UangKlaim
  sumReasured?: UangKlaim
  shareRetro?: UangKlaim
  retrocededShare?: UangKlaim
  payableTo?: string
  namaBank?: string
  cabangBank?: string
  nomorRekening?: string
}

/** Dokumen peserta. `id` TEKS (17 angka - melewati batas bilangan aman JS); `tanggal` null = kosong. */
export interface DokumenKlaim {
  id: string
  namaFile: string
  kategori2: string
  tanggal: string | null
}

/** Diagnosa peserta (grid DIAGNOSE). */
export interface DiagnosaKlaim {
  id: number
  nama: string
}

/** Peserta klaim induk (PesertaDetail). Tanggal "YYYY-MM-DD[ HH24:MI:SS]". */
export interface PesertaKlaim {
  id: string
  idPega?: string
  nomorPremiList: string
  nomorPolis: string
  nomorSertifikat: string
  namaTertanggung?: string
  umur?: string
  jenisKelamin?: string
  tanggalLahir?: string
  tanggalMulai?: string
  tanggalExpired?: string
  tanggalEfektif?: string
  valuasiMulai?: string
  valuasiSelesai?: string
  tanggalTerimaKlaim: string
  tanggalDokumenLengkap: string
  tanggalKejadian: string
  rekomendasi?: string
  catatanAnalis?: string
  baris: BarisAdjustmentKlaim[]
  dokumen: DokumenKlaim[]
  diagnosa: DiagnosaKlaim[]
}

/** Klaim induk — `GET /api/klaim-life/{id}`. `claimRetro` di tingkat header klaim. */
export interface KlaimInduk {
  id: string
  nomorKlaim: string
  claimRetro?: UangKlaim
  peserta: PesertaKlaim[]
}

export async function ambilKlaimInduk(klaimID: string): Promise<KlaimInduk> {
  return minta<KlaimInduk>(`/api/klaim-life/${encodeURIComponent(klaimID)}`)
}

/** Header polis peserta (Policy Detail) — `GET /api/polis-life/ringkas?idpega=`. 404 = header tidak ditemukan. */
export interface PolisKlaim {
  /** `PolicyDataLife.Type` - Claim Retro / Share Retro / Retroceded Share hanya untuk TP / TR (ShowTransfer b15429). */
  type: string
  cedingCoName: string
  businessName: string
  policyHolderName: string
  sobName: string
}

export async function ambilPolisKlaim(idPega: string): Promise<PolisKlaim> {
  return minta<PolisKlaim>(`/api/polis-life/ringkas?idpega=${encodeURIComponent(idPega)}`)
}

/** Lapisan spreading PremiumList Life peserta (grid Treaty Type | Retroceded Share). */
export interface SpreadingKlaim {
  lapisan: { treatyTypeName: string; retrocededShare: UangKlaim }[]
}

export async function ambilSpreadingKlaim(klaimID: string, pesertaID: string): Promise<SpreadingKlaim> {
  return minta<SpreadingKlaim>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/peserta/${encodeURIComponent(pesertaID)}/spreading-polis`,
  )
}

/** Isi dokumen (View) — `GET /api/klaim-life/{id}/dokumen/{dokId}/isi`; dokumen klaim lain 404. Pola Claim Life. */
export async function ambilIsiDokumenKlaim(klaimID: string, dokID: string): Promise<Blob> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jalur = `/api/klaim-life/${encodeURIComponent(klaimID)}/dokumen/${encodeURIComponent(dokID)}/isi`
    const jawab = await fetch(rakitURL(jalur), {
      method: 'GET',
      headers: { ...headerIdentitas() },
      signal: kendali.signal,
    })
    if (!jawab.ok) throw kegagalanDari(jawab.status, await jawab.text())
    return await jawab.blob()
  } finally {
    clearTimeout(jam)
  }
}
