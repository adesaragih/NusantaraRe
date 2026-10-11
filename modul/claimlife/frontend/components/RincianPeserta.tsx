// Rincian satu peserta klaim — panel master-detail `ViewClaimDetailLifeGCNM` → `ClaimLifeDetailGCNM`.
//
// ⛔ DIPINDAH, bukan ditulis ulang (10-10-2026): blok per peserta yang dulu hidup di `pages/KlaimLife.tsx` kini di sini,
// dipakai BERSAMA layar Detail (KlaimLife) dan baris terbuka grid Participant Details layar Outstanding.
//
// ⭐ TATA LETAK 11-10-2026 mengikuti screenshot Pega CLMLF-9293 (keputusan work owner "saat pesertanya diklik bisa
// menampilkan detail pesertanya dan bisa isi adjustment-nya seperti di pega"); label VERBATIM `ClaimLifeDetailGCNM.xml`
// (nomor baris di labels.ts OS_PEGA.rincian):
//   kiri    CERTIFICATE NO, POLICY NO, PLAN, AGE, SEX; kotak LOCATION OF DEATH (TIDAK ada di korpus — baca saja "—");
//           ADMIN NOTES (migrasi 027; dibekukan bila peserta diputus, b2628)
//   tengah  DOB, BEGIN/EXPIRED/EFFECTIVE DATE, VALUATION BEGIN/EXPIRED DATE, RECEIVED DATE, STNC TREATY
//   kanan   DATE OF LOSS, CLAIM RECEIVED DATE, MAX EXPIRED CLAIM, MAX CLAIM RECEIVED, DOCUMENT COMPLETE DATE, STNC DATA,
//           MAX DATA RECEIVED, CONFIRMATION DATE, tombol Edit Date (membuka isian tanggal yang sudah ada)
//   CLAIM ANALYST (NOTES, baca saja); kotak MEDICAL ADVISOR (diagnosa + RECOMMENDATION) di tahap Medical Check dan
//   Claim Analis — dapat diubah hanya Medical Advisor di tahap Medical Check (11-10-2026), selain itu baca saja
//   sub-tab Adjustment (b16543) | Document (b23293)
//
// ⛔ Nilai uang tidak pernah disentuh aritmetika (ADR-U-0003, ADR-U-0016); Percent Claim dan Claim Paid diturunkan SERVER.

import { useEffect, useState, type ReactNode } from 'react'

