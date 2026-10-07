// Tab Correspondence (tiket 56) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` kontainer S95 ->
// `Correspondence` (grid `.OfferFacIn.CorrespondenceList`) + editor baris `CorrespondenceContent`, tata letak = gambar
// layar Pega dari work owner 06-10-2026:
// - grid Date · Sender · Subject · Recipient · Status, Add (draf baru: tanggal + pengirim dari akun login, server),
//   Remove (hanya "Not Sent");
// - ▸ membuka editor: Sender, Sender Mail Address (baca-saja); Recipients: (saran akun ber-email, atau diketik bebas) +
//   Add / Remove; Subject; Message (teks berformat); Save dan Send Email (hanya "Not Sent");
// - sesudah terkirim: status "Sent", seluruh isi baca-saja.
// Lampiran tidak dibuat (Pega menyembunyikannya, `pyContainerVisibleWhen` = `1=2`).
//
// Keputusan work owner 06-10-2026 ("iya kirim email sungguhan ... setuju dengan rekomendasi kamu"): email dikirim
// sungguhan lewat SMTP konfigurasi server; To = penerima (Pega: To = pengirim sendiri, penerima di CC - dibetulkan);
// kolom Recipient grid = nama penerima; daftar penerima dari akun sistem baru; tabel baru.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilKorespondensi,
  hapusKorespondensi,
  kirimKorespondensi,
  saranPenerimaKorespondensi,
  simpanKorespondensi,
  tambahKorespondensi,
  type Korespondensi,
  type PenerimaKorespondensi,
  type TampilanKorespondensi,
} from '../api'
import { KORESPONDENSI as K, TEKS_KORESPONDENSI as T } from '../labels'
import EditorTeksKaya from './EditorTeksKaya'
import { BATAS_BITA_PESAN, bitaPesan } from './htmlKorespondensi'
import IsianSaran from './IsianSaran'

/** Kabel `DD-MM-YYYY HH:MM` -> tampil `DD/MM/YY HH:MM` (gambar layar Pega "06/10/26 17:28"). */
export function tanggalKorespondensi(kabel: string): string {
  const m = /^(\d{2})-(\d{2})-\d{2}(\d{2}) (\d{2}:\d{2})$/.exec(kabel.trim())
  return m ? `${m[1]}/${m[2]}/${m[3]} ${m[4]}` : kabel
}

/** Kolom Recipient grid: nama penerima (atau alamatnya bila tanpa nama), dipisah koma. */
export const daftarNamaPenerima = (p: PenerimaKorespondensi[]) =>
  p
    .map((x) => x.recipient.trim() || x.recipientMail.trim())
    .filter((x) => x !== '')
    .join(', ')

