// Portal EDM Treaty In - `Harness/SFAPortal_Endorsement_Treaty` + `Section/SFAPortal_Endorsement_Treaty`. Asal pola:
// `modul/nbtreatyin/frontend/pages/PortalNBTreatyIn.tsx` (06-10-2026: kepala + tombol, kapsul saring, tabel portal).
//
//   S126 "Workarea header"   LABEL "Addendum Treaty" + tombol "Create New Addendum Treaty" (showHarness
//                            `TreatyCreateEdm` -> `components/BuatEDM.tsx`)
//   S116 / S117              `.FilterTermForEndorsement` (placeholder "Policy Number") + tombol Filter (click ->
//                            refresh; Esc -> setValue "" -> refresh)
//   S119                     tombol Refresh (tooltip "Refresh EDM Grid", click -> refresh)
//   S122 grid                RD `InboxEDM_RD2` 10 kolom VERBATIM (`KOLOM_PORTAL`), pyPageSize 50 Numeric; tautan
//                            "EDM Number" nonaktif bila `.pxPages(A).PositionNote != 'ReasTreatyInAdmin'`
//
// Switch In Progress / Resolved (StripTab di atas kapsul saring) - aturan portal NB Treaty In berlaku untuk EDM
// (keputusan work owner 07-10-2026; XML hanya satu grid ber-filter C `pyStatusWork != "Resolved-Completed"`):
// In Progress = berkas buatan akun yang masih proses (tautan aktif hanya di posisi admin, XML); Resolved = semua
// berkas selesai, dibuka hanya-baca. Status dipegang `rute.tsx` (bertahan sesudah Back; klik menu = In Progress).
//
// Tombol Copy Old di samping Create (perintah work owner 07-10-2026 "SAMA SEPERTI MASTER PRODUCTNAME LIFE, KHUSUS BUAT
// SUPERUSER"; bukan layar Pega): tampil bila `GET /hak` menyatakan superadmin ber-hak penuh; popup
// `components/DialogCopyOld.tsx`; portal dimuat ulang bila ada dokumen yang tersalin.

import { useState } from 'react'

import { Gagal, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbil } from '../ambil'
import { POSISI_ADMIN, ambilAcuan, ambilHak, daftarKasus, type RingkasanKasus } from '../api'
import BuatEDM from '../components/BuatEDM'
import DialogCopyOld from '../components/DialogCopyOld'
import Paginasi from '../components/Paginasi'
import {
  COPY_OLD,
  JUDUL,
  KOLOM_PORTAL,
  KOLOM_RESOLVED,
  PORTAL,
  STATUS_PORTAL,
  TOMBOL,
  labelJenisEDM,
  type StatusPortal,
} from '../labels'
import { tampilCopyOld } from '../lama'
import { idTampil, sajikan } from '../sajian'
import { BARIS_PER_HALAMAN_GRID, irisan } from '../paginasi'

/** Tautan sel 1: In Progress aktif hanya bagi kasus di posisi admin (pxLink `pyDisabled` bila PositionNote != ADM);
 *  Resolved selalu aktif - berkas selesai dibuka hanya-baca (aturan portal NB, WO 07-10-2026). */
export const tautanAktif = (b: RingkasanKasus, selesai = false) => selesai || b.positionNote === POSISI_ADMIN

/** Teks kolom Status: berkas tertutup (`Resolved-Completed` / `Resolved-Rejected`) = STATUS_WORK, selain itu NBStatus
 *  (perintah work owner 07-10-2026, pola portal NB). */
const statusPortal = (b: RingkasanKasus) => (b.statusWork.startsWith('Resolved-') ? b.statusWork : b.nbStatus)

