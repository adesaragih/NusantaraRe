// Kasus Komite — tiket 01 (baca) + tiket 02 (keputusan).
//
// Layar keputusan komite = `Section/ShowTransfer.xml` (korpus kelvin/KomiteClaimLife (Done), lewat
// `FlowAction/ViewTransferDtl.xml` b91). Formulir hanya tampil pada GILIRAN pelaku.
//
// ⭐ 12-10-2026 (work owner "perbaiki tampilannya", lalu "perbaiki tampilannya lagi, berikut referensi dari pega"):
//   - kulit layar Claim Life (token `--kcl-*` disalin, tanpa impor CSS claimlife); kepala: Back ghost kecil, judul
//     "Committee Case …", chip giliran;
//   - isi mengikuti URUTAN kotak ShowTransfer di bawah judul "Form of Internal Claim Requisition": Policy Detail, Claim
//     Detail, Medical Advisor and Claim Analyst, Adjustment Detail (+ grid Treaty Type | Retroceded Share), Document Claim
//     (Refresh, View), Submission Date / Remarks, List of Committee, keputusan;
//   - keputusan: "Are you sure to accept this document?" (VERBATIM, label dropdown Approve/Reject) dan `Note` WAJIB
//     (b32893-b32940 `pyRequired true`; server menegakkannya), Cancel kiri, Submit kanan → modal konfirmasi.
//
// ⛔ Data klaim induk BACA SAJA lewat rute Claim Life yang dipinjam (`api.ts`, `cmd/api/rakit.go` `ruteDipinjam`), hanya
// peserta + baris adjustment yang diserahkan ke komite. Nol penulisan baru; aksi komite tidak berubah.
//
// ⚠️ Syarat tampil ShowTransfer yang dibawa: Claim Retro / Share Retro / Retroceded Share hanya untuk polis Type TP / TR
// (b15429, b16998, b17166). `IsTreatyIn==0` (×12) tidak dibawa - Claim Life tidak memuat kasus Treaty In; `Swift Code`
// (tampil bila terisi, b18945) tidak dibaca kontrak JSON Claim Life, jadi tidak tampil.

import { useCallback, useEffect, useState, type ReactNode } from 'react'

import { PERAN, type KodePeran } from '../../../../inti/frontend/labels'
import { simpanBlob } from '../../../../inti/frontend/lib/simpanBlob'
import {
  EFEK_KOMITE,
  ESKALASI_KOMITE,
  FORM_KOMITE,
  KASUS_KOMITE,
  KEPUTUSAN_KOMITE,
  KOLOM_INBOX_KOMITE,
} from '../labels'
import { Gagal, IkonTutup, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilIsiDokumenKlaim,
  ambilKasusKomite,
  ambilKlaimInduk,
  ambilPolisKlaim,
  ambilRiwayatKomite,
  ambilSpreadingKlaim,
  eskalasiKomite,
  putuskanKomite,
  type BarisAdjustmentKlaim,
  type DokumenKlaim,
  type KasusKomite as Kasus,
  type KeputusanAsliKomite,
  type KlaimInduk,
  type PesertaKlaim,
  type PolisKlaim,
  type RiwayatKomite,
  type SpreadingKlaim,
} from '../api'
import {
  kataEfekKomite,
  kataStatusKomite,
  nadaStatusKomite,
  selKomite,
  tanggalKomite,
  teksKomite,
  tingkatKomite,
  uangKomite,
  waktuKomite,
} from '../nilai'

/** Enum tertutup `{1 Approve, 2 Reject}` — urut dropdown. */
export const PILIHAN_KEPUTUSAN = [
  { kode: '1', label: KEPUTUSAN_KOMITE.setuju },
  { kode: '2', label: KEPUTUSAN_KOMITE.tolak },
] as const

/**
 * Keputusan asli tingkat yang tertimpa langkah 5.1 — MURNI (OQ-K-05,
 * GILIRAN-17): `Approved — komentar`, atau tanda kosong bila tidak tertimpa.
 */
export function teksKeputusanAsli(asli: KeputusanAsliKomite | undefined): string {
  if (asli === undefined) return '—'
  const kata = kataStatusKomite(asli.status)
  return asli.comment.trim() === '' ? kata : `${kata} — ${asli.comment}`
}

/** Kolom "Before Final Reject" hanya tampil bila ada tingkat yang tertimpa Reject akhir. */
export function adaKeputusanAsli(riwayat: RiwayatKomite | null): boolean {
  return riwayat !== null && riwayat.tangga.some((a) => a.asli !== undefined)
}

