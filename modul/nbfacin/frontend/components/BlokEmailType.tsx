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

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ApiFailure } from '../../../../inti/frontend/klien'
import {
  ajukanSubmit,
  ambilSubmit,
  periksaSubmit,
  type HasilSubmit,
  type KodeEmailType,
  type PilihanSebabTc,
  type TampilanSubmit,
} from '../api'
import { EMAIL_TYPE as E, TEKS_EMAIL_TYPE as T } from '../labels'
import PopupSebabDecline from './PopupSebabDecline'

export interface IsianSubmit {
  emailType: '' | KodeEmailType
  comment: string
  bandingTo: string
  sebabTc: PilihanSebabTc[]
  sebabNonTc: string
}

export const isianSubmitKosong = (): IsianSubmit => ({ emailType: '', comment: '', bandingTo: '', sebabTc: [], sebabNonTc: '' })

/** Kabel `DD-MM-YYYY HH:MM` -> `DD/MM/YYYY HH:MM` (format grid Pega `dd/MM/yyyy`). */
export const tanggalRiwayat = (v: string) => v.replace(/^(\d{2})-(\d{2})-(\d{4})/, '$1/$2/$3')

/** Pesan kelengkapan isian sebelum Submit; '' = lengkap. */
export function galatIsianSubmit(isi: IsianSubmit): string {
  if (isi.emailType === '') return T.wajibType
  if (isi.emailType === '4' && isi.bandingTo === '') return T.wajibBanding
  if (isi.emailType === '7' && isi.sebabNonTc === '' && !isi.sebabTc.some((s) => s.prioritas === 'I')) return T.wajibSebab
  if (isi.comment.trim() === '') return T.wajibComment
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
  const [memeriksa, setMemeriksa] = useState(false)
  const [mengajukan, setMengajukan] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [berhasil, setBerhasil] = useState<HasilSubmit | null>(null)

  const muat = useCallback(() => {
    ambilSubmit(caseId).then(setV, (err: unknown) => setGalat(err))
  }, [caseId])

  useEffect(() => {
    muat()
  }, [muat])

  async function pilihType(kode: KodeEmailType) {
    setIsi((x) => ({ ...x, emailType: kode, bandingTo: '', sebabTc: kode === '7' ? x.sebabTc : [], sebabNonTc: kode === '7' ? x.sebabNonTc : '' }))
    setProteksi([])
    setBerhasil(null)
    if (kode !== '1') return
    setMemeriksa(true)
    setGalat(null)
    try {
      setProteksi((await periksaSubmit(caseId, kode)).pesan)
    } catch (err) {
      setGalat(err)
    } finally {
      setMemeriksa(false)
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

  const galatIsi = galatIsianSubmit(isi)
  const bisaSubmit =
    v !== null && v.bolehSubmit && galatIsi === '' && !memeriksa && !mengajukan && (isi.emailType !== '1' || proteksi.length === 0)

  /** Tombol Submit kaki halaman: Reject / Decline lewat konfirmasi dulu. */
  const submit = () => {
    if (!bisaSubmit) return
    if (isi.emailType === '2' || isi.emailType === '7') setKonfirmasi(true)
    else void kirim()
  }

  return { v, isi, setIsi, proteksi, memeriksa, mengajukan, galat, konfirmasi, setKonfirmasi, berhasil, galatIsi, bisaSubmit, pilihType, kirim, submit }
}

export type SubmitMarketing = ReturnType<typeof useSubmitMarketing>

export default function BlokEmailType({ s }: { s: SubmitMarketing }) {
  const [popup, setPopup] = useState(false)
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
      {!v.bolehSubmit && (
        <div className="alert nbf-et__info">
          <span>{T.sudahDiproses}</span>
          {v.alasan !== '' && <small>{T.tidakBoleh + v.alasan}</small>}
        </div>
      )}

      {v.bolehSubmit && (
        <>
      <div className="nbf-et__baris">
        <span className="nbf-et__label">
          {E.label}
          <span className="field__req">*</span>
        </span>
        <div className="nbf-et__opsi" role="radiogroup" aria-label={E.label}>
          {E.opsi.map((o) => (
            <label key={o.kode}>
              <input
                type="radio"
                name="nbfacin-email-type"
                checked={isi.emailType === o.kode}
                disabled={!bisaUbah}
                onChange={() => void s.pilihType(o.kode)}
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

      {isi.emailType === '7' && (
        <div className="nbf-et__baris">
          <span className="nbf-et__label" />
          <button type="button" className="btn btn--sm" disabled={!bisaUbah} onClick={() => setPopup(true)}>
            {E.pilihSebab}
          </button>
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
          disabled={!bisaUbah || isi.emailType === '7'}
          onChange={(e) => setIsi({ ...isi, comment: e.target.value })}
        />
      </div>

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
            {v.riwayat.map((r, i) => (
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
            <button type="button" className="btn btn--primary" onClick={() => void s.kirim()}>
              {T.ya}
            </button>
          }
        >
          <p>{isi.emailType === '2' ? E.yakinReject : E.yakinDecline}</p>
        </Modal>
      )}
    </section>
  )
}
