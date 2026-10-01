// Layar kasus Endorsement Life - FlowAction `InputEDMLife` (b166) → `Section/InputEDMLife.xml`.
//
// Kepala seluruhnya baca-saja (sel `ro`); grid peserta b11899 (EdmType 1) / b17500 (EdmType 3).
// `Save` b37202 → `SetPremi_EDM` (tiket 05/06): centang `.EdmBatal` b15753 dan `DELETE ALL` b13607
// hanya di grid Perubahan Data; rekap mata uang per `.Type` (b23064 …) sesudah simpan.
// Unggah CSV (tiket 07) - wadah b8698 `.EdmType=1 && .EditInput=1`: `Upload CSV` b8973 dan
// `Add CSV Data` b10405 mati bila `.EditInput1=1` (b8965/b10403, ada baris `New`); `View Upload` b9340.
// Keputusan (tiket 08) - `ConfirmSection` tampil sesudah `Save` (b35518 `.IsJsonPolis=1`); Confirm
// menutup dengan `ConfirmSubmitEDM`, Decline kembali ke kotak masuk (`finishAssignment` b38385).

import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilKasus,
  ambilPeserta,
  simpanKasus,
  tambahCSV,
  type HalamanEDM,
  type HasilPutusanEDM,
  type KasusEDM,
  type PeriksaCSVEDM,
  type PesertaEDM,
} from '../api'
import HasilCSV from '../components/HasilCSV'
import KeputusanEDM from '../components/KeputusanEDM'
import PolisLama, { tipePopup } from '../components/PolisLama'
import RincianPeserta from '../components/RincianPeserta'
import { TabelKorpus } from '../components/TabelKorpus'
import TerimaKasih from '../components/TerimaKasih'
import UnggahCSV, { TabelPesanCSV } from '../components/UnggahCSV'
import { barisUnggahan } from '../csv'
import { MATA_UANG } from '../kolomKorpus'
import { BUAT_EDM, GRID_EDM, KASUS_EDM, POLIS_LAMA_EDM, SIMPAN_EDM, UMUM_EDM, UNGGAH_EDM } from '../labels'
import {
  PILIHAN_KOSONG,
  UKURAN_HALAMAN_EDM,
  alihBaris,
  alihSemua,
  bolehCentang,
  labelEdmType,
  labelTypeCeding,
  sel,
  selAngka,
  selTanggal,
  tercentang,
  type PilihanHapusEDM,
} from '../tampilan'
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
  centang,
  terbuka,
  onBuka,
  children,
}: {
  p: PesertaEDM
  kolom: ReturnType<typeof kolomGrid>
  /** `null` = grid Batal b17500, tanpa kolom kotak centang. */
  centang: { boleh: boolean; aktif: boolean; onAlih: () => void } | null
  terbuka: boolean
  onBuka: () => void
  children: ReactNode
}) {
  return (
    <>
      <tr>
        {centang !== null && (
          <td>
            <input
              type="checkbox"
              aria-label={UMUM_EDM.tandaiHapus}
              checked={centang.boleh ? centang.aktif : p.edmStatus === 'Delete'}
              disabled={!centang.boleh}
              onChange={centang.onAlih}
            />
          </td>
        )}
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
          <td colSpan={kolom.length + (centang === null ? 1 : 2)}>{children}</td>
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
  const [pilihan, setPilihan] = useState<PilihanHapusEDM>(PILIHAN_KOSONG)
  const [menyimpan, setMenyimpan] = useState(false)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [jendelaCSV, setJendelaCSV] = useState<'unggah' | 'hasil' | null>(null)
  const [berkasCSV, setBerkasCSV] = useState<File | null>(null)
  const [barisCSV, setBarisCSV] = useState<Array<Record<string, string>>>([])
  const [periksa, setPeriksa] = useState<PeriksaCSVEDM | null>(null)
  const [menambah, setMenambah] = useState(false)
  const [disimpanCSV, setDisimpanCSV] = useState<number | null>(null)
  const [putusan, setPutusan] = useState<HasilPutusanEDM | null>(null)
  // Berkas CSV yang berlaku sekarang - hasil baca berkas lama tidak boleh menimpa berkas yang lebih baru.
  const berkasKini = useRef<File | null>(null)

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

  const simpan = async () => {
    setMenyimpan(true)
    setGalatSimpan(null)
    try {
      await simpanKasus(kasusId, pilihan)
      setPilihan(PILIHAN_KOSONG)
      await muatKasus()
      await muatPeserta(halaman)
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setMenyimpan(false)
    }
  }

  const pilihCSV = (berkas: File, hasil: PeriksaCSVEDM) => {
    berkasKini.current = berkas
    setBerkasCSV(berkas)
    setPeriksa(hasil)
    setDisimpanCSV(null)
    setBarisCSV([])
  }

  // `View Upload`: berkas diurai di peramban baru saat popup dibuka (unggahan tanpa batas baris).
  const bukaHasilCSV = () => {
    setJendelaCSV('hasil')
    const b = berkasKini.current
    if (b === null) return
    void b.text().then((t) => {
      if (berkasKini.current === b) setBarisCSV(barisUnggahan(t))
    })
  }

  const tambah = async () => {
    if (berkasCSV === null) return
    setMenambah(true)
    setGalatSimpan(null)
    try {
      const j = await tambahCSV(kasusId, berkasCSV)
      if (j.jenis === 'ditolak') {
        setPeriksa(j.periksa)
        return
      }
      setDisimpanCSV(j.hasil.disimpan)
      setPeriksa(null)
      berkasKini.current = null
      setBerkasCSV(null)
      setBarisCSV([])
      await muatKasus()
      await muatPeserta(halaman)
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setMenambah(false)
    }
  }

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
  const denganCentang = kasus.edmType === '1'
  const terbukaKasus = kasus.status === ''
  // `DELETE ALL` mengikuti kotak centang: hidup selama kasus terbuka dan belum disimpan.
  const bolehHapus = bolehCentang({ edmStatus: 'Old' }, kasus)
  const kolomUang = MATA_UANG[tipePopup(k.TYPE ?? '')]

  return (
    <section className="panel edm-kasus">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{KASUS_EDM.judul}</h2>
        <div className="edm-aksi">
          {/* `View Old Policy` - dipindah dari `ShowLifePremiumSummary_EDM` (R02); popup sesuai `.Type`. */}
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setPolisLama(true)}>
            {POLIS_LAMA_EDM.viewOldPolicy}
          </button>
          {/* `Save` b37202: `VIS .EditInput=1` (kasus terbuka), mati bila `.IsJsonPolis=1` b37200. */}
          {terbukaKasus && (
            <button type="button" className="btn btn--sm" disabled={kasus.sudahSimpan || menyimpan} onClick={() => void simpan()}>
              {menyimpan ? UMUM_EDM.menyimpan : SIMPAN_EDM.save}
            </button>
          )}
          <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
            {UMUM_EDM.kembali}
          </button>
        </div>
      </header>
      {galatSimpan !== null && <Gagal galat={galatSimpan} />}
      <p className="edm-catatan">{UMUM_EDM.terkunci}</p>
      <dl className="edm-kepala">
        {medan.map(([label, nilai]) => (
          <div key={label} className="edm-kepala__medan">
            <dt>{label}</dt>
            <dd>{sel(nilai)}</dd>
          </div>
        ))}
      </dl>
      {denganCentang && terbukaKasus && (
        <div className="edm-aksi">
          <button type="button" className="btn btn--ghost btn--sm" disabled={kasus.csvTerkunci} onClick={() => setJendelaCSV('unggah')}>
            {UNGGAH_EDM.uploadCsv}
          </button>
          <button type="button" className="btn btn--ghost btn--sm" onClick={bukaHasilCSV}>
            {UNGGAH_EDM.viewUpload}
          </button>
          <button
            type="button"
            className="btn btn--sm"
            disabled={kasus.csvTerkunci || berkasCSV === null || menambah}
            onClick={() => void tambah()}
          >
            {menambah ? UMUM_EDM.mengunggah : UNGGAH_EDM.addCsvData}
          </button>
          {disimpanCSV !== null && (
            <span className="edm-catatan">
              {UMUM_EDM.barisDisimpan}: {disimpanCSV}
            </span>
          )}
        </div>
      )}
      {periksa !== null && jendelaCSV === null && <TabelPesanCSV hasil={periksa} />}
      {peserta !== null && peserta.baris.length === 0 && <Kosong pesan={UMUM_EDM.kosong} />}
      {peserta !== null && peserta.baris.length > 0 && (
        <div className="edm-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {denganCentang && (
                  <th>
                    <button type="button" className="btn btn--ghost btn--sm" disabled={!bolehHapus} onClick={() => setPilihan(alihSemua(pilihan))}>
                      {SIMPAN_EDM.deleteAll}
                    </button>
                  </th>
                )}
                {kolom.map(([label]) => (
                  <th key={label}>{label}</th>
                ))}
                <th aria-label={UMUM_EDM.rinci} />
              </tr>
            </thead>
            <tbody>
              {peserta.baris.map((p) => (
                <Baris
                  key={p.id}
                  p={p}
                  kolom={kolom}
                  centang={
                    denganCentang
                      ? { boleh: bolehCentang(p, kasus), aktif: tercentang(pilihan, p.id), onAlih: () => setPilihan(alihBaris(pilihan, p.id)) }
                      : null
                  }
                  terbuka={buka === p.id}
                  onBuka={() => setBuka(buka === p.id ? '' : p.id)}
                >
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
      {kasus.rekap.length > 0 && (
        <TabelKorpus kolom={kolomUang} baris={kasus.rekap} kunci={(b, i) => `${b.CURRENCY ?? ''}-${i}`} />
      )}
      {terbukaKasus && kasus.sudahSimpan && (
        <KeputusanEDM
          kasusId={kasusId}
          riwayat={kasus.riwayat}
          onSelesai={(h) => (h.noEndorsement !== '' ? setPutusan(h) : onTutup())}
        />
      )}
      {putusan !== null && <TerimaKasih hasil={putusan} onTutup={onTutup} />}
      {jendelaCSV === 'unggah' && <UnggahCSV kasusId={kasusId} onBerkas={pilihCSV} onTutup={() => setJendelaCSV(null)} />}
      {jendelaCSV === 'hasil' && <HasilCSV baris={barisCSV} onTutup={() => setJendelaCSV(null)} />}
      {polisLama && <PolisLama kasusId={kasusId} tipe={k.TYPE ?? ''} onTutup={() => setPolisLama(false)} />}
    </section>
  )
}