/**
 * Kontrol eskalasi hanya untuk admin, dan hanya bila masih ada tingkat di atas.
 *
 * ⚠️ `[asumsi — OQ-007/OQ-021]` admin komite = `ReasLifeAdmin`. Server tetap
 * yang menegakkan; ini hanya menyembunyikan tombol yang pasti ditolak.
 */
export function bolehEskalasi(peran: readonly KodePeran[], k: Kasus | null): boolean {
  if (k === null) return false
  return peran.includes(PERAN.admin) && k.kasus.tingkatBerjalan > 0 &&
    k.kasus.tingkatBerjalan < k.kasus.komiteLoop
}

/** Kalimat hasil — satu tempat. */
export function kalimatHasilKeputusan(h: {
  kataKeputusan: string
  tingkatDiputus: number
  berlanjut: boolean
  tingkatBerikut: number
  nomorAkseptasi?: string
}): string {
  const dasar = `${kataStatusKomite(h.kataKeputusan)} recorded at level ${String(h.tingkatDiputus)}.`
  if (h.berlanjut) return `${dasar} The case moves up to level ${String(h.tingkatBerikut)}.`
  const nomor = h.nomorAkseptasi ?? ''
  return nomor === '' ? `${dasar} The approval ladder has ended.` : `${dasar} Acceptance No ${nomor}.`
}

/** Chip giliran di kepala: giliran pelaku, tingkat pertama yang masih menunggu, atau tangga selesai. */
export function teksGiliran(k: Kasus, riwayat: RiwayatKomite | null): string {
  if (k.giliranSaya) return KASUS_KOMITE.giliranAnda
  const menunggu = riwayat?.tangga.find((a) => nadaStatusKomite(a.status) === 'menunggu')
  if (menunggu !== undefined) return `${KASUS_KOMITE.menungguTingkat} ${String(menunggu.urut)}`
  if (riwayat === null && k.kasus.tingkatBerjalan > 0 && k.kasus.tingkatBerjalan <= k.kasus.komiteLoop) {
    return `${KASUS_KOMITE.menungguTingkat} ${String(k.kasus.tingkatBerjalan)}`
  }
  return KASUS_KOMITE.tanggaSelesai
}

/** Submit keputusan hanya bila keputusan DAN Note terisi (ShowTransfer b32626 / b32940, keduanya wajib). */
export function bolehKirimKeputusan(keputusan: string, note: string): boolean {
  return keputusan !== '' && note.trim() !== ''
}

/** Peserta + baris adjustment yang DISERAHKAN ke komite (baris dengan id = adjustment kasus); null bila tidak ada. */
export function pesertaDiserahkan(
  klaim: KlaimInduk,
  adjID: string,
): { peserta: PesertaKlaim; baris: BarisAdjustmentKlaim } | null {
  for (const peserta of klaim.peserta) {
    const baris = peserta.baris.find((b) => b.id === adjID)
    if (baris !== undefined) return { peserta, baris }
  }
  return null
}

/** Claim Retro / Share Retro / Retroceded Share — polis Type TP / TR (ShowTransfer b15429); polis tak terbaca = tampil. */
export function tampilRetro(polis: PolisKlaim | null): boolean {
  if (polis === null) return true
  const t = polis.type.trim().toUpperCase()
  return t === 'TP' || t === 'TR'
}

/** Data klaim induk yang sudah dibaca untuk layar ini (BACA SAJA). */
export interface IndukKomite {
  klaim: KlaimInduk
  /** null = baris yang diserahkan tidak ada di klaim induk. */
  diserahkan: { peserta: PesertaKlaim; baris: BarisAdjustmentKlaim } | null
  /** null = tidak terbaca (pesan `polisHilang`). */
  polis: PolisKlaim | null
  /** null = tidak terbaca. */
  spreading: SpreadingKlaim | null
}

/** Chip status tangga (Waiting / Approved / Rejected / …), warna menurut nada. */
export function ChipStatus({ kata }: { kata: string }) {
  return (
    <span className={`komiteclaimlife__chip komiteclaimlife__chip--${nadaStatusKomite(kata)}`}>
      {kataStatusKomite(kata)}
    </span>
  )
}

/** Satu pasangan label-nilai dalam grid kartu. */
function Medan({ label, nilai, penuh, teks }: { label: string; nilai: ReactNode; penuh?: boolean; teks?: boolean }) {
  return (
    <div className={penuh === true ? 'komiteclaimlife__grid-penuh' : undefined}>
      <dt>{label}</dt>
      <dd className={teks === true ? 'komiteclaimlife__teks-panjang' : undefined}>{nilai}</dd>
    </div>
  )
}

