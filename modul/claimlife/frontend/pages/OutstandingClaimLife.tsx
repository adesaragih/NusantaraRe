// Layar Outstanding Claim — `FlowAction/OSClaimLife.xml` → `Section/InputOSClaimLife.xml`.
//
// ⭐ DIROMBAK 10-10-2026 mengikuti layar Pega (permintaan work owner "kenapa tidak ada detail peserta nya seperti yang
// di pega?", screenshot CLMLF-9293). Korpus: `kelvin/InboxClaimLife (Done)/Section/InputOSClaimLife.xml`; bNNN = nomor
// baris berkas itu.
//
//   Kepala   "Claim Life (<Case ID>)" + Claim No `.ClaimData.PremiumListSummary.CLAIM_NO` b958 (baca saja; kosong
//            sampai Save to RNM).
//   Tab      b1607: Register Claim b1956 · Policy Details b5909 · Claim Details b14576.
//   Register PL Number b3054 · Catastrophe b3475 (radio; nonaktif bila CLAIM_NO terisi b3515) · Is KPR b3781 ·
//            STNC Treaty b4726 (baca saja) · Catastrophe Note b4905 (tampil bila Catastrophe; baca saja bila CLAIM_NO
//            terisi b5059).
//   Polis    PanelPolisOutstanding tiga kolom. View Detail (DetailPolisLife b13267) dicabut 11-10-2026 (keputusan
//            work owner); rute baca GET …/peserta-polis tetap ada tanpa pemakai di layar.
//   Klaim    Participant Details b15524: grid `PremiumListDetail` b15924 — PL NUMBER b16101, NAME OF INSURED b16256,
//            DATE OF BIRTH b16406, Accept `.IsAccept` b17582 (kepala SELECT ALL b16609 → SelectAllClaimLife_act),
//            DELETE b17865 (tampil bila CLAIM_NO kosong b18082). Baris dibuka = ViewClaimDetailLifeGCNM b18252 →
//            RincianPeserta. Save b19413.
//   Bawah    Save to RNM b21102 (nonaktif sesudah tersimpan b21095) · Send Back to Register b21404 (tampil bila CLAIM_NO
//            kosong b21699) · Send to Medical Check b21839 (nonaktif bila CLAIM_NO kosong b21832; tampil bila QP/QR
//            b22123) · Close Claim b22808.
//
// ⛔ Setiap kondisi tampil/nonaktif di sini DITEGAKKAN JUGA DI SERVER (Pindah 409, Save 409); layar hanya kenyamanan.
//
// ⛔ BUTIR aw - dua tombol perpindahan. Di Pega `SendtoAdmin_Act` dan `SendtoAdmin_Act1` sama-sama berprasyarat
// `pyPosition=="ReasLifeMedicalAdvisor"` (WhenFalse=3), padahal layar ini dipegang `ReasLifeAdmin`. Cacat rule warisan;
// yang ditiru MAKSUDnya - lihat `OQ-untuk-tim.md` OQ-C. Daftarnya satu sumber: `PanelPindahTahap`.
//
// ⛔ Tombol "Detail klaim" lama DIHAPUS: seluruh isinya (tanggal, total, dokumen, diagnosa, adjustment, Reject, Close
// Claim) kini di tab Claim Details lewat komponen bersama.
//
// ⭐ 11-10-2026 (work owner): layar ini dipakai SEMUA tahap yang dibuka dari Inbox — pita kuning peringatan tahap
// dan tombol ke layar Detail DIHAPUS. Tata letak sama; yang berbeda hanya isian dan tombol bawah:
//   Outstanding   isian Admin + Save to RNM / Send to Medical Check (di atas).
//   Medical Check `FlowAction/MedicalCheckClaimLife` — data klaim BACA SAJA; bagian Medical Advisor (diagnosa,
//                 Recommendation) diisi Medical Advisor; tombol Save · Send Back to Admin · Send to Claim Analyst.
//   Claim Analis  Send Back to Admin · Send Back to Medical Check · Close Claim (PanelTutupKlaim).

import { useCallback, useEffect, useState, type ReactNode } from 'react'

