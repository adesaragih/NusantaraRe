// Halaman depan NB FacIn - portal Opportunity (tiket 25; daftar case NB tiket 32).
//
// Port harness `SFAPortalOpportunities` (`D:\migrasi\RNM\NB FacIn\Harness\SFAPortalOpportunities.xml`,
// PEGACRM-PORTAL!SFAPORTALOPPORTUNITIES) beserta section `SFAPortalOpportunitiesHeader` dan
// `SFAPortal_OpportunitiesList`. Hanya unsur yang TAMPIL di Pega yang diport; yang tersembunyi
// permanen (`1=2` / `NEVER`: Stage view, List view, All/Individual/Corporate, Phase, Export, Refresh)
// tidak dirender sama sekali.
//
// Grid = `GetListOpportunityF` (keputusan agent B-1): Offer No · Name (tautan) · (View) · Group Business ·
// Insured Name · Marketing · (NBStatus) · Status. Datanya = case NB yang dibuat Create opportunity
// (permintaan work owner 03-10-2026: "nb yang sudah di create, muncul disini … harus ada case id nya, group
// business, insured name, status") lewat `GET /api/nbfacin/opportunity` (tiket 32): Offer No = nomor case,
// Status = `T_WORK_POLIS.STATUS_WORK`. Kotak saring (placeholder Pega "NB-1234 or Name") + Filter / Enter
// mencari "mengandung"; ikon ✕ mengosongkan. Klik Name membuka case (Pega `openWorkByHandle`). Kolom View dan
// NBStatus (judul kosong di Pega) tetap ada tetapi kosong - tautan View (`ViewOutstandingCase`) belum diport.
//
// `Create opportunity` membuka form Opportunity (B-3, tiket 26). Nol catatan pengembang di layar (B-4).

import { useEffect, useRef, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { daftarCaseNB, type HalamanCaseNB } from '../api'
import { KEPALA_PORTAL, KOLOM_PORTAL, SARING_PORTAL, TEKS_PORTAL } from '../labels'

export default function PortalOpportunity({ onBuat, onBuka }: { onBuat: () => void; onBuka: (caseId: string) => void }) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<HalamanCaseNB | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  // Nomor permintaan terakhir: jawaban lama dibuang (pola popup ChooseAccount).
  const nomorPermintaan = useRef(0)
  const kunciCari = useRef('')

  async function muat(cari: string, halaman: number) {
    const nomor = ++nomorPermintaan.current
    kunciCari.current = cari
    setMemuat(true)
    setGalat(null)
    try {
      const h = await daftarCaseNB(cari, halaman)
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

  return (
    <div className="nbfacin">
      <section className="panel">
        <div className="nbf-kepala">
          <h4 className="panel__title">{KEPALA_PORTAL.judul.label}</h4>
          <button type="button" className="btn btn--primary" onClick={onBuat}>
            {KEPALA_PORTAL.buat.label}
          </button>
        </div>
        <form
          className="nbf-saring"
          onSubmit={(e) => {
            e.preventDefault()
            void muat(kotak, 1)
          }}
        >
          <input
            className="field__input nbf-saring__kotak"
            type="text"
            aria-label={SARING_PORTAL.label.label}
            placeholder={SARING_PORTAL.placeholder.label}
            value={kotak}
            onChange={(e) => setKotak(e.target.value)}
          />
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            aria-label={TEKS_PORTAL.hapusIsian}
            onClick={() => {
              setKotak('')
              void muat('', 1)
            }}
          >
            ✕
          </button>
          <button type="submit" className="btn btn--sm" disabled={memuat}>
            {SARING_PORTAL.tombol.label}
          </button>
        </form>
        {hasil && hasil.baris.length > 0 && (
          <Halaman halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(kunciCari.current, h)} />
        )}
        <div className="table-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {KOLOM_PORTAL.map((k) => (
                  <th key={k.sel} scope="col">
                    {k.label}
                  </th>
                ))}
              </tr>
            </thead>
            {hasil && hasil.baris.length > 0 && (
              <tbody>
                {hasil.baris.map((b) => (
                  <tr key={b.caseId}>
                    <td>{b.caseId}</td>
                    <td>
                      <button type="button" className="nbf-tautan" onClick={() => onBuka(b.caseId)}>
                        {b.name || b.caseId}
                      </button>
                    </td>
                    <td />
                    <td>{b.groupBusiness}</td>
                    <td>{b.insuredName}</td>
                    <td>{b.marketing}</td>
                    <td />
                    <td>{b.status}</td>
                  </tr>
                ))}
              </tbody>
            )}
          </table>
        </div>
        {memuat && !hasil && <Memuat />}
        {hasil && hasil.baris.length === 0 && <Kosong pesan={TEKS_PORTAL.tanpaCase} />}
        <Gagal galat={galat} />
      </section>
    </div>
  )
}