/** Kartu berjudul. */
function Kartu({ id, judul, children }: { id: string; judul?: string; children: ReactNode }) {
  return (
    <section
      className="komiteclaimlife__kartu"
      data-kartu={id}
      aria-labelledby={judul === undefined ? undefined : `komiteclaimlife-${id}`}
    >
      {judul !== undefined && (
        <h3 id={`komiteclaimlife-${id}`} className="komiteclaimlife__kartu-judul">
          {judul}
        </h3>
      )}
      {children}
    </section>
  )
}

/** Pita galat di atas layar (pola PitaGalat Claim Life, disalin): kalimat atau objek galat, dapat ditutup. */
function PitaGalat({ galat, onTutup }: { galat: unknown; onTutup: () => void }) {
  if (galat === null || galat === undefined || galat === '') return null
  return (
    <div className="komiteclaimlife__pita">
      {typeof galat === 'string' ? (
        <div role="alert" className="alert alert--error">
          {galat}
        </div>
      ) : (
        <Gagal galat={galat} />
      )}
      <button type="button" className="komiteclaimlife__pita-tutup" aria-label="Dismiss" title="Dismiss" onClick={onTutup}>
        <IkonTutup ukuran={14} />
      </button>
    </div>
  )
}

/** Modal konfirmasi keputusan: keputusan + Note yang akan dikirim, peringatan, Cancel / Yes. */
export function KonfirmasiKeputusan({
  keputusan,
  note,
  kirim,
  onBatal,
  onYa,
}: {
  keputusan: string
  note: string
  kirim: boolean
  onBatal: () => void
  onYa: () => void
}) {
  const label = PILIHAN_KEPUTUSAN.find((p) => p.kode === keputusan)?.label ?? '—'
  return (
    <Modal
      judul={KEPUTUSAN_KOMITE.judulKonfirmasi}
      onTutup={onBatal}
      tanpaTutup={kirim}
      labelBatal={KEPUTUSAN_KOMITE.cancel}
      aksi={
        <button type="button" className="btn btn--primary" disabled={kirim} onClick={onYa}>
          {kirim ? KEPUTUSAN_KOMITE.mengirim : KEPUTUSAN_KOMITE.ya}
        </button>
      }
    >
      <div className="komiteclaimlife__konfirmasi">
        <p className="komiteclaimlife__tanya">{KEPUTUSAN_KOMITE.tanyaKonfirmasi}</p>
        <dl className="komiteclaimlife__grid">
          <Medan label={KEPUTUSAN_KOMITE.keputusan} nilai={<strong>{label}</strong>} />
          <Medan label={KEPUTUSAN_KOMITE.komentar} nilai={teksKomite(note)} penuh teks />
        </dl>
        <p role="note" className="alert alert--warn">
          {KEPUTUSAN_KOMITE.peringatan}
        </p>
      </div>
    </Modal>
  )
}

