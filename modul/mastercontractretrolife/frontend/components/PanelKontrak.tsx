// Panel kontrak dan batas proteksi - tiket 02/04. Harness `InboxRetroLimitReinsurers` → section
// `InputRetroLimitReinsurers.xml`, dibuka tombol baris `ReinsType` halaman awal
// (`InputRetrocessionLife.xml` b11988, `showHarness` popup b12720).
//
// Kepala `Reins Type` b601 + `ID Treaty Year` b1058. Grid RD `BrowseTreatyContract_Life_RD` (param
// `IDTREATYYEAR`, urut `.ID ASC` - server). Form (wadah b1566 `DATASHOW = 1`): `Inputor` b2070/b2290,
// `Modified Date` b2499, `REINS TYPE` dropdown b3242 (nilai `.ID`, tampil `.Note`), `TREATY START`
// b3604 / `TREATY END` b3894 BACA-SAJA, empat batas b4184-b4936, `Save` b5773, `Cancel` b6066. Baris:
// `Edit` b12621, `Business List` b12999, `Reinsurer List` b14034, `Delete` b15161. `Add` b8596.
//
// ⛔ K4: `TREATY START`/`TREATY END` tidak pernah diketik - server menyalinnya dari tahun treaty induk
// (`SetValueRetroLimit_TreatyYearLife` b485/b507); layar menampilkan tanggal tahun itu.

import { useCallback, useEffect, useState } from 'react'