import { OS_PEGA, TAHAP, TOMBOL_MEDIS, TOMBOL_OS } from '../labels'
import { PanelPolisOutstanding, tanggalTampil } from '../components/PanelDataPolis'
import { bacaPolisPeserta } from '../pilihanPeserta'
import { PanelPindahTahap, type TombolPindah } from '../components/PanelPindahTahap'
import { PanelTutupKlaim } from '../components/PanelTutupKlaim'
import {
  claimRNMKosong,
  galatClaimRNM,
  galatDOL,
  isianBarisAwal,
  KOLOM_CLAIM_RNM,
  PESAN_CLAIM_RNM_BUKAN_ANGKA,
  RincianPeserta,
  tanggalAwal,
} from '../components/RincianPeserta'
import { Modal, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { JUDUL_KEKURANGAN_MEDIS, PitaGalat } from '../components/PitaGalat'
import { PERAN, type KodePeran } from '../../../../inti/frontend/labels'
import { jungkitPilihSemua, type KeadaanPilihSemua } from '../../../../inti/frontend/lib/pilihSemua'
import {
  ambilDataPolisIdPega,
  ambilKlaimLife,
  aturanTombolOutstanding,
  bernomor,
  bolehDeletePeserta,
  bolehIsianOutstanding,
  cabutPeserta,
  GagalSimpanRNM,
  kasusTertutup,
  simpanMedis,
  pindahTahap,
  PILIHAN_KATASTROFE,
  PILIHAN_KPR,
  simpanIsianOutstanding,
  simpanKeRNM,
  TAHAP_JALUR,
  type TahapJalur,
  type HasilSimpanRNM,
  type IsianOutstanding,
  type KekuranganRNM,
  type Klaim,
  type PolicyDataLife,
} from '../api'

export interface OutstandingProps {
  /** Peran pelaku (sesi) — Admin Notes hanya dapat diisi Admin (11-10-2026). Tanpa ini: baca saja. */
  peran?: readonly KodePeran[]
  /** Pengenal work kasus yang sedang dibuka. */
  klaimID: string
  /** Dipanggil sesudah kasus berpindah / ditutup — layar kembali ke Inbox; `kabar` = pita sukses di Inbox (11-10-2026). */
  onPindah: (kabar?: string) => void
}

/** Tab layar, urut b1607. */
export const TAB_OS = ['register', 'polis', 'klaim'] as const
export type TabOS = (typeof TAB_OS)[number]
export function labelTabOS(t: TabOS): string {
  return t === 'register' ? OS_PEGA.tabRegister : t === 'polis' ? OS_PEGA.tabPolis : OS_PEGA.tabKlaim
}

/** Kepala kolom grid Participant Details, urut b15924. */
export const KOLOM_PESERTA = ['No.', OS_PEGA.kolomPL, OS_PEGA.kolomNama, OS_PEGA.kolomLahir, OS_PEGA.accept, OS_PEGA.hapus] as const

/** Saringan nama tertanggung grid (11-10-2026, "dengan saringan bila mudah"): huruf besar-kecil diabaikan. */
export function saringPeserta<T extends { namaTertanggung?: string }>(peserta: readonly T[], kata: string): T[] {
  const k = kata.trim().toUpperCase()
  return k === '' ? [...peserta] : peserta.filter((p) => (p.namaTertanggung ?? '').toUpperCase().includes(k))
}

/** Isian awal tombol Save dari klaim tersimpan; Catastrophe kosong = bawaan Pega "Non-Catastrophe" (b3728). */
export function isianAwal(klaim: Klaim): IsianOutstanding {
  const accept: Record<string, string> = {}
  for (const p of klaim.peserta) accept[p.id] = p.isAccept === 'true' ? 'true' : 'false'
  return {
    stsKatastrofe: klaim.stsKatastrofe !== undefined && klaim.stsKatastrofe !== '' ? klaim.stsKatastrofe : 'Non-Catastrophe',
    // Bawaan Non-KPR untuk kolom kosong (keputusan work owner 11-10-2026); ikut terkirim pada Save grid. Nilai warisan asing
    // tetap apa adanya ("Stored value").
    isKpr: klaim.isKpr !== undefined && klaim.isKpr !== '' ? klaim.isKpr : 'Non-KPR',
    katastrofeNote: klaim.katastrofeNote ?? '',
    accept,
    // Admin Notes dan tanggal hanya dikirim untuk peserta yang diubah (11-10-2026).
    catatanAdmin: {},
    tanggal: {},
    adjustment: {},
  }
}

/**
 * Isian yang TIDAK DAPAT DISIMPAN (format rusak) — hanya itu yang menonaktifkan Save grid. Keputusan work owner 11-10-2026
 * (terbaru): "hanya save progress saja" — DOL di luar valuation, Claim RNM > Share, dan kolom wajib kosong TETAP tampil
 * merah sebagai peringatan, tetapi Save tetap aktif dan menyimpan; penegakannya di Save to RNM. Tanggal berasal dari
 * <input type="date"> (selalu sah atau kosong); teks lebar dibatasi maxLength / server 400.
 */
export function adaIsianTidakSah(klaim: Klaim, isian: IsianOutstanding): boolean {
  for (const [adjID, isi] of Object.entries(isian.adjustment ?? {})) {
    const baris = klaim.peserta.flatMap((p) => p.baris).find((b) => b.id === adjID)
    if (baris !== undefined && galatClaimRNM(isi.jumlah, baris.shareNusantaraRe) === PESAN_CLAIM_RNM_BUKAN_ANGKA) return true
  }
  return false
}

/**
 * Pra-periksa Send to Medical Check (keputusan work owner 11-10-2026): setiap peserta yang di-Accept (TERSIMPAN) punya >= 1
 * dokumen — data yang dimuat; dokumen berstatus apa pun dihitung. Kalimat SAMA dengan server
 * (services.PeriksaDokumenAcceptMedis); server tetap lapis terakhir (422 berdaftar).
 */
export function kekuranganDokumenMedis(klaim: Klaim): KekuranganRNM[] {
  const out: KekuranganRNM[] = []
  klaim.peserta.forEach((p, i) => {
    if (p.isAccept !== 'true' || p.dokumen.length > 0) return
    const no = i + 1
    const nama = (p.namaTertanggung ?? '').trim() || p.nomorSertifikat
    out.push({
      nomorPeserta: no,
      pesertaId: p.id,
      sertifikat: p.nomorSertifikat,
      nama: p.namaTertanggung ?? '',
      kolom: 'Document',
      pesan: `Insured number ${no} (${nama}): at least one claim document is required because the insured is accepted.`,
    })
  })
  return out
}

/** Petunjuk Save to RNM nonaktif selama ada isian yang belum disimpan. */
export const PESAN_SIMPAN_DULU = 'Save your changes first.'

/** Bagian layar yang punya perubahan belum disimpan — untuk petunjuk di tab Claim Details (11-10-2026). */
export function bagianBerubah(klaim: Klaim, isian: IsianOutstanding): string[] {
  const out: string[] = []
  if (
    isian.stsKatastrofe !== (klaim.stsKatastrofe ?? '') ||
    isian.isKpr !== (klaim.isKpr ?? '') ||
    isian.katastrofeNote !== (klaim.katastrofeNote ?? '')
  ) {
    out.push(OS_PEGA.tabRegister)
  }
  if (adaPerubahan({ ...klaim, stsKatastrofe: isian.stsKatastrofe, isKpr: isian.isKpr, katastrofeNote: isian.katastrofeNote }, isian)) {
    out.push(OS_PEGA.pesertaDetail)
  }
  return out
}

/** Kalimat petunjuk: "You have unsaved changes (Register Claim). Click Save first."; tanpa bagian = PESAN_SIMPAN_DULU. */
export function pesanSimpanDulu(bagian: readonly string[] | undefined): string {
  return bagian !== undefined && bagian.length > 0
    ? `You have unsaved changes (${bagian.join(', ')}). Click Save first.`
    : PESAN_SIMPAN_DULU
}

function samaIsi<T extends object>(a: T, b: T): boolean {
  return (Object.keys(b) as (keyof T)[]).every((k) => a[k] === b[k])
}

/**
 * Ada isian yang BELUM disimpan tombol Save grid (11-10-2026, pilihan (a)): Save to RNM membaca data TERSIMPAN, jadi
 * selama ada perubahan ia nonaktif ("Save your changes first.") — tidak berjalan diam-diam di atas data lama. Pilihan (b)
 * "simpan dulu lalu Save to RNM" tidak dipakai: penomoran Save to RNM bertransaksi tersendiri, dua langkah itu tidak
 * dapat dibuat satu.
 *
 * Dibanding nilai MENTAH tersimpan, bukan isianAwal: bawaan Non-Catastrophe / Non-KPR pada klaim lama yang kolomnya
 * kosong belum tersimpan, dan Is KPR wajib di server.
 */
export function adaPerubahan(klaim: Klaim, isian: IsianOutstanding): boolean {
  if (
    isian.stsKatastrofe !== (klaim.stsKatastrofe ?? '') ||
    isian.isKpr !== (klaim.isKpr ?? '') ||
    isian.katastrofeNote !== (klaim.katastrofeNote ?? '')
  ) {
    return true
  }
  for (const p of klaim.peserta) {
    if ((isian.accept[p.id] ?? 'false') !== (p.isAccept === 'true' ? 'true' : 'false')) return true
    const c = isian.catatanAdmin?.[p.id]
    if (c !== undefined && c !== (p.catatanAdmin ?? '')) return true
    const t = isian.tanggal?.[p.id]
    if (t !== undefined && !samaIsi(t, tanggalAwal(p))) return true
    for (const b of p.baris) {
      const a = isian.adjustment?.[b.id]
      if (a !== undefined && !samaIsi(a, isianBarisAwal(b))) return true
    }
  }
  return false
}

function kosongTeks(v: string | undefined): boolean {
  return (v ?? '').trim() === ''
}

/**
 * Pre-check Save to RNM di layar (11-10-2026): kolom wajib yang KOSONG, DOL di luar periode valuation, dan Claim RNM >
 * Share, untuk SEMUA peserta — cermin `services.PeriksaSimpanRNM` (nama kolom dan kalimat sama). DOB sumber dan klaim
 * ganda hanya diperiksa server; server tetap lapis terakhir.
 */
export function kekuranganWajib(klaim: Klaim, isian: IsianOutstanding): KekuranganRNM[] {
  const out: KekuranganRNM[] = []
  if (kosongTeks(isian.isKpr)) {
    out.push({ nomorPeserta: 0, pesertaId: '', sertifikat: '', nama: '', kolom: 'Is KPR', pesan: 'Is KPR is required.' })
  }
  // Treaty (TP/TR) tidak dituntut berdokumen - langkah 3 XML (services.tipeTanpaGerbangDokumen).
  const tipe = (klaim.type ?? '').trim().toUpperCase()
  const dokumenWajib = tipe !== '' && tipe !== 'TP' && tipe !== 'TR'
  klaim.peserta.forEach((p, i) => {
    const no = i + 1
    const tambah = (kolom: string) =>
      out.push({
        nomorPeserta: no,
        pesertaId: p.id,
        sertifikat: p.nomorSertifikat,
        nama: p.namaTertanggung ?? '',
        kolom,
        pesan: `${kolom} is required for insured number ${no}`,
      })
    const t = isian.tanggal?.[p.id] ?? tanggalAwal(p)
    if (dokumenWajib && isian.accept[p.id] === 'true' && p.dokumen.length === 0) tambah('Document')
    const dolLuar = kosongTeks(t.kejadian) ? null : galatDOL(klaim, p, t.kejadian)
    if (kosongTeks(t.kejadian)) tambah('Date of Loss')
    else if (dolLuar !== null) {
      // Kalimat SAMA dengan server (GalatDOL.Kalimat) - aturan bisnis yang kini ditegakkan Save to RNM, bukan Save grid.
      out.push({ nomorPeserta: no, pesertaId: p.id, sertifikat: p.nomorSertifikat, nama: p.namaTertanggung ?? '',
        kolom: 'Date of Loss', pesan: `Insured number ${no}: ${dolLuar}` })
    }
    if (kosongTeks(p.tanggalMulai)) tambah('Begin Date')
    if (kosongTeks(p.tanggalExpired)) tambah('Expired Date')
    if (kosongTeks(p.nomorPolis)) tambah('Policy No')
    if (kosongTeks(p.nomorSertifikat)) tambah('Certificate No')
    if (kosongTeks(t.terimaKlaim)) tambah('Claim Received Date')
    if (kosongTeks(t.konfirmasi)) tambah('Confirmation Date')
    const jumlahBaris = (b: (typeof p.baris)[number]) => (isian.adjustment?.[b.id] ?? isianBarisAwal(b)).jumlah
    if (p.baris.some((b) => claimRNMKosong(jumlahBaris(b)))) tambah(KOLOM_CLAIM_RNM)
    else {
      // Claim RNM > Share Nusantara Re (baris belum berstatus) - kalimat SAMA dengan server.
      const lebih = p.baris
        .filter((b) => (b.kodeStatus ?? '').trim() === '')
        .map((b) => galatClaimRNM(jumlahBaris(b), b.shareNusantaraRe))
        .find((g) => g !== null && g !== PESAN_CLAIM_RNM_BUKAN_ANGKA)
      if (lebih !== undefined && lebih !== null) {
        out.push({ nomorPeserta: no, pesertaId: p.id, sertifikat: p.nomorSertifikat, nama: p.namaTertanggung ?? '',
          kolom: KOLOM_CLAIM_RNM, pesan: `Insured number ${no}: ${lebih}` })
      }
    }
  })
  return out
}

/** Kolom yang kurang untuk satu peserta — penanda merah di rinciannya. */
export function kolomKurangPeserta(kekurangan: readonly KekuranganRNM[] | null, pesertaID: string): Set<string> {
  return new Set((kekurangan ?? []).filter((k) => k.pesertaId === pesertaID).map((k) => k.kolom))
}

/**
 * SELECT ALL (`SelectAllClaimLife_act`): keadaan berikutnya menurut `jungkitPilihSemua` (b247), lalu SETIAP peserta
 * di-Accept dengan nilai itu (b429 `.IsAccept = Select.CARI1`).
 */
export function pilihSemuaAccept(
  keadaan: KeadaanPilihSemua,
  accept: Record<string, string>,
): { keadaan: KeadaanPilihSemua; accept: Record<string, string> } {
  const berikut = jungkitPilihSemua(keadaan)
  const baru: Record<string, string> = {}
  for (const id of Object.keys(accept)) baru[id] = berikut
  return { keadaan: berikut, accept: baru }
}

// ---------------------------------------------------------------------------
// Tab Register Claim
// ---------------------------------------------------------------------------

export function TabRegisterClaim({
  klaim,
  isian,
  onUbah,
  boleh,
}: {
  klaim: Klaim
  isian: IsianOutstanding
  onUbah: (i: IsianOutstanding) => void
  /** Kasus di tahap Outstanding dan terbuka. */
  boleh: boolean
}) {
  const terkunci = bernomor(klaim)
  const kprWarisan = isian.isKpr !== '' && !(PILIHAN_KPR as readonly string[]).includes(isian.isKpr)
  // Tata letak SAMA dengan tab Policy Details (PanelPolisOutstanding): os__form + dl.os__kolom + baris dt/dd
  // (keputusan work owner 11-10-2026, "nilai tidak boleh meluber ke kolom sebelah").
  return (
    <section className="os__form os__form--dua" aria-label={OS_PEGA.tabRegister}>
      <dl className="os__kolom">
        <div className="os__baris-form">
          <dt>{OS_PEGA.plNumber}</dt>
          <dd className="os__nilai-kode">{klaim.peserta[0]?.nomorPremiList || '—'}</dd>
        </div>
        {/* ⛔ radiogroup (bukan fieldset/legend): legend bukan butir grid yang andal. */}
        <div className="os__baris-form" role="radiogroup" aria-label={OS_PEGA.katastrofe}>
          <dt>{OS_PEGA.katastrofe}</dt>
          <dd className="os__radio-datar">
            {PILIHAN_KATASTROFE.map((v) => (
              <label key={v} className="os__radio">
                <input
                  type="radio"
                  name="stsKatastrofe"
                  value={v}
                  checked={isian.stsKatastrofe === v}
                  disabled={!boleh || terkunci}
                  onChange={() => onUbah({ ...isian, stsKatastrofe: v })}
                />
                {v}
              </label>
            ))}
          </dd>
        </div>
        <div className="os__baris-form" role="radiogroup" aria-label={OS_PEGA.isKpr}>
          <dt>
            {OS_PEGA.isKpr}
            <span className="os__wajib" aria-hidden="true">
              *
            </span>
          </dt>
          <dd className="os__radio-datar">
            {PILIHAN_KPR.map((v) => (
              <label key={v} className="os__radio">
                <input
                  type="radio"
                  name="isKpr"
                  value={v}
                  checked={isian.isKpr === v}
                  disabled={!boleh}
                  required
                  onChange={() => onUbah({ ...isian, isKpr: v })}
                />
                {v}
              </label>
            ))}
            {/* Data warisan di luar KPR / Non-KPR (mis. Yes/Y/1) ditampilkan apa adanya, tidak ditimpa diam-diam. */}
            {kprWarisan && <small className="os__catatan">Stored value: {isian.isKpr}</small>}
          </dd>
        </div>
      </dl>
      <dl className="os__kolom">
        <div className="os__baris-form">
          <dt>{OS_PEGA.stnc}</dt>
          <dd className="os__nilai-kode">{klaim.stncClaim || '—'}</dd>
        </div>
        {/* Catastrophe Note di BAWAH STNC Treaty (keputusan work owner 11-10-2026). */}
        {isian.stsKatastrofe === 'Catastrophe' && (
          <div className="os__baris-form os__catatan-katastrofe">
            <dt>{OS_PEGA.catatanKatastrofe}</dt>
            <dd>
              <textarea
                className="field__input os__isian-area"
                rows={3}
                maxLength={1000}
                value={isian.katastrofeNote}
                readOnly={!boleh || terkunci}
                placeholder="Describe the catastrophe…"
                aria-label={OS_PEGA.catatanKatastrofe}
                onChange={(e) => onUbah({ ...isian, katastrofeNote: e.target.value })}
              />
            </dd>
          </div>
        )}
      </dl>
    </section>
  )
}

// ---------------------------------------------------------------------------
// Grid Participant Details
// ---------------------------------------------------------------------------

export function GridPeserta({
  klaim,
  accept,
  onAccept,
  onPilihSemua,
  terbuka,
  onBuka,
  onHapus,
  boleh,
  rincian,
  kakiKanan,
}: {
  klaim: Klaim
  accept: Record<string, string>
  onAccept: (pesertaID: string, nilai: string) => void
  onPilihSemua: () => void
  terbuka: string | null
  onBuka: (pesertaID: string | null) => void
  onHapus: (pesertaID: string) => void
  boleh: boolean
  /** Isi baris yang dibuka (ClaimLifeDetailGCNM) - dipisah supaya grid dapat diuji tanpa server. */
  rincian: (pesertaID: string) => ReactNode
  /** Tombol Save kecil di kanan bawah grid (b19413). */
  kakiKanan?: ReactNode
}) {
  const [kata, setKata] = useState('')
  const tampilHapus = bolehDeletePeserta(klaim)
  const lebar = tampilHapus ? 6 : 5
  // Nomor urut 1..n mengikuti urutan peserta klaim, tidak berubah oleh saringan.
  const nomor = new Map(klaim.peserta.map((p, i) => [p.id, i + 1]))
  return (
    <>
      <label className="os__saring">
        <span className="os__label">{OS_PEGA.kolomNama}</span>
        <input
          type="search"
          className="field__input os__isian-teks"
          value={kata}
          placeholder="Filter"
          onChange={(e) => setKata(e.target.value)}
        />
      </label>
      <div className="inbox__gulir">
        <table className="inbox__tabel os__grid">
          <thead>
            <tr>
              <th>No.</th>
              <th>{OS_PEGA.kolomPL}</th>
              <th>{OS_PEGA.kolomNama}</th>
              <th>{OS_PEGA.kolomLahir}</th>
              <th>
                {/* SELECT ALL b16609 — tautan kepala kolom Accept (SelectAllClaimLife_act). */}
                <button type="button" className="os__tautan" disabled={!boleh} onClick={onPilihSemua}>
                  {OS_PEGA.pilihSemua}
                </button>
              </th>
              {tampilHapus && <th>{OS_PEGA.hapus}</th>}
            </tr>
          </thead>
          <tbody>
            {saringPeserta(klaim.peserta, kata).map((p) => {
              const buka = terbuka === p.id
              return [
                // Klik baris membuka rincian DI BAWAH baris itu (master-detail b18252) dan menutup baris lain.
                <tr
                  key={p.id}
                  className={'inbox__baris' + (buka ? ' os__baris--pilih' : '')}
                  onClick={() => onBuka(buka ? null : p.id)}
                >
                  <td>
                    <button
                      type="button"
                      className="os__buka"
                      aria-expanded={buka}
                      onClick={(e) => {
                        e.stopPropagation()
                        onBuka(buka ? null : p.id)
                      }}
                    >
                      {nomor.get(p.id)}
                    </button>
                  </td>
                  <td>{p.nomorPremiList || '—'}</td>
                  <td>{p.namaTertanggung || '—'}</td>
                  <td>{tanggalTampil(p.tanggalLahir)}</td>
                  <td onClick={(e) => e.stopPropagation()}>
                    <label className="os__radio">
                      <input
                        type="checkbox"
                        checked={accept[p.id] === 'true'}
                        disabled={!boleh}
                        onChange={(e) => onAccept(p.id, e.target.checked ? 'true' : 'false')}
                      />{' '}
                      {OS_PEGA.accept}
                    </label>
                  </td>
                  {tampilHapus && (
                    <td onClick={(e) => e.stopPropagation()}>
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => onHapus(p.id)}>
                        {OS_PEGA.hapus}
                      </button>
                    </td>
                  )}
                </tr>,
                buka ? (
                  <tr key={`${p.id}:rincian`} className="os__rincian">
                    <td colSpan={lebar}>{rincian(p.id)}</td>
                  </tr>
                ) : null,
              ]
            })}
          </tbody>
        </table>
      </div>
      {kakiKanan !== undefined && <p className="os__simpan os__simpan--kanan">{kakiKanan}</p>}
    </>
  )
}

