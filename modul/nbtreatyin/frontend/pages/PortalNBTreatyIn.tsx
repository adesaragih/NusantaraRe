// Portal NB Treaty In - Harness `SFAPortalOpportunities`: daftar berkas
// (`Section/SFAPortal_OpportunitiesList`) dengan saringan teks dan tombol
// "Create opportunity" (`SFAPortalOpportunitiesHeader`, `createWork`).
//
// ⛔ Bagian portal CRM yang lain - ringkasan prospek (`SFAPortal_
// OpportunitiesList_Header`, bersyarat tampil `1=2`), Stage view/List view -
// milik modul CRM dan tidak dibangun.

import { useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbil } from '../ambil'
import { buatKasus, daftarKasus } from '../api'
import { JUDUL, JUDUL_POSISI, KOLOM_PORTAL, PORTAL, TOMBOL } from '../labels'

export default function PortalNBTreatyIn({ onBuka, pesan }: { onBuka: (id: string) => void; pesan?: string }) {
  const [cari, setCari] = useState('')
  const [kueri, setKueri] = useState('')
  const { data: baris, galat: galatDaftar } = useAmbil(() => daftarKasus(kueri), [kueri])
  const [galatBuat, setGalatBuat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const galat = galatBuat ?? galatDaftar

  const buat = async () => {
    setSibuk(true)
    setGalatBuat(null)
    try {
      const k = await buatKasus()
      onBuka(k.id)
    } catch (e: unknown) {
      setGalatBuat(e)
    } finally {
      setSibuk(false)
    }
  }

  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{JUDUL.portal}</h2>
        <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void buat()}>
          {TOMBOL.create}
        </button>
      </header>
      {pesan && <div className="alert alert--ok">{pesan}</div>}
      <form
        className="nbti__saring"
        onSubmit={(e) => {
          e.preventDefault()
          setKueri(cari.trim())
        }}
      >
        <input
          className="field__input"
          aria-label={PORTAL.filter}
          placeholder={PORTAL.filter}
          value={cari}
          onChange={(e) => setCari(e.target.value)}
          // `.FilterTermForOpportunity` esc -> setValue "" -> refresh (enter = submit form)
          onKeyDown={(e) => {
            if (e.key === 'Escape') {
              setCari('')
              setKueri('')
            }
          }}
        />
        <button type="submit" className="btn">
          {TOMBOL.filter}
        </button>
      </form>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && baris.length === 0 && <Kosong pesan={PORTAL.kosong} petunjuk={PORTAL.kosongPetunjuk} />}
      {baris !== null && baris.length > 0 && (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_PORTAL.id}</th>
                <th scope="col">{KOLOM_PORTAL.bisnis}</th>
                <th scope="col">{KOLOM_PORTAL.tertanggung}</th>
                <th scope="col">{KOLOM_PORTAL.marketing}</th>
                <th scope="col">{KOLOM_PORTAL.status}</th>
                <th scope="col">{KOLOM_PORTAL.posisi}</th>
                <th scope="col">{KOLOM_PORTAL.nopol}</th>
              </tr>
            </thead>
            <tbody>
              {baris.map((b) => (
                <tr key={b.id}>
                  <td>
                    <button type="button" className="nbti__tautan" onClick={() => onBuka(b.id)}>
                      {b.id}
                    </button>
                  </td>
                  <td>{b.businessName}</td>
                  <td>{b.insuredName}</td>
                  <td>{b.marketingName}</td>
                  <td>{b.nbStatus}</td>
                  <td>{JUDUL_POSISI[b.positionNote] ?? b.positionNote}</td>
                  <td>{b.noPolis}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
