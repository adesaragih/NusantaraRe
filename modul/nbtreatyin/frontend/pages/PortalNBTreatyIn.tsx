// Portal NB Treaty In - Harness `SFAPortalOpportunities`: daftar berkas
// (`Section/SFAPortal_OpportunitiesList`) dengan saringan teks dan tombol
// "Create opportunity" (`SFAPortalOpportunitiesHeader`, `createWork`).
//
// ⛔ Bagian portal CRM yang lain - ringkasan prospek (`SFAPortal_
// OpportunitiesList_Header`, bersyarat tampil `1=2`), Stage view/List view -
// milik modul CRM dan tidak dibangun.
//
// W6 audit silang P3 (dibaca ulang 04-10-2026): grid `GetListOpportunity` -
// LABEL "Offer No" (`.TextNoQuotation`), "Name" (`.Name`), "Group Business",
// "Insured Name", "Marketing", "Status". Kolom "Position" / "No Polis" (tanpa sel
// XML) dibuang. ⛔ RALAT tiket 11: kolom "Name" (pxLink `.Name` -> `openWorkByHandle
// .pzInsKey`) tidak dirender - `.Name` ditulis nol rule korpus, tak berkolom, jadi
// tautannya selalu tanpa teks; pembuka berkas (kunci yang sama) dipasang di sel
// "Offer No". Teks daftar kosong = komponen `Kosong` inti.

import { useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbil } from '../ambil'
import { buatKasus, daftarKasus, type RingkasanKasus } from '../api'
import { KOLOM_PORTAL, PORTAL, TOMBOL } from '../labels'

/** Grid `GetListOpportunity` - kolom VERBATIM `SFAPortal_OpportunitiesList` C[1.x]/C[2.x]. */
export function TabelPortal({ baris, onBuka }: { baris: RingkasanKasus[]; onBuka: (id: string) => void }) {
  return (
    <div className="table-wrap nbti__tabel-portal-wadah">
      <table className="nbti__tabel-portal">
        <thead>
          <tr>
            <th scope="col">{KOLOM_PORTAL.id}</th>
            <th scope="col">{KOLOM_PORTAL.bisnis}</th>
            <th scope="col">{KOLOM_PORTAL.tertanggung}</th>
            <th scope="col">{KOLOM_PORTAL.marketing}</th>
            <th scope="col">{KOLOM_PORTAL.status}</th>
          </tr>
        </thead>
        <tbody>
          {baris.map((b) => (
            <tr key={b.id}>
              <td data-label={KOLOM_PORTAL.id}>
                {/* `.TextNoQuotation`; tautan openWorkByHandle (dari sel `.Name`, RALAT tiket 11) */}
                <button type="button" className="nbti__tautan" onClick={() => onBuka(b.id)}>
                  {b.id}
                </button>
              </td>
              <td data-label={KOLOM_PORTAL.bisnis}>{b.businessName}</td>
              <td data-label={KOLOM_PORTAL.tertanggung}>{b.insuredName}</td>
              <td data-label={KOLOM_PORTAL.marketing}>{b.marketingName}</td>
              <td data-label={KOLOM_PORTAL.status}>{b.nbStatus !== '' && <span className="nbti__status">{b.nbStatus}</span>}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

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
    <div className="inbox nbti__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{PORTAL.judul}</h2>
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
          placeholder={PORTAL.placeholder}
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
        {/* ikon pengosong C[1.2]: click -> setValue .FilterTermForOpportunity "" -> postValue, TANPA refresh */}
        <button type="button" className="btn btn--ghost" aria-label={PORTAL.bersihkan} onClick={() => setCari('')}>
          ×
        </button>
        <button type="submit" className="btn">
          {TOMBOL.filter}
        </button>
      </form>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && baris.length === 0 && <Kosong pesan={PORTAL.kosong} petunjuk={PORTAL.kosongPetunjuk} />}
      {baris !== null && baris.length > 0 && <TabelPortal baris={baris} onBuka={onBuka} />}
    </div>
  )
}
