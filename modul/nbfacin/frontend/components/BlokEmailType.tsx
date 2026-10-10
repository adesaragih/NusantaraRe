// Blok Email Type tahap marketing + Submit (tiket 57) - port `NB FacIn\Section\EmailSection.xml` (lewat
// `InputInwardFacultative` -> `ComfirmPolis_NonLife`, tahap Offer) dengan gambar layar work owner 06-10-2026:
// - Email Type* (Accept / Reject / Banding / Decline); memilih Accept menjalankan proteksi (`Protection_Act`, server) dan
//   pesannya menahan Submit; Banding memunculkan "Banding to"* (`GetAksepBanding_SQL`); Decline memunculkan "Choose Cause
//   of decline" dan Comment terkunci (diisi popup);
// - Comment* (`.OfferFacIn.Comment`);
// - grid riwayat Date · PIC · Approval · Status · Comment (terbaru dulu);
// - Submit di kaki halaman (hook `useSubmitMarketing`); Reject / Decline meminta konfirmasi `RejectNotificationSection`.
// Alur sesudah Submit di backend: Accept -> Team Leader, Reject -> tetap di marketing, Banding -> jabatan Banding to,
// Decline -> selesai; riwayat ke tabel lama (keputusan work owner 07-10-2026 "setuju" atas tiket 57).
//
// Tiket 61 - tahap Underwriting (`tahap === 'uw'`, `Section\EmailSection_IsUW.xml`): opsi `.EmailTypeUW` (Accept / Reject /
// Ask / Decline), Comment wajib, Validity Of Acceptance (hari) + "days, or up to" + Waiting Bind (dihitung server,
// `GetDateValidity_ACT`). Decline TANPA Cause of decline (tidak ada di EmailSection_IsUW) dan Comment tetap diisi sendiri;
// konfirmasi hanya Decline (tombol Reject UW langsung `SetRejectProposal`, tanpa modal).

import { useCallback, useEffect, useRef, useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ApiFailure } from '../../../../inti/frontend/klien'
import {
  ajukanSubmit,
  ambilSubmit,
  hitungValiditas,
  periksaSubmit,
  tekanBinding,
  type HasilSubmit,
  type KodeEmailType,
  type PilihanSebabTc,
  type TampilanSubmit,
} from '../api'
import {
  EMAIL_TYPE as E,
  EMAIL_TYPE_BINDING as EB,
  EMAIL_TYPE_TL as ETL,
  EMAIL_TYPE_UW as EUW,
  EMAIL_TYPE_UW_POLICY as EUWP,
  TEKS_EMAIL_TYPE as T,
} from '../labels'
import HalamanPopup from './HalamanPopup'
import PopupSebabDecline from './PopupSebabDecline'

/** Baris riwayat akseptasi per halaman (work owner 08-10-2026: "ini dibuatin page dong, 6 list perhalaman"). */
export const UKURAN_RIWAYAT = 6

/** Potongan riwayat untuk halaman `hal` (1-based; dijepit ke halaman terakhir yang ada). */
export function potongRiwayat<T>(baris: readonly T[], hal: number): { baris: T[]; kini: number } {
  const totalHal = Math.max(1, Math.ceil(baris.length / UKURAN_RIWAYAT))
  const kini = Math.min(Math.max(1, hal), totalHal)
  return { baris: baris.slice((kini - 1) * UKURAN_RIWAYAT, kini * UKURAN_RIWAYAT), kini }
}

export interface IsianSubmit {
  emailType: '' | KodeEmailType
  comment: string
  bandingTo: string
  sebabTc: PilihanSebabTc[]
  sebabNonTc: string
  /** Team Leader (tiket 58). */
  deductibleAcceptance: boolean
  commentDeductible: string
  /** Underwriting (tiket 61): Validity Of Acceptance (hari). */
  daysValidity: string
  /** Binding (tiket 64): centang Received RI Slip. */
  receivedRiSlip: boolean
}

export const isianSubmitKosong = (): IsianSubmit => ({
  emailType: '',
  comment: '',
  bandingTo: '',
  sebabTc: [],
  sebabNonTc: '',
  deductibleAcceptance: false,
  commentDeductible: '',
  daysValidity: '',
  receivedRiSlip: false,
})