/** Grid `InboxEDM_RD2` - kolom VERBATIM `SFAPortal_Endorsement_Treaty` S122. */
export function TabelPortal({
  baris,
  onBuka,
  selesai = false,
}: {
  baris: RingkasanKasus[]
  onBuka: (id: string) => void
  selesai?: boolean
}) {
  const [hal, setHal] = useState(1)
  return (
    <>
      <div className="table-wrap edmt__tabel-portal-wadah">
        <table className="edmt__tabel-portal">
          <thead>
            <tr>
              {KOLOM_PORTAL.map((k) => (
                <th key={k.kunci} scope="col">
                  {k.judul}
                </th>
              ))}
              {selesai && <th scope="col">{KOLOM_RESOLVED.tglProd}</th>}
            </tr>
          </thead>
          <tbody>
            {irisan(baris, hal, BARIS_PER_HALAMAN_GRID).map((b) => (
              <tr key={b.id}>
                {KOLOM_PORTAL.map((k) => (
                  <td key={k.kunci} data-label={k.judul}>
                    {k.kunci === 'id' ? (
                      // pxLink "Open Assignment" (openAssignment .pzInsKey)
                      <button
                        type="button"
                        className="edmt__tautan"
                        disabled={!tautanAktif(b, selesai)}
                        onClick={() => onBuka(b.id)}
                      >
                        {idTampil(b.id)}
                      </button>
                    ) : k.kunci === 'nbStatus' ? (
                      // pxDisplayText pyVisible NOTBLANK. WO 07-10-2026 "KALO DAH RESOLVE STATUS NYA PAKE STATUS
                      // RESOLVE" (sama dengan portal NB): berkas Resolved = T_WORK_POLIS.STATUS_WORK, bukan NBStatus
                      statusPortal(b) !== '' && <span className="edmt__status">{statusPortal(b)}</span>
                    ) : k.kunci === 'edmType' ? (
                      // pxDropdown: teks DT TreatyEDMListType (screenshot work owner 07-10-2026)
                      labelJenisEDM(b.edmType)
                    ) : (
                      // nilai DB apa adanya (Proportional Type tanpa prompt values di korpus)
                      b[k.kunci]
                    )}
                  </td>
                ))}
                {selesai && (
                  // WO 07-10-2026: tanggal produksi DD-MM-YYYY (Policy Number sudah kolom XML ke-3)
                  <td data-label={KOLOM_RESOLVED.tglProd}>{sajikan(b.tglProd ?? '', 'tanggal')}</td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {/* pyPageMode Numeric, pyPageSize 50 */}
      <Paginasi jumlahBaris={baris.length} ukuran={BARIS_PER_HALAMAN_GRID} hal={hal} onHal={setHal} />
    </>
  )
}

/** Kepala portal: judul, tombol Copy Old (superadmin) di samping tombol Create. */
export function KepalaPortal({
  copyOld,
  onCopyOld,
  onBuat,
}: {
  copyOld: boolean
  onCopyOld: () => void
  onBuat: () => void
}) {
  return (
    <header className="inbox__kepala">
      <h2 className="inbox__judul">{JUDUL.portal}</h2>
      <div className="edmt__kepala-aksi">
        {copyOld && (
          <button type="button" className="btn" onClick={onCopyOld}>
            {COPY_OLD.tombol}
          </button>
        )}
        <button type="button" className="btn btn--primary" onClick={onBuat}>
          {TOMBOL.create}
        </button>
      </div>
    </header>
  )
}

export default function PortalEDMTreatyIn({
  onBuka,
  pesan,
  status,
  onStatus,
}: {
  onBuka: (id: string) => void
  pesan?: string
  /** Switch In Progress / Resolved - dipegang `rute.tsx`; tanpa itu dipegang lokal (uji render). */
  status?: StatusPortal
  onStatus?: (s: StatusPortal) => void
}) {
  const [cari, setCari] = useState('')
  const [kueri, setKueri] = useState('')
  // refresh thisSection: muat ulang walau saringannya sama
  const [ketuk, setKetuk] = useState(0)
  const [statusLokal, setStatusLokal] = useState<StatusPortal>(STATUS_PORTAL[0])
  const aktif = status ?? statusLokal
  const selesai = aktif === STATUS_PORTAL[1]
  const { data: baris, galat } = useAmbil(() => daftarKasus(kueri, selesai), [kueri, ketuk, selesai])
  const [buat, setBuat] = useState(false)
  const { data: acuan } = useAmbil(() => (buat ? ambilAcuan() : Promise.resolve(null)), [buat])
  const { data: hak } = useAmbil(() => ambilHak(), [])
  const [copyOld, setCopyOld] = useState(false)

  const muatUlang = (teks: string) => {
    setKueri(teks)
    setKetuk((k) => k + 1)
  }

  return (
    <div className="inbox edmt__akar">
      <KepalaPortal copyOld={tampilCopyOld(hak)} onCopyOld={() => setCopyOld(true)} onBuat={() => setBuat(true)} />
      {pesan && <div className="alert alert--ok">{pesan}</div>}
      <StripTab tab={STATUS_PORTAL} aktif={aktif} onPilih={onStatus ?? setStatusLokal} />
      <div className="edmt__aksi edmt__saring-baris">
        {/* satu kapsul: ikon cari, isian, pengosong, tombol Filter (pola NB) */}
        <form
          className="edmt__saring"
          role="search"
          onSubmit={(e) => {
            e.preventDefault()
            muatUlang(cari.trim())
          }}
        >
          <svg className="edmt__saring-ikon" viewBox="0 0 20 20" aria-hidden="true">
            <circle cx="9" cy="9" r="6" />
            <path d="m13.5 13.5 4 4" />
          </svg>
          <input
            className="edmt__saring-isian"
            aria-label={PORTAL.filter}
            placeholder={PORTAL.placeholder}
            value={cari}
            onChange={(e) => setCari(e.target.value)}
            // tombol Filter keyboard Esc -> setValue .FilterTermForEndorsement "" -> refresh thisSection
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                setCari('')
                muatUlang('')
              }
            }}
          />
          <button
            type="button"
            className="edmt__saring-hapus"
            aria-label={PORTAL.bersihkan}
            disabled={cari === ''}
            onClick={() => setCari('')}
          >
            ×
          </button>
          <button type="submit" className="edmt__saring-tombol">
            {TOMBOL.filter}
          </button>
        </form>
        {/* S119: click -> refresh thisSection (isian saring ikut terkirim) */}
        <button type="button" className="btn" title={TOMBOL.refreshTooltip} onClick={() => muatUlang(cari.trim())}>
          {TOMBOL.refresh}
        </button>
      </div>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && baris.length === 0 && <Kosong pesan={PORTAL.kosong} petunjuk={PORTAL.kosongPetunjuk} />}
      {baris !== null && baris.length > 0 && <TabelPortal baris={baris} onBuka={onBuka} selesai={selesai} />}

      {buat && (
        <BuatEDM
          opsiJenis={(acuan?.jenisEdm ?? []).map((p) => ({ value: p.nilai, label: p.label }))}
          onTutup={() => setBuat(false)}
          onBuka={(id) => {
            setBuat(false)
            onBuka(id)
          }}
        />
      )}
      {copyOld && (
        <DialogCopyOld
          onTutup={(adaYangDisalin) => {
            setCopyOld(false)
            if (adaYangDisalin) muatUlang(kueri)
          }}
        />
      )}
    </div>
  )
}
