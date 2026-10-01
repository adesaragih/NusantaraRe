// Layar kasus Endorsement Life - FlowAction `InputEDMLife` (b166) → `Section/InputEDMLife.xml`.
//
// Kepala seluruhnya baca-saja (sel `ro`); grid peserta b11899 (EdmType 1) / b17500 (EdmType 3).

import { useCallback, useEffect, useState, type ReactNode } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilKasus, ambilPeserta, type HalamanEDM, type KasusEDM, type PesertaEDM } from '../api'
import PolisLama from '../components/PolisLama'
import RincianPeserta from '../components/RincianPeserta'
import { BUAT_EDM, GRID_EDM, KASUS_EDM, POLIS_LAMA_EDM, UMUM_EDM } from '../labels'
import { UKURAN_HALAMAN_EDM, labelEdmType, labelTypeCeding, sel, selAngka, selTanggal } from '../tampilan'
import '../endorsementlife.css'

/** Kolom grid peserta, urut korpus; `holder` false = grid Batal b17500 (tanpa `POLICY HOLDER`). */
export function kolomGrid(denganHolder: boolean): ReadonlyArray<readonly [label: string, kolom: string, jenis: 't' | 'd' | 'n']> {
  const semua: Array<readonly [string, string, 't' | 'd' | 'n']> = [
    [GRID_EDM.policyNo, 'POLICY_NO', 't'],
    [GRID_EDM.policyHolder, 'POLICY_HOLDER', 't'],
    [GRID_EDM.certificateNo, 'CERTIFICATE_NO', 't'],
    [GRID_EDM.nameOfInsured, 'NAME_OF_INSURED', 't'],
    [GRID_EDM.sex, 'SEX', 't'],
    [GRID_EDM.dateOfBirth, 'DOB', 'd'],
    [GRID_EDM.entryAge, 'ENTRY_AGE', 'n'],
    [GRID_EDM.plan, 'PLAN', 't'],
    [GRID_EDM.beginDate, 'BEGIN_DATE', 'd'],
    [GRID_EDM.effectiveDate, 'EFFECTIVE_DATE', 'd'],
    [GRID_EDM.expiredDate, 'EXPIRED_DATE', 'd'],
  ]
  return denganHolder ? semua : semua.filter(([, k]) => k !== 'POLICY_HOLDER')
}

function isiSel(p: PesertaEDM, kolom: string, jenis: 't' | 'd' | 'n'): string {
  const v = p.nilai[kolom] ?? ''
  if (jenis === 'd') return selTanggal(v)
  if (jenis === 'n') return selAngka(v)
  return sel(v)
}

/**
 * Satu baris grid peserta; `Details` membuka rincian di bawahnya (`pyEditingMode` `expandPane`,
 * `pyEditAction` `PL_DetailAction`).
 */
function Baris({
  p,
  kolom,
  terbuka,
  onBuka,
  children,
}: {
  p: PesertaEDM
  kolom: ReturnType<typeof kolomGrid>
  terbuka: boolean
  onBuka: () => void
  children: ReactNode
}) {
  return (
    <>
      <tr>
        {kolom.map(([label, kol, jenis]) => (
          <td key={label}>{isiSel(p, kol, jenis)}</td>
        ))}
        <td>
          <button type="button" className="btn btn--ghost btn--sm" aria-expanded={terbuka} onClick={onBuka}>
            {UMUM_EDM.rinci}
          </button>
        </td>
      </tr>
      {terbuka && (
        <tr className="edm-baris-rinci">
          <td colSpan={kolom.length + 1}>{children}</td>
        </tr>
      )}
    </>
  )
}

