// Kotak masuk Endorsement Life - `Section/InboxEndorsementLife.xml` (harness portal `InboxEndorsementLife`).
//
// Grid b8284 dari RD `InboxEDMLife`: kasus yang belum `Resolved-Completed` dan belum `Resolved-Rejected`,
// terbaru dahulu. Tautan `Case ID` (b10202 → `openAssignment` b10349) membuka layar kasus - sejak 08-10-2026
// (keputusan work owner) SELURUH baris yang membukanya, seperti kotak masuk PremiumList Life (`inbox__baris`).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilInbox, type BarisInbox, type HalamanEDM } from '../api'
import { INBOX_EDM, UMUM_EDM } from '../labels'
import { kasusEDMBisaDibuka, UKURAN_HALAMAN_EDM, labelEdmType, sel, selTanggal, selWorkStatus } from '../tampilan'
import '../endorsementlife.css'

/**
 * Kepala grid VERBATIM. ⚠️ Urutan TAMPILAN menyimpang dari korpus (b8418 … b10052, `Create Date` b9644 di
 * antara `Marketing Name` dan `Create Operator`): `Create Date` dipindah tepat sesudah `Case ID` -
 * keputusan work owner 07-10-2026. Juga 07-10-2026: kolom SOB, Policy Holder, Marketing Name, EDM Type Batal
 * DIHAPUS, `Endorsement No` diisi `NO_ENDORS` (kosong selama kasus terbuka), bukan nomor polis (R25), dan
 * kepala `Status` bernama `Work Status` (isi tetap `STATUSS`).
 * Sel `<td>` mengikuti urutan ini.
 */
export const KOLOM_INBOX: readonly string[] = [
  INBOX_EDM.kolomCaseId,
  INBOX_EDM.kolomCreateDate,
  INBOX_EDM.kolomEndorsementNo,
  INBOX_EDM.kolomType,
  INBOX_EDM.kolomEdmType,
  INBOX_EDM.kolomPolicyNo,
  INBOX_EDM.kolomCeding,
  INBOX_EDM.kolomCreateOperator,
  INBOX_EDM.kolomStatus,
]

/** Tombol yang membuka kasus dari baris yang terfokus. */
export function tombolBukaKasus(key: string): boolean {
  return key === 'Enter' || key === ' '
}

/**
 * Satu baris kotak masuk - SATU pembuka: baris itu sendiri (klik, atau Enter/Space saat terfokus), keputusan work
 * owner 08-10-2026. `Case ID` teks biasa (bukan tombol, tanpa garis bawah sejak 08-10-2026), jadi tidak ada dua pembuka bertumpuk
 * dan `onBuka` tidak pernah terpanggil dua kali.
 */
export function BarisKotakMasuk({ b, onBuka }: { b: BarisInbox; onBuka: (id: string) => void }) {
  // Kasus tertutup TETAP terdaftar tetapi tidak dapat dibuka (keputusan work owner 09-10-2026, pola PremiumList):
  // abu, kursor biasa, tanpa onClick/onKeyDown/tabIndex.
  if (!kasusEDMBisaDibuka(b.status)) {
    return (
      <tr className="edm-inbox__tutup" title={INBOX_EDM.tutup}>
        <SelBaris b={b} />
      </tr>
    )
  }
  return (
    <tr
      className="inbox__baris edm-baris-kasus"
      tabIndex={0}
      aria-label={`${INBOX_EDM.kolomCaseId} ${b.caseId}`}
      title={`${INBOX_EDM.kolomCaseId} ${b.caseId}`}
      onClick={() => onBuka(b.caseId)}
      onKeyDown={(e) => {
        if (!tombolBukaKasus(e.key)) return
        e.preventDefault()
        onBuka(b.caseId)
      }}
    >
      <SelBaris b={b} />
    </tr>
  )
}

/** Sembilan sel satu baris kotak masuk (terbuka maupun tertutup). */
function SelBaris({ b }: { b: BarisInbox }) {
  return (
    <>
      {/* Case ID teks biasa tanpa garis bawah, seperti kotak masuk PremiumList (keputusan work owner 08-10-2026). */}
      <td>{b.caseId}</td>
      {/* dd/mm/yyyy tanpa jam (07-10-2026); urutan data tetap TGL_INPUT DESC. */}
      <td>{selTanggal(b.createDate)}</td>
      <td>{sel(b.endorsementNo)}</td>
      <td>{sel(b.type)}</td>
      <td>{labelEdmType(b.edmType)}</td>
      <td>{sel(b.policyNo)}</td>
      <td>{sel(b.ceding)}</td>
      <td>{sel(b.createOperator)}</td>
      <td>{selWorkStatus(b.status)}</td>
    </>
  )
}

export default function InboxEndorsementLife({ onBuka, onBuat }: { onBuka: (id: string) => void; onBuat: () => void }) {
  const [halaman, setHalaman] = useState(1)
  const [isi, setIsi] = useState<HalamanEDM<BarisInbox> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  const muat = useCallback(async (h: number) => {
    setGalat(null)
    try {
      setIsi(await ambilInbox(h))
    } catch (e) {
      setGalat(e)
    }
  }, [])

  useEffect(() => {
    void muat(halaman)
  }, [muat, halaman])

  return (
    <section className="inbox edm-inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{INBOX_EDM.judul}</h2>
        {/* `Create Addendum` b6620 → `showHarness` b6984 `EndorsmentLife_harnes` (`pyTarget=current`). */}
        <button type="button" className="btn btn--primary btn--sm" onClick={onBuat}>
          {INBOX_EDM.createAddendum}
        </button>
      </header>
      {galat !== null && <Gagal galat={galat} />}
      {isi === null && galat === null && <Memuat pesan={UMUM_EDM.memuat} />}
      {isi !== null && isi.baris.length === 0 && <Kosong pesan={UMUM_EDM.kosong} />}
      {isi !== null && isi.baris.length > 0 && (
        <div className="edm-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {KOLOM_INBOX.map((k) => (
                  <th key={k}>{k}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {isi.baris.map((b) => (
                <BarisKotakMasuk key={b.caseId} b={b} onBuka={onBuka} />
              ))}
            </tbody>
          </table>
        </div>
      )}
      {isi !== null && isi.total > 0 && (
        <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_EDM} total={isi.total} onPindah={setHalaman} />
      )}
    </section>
  )
}