// ---------------------------------------------------------------------------
// Tombol bawah
// ---------------------------------------------------------------------------

/**
 * Tombol bawah layar Outstanding — HANYA Save to RNM dan Send to Medical Check (keputusan work owner 11-10-2026:
 * "Tombol send back to register, check readiness claim, dan close claim hilangkan saja"). Rute dan aksi ketiganya TIDAK
 * dihapus dari server; Close Claim tetap di layar Detail / Claim Analis (PanelTutupKlaim).
 */
export function TombolBawahOutstanding({
  klaim,
  sibukRNM,
  onSimpanRNM,
  sesudahPindah,
  adaPerubahan,
  bagian,
  sebelumPindah,
  tanganiGalat,
}: {
  klaim: Klaim
  sibukRNM: boolean
  onSimpanRNM: () => void
  sesudahPindah: () => void
  /** Ada isian belum disimpan — Save to RNM nonaktif dengan petunjuk (11-10-2026). */
  adaPerubahan?: boolean
  /** Bagian yang berubah (bagianBerubah) — disebut dalam petunjuk. */
  bagian?: readonly string[]
  /** Pra-periksa Send to Medical Check (false = tidak memindah). */
  sebelumPindah?: (t: TombolPindah) => boolean
  /** Galat perpindahan berdaftar ke tempat galat halaman; true = sudah ditangani. */
  tanganiGalat?: (e: unknown) => boolean
}) {
  const aturan = aturanTombolOutstanding(klaim)
  // Perubahan belum disimpan menonaktifkan Save to RNM DAN Send to Medical Check (11-10-2026: Accept harus tersimpan).
  const medisAktif = aturan.kirimMedisTampil && aturan.kirimMedisAktif
  const tungguSimpan = adaPerubahan === true && (aturan.simpanRNMAktif || medisAktif)
  const aturTombol = (t: TombolPindah) =>
    t.tujuan === TAHAP_JALUR.inputRegister
      ? // Send Back to Register (b21404) disembunyikan di layar ini (11-10-2026).
        { tampil: false, aktif: false }
      : t.tujuan === TAHAP_JALUR.medicalCheck
        ? { tampil: aturan.kirimMedisTampil, aktif: aturan.kirimMedisAktif && adaPerubahan !== true }
        : { tampil: true, aktif: true }
  return (
    <>
      {/* Petunjuk sebagai teks bantuan DI ATAS deretan tombol (11-10-2026), bukan menempel di tombol. */}
      {tungguSimpan && (
        <p role="note" className="alert alert--info os__bantuan">
          {pesanSimpanDulu(bagian)}
        </p>
      )}
      <div className="os__aksi">
        {/* Save to RNM b21102 — gerbangnya seluruhnya di server; nonaktif sesudah tersimpan (b21095, OQ-N1). */}
        <button
          type="button"
          className="btn btn--primary"
          disabled={sibukRNM || !aturan.simpanRNMAktif || tungguSimpan}
          title={tungguSimpan ? pesanSimpanDulu(bagian) : undefined}
          onClick={onSimpanRNM}
        >
          {sibukRNM ? 'Saving…' : TOMBOL_OS.simpanRNM}
        </button>
        {/* ⛔ Perpindahan TIDAK disalin di sini (GILIRAN-11): satu daftar, satu dialog — PanelPindahTahap. */}
        <PanelPindahTahap
          klaimID={klaim.id}
          tahap={TAHAP.outstandingClaim}
          sesudahPindah={sesudahPindah}
          aturTombol={aturTombol}
          tanpaJudul
          sebelumPindah={sebelumPindah}
          tanganiGalat={tanganiGalat}
        />
      </div>
    </>
  )
}

