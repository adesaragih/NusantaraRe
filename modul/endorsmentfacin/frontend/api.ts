// Klien modul Endorsment Fac In - `/api/endorsmentfacin` (`modul/endorsmentfacin/backend/handlers`). Irisan I2
// (keputusan work owner 11-10-2026, W-EDM-01): pemeriksaan sebelum kasus endorsement dibuat - NOL tulis.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_EDMFACIN = '/api/endorsmentfacin'

/** Badan `POST /api/endorsmentfacin/periksa` - isian form Create Addendum yang dibaca kedua pemeriksaan. */
export interface IsianPeriksa {
  /** `.PolicyNo` - teks bebas (pxTextInput, tanpa autocomplete di `Section/WorkPrimaryDetails.xml`). */
  policyNo: string
  /** `.Quotation.EdmType` - kode `DDL\EdmType.xml`: `1`, `2`, `4`. */
  edmType: string
  /** `.Quotation.EdmTypeNew` - hanya bila EdmType `4`. */
  edmTypeNew: string
  /** `.EndorsementDate` - DD-MM-YYYY; kosong = pemeriksaan tanggal belum dijalankan. */
  endorsementDate: string
  /** `.EndorsementInternalRetro` - `0` / `1`. */
  endorsementInternalRetro: string
}

/** Hasil satu gerbang `SetErrorBatalEndorsement_Act` langkah 5. */
export type HasilGerbang = 'lolos' | 'tolak' | 'diklep' | 'dilewati' | 'belum-diperiksa'

export interface PesanMedan {
  /** `PolicyNo` atau `EndorsementDate` - medan tempat Pega memasang pesan. */
  medan: string
  /** Teks pesan verbatim korpus (bahasa Inggris/Indonesia apa adanya). */
  teks: string
}

/** Jawaban `POST /api/endorsmentfacin/periksa`. */
export interface HasilPeriksa {
  /** Ada pesan dari gerbang, pemeriksaan bisnis, atau tanggal. */
  ditolak: boolean
  /** false bila salah satu gerbang `belum-diperiksa` (sumber datanya belum terverifikasi atau tak terjangkau). */
  lengkap: boolean
  pesan: PesanMedan[]
  gerbang: { no: number; hasil: HasilGerbang; catatan: string }[]
  /** Langkah 8-18: jenis bisnis polis (`OutData.pxResults(1).CARI2`); `Life` menyembunyikan sebagian medan. */
  bisnis: { hasil: 'lolos' | 'tolak'; jenisBisnis: string }
  /** `CheckEDMPolisDate`; tanpa tanggal = `null`. */
  tanggal: { hasil: 'lolos' | 'tolak' | 'dilewati'; mulai: string; akhir: string; catatan?: string } | null
}

export function periksaAddendum(isi: IsianPeriksa): Promise<HasilPeriksa> {
  return minta(`${PREFIX_EDMFACIN}/periksa`, { metode: 'POST', badan: isi })
}

/** Badan `POST /api/endorsmentfacin/kasus` - isian pemeriksaan + medan form lain (`WorkPrimaryDetails`). Irisan I3. */
export interface IsianKasus extends IsianPeriksa {
  /** `.Note` - disimpan sebagai `OfferFacIn.EndorsmentReason` (`SetValueToEDMWork` langkah 13). */
  note: string
  /** `.EndorsementSource` - kode tersembunyi autocomplete; teks bebas sementara = kosong. */
  endorsementSource: string
  /** `.EndorsementSourceNote` - teks yang tampil di medan Endorsement Source. */
  endorsementSourceNote: string
  /** `.EndorsementStatus` - belum dirender (prompt list tidak ada di korpus): kosong. */
  endorsementStatus: string
  /** `.Quotation.Type` - Detail EDM Type. */
  detailType: string
  adminCharge: boolean
  survey: boolean
}

/**
 * Jawaban `POST /api/endorsmentfacin/kasus`: 201 = kasus lahir (`caseId` terisi); 200 dengan `caseId` kosong = DITOLAK
 * atau pemeriksaan belum lengkap - hasil bisnis yang sah (spec 07 §7.2), rinciannya di `periksa`. Masukan tidak sah = 422.
 */
export interface HasilKasus {
  caseId: string
  pesan?: PesanMedan[]
  periksa?: HasilPeriksa
}

export function buatKasus(isi: IsianKasus): Promise<HasilKasus> {
  return minta(`${PREFIX_EDMFACIN}/kasus`, { metode: 'POST', badan: isi })
}
