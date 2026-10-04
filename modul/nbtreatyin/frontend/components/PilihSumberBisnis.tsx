// Popup Harness `SOB` -> `Section/SourceHierarki` (pemilih Source Of Business,
// hanya untuk ClaimType 'XOL Retro'). Isi grid = RD `BrowseAgentHierarkiList_RD`
// sebagaimana efektif di NB: daftar akar (`LEADER0 IS NULL`), satu kolom
// "Source of Business Name" (`.ClientName`).
//
//   klik nama baris  = pyRowEditing masterDetail -> flow action `AgentSourceBizDetails`
//                      -> pra-proses `SearchHierarkiSourceBizAgent_PostDT` (backend
//                      `POST .../pilih-sumber-bisnis`, TANPA simpan - F4); simpul
//                      beranak MENGOSONGKAN sumber bisnis - apa adanya
//   tombol Choose    = tampil `.ChildCount = 0`; `window.close` + `opener.location.reload`
//                      -> popup ditutup
//
// ⭐ F4 (IKUTI XML): PostDT hanya menulis clipboard (`pyWorkPage.Quotation.*`, nol
// Obj-Save). Hasilnya DIPEGANG state layar (`pegangSumberBisnis`) dan ikut terkirim
// pada Save/Submit (dan refresh) - disimpan saat itu, bukan saat klik. XML tidak
// menjalankan refresh berhitung sesudah klik/Choose: layar pun tidak.
//
// ⛔ Tidak dibangun, dengan bukti XML:
//   - autocomplete "Search" (`SearchSOB.CARI1`, RD `BrowseAgentNusaRe_RD`): nilainya
//     dibaca NOL rule - `btnSOB_DT`/`btnCedingCO_DT` hanya mengosongkannya, grid tidak
//     memakainya;
//   - perluasan simpul pohon: `AgentSourceBizTreatyIn_Act` hanya mengisi `Param.Leader`
//     bila `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="SOB"` - properti yang
//     ditulis nol rule - sehingga setiap simpul mengembalikan daftar akar yang sama;
//   - aktivitas Choose `SearchHierarkiSourceBizAgentTreatyIn_Act`: menulis
//     `pyWorkPage.OfferTreatyIn.QuotationData.*` (+ RDB `BrowseClientEmail_SQL`) yang
//     dibaca nol rule NB.

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbil } from '../ambil'
import {
  KLAIM_XOL_RETRO,
  POLIS,
  daftarSumberBisnis,
  nilai,
  setel,
  type BarisAgen,
  type Halaman,
  type HasilSumberBisnis,
} from '../api'
import { JUDUL, KOLOM_SOB, TOMBOL } from '../labels'

/** pyVisible tombol `Select Source Of Business`: `.ClaimType = 'XOL Retro'`. */
export function tampilTombolSOB(h: Halaman): boolean {
  return nilai(h, POLIS + 'ClaimType') === KLAIM_XOL_RETRO
}

/** Jalur halaman tiap medan yang ditulis `SearchHierarkiSourceBizAgent_PostDT`
 *  langkah 1.1, 1.2, 4, 5. */
export const JALUR_SUMBER_BISNIS: Record<keyof HasilSumberBisnis, string> = {
  sourceOfBusiness: 'Quotation.SourceOfBusiness',
  sobName: 'Quotation.SobName',
  sobLeader0: 'Quotation.SobLeader0',
  sobLeader1: 'Quotation.SobLeader1',
}

/** Halaman layar dengan hasil klik dipegang: HANYA keempat jalur PostDT ditulis;
 *  medan lain (termasuk `PolicyTreatyIn.*`) tidak disentuh. */
export function pegangSumberBisnis(h: Halaman, hasil: HasilSumberBisnis): Halaman {
  let x = h
  for (const [k, j] of Object.entries(JALUR_SUMBER_BISNIS) as [keyof HasilSumberBisnis, string][]) {
    x = setel(x, j, hasil[k] ?? '')
  }
  return x
}

/** pyVisible tombol `Choose`: `.ChildCount = 0` (kosong = 0). */
export function tampilChoose(b: BarisAgen): boolean {
  return Number(b.childCount.trim() || '0') === 0
}

export default function PilihSumberBisnis({
  terpilih,
  sibuk,
  onPilih,
  onTutup,
}: {
  /** `Quotation.SourceOfBusiness` yang dipegang layar - baris yang sedang terpilih. */
  terpilih: string
  sibuk: boolean
  /** Klik baris; `tutup` = lewat tombol Choose (popup ditutup sesudahnya). */
  onPilih: (idAgen: string, tutup: boolean) => void
  onTutup: () => void
}) {
  const { data: baris, galat } = useAmbil(daftarSumberBisnis, [])

  return (
    <Modal judul={JUDUL.sumberBisnis} onTutup={onTutup} penuh>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && baris.length === 0 && <Kosong pesan="—" />}
      {baris !== null && baris.length > 0 && (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_SOB}</th>
                <th scope="col" />
              </tr>
            </thead>
            <tbody>
              {baris.map((b, i) => (
                <tr key={`${b.id}-${i}`} aria-selected={b.id !== '' && b.id === terpilih}>
                  <td>
                    <button
                      type="button"
                      className="nbti__tautan"
                      disabled={sibuk}
                      onClick={() => onPilih(b.id, false)}
                    >
                      {b.id !== '' && b.id === terpilih ? <strong>{b.clientName}</strong> : b.clientName}
                    </button>
                  </td>
                  <td>
                    {tampilChoose(b) && (
                      <button
                        type="button"
                        className="btn btn--primary"
                        disabled={sibuk}
                        onClick={() => onPilih(b.id, true)}
                      >
                        {TOMBOL.choose}
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Modal>
  )
}
