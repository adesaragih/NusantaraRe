// Popup Following Policy Number (tiket 70) - port harness `NB FacIn\Harness\ViewFollowingNB.xml` -> section
// `ViewOpenFollowingPolicy` (View Option + Filter + Refresh, judul "Policy List") + grid `GridViewFollowingNB` (20 per halaman;
// link View Details / Copy All / Copy Nopolis; View Details & Copy Nopolis disabled bila Policy Number kosong).
//
// Sistem baru (work owner 08-10-2026 K70-1 "disistem baru kan tidak pakai json_data lgi, buatin sesuai sistem baru dong"):
// - cari = tabel FACINPRODUCTION (backend); Pega mencari saat Enter - di sini juga saat mengetik / ganti View Option (pola
//   PopupPilihAgent: laporan work owner 03-10-2026 tombol search "belum berfungsi" bila hanya Enter);
// - Copy Nopolis / Copy All ditulis backend lalu layar pembuka memuat ulang case (Pega: opener.location.reload + window.close);
//   Copy All mengganti seluruh lokasi -> konfirmasi (A70-2);
// - View Details = isi case sumber baca-saja dari tabel (Object / Coverage / Attachment File), tombol Back kembali ke daftar
//   (A70-3; Pega membaca JSON_POLIS dan mengganti pyWorkPage).

import { useEffect, useRef, useState, type ReactNode } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  cariFollowing,
  salinNopolisFollowing,
  salinSemuaFollowing,
  type BarisFollowing,
  type HalamanFollowing,
  type KasusNB,
  type OpsiFollowing,
} from '../api'
import { FOLLOWING as F, TEKS_FOLLOWING as T } from '../labels'
import HalamanPopup from './HalamanPopup'
import { JEDA_CARI_MS } from './PopupPilihAgent'

/** Tanggal kabel `DD-MM-YYYY` -> tampil `DD/MM/YYYY` (grid Pega pxDateTime). */
export const tanggalFollowing = (v: string) => v.replace(/^(\d{2})-(\d{2})-(\d{4})/, '$1/$2/$3')

