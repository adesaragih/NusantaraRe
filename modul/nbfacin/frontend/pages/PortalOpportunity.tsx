// Halaman depan NB FacIn - portal Opportunity (tiket 25; daftar case NB tiket 32).
//
// Port harness `SFAPortalOpportunities` (`D:\migrasi\RNM\NB FacIn\Harness\SFAPortalOpportunities.xml`,
// PEGACRM-PORTAL!SFAPORTALOPPORTUNITIES) beserta section `SFAPortalOpportunitiesHeader` dan
// `SFAPortal_OpportunitiesList`. Hanya unsur yang TAMPIL di Pega yang diport; yang tersembunyi
// permanen (`1=2` / `NEVER`: Stage view, List view, All/Individual/Corporate, Phase, Export, Refresh)
// tidak dirender sama sekali.
//
// Grid = `GetListOpportunityF` (keputusan agent B-1): Offer No · Name · (View) · Group Business ·
// Insured Name · Marketing · NB Status · Status. Datanya = case NB yang dibuat Create opportunity
// (permintaan work owner 03-10-2026) lewat `GET /api/nbfacin/opportunity` (tiket 32): Offer No = nomor case,
// Status = `T_WORK_POLIS.STATUS_WORK`. Kotak saring (placeholder Pega "NB-1234 or Name") + Filter / Enter
// mencari "mengandung"; ✕ mengosongkan. Kolom View (judul kosong di Pega) tetap ada tetapi kosong - tautan View
// (`ViewOutstandingCase`) belum diport.
//
// Tata letak = halaman depan NB Treaty In (permintaan work owner 11-10-2026 "tampilan depan nb facin, samain dong kaya nb
// treaty in", `modul/nbtreatyin/frontend/pages/PortalNBTreatyIn.tsx`): kepala judul + tombol Create di kanan, switch
// In Progress / Resolved (bawaan In Progress; saringan `selesai` di server), saringan SATU kapsul (ikon cari, isian,
// pengosong, tombol Filter), nomor case = tautan pembuka case, status = chip. Tab Resolved menambah Policy Number sesudah
// Offer No dan Production Date di akhir (pola NB Treaty In; kosong bila server tidak mengirimnya). Tombol Copy Old NB Treaty
// In TIDAK dibawa - fitur penyalin dokumen lama milik NB Treaty In, tanpa padanan di NB FacIn.

