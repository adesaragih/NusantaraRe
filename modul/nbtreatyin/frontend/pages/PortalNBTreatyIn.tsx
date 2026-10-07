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
//
// RALAT 06-10-2026 (keputusan work owner): daftar HANYA berkas buatan akun ini (filter pembuat, LINI non-life tetap),
// dengan switch In Progress / Resolved di atasnya, bawaan In Progress.
//
// Tombol Copy Old di samping Create (perintah work owner 07-10-2026, khusus superuser; bukan layar Pega): tampil bila
// `GET /hak` menyatakan superadmin ber-hak penuh; popup `components/DialogCopyOld.tsx`; portal dimuat ulang bila ada
// dokumen yang tersalin.

import { useState } from 'react'

import { Gagal, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbil } from '../ambil'
import { ambilHak, buatKasus, daftarKasus, type RingkasanKasus } from '../api'
import DialogCopyOld from '../components/DialogCopyOld'
import {
  COPY_OLD,
  KOLOM_PORTAL,
  KOLOM_PORTAL_SELESAI,
  PORTAL,
  STATUS_PORTAL,
  TOMBOL,
  type StatusPortal,
} from '../labels'
import { tampilCopyOld } from '../lama'
import { sajikan } from '../sajian'

/** Grid `GetListOpportunity` - kolom VERBATIM `SFAPortal_OpportunitiesList` C[1.x]/C[2.x]; tab Resolved (`selesai`)
 *  menambah Policy Number sesudah Offer No dan Production Date di akhir (keputusan work owner 07-10-2026). */