/** Kalimat modal konfirmasi Save to RNM (keputusan work owner 11-10-2026). */
export const KONFIRMASI_RNM = {
  judul: 'Save to RNM',
  tanya: 'Are you sure you want to Save to RNM?',
  peringatan:
    'This action cannot be undone. A claim number will be generated, the insured will be set to Outstanding, and the claim data will be locked.',
  batal: 'Cancel',
  ya: 'Yes, Save to RNM',
  proses: 'Saving…',
} as const

/**
 * Modal konfirmasi Save to RNM — Modal inti, pola KonfirmasiDaftar (keputusan work owner 11-10-2026: "buat pop up pesan
 * validasi apakah yakin ingin save to rnm dan info tidak bisa dikembalikan lagi"). Dibuka HANYA sesudah pre-check lolos.
 * Cancel / X / Escape menutup tanpa memanggil server; selama proses tombolnya nonaktif ("Saving…") dan modal tidak dapat
 * ditutup — satu klik, satu permintaan.
 */
export function KonfirmasiSimpanRNM({
  klaim,
  sibuk,
  onLanjut,
  onTutup,
}: {
  klaim: Klaim
  sibuk: boolean
  onLanjut: () => void
  onTutup: () => void
}) {
  const pl = [...new Set(klaim.peserta.map((p) => p.nomorPremiList).filter((x) => x !== ''))]
  return (
    <Modal
      judul={KONFIRMASI_RNM.judul}
      onTutup={onTutup}
      tanpaTutup={sibuk}
      labelBatal={KONFIRMASI_RNM.batal}
      aksi={
        <button type="button" className="btn btn--primary" onClick={onLanjut} disabled={sibuk}>
          {sibuk ? KONFIRMASI_RNM.proses : KONFIRMASI_RNM.ya}
        </button>
      }
    >
      <p>{KONFIRMASI_RNM.tanya}</p>
      <p role="note" className="alert alert--warn">
        {KONFIRMASI_RNM.peringatan}
      </p>
      <dl className="os__ringkas-rnm">
        <div className="os__baris-form">
          <dt>{OS_PEGA.pesertaDetail}</dt>
          <dd>{klaim.peserta.length}</dd>
        </div>
        <div className="os__baris-form">
          <dt>{OS_PEGA.plNumber}</dt>
          <dd className="os__nilai-kode">{pl.length > 0 ? pl.join(', ') : '—'}</dd>
        </div>
      </dl>
    </Modal>
  )
}