/** Team Leader tahap Offer (tiket 58) atau Policy (tiket 69): layar & isian sama (`EmailSection.xml`, `.EmailTypeTL`). */
export const tahapTL = (tahap: TampilanSubmit['tahap'] | undefined) => tahap === 'teamLeader' || tahap === 'teamLeaderPolicy'
/** Underwriting tahap Offer (tiket 61) atau Policy (tiket 69): `EmailSection_IsUW.xml` - Decline tanpa Cause of decline. */
export const tahapUW = (tahap: TampilanSubmit['tahap'] | undefined) => tahap === 'uw' || tahap === 'uwPolicy'

/** Opsi Email Type menurut tahap case: marketing (57), Team Leader (58 / 69), Underwriting (61), UW Policy (69), Binding (64). */
export const opsiEmailType = (tahap: TampilanSubmit['tahap'] | undefined) =>
  tahapTL(tahap)
    ? ETL.opsi
    : tahap === 'uw'
      ? EUW.opsi
      : tahap === 'uwPolicy'
        ? EUWP.opsi
        : tahap === 'binding'
          ? EB.opsi
          : E.opsi

/** Pilihan yang melewati proteksi server (Accept, Banding) - Protection_Act langkah 8. */
export const diproteksi = (kode: string) => kode === '1' || kode === '4'

/** Validity Of Acceptance: kosong atau angka bulat hari. */
export const hariSah = (v: string) => /^\d{0,4}$/.test(v.trim())


/** Kabel `DD-MM-YYYY HH:MM` -> `DD/MM/YYYY HH:MM` (format grid Pega `dd/MM/yyyy`). */
export const tanggalRiwayat = (v: string) => v.replace(/^(\d{2})-(\d{2})-(\d{4})/, '$1/$2/$3')

/** Pesan kelengkapan isian sebelum Submit; '' = lengkap. Tahap UW: Decline tanpa Cause of decline (tiket 61). */
export function galatIsianSubmit(
  isi: IsianSubmit,
  tahap: TampilanSubmit['tahap'] | undefined = undefined,
  confirmBinding = false,
): string {
  if (isi.emailType === '') return T.wajibType
  // Binding: Accept butuh tombol Binding (Pega: Send Email Accept butuh ConfirmBinding = 1).
  if (tahap === 'binding' && isi.emailType === '1' && !confirmBinding) return T.wajibBinding
  if (isi.emailType === '4' && isi.bandingTo === '') return T.wajibBanding
  if (!tahapUW(tahap) && isi.emailType === '7' && isi.sebabNonTc === '' && !isi.sebabTc.some((s) => s.prioritas === 'I')) return T.wajibSebab
  if (isi.comment.trim() === '') return T.wajibComment
  if (tahapUW(tahap) && !hariSah(isi.daysValidity)) return T.hariAngka
  return ''
}

/**
 * Pesan proteksi dari penolakan 422 backend (`{"galat": pesan digabung "; ", "pesan": [...]}`): klien inti hanya menyimpan
 * `galat`, jadi daftarnya dipulihkan dengan memecah "; ". null = bukan penolakan proteksi.
 */
export function pesanProteksiDari(err: unknown): string[] | null {
  if (!(err instanceof ApiFailure) || err.status !== 422) return null
  return err.message
    .split('; ')
    .map((p) => p.trim())
    .filter((p) => p !== '')
}