export default function PopupFollowing({
  caseId,
  onTutup,
  onTersalin,
  tampilDetail,
}: {
  caseId: string
  onTutup: () => void
  /**
   * View Details: layar case sumber LENGKAP mode hanya lihat (General + semua tab + riwayat + lampiran), dipasok pembuka
   * (InwardFacultative) supaya modul ini tidak mengimpor halamannya sendiri. `kembali` = tombol Back.
   */
  tampilDetail: (baris: BarisFollowing, kembali: () => void) => ReactNode
  /** Case sudah diubah server (Copy Nopolis / Copy All) - pembuka memuat ulang. */
  onTersalin: (k: KasusNB) => void
}) {
  const [opsi, setOpsi] = useState<OpsiFollowing | ''>('')
  const [filter, setFilter] = useState('')
  const [hasil, setHasil] = useState<HalamanFollowing | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  const [menyalin, setMenyalin] = useState(false)
  const [yakinSemua, setYakinSemua] = useState<BarisFollowing | null>(null)
  const [detail, setDetail] = useState<BarisFollowing | null>(null)
  const nomorPermintaan = useRef(0)

  async function muat(o: OpsiFollowing | '', f: string, halaman: number) {
    const nomor = ++nomorPermintaan.current
    // GetListViewFollowing_act langkah 3: filter kosong -> tidak mencari.
    if (o === '' || f.trim() === '') {
      setHasil(null)
      setMemuat(false)
      return
    }
    setMemuat(true)
    setGalat(null)
    try {
      const h = await cariFollowing(caseId, o, f, halaman)
      if (nomor === nomorPermintaan.current) setHasil(h)
    } catch (err) {
      if (nomor === nomorPermintaan.current) {
        setHasil(null)
        setGalat(err)
      }
    } finally {
      if (nomor === nomorPermintaan.current) setMemuat(false)
    }
  }

  useEffect(() => {
    const jadwal = window.setTimeout(() => void muat(opsi, filter, 1), JEDA_CARI_MS)
    return () => window.clearTimeout(jadwal)
  }, [opsi, filter])

  async function salin(b: BarisFollowing, semua: boolean) {
    setMenyalin(true)
    setGalat(null)
    try {
      const k = await (semua ? salinSemuaFollowing : salinNopolisFollowing)(caseId, b)
      onTersalin(k)
      onTutup()
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyalin(false)
      setYakinSemua(null)
    }
  }

  if (detail !== null)
    return (
      <Modal judul={T.judulDetail(detail.policyNo)} onTutup={onTutup} penuh>
        <div className="nbf-following__detail">{tampilDetail(detail, () => setDetail(null))}</div>
      </Modal>
    )

  return (
    <Modal judul={F.judul} onTutup={onTutup} onKirim={() => void muat(opsi, filter, 1)} penuh>
      {/* Dirapikan (work owner 09-10-2026 "ini juga rapikan biar lebih user friendly"): pola halaman View Policy - bilah saring
          sebaris, tabel kartu tema Kelola User, pesan kosong di dalam tabel, tiga aksi sebaris per baris. */}
      <div className="toolbar nbf-following__cari">
        <select
          className="field__input nbf-vp__opsi"
          aria-label={F.viewOption}
          value={opsi}
          onChange={(e) => setOpsi(e.target.value as OpsiFollowing | '')}
        >
          {/* Pilihan kosong = "Choose" saja seperti Pega (work owner 11-10-2026); nama "View Option" tetap di aria-label. */}
          <option value="">{F.pilihKosong}</option>
          {F.opsi.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
        <input
          className="field__input nbf-saring__kotak"
          type="search"
          aria-label={F.filter}
          placeholder={F.placeholder}
          value={filter}
          disabled={opsi === ''}
          onChange={(e) => setFilter(e.target.value)}
        />
        {/* ResetFilterValue_DT: kosongkan Filter (View Option tetap). */}
        <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setFilter('')}>
          {F.refresh}
        </button>
      </div>
      <div className="nbf-following__kepala">
        <h4 className="nbf-following__judul">
          {F.daftar}
          {hasil !== null && <span className="nbf-following__jumlah">{hasil.total}</span>}
        </h4>
        {hasil && hasil.baris.length > 0 && (
          <div className="nbf-tema__pager">
            <HalamanPopup halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(opsi, filter, h)} />
          </div>
        )}
      </div>
      <div className="table-wrap">
        <table className="nbf-tabel nbf-tema__tabel">
          <thead>
            <tr>
              {F.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" />
            </tr>
          </thead>
          <tbody>
            {(memuat || hasil === null || hasil.baris.length === 0) && (
              <tr>
                <td colSpan={F.kolom.length + 1} className="nbf-following__kosong">
                  {memuat ? <Memuat /> : hasil === null ? T.isiFilter : T.kosong}
                </td>
              </tr>
            )}
            {!memuat &&
              (hasil?.baris ?? []).map((b) => (
                <tr key={`${b.caseId}-${b.policyNo}`}>
                  <td>{b.caseId}</td>
                  <td className="nbf-following__nopolis">{b.policyNo}</td>
                  <td>{b.reffNumber}</td>
                  <td>{b.insuredName}</td>
                  <td>{b.sourceOfBusiness}</td>
                  <td>{b.qqName}</td>
                  <td>{tanggalFollowing(b.startPeriod)}</td>
                  <td>{tanggalFollowing(b.endPeriod)}</td>
                  <td className="table__actions">
                    <span className="nbf-following__aksi">
                      <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" disabled={b.policyNo === ''} onClick={() => setDetail(b)}>
                        {F.lihat}
                      </button>
                      <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" disabled={menyalin} onClick={() => setYakinSemua(b)}>
                        {F.salinSemua}
                      </button>
                      <button
                        type="button"
                        className="btn nbf-tema__tombol btn--ghost btn--sm"
                        disabled={menyalin || b.policyNo === ''}
                        onClick={() => void salin(b, false)}
                      >
                        {F.salinNopolis}
                      </button>
                    </span>
                  </td>
                </tr>
              ))}
          </tbody>
        </table>
      </div>
      {menyalin && <p className="muted">{T.menyalin}</p>}
      <Gagal galat={galat} />
      {yakinSemua && (
        <Modal
          judul={F.salinSemua}
          onTutup={() => setYakinSemua(null)}
          aksi={
            <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" disabled={menyalin} onClick={() => void salin(yakinSemua, true)}>
              {T.ya}
            </button>
          }
        >
          <p>{T.yakinSalinSemua(yakinSemua.policyNo)}</p>
        </Modal>
      )}
    </Modal>
  )
}