/** Kalimat modal konfirmasi Send to Medical Check (keputusan work owner 11-10-2026: "iya perlu pop up konfirmasi"). */
export const KONFIRMASI_MEDIS = {
  judul: 'Send to Medical Check',
  tanya: 'Are you sure you want to send this claim to Medical Check?',
  peringatan:
    'After sending, the claim moves to the Medical Advisor queue and can no longer be edited from the Outstanding screen.',
  batal: 'Cancel',
  ya: 'Yes, Send to Medical Check',
  proses: 'Sending…',
  claimNo: 'Claim No',
  peserta: 'Insured',
  accept: 'Accepted insured',
} as const

/** Pita sukses Inbox sesudah Send to Medical Check (pola Register → Inbox). */
export function kabarTerkirimMedis(klaimID: string): string {
  return `Claim ${klaimID} sent to Medical Check.`
}

/** Kabar Inbox sesudah perpindahan dari layar kasus, per tujuan (11-10-2026). */
export function kabarPindah(klaimID: string, tujuan: TahapJalur): string {
  switch (tujuan) {
    case TAHAP_JALUR.outstanding:
      return `Claim ${klaimID} sent back to Admin.`
    case TAHAP_JALUR.medicalCheck:
      return kabarTerkirimMedis(klaimID)
    case TAHAP_JALUR.claimAnalis:
      return `Claim ${klaimID} sent to Claim Analyst.`
    default:
      return `Claim ${klaimID} sent back to Register.`
  }
}

/** Kabar Inbox sesudah Close Claim dari layar kasus. */
export function kabarTutup(nomor: string): string {
  return `Claim ${nomor} closed.`
}

/**
 * Tombol bawah tahap Medical Check (`MedicalCheckClaimLife`: Save b18572 · Send Back to Admin b20256 · Send to Claim
 * Analyst b21151) dan Claim Analis (Send Back to Admin · Send Back to Medical Check · Close Claim). Perpindahan TIDAK
 * disalin: satu daftar, satu dialog — PanelPindahTahap. Gerbang pemegang tahap ditegakkan server (403/409).
 */
export function TombolBawahTahap({
  klaim,
  onSimpan,
  sibukSimpan = false,
  bolehSimpan = false,
  sesudahPindah,
  sesudahTutup,
  adaPerubahan = false,
}: {
  klaim: Klaim
  /** Recommendation belum disimpan — perpindahan nonaktif dengan petunjuk (pola Outstanding, 11-10-2026). */
  adaPerubahan?: boolean
  /** Save Medical Check (Recommendation); tanpa ini tombol Save tidak tampil. */
  onSimpan?: () => void
  sibukSimpan?: boolean
  /** Save aktif (ada perubahan dan pelaku Medical Advisor). */
  bolehSimpan?: boolean
  sesudahPindah: (tujuan: TahapJalur) => void
  sesudahTutup: (nomor: string) => void
}) {
  if (kasusTertutup(klaim)) return null
  return (
    <>
      {adaPerubahan && (
        <p role="note" className="alert alert--info os__bantuan">
          Save your changes first.
        </p>
      )}
      <div className="os__aksi">
        {klaim.tahap === TAHAP.medicalCheck && onSimpan !== undefined && (
          <button type="button" className="btn btn--primary" disabled={sibukSimpan || !bolehSimpan} onClick={onSimpan}>
            {sibukSimpan ? 'Saving…' : TOMBOL_MEDIS.simpan}
          </button>
        )}
        <PanelPindahTahap
          klaimID={klaim.id}
          tahap={klaim.tahap}
          sesudahPindah={sesudahPindah}
          aturTombol={() => ({ tampil: true, aktif: !adaPerubahan })}
          tanpaJudul
        />
        {klaim.tahap === TAHAP.claimAnalis && <PanelTutupKlaim klaim={klaim} sesudahTutup={sesudahTutup} />}
      </div>
    </>
  )
}