import {
  ambilJenisReasuransi,
  ambilKontrak,
  simpanKontrak,
  type JawabanKontrak,
  type JenisReasuransi,
  type Kontrak,
  type KontrakMasuk,
  type TahunTreaty,
} from '../api'
import { KONTRAK_MCRL, TAHUN_MCRL, UMUM_MCRL } from '../labels'
import { operatorKini, sel, selAngka, selTanggal, selWaktu, waktuKini } from '../tampilan'
import { Field, Gagal, Kosong, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { KepalaPanel } from './Bingkai'
import KonfirmasiHapus from './KonfirmasiHapus'
import PanelBusiness from './PanelBusiness'
import PanelReinsurer from './PanelReinsurer'
import { useHapus } from './useHapus'

/** Isian form - `InputTreatyContract.*`. */
export interface FormKontrak {
  id: string
  reinsTypeId: string
  /** `MINIMUM LIMIT (IDR)` → `B_IDR`. */
  bIdr: string
  /** `MAXIMUM LIMIT (IDR)` → `IDR`. */
  idr: string
  /** `MINIMUM LIMIT (USD)` → `B_USD`. */
  bUsd: string
  /** `MAXIMUM LIMIT (USD)` → `USD`. */
  usd: string
  /** `Inputor` / `Modified Date` (ro). */
  userId: string
  tglUpdate: string
}

/**
 * `NewInputTreatyLimit_Life` b402: kosongkan `ID, REINSTYPEID, REINSTYPENAME, IDR, USD, B_IDR, B_USD`,
 * `TGLUPDATE ← @CurrentDateTime()`, `USERID ← OperatorID`.
 */
export function formKontrakBaru(operator: string, kini: string): FormKontrak {
  return { id: '', reinsTypeId: '', bIdr: '', idr: '', bUsd: '', usd: '', userId: operator, tglUpdate: kini }
}

/**
 * `SetRetroListLife_Act` b597: `USERID ← OperatorID`, `TGLUPDATE ← @CurrentDateTime()`, salin
 * `REINSTYPEID`, `IDR`, `USD`, `B_IDR`, `B_USD`, `ID`.
 */
export function formKontrakDari(k: Kontrak, operator: string, kini: string): FormKontrak {
  return { id: k.id, reinsTypeId: k.reinsTypeId, bIdr: k.bIdr, idr: k.idr, bUsd: k.bUsd, usd: k.usd, userId: operator, tglUpdate: kini }
}

/** Badan simpan; batas dikirim apa adanya (di-trim) - server menormalkan dan memeriksa urutan layer. */
export function keKontrakMasuk(f: FormKontrak): KontrakMasuk {
  return { id: f.id, reinsTypeId: f.reinsTypeId, idr: f.idr.trim(), usd: f.usd.trim(), bIdr: f.bIdr.trim(), bUsd: f.bUsd.trim() }
}

/** Opsi dropdown `REINS TYPE` - nilai `.ID`, tampil `.Note` (`pyPrompt`). */
export function opsiJenis(daftar: readonly JenisReasuransi[]): { value: string; label: string }[] {
  return daftar.map((j) => ({ value: j.id, label: j.note.trim() === '' ? j.id : j.note }))
}

/**
 * Jenis yang BELUM dipakai kontrak lain di tahun treaty ini - satu REINS TYPE
 * satu kontrak per tahun (keputusan work owner 04-10-2026). Jenis milik kontrak
 * yang sedang diubah (`idSendiri`) tetap tersedia. Server menolak hal yang sama
 * (`PesanJenisGanda`).
 */
export function jenisTersedia(
  daftarJenis: readonly JenisReasuransi[],
  kontrak: readonly Kontrak[],
  idSendiri: string,
): JenisReasuransi[] {
  const terpakai = new Set(kontrak.filter((k) => k.id !== idSendiri).map((k) => k.reinsTypeId.trim()))
  return daftarJenis.filter((j) => !terpakai.has(j.id.trim()))
}

/** Popup yang dibuka dari satu baris kontrak - `Reinsurer List` atau `Business List`. */
type Anak = { jenis: 'reinsurer' | 'business'; kontrak: Kontrak }

export default function PanelKontrak({ tahun, onTutup }: { tahun: TahunTreaty; onTutup: () => void }) {
  const [jawab, setJawab] = useState<JawabanKontrak | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [form, setForm] = useState<FormKontrak | null>(null)
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [pesan, setPesan] = useState<string | null>(null)
  const [jenis, setJenis] = useState<JenisReasuransi[]>([])
  const [galatJenis, setGalatJenis] = useState<unknown>(null)
  // Setiap `Add`/`Edit` memuat ulang master jenis - galat sebelumnya tidak hilang tanpa dicoba lagi.
  const [muatJenisKe, setMuatJenisKe] = useState(0)
  const [anak, setAnak] = useState<Anak | null>(null)

  const muat = useCallback(async () => {
    try {
      const j = await ambilKontrak(tahun.id)
      setJawab(j)
      setGalat(null)
    } catch (e) {
      setGalat(e)
    }
  }, [tahun.id])

  useEffect(() => {
    void muat()
  }, [muat])

  // Master jenis dimuat saat form dibuka - kosong/tak terbaca dinyatakan (503 berkalimat).
  const formTerbuka = form !== null
  useEffect(() => {
    if (!formTerbuka) return
    let hidup = true
    ambilJenisReasuransi()
      .then((d) => {
        if (hidup) setJenis(d.daftar)
      })
      .catch((e: unknown) => {
        if (hidup) setGalatJenis(e)
      })
    return () => {
      hidup = false
    }
  }, [formTerbuka, muatJenisKe])

  const hapusan = useHapus('kontrak', async (id, p) => {
    if (form?.id === id) setForm(null)
    setPesan(p)
    await muat()
  })

  function buka(f: FormKontrak): void {
    setGalatForm(null)
    setGalatJenis(null)
    setMuatJenisKe((k) => k + 1)
    setPesan(null)
    setForm(f)
  }

  function ubah(k: 'bIdr' | 'idr' | 'bUsd' | 'usd' | 'reinsTypeId') {
    return (v: string) => {
      setForm((f) => (f === null ? f : { ...f, [k]: v }))
    }
  }

  async function simpan(): Promise<void> {
    if (form === null || menyimpan) return
    setMenyimpan(true)
    setGalatForm(null)
    try {
      await simpanKontrak(tahun.id, keKontrakMasuk(form))
      // `SaveTreatyLimit_Act` langkah 6 b1275: `DATASHOW3 = ""`, `DATASHOW = ""` - form tertutup.
      setForm(null)
      await muat()
    } catch (e) {
      setGalatForm(e)
    } finally {
      setMenyimpan(false)
    }
  }

  if (anak !== null) {
    const tutup = () => {
      setAnak(null)
    }
    return anak.jenis === 'reinsurer' ? (
      <PanelReinsurer key={anak.kontrak.id} kontrak={anak.kontrak} onTutup={tutup} />
    ) : (
      <PanelBusiness key={anak.kontrak.id} kontrak={anak.kontrak} onTutup={tutup} />
    )
  }

  const induk = jawab?.tahun ?? tahun
  const daftar = jawab?.daftar ?? []

  return (
    <section className="panel">
      {/* ID Treaty Year PALING DEPAN, lalu Underwriting Year dan Transaction Year
          (keputusan work owner 04-10-2026) - label sama dengan halaman depan. */}
      <KepalaPanel
        judul={KONTRAK_MCRL.judul}
        medan={[
          [KONTRAK_MCRL.idTreatyYear, induk.id],
          [TAHUN_MCRL.formUnderwritingYear, induk.underwritingYear],
          [TAHUN_MCRL.formTransactionYear, induk.treatyYear],
        ]}
        onTutup={onTutup}
      />

      {form !== null && (
        <div className="panel">
          {galatForm !== null && <Gagal galat={galatForm} />}
          {galatJenis !== null && <Gagal galat={galatJenis} />}
          <div className="form-grid">
            <Pilih
              label={KONTRAK_MCRL.formReinsType}
              value={form.reinsTypeId}
              onChange={ubah('reinsTypeId')}
              opsi={opsiJenis(jenisTersedia(jenis, daftar, form.id))}
              required
            />
            {/*
              TREATY START / TREATY END SENGAJA tidak ditampilkan di form (keputusan
              work owner 04-10-2026). Nilainya tetap disalin server dari tahun treaty
              (`KontrakMasuk.keModel`) dan tetap tampil di kolom grid.
            */}
            <Field label={KONTRAK_MCRL.formMinIdr} value={form.bIdr} onChange={ubah('bIdr')} required />
            <Field label={KONTRAK_MCRL.formMaxIdr} value={form.idr} onChange={ubah('idr')} required />
            <Field label={KONTRAK_MCRL.formMinUsd} value={form.bUsd} onChange={ubah('bUsd')} required />
            <Field label={KONTRAK_MCRL.formMaxUsd} value={form.usd} onChange={ubah('usd')} required />
            {/* Modified Date dan Inputor di PALING AKHIR (keputusan work owner 04-10-2026). */}
            <Field label={KONTRAK_MCRL.formModifiedDate} value={selWaktu(form.tglUpdate)} onChange={() => undefined} readOnly />
            <Field label={KONTRAK_MCRL.formInputor} value={form.userId} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={menyimpan} onClick={() => void simpan()}>
              {KONTRAK_MCRL.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalatForm(null)
              }}
            >
              {KONTRAK_MCRL.cancel}
            </button>
          </div>
        </div>
      )}

      {pesan !== null && <p role="status">{pesan}</p>}

      <div className="aksi-baris mcrl-aksi-grid">
        <button type="button" className="btn btn--primary" onClick={() => buka(formKontrakBaru(operatorKini(), waktuKini(new Date())))}>
          {KONTRAK_MCRL.add}
        </button>
      </div>
      {jawab === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {jawab !== null && daftar.length === 0 && <Kosong pesan={UMUM_MCRL.kosong} />}
      {daftar.length > 0 && (
        <div className="mcrl-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{KONTRAK_MCRL.kolomId}</th>
                <th>{KONTRAK_MCRL.kolomReinsType}</th>
                <th>{KONTRAK_MCRL.kolomTreatyStart}</th>
                <th>{KONTRAK_MCRL.kolomTreatyEnd}</th>
                <th>{KONTRAK_MCRL.kolomMinIdr}</th>
                <th>{KONTRAK_MCRL.kolomMaxIdr}</th>
                <th>{KONTRAK_MCRL.kolomMinUsd}</th>
                <th>{KONTRAK_MCRL.kolomMaxUsd}</th>
                <th className="table__actions" />
              </tr>
            </thead>
            <tbody>
              {daftar.map((k) => (
                <tr key={k.id} className="inbox__baris">
                  <td>{sel(k.id)}</td>
                  <td>{sel(k.reinsTypeName)}</td>
                  <td>{selTanggal(k.treatyStartDate)}</td>
                  <td>{selTanggal(k.treatyEndDate)}</td>
                  <td>{selAngka(k.bIdr)}</td>
                  <td>{selAngka(k.idr)}</td>
                  <td>{selAngka(k.bUsd)}</td>
                  <td>{selAngka(k.usd)}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formKontrakDari(k, operatorKini(), waktuKini(new Date())))}>
                      {KONTRAK_MCRL.edit}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setAnak({ jenis: 'business', kontrak: k })
                      }}
                    >
                      {KONTRAK_MCRL.businessList}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setAnak({ jenis: 'reinsurer', kontrak: k })
                      }}
                    >
                      {KONTRAK_MCRL.reinsurerList}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setPesan(null)
                        hapusan.minta(k.id, k.reinsTypeName || k.id)
                      }}
                    >
                      {KONTRAK_MCRL.delete}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {hapusan.konfirmasi !== null && (
        <KonfirmasiHapus
          judul={KONTRAK_MCRL.delete}
          labelBatal={KONTRAK_MCRL.cancel}
          jenis="kontrak"
          nama={hapusan.konfirmasi.nama}
          dampak={hapusan.konfirmasi.dampak}
          galat={hapusan.konfirmasi.galat}
          sibuk={hapusan.sibuk}
          onYa={() => void hapusan.ya()}
          onBatal={hapusan.batal}
        />
      )}
    </section>
  )
}
