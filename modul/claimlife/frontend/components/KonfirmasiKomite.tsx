// Pop-up Send Claim to Committee — `Section/ClaimComite.xml` (Harness `Committe_Life`), keputusan work owner 12-10-2026
// opsi A: "Remarks disimpan di jejak audit, dan Date memakai waktu server sehingga tidak bisa diubah".
//
//   Date        `.DataCommitteeTreaty.TanggalComitee` — di sini waktu server (BACA SAJA; yang tersimpan WAKTU jejak).
//   PIC         `OperatorID.pyUserName` (baca saja) — akun pelaku.
//   Adjustment  label (b2018) + keterangan baris: Claim RNM dan mata uangnya.
//   Remarks*    `.DataCommitteeTreaty.Remarks` (pxTextArea, wajib) → T_CLAIMLF_JEJAK.KOMENTAR, ≤ 4000 byte.
//   Grid        `.KomiteList`: Committee Name | Status | Date Approve | Comment — PRATINJAU roster EMAILKOMITE (status
//               awal "Waiting").
//   Tombol      "Send Claim to Committee" pyDisabledWhen `.KomiteList(1).IDKomite = '' || .DataCommitteeTreaty.Remarks=''
//               || pyWorkPage.ClaimData.IsKPR=''`; TANPA konfirmasi kedua (Pega: runActivity CreateKMTLife_Act). Cancel.

import { useEffect, useState } from 'react'

import { TOMBOL_KOMITE } from '../labels'
import { ambilPratinjauKomite, serahkanKeKomite, type Klaim, type PratinjauKomite } from '../api'
import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { formatNumber, DESIMAL_TAK_DIBATASI } from '../../../../inti/frontend/lib/format'
import { PitaGalat } from './PitaGalat'

/** Batas Remarks — T_CLAIMLF_JEJAK.KOMENTAR VARCHAR2(4000); server memeriksa byte. */
export const BATAS_REMARKS_KOMITE = 4000

/** Kalimat pop-up (English, label korpus VERBATIM ClaimComite). */
export const POPUP_KOMITE = {
  judul: TOMBOL_KOMITE.serahkan,
  tanggal: 'Date',
  pic: 'PIC',
  adjustment: 'Adjustment',
  remarks: 'Remarks',
  kolom: ['Committee Name', 'Status', 'Date Approve', 'Comment'],
  menunggu: 'Waiting',
  memuat: 'Loading the committee list…',
  tanpaKpr: 'Is KPR of this claim is empty — fill it before sending to the committee.',
  tanpaAnggota: 'No committee members found.',
  judulKomite: 'Committee',
  wajibRemarks: 'Remarks is required.',
  isiRemarks: 'Fill in Remarks to continue.',
  proses: 'Sending…',
} as const

/** Waktu server ISO → dd/mm/yyyy HH:mm (waktu lokal peramban). */
export function tanggalWaktuKomite(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  const dua = (n: number) => String(n).padStart(2, '0')
  return `${dua(d.getDate())}/${dua(d.getMonth() + 1)}/${d.getFullYear()} ${dua(d.getHours())}:${dua(d.getMinutes())}`
}

/** Penghitung karakter Remarks "n / 4000". */
export function hitungRemarks(remarks: string): string {
  return `${remarks.length} / ${BATAS_REMARKS_KOMITE}`
}

/** Aturan XML tombol kirim: anggota roster ada, Remarks terisi, Is KPR klaim terisi. */
export function bolehKirimKomite(pratinjau: PratinjauKomite | null, remarks: string, isKpr: string | undefined): boolean {
  return pratinjau !== null && pratinjau.anggota.length > 0 && remarks.trim() !== '' && (isKpr ?? '').trim() !== ''
}