/**
 * Modal konfirmasi Send to Medical Check — Modal inti, pola KonfirmasiSimpanRNM. Dibuka HANYA sesudah pra-periksa dokumen
 * lolos; Cancel / X / Escape menutup tanpa memanggil server; selama proses nonaktif ("Sending…") dan tak dapat ditutup.
 */
export function KonfirmasiKirimMedis({
  klaim,
  sibuk,
  onLanjut,
  onTutup,
}: {
  klaim: Klaim
  sibuk: boolean
  onLanjut: () => void
  onTutup: () => void
}) {
  const diterima = klaim.peserta.filter((p) => p.isAccept === 'true').length
  return (
    <Modal
      judul={KONFIRMASI_MEDIS.judul}
      onTutup={onTutup}
      tanpaTutup={sibuk}
      labelBatal={KONFIRMASI_MEDIS.batal}
      aksi={
        <button type="button" className="btn btn--primary" onClick={onLanjut} disabled={sibuk}>
          {sibuk ? KONFIRMASI_MEDIS.proses : KONFIRMASI_MEDIS.ya}
        </button>
      }
    >
      <p>{KONFIRMASI_MEDIS.tanya}</p>
      <p role="note" className="alert alert--warn">
        {KONFIRMASI_MEDIS.peringatan}
      </p>
      <dl className="os__ringkas-rnm">
        <div className="os__baris-form">
          <dt>{KONFIRMASI_MEDIS.claimNo}</dt>
          <dd className="os__nilai-kode">{klaim.nomorKlaim || '—'}</dd>
        </div>
        <div className="os__baris-form">
          <dt>{KONFIRMASI_MEDIS.peserta}</dt>
          <dd>{klaim.peserta.length}</dd>
        </div>
        <div className="os__baris-form">
          <dt>{KONFIRMASI_MEDIS.accept}</dt>
          <dd>{diterima}</dd>
        </div>
      </dl>
    </Modal>
  )
}

// ---------------------------------------------------------------------------
// Halaman
// ---------------------------------------------------------------------------

/**
 * Admin Notes (dan Location of Death, bila kelak berkolom) dapat diubah HANYA oleh Admin selama kasus di inbox Admin —
 * tahap Outstanding, kasus terbuka (keputusan work owner 11-10-2026: "hanya bisa diedit oleh admin dan jika masih di inbox
 * admin. di inbox lain hanya read only"). Server menegakkan yang sama pada Save grid: 409 bukan tahap Outstanding, 403
 * bukan pemegang tahap (ReasLifeAdmin).
 */
export function bolehUbahCatatanAdmin(klaim: Klaim | null, peran: readonly KodePeran[] | undefined): boolean {
  return bolehIsianOutstanding(klaim) && (peran ?? []).includes(PERAN.admin)
}

/**
 * Kotak MEDICAL ADVISOR (diagnosa, Recommendation) dapat diubah HANYA di tahap Medical Check oleh Medical Advisor, kasus
 * terbuka (11-10-2026; b4908/b6401 pyPosition ReasLifeMedicalAdvisor). Server menegakkan yang sama (409 tahap, 403
 * pemegang tahap); gerbang STS_REJECT per peserta di RincianPeserta.
 */
/** Recommendation yang benar-benar berbeda dari nilai tersimpan (peserta aktif saja). */
export function rekomendasiBerubah(klaim: Klaim, isian: Record<string, string>): Record<string, string> {
  const out: Record<string, string> = {}
  for (const p of klaim.peserta) {
    const v = isian[p.id]
    if (v !== undefined && v !== (p.rekomendasi ?? '')) out[p.id] = v
  }
  return out
}

export function bolehUbahMedis(klaim: Klaim | null, peran: readonly KodePeran[] | undefined): boolean {
  return klaim !== null && !kasusTertutup(klaim) && klaim.tahap === TAHAP.medicalCheck && (peran ?? []).includes(PERAN.medis)
}