/** Status dan aksi Submit tahap marketing - dipakai bersama blok Email Type dan tombol Submit kaki halaman. */
export function useSubmitMarketing(caseId: string, saatBerhasil: (h: HasilSubmit) => void) {
  const [v, setV] = useState<TampilanSubmit | null>(null)
  const [isi, setIsi] = useState<IsianSubmit>(isianSubmitKosong)
  const [proteksi, setProteksi] = useState<string[]>([])
  // Tahap Binding: pesan proteksi yang menahan tombol Binding (subset `proteksi`, ditentukan server menurut flag Pega).
  const [proteksiBinding, setProteksiBinding] = useState<string[]>([])
  const [memeriksa, setMemeriksa] = useState(false)
  const [mengajukan, setMengajukan] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [berhasil, setBerhasil] = useState<HasilSubmit | null>(null)
  // Tahap UW (tiket 61): "days, or up to" / Waiting Bind terakhir dari server.
  const [validitas, setValiditas] = useState({ dateValidity: '', waitingBindDate: '' })
  const nomorValiditas = useRef(0)
  const [membinding, setMembinding] = useState(false)

  const muat = useCallback(() => {
    ambilSubmit(caseId).then(
      (h) => {
        setV(h)
        // Isian TL tersimpan (Pega menyimpan centang Deductible Acceptance saat diubah); Validity UW dari server.
        // Validity UW: nilai TERSIMPAN saja ('' bila belum); default dihitung saat Email Type diklik (Pega GetDateValidity_ACT).
        setIsi((x) => ({
          ...x,
          deductibleAcceptance: h.deductibleAcceptance,
          commentDeductible: h.commentDeductible,
          daysValidity: h.daysValidity,
          receivedRiSlip: h.receivedRiSlip,
        }))
        setValiditas({ dateValidity: h.dateValidity, waitingBindDate: h.waitingBindDate })
      },
      (err: unknown) => setGalat(err),
    )
  }, [caseId])

  /**
   * Proteksi di server (Pega `Protection_Act`) untuk Accept DAN Banding: langkah 8 keluar hanya untuk Reject / Decline / Ask /
   * Revise (EmailType 2/7/3, EmailTypeBinding 9/2/3/7), jadi Banding ikut diproteksi; tombol Send Email Banding EmailSection
   * nonaktif bila ProtectSpreading / ProtectPayment / ProtectMarine / ProtectNeRate (work owner 08-10-2026: "iya disamakan").
   */
  async function periksa(baru: IsianSubmit) {
    if (!diproteksi(baru.emailType)) return
    setMemeriksa(true)
    setGalat(null)
    try {
      const h = await periksaSubmit(caseId, {
        emailType: baru.emailType as KodeEmailType,
        deductibleAcceptance: baru.deductibleAcceptance,
        commentDeductible: baru.commentDeductible,
        daysValidity: baru.daysValidity.trim(),
      })
      setProteksi(h.pesan)
      setProteksiBinding(h.pesanBinding ?? [])
    } catch (err) {
      setGalat(err)
    } finally {
      setMemeriksa(false)
    }
  }

  useEffect(() => {
    muat()
  }, [muat])

  async function pilihType(kode: KodeEmailType) {
    const baru: IsianSubmit = {
      ...isi,
      emailType: kode,
      bandingTo: '',
      sebabTc: kode === '7' ? isi.sebabTc : [],
      sebabNonTc: kode === '7' ? isi.sebabNonTc : '',
    }
    setIsi(baru)
    setProteksi([])
    setBerhasil(null)
    // Tahap UW: klik Email Type menjalankan GetDateValidity_ACT mode default (EmailSection_IsUW) - hari + dua tanggal diisi.
    if (tahapUW(v?.tahap)) void muatValiditas('')
    await periksa(baru)
  }

  /** Isian Deductible TL berubah: simpan di layar; Accept diperiksa ulang bila `periksaLagi` (centang / lepas fokus). */
  function ubahDeductible(isiTl: Partial<Pick<IsianSubmit, 'deductibleAcceptance' | 'commentDeductible'>>, periksaLagi: boolean) {
    const baru: IsianSubmit = { ...isi, ...isiTl }
    if (baru.deductibleAcceptance === false) baru.commentDeductible = ''
    setIsi(baru)
    if (periksaLagi) void periksa(baru)
  }

  /** Hitung Validity di server: hari '' = mode default (klik Email Type), angka = mode edit. Jawaban lama dibuang. */
  async function muatValiditas(hari: string) {
    const n = ++nomorValiditas.current
    try {
      const h = await hitungValiditas(caseId, hari)
      if (n !== nomorValiditas.current) return
      setValiditas({ dateValidity: h.dateValidity, waitingBindDate: h.waitingBindDate })
      if (hari === '') setIsi((x) => ({ ...x, daysValidity: h.daysValidity }))
    } catch (err) {
      if (n === nomorValiditas.current) setGalat(err)
    }
  }

  /** Validity Of Acceptance UW diubah (lepas fokus): mode edit. */
  async function ubahHari() {
    if (!hariSah(isi.daysValidity) || isi.daysValidity.trim() === '') return
    await muatValiditas(isi.daysValidity.trim())
  }

  /** Tombol "Binding" (tiket 64): ditulis server SEBELUM Submit (seperti Pega); tampilan Submit dimuat ulang dari jawabannya. */
  async function binding() {
    setMembinding(true)
    setGalat(null)
    try {
      const h = await tekanBinding(caseId)
      setV(h)
      setIsi((x) => ({ ...x, receivedRiSlip: h.receivedRiSlip }))
    } catch (err) {
      setGalat(err)
    } finally {
      setMembinding(false)
    }
  }

  async function kirim() {
    if (isi.emailType === '') return
    setKonfirmasi(false)
    setMengajukan(true)
    setGalat(null)
    try {
      const h = await ajukanSubmit(caseId, {
        emailType: isi.emailType,
        comment: isi.comment,
        bandingTo: isi.bandingTo,
        sebabTc: isi.sebabTc.filter((s) => s.prioritas !== ''),
        sebabNonTc: isi.sebabNonTc,
        deductibleAcceptance: isi.deductibleAcceptance,
        commentDeductible: isi.commentDeductible,
        daysValidity: tahapUW(v?.tahap) ? isi.daysValidity.trim() : '',
        receivedRiSlip: v?.tahap === 'binding' ? isi.receivedRiSlip : false,
      })
      if (h.pesan.length > 0) {
        setProteksi(h.pesan)
        return
      }
      setBerhasil(h)
      setIsi(isianSubmitKosong())
      setProteksi([])
      muat()
      saatBerhasil(h)
    } catch (err) {
      const p = pesanProteksiDari(err)
      if (p) setProteksi(p)
      else setGalat(err)
    } finally {
      setMengajukan(false)
    }
  }

  const galatIsi = galatIsianSubmit(isi, v?.tahap, v?.confirmBinding ?? false)
  const bisaSubmit =
    v !== null && v.bolehSubmit && galatIsi === '' && !memeriksa && !mengajukan && (!diproteksi(isi.emailType) || proteksi.length === 0)

  /** Tombol Submit kaki halaman: Reject / Decline lewat konfirmasi dulu (tahap UW: Decline saja). */
  const submit = () => {
    if (!bisaSubmit) return
    if (isi.emailType === '7' || isi.emailType === '9' || (isi.emailType === '2' && !tahapUW(v?.tahap))) setKonfirmasi(true)
    else void kirim()
  }

  return {
    v,
    isi,
    setIsi,
    proteksi,
    proteksiBinding,
    memeriksa,
    mengajukan,
    galat,
    konfirmasi,
    setKonfirmasi,
    berhasil,
    galatIsi,
    bisaSubmit,
    pilihType,
    ubahDeductible,
    validitas,
    ubahHari,
    membinding,
    binding,
    periksaUlang: () => void periksa(isi),
    kirim,
    submit,
  }
}