export function KonfirmasiKomite({
  klaim,
  pesertaID,
  adjID,
  onTutup,
  onTerkirim,
  awal,
  jumlahKlaim = '',
  mataUang = '',
  galatAwal,
}: {
  klaim: Klaim
  pesertaID: string
  adjID: string
  onTutup: () => void
  /** Sesudah berhasil: pemanggil menutup modal, memuat ulang, dan menampilkan pita sukses. */
  onTerkirim: () => void | Promise<void>
  /** Pratinjau yang sudah ada (uji); tanpa ini dibaca dari server saat modal dibuka. */
  awal?: PratinjauKomite
  /** Claim RNM + mata uang baris dari layar — konteks bila pratinjau ditolak server (12-10-2026). */
  jumlahKlaim?: string
  mataUang?: string
  /** Galat pratinjau yang sudah ada (uji keadaan pratinjau ditolak); tanpa ini dibaca dari server. */
  galatAwal?: unknown
}) {
  const [pratinjau, setPratinjau] = useState<PratinjauKomite | null>(awal ?? null)
  const [remarks, setRemarks] = useState('')
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(galatAwal ?? null)
  // Remarks disentuh (blur) - petunjuk wajib baru tampil sesudahnya.
  const [disentuh, setDisentuh] = useState(false)
  const remarksKosong = remarks.trim() === ''

  useEffect(() => {
    if (awal !== undefined || galatAwal !== undefined) return
    let batal = false
    ambilPratinjauKomite(klaim.id, pesertaID, adjID)
      .then((h) => {
        if (!batal) setPratinjau(h)
      })
      .catch((e: unknown) => {
        if (!batal) setGalat(e)
      })
    return () => {
      batal = true
    }
  }, [awal, galatAwal, klaim.id, pesertaID, adjID])

  /** SATU-SATUNYA pemanggil serahkanKeKomite; klik ganda ditahan. */
  async function kirim(): Promise<void> {
    if (sibuk || !bolehKirimKomite(pratinjau, remarks, klaim.isKpr)) return
    setSibuk(true)
    setGalat(null)
    try {
      await serahkanKeKomite(klaim.id, pesertaID, adjID, remarks.trim())
      await onTerkirim()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  const tanpaKpr = (klaim.isKpr ?? '').trim() === ''
  // Pratinjau DITOLAK (mis. roster kosong): konteks tetap tampil di bawah pita galat (keputusan work owner 12-10-2026) —
  // Date = waktu peramban saat ini (bukan yang tersimpan), PIC tak diketahui, Adjustment dari baris layar, Remarks nonaktif,
  // grid kosong; tombol kirim tetap nonaktif (tanpa anggota).
  const gagalPratinjau = pratinjau === null && galat !== null
  const tampil: PratinjauKomite | null =
    pratinjau ?? (gagalPratinjau ? { waktu: new Date().toISOString(), pic: '', jumlahKlaim, mataUang, anggota: [] } : null)
  return (
    <Modal
      judul={POPUP_KOMITE.judul}
      onTutup={onTutup}
      tanpaTutup={sibuk}
      labelBatal={TOMBOL_KOMITE.batal}
      aksi={
        <>
          {/* Sebab tombol nonaktif karena Remarks (sebab lain sudah dinyatakan di isi modal). */}
          {!sibuk && remarksKosong && tampil !== null && !gagalPratinjau && (
            <span className="komite-popup__bantu">{POPUP_KOMITE.isiRemarks}</span>
          )}
          <button
            type="button"
            className="btn btn--primary"
            disabled={sibuk || !bolehKirimKomite(pratinjau, remarks, klaim.isKpr)}
            onClick={() => void kirim()}
          >
            {sibuk ? POPUP_KOMITE.proses : TOMBOL_KOMITE.serahkan}
          </button>
        </>
      }
    >
      {/* 12-10-2026 ("rapikan tampilan send to committee ini"): wadah berkelas untuk lebar modal ≈720px dan jarak 16px. */}
      <div className="komite-popup">
        <PitaGalat galat={galat} onTutup={() => setGalat(null)} />
        {tampil === null && <p role="status">{POPUP_KOMITE.memuat}</p>}
        {tampil !== null && (
          <>
            {/* Tiga kolom kecil berdampingan: Date | PIC | Adjustment (satu kolom di layar sempit). */}
            <dl className="komite-popup__info">
              <div>
                <dt>{POPUP_KOMITE.tanggal}</dt>
                <dd>{tanggalWaktuKomite(tampil.waktu)}</dd>
              </div>
              <div>
                <dt>{POPUP_KOMITE.pic}</dt>
                <dd>{tampil.pic || '—'}</dd>
              </div>
              <div>
                <dt>{POPUP_KOMITE.adjustment}</dt>
                <dd>
                  {tampil.jumlahKlaim === ''
                    ? '—'
                    : `${formatNumber(tampil.jumlahKlaim, DESIMAL_TAK_DIBATASI)} ${tampil.mataUang}`.trim()}
                </dd>
              </div>
            </dl>
            <label className="komite-popup__remarks">
              <span className="os__label">
                {POPUP_KOMITE.remarks} <span aria-hidden="true">*</span>
              </span>
              <textarea
                className="field__input"
                rows={4}
                required
                maxLength={BATAS_REMARKS_KOMITE}
                value={remarks}
                disabled={sibuk || gagalPratinjau}
                aria-invalid={disentuh && remarksKosong}
                onChange={(e) => setRemarks(e.target.value)}
                onBlur={() => setDisentuh(true)}
              />
              <span className="komite-popup__kaki-isian">
                {disentuh && remarksKosong ? (
                  <span role="alert" className="komite-popup__galat">
                    {POPUP_KOMITE.wajibRemarks}
                  </span>
                ) : (
                  <span />
                )}
                <span className="komite-popup__hitung" aria-live="polite">
                  {hitungRemarks(remarks)}
                </span>
              </span>
            </label>
            {tanpaKpr && (
              <p role="note" className="alert alert--warn">
                {POPUP_KOMITE.tanpaKpr}
              </p>
            )}
            <section className="komite-popup__komite" aria-label={POPUP_KOMITE.judulKomite}>
              <h4 className="komite-popup__judul">{POPUP_KOMITE.judulKomite}</h4>
              {/* Kartu tabel = pembungkus sendiri (tepi + sudut + overflow hidden), tabel di dalamnya polos - tepi kartu tidak
                  pernah terpotong area gulir badan modal (12-10-2026). */}
              <div className="komite-popup__kartu">
              <table className="inbox__tabel komite-popup__tabel">
                <thead>
                  <tr>
                    {POPUP_KOMITE.kolom.map((k) => (
                      <th key={k}>{k}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {tampil.anggota.length === 0 && (
                    <tr>
                      <td colSpan={4}>{POPUP_KOMITE.tanpaAnggota}</td>
                    </tr>
                  )}
                  {tampil.anggota.map((a) => (
                    <tr key={`${a.urut}:${a.namaKomite}`}>
                      <td>{a.namaKomite || '—'}</td>
                      <td>
                        <span className="komite-popup__chip">{POPUP_KOMITE.menunggu}</span>
                      </td>
                      <td>—</td>
                      <td>—</td>
                    </tr>
                  ))}
                </tbody>
              </table>
              </div>
            </section>
          </>
        )}
      </div>
    </Modal>
  )
}