import { DETAIL, OS_PEGA, TAHAP, TOMBOL, TOMBOL_KOMITE } from '../labels'
import { DialogTolakOutstanding } from './DialogTolakOutstanding'
import { PanelDokumenPeserta } from './PanelDokumenPeserta'
import { GridDiagnosa } from './GridDiagnosa'
import { keIsianTanggal, PanelTanggalKlaim } from './PanelTanggalKlaim'
import { tanggalTampil } from './PanelDataPolis'
import { Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { DESIMAL_TAK_DIBATASI, formatNumber, formatPersen } from '../../../../inti/frontend/lib/format'
import { bandingDesimal, normalisasiDesimal, persenKlaimTeks, saringKetikanAngka } from '../desimal'
import {
  ambilSpreadingPolis,
  BATAS_REKOMENDASI,
  bolehAddAdjustment,
  bolehUbahDiagnosa,
  bolehSerahkanDiLayar,
  bolehSimpanAdjustment,
  bolehSuntingClaimRNM,
  bolehUbahTanggalKlaim,
  jumlahUang,
  kataStatusBaris,
  serahkanKeKomite,
  simpanAdjustment,
  STATUS_OUTSTANDING,
  tambahPutaran,
  tolakBarisAdjustment,
  ubahTanggalKejadian,
  type BarisAdjustment,
  type Klaim,
  type Peserta,
  type IsianBarisAdjustment,
  type IsianTanggal,
  type Rasio,
  type SpreadingBaris,
  type Uang,
} from '../api'
import { kodeStatusGalat, pesanGalat } from '../../../../inti/frontend/klien'

const R = OS_PEGA.rincian
const A = OS_PEGA.adjustment

export interface RincianPesertaProps {
  klaim: Klaim
  /** Peserta yang dirinci; `p` supaya nama dan ikatannya sama dengan blok asalnya. */
  p: Peserta
  /** Membaca ulang klaim dari server sesudah sebuah aksi — status klaim adalah turunan, tidak ditebak di layar. */
  onMuatUlang: () => Promise<void>
  /**
   * ADMIN NOTES yang sedang diisi (layar Outstanding menyimpannya lewat tombol Save grid). Tanpa `onUbahCatatan`
   * medannya baca saja (layar Detail).
   */
  catatanAdmin?: string
  onUbahCatatan?: (v: string) => void
  /**
   * Lima tanggal yang diisi pemakai (layar Outstanding, keputusan work owner 11-10-2026). Bila `onUbahTanggal` ada,
   * kolom kanan menjadi ISIAN tanggal (sebelum Save to RNM) dan tombol Edit Date tidak ditampilkan; tanpa itu (layar
   * Detail) nilainya baca saja dan Edit Date tetap.
   */
  tanggal?: IsianTanggal
  onUbahTanggal?: (t: IsianTanggal) => void
  /**
   * Claim RNM + data bank per baris adjustment (layar Outstanding, disimpan tombol Save grid — 11-10-2026). Tanpa
   * `onUbahAdjustment` isian adjustment baca saja.
   */
  isianAdjustment?: Record<string, IsianBarisAdjustment>
  onUbahAdjustment?: (adjID: string, isi: IsianBarisAdjustment) => void
  /** Kolom wajib peserta ini yang kosong saat Save to RNM (nama kolom `KekuranganRNM.kolom`) — ditandai merah. */
  kolomKurang?: ReadonlySet<string>
  /**
   * Kotak MEDICAL ADVISOR dapat diubah (tahap Medical Check, pelaku Medical Advisor, kasus terbuka — dihitung pemanggil,
   * 11-10-2026). Gerbang `STS_REJECT` peserta ditambahkan di sini. Tanpa ini baca saja.
   */
  bolehMedis?: boolean
  /** RECOMMENDATION yang sedang diisi (disimpan tombol Save Medical Check). Tanpa `onUbahRekomendasi` baca saja. */
  rekomendasi?: string
  onUbahRekomendasi?: (v: string) => void
}

/** Isian tanggal awal dari nilai tersimpan peserta. */
export function tanggalAwal(p: Peserta): IsianTanggal {
  return {
    kejadian: keIsianTanggal(p.tanggalKejadian),
    terimaKlaim: keIsianTanggal(p.tanggalTerimaKlaim),
    dokumenLengkap: keIsianTanggal(p.tanggalDokumenLengkap),
    maksData: keIsianTanggal(p.tanggalMaksData ?? ''),
    konfirmasi: keIsianTanggal(p.tanggalKonfirmasi),
  }
}

/** Satu baris label-nilai bentuk form Pega. */
function Baris({ label, nilai, angka, kode }: { label: string; nilai: string | undefined; angka?: boolean; kode?: boolean }) {
  return (
    <div className="os__baris-form">
      <dt>{label}</dt>
      {/* kode: nomor panjang (Certificate No / Policy No) membungkus di dalam kolom nilai, TANPA gulir (11-10-2026). */}
      <dd className={angka === true ? 'os__angka' : kode === true ? 'rincian__nilai-kode' : undefined}>{nilai === undefined || nilai === '' ? '—' : nilai}</dd>
    </div>
  )
}

/** Peserta diputus (STS_REJECT 1/2) membekukan isian catatan (b2628). */
function diputus(p: Peserta): boolean {
  return p.kodeStatus === '1' || p.kodeStatus === '2'
}

type SubTab = 'adjustment' | 'dokumen'

/** Satu isian tanggal di kolom kanan; wajib bertanda bintang; `salah` = kosong saat Save to RNM (merah). */
function IsianTanggalBaris({
  label,
  nilai,
  wajib,
  salah,
  pesan,
  boleh,
  onUbah,
}: {
  label: string
  nilai: string
  wajib?: boolean
  salah?: boolean
  /** Pesan galat isian (mis. DOL di luar periode valuation) — merah, tepat di bawah isian. */
  pesan?: string | null
  boleh: boolean
  onUbah: (v: string) => void
}) {
  const adaPesan = pesan !== undefined && pesan !== null && pesan !== ''
  const merah = adaPesan || (salah === true && nilai.trim() === '')
  return (
    <label className="os__baris-form">
      <span className="os__label">
        {label}
        {wajib === true && (
          <span className="os__wajib" aria-hidden="true">
            *
          </span>
        )}
      </span>
      {boleh ? (
        <span className="os__isian-bungkus">
          <input
            type="date"
            className={'field__input os__isian-tanggal' + (merah ? ' field__input--error' : '')}
            value={nilai}
            required={wajib}
            aria-label={label}
            aria-invalid={merah}
            onChange={(e) => onUbah(e.target.value)}
          />
          {adaPesan && (
            <span role="alert" className="os__galat-isian">
              {pesan}
            </span>
          )}
        </span>
      ) : (
        <span>{tanggalTampil(nilai)}</span>
      )}
    </label>
  )
}

export function RincianPeserta({
  klaim,
  p,
  onMuatUlang,
  catatanAdmin,
  onUbahCatatan,
  tanggal,
  onUbahTanggal,
  isianAdjustment,
  onUbahAdjustment,
  kolomKurang,
  bolehMedis = false,
  rekomendasi,
  onUbahRekomendasi,
}: RincianPesertaProps) {
  const kurang = (kolom: string) => kolomKurang?.has(kolom) === true
  const [galat, setGalat] = useState<string | null>(null)
  const [galatTanggal, setGalatTanggal] = useState<string | null>(null)
  const [tanggalSibuk, setTanggalSibuk] = useState(false)
  const [menolak, setMenolak] = useState<string | null>(null)
  const [dialogTolak, setDialogTolak] = useState<string | null>(null)
  const [galatTolak, setGalatTolak] = useState<string | null>(null)
  const [menyerahkan, setMenyerahkan] = useState<string | null>(null)
  const [memutar, setMemutar] = useState(false)
  const [mengaksep, setMengaksep] = useState(false)
  const [editTanggal, setEditTanggal] = useState(false)
  // Kekurangan dokumen (Send to Medical Check, 11-10-2026): sub-tab Document langsung terbuka.
  const [subTab, setSubTab] = useState<SubTab>(() =>
    kolomKurang?.has('Document') === true && kolomKurang.has(KOLOM_CLAIM_RNM) !== true ? 'dokumen' : 'adjustment',
  )
  // Claim RNM kurang saat Save to RNM: baris pertama yang kosong langsung dibuka (11-10-2026).
  const [barisBuka, setBarisBuka] = useState<string | null>(() =>
    kolomKurang?.has(KOLOM_CLAIM_RNM) === true
      ? (p.baris.find((b) => claimRNMKosong(isianAdjustment?.[b.id]?.jumlah ?? jumlahUang(b.jumlahKlaim)))?.id ?? null)
      : null,
  )
  // Spreading PremiumList peserta ini (11-10-2026): TERLIPAT bawaan; dibaca SEKALI saat "View Spreading" pertama dibuka
  // (mintaSpreading), lalu disimpan di state peserta ini. null = sedang dimuat.
  const [spreading, setSpreading] = useState<SpreadingBaris | null>(null)
  const [galatSpreading, setGalatSpreading] = useState<string | null>(null)
  const [mintaSpreading, setMintaSpreading] = useState(false)
  useEffect(() => {
    if (!mintaSpreading) return
    let batal = false
    ambilSpreadingPolis(klaim.id, p.id)
      .then((s) => {
        if (!batal) setSpreading(s)
      })
      .catch((e: unknown) => {
        if (!batal) setGalatSpreading(pesanGalat(e) ?? 'Failed to load the spreading.')
      })
    return () => {
      batal = true
    }
  }, [mintaSpreading, klaim.id, p.id])

  /** Mengirim tanggal kejadian peserta ini; validasinya di server (`ValidasiDOL_Act`). */
  async function ubahTanggal(nilai: string): Promise<void> {
    // Galat lama dibersihkan LEBIH DULU, sebelum keluar karena kotak dikosongkan.
    setGalatTanggal(null)
    // Kosong tidak dikirim: backend menjawab 400 "tanggal kejadian belum diisi".
    if (nilai === '') return
    // ⛔ Menolak permintaan kedua selagi yang pertama terbang (pembacaan ulang yang kedaluwarsa menimpa layar).
    if (tanggalSibuk) return
    setTanggalSibuk(true)
    try {
      await ubahTanggalKejadian(klaim.id, p.id, nilai)
      await onMuatUlang()
    } catch (e) {
      // Kalimat server apa adanya: 422 menyebut jendela valuasi mana yang dilanggar.
      setGalatTanggal(pesanGalat(e) ?? 'Date of loss rejected.')
    } finally {
      setTanggalSibuk(false)
    }
  }

  async function aksep(): Promise<void> {
    setGalat(null)
    setMengaksep(true)
    try {
      const nomor = await simpanAdjustment(klaim.id, p.id)
      await onMuatUlang()
      setGalat(`Row accepted with number ${nomor}.`)
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Only the holder of this claim stage can accept its rows.'
          : kode === 409
            ? (pesanGalat(err) ?? 'The row cannot be accepted in its current state.')
            : kode === 501
              ? (pesanGalat(err) ?? 'Acceptance cannot be saved yet.')
              : 'Failed to accept the row.',
      )
    } finally {
      setMengaksep(false)
    }
  }

  async function putaranBaru(): Promise<void> {
    setGalat(null)
    setMemutar(true)
    try {
      await tambahPutaran(klaim.id, p.id)
      await onMuatUlang()
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Only ReasLifeSPV can add an adjustment row.'
          : kode === 409
            ? (pesanGalat(err) ?? 'An adjustment row cannot be added in this state.')
            : kode === 501
              ? 'The audit trail cannot be recorded yet.'
              : 'Failed to add an adjustment row.',
      )
    } finally {
      setMemutar(false)
    }
  }

  async function serahkan(adjID: string): Promise<void> {
    setGalat(null)
    setMenyerahkan(adjID)
    try {
      await serahkanKeKomite(klaim.id, p.id, adjID)
      await onMuatUlang()
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Your role cannot submit rows of this type to the Committee.'
          : kode === 409
            ? (pesanGalat(err) ?? 'The row was already submitted, or is no longer Outstanding.')
            : kode === 422
              ? // ⛔ Pesan 422 dicetak APA ADANYA ("Name of bank cannot be empty" - kalimat sistem lama).
                (pesanGalat(err) ?? 'Submission rejected.')
              : kode === 501
                ? 'Submission cannot be saved yet.'
                : 'Failed to submit the row to the Committee.',
      )
    } finally {
      setMenyerahkan(null)
    }
  }

  // OQ-M5 (GILIRAN-17): `Reject Outstanding` membuka dialog `RejectOSClaimLife_Sec`; dikirim bersama `Remarks`.
  async function tolak(adjID: string, remarks: string): Promise<void> {
    // Cegah klik ganda: satu penolakan sekali jalan.
    if (menolak !== null) return
    setGalatTolak(null)
    setMenolak(adjID)
    try {
      await tolakBarisAdjustment(klaim.id, adjID, remarks)
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalatTolak(
        kode === 400
          ? (pesanGalat(err) ?? 'Remarks are required.')
          : kode === 403
            ? 'Only ReasLifeAdmin can reject a row.'
            : kode === 409
              ? (pesanGalat(err) ?? 'The row is already decided and cannot be rejected again.')
              : kode === 422
                ? 'The claim has no claim number yet.'
                : kode === 501
                  ? 'The audit trail cannot be recorded yet.'
                  : 'Failed to reject the row.',
      )
      setMenolak(null)
      return
    }
    // Penolakan SUDAH tersimpan: dialog ditutup; kegagalan membaca ulang dilaporkan di layar.
    setDialogTolak(null)
    try {
      await onMuatUlang()
    } catch {
      setGalat('The row was rejected, but the claim could not be reloaded.')
    } finally {
      setMenolak(null)
    }
  }

  const bolehTanggal = bolehUbahTanggalKlaim(klaim)
  // Admin Notes dapat diisi HANYA bila halaman memberi onUbahCatatan (Admin + tahap Outstanding, bolehUbahCatatanAdmin)
  // DAN peserta belum diputus (b2628 `.STS_REJECT=='1' || '2'`). Sesudah Save to RNM TETAP terbuka: XML b2588/b2628 tidak
  // menguncinya (beda dengan tanggal, yang memakai `CLAIM_NO!=''`).
  const catatanBoleh = onUbahCatatan !== undefined && !diputus(p)
  // Isian tanggal langsung (layar Outstanding): terbuka hanya sebelum Save to RNM (OQ-M1) dan klaim belum bernomor.
  const isianTanggal = onUbahTanggal !== undefined ? (tanggal ?? tanggalAwal(p)) : null
  const tanggalBoleh = isianTanggal !== null && bolehTanggal && klaim.nomorKlaim.trim() === ''

  return (
    <div className="rincian">
      {dialogTolak !== null && (
        <DialogTolakOutstanding
          subjek={p.namaTertanggung || p.nomorSertifikat}
          sibuk={menolak === dialogTolak}
          galat={galatTolak}
          onKirim={(remarks) => void tolak(dialogTolak, remarks)}
          onBatal={() => {
            setDialogTolak(null)
          }}
        />
      )}
      {galat !== null && (
        <p role="alert" className="alert alert--warn">
          {galat}
        </p>
      )}

      {/* Tata letak 11-10-2026 ("tidak perlu ikuti pega 100%, hanya jadikan referensi"): TIGA kelompok berjudul, kolom sama
          lebar, label selebar sama (nilai sejajar satu garis); di bawahnya Location of Death | Admin Notes sama tinggi.
          Responsif 3 -> 2 -> 1 kolom (claimlife.css .rincian__grid). */}
      <div className="rincian__grid">
        <section className="rincian__kelompok" aria-label={R.kelompokTertanggung}>
          <h4 className="rincian__judul">{R.kelompokTertanggung}</h4>
          <dl className="os__kolom">
            <Baris label={R.sertifikat} nilai={p.nomorSertifikat} kode />
            <Baris label={R.polis} nilai={p.nomorPolis} kode />
            <Baris label={R.plan} nilai={p.plan} />
            <Baris label={R.umur} nilai={p.umur} />
            <Baris label={R.sex} nilai={p.jenisKelamin} />
            <Baris label={R.dob} nilai={tanggalTampil(p.tanggalLahir)} />
          </dl>
          {/* Susunan lama (keputusan work owner 11-10-2026): LOCATION OF DEATH lalu ADMIN NOTES DI BAWAHNYA, bertumpuk di kolom
              kiri di bawah data Insured. Keduanya dapat diubah HANYA oleh Admin selama kasus di inbox Admin (tahap
              Outstanding) — selain itu baca saja. */}
          {/* ⚠️ LOCATION OF DEATH ada di screenshot 11-10-2026 tetapi TIDAK di korpus XML dan tanpa kolom: baca saja "—"
              sampai kolomnya disetujui (jatah migrasi Claim Life habis). */}
          <fieldset className="rincian__kotak">
            <legend className="os__label">{R.lokasiKematian}</legend>
            <dl className="os__kolom">
              <Baris label={R.tempatKematian} nilai="" />
              <Baris label={R.negara} nilai="" />
              <Baris label={R.provinsi} nilai="" />
              <Baris label={R.kota} nilai="" />
            </dl>
          </fieldset>
          <div className="rincian__kotak rincian__catatan">
            {catatanBoleh ? (
              <label className="rincian__catatan-isian">
                <span className="os__label">{R.catatanAdmin}</span>
                <textarea
                  className="field__input os__isian-area"
                  rows={3}
                  maxLength={1000}
                  value={catatanAdmin ?? p.catatanAdmin ?? ''}
                  onChange={(e) => onUbahCatatan?.(e.target.value)}
                />
              </label>
            ) : (
              <>
                <span className="os__label">{R.catatanAdmin}</span>
                {/* Baca saja: teksnya saja, tanpa kotak isian. */}
                <p className="rincian__catatan-teks">{(catatanAdmin ?? p.catatanAdmin ?? '') || '—'}</p>
              </>
            )}
          </div>
        </section>
        <section className="rincian__kelompok" aria-label={R.kelompokPertanggungan}>
          <h4 className="rincian__judul">{R.kelompokPertanggungan}</h4>
          <dl className="os__kolom">
            <Baris label={R.mulai} nilai={tanggalTampil(p.tanggalMulai)} />
            <Baris label={R.expired} nilai={tanggalTampil(p.tanggalExpired)} />
            <Baris label={R.efektif} nilai={tanggalTampil(p.tanggalEfektif)} />
            <Baris label={R.valuasiMulai} nilai={tanggalTampil(p.valuasiMulai)} />
            <Baris label={R.valuasiSelesai} nilai={tanggalTampil(p.valuasiSelesai)} />
            <Baris label={R.terima} nilai={tanggalTampil(p.tanggalTerima)} />
            <Baris label={R.stncTreaty} nilai={tanggalTampil(p.stnc)} />
          </dl>
        </section>
        <section className="rincian__kelompok" aria-label={R.kelompokTanggalKlaim}>
          <h4 className="rincian__judul">{R.kelompokTanggalKlaim}</h4>
          <dl className="os__kolom">
            {/* DATE OF LOSS*, CLAIM RECEIVED DATE*, CONFIRMATION DATE*, DOCUMENT COMPLETE DATE, MAX DATA RECEIVED (migrasi
                028) DIISI PEMAKAI sebelum Save to RNM; sesudahnya baca saja. Bintang = wajib di Save to RNM. DOL diperiksa
                LANGSUNG terhadap periode valuation (galatDOL). */}
            {isianTanggal !== null ? (
              <IsianTanggalBaris label={R.dol} nilai={isianTanggal.kejadian} wajib salah={kurang('Date of Loss')}
                pesan={tanggalBoleh ? galatDOL(klaim, p, isianTanggal.kejadian) : null} boleh={tanggalBoleh}
                onUbah={(v) => onUbahTanggal?.({ ...isianTanggal, kejadian: v })} />
            ) : (
              <Baris label={R.dol} nilai={tanggalTampil(p.tanggalKejadian)} />
            )}
            {isianTanggal !== null ? (
              <IsianTanggalBaris label={R.terimaKlaim} nilai={isianTanggal.terimaKlaim} wajib salah={kurang('Claim Received Date')}
                boleh={tanggalBoleh}
                onUbah={(v) => onUbahTanggal?.({ ...isianTanggal, terimaKlaim: v })} />
            ) : (
              <Baris label={R.terimaKlaim} nilai={tanggalTampil(p.tanggalTerimaKlaim)} />
            )}
            {isianTanggal !== null ? (
              <IsianTanggalBaris label={R.konfirmasi} nilai={isianTanggal.konfirmasi} wajib salah={kurang('Confirmation Date')}
                boleh={tanggalBoleh}
                onUbah={(v) => onUbahTanggal?.({ ...isianTanggal, konfirmasi: v })} />
            ) : (
              <Baris label={R.konfirmasi} nilai={tanggalTampil(p.tanggalKonfirmasi)} />
            )}
            {isianTanggal !== null ? (
              <IsianTanggalBaris label={R.dokumenLengkap} nilai={isianTanggal.dokumenLengkap} boleh={tanggalBoleh}
                onUbah={(v) => onUbahTanggal?.({ ...isianTanggal, dokumenLengkap: v })} />
            ) : (
              <Baris label={R.dokumenLengkap} nilai={tanggalTampil(p.tanggalDokumenLengkap)} />
            )}
            {isianTanggal !== null ? (
              <IsianTanggalBaris label={R.maksData} nilai={isianTanggal.maksData} boleh={tanggalBoleh}
                onUbah={(v) => onUbahTanggal?.({ ...isianTanggal, maksData: v })} />
            ) : (
              <Baris label={R.maksData} nilai={tanggalTampil(p.tanggalMaksData)} />
            )}
            <Baris label={R.maksKadaluarsa} nilai={p.maksKlaimKadaluarsa} />
            {/* Butir bk — dihitung backend saat baca. ⛔ Alasan teknis (penandaTerimaKlaimAlasan) TIDAK dirender (11-10-2026):
                bila tidak dihitung, tampil "—" dengan keterangan ringkas di tooltip. */}
            {p.penandaTerimaKlaim === '' && p.penandaTerimaKlaimAlasan !== '' ? (
              <div className="os__baris-form">
                <dt>{R.maksTerima}</dt>
                <dd title={TIDAK_TERSEDIA_AMBANG}>
                  — <span className="os__info" aria-label={TIDAK_TERSEDIA_AMBANG}>ⓘ</span>
                </dd>
              </div>
            ) : (
              <Baris label={R.maksTerima} nilai={p.penandaTerimaKlaim} />
            )}
            {/* STNC DATA dihitung ValidasiSTNC_Act di Pega — belum ditiru. */}
            <Baris label={R.stncData} nilai="" />
            {/* Edit Date hanya di layar tanpa isian langsung (layar Detail). */}
            {isianTanggal === null && (
              <p>
                <button type="button" className="btn btn--ghost btn--sm" onClick={() => setEditTanggal(!editTanggal)}>
                  {R.editDate}
                </button>
              </p>
            )}
          </dl>
        </section>
      </div>
      {/* Edit Date b14115 → `EditDateClaimLife_Section`; validasinya `ValidasiDOL_Act` di SERVER. ⛔ Tanggal kejadian
          milik PESERTA. Isiannya terkunci sesudah Save to RNM (OQ-M1). */}
      {editTanggal && (
        <div className="register__kartu rincian__edit">
          <label className="os__baris-form">
            <span className="os__label">{DETAIL.ubahTanggal}</span>
            <input
              type="date"
              // ⛔ `key` memaksa kotaknya LAHIR ULANG ketika tanggal tersimpan berubah (defaultValue dibaca sekali).
              key={`${p.id}:${p.tanggalKejadian}`}
              // ⛔ Backend memulangkan tanggal BERJAM; kotak tanggal menolaknya diam-diam.
              defaultValue={keIsianTanggal(p.tanggalKejadian)}
              // ⛔ Butir bj: DOL bergerbang SAMA dengan tiga tanggal lainnya (b1000) - Admin, tahap Outstanding.
              disabled={tanggalSibuk || !bolehTanggal}
              onChange={(e) => {
                void ubahTanggal(e.target.value)
              }}
            />
          </label>
          {galatTanggal !== null && <span role="alert"> {galatTanggal}</span>}
          {/* Tiga tanggal lain dialog yang sama — `Save` b1910. */}
          <PanelTanggalKlaim
            key={`${p.id}:${p.tanggalTerimaKlaim}:${p.tanggalDokumenLengkap}:${p.tanggalKonfirmasi}`}
            klaimID={klaim.id}
            peserta={p}
            boleh={bolehTanggal}
            sesudahSimpan={() => {
              void onMuatUlang()
            }}
          />
        </div>
      )}

      {/* CLAIM ANALYST hanya di tahap Claim Analis (keputusan work owner 11-10-2026: "kolom claim analyst hanya tampil di
          inbox claim analyst, di admin hide saja"); di Pega bagian ini tampil bagi pyPosition ReasLifeSPV (b15474). */}
      {klaim.tahap === TAHAP.claimAnalis && (
        <div className="os__baris-form os__baris-form--tumpuk">
          <span className="os__label">{R.analis}</span>
          <span>{p.catatanAnalis || '—'}</span>
        </div>
      )}

      {/* Kotak MEDICAL ADVISOR (11-10-2026) — grid diagnosa butir bd. Tampil di Medical Check dan Claim Analis; dapat diubah
          hanya Medical Advisor di Medical Check (b4908, b6401) selama peserta belum diputus (STS_REJECT b4682). */}
      {(klaim.tahap === TAHAP.medicalCheck || klaim.tahap === TAHAP.claimAnalis) && (
        <fieldset className="rincian__kotak rincian__medis">
          <legend className="os__label">{R.medis}</legend>
          <GridDiagnosa
            klaimID={klaim.id}
            peserta={p}
            onBerubah={onMuatUlang}
            boleh={bolehMedis && bolehUbahDiagnosa(p)}
          />
          {bolehMedis && !bolehUbahDiagnosa(p) && (
            <p className="diagnosa-grid__catatan">This insured has been decided; the medical data is read-only.</p>
          )}
          {/* RECOMMENDATION b7316/b7323 di bawah grid diagnosa; dapat diisi bila bolehMedis dan peserta belum diputus
              (pyDisabledWhen STS_REJECT), disimpan tombol Save Medical Check. Selain itu teks saja. */}
          {bolehMedis && bolehUbahDiagnosa(p) && onUbahRekomendasi !== undefined ? (
            <label className="rincian__rekomendasi">
              <span className="os__label">{R.rekomendasi}</span>
              <textarea
                className="field__input os__isian-area"
                rows={3}
                maxLength={BATAS_REKOMENDASI}
                value={rekomendasi ?? p.rekomendasi ?? ''}
                onChange={(e) => onUbahRekomendasi(e.target.value)}
              />
            </label>
          ) : (
            <div className="rincian__rekomendasi">
              <span className="os__label">{R.rekomendasi}</span>
              <p className="rincian__catatan-teks">{(rekomendasi ?? p.rekomendasi ?? '') || '—'}</p>
            </div>
          )}
        </fieldset>
      )}

      <StripTab<SubTab>
        tab={['adjustment', 'dokumen']}
        aktif={subTab}
        onPilih={setSubTab}
        label={(t) => (t === 'adjustment' ? R.tabAdjustment : R.tabDokumen)}
      />

      {subTab === 'dokumen' && (
        // Dokumen pendukung peserta — `Section/DocumentLife.xml` + `LoadDocumentLife_ACT`.
        <PanelDokumenPeserta klaimID={klaim.id} pesertaID={p.id} dokumen={p.dokumen} onBerubah={onMuatUlang} />
      )}

      {subTab === 'adjustment' && (
        <>
          {p.baris.length === 0 ? (
            <p>No adjustment rows yet.</p>
          ) : (
            <div className="inbox__gulir">
              <table className="inbox__tabel">
                <thead>
                  <tr>
                    {R.kolomAdjustment.map((k) => (
                      <th key={k}>{k}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {p.baris.map((b) => {
                    const buka = barisBuka === b.id
                    return [
                      // Grid `.AdjustmentList` b17126: Status (teks), OutStanding Date (jejak), Acceptation No / Date.
                      <tr
                        key={b.id}
                        className={
                          'inbox__baris' +
                          (buka ? ' os__baris--pilih' : '') +
                          (kurang(KOLOM_CLAIM_RNM) && claimRNMKosong(isianAdjustment?.[b.id]?.jumlah ?? jumlahUang(b.jumlahKlaim))
                            ? ' os__baris--kurang'
                            : '')
                        }
                        onClick={() => setBarisBuka(buka ? null : b.id)}
                      >
                        <td>{kataStatusBaris(b.kodeStatus)}</td>
                        <td>{tanggalTampil(b.tanggalOutstanding)}</td>
                        <td>{b.nomorAkseptasi || '—'}</td>
                        <td>{tanggalTampil(b.tanggalAkseptasi)}</td>
                      </tr>,
                      buka ? (
                        <tr key={`${b.id}:isian`} className="os__rincian">
                          <td colSpan={4}>
                            <IsianAdjustment
                              klaim={klaim}
                              b={b}
                              isi={isianAdjustment?.[b.id]}
                              spreading={spreading}
                              galatSpreading={galatSpreading}
                              onMuatSpreading={() => setMintaSpreading(true)}
                              onUbah={onUbahAdjustment === undefined ? undefined : (isi) => onUbahAdjustment(b.id, isi)}
                              aksi={
                                <>
                                  {/* Reject Outstanding b15101: Admin && CLAIM_NO != '' && STS_REJECT == 0 (b15399). */}
                                  {b.status === STATUS_OUTSTANDING && klaim.nomorKlaim ? (
                                    <button
                                      type="button"
                                      className="btn btn--ghost btn--sm"
                                      disabled={menolak === b.id}
                                      onClick={() => {
                                        setGalatTolak(null)
                                        setDialogTolak(b.id)
                                      }}
                                    >
                                      {menolak === b.id ? 'Rejecting…' : TOMBOL.tolakOutstanding}
                                    </button>
                                  ) : null}
                                  {/* Send To Komite b15480 — HANYA di tahap Claim Analis (keputusan work owner 11-10-2026: "iya itu hanya
                                      di claim analis saja"); server menolak 409 di tahap lain. */}
                                  {klaim.tahap === TAHAP.claimAnalis && bolehSerahkanDiLayar(b) ? (
                                    <button
                                      type="button"
                                      className="btn btn--ghost btn--sm"
                                      disabled={menyerahkan === b.id}
                                      onClick={() => void serahkan(b.id)}
                                    >
                                      {menyerahkan === b.id ? 'Submitting…' : TOMBOL_KOMITE.serahkan}
                                    </button>
                                  ) : null}
                                </>
                              }
                            />
                          </td>
                        </tr>
                      ) : null,
                    ]
                  })}
                </tbody>
              </table>
            </div>
          )}

          {/* ⭐ Jalur akseptasi Claim Life — `SaveAdjustment_Act`. */}
          {bolehSimpanAdjustment(p, p.baris) && (
            <p>
              <button type="button" className="btn btn--primary btn--sm" disabled={mengaksep} onClick={() => void aksep()}>
                {mengaksep ? 'Saving…' : 'Save Adjustment'}
              </button>
            </p>
          )}

          {/* `Add` b17937 hanya bila `pyPosition=='ReasLifeSPV'` (b18160): jalur PUTARAN. ⛔ `Delete` b19120 TIDAK
              DIRENDER (OQ-N7, ADR-U-0031). */}
          {bolehAddAdjustment(klaim, p) && (
            <p>
              <button type="button" className="btn btn--ghost btn--sm" disabled={memutar} onClick={() => void putaranBaru()}>
                {DETAIL.tambahAdjustment}
              </button>
            </p>
          )}

          {/* ⛔ Enam total peserta TIDAK ditampilkan di layar mana pun (keputusan work owner 11-10-2026 "buang saja dari semua
              halaman"; di Pega wadahnya `pyContainerVisibleWhen 1=2`, ClaimLifeDetailGCNM b20403). Backend
              models.HitungTotalPeserta dan gerbang mata uangnya tetap. */}
        </>
      )}
    </div>
  )
}

/** Angka uang bergaya Indonesia (ribuan titik, desimal koma) lewat pemformat inti — tanpa float. */
function angka(v: Uang | undefined): string {
  const t = formatNumber(jumlahUang(v), DESIMAL_TAK_DIBATASI)
  return t === '' ? '—' : t
}

/** Rasio (persen/rate) bergaya Indonesia lewat pemformat inti. */
function rasio(r: Rasio | undefined): string {
  const t = formatNumber(r?.value ?? '', DESIMAL_TAK_DIBATASI)
  return t === '' ? '—' : t
}

/**
 * Galat isian Claim RNM di layar (11-10-2026): bukan angka, atau melebihi Share Nusantara Re (nilai share berformat).
 * Perbandingan desimal TEKS (desimal.ts), tanpa float. null = sah. Server tetap menolak 422 sebagai lapis kedua.
 */
export function galatClaimRNM(nilai: string, share: Uang | undefined): string | null {
  // Keputusan work owner 11-10-2026: Claim RNM WAJIB dan > KLAIM_RNM_MINIMAL. Tidak menghalangi Save grid
  // (adaIsianTidakSah melewatkannya); ditegakkan Save to RNM.
  if (nilai.trim() === '') return PESAN_CLAIM_RNM_WAJIB
  const n = normalisasiDesimal(nilai)
  if (n === null) return PESAN_CLAIM_RNM_BUKAN_ANGKA
  if (bandingDesimal(n, KLAIM_RNM_MINIMAL) <= 0) return PESAN_CLAIM_RNM_WAJIB
  const s = normalisasiDesimal(jumlahUang(share))
  if (s !== null && bandingDesimal(n, s) > 0) {
    return `Claim RNM cannot be greater than Share Nusantara Re (${angka(share)}).`
  }
  return null
}

/** Claim RNM bukan angka — SATU-SATUNYA galat Claim RNM yang menghalangi Save grid (format rusak, tak dapat disimpan). */
export const PESAN_CLAIM_RNM_BUKAN_ANGKA = 'Claim RNM must be a number.'

/** Pesan Claim RNM kosong / nol (keputusan work owner 11-10-2026). */
export const PESAN_CLAIM_RNM_WAJIB = 'Claim RNM is required.'
/** Claim RNM harus LEBIH BESAR dari ini — SAMA dengan backend `services.KlaimRNMMinimal` (ubah keduanya bila 0 kelak boleh). */
export const KLAIM_RNM_MINIMAL = '0'

/**
 * Percent Claim dan Claim Paid DI LAYAR saat Claim RNM diketik (keputusan work owner 11-10-2026: "otomatis terisi saat user
 * mengisi Claim RNM"). Isian sah: Claim Paid = Claim RNM, Percent = persenKlaimTeks (aturan sama dengan backend).
 * Isian TIDAK sah: keduanya "—" (bukan nilai lama) — angka yang tampil selalu milik isian yang sedang terlihat.
 */
export function turunanClaimRNM(nilai: string, share: Uang | undefined): { persen: string; dibayar: string } {
  if (galatClaimRNM(nilai, share) !== null) return { persen: '—', dibayar: '—' }
  const n = normalisasiDesimal(nilai) ?? ''
  const s = normalisasiDesimal(jumlahUang(share))
  const persen = s === null ? '' : persenKlaimTeks(n, s)
  return {
    persen: persen === '' ? '—' : formatPersen(persen, DESIMAL_TAK_DIBATASI),
    dibayar: formatNumber(n, DESIMAL_TAK_DIBATASI) || '—',
  }
}

/**
 * DOL di luar periode valuation — validasi LANGSUNG di layar (keputusan work owner 11-10-2026: "date of loss harus ada di
 * periode valuation date, tidak boleh di luar. jika di luar tidak bisa lanjut dan ada pesan").
 *
 * Aturan SAMA dengan backend `services.DOLDalamJendela` (Save grid, Edit Date, Save to RNM) — keputusan work owner
 * 11-10-2026 "tanggal awal valuation masih termasuk": INKLUSIF di kedua ujung, tanpa geser, per tanggal:
 *   QR / QP  jendela GROSS:  BEGIN <= DOL <= EXPIRED
 *   TR / TP  jendela RETRO:  BEGIN <= DOL <= EXPIRED
 * Type lain / tanggal valuasi kosong: null (server yang memutuskan). Kalimat = `services.KalimatDOLDiLuar`.
 */
export function galatDOL(klaim: Klaim, p: Peserta, dol: string): string | null {
  const d = dol.trim().slice(0, 10)
  if (!/^\d{4}-\d{2}-\d{2}$/.test(d)) return null
  const tipe = (klaim.type ?? '').trim().toUpperCase()
  const gross = tipe === 'QR' || tipe === 'QP'
  if (!gross && tipe !== 'TR' && tipe !== 'TP') return null
  const awal = keIsianTanggal((gross ? p.valuasiMulai : p.valuasiRetroMulai) ?? '')
  const akhir = keIsianTanggal((gross ? p.valuasiSelesai : p.valuasiRetroSelesai) ?? '')
  if (awal === '' || akhir === '') return null
  if (d >= awal && d <= akhir) return null
  return `Date of Loss must be within the valuation period (${tanggalTampil(awal)} – ${tanggalTampil(akhir)}).`
}

/** Keterangan MAX CLAIM RECEIVED yang tidak dapat dihitung (ambang produk belum diisi / tak terbaca) — tanpa teks galat. */
export const TIDAK_TERSEDIA_AMBANG = 'Not available: product threshold not set'

/** Nama kolom Claim RNM di `KekuranganRNM.kolom` (backend services.PeriksaSimpanRNM). */
export const KOLOM_CLAIM_RNM = 'Claim RNM'

/** Claim RNM kosong / tidak > minimal — kekurangan wajib Save to RNM. */
export function claimRNMKosong(nilai: string): boolean {
  return galatClaimRNM(nilai, undefined) === PESAN_CLAIM_RNM_WAJIB
}

/** Isian awal satu baris dari nilai tersimpan. */
export function isianBarisAwal(b: BarisAdjustment): IsianBarisAdjustment {
  return {
    jumlah: jumlahUang(b.jumlahKlaim),
    namaBank: b.namaBank ?? '',
  }
}

/**
 * Bagian Spreading isian adjustment — bentuk tab Spreading PremiumList Life (disalin bentuknya, bukan diimpor): rincian
 * per lapisan (Reins Type, Retroceded Share hanya di baris pertama lapisan) lalu per reinsurer, ditambah rekap per
 * reinsurer (dijumlahkan server). Angka rata kanan, bergaya Indonesia.
 */
export function BagianSpreading({ s, galat }: { s: SpreadingBaris | null | undefined; galat?: string | null }) {
  if (galat !== undefined && galat !== null) {
    return (
      <p role="alert" className="alert alert--warn">
        {galat}
      </p>
    )
  }
  if (s === null) return <Memuat />
  const lapisan = s?.lapisan ?? []
  if (lapisan.length === 0) {
    // Polis warisan Pega / peserta yang belum di-spreading di PremiumList Life.
    return <Kosong pesan="No spreading data for this insured in the policy." />
  }
  return (
    <>
      <h6 className="os__subjudul">{A.spreadingDetail}</h6>
      <div className="inbox__gulir">
        <table className="inbox__tabel os__tabel-angka">
          <thead>
            <tr>
              {A.kolomSpreading.map((k, i) => (
                <th key={k} className={A.kolomSpreadingAngka[i] === true ? 'os__angka' : undefined}>
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {lapisan.flatMap((l, li) =>
              (l.reinsurer.length === 0 ? [null] : l.reinsurer).map((r, ri) => (
                <tr key={`${li}:${ri}`} className={ri === 0 ? 'os__lapisan-awal' : undefined}>
                  <td>{ri === 0 ? l.treatyTypeName || '—' : ''}</td>
                  <td className="os__angka">{ri === 0 ? angka(l.retrocededShare) : ''}</td>
                  <td>{r?.nama || '—'}</td>
                  <td className="os__angka">{r === null ? '—' : rasio(r.pctShare)}</td>
                  <td className="os__angka">{r === null ? '—' : angka(r.amount)}</td>
                  <td className="os__angka">{r === null ? '—' : rasio(r.rate)}</td>
                  <td className="os__angka">{r === null ? '—' : rasio(r.komisi)}</td>
                  <td className="os__angka">{r === null ? '—' : rasio(r.ovrComm)}</td>
                  <td className="os__angka">{r === null ? '—' : angka(r.gross)}</td>
                  <td className="os__angka">{r === null ? '—' : angka(r.net)}</td>
                </tr>
              )),
            )}
          </tbody>
        </table>
      </div>
      <h6 className="os__subjudul">{A.spreadingSummary}</h6>
      <div className="inbox__gulir">
        <table className="inbox__tabel os__tabel-angka">
          <thead>
            <tr>
              <th>{A.kolomRekap[0]}</th>
              <th>{A.kolomRekap[1]}</th>
              <th className="os__angka">{A.kolomRekap[2]}</th>
              <th className="os__angka">{A.kolomRekap[3]}</th>
            </tr>
          </thead>
          <tbody>
            {(s?.rekap ?? []).map((r, i) => (
              <tr key={i}>
                <td>{r.reinsurer || '—'}</td>
                <td>{r.currency || '—'}</td>
                <td className="os__angka">{angka(r.gross)}</td>
                <td className="os__angka">{angka(r.net)}</td>
              </tr>
            ))}
          </tbody>
          {(s?.total ?? []).length > 0 && (
            <tfoot>
              {(s?.total ?? []).map((t, i) => (
                <tr key={i} className="os__total">
                  <th scope="row">{A.total}</th>
                  <td>{t.currency || '—'}</td>
                  <td className="os__angka">{angka(t.gross)}</td>
                  <td className="os__angka">{angka(t.net)}</td>
                </tr>
              ))}
            </tfoot>
          )}
        </table>
      </div>
    </>
  )
}

/**
 * Isian satu baris adjustment — dialog `AdjustmentDetail_Section` (FlowAction `Adjustment_Detail`), tata letak screenshot
 * Pega 11-10-2026:
 *   kiri    Currency b1484, Percent Claim (%) b1649, Claim Paid b1972
 *   tengah  Share Nusantara Re b2679, Ceding Retention b2843, Claim RNM b3005 (SATU-SATUNYA angka yang dapat diisi)
 *   kanan   Sum Reasured b4133, Sum Insured b4296
 *   bawah   Payable To b5929 (baca saja), Name of Bank b6122, Branch of Bank b6870, Account No b7049 (isian TEKS — master
 *           BANKACCOUNT belum diketahui); bagian Spreading = spreading PremiumList Life peserta (11-10-2026, tanpa hitung)
 * ⭐ SEBELUM Save to RNM Claim RNM dan data bank dapat diisi dan DISIMPAN TOMBOL SAVE GRID (tombol "Save to Outstanding"
 * dicabut 11-10-2026); SESUDAHNYA seluruhnya baca saja.
 */
export function IsianAdjustment({
  klaim,
  b,
  isi,
  onUbah,
  aksi,
  spreading,
  galatSpreading,
  onMuatSpreading,
  spreadingTerbuka,
}: {
  klaim: Klaim
  b: BarisAdjustment
  /** Spreading PremiumList peserta (dimuat RincianPeserta); null = sedang dimuat. */
  spreading?: SpreadingBaris | null
  galatSpreading?: string | null
  /** Dipanggil saat "View Spreading" dibuka — RincianPeserta memuat spreading SEKALI per peserta. */
  onMuatSpreading?: () => void
  /** Keadaan awal bagian Spreading (bawaan TERLIPAT). */
  spreadingTerbuka?: boolean
  /** Isian yang sedang diubah (dipegang halaman); kosong = nilai tersimpan. */
  isi?: IsianBarisAdjustment
  /** Tanpa ini isian baca saja (layar Detail). */
  onUbah?: (isi: IsianBarisAdjustment) => void
  /** Tombol baris (Reject Outstanding, Send To Komite) — kaki dialog. */
  aksi?: ReactNode
}) {
  const nilai = isi ?? isianBarisAwal(b)
  const boleh = onUbah !== undefined && bolehSuntingClaimRNM(klaim, b)

  const galatRnm = boleh ? galatClaimRNM(nilai.jumlah, b.shareNusantaraRe) : null
  // Isian terbuka: Percent / Claim Paid dihitung ulang dari ketikan; baca saja: nilai turunan server.
  const turunan = boleh
    ? turunanClaimRNM(nilai.jumlah, b.shareNusantaraRe)
    : { persen: formatPersen(b.persenKlaim ?? '', DESIMAL_TAK_DIBATASI), dibayar: angka(b.klaimDibayar) }
  const [lihatSpreading, setLihatSpreading] = useState(spreadingTerbuka === true)
  const teks = (label: string, kunci: keyof IsianBarisAdjustment, mode?: 'decimal') => {
    const angkaIsian = mode === 'decimal'
    const salah = angkaIsian && galatRnm !== null
    return (
      <label className="os__baris-form">
        <span className="os__label">
          {label}
          {/* Claim RNM WAJIB sebelum Save to RNM (keputusan work owner 11-10-2026). */}
          {kunci === 'jumlah' && (
            <span className="os__wajib" aria-hidden="true">
              *
            </span>
          )}
        </span>
        {boleh ? (
          <span className="os__isian-bungkus">
            <input
              type="text"
              // Kelas isian inti (field__input) — tinggi, border, radius, padding seragam dengan isian lain.
              className={
                'field__input ' + (angkaIsian ? 'os__isian-angka' : 'os__isian-teks') + (salah ? ' field__input--error' : '')
              }
              inputMode={mode}
              value={nilai[kunci]}
              aria-label={label}
              aria-invalid={salah}
              onChange={(e) =>
                onUbah?.({ ...nilai, [kunci]: angkaIsian ? saringKetikanAngka(e.target.value) : e.target.value })
              }
            />
            {salah && (
              <span role="alert" className="os__galat-isian">
                {galatRnm}
              </span>
            )}
          </span>
        ) : (
          <span className={angkaIsian ? 'os__angka' : undefined}>
            {kunci === 'jumlah' ? angka(b.jumlahKlaim) : nilai[kunci] || '—'}
          </span>
        )}
      </label>
    )
  }

  return (
    <div className="rincian__adjustment">
      <div className="os__form os__form--tiga os__form--angka">
        <dl className="os__kolom">
          <Baris label={A.mataUang} nilai={typeof b.jumlahKlaim?.currency === 'string' ? b.jumlahKlaim.currency : ''} />
          {/* Teks rata kiri, angka rata kanan — keduanya dalam kolom nilai selebar tetap (os__form--angka). */}
          <Baris label={A.persen} nilai={turunan.persen} angka />
          <Baris label={A.dibayar} nilai={turunan.dibayar} angka />
        </dl>
        <dl className="os__kolom">
          <Baris label={A.shareRnm} nilai={angka(b.shareNusantaraRe)} angka />
          <Baris label={A.retensi} nilai={angka(b.cedingRetention)} angka />
          {teks(A.klaimRnm, 'jumlah', 'decimal')}
        </dl>
        <dl className="os__kolom">
          <Baris label={A.sumReasured} nilai={angka(b.sumReasured)} angka />
          <Baris label={A.sumInsured} nilai={angka(b.sumInsured)} angka />
        </dl>
      </div>
      <div className="os__form os__form--dua os__form--bank">
        <dl className="os__kolom">
          <Baris label={A.payableTo} nilai={b.payableTo} />
          {teks(A.namaBank, 'namaBank')}
        </dl>
        <dl className="os__kolom">
          {/* "semua read only kecuali name of bank" (11-10-2026): Payable To, Branch of Bank, Account No baca saja. */}
          <Baris label={A.cabangBank} nilai={b.cabangBank} />
          <Baris label={A.nomorRekening} nilai={b.nomorRekening} />
        </dl>
      </div>
      {/* Spreading TERLIPAT bawaan (11-10-2026, "seperti adjustment saat diklik baru muncul"); dimuat saat pertama dibuka. */}
      <section className="os__spreading" aria-label={A.spreading}>
        <button
          type="button"
          className="btn btn--ghost btn--sm os__lipat"
          aria-expanded={lihatSpreading}
          onClick={() => {
            if (!lihatSpreading) onMuatSpreading?.()
            setLihatSpreading(!lihatSpreading)
          }}
        >
          {lihatSpreading ? 'Hide Spreading ▾' : 'View Spreading ▸'}
        </button>
        {lihatSpreading && <BagianSpreading s={spreading} galat={galatSpreading} />}
      </section>
      {aksi !== undefined && <p className="os__aksi-baris">{aksi}</p>}
    </div>
  )
}