export default function OutstandingClaimLife({ klaimID, onPindah, peran }: OutstandingProps) {
  const [klaim, setKlaim] = useState<Klaim | null>(null)
  // ⛔ Klaim yang GAGAL dimuat dinyatakan, bukan ditelan (temuan /code-review).
  const [galatMuat, setGalatMuat] = useState<unknown>(null)
  const [polis, setPolis] = useState<PolicyDataLife | null>(null)
  const [polisTidakDitemukan, setPolisTidakDitemukan] = useState<{ idPega: string } | undefined>(undefined)
  const [tab, setTab] = useState<TabOS>('register')
  const [isian, setIsian] = useState<IsianOutstanding | null>(null)
  const [keadaanSemua, setKeadaanSemua] = useState<KeadaanPilihSemua>('')
  const [terbuka, setTerbuka] = useState<string | null>(null)
  const [hapus, setHapus] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)
  const [sibukRNM, setSibukRNM] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [tersimpan, setTersimpan] = useState<HasilSimpanRNM | null>(null)
  // Kekurangan kolom wajib Save to RNM (pre-check layar atau 422 server) - satu kotak, bukan galat biasa.
  const [kekurangan, setKekurangan] = useState<KekuranganRNM[] | null>(null)
  // Judul kotak kekurangan: Save to RNM (bawaan) atau Send to Medical Check.
  const [judulKekurangan, setJudulKekurangan] = useState<string | undefined>(undefined)
  // Modal konfirmasi Save to RNM (sesudah pre-check lolos).
  const [konfirmasiRNM, setKonfirmasiRNM] = useState(false)
  // Modal konfirmasi Send to Medical Check (sesudah pra-periksa dokumen lolos, 11-10-2026).
  const [konfirmasiMedis, setKonfirmasiMedis] = useState(false)
  const [sibukMedis, setSibukMedis] = useState(false)
  // Recommendation yang diubah di layar (Medical Check, 11-10-2026) — hanya peserta yang berubah; dibuang sesudah Save.
  const [rekomendasi, setRekomendasi] = useState<Record<string, string>>({})
  const [sibukSimpanMedis, setSibukSimpanMedis] = useState(false)

  /**
   * Isian berubah: kotak kekurangan lama DIBERSIHKAN (11-10-2026) — kotak itu hanya milik klik Save to RNM atas data
   * TERSIMPAN; selama ada perubahan Save to RNM nonaktif dan yang tampil hanya petunjuk "Save your changes first.".
   */
  function ubahIsian(i: IsianOutstanding): void {
    setIsian(i)
    setKekurangan(null)
  }

  const muatUlang = useCallback(async (): Promise<void> => {
    const k = await ambilKlaimLife(klaimID)
    setKlaim(k)
    setIsian(isianAwal(k))
  }, [klaimID])

  useEffect(() => {
    let batal = false
    void (async () => {
      setGalatMuat(null)
      let idPega: string
      try {
        const k = await ambilKlaimLife(klaimID)
        if (batal) return
        setKlaim(k)
        setIsian(isianAwal(k))
        // ⛔ 10-10-2026: header polis dari IDPEGA peserta klaim (server menurunkannya lewat SOURCE_ID).
        idPega = k.peserta?.[0]?.idPega ?? ''
      } catch (e) {
        if (!batal) setGalatMuat(e)
        return
      }
      // Header yang tidak ditemukan BUKAN kerusakan layar ini: panelnya menyatakannya, Save to RNM menolaknya terang.
      const h = await bacaPolisPeserta(ambilDataPolisIdPega, idPega)
      if (!batal) {
        setPolis(h.polis)
        setPolisTidakDitemukan(h.tidakDitemukan ? { idPega: h.idPega } : undefined)
      }
    })()
    return () => {
      batal = true
    }
  }, [klaimID])

  /** Tombol Save b19413 — Accept peserta dan medan tab Register Claim. */
  async function simpan(): Promise<void> {
    if (isian === null || sibuk) return
    setSibuk(true)
    // Setiap aksi mulai bersih: galat, kotak kekurangan, dan kabar lama dibuang.
    setGalat(null)
    setKekurangan(null)
    setPesan(null)
    setTersimpan(null)
    try {
      await simpanIsianOutstanding(klaimID, isian)
      await muatUlang()
      // Save = simpan progres saja (11-10-2026).
      setPesan('Progress saved.')
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  /** Tombol Save layar Medical Check — Recommendation yang berubah; galat di tempat galat yang sama. */
  async function simpanRekomendasi(): Promise<void> {
    if (sibukSimpanMedis || klaim === null) return
    const berubah = rekomendasiBerubah(klaim, rekomendasi)
    if (Object.keys(berubah).length === 0) return
    setSibukSimpanMedis(true)
    setGalat(null)
    setKekurangan(null)
    setPesan(null)
    try {
      await simpanMedis(klaimID, { rekomendasi: berubah })
      await muatUlang()
      setRekomendasi({})
      setPesan('Medical Check saved.')
    } catch (e) {
      setGalat(e)
    } finally {
      setSibukSimpanMedis(false)
    }
  }

  /** Menampilkan kekurangan dan membuka peserta bermasalah PERTAMA (tab Claim Details), atau tab Register Claim. */
  function tunjukkanKekurangan(k: KekuranganRNM[], judul?: string): void {
    setJudulKekurangan(judul)
    setKekurangan(k)
    const pertama = k.find((x) => x.pesertaId !== '')
    if (pertama !== undefined) {
      setTab('klaim')
      setTerbuka(pertama.pesertaId)
    } else {
      setTab('register')
    }
  }

  /**
   * `Save to RNM` b21102. Pre-check kolom wajib di layar LEBIH DULU (tanpa permintaan bila ada yang kosong); server tetap
   * memeriksa seluruh gerbang dan mengembalikan SEMUA kekurangan sekaligus (11-10-2026).
   */
  async function simpanRNM(): Promise<void> {
    // Selama ada perubahan belum disimpan Save to RNM nonaktif — dan pre-check TIDAK berjalan (11-10-2026).
    if (sibukRNM || klaim === null || isian === null || adaPerubahan(klaim, isian)) return
    setGalat(null)
    setKekurangan(null)
    setPesan(null)
    setTersimpan(null)
    // Pre-check atas data TERSIMPAN (isian = nilai tersimpan, sebab tidak ada perubahan).
    const kurang = kekuranganWajib(klaim, isian)
    if (kurang.length > 0) {
      tunjukkanKekurangan(kurang)
      return
    }
    // Lolos: konfirmasi dulu (11-10-2026). Server baru dipanggil dari tombol "Yes, Save to RNM".
    setKonfirmasiRNM(true)
  }

  /** "Yes, Save to RNM" — SATU permintaan; galat server tampil di tempat galat biasa sesudah modal ditutup. */
  async function jalankanSimpanRNM(): Promise<void> {
    if (sibukRNM) return
    setSibukRNM(true)
    try {
      setTersimpan(await simpanKeRNM(klaimID))
      await muatUlang()
    } catch (e) {
      setTersimpan(null)
      if (e instanceof GagalSimpanRNM) tunjukkanKekurangan(e.kekurangan)
      else setGalat(e)
    } finally {
      setSibukRNM(false)
      setKonfirmasiRNM(false)
    }
  }

  /**
   * "Yes, Send to Medical Check" — SATU-SATUNYA pemanggil pindahTahap di layar ini. Berhasil: kembali ke Inbox dengan pita
   * sukses (layar ini milik tahap Outstanding). Gagal: galat di tempat galat biasa sesudah modal ditutup (422 berdaftar ke
   * kotak kekurangan).
   */
  async function jalankanKirimMedis(): Promise<void> {
    if (sibukMedis) return
    setSibukMedis(true)
    try {
      await pindahTahap(klaimID, TAHAP_JALUR.medicalCheck)
      setKonfirmasiMedis(false)
      onPindah(kabarTerkirimMedis(klaimID))
    } catch (e) {
      setKonfirmasiMedis(false)
      if (e instanceof GagalSimpanRNM) tunjukkanKekurangan(e.kekurangan, JUDUL_KEKURANGAN_MEDIS)
      else setGalat(e)
    } finally {
      setSibukMedis(false)
    }
  }

  /** DELETE b17865 → cabut peserta (penanda STS_HAPUS, ADR-U-0031) sesudah konfirmasi. */
  async function jalankanHapus(pesertaID: string): Promise<void> {
    setSibuk(true)
    setGalat(null)
    setKekurangan(null)
    setPesan(null)
    try {
      await cabutPeserta(klaimID, pesertaID)
      setHapus(null)
      if (terbuka === pesertaID) setTerbuka(null)
      await muatUlang()
    } catch (e) {
      setHapus(null)
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  const boleh = bolehIsianOutstanding(klaim)
  const hapusPeserta = klaim?.peserta.find((p) => p.id === hapus)

  return (
    <section className="os">
      <header className="os__kepala">
        <div>
          <h2 className="os__judul">
            {OS_PEGA.judulKasus} ({klaimID})
          </h2>
          <p className="os__nomor">
            <span className="os__label">Claim No</span>
            <strong>{klaim !== null && bernomor(klaim) ? klaim.nomorKlaim : '—'}</strong>
          </p>
        </div>
      </header>

      {/* SATU tempat galat & kabar layar ini (11-10-2026, "rapi seperti premium list"): di ATAS isi, di bawah kepala —
          pola PremiumListDetail / InputOffer (Gagal di atas panel). Galat muat, galat Save grid / Save to RNM / Delete,
          dan kotak kekurangan memakai komponen yang sama (PitaGalat), dapat ditutup, dibersihkan setiap aksi. */}
      <div className="cl-tempat-galat">
        <PitaGalat galat={galatMuat} />
        <PitaGalat
          galat={galat}
          kekurangan={kekurangan}
          judul={judulKekurangan}
          onTutup={() => {
            setGalat(null)
            setKekurangan(null)
          }}
        />
        {pesan !== null && (
          <p role="status" className="alert alert--ok">
            {pesan}
          </p>
        )}
        {tersimpan !== null && (
          <p role="status" className="alert alert--ok">
            Saved to RNM — claim number {tersimpan.nomorKlaim}
            {tersimpan.nomorBaru ? ' (newly issued)' : ''}; {tersimpan.barisDitandai} row(s) marked Outstanding.
          </p>
        )}
      </div>

      {klaim !== null && isian !== null && (
        <>
          <StripTab tab={TAB_OS} aktif={tab} onPilih={setTab} label={labelTabOS} />

          {tab === 'register' && <TabRegisterClaim klaim={klaim} isian={isian} onUbah={ubahIsian} boleh={boleh} />}

          {/* Panel yang SAMA dengan layar Register — polis yang sama, tahap berbeda. */}
          {tab === 'polis' && (
            <>
              {/* Tiga kolom seperti Pega (11-10-2026); PanelDataPolis Register tetap tanpa System Reinsurance / Premium
                  Method. */}
              <PanelPolisOutstanding
                polis={polis}
                plNumber={klaim.peserta[0]?.nomorPremiList ?? ''}
                polisTidakDitemukan={polisTidakDitemukan}
              />
              {/* ⛔ View Detail (DetailPolisLife b13267) DICABUT dari tab ini - keputusan work owner 11-10-2026. */}
            </>
          )}

          {tab === 'klaim' && (
            <section className="register__kartu os__kartu" aria-label={OS_PEGA.pesertaDetail}>
              <h3>{OS_PEGA.pesertaDetail}</h3>
              <GridPeserta
                klaim={klaim}
                accept={isian.accept}
                onAccept={(id, v) => {
                  ubahIsian({ ...isian, accept: { ...isian.accept, [id]: v } })
                }}
                onPilihSemua={() => {
                  const h = pilihSemuaAccept(keadaanSemua, isian.accept)
                  setKeadaanSemua(h.keadaan)
                  ubahIsian({ ...isian, accept: h.accept })
                }}
                terbuka={terbuka}
                onBuka={setTerbuka}
                onHapus={setHapus}
                boleh={boleh}
                rincian={(id) => {
                  const p = klaim.peserta.find((x) => x.id === id)
                  return p === undefined ? null : (
                    <RincianPeserta
                      klaim={klaim}
                      p={p}
                      onMuatUlang={muatUlang}
                      catatanAdmin={isian.catatanAdmin?.[p.id] ?? p.catatanAdmin ?? ''}
                      onUbahCatatan={
                        bolehUbahCatatanAdmin(klaim, peran)
                          ? (v) => ubahIsian({ ...isian, catatanAdmin: { ...isian.catatanAdmin, [p.id]: v } })
                          : undefined
                      }
                      // Lima tanggal diisi langsung dan disimpan tombol Save grid (11-10-2026); sesudah Save to RNM baca saja.
                      tanggal={isian.tanggal?.[p.id] ?? tanggalAwal(p)}
                      onUbahTanggal={(t) => ubahIsian({ ...isian, tanggal: { ...isian.tanggal, [p.id]: t } })}
                      // Claim RNM + bank dilebur ke tombol Save grid (11-10-2026).
                      isianAdjustment={isian.adjustment}
                      kolomKurang={kolomKurangPeserta(kekurangan, p.id)}
                      bolehMedis={bolehUbahMedis(klaim, peran)}
                      rekomendasi={rekomendasi[p.id]}
                      onUbahRekomendasi={
                        bolehUbahMedis(klaim, peran)
                          ? (v) => setRekomendasi((r) => ({ ...r, [p.id]: v }))
                          : undefined
                      }
                      onUbahAdjustment={
                        boleh
                          ? (adjID, isi) => ubahIsian({ ...isian, adjustment: { ...isian.adjustment, [adjID]: isi } })
                          : undefined
                      }
                    />
                  )
                }}
                kakiKanan={
                  boleh ? (
                    <button
                      type="button"
                      className="btn btn--primary btn--sm"
                      disabled={sibuk || adaIsianTidakSah(klaim, isian)}
                      onClick={() => void simpan()}
                    >
                      {sibuk ? 'Saving…' : OS_PEGA.simpan}
                    </button>
                  ) : undefined
                }
              />
              {/* Save to RNM + Send to Medical Check HANYA di tab Claim Details, di kaki tab sesudah Save grid (keputusan work
                  owner 11-10-2026). Perubahan tab Register tetap terdeteksi dan disebut di petunjuk. */}
              {klaim.tahap === TAHAP.outstandingClaim && (
                <div className="os__kaki-tab">
                  <TombolBawahOutstanding
                    klaim={klaim}
                    sibukRNM={sibukRNM}
                    adaPerubahan={adaPerubahan(klaim, isian)}
                    bagian={bagianBerubah(klaim, isian)}
                    onSimpanRNM={() => {
                      void simpanRNM()
                    }}
                    sebelumPindah={(t) => {
                      if (t.tujuan !== TAHAP_JALUR.medicalCheck) return true
                      // Pra-periksa Send to Medical Check (11-10-2026): peserta Accept tanpa dokumen → kotak kekurangan,
                      // peserta pertama dibuka di sub-tab Document; tanpa permintaan ke server.
                      setGalat(null)
                      setPesan(null)
                      const kurang = kekuranganDokumenMedis(klaim)
                      if (kurang.length > 0) {
                        tunjukkanKekurangan(kurang, JUDUL_KEKURANGAN_MEDIS)
                        return false
                      }
                      setKekurangan(null)
                      // Lolos: konfirmasi dulu (11-10-2026); perpindahan dijalankan tombol modal, bukan panel.
                      setKonfirmasiMedis(true)
                      return false
                    }}
                    tanganiGalat={(e) => {
                      if (e instanceof GagalSimpanRNM) {
                        tunjukkanKekurangan(e.kekurangan, JUDUL_KEKURANGAN_MEDIS)
                        return true
                      }
                      return false
                    }}
                    sesudahPindah={() => onPindah()}
                  />
                </div>
              )}
              {/* Medical Check / Claim Analis (11-10-2026): tombol tahap itu sendiri, di tempat yang sama. */}
              {(klaim.tahap === TAHAP.medicalCheck || klaim.tahap === TAHAP.claimAnalis) && (
                <div className="os__kaki-tab">
                  <TombolBawahTahap
                    klaim={klaim}
                    onSimpan={() => {
                      void simpanRekomendasi()
                    }}
                    sibukSimpan={sibukSimpanMedis}
                    bolehSimpan={bolehUbahMedis(klaim, peran) && Object.keys(rekomendasiBerubah(klaim, rekomendasi)).length > 0}
                    adaPerubahan={Object.keys(rekomendasiBerubah(klaim, rekomendasi)).length > 0}
                    sesudahPindah={(tujuan) => onPindah(kabarPindah(klaimID, tujuan))}
                    sesudahTutup={(nomor) => onPindah(kabarTutup(nomor))}
                  />
                </div>
              )}
            </section>
          )}

        </>
      )}


      {konfirmasiMedis && klaim !== null && (
        <KonfirmasiKirimMedis
          klaim={klaim}
          sibuk={sibukMedis}
          onLanjut={() => {
            void jalankanKirimMedis()
          }}
          onTutup={() => {
            if (!sibukMedis) setKonfirmasiMedis(false)
          }}
        />
      )}

      {konfirmasiRNM && klaim !== null && (
        <KonfirmasiSimpanRNM
          klaim={klaim}
          sibuk={sibukRNM}
          onLanjut={() => {
            void jalankanSimpanRNM()
          }}
          onTutup={() => {
            if (!sibukRNM) setKonfirmasiRNM(false)
          }}
        />
      )}

      {hapus !== null && (
        <Modal
          judul={OS_PEGA.konfirmasiHapus}
          onTutup={() => setHapus(null)}
          aksi={
            <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void jalankanHapus(hapus)}>
              {OS_PEGA.hapus}
            </button>
          }
        >
          <p>
            {hapusPeserta?.namaTertanggung || hapusPeserta?.nomorSertifikat || hapus} — {hapusPeserta?.nomorPremiList}
          </p>
        </Modal>
      )}
    </section>
  )
}