export default function InputEDMLife({ kasusId, onTutup }: { kasusId: string; onTutup: () => void }) {
  const [kasus, setKasus] = useState<KasusEDM | null>(null)
  const [peserta, setPeserta] = useState<HalamanEDM<PesertaEDM> | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [buka, setBuka] = useState('')
  const [polisLama, setPolisLama] = useState(false)

  const muatKasus = useCallback(async () => {
    setGalat(null)
    try {
      setKasus(await ambilKasus(kasusId))
    } catch (e) {
      setGalat(e)
    }
  }, [kasusId])

  const muatPeserta = useCallback(
    async (h: number) => {
      try {
        setPeserta(await ambilPeserta(kasusId, h))
      } catch (e) {
        setGalat(e)
      }
    },
    [kasusId],
  )

  useEffect(() => {
    void muatKasus()
  }, [muatKasus])
  useEffect(() => {
    void muatPeserta(halaman)
  }, [muatPeserta, halaman])

  if (galat !== null) return <Gagal galat={galat} />
  if (kasus === null) return <Memuat pesan={UMUM_EDM.memuat} />

  const k = kasus.kepala
  const medan: ReadonlyArray<readonly [string, string]> = [
    [BUAT_EDM.policyNo, kasus.policyNo],
    [KASUS_EDM.productName, k.PRODUCT_NAME ?? ''],
    [KASUS_EDM.productNameId, k.PRODUCT_NAME_ID ?? ''],
    [KASUS_EDM.type, k.TYPE ?? ''],
    [KASUS_EDM.reinsuranceSystem, labelTypeCeding(k.TYPE_CEDING ?? '')],
    [KASUS_EDM.classOfBusiness, k.BUSINESS_NAME ?? ''],
    [KASUS_EDM.sob, k.SOB_NAME ?? ''],
    [KASUS_EDM.policyHolder, k.POLICY_HOLDER_NAME ?? ''],
    [KASUS_EDM.premiumMethod, k.PRO_RATE_TYPE ?? ''],
    [KASUS_EDM.ceding, k.CEDING_CO_NAME ?? ''],
    [KASUS_EDM.marketingOfficer, k.MARKETING_NAME ?? ''],
    [KASUS_EDM.edmType, labelEdmType(kasus.edmType)],
    [KASUS_EDM.description, kasus.description],
    [BUAT_EDM.edmDate, selTanggal(kasus.edmDate)],
  ]
  const kolom = kolomGrid(kasus.edmType !== '3')

  return (
    <section className="panel edm-kasus">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{KASUS_EDM.judul}</h2>
        <div className="edm-aksi">
          {/* `View Old Policy` - dipindah dari `ShowLifePremiumSummary_EDM` (R02); popup sesuai `.Type`. */}
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setPolisLama(true)}>
            {POLIS_LAMA_EDM.viewOldPolicy}
          </button>
          <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
            {UMUM_EDM.kembali}
          </button>
        </div>
      </header>
      <p className="edm-catatan">{UMUM_EDM.terkunci}</p>
      <dl className="edm-kepala">
        {medan.map(([label, nilai]) => (
          <div key={label} className="edm-kepala__medan">
            <dt>{label}</dt>
            <dd>{sel(nilai)}</dd>
          </div>
        ))}
      </dl>
      {peserta !== null && peserta.baris.length === 0 && <Kosong pesan={UMUM_EDM.kosong} />}
      {peserta !== null && peserta.baris.length > 0 && (
        <div className="edm-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {kolom.map(([label]) => (
                  <th key={label}>{label}</th>
                ))}
                <th aria-label={UMUM_EDM.rinci} />
              </tr>
            </thead>
            <tbody>
              {peserta.baris.map((p) => (
                <Baris key={p.id} p={p} kolom={kolom} terbuka={buka === p.id} onBuka={() => setBuka(buka === p.id ? '' : p.id)}>
                  <RincianPeserta kasusId={kasusId} pesertaId={p.id} tipe={k.TYPE ?? ''} />
                </Baris>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {peserta !== null && peserta.total > 0 && (
        <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_EDM} total={peserta.total} onPindah={setHalaman} />
      )}
      {polisLama && <PolisLama kasusId={kasusId} tipe={k.TYPE ?? ''} onTutup={() => setPolisLama(false)} />}
    </section>
  )
}