const POLA_EMAIL = /^[^\s@<>(),;:"]+@[^\s@<>(),;:"]+\.[^\s@<>(),;:"]+$/

/** Pesan galat sebelum Send Email: minimal satu penerima ber-email, dan semua alamat yang diisi sah. "" = boleh. */
export function galatPenerima(p: PenerimaKorespondensi[]): string {
  const diisi = p.map((x) => x.recipientMail.trim()).filter((x) => x !== '')
  if (diisi.length === 0) return T.penerimaKosong
  const salah = diisi.filter((x) => !POLA_EMAIL.test(x))
  return salah.length > 0 ? T.emailSalah + salah.join(', ') : ''
}

const isiDari = (k: Korespondensi) => ({ subject: k.subject, message: k.message, penerima: k.penerima })

export default function TabCorrespondence({ caseId, awal = null }: { caseId: string; awal?: TampilanKorespondensi | null }) {
  const [v, setV] = useState<TampilanKorespondensi | null>(awal)
  const [galat, setGalat] = useState<unknown>(null)
  const [terbuka, setTerbuka] = useState<string[]>([])
  const [sibuk, setSibuk] = useState<{ no: string; aksi: 'simpan' | 'kirim' | 'hapus' | 'tambah' } | null>(null)
  const [info, setInfo] = useState<{ no: string; teks: string } | null>(null)
  const [galatKirim, setGalatKirim] = useState<{ no: string; teks: string } | null>(null)

  useEffect(() => {
    let batal = false
    ambilKorespondensi(caseId).then(
      (h) => {
        if (!batal) setV(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId])

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  const ganti = (k: Korespondensi) => setV({ ...v, daftar: v.daftar.map((x) => (x.no === k.no ? k : x)) })
  const ubah = (no: string, isi: Partial<Korespondensi>) => {
    setInfo(null)
    setV({ ...v, daftar: v.daftar.map((x) => (x.no === no ? { ...x, ...isi } : x)) })
  }

  async function tambah() {
    if (v === null) return
    setSibuk({ no: '', aksi: 'tambah' })
    setGalat(null)
    try {
      const k = await tambahKorespondensi(caseId)
      setV({ ...v, daftar: [...v.daftar, k] })
      setTerbuka((t) => [...t, k.no])
    } catch (err) {
      setGalat(err)
    } finally {
      setSibuk(null)
    }
  }

  async function hapus(k: Korespondensi) {
    if (v === null) return
    setSibuk({ no: k.no, aksi: 'hapus' })
    setGalat(null)
    try {
      await hapusKorespondensi(caseId, k.no)
      setV({ ...v, daftar: v.daftar.filter((x) => x.no !== k.no) })
      setTerbuka((t) => t.filter((x) => x !== k.no))
    } catch (err) {
      setGalat(err)
    } finally {
      setSibuk(null)
    }
  }

  async function simpan(k: Korespondensi) {
    setSibuk({ no: k.no, aksi: 'simpan' })
    setGalat(null)
    setGalatKirim(null)
    try {
      ganti(await simpanKorespondensi(caseId, k.no, isiDari(k)))
      setInfo({ no: k.no, teks: T.tersimpan })
    } catch (err) {
      setGalat(err)
    } finally {
      setSibuk(null)
    }
  }

  async function kirim(k: Korespondensi) {
    const g = galatPenerima(k.penerima)
    if (g !== '') {
      setGalatKirim({ no: k.no, teks: g })
      return
    }
    setSibuk({ no: k.no, aksi: 'kirim' })
    setGalat(null)
    setGalatKirim(null)
    try {
      ganti(await kirimKorespondensi(caseId, k.no, isiDari(k)))
      setInfo({ no: k.no, teks: T.terkirimInfo })
    } catch (err) {
      setGalat(err)
    } finally {
      setSibuk(null)
    }
  }

  const ubahPenerima = (k: Korespondensi, i: number, p: Partial<PenerimaKorespondensi>) =>
    ubah(k.no, { penerima: k.penerima.map((x, j) => (j === i ? { ...x, ...p } : x)) })

  const editor = (k: Korespondensi) => {
    const kunci = k.terkirim
    const sedang = sibuk?.no === k.no ? sibuk.aksi : null
    const bita = bitaPesan(k.message)
    const kepanjangan = bita > BATAS_BITA_PESAN
    return (
      <div className="nbf-kor__editor">
        {kunci && <div className="alert">{T.terkunci}</div>}
        {info?.no === k.no && <div className="alert alert--ok">{info.teks}</div>}
        <dl className="nbf-kor__pengirim">
          <div>
            <dt>{K.sender}</dt>
            <dd>{k.sender || '—'}</dd>
          </div>
          <div>
            <dt>{K.senderMail}</dt>
            <dd>{k.senderMail || '—'}</dd>
          </div>
        </dl>

        <section className="nbf-kor__penerima">
          <div className="nbf-kor__penerima-kepala">
            <h6>{K.penerima}</h6>
            {!kunci && (
              <button
                type="button"
                className="btn btn--ghost btn--sm"
                onClick={() => ubah(k.no, { penerima: [...k.penerima, { recipient: '', recipientMail: '' }] })}
              >
                ＋ {K.tambah}
              </button>
            )}
          </div>
          {k.penerima.length === 0 && <p className="nbf-kor__kosong">{T.kosong}</p>}
          {k.penerima.map((p, i) =>
            kunci ? (
              <div key={i} className="nbf-kor__penerima-baris nbf-kor__penerima-baris--baca">
                <span>{p.recipient || '—'}</span>
                <span>{p.recipientMail}</span>
              </div>
            ) : (
              <div key={i} className="nbf-kor__penerima-baris">
                <IsianSaran
                  label={T.namaPenerima}
                  value={p.recipient}
                  onKetik={(x) => ubahPenerima(k, i, { recipient: x })}
                  cari={saranPenerimaKorespondensi}
                  onPilih={(s) => ubahPenerima(k, i, { recipient: s.label, recipientMail: s.ekstra ?? '' })}
                  tampil={(s) => (s.ekstra ? `${s.label} — ${s.ekstra}` : s.label)}
                  muatSaatFokus
                />
                <div>
                  <label htmlFor={`kor-${k.no}-mail-${i}`}>{T.emailPenerima}</label>
                  <input
                    id={`kor-${k.no}-mail-${i}`}
                    type="email"
                    value={p.recipientMail}
                    onChange={(e) => ubahPenerima(k, i, { recipientMail: e.target.value.trim() })}
                  />
                </div>
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => ubah(k.no, { penerima: k.penerima.filter((_, j) => j !== i) })}
                >
                  ✕ {K.hapus}
                </button>
              </div>
            ),
          )}
        </section>

        <div className="nbf-kor__subjek">
          <label htmlFor={`kor-${k.no}-subjek`}>{K.subject}</label>
          <input id={`kor-${k.no}-subjek`} value={k.subject} disabled={kunci} onChange={(e) => ubah(k.no, { subject: e.target.value })} />
        </div>

        <EditorTeksKaya key={`${k.no}-${kunci ? 'baca' : 'ubah'}`} label={K.message} awal={k.message} bacaSaja={kunci} onChange={(html) => ubah(k.no, { message: html })} />
        {!kunci && (
          <p className={'nbf-kor__ukuran' + (kepanjangan ? ' nbf-kor__ukuran--lebih' : '')}>
            {bita}
            {T.ukuran}
            {kepanjangan && ` - ${T.terlaluPanjang}`}
          </p>
        )}

        {galatKirim?.no === k.no && <div className="alert alert--warn">{galatKirim.teks}</div>}
        {!kunci && (
          <div className="nbf-kor__tombol">
            <button type="button" className="btn btn--primary" disabled={sibuk !== null || kepanjangan} onClick={() => void simpan(k)}>
              {sedang === 'simpan' ? T.menyimpan : K.simpan}
            </button>
            <button type="button" className="btn" disabled={sibuk !== null || kepanjangan} onClick={() => void kirim(k)}>
              {sedang === 'kirim' ? T.mengirim : K.kirim}
            </button>
          </div>
        )}
      </div>
    )
  }

  return (
    <div className="nbf-objek nbf-kor">
      <Gagal galat={galat} />
      {v.pesan.map((p) => (
        <div key={p} className="alert alert--warn">
          {p}
        </div>
      ))}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-kor__grid">
          <thead>
            <tr>
              <th scope="col" />
              {K.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col">
                <button type="button" className="btn btn--ghost btn--sm" disabled={sibuk !== null} onClick={() => void tambah()}>
                  ＋ {K.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {v.daftar.length === 0 && (
              <tr>
                <td colSpan={7}>{T.kosong}</td>
              </tr>
            )}
            {v.daftar.map((k) => {
              const buka = terbuka.includes(k.no)
              return [
                <tr key={k.no} className={buka ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <button
                      type="button"
                      className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
                      aria-label={T.bukaBaris}
                      aria-expanded={buka}
                      onClick={() => setTerbuka((t) => (buka ? t.filter((x) => x !== k.no) : [...t, k.no]))}
                    >
                      ▸
                    </button>
                  </td>
                  <td>{tanggalKorespondensi(k.tanggal)}</td>
                  <td>{k.sender}</td>
                  <td>{k.subject}</td>
                  <td>{daftarNamaPenerima(k.penerima)}</td>
                  <td>
                    <span className={'nbf-kor__status' + (k.terkirim ? ' nbf-kor__status--terkirim' : '')}>
                      {k.terkirim ? K.terkirim : K.belumTerkirim}
                    </span>
                  </td>
                  <td>
                    {!k.terkirim && (
                      <button type="button" className="btn btn--ghost btn--sm" disabled={sibuk !== null} onClick={() => void hapus(k)}>
                        ✕ {K.hapus}
                      </button>
                    )}
                  </td>
                </tr>,
                buka && (
                  <tr key={`${k.no}-d`} className="nbf-objek__detail">
                    <td colSpan={7}>{editor(k)}</td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