/** Kotak-kotak klaim induk ShowTransfer (Policy Detail … Document Claim) — MURNI dari props. */
export function FormKlaimKomite({
  induk,
  memuat,
  galat,
  onRefresh,
  onLihat,
}: {
  induk: IndukKomite | null
  memuat: boolean
  galat: unknown
  onRefresh: () => void
  onLihat: (d: DokumenKlaim) => void
}) {
  if (induk === null) {
    if (galat !== null && galat !== undefined) {
      return (
        <Kartu id="klaim-induk">
          <p role="alert" className="alert alert--warn">
            {FORM_KOMITE.gagalKlaim}
          </p>
          <Gagal galat={galat} />
        </Kartu>
      )
    }
    return memuat ? <Memuat pesan={FORM_KOMITE.memuatKlaim} /> : null
  }
  if (induk.diserahkan === null) {
    return (
      <Kartu id="klaim-induk">
        <p role="alert" className="alert alert--warn">
          {FORM_KOMITE.pesertaHilang}
        </p>
      </Kartu>
    )
  }
  const { peserta: p, baris: b } = induk.diserahkan
  const polis = induk.polis
  const retro = tampilRetro(polis)
  const lapisan = induk.spreading?.lapisan ?? []
  return (
    <>
      <Kartu id="polis" judul={FORM_KOMITE.polis}>
        {polis === null && (
          <p role="note" className="alert alert--warn">
            {FORM_KOMITE.polisHilang}
          </p>
        )}
        <div className="komiteclaimlife__kolom komiteclaimlife__kolom--2">
          <dl className="komiteclaimlife__grid komiteclaimlife__grid--tumpuk">
            <Medan label={FORM_KOMITE.cedingCompany} nilai={teksKomite(polis?.cedingCoName)} />
            <Medan label={FORM_KOMITE.plan} nilai={teksKomite(polis?.businessName)} />
            <Medan label={FORM_KOMITE.plNumber} nilai={teksKomite(p.nomorPremiList)} />
          </dl>
          <dl className="komiteclaimlife__grid komiteclaimlife__grid--tumpuk">
            <Medan label={FORM_KOMITE.policyHolder} nilai={teksKomite(polis?.policyHolderName)} />
            <Medan label={FORM_KOMITE.sobName} nilai={teksKomite(polis?.sobName)} />
          </dl>
        </div>
      </Kartu>
      <Kartu id="klaim" judul={FORM_KOMITE.klaim}>
        <div className="komiteclaimlife__kolom komiteclaimlife__kolom--3">
          <dl className="komiteclaimlife__grid komiteclaimlife__grid--tumpuk">
            <Medan label={FORM_KOMITE.namaTertanggung} nilai={teksKomite(p.namaTertanggung)} />
            <Medan label={FORM_KOMITE.sertifikat} nilai={teksKomite(p.nomorSertifikat)} />
            <Medan label={FORM_KOMITE.nomorPolis} nilai={teksKomite(p.nomorPolis)} />
            <Medan label={FORM_KOMITE.umur} nilai={teksKomite(p.umur)} />
            <Medan label={FORM_KOMITE.jenisKelamin} nilai={teksKomite(p.jenisKelamin)} />
          </dl>
          <dl className="komiteclaimlife__grid komiteclaimlife__grid--tumpuk">
            <Medan label={FORM_KOMITE.tanggalLahir} nilai={tanggalKomite(p.tanggalLahir)} />
            <Medan label={FORM_KOMITE.mulai} nilai={tanggalKomite(p.tanggalMulai)} />
            <Medan label={FORM_KOMITE.expired} nilai={tanggalKomite(p.tanggalExpired)} />
            <Medan label={FORM_KOMITE.efektif} nilai={tanggalKomite(p.tanggalEfektif)} />
            <Medan label={FORM_KOMITE.valuasiMulai} nilai={tanggalKomite(p.valuasiMulai)} />
            <Medan label={FORM_KOMITE.valuasiSelesai} nilai={tanggalKomite(p.valuasiSelesai)} />
          </dl>
          <dl className="komiteclaimlife__grid komiteclaimlife__grid--tumpuk">
            <Medan label={FORM_KOMITE.terimaKlaim} nilai={tanggalKomite(p.tanggalTerimaKlaim)} />
            <Medan label={FORM_KOMITE.dokumenLengkap} nilai={tanggalKomite(p.tanggalDokumenLengkap)} />
            <Medan label={FORM_KOMITE.kejadian} nilai={tanggalKomite(p.tanggalKejadian)} />
          </dl>
        </div>
      </Kartu>
      <Kartu id="medis" judul={FORM_KOMITE.medis}>
        <div className="komiteclaimlife__gulir">
          <table className="komiteclaimlife__tabel komiteclaimlife__tabel--sempit">
            <thead>
              <tr>
                <th>{FORM_KOMITE.diagnosa}</th>
              </tr>
            </thead>
            <tbody>
              {p.diagnosa.length === 0 && (
                <tr>
                  <td className="komiteclaimlife__redup">{FORM_KOMITE.tanpaItem}</td>
                </tr>
              )}
              {p.diagnosa.map((d) => (
                <tr key={d.id}>
                  <td>{teksKomite(d.nama)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <dl className="komiteclaimlife__grid">
          <Medan label={FORM_KOMITE.rekomendasi} nilai={teksKomite(p.rekomendasi)} penuh teks />
          <Medan label={FORM_KOMITE.analis} nilai={teksKomite(p.catatanAnalis)} penuh teks />
        </dl>
      </Kartu>
      <Kartu id="adjustment" judul={FORM_KOMITE.adjustment}>
        <dl className="komiteclaimlife__grid komiteclaimlife__grid--4">
          <Medan label={FORM_KOMITE.mataUang} nilai={teksKomite(b.jumlahKlaim.currency)} />
        </dl>
        <dl className="komiteclaimlife__grid komiteclaimlife__grid--4">
          <Medan label={FORM_KOMITE.shareRnm} nilai={uangKomite(b.shareNusantaraRe?.amount)} />
          <Medan label={FORM_KOMITE.cedingRetention} nilai={uangKomite(b.cedingRetention?.amount)} />
          <Medan label={FORM_KOMITE.claimRnm} nilai={uangKomite(b.jumlahKlaim.amount)} />
          {retro && <Medan label={FORM_KOMITE.claimRetro} nilai={uangKomite(induk.klaim.claimRetro?.amount)} />}
        </dl>
        <dl className="komiteclaimlife__grid komiteclaimlife__grid--4">
          <Medan label={FORM_KOMITE.sumReinsured} nilai={uangKomite(b.sumReasured?.amount)} />
          <Medan label={FORM_KOMITE.sumInsured} nilai={uangKomite(b.sumInsured?.amount)} />
          {retro && <Medan label={FORM_KOMITE.shareRetro} nilai={uangKomite(b.shareRetro?.amount)} />}
          {retro && <Medan label={FORM_KOMITE.retrocededShare} nilai={uangKomite(b.retrocededShare?.amount)} />}
        </dl>
        <dl className="komiteclaimlife__grid komiteclaimlife__grid--4">
          <Medan label={FORM_KOMITE.payableTo} nilai={teksKomite(b.payableTo)} />
          <Medan label={FORM_KOMITE.namaBank} nilai={teksKomite(b.namaBank)} />
          <Medan label={FORM_KOMITE.cabangBank} nilai={teksKomite(b.cabangBank)} />
          <Medan label={FORM_KOMITE.nomorRekening} nilai={teksKomite(b.nomorRekening)} />
        </dl>
        <div className="komiteclaimlife__gulir">
          <table className="komiteclaimlife__tabel komiteclaimlife__tabel--sempit">
            <thead>
              <tr>
                <th>{FORM_KOMITE.treatyType}</th>
                <th className="komiteclaimlife__kanan">{FORM_KOMITE.retrocededShare}</th>
              </tr>
            </thead>
            <tbody>
              {lapisan.length === 0 && (
                <tr>
                  <td colSpan={2} className="komiteclaimlife__redup">
                    {FORM_KOMITE.tanpaItem}
                  </td>
                </tr>
              )}
              {lapisan.map((l, i) => (
                <tr key={`${l.treatyTypeName}:${String(i)}`}>
                  <td>{teksKomite(l.treatyTypeName)}</td>
                  <td className="komiteclaimlife__kanan">{uangKomite(l.retrocededShare.amount)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Kartu>
      <Kartu id="dokumen" judul={FORM_KOMITE.dokumen}>
        <div className="komiteclaimlife__aksi komiteclaimlife__aksi--kiri">
          <button type="button" className="btn btn--ghost btn--sm" disabled={memuat} onClick={onRefresh}>
            {FORM_KOMITE.refresh}
          </button>
        </div>
        <div className="komiteclaimlife__gulir">
          <table className="komiteclaimlife__tabel">
            <thead>
              <tr>
                <th>{FORM_KOMITE.namaFile}</th>
                <th>{FORM_KOMITE.jenisLampiran}</th>
                <th>{FORM_KOMITE.tanggalDokumen}</th>
              </tr>
            </thead>
            <tbody>
              {p.dokumen.length === 0 && (
                <tr>
                  <td colSpan={3} className="komiteclaimlife__redup">
                    {FORM_KOMITE.tanpaItem}
                  </td>
                </tr>
              )}
              {p.dokumen.map((d) => (
                <tr key={d.id}>
                  <td className="komiteclaimlife__teks">
                    <button
                      type="button"
                      className="komiteclaimlife__tautan"
                      title={`${FORM_KOMITE.lihat} ${d.namaFile}`}
                      onClick={() => onLihat(d)}
                    >
                      {teksKomite(d.namaFile)}
                    </button>
                  </td>
                  <td>{teksKomite(d.kategori2)}</td>
                  <td className="komiteclaimlife__waktu">{waktuKomite(d.tanggal ?? '')}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Kartu>
    </>
  )
}

/** Isi layar — MURNI dari props (diuji tanpa jaringan); keadaan dan panggilan API di `KasusKomite`. */
export function TampilanKasusKomite({
  kasusID,
  k,
  riwayat,
  induk,
  memuatInduk,
  galatInduk,
  peran,
  sibuk,
  kirim,
  galat,
  kabar,
  keputusan,
  komentar,
  konfirmasi,
  onKembali,
  onTutupGalat,
  onKeputusan,
  onKomentar,
  onMintaKonfirmasi,
  onBatalKonfirmasi,
  onSubmit,
  onEskalasi,
  onRefresh,
  onLihat,
}: {
  kasusID: string
  k: Kasus | null
  riwayat: RiwayatKomite | null
  induk: IndukKomite | null
  memuatInduk: boolean
  galatInduk: unknown
  peran: readonly KodePeran[]
  sibuk: boolean
  kirim: boolean
  galat: unknown
  kabar: string
  keputusan: string
  komentar: string
  konfirmasi: boolean
  onKembali: () => void
  onTutupGalat: () => void
  onKeputusan: (v: string) => void
  onKomentar: (v: string) => void
  onMintaKonfirmasi: () => void
  onBatalKonfirmasi: () => void
  onSubmit: () => void
  onEskalasi: () => void
  onRefresh: () => void
  onLihat: (d: DokumenKlaim) => void
}) {
  const kolomAsli = adaKeputusanAsli(riwayat)
  return (
    <section className="komiteclaimlife__akar komite-kasus" aria-labelledby="komiteclaimlife-judul">
      <header className="komiteclaimlife__kepala">
        <button type="button" className="btn btn--ghost btn--sm komiteclaimlife__kembali" onClick={onKembali}>
          {KASUS_KOMITE.kembali}
        </button>
        <div className="komiteclaimlife__judul-baris">
          <h2 id="komiteclaimlife-judul" className="komiteclaimlife__judul">
            {KASUS_KOMITE.judul} {kasusID}
          </h2>
          {k !== null && (
            <span
              role="status"
              className={`komiteclaimlife__chip komiteclaimlife__chip--${k.giliranSaya ? 'giliran' : 'lain'}`}
            >
              {teksGiliran(k, riwayat)}
            </span>
          )}
        </div>
        {k !== null && (
          <dl className="komiteclaimlife__meta komite-kasus__kepala">
            <Medan label={KOLOM_INBOX_KOMITE.nomorKlaim} nilai={selKomite(k.kasus.nomorKlaim)} />
            <Medan label={KOLOM_INBOX_KOMITE.tingkat} nilai={tingkatKomite(k.kasus)} />
            <Medan
              label={KOLOM_INBOX_KOMITE.statusBaris}
              nilai={<span className="komiteclaimlife__chip komiteclaimlife__chip--lain">{selKomite(k.kasus.statusBaris)}</span>}
            />
          </dl>
        )}
      </header>
      <PitaGalat galat={galat} onTutup={onTutupGalat} />
      {sibuk && <Memuat />}
      {kabar !== '' && (
        <div role="status" className="alert alert--ok">
          {kabar}
        </div>
      )}
      {k !== null && (
        <div className="komiteclaimlife__kartu-kartu">
          <h3 className="komiteclaimlife__judul-form">{FORM_KOMITE.judul}</h3>
          <FormKlaimKomite induk={induk} memuat={memuatInduk} galat={galatInduk} onRefresh={onRefresh} onLihat={onLihat} />
          {/* Remarks penyerahan dari Claim Life (12-10-2026, opsi A) — baca saja, dari jejak klaim induk. */}
          {riwayat?.penyerahan !== undefined && (
            <section className="komiteclaimlife__kartu komite-kasus__penyerahan" data-kartu="penyerahan" aria-label={KASUS_KOMITE.penyerahan}>
              <dl className="komiteclaimlife__grid">
                <Medan label={KASUS_KOMITE.penyerahanTanggal} nilai={waktuKomite(riwayat.penyerahan.tanggal)} />
                <div className="komiteclaimlife__grid-penuh">
                  <dt>{KASUS_KOMITE.penyerahanRemarks}</dt>
                  <dd className="komite-kasus__remarks">{riwayat.penyerahan.remarks || '—'}</dd>
                </div>
              </dl>
            </section>
          )}
          <Kartu id="tangga" judul={KASUS_KOMITE.tangga}>
            {riwayat !== null && (
              <div className="komiteclaimlife__gulir">
                <table className="komiteclaimlife__tabel komite-kasus__tangga">
                  <thead>
                    <tr>
                      <th>{KASUS_KOMITE.committee}</th>
                      <th>{KASUS_KOMITE.status}</th>
                      <th>{KASUS_KOMITE.dateApprove}</th>
                      <th>{KASUS_KOMITE.comment}</th>
                      {kolomAsli && <th>{KASUS_KOMITE.asli}</th>}
                    </tr>
                  </thead>
                  <tbody>
                    {riwayat.tangga.length === 0 && (
                      <tr>
                        <td colSpan={kolomAsli ? 5 : 4}>{KASUS_KOMITE.tanpaTangga}</td>
                      </tr>
                    )}
                    {riwayat.tangga.map((a) => (
                      <tr key={a.urut}>
                        <td>{selKomite(a.committee)}</td>
                        <td>
                          <ChipStatus kata={a.status} />
                        </td>
                        <td className="komiteclaimlife__waktu">{waktuKomite(a.dateApprove)}</td>
                        <td className="komiteclaimlife__teks">{selKomite(a.comment)}</td>
                        {kolomAsli && <td className="komiteclaimlife__teks">{teksKeputusanAsli(a.asli)}</td>}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
            {riwayat !== null && riwayat.eskalasi.length > 0 && (
              <>
                <h4 className="komiteclaimlife__sub-judul">{KASUS_KOMITE.eskalasi}</h4>
                <ul className="komite-kasus__eskalasi-riwayat">
                  {riwayat.eskalasi.map((e, i) => (
                    <li key={`${e.waktu}-${String(i)}`}>
                      Level {e.dariTingkat} → {e.keTingkat} by {e.oleh} on {waktuKomite(e.waktu)}
                    </li>
                  ))}
                </ul>
              </>
            )}
            {bolehEskalasi(peran, k) && (
              <p className="komite-kasus__eskalasi">
                <button type="button" className="btn btn--ghost btn--sm" disabled={kirim} onClick={onEskalasi}>
                  {ESKALASI_KOMITE.tombol}
                </button>{' '}
                <span className="komiteclaimlife__redup">{ESKALASI_KOMITE.keterangan}</span>
              </p>
            )}
          </Kartu>
          {k.efek.efek.length > 0 && (
            <Kartu id="efek" judul={`${EFEK_KOMITE.judul} — ${kataEfekKomite(k.efek.keadaan)}`}>
              <div className="komiteclaimlife__gulir">
                <table className="komiteclaimlife__tabel komite-kasus__efek">
                  <thead>
                    <tr>
                      <th>{EFEK_KOMITE.jenis}</th>
                      <th>{EFEK_KOMITE.keadaan}</th>
                      <th className="komiteclaimlife__tengah">{EFEK_KOMITE.percobaan}</th>
                      <th>{EFEK_KOMITE.sejak}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {k.efek.efek.map((e, i) => (
                      <tr key={`${e.jenis}-${String(i)}`}>
                        <td>{e.jenis}</td>
                        <td>{kataEfekKomite(e.keadaan)}</td>
                        <td className="komiteclaimlife__tengah">{e.percobaan}</td>
                        <td className="komiteclaimlife__waktu">{waktuKomite(e.sejak)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </Kartu>
          )}
          {k.giliranSaya && (
            <Kartu id="keputusan" judul={KEPUTUSAN_KOMITE.judul}>
              <form
                className="komite-kasus__keputusan komiteclaimlife__form"
                onSubmit={(e) => {
                  e.preventDefault()
                  onMintaKonfirmasi()
                }}
              >
                <label className="komiteclaimlife__isian">
                  <span className="komiteclaimlife__label">
                    {KEPUTUSAN_KOMITE.konfirmasi} <span aria-hidden="true">*</span>
                  </span>
                  <select
                    required
                    className="field__input"
                    value={keputusan}
                    disabled={kirim}
                    onChange={(e) => {
                      onKeputusan(e.target.value)
                    }}
                  >
                    <option value="">{KEPUTUSAN_KOMITE.pilih}</option>
                    {PILIHAN_KEPUTUSAN.map((p) => (
                      <option key={p.kode} value={p.kode}>
                        {p.label}
                      </option>
                    ))}
                  </select>
                </label>
                <label className="komiteclaimlife__isian komiteclaimlife__isian--penuh">
                  <span className="komiteclaimlife__label">
                    {KEPUTUSAN_KOMITE.komentar} <span aria-hidden="true">*</span>
                  </span>
                  <textarea
                    required
                    className="field__input"
                    rows={4}
                    value={komentar}
                    disabled={kirim}
                    onChange={(e) => {
                      onKomentar(e.target.value)
                    }}
                  />
                </label>
                <div className="komiteclaimlife__aksi komiteclaimlife__aksi--antara">
                  <button type="button" className="btn btn--ghost" disabled={kirim} onClick={onKembali}>
                    {KEPUTUSAN_KOMITE.cancel}
                  </button>
                  <button
                    type="submit"
                    className="btn btn--primary"
                    disabled={kirim || !bolehKirimKeputusan(keputusan, komentar)}
                  >
                    {kirim ? KEPUTUSAN_KOMITE.mengirim : KEPUTUSAN_KOMITE.submit}
                  </button>
                </div>
              </form>
              {konfirmasi && (
                <KonfirmasiKeputusan
                  keputusan={keputusan}
                  note={komentar}
                  kirim={kirim}
                  onBatal={onBatalKonfirmasi}
                  onYa={onSubmit}
                />
              )}
            </Kartu>
          )}
        </div>
      )}
    </section>
  )
}

/** Baca klaim induk BACA SAJA: klaim, lalu polis + spreading peserta yang diserahkan (gagal = null, bukan galat layar). */
async function bacaInduk(klaimID: string, adjID: string): Promise<IndukKomite> {
  const klaim = await ambilKlaimInduk(klaimID)
  const diserahkan = pesertaDiserahkan(klaim, adjID)
  if (diserahkan === null) return { klaim, diserahkan, polis: null, spreading: null }
  const idPega = (diserahkan.peserta.idPega ?? '').trim()
  const [polis, spreading] = await Promise.all([
    idPega === '' ? Promise.resolve(null) : ambilPolisKlaim(idPega).catch(() => null),
    ambilSpreadingKlaim(klaimID, diserahkan.peserta.id).catch(() => null),
  ])
  return { klaim, diserahkan, polis, spreading }
}

export default function KasusKomite({
  kasusID,
  peran,
  onKembali,
}: {
  kasusID: string
  /** Peran pelaku — hanya untuk menampilkan kontrol eskalasi. */
  peran: readonly KodePeran[]
  onKembali: () => void
}) {
  const [k, setK] = useState<Kasus | null>(null)
  const [riwayat, setRiwayat] = useState<RiwayatKomite | null>(null)
  const [induk, setInduk] = useState<IndukKomite | null>(null)
  const [memuatInduk, setMemuatInduk] = useState(false)
  const [galatInduk, setGalatInduk] = useState<unknown>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)
  const [keputusan, setKeputusan] = useState('')
  const [komentar, setKomentar] = useState('')
  const [kirim, setKirim] = useState(false)
  const [kabar, setKabar] = useState('')
  const [konfirmasi, setKonfirmasi] = useState(false)

  const muat = useCallback(async () => {
    setSibuk(true)
    setGalat(null)
    try {
      const [kasus, rw] = await Promise.all([ambilKasusKomite(kasusID), ambilRiwayatKomite(kasusID)])
      setK(kasus)
      setRiwayat(rw)
    } catch (e) {
      // ⚠️ 403 bila pelaku bukan anggota tangga — pesan server tampil apa adanya.
      setGalat(e)
      setK(null)
    } finally {
      setSibuk(false)
    }
  }, [kasusID])

  useEffect(() => {
    void muat()
  }, [muat])

  const klaimID = k?.kasus.klaimId ?? ''
  const adjID = k?.adjustmentId ?? ''
  const muatInduk = useCallback(async () => {
    if (klaimID === '') return
    setMemuatInduk(true)
    setGalatInduk(null)
    try {
      setInduk(await bacaInduk(klaimID, adjID))
    } catch (e) {
      setInduk(null)
      setGalatInduk(e)
    } finally {
      setMemuatInduk(false)
    }
  }, [klaimID, adjID])

  useEffect(() => {
    void muatInduk()
  }, [muatInduk])

  async function lihat(d: DokumenKlaim): Promise<void> {
    try {
      simpanBlob(await ambilIsiDokumenKlaim(klaimID, d.id), d.namaFile)
    } catch (e) {
      setGalat(e)
    }
  }

  async function eskalasi(): Promise<void> {
    if (kirim) return
    setKirim(true)
    setGalat(null)
    setKabar('')
    try {
      const h = await eskalasiKomite(kasusID)
      setKabar(`Escalated from level ${String(h.dariTingkat)} to level ${String(h.keTingkat)}.`)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setKirim(false)
    }
  }

  function mintaKonfirmasi(): void {
    if (kirim) return
    if (keputusan === '') {
      setGalat(KEPUTUSAN_KOMITE.wajibPilih)
      return
    }
    if (komentar.trim() === '') {
      setGalat(KEPUTUSAN_KOMITE.wajibNote)
      return
    }
    setGalat(null)
    setKonfirmasi(true)
  }

  async function submit(): Promise<void> {
    if (kirim) return
    setKirim(true)
    setGalat(null)
    setKabar('')
    try {
      const h = await putuskanKomite(kasusID, keputusan, komentar)
      setKonfirmasi(false)
      setKabar(
        kalimatHasilKeputusan(h) +
          (h.efekTertunda.length > 0
            ? ` Outgoing effects queued, not yet complete: ${h.efekTertunda.join(', ')}.`
            : ''),
      )
      setKeputusan('')
      setKomentar('')
      await muat()
    } catch (e) {
      // Galat di pita atas layar; modal ditutup supaya pita terlihat.
      setKonfirmasi(false)
      setGalat(e)
    } finally {
      setKirim(false)
    }
  }

  return (
    <TampilanKasusKomite
      kasusID={kasusID}
      k={k}
      riwayat={riwayat}
      induk={induk}
      memuatInduk={memuatInduk}
      galatInduk={galatInduk}
      peran={peran}
      sibuk={sibuk}
      kirim={kirim}
      galat={galat}
      kabar={kabar}
      keputusan={keputusan}
      komentar={komentar}
      konfirmasi={konfirmasi}
      onKembali={onKembali}
      onTutupGalat={() => setGalat(null)}
      onKeputusan={setKeputusan}
      onKomentar={setKomentar}
      onMintaKonfirmasi={mintaKonfirmasi}
      onBatalKonfirmasi={() => setKonfirmasi(false)}
      onSubmit={() => void submit()}
      onEskalasi={() => void eskalasi()}
      onRefresh={() => void muatInduk()}
      onLihat={(d) => void lihat(d)}
    />
  )
}