import { useEffect, useRef, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { daftarCaseNB, type HalamanCaseNB } from '../api'
import { KEPALA_PORTAL, KOLOM_PORTAL, KOLOM_PORTAL_SELESAI, SARING_PORTAL, STATUS_PORTAL, TEKS_PORTAL, type StatusPortal } from '../labels'

/**
 * `pesan` - konfirmasi sesudah Submit case berhasil (tiket 57: case ditutup dan pengguna kembali ke portal, seperti
 * `finishAssignment` Pega); `onTutupPesan` menghapusnya. `status` / `onStatus` - posisi switch dipegang rute supaya
 * bertahan sesudah kembali dari layar case (pola NB Treaty In); tanpa itu keadaan lokal.
 */
export default function PortalOpportunity({
  onBuat,
  onBuka,
  pesan = '',
  onTutupPesan,
  status,
  onStatus,
}: {
  onBuat: () => void
  onBuka: (caseId: string) => void
  pesan?: string
  onTutupPesan?: () => void
  status?: StatusPortal
  onStatus?: (s: StatusPortal) => void
}) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<HalamanCaseNB | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  const [statusLokal, setStatusLokal] = useState<StatusPortal>(STATUS_PORTAL[0])
  const aktif = status ?? statusLokal
  const selesai = aktif === STATUS_PORTAL[1]
  // Nomor permintaan terakhir: jawaban lama dibuang (pola popup ChooseAccount).
  const nomorPermintaan = useRef(0)
  const kunciCari = useRef('')

  async function muat(cari: string, halaman: number, sudahSelesai: boolean = selesai) {
    const nomor = ++nomorPermintaan.current
    kunciCari.current = cari
    setMemuat(true)
    setGalat(null)
    try {
      const h = await daftarCaseNB(cari, halaman, sudahSelesai)
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
    void muat('', 1)
  }, [])

  const pilihStatus = (s: StatusPortal) => {
    if (onStatus) onStatus(s)
    else setStatusLokal(s)
    setHasil(null)
    void muat(kunciCari.current, 1, s === STATUS_PORTAL[1])
  }

  return (
    <div className="nbfacin">
      {pesan !== '' && (
        <div className="alert alert--ok nbf-portal__pesan">
          <span>{pesan}</span>
          {onTutupPesan && (
            <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" aria-label="Tutup" onClick={onTutupPesan}>
              ✕
            </button>
          )}
        </div>
      )}
      <section className="inbox nbf-tema">
        <header className="inbox__kepala nbf-portal__kepala">
          <h2 className="inbox__judul nbf-tema__judul">{KEPALA_PORTAL.judul.label}</h2>
          <div className="nbf-portal__aksi">
            <button type="button" className="btn btn--primary nbf-tema__tombol nbf-tema__tombol--utama" onClick={onBuat}>
              {KEPALA_PORTAL.buat.label}
            </button>
          </div>
        </header>
        <StripTab tab={STATUS_PORTAL} aktif={aktif} onPilih={pilihStatus} />
        {/* satu kapsul: ikon cari, isian, pengosong, tombol Filter (pola saringan NB Treaty In) */}
        <form
          className="nbf-saring"
          role="search"
          onSubmit={(e) => {
            e.preventDefault()
            void muat(kotak, 1)
          }}
        >
          <svg className="nbf-saring__ikon" viewBox="0 0 20 20" aria-hidden="true">
            <circle cx="9" cy="9" r="6" />
            <path d="m13.5 13.5 4 4" />
          </svg>
          <input
            className="nbf-saring__isian"
            aria-label={SARING_PORTAL.label.label}
            placeholder={SARING_PORTAL.placeholder.label}
            value={kotak}
            onChange={(e) => setKotak(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                setKotak('')
                void muat('', 1)
              }
            }}
          />
          <button
            type="button"
            className="nbf-saring__hapus"
            aria-label={TEKS_PORTAL.hapusIsian}
            disabled={kotak === ''}
            onClick={() => {
              setKotak('')
              void muat('', 1)
            }}
          >
            ×
          </button>
          <button type="submit" className="nbf-saring__tombol" disabled={memuat}>
            {SARING_PORTAL.tombol.label}
          </button>
        </form>
        {hasil && hasil.baris.length > 0 && (
          <div className="nbf-tema__pager">
            <Halaman halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(kunciCari.current, h)} />
          </div>
        )}
        <table className="inbox__tabel nbf-tema__tabel">
          <thead>
            <tr>
              {KOLOM_PORTAL.map((k, i) => [
                <th key={k.sel} scope="col">
                  {k.label}
                </th>,
                selesai && i === 0 && (
                  <th key="noPolis" scope="col">
                    {KOLOM_PORTAL_SELESAI.noPolis}
                  </th>
                ),
              ])}
              {selesai && <th scope="col">{KOLOM_PORTAL_SELESAI.tglProduksi}</th>}
            </tr>
          </thead>
          {hasil && hasil.baris.length > 0 && (
            <tbody>
              {hasil.baris.map((b) => (
                <tr key={b.caseId} className="inbox__baris" onClick={() => onBuka(b.caseId)}>
                  <td>
                    {/* nomor case = tautan pembuka case (Pega openWorkByHandle; pola Offer No NB Treaty In) */}
                    <button
                      type="button"
                      className="nbf-tautan nbf-portal__nomor"
                      onClick={(e) => {
                        // Baris sudah membuka case; cegah panggilan ganda.
                        e.stopPropagation()
                        onBuka(b.caseId)
                      }}
                    >
                      {b.caseId}
                    </button>
                  </td>
                  {selesai && <td>{b.noPolis ?? ''}</td>}
                  <td>{b.name}</td>
                  <td />
                  <td>{b.groupBusiness}</td>
                  <td>{b.insuredName}</td>
                  <td>{b.marketing}</td>
                  <td>{b.nbStatus !== '' && <span className="nbf-portal__status">{b.nbStatus}</span>}</td>
                  <td>{b.status !== '' && <span className="nbf-portal__status">{b.status}</span>}</td>
                  {selesai && <td>{b.tglProduksi ?? ''}</td>}
                </tr>
              ))}
            </tbody>
          )}
        </table>
        {memuat && !hasil && <Memuat />}
        {hasil && hasil.baris.length === 0 && <Kosong pesan={TEKS_PORTAL.tanpaCase} />}
        <Gagal galat={galat} />
      </section>
    </div>
  )
}