export type SubmitMarketing = ReturnType<typeof useSubmitMarketing>

/** `hanyaRiwayat`: mode hanya lihat (View Details / View Policy) - form Email Type disembunyikan, riwayat akseptasi tetap. */
export default function BlokEmailType({ s, hanyaRiwayat = false }: { s: SubmitMarketing; hanyaRiwayat?: boolean }) {
  const [popup, setPopup] = useState(false)
  const [halRiwayat, setHalRiwayat] = useState(1)
  const { v, isi, setIsi } = s
  if (v === null) return s.galat ? <Gagal galat={s.galat} /> : null
  const bisaUbah = v.bolehSubmit && !s.mengajukan

  return (
    <section className="panel nbf-et">
      <Gagal galat={s.galat} />
      {s.berhasil && (
        <div className="alert alert--ok">
          {T.terkirim} {s.berhasil.nbStatus}
        </div>
      )}
      {/* Case di luar tahap Submit: form disembunyikan TANPA keterangan (work owner 08-10-2026: "PESAN INI TIDAK USAH
          DIMUNCULKAN"); riwayat akseptasi tetap tampil. */}
      {v.bolehSubmit && !hanyaRiwayat && (
        <>
      <div className="nbf-et__baris">
        <span className="nbf-et__label">
          {E.label}
          <span className="field__req">*</span>
        </span>
        <div className="nbf-et__opsi" role="radiogroup" aria-label={E.label}>
          {opsiEmailType(v.tahap).map((o) => (
            <label key={o.kode}>
              <input
                type="radio"
                name="nbfacin-email-type"
                checked={isi.emailType === o.kode}
                disabled={!bisaUbah}
                onChange={(e) => {
                  // Fokus dilepas sesudah memilih: tombol panah (mis. menggulir ke Comment) tidak lagi memindahkan pilihan
                  // radio (work owner 08-10-2026: "pilihannya mau keganti sendiri ... awal dipilih Accept ... jadi Ask").
                  e.currentTarget.blur()
                  void s.pilihType(o.kode)
                }}
              />{' '}
              {o.label}
            </label>
          ))}
        </div>
      </div>

      {isi.emailType === '4' && (
        <div className="nbf-et__baris">
          <label className="nbf-et__label" htmlFor="nbf-et-banding">
            {E.bandingTo}
            <span className="field__req">*</span>
          </label>
          {v.opsiBanding.length === 0 ? (
            <span className="nbf-et__catatan">{T.bandingTanpaOpsi}</span>
          ) : (
            <select id="nbf-et-banding" value={isi.bandingTo} disabled={!bisaUbah} onChange={(e) => setIsi({ ...isi, bandingTo: e.target.value })}>
              <option value="">{E.bandingKosong}</option>
              {v.opsiBanding.map((b) => (
                <option key={b} value={b}>
                  {b}
                </option>
              ))}
            </select>
          )}
        </div>
      )}

      {/* Tahap Binding (tiket 64, EmailSectionCeding): tombol Binding tampil saat Accept / sudah di-bind; Received RI Slip baru
          dapat dicentang sesudah Binding (disabled bila ConfirmBinding != 1). */}
      {v.tahap === 'binding' && (isi.emailType === '1' || v.confirmBinding) && (
        <div className="nbf-et__binding">
          <button
            type="button"
            className="btn nbf-tema__tombol btn--sm"
            // Pega EmailSectionCeding: nonaktif bila ConfirmBinding = 1 atau ProtectSpreading / ProtectDeductible - BUKAN semua
            // pesan proteksi Accept (work owner 08-10-2026: "tombol binding nya kenapa blm bisa ditekan").
            disabled={!bisaUbah || v.confirmBinding || s.membinding || s.memeriksa || s.proteksiBinding.length > 0}
            onClick={() => void s.binding()}
          >
            {s.membinding ? T.membinding : EB.tombolBinding}
          </button>
          <label className="nbf-et__centang">
            <input
              type="checkbox"
              checked={isi.receivedRiSlip}
              disabled={!bisaUbah || !v.confirmBinding}
              onChange={(e) => setIsi({ ...isi, receivedRiSlip: e.target.checked })}
            />{' '}
            {EB.receivedRiSlip}
          </label>
        </div>
      )}
      {/* Penahan tombol Binding yang tidak ada di daftar penahan Accept (mis. deductible, ProtectDeductible.CARI2 / CARI4). */}
      {v.tahap === 'binding' && s.proteksiBinding.some((p) => !s.proteksi.includes(p)) && (
        <ul className="nbf-et__proteksi">
          {s.proteksiBinding
            .filter((p) => !s.proteksi.includes(p))
            .map((p) => (
              <li key={p}>{p}</li>
            ))}
        </ul>
      )}

      {isi.emailType === '7' && !tahapUW(v.tahap) && (
        <div className="nbf-et__baris">
          <span className="nbf-et__label" />
          <button type="button" className="btn nbf-tema__tombol btn--sm" disabled={!bisaUbah} onClick={() => setPopup(true)}>
            {E.pilihSebab}
          </button>
        </div>
      )}

      {tahapTL(v.tahap) && v.deductibleBedaMataUang && (
        <div className="nbf-et__deductible">
          <label className="nbf-et__centang">
            <input
              type="checkbox"
              checked={isi.deductibleAcceptance}
              disabled={!bisaUbah}
              onChange={(e) => s.ubahDeductible({ deductibleAcceptance: e.target.checked }, true)}
            />{' '}
            {ETL.deductibleAcceptance}
          </label>
          {isi.deductibleAcceptance && (
            <div className="nbf-et__komentar">
              <label htmlFor="nbf-et-comment-deductible">
                {ETL.commentDeductible}
                <span className="field__req">*</span>
              </label>
              <textarea
                id="nbf-et-comment-deductible"
                rows={3}
                value={isi.commentDeductible}
                disabled={!bisaUbah}
                onChange={(e) => s.ubahDeductible({ commentDeductible: e.target.value }, false)}
                onBlur={() => s.periksaUlang()}
              />
            </div>
          )}
        </div>
      )}

      <div className="nbf-et__komentar">
        <label htmlFor="nbf-et-comment">
          {E.comment}:<span className="field__req">*</span>
        </label>
        <textarea
          id="nbf-et-comment"
          rows={5}
          value={isi.comment}
          disabled={!bisaUbah || (isi.emailType === '7' && !tahapUW(v.tahap))}
          onChange={(e) => setIsi({ ...isi, comment: e.target.value })}
        />
      </div>

      {tahapUW(v.tahap) && (
        <div className="nbf-et__validitas">
          <label htmlFor="nbf-et-hari">{EUW.validitas}</label>
          <input
            id="nbf-et-hari"
            inputMode="numeric"
            value={isi.daysValidity}
            disabled={!bisaUbah}
            onChange={(e) => setIsi({ ...isi, daysValidity: e.target.value })}
            onBlur={() => void s.ubahHari()}
          />
          <span>{EUW.hariAtau}</span>
          <output aria-label={EUW.hariAtau}>{tanggalRiwayat(s.validitas.dateValidity)}</output>
          <span>{EUW.waitingBind}</span>
          <output aria-label={EUW.waitingBind}>{tanggalRiwayat(s.validitas.waitingBindDate)}</output>
        </div>
      )}

      {s.memeriksa && <p className="muted">{T.memeriksa}</p>}
      {s.proteksi.length > 0 && (
        <ul className="nbf-et__proteksi">
          {s.proteksi.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      )}
        </>
      )}

      {/* Pager di ATAS tabel riwayat (work owner 08-10-2026: "utk page nya taro diatasnya saja bukan dibawah"). */}
      {v.riwayat.length > UKURAN_RIWAYAT && (
        <HalamanPopup
          halaman={potongRiwayat(v.riwayat, halRiwayat).kini}
          ukuran={UKURAN_RIWAYAT}
          total={v.riwayat.length}
          onPindah={setHalRiwayat}
        />
      )}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-et__riwayat">
          <thead>
            <tr>
              {E.kolomRiwayat.map((k, i) => (
                <th key={`${k}-${i}`} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {v.riwayat.length === 0 && (
              <tr>
                <td colSpan={5}>{T.riwayatKosong}</td>
              </tr>
            )}
            {potongRiwayat(v.riwayat, halRiwayat).baris.map((r, i) => (
              <tr key={i}>
                <td>{tanggalRiwayat(r.tanggal)}</td>
                <td>{r.pic}</td>
                <td>{r.approval}</td>
                <td>{r.status}</td>
                <td className="nbf-et__teks">{r.comment}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {popup && (
        <PopupSebabDecline
          tc={v.sebabTc}
          nonTc={v.sebabNonTc}
          pctOtomatis={v.pctRnmOtomatis}
          awal={isi.sebabTc.length > 0 || isi.sebabNonTc !== '' ? { pilihan: isi.sebabTc, nonTc: isi.sebabNonTc } : null}
          onTutup={() => setPopup(false)}
          onSubmit={(h) => {
            setIsi({ ...isi, sebabTc: h.pilihan, sebabNonTc: h.nonTc, comment: h.komentar })
            setPopup(false)
          }}
        />
      )}

      {s.konfirmasi && (
        <Modal
          judul={E.label}
          onTutup={() => s.setKonfirmasi(false)}
          labelBatal={T.batal}
          aksi={
            <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => void s.kirim()}>
              {T.ya}
            </button>
          }
        >
          <p>{isi.emailType === '2' ? E.yakinReject : isi.emailType === '9' ? EB.yakinRevise : E.yakinDecline}</p>
        </Modal>
      )}
    </section>
  )
}
