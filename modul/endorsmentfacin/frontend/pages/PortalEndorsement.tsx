// Halaman depan Endorsment Fac In - portal endorsement (irisan I0, keputusan work owner 11-10-2026 W-EDM-01).
//
// Port harness `SFAPortalEndorsement` (`D:\migrasi\RNM\Endorsment Fac In\Harness\SFAPortalEndorsement.xml`,
// DATA-PORTAL!SFAPORTALENDORSEMENT) beserta section `SFAPortal_EndorsementHeader`: judul `Addendum` (L1757) dan tombol
// `Create Addendum` (L1053-1072, `createWork` kelas `ASM-SFAGIS-Work-Endorsement`).
//
// ⚠️ `[terverifikasi]` Harness ini TIDAK memuat grid: setiap `pyGridRDName` / `pyGridRDClass` kosong, dan sumber daftar
// kasus endorsement (tombol `Refresh` bertooltip "Refresh EDM Grid") tidak ada di korpus - `[pertanyaan terbuka]`. Karena
// itu kotak saring `Policy Number` + `Filter` dan `Refresh` (yang hanya memuat ulang grid itu) belum dirender: tanpa
// grid keduanya tidak melakukan apa pun.

import { Kosong } from '../../../../inti/frontend/components/ui/dasar'
import { PERIKSA, PORTAL } from '../labels'

/** `pesan` - konfirmasi sesudah kasus `EDM-` lahir (irisan I3); `onTutupPesan` menghapusnya. */
export default function PortalEndorsement({
  onBuat,
  pesan = '',
  onTutupPesan,
}: {
  onBuat: () => void
  pesan?: string
  onTutupPesan?: () => void
}) {
  return (
    <section className="inbox edmf-tema edmf-portal">
      {pesan !== '' && (
        <div className="alert alert--ok edmf-portal__pesan" role="status">
          <span>{pesan}</span>
          {onTutupPesan && (
            <button type="button" className="btn edmf-tombol btn--ghost btn--sm" aria-label={PERIKSA.tutup} onClick={onTutupPesan}>
              ✕
            </button>
          )}
        </div>
      )}
      <header className="inbox__kepala edmf-portal__kepala">
        <h2 className="inbox__judul edmf-tema__judul">{PORTAL.judul}</h2>
        <button type="button" className="btn btn--primary edmf-tombol edmf-tombol--utama" onClick={onBuat}>
          {PORTAL.buat}
        </button>
      </header>
      <Kosong pesan={PORTAL.tanpaDaftar} petunjuk={PORTAL.petunjukTanpaDaftar} />
    </section>
  )
}