export function TabelPortal({
  baris,
  onBuka,
  selesai = false,
}: {
  baris: RingkasanKasus[]
  onBuka: (id: string) => void
  selesai?: boolean
}) {
  return (
    <div className="table-wrap nbti__tabel-portal-wadah">
      <table className="nbti__tabel-portal">
        <thead>
          <tr>
            <th scope="col">{KOLOM_PORTAL.id}</th>
            {selesai && <th scope="col">{KOLOM_PORTAL_SELESAI.noPolis}</th>}
            <th scope="col">{KOLOM_PORTAL.jenis}</th>
            <th scope="col">{KOLOM_PORTAL.bisnis}</th>
            <th scope="col">{KOLOM_PORTAL.tertanggung}</th>
            <th scope="col">{KOLOM_PORTAL.marketing}</th>
            <th scope="col">{KOLOM_PORTAL.status}</th>
            <th scope="col">{KOLOM_PORTAL.pembuat}</th>
            <th scope="col">{KOLOM_PORTAL.tanggal}</th>
            {selesai && <th scope="col">{KOLOM_PORTAL_SELESAI.tglProduksi}</th>}
          </tr>
        </thead>
        <tbody>
          {baris.map((b) => {
            // berkas Resolved: status T_WORK_POLIS.STATUS_WORK (Resolved-Completed / Resolved-Rejected), bukan NBStatus
            // terakhir sebelum ditutup (perintah work owner 06-10-2026: "pake status resolve di work_polis aja")
            const status = b.statusWork.startsWith('Resolved-') ? b.statusWork : b.nbStatus
            return (
              <tr key={b.id}>
                <td data-label={KOLOM_PORTAL.id}>
                  {/* `.TextNoQuotation`; tautan openWorkByHandle (dari sel `.Name`, RALAT tiket 11) */}
                  <button type="button" className="nbti__tautan" onClick={() => onBuka(b.id)}>
                    {b.id}
                  </button>
                </td>
                {selesai && <td data-label={KOLOM_PORTAL_SELESAI.noPolis}>{b.noPolis}</td>}
                {/* T_POLIS_QUOTATION.PROPORTIONAL_TYPE apa adanya - tidak disingkat (perintah work owner 06-10-2026) */}
                <td data-label={KOLOM_PORTAL.jenis}>{b.proportionalType}</td>
                <td data-label={KOLOM_PORTAL.bisnis}>{b.businessName}</td>
                <td data-label={KOLOM_PORTAL.tertanggung}>{b.insuredName}</td>
                <td data-label={KOLOM_PORTAL.marketing}>{b.marketingName}</td>
                <td data-label={KOLOM_PORTAL.status}>
                  {status !== '' && <span className="nbti__status">{status}</span>}
                </td>
                <td data-label={KOLOM_PORTAL.pembuat}>{b.createOpName}</td>
                <td data-label={KOLOM_PORTAL.tanggal}>{sajikan(b.tglCreate, 'tanggal')}</td>
                {selesai && (
                  <td data-label={KOLOM_PORTAL_SELESAI.tglProduksi}>{sajikan(b.productionDate, 'tanggal')}</td>
                )}
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

/** Kepala portal: judul, tombol Copy Old (superadmin, perintah work owner 07-10-2026) di samping tombol Create. */
export function KepalaPortal({
  copyOld,
  sibuk,
  onCopyOld,
  onBuat,
}: {
  copyOld: boolean
  sibuk: boolean
  onCopyOld: () => void
  onBuat: () => void
}) {
  return (
    <header className="inbox__kepala">
      <h2 className="inbox__judul">{PORTAL.judul}</h2>
      <div className="nbti__kepala-aksi">
        {copyOld && (
          <button type="button" className="btn" onClick={onCopyOld}>
            {COPY_OLD.tombol}
          </button>
        )}
        <button type="button" className="btn btn--primary" disabled={sibuk} onClick={onBuat}>
          {TOMBOL.create}
        </button>
      </div>
    </header>
  )
}

export default function PortalNBTreatyIn({
  onBuka,
  pesan,
  status,
  onStatus,
}: {
  onBuka: (id: string) => void
  pesan?: string
  /** Posisi switch; dipegang rute supaya bertahan sesudah Back dari layar kasus. Tanpa ini: keadaan lokal. */
  status?: StatusPortal
  onStatus?: (s: StatusPortal) => void
}) {
  const [cari, setCari] = useState('')
  const [kueri, setKueri] = useState('')
  const [statusLokal, setStatusLokal] = useState<StatusPortal>(STATUS_PORTAL[0])
  const aktif = status ?? statusLokal
  const selesai = aktif === STATUS_PORTAL[1]
  // ketuk - muat ulang sesudah Copy Old menyalin (saringan sama)
  const [ketuk, setKetuk] = useState(0)
  const { data: baris, galat: galatDaftar } = useAmbil(() => daftarKasus(kueri, '', selesai), [kueri, selesai, ketuk])
  const { data: hak } = useAmbil(() => ambilHak(), [])
  const [copyOld, setCopyOld] = useState(false)
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
      <KepalaPortal
        copyOld={tampilCopyOld(hak)}
        sibuk={sibuk}
        onCopyOld={() => setCopyOld(true)}
        onBuat={() => void buat()}
      />
      {pesan && <div className="alert alert--ok">{pesan}</div>}
      {/* switch In Progress / Resolved (keputusan work owner 06-10-2026, bawaan In Progress) */}
      <StripTab tab={STATUS_PORTAL} aktif={aktif} onPilih={onStatus ?? setStatusLokal} />
      {/* satu kapsul: ikon cari, isian, pengosong, tombol Filter (perintah work owner 06-10-2026: "rapihkan, keren, elegan") */}
      <form
        className="nbti__saring"
        role="search"
        onSubmit={(e) => {
          e.preventDefault()
          setKueri(cari.trim())
        }}
      >
        <svg className="nbti__saring-ikon" viewBox="0 0 20 20" aria-hidden="true">
          <circle cx="9" cy="9" r="6" />
          <path d="m13.5 13.5 4 4" />
        </svg>
        <input
          className="nbti__saring-isian"
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
        <button
          type="button"
          className="nbti__saring-hapus"
          aria-label={PORTAL.bersihkan}
          disabled={cari === ''}
          onClick={() => setCari('')}
        >
          ×
        </button>
        <button type="submit" className="nbti__saring-tombol">
          {TOMBOL.filter}
        </button>
      </form>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null &&
        baris.length === 0 &&
        (selesai ? (
          <Kosong pesan={PORTAL.kosongSelesai} />
        ) : (
          <Kosong pesan={PORTAL.kosong} petunjuk={PORTAL.kosongPetunjuk} />
        ))}
      {baris !== null && baris.length > 0 && <TabelPortal baris={baris} onBuka={onBuka} selesai={selesai} />}
      {copyOld && (
        <DialogCopyOld
          onTutup={(adaYangDisalin) => {
            setCopyOld(false)
            if (adaYangDisalin) setKetuk((k) => k + 1)
          }}
        />
      )}
    </div>
  )
}
