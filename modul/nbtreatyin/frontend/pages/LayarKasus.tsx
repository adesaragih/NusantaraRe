// Layar realisasi satu kasus - flow action posisinya:
//   Admin   `InboxPolicyTreatyIn`  -> `GeneralPolicyTreatyIn` / `DetailPolicyTreatyIn`
//   Atasan  `DeptHeadTreatyIn_UW`  -> `GeneralDeptHeadTreatyIn_UW` / `DetailDeptHeadTreatyIn_UW`
// beserta grid `.SpreadingRiskList`, jadwal angsuran `.ListInstallment`, dan
// `ListSuggest`.
//
// ⛔ Setiap refresh berhitung dikirim ke backend (`POST .../hitung`); layar ini
// tidak menghitung apa pun. Action set sel yang memuat lebih dari satu refresh
// dikirim SEKALI sebagai `urutan` (dijalankan berurutan atas halaman yang sama).
// Medan wajib datang dari backend (`medanWajib`). Padanan setiap tombol/aksi
// dengan rule XML: `docs/alat/tombol.json`.

import { useCallback, useState } from 'react'

import { Gagal, Memuat, Modal, Panel, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbilBatal } from '../ambil'
import {
  POLIS,
  ambilAcuan,
  bukaKasus,
  daftar,
  hitung,
  kirimKasus,
  nilai,
  pilihBisnis,
  pilihSumberBisnis,
  riwayatKasus,
  setel,
  setelDaftar,
  simpanKasus,
  terbitkanNomor,
  type Acuan,
  type Baris,
  type Halaman,
  type Layar,
  type NomorPolis,
  type Riwayat,
} from '../api'
import DetailNonProp from '../components/DetailNonProp'
import InputAngka from '../components/InputAngka'
import KotakMedan from '../components/KotakMedan'
import PilihBisnis from '../components/PilihBisnis'
import PilihSumberBisnis, { tampilTombolSOB } from '../components/PilihSumberBisnis'
import {
  BAGIAN,
  JUDUL,
  JUDUL_POSISI,
  KOLOM_ANGSURAN,
  KOLOM_SPREADING,
  KOLOM_USULAN,
  KONFIRMASI_TOLAK,
  NOMOR_DIAKSEP,
  PESAN,
  PILIHAN_APPROVAL,
  PORTAL,
  TOMBOL,
} from '../labels'
import {
  MEDAN_ADMIN_UANG,
  MEDAN_ADMIN_UMUM,
  MEDAN_ATASAN_TOTAL,
  MEDAN_ATASAN_UANG,
  MEDAN_ATASAN_UMUM,
  SAJIAN_ANGSURAN,
  SAJIAN_SPREADING,
  bukanXOLRetro,
  medanTampil,
  wadahUangAdmin,
  wadahUangAtasan,
  type Aksi,
  type Medan,
} from '../medan'
import { tampilNonProp } from '../nonprop'
import { sajikan, type Sajian } from '../sajian'
import { tampilTanggalProduksi } from '../tempat'

const SPREADING = POLIS + 'SpreadingRiskList'
const ANGSURAN = POLIS + 'ListInstallment'
const USULAN = POLIS + 'SuggestList'

function opsi(p: { nilai: string; label: string }[] | null | undefined): Opsi[] {
  return (p ?? []).map((x) => ({ value: x.nilai, label: x.label }))
}

export default function LayarKasus({ id, onKembali }: { id: string; onKembali: (pesan?: string) => void }) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [h, setH] = useState<Halaman | null>(null)
  const [acuan, setAcuan] = useState<Acuan | null>(null)
  const [riwayat, setRiwayat] = useState<Riwayat[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [info, setInfo] = useState('')
  const [popupBisnis, setPopupBisnis] = useState(false)
  const [popupSOB, setPopupSOB] = useState(false)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [nomor, setNomor] = useState<NomorPolis | null>(null)

  const terima = useCallback((ly: Layar) => {
    setLayar(ly)
    setH(ly.halaman)
  }, [])

  useAmbilBatal(
    () => Promise.all([bukaKasus(id), ambilAcuan(), riwayatKasus(id)]),
    ([ly, a, r]) => {
      terima(ly)
      setAcuan(a)
      setRiwayat(r ?? [])
    },
    setGalat,
    [id, terima],
  )

  if (galat !== null && layar === null) return <Gagal galat={galat} />
  if (layar === null || h === null) return <Memuat />

  const posisi = layar.kasus.positionNote
  const admin = posisi === 'ReasTreatyInAdmin'
  const boleh = layar.bolehKerja
  const wajib = new Set(layar.medanWajib ?? [])

  async function jalankan<T>(f: () => Promise<T>, lanjut: (x: T) => void) {
    setSibuk(true)
    setGalat(null)
    setInfo('')
    try {
      lanjut(await f())
    } catch (e: unknown) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  const ubah = (jalur: string, v: string) => setH((x) => (x ? setel(x, jalur, v) : x))
  /** Action set satu sel - SATU bentuk permintaan: `urutan` satu refresh atau lebih. */
  const refresh = (urutan: Aksi[], indeks?: number, halaman?: Halaman) => {
    if (urutan.length === 0) return
    void jalankan(() => hitung(id, { urutan, indeks, halaman: halaman ?? h }), terima)
  }

  const selesai = (m: Medan, v: string) => {
    const baru = setel(h, m.jalur, v)
    setH(baru)
    if (m.aksi && boleh) refresh(m.aksi, undefined, baru)
  }

  /** Sel angka hanya-baca grid: nilai berformat, kode mata uang bila uang (AC 85). */
  const kodeMU = nilai(h, POLIS + 'Currency')
  const sel = (v: string | undefined, s: Sajian, uang = false) => (
    <>
      {uang && kodeMU && (v ?? '') !== '' && <span className="nbti__kode">{kodeMU}</span>}
      {sajikan(v ?? '', s)}
    </>
  )

  const kotak = (m: Medan, i: number) => (
    <KotakMedan
      key={`${m.jalur}-${i}`}
      medan={m}
      halaman={h}
      wajib={wajib.has(m.jalur)}
      hanyaBaca={!boleh}
      opsiMataUang={opsi(acuan?.mataUang)}
      opsiMO={opsi(acuan?.mo)}
      onUbah={ubah}
      onSelesai={selesai}
    />
  )

  const ubahBaris = (jalur: string, i: number, kunci: string, v: string) =>
    setH((x) => (x ? setelDaftar(x, jalur, daftar(x, jalur).map((b, j) => (j === i ? { ...b, [kunci]: v } : b))) : x))

  const kirim = () =>
    void jalankan(
      () => kirimKasus(id, h),
      (hasil) => onKembali(hasil.pesanKonversi ? `${PESAN.terkirim} ${hasil.pesanKonversi}` : PESAN.terkirim),
    )

  const tombolKirim = () => {
    if (!boleh) return null
    switch (layar.tombol) {
      case 'kirim':
        return (
          <button type="button" className="btn btn--primary" disabled={sibuk} onClick={kirim}>
            {TOMBOL.submit}
          </button>
        )
      case 'konfirmasi-tolak':
        return (
          <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => setKonfirmasi(true)}>
            {TOMBOL.submit}
          </button>
        )
      case 'nomor-polis':
        return (
          <button
            type="button"
            className="btn btn--primary"
            disabled={sibuk}
            onClick={() => void jalankan(() => terbitkanNomor(id, h), setNomor)}
          >
            {TOMBOL.submit}
          </button>
        )
    }
    return null
  }

  const spreading = daftar(h, SPREADING)
  const angsuran = daftar(h, ANGSURAN)
  const usulan = daftar(h, USULAN)
  const tanggalProduksi = tampilTanggalProduksi(h, layar.tempat)
  // Wadah bagian uang / spreading / angsuran (`pyContainerVisibleWhen`).
  const wadahUang = admin ? wadahUangAdmin(h) : wadahUangAtasan(h)
  const ubahAdmin = admin && boleh

  /** Sel persen spreading: tersunting admin, selain itu hanya-baca berformat. */
  const persenSpreading = (b: Baris, i: number, k: 'SharePercentage' | 'ClaimPercentage', label: string) =>
    ubahAdmin ? (
      // change -> refresh CountSpreading_Act(Index=.pxListSubscript)
      <InputAngka
        label={label}
        value={b[k] ?? ''}
        sajian={SAJIAN_SPREADING.persen}
        onChange={(v) => ubahBaris(SPREADING, i, k, v)}
        onBlur={() => refresh([{ aksi: 'CountSpreading' }], i + 1)}
      />
    ) : (
      sel(b[k], SAJIAN_SPREADING.persen)
    )

  return (
    <div className="inbox nbti__layar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">
          {JUDUL_POSISI[posisi] ?? JUDUL.portal} — {layar.kasus.id}
        </h2>
        <button type="button" className="btn btn--ghost" onClick={() => onKembali()}>
          {TOMBOL.kembali}
        </button>
      </header>
      {!boleh && <div className="alert alert--info">{PORTAL.hanyaBaca}</div>}
      {nilai(h, 'NBStatus') && <div className="nbti__status">{nilai(h, 'NBStatus')}</div>}
      {galat !== null && <Gagal galat={galat} />}
      {info && <div className="alert alert--ok">{info}</div>}
      {(layar.pesan ?? []).length > 0 && (
        <div className="alert alert--warn">
          {(layar.pesan ?? []).map((p) => (
            <div key={p}>{p}</div>
          ))}
        </div>
      )}

      <Panel judul={BAGIAN.umum}>
        {ubahAdmin && (
          <div className="nbti__aksi">
            {/* wadah `.ClaimType != 'XOL Retro'`; click -> showHarness BusinessAndSOBList -> refresh */}
            {bukanXOLRetro(h) && (
              <button type="button" className="btn" onClick={() => setPopupBisnis(true)}>
                {TOMBOL.chooseBusiness}
              </button>
            )}
            {/* `.ClaimType = 'XOL Retro'`; click -> showHarness SOB (paket P3) */}
            {tampilTombolSOB(h) && (
              <button type="button" className="btn" onClick={() => setPopupSOB(true)}>
                {TOMBOL.selectSOB}
              </button>
            )}
            {/* `.TreatyType='XOL'`; click -> refresh (pra-DT TreatyEnableDisableInput) */}
            {nilai(h, POLIS + 'TreatyType') === 'XOL' && (
              <button type="button" className="btn" onClick={() => refresh([{ aksi: 'TreatyEnableDisableInput' }])}>
                {TOMBOL.enableDisable}
              </button>
            )}
          </div>
        )}
        <div className="form-grid">{medanTampil(admin ? MEDAN_ADMIN_UMUM : MEDAN_ATASAN_UMUM, h).map(kotak)}</div>
      </Panel>

      {/* subsection `DetailPoliciesNonProportional` - wadah `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'` (K8) */}
      {tampilNonProp(h) && (
        <DetailNonProp
          halaman={h}
          sunting={admin && boleh}
          tempat={layar.tempat}
          // `SpreadingRiskList` .TreatyType: pyListSource reportdefinition BrowseReinsuranceType_RD (.ID / .Note)
          opsiSpreading={acuan?.jenisReas ?? []}
          onUbahBaris={ubahBaris}
          onSetelDaftar={(j, b) => setH(setelDaftar(h, j, b))}
          onRefresh={(aksi, indeks) => refresh([{ aksi }], indeks)}
        />
      )}

      {wadahUang && (
        <Panel judul={BAGIAN.uang}>
          <div className="form-grid">{medanTampil(admin ? MEDAN_ADMIN_UANG : MEDAN_ATASAN_UANG, h).map(kotak)}</div>
        </Panel>
      )}

      {wadahUang && (
        <Panel judul={BAGIAN.spreading}>
          {ubahAdmin && (
            // click -> addRow
            <button
              type="button"
              className="btn btn--sm"
              onClick={() => setH(setelDaftar(h, SPREADING, [...spreading, {} as Baris]))}
            >
              {TOMBOL.add}
            </button>
          )}
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th scope="col">{KOLOM_SPREADING.treatyType}</th>
                  <th scope="col">{KOLOM_SPREADING.share}</th>
                  <th scope="col">{KOLOM_SPREADING.premium}</th>
                  <th scope="col">{KOLOM_SPREADING.claimPct}</th>
                  <th scope="col">{KOLOM_SPREADING.claim}</th>
                  {ubahAdmin && <th scope="col" />}
                </tr>
              </thead>
              <tbody>
                {spreading.map((b, i) => (
                  <tr key={i}>
                    <td>
                      {ubahAdmin ? (
                        <select
                          className="field__input"
                          aria-label={KOLOM_SPREADING.treatyType}
                          value={b.TreatyType ?? ''}
                          onChange={(e) => ubahBaris(SPREADING, i, 'TreatyType', e.target.value)}
                        >
                          {/* pyNoSelectionText */}
                          <option value="">{TOMBOL.pilihKosong}</option>
                          {(acuan?.spreading ?? []).map((o) => (
                            <option key={o.nilai} value={o.nilai}>
                              {o.label}
                            </option>
                          ))}
                        </select>
                      ) : admin ? (
                        b.TreatyName || b.TreatyType || ''
                      ) : (
                        // DetailDeptHeadTreatyIn_UW .TreatyType: dropdown ro, BrowseReinsuranceType_RD (.ID -> .Note)
                        (acuan?.jenisReas ?? []).find((o) => o.nilai === b.TreatyType)?.label || b.TreatyType || ''
                      )}
                    </td>
                    <td>{persenSpreading(b, i, 'SharePercentage', KOLOM_SPREADING.share)}</td>
                    <td>{sel(b.PremiumSpreaded, SAJIAN_SPREADING.uang, true)}</td>
                    <td>{persenSpreading(b, i, 'ClaimPercentage', KOLOM_SPREADING.claimPct)}</td>
                    <td>{sel(b.ClaimSpreaded, SAJIAN_SPREADING.uang, true)}</td>
                    {ubahAdmin && (
                      <td>
                        {/* click -> deleteRow */}
                        <button
                          type="button"
                          className="btn btn--sm btn--danger"
                          onClick={() => setH(setelDaftar(h, SPREADING, spreading.filter((_, j) => j !== i)))}
                        >
                          {TOMBOL.delete}
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
              <tfoot>
                <tr>
                  <td>{KOLOM_SPREADING.totalShare}</td>
                  <td>{sel(nilai(h, POLIS + 'TotalSharePercentagePremium'), SAJIAN_SPREADING.total)}</td>
                  <td>{sel(nilai(h, POLIS + 'TotalPremium'), SAJIAN_SPREADING.total, true)}</td>
                  <td>{sel(nilai(h, POLIS + 'TotalSharePercentageClaim'), SAJIAN_SPREADING.total)}</td>
                  <td>{sel(nilai(h, POLIS + 'TotalClaim'), SAJIAN_SPREADING.total, true)}</td>
                </tr>
              </tfoot>
            </table>
          </div>
          {!admin && <div className="form-grid">{medanTampil(MEDAN_ATASAN_TOTAL, h).map(kotak)}</div>}
        </Panel>
      )}

      {wadahUang && (
        <Panel judul={BAGIAN.angsuran}>
          <div className="nbti__aksi">
            <label className="field__label">{KOLOM_ANGSURAN.installment}</label>
            {ubahAdmin ? (
              // change -> refresh FillPaymentInstallment(Installment=.Installment)
              <input
                className="field__input nbti__pendek"
                aria-label={KOLOM_ANGSURAN.installment}
                value={nilai(h, POLIS + 'Installment')}
                onChange={(e) => ubah(POLIS + 'Installment', e.target.value)}
                onBlur={() => refresh([{ aksi: 'FillPaymentInstallment' }])}
              />
            ) : (
              <span>{nilai(h, POLIS + 'Installment')}</span>
            )}
          </div>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th scope="col">{KOLOM_ANGSURAN.no}</th>
                  <th scope="col">{KOLOM_ANGSURAN.dueDate}</th>
                  <th scope="col">{KOLOM_ANGSURAN.pct}</th>
                  <th scope="col">{KOLOM_ANGSURAN.premium}</th>
                  <th scope="col">{KOLOM_ANGSURAN.total}</th>
                </tr>
              </thead>
              <tbody>
                {angsuran.map((b, i) => (
                  <tr key={i}>
                    <td>{b.InstallmentNo ?? ''}</td>
                    {/* Grid S45 `.ListInstallment` kedua layar: pyEditingMode/pyRowEditing
                        readOnly - sel ber-aksi SetValidateInstallment_Act / CountPctInstallment_Act
                        tidak pernah terpicu; baris diisi FillPaymentInstallment (audit silang P3). */}
                    <td>{sajikan(b.DueDate ?? '', SAJIAN_ANGSURAN.dueDate)}</td>
                    <td>{sel(b.InstallmentPercentage, SAJIAN_ANGSURAN.persen)}</td>
                    <td>{sel(b.Premium, SAJIAN_ANGSURAN.premium, true)}</td>
                    <td>{sel(b.PaymentTotal, SAJIAN_ANGSURAN.total, true)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Panel>
      )}

      <Panel judul={BAGIAN.usulan}>
        {boleh && (
          <div className="form-grid">
            <div className="field">
              <span className="field__label">
                {KOLOM_USULAN.putusan}
                <span className="field__req">*</span>
              </span>
              <div className="nbti__radio">
                {PILIHAN_APPROVAL.map((o) => (
                  <label key={o.value}>
                    <input
                      type="radio"
                      name="nbti-approval"
                      checked={nilai(h, POLIS + 'IsApproved') === o.value}
                      onChange={() => {
                        const baru = setel(h, POLIS + 'IsApproved', o.value)
                        setH(baru)
                        // change -> runActivity SetDueTo_act (CekLimitTreatyAcc_Act: K2;
                        // Protection_Act: varian tak ada di korpus, INVENTARIS 1.2)
                        refresh([{ aksi: 'SetDueTo' }], undefined, baru)
                      }}
                    />
                    {o.label}
                  </label>
                ))}
              </div>
            </div>
            {tanggalProduksi && (
              <KotakMedan
                medan={{ jalur: POLIS + 'ProductionDate', label: 'Production Date', jenis: 'tanggal' }}
                halaman={h}
                wajib={wajib.has(POLIS + 'ProductionDate')}
                hanyaBaca={false}
                opsiMataUang={[]}
                opsiMO={[]}
                onUbah={ubah}
                onSelesai={() => undefined}
              />
            )}
            <KotakMedan
              medan={{ jalur: POLIS + 'Suggest', label: 'Suggest', jenis: 'area' }}
              halaman={h}
              wajib
              hanyaBaca={false}
              opsiMataUang={[]}
              opsiMO={[]}
              onUbah={ubah}
              onSelesai={() => undefined}
            />
          </div>
        )}
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_USULAN.tanggal}</th>
                <th scope="col">{KOLOM_USULAN.operator}</th>
                <th scope="col">{KOLOM_USULAN.putusan}</th>
                <th scope="col">{KOLOM_USULAN.catatan}</th>
              </tr>
            </thead>
            <tbody>
              {usulan.map((b, i) => (
                <tr key={i}>
                  <td>{sajikan(b.Date ?? '', 'tanggal')}</td>
                  <td>{b.OperatorName ?? ''}</td>
                  <td>{PILIHAN_APPROVAL.find((o) => o.value === b.IsApproved)?.label ?? b.IsApproved ?? ''}</td>
                  <td className="nbti__catatan">{b.Suggest ?? ''}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>

      {riwayat.length > 0 && (
        <Panel judul={JUDUL.riwayat}>
          <div className="table-wrap">
            <table>
              <tbody>
                {riwayat.map((r, i) => (
                  <tr key={i}>
                    <td>{r.tglTransfer}</td>
                    <td>{r.workbasket}</td>
                    <td>{r.status}</td>
                    <td>{r.username}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Panel>
      )}

      {boleh && (
        <div className="nbti__kaki">
          {admin && (
            <button
              type="button"
              className="btn"
              disabled={sibuk}
              onClick={() =>
                void jalankan(
                  () => simpanKasus(id, h),
                  (ly) => {
                    terima(ly)
                    setInfo(PESAN.tersimpan)
                  },
                )
              }
            >
              {TOMBOL.save}
            </button>
          )}
          {tombolKirim()}
        </div>
      )}

      {popupBisnis && (
        <PilihBisnis
          onTutup={() => setPopupBisnis(false)}
          onPilih={(idDetail) => {
            setPopupBisnis(false)
            void jalankan(() => pilihBisnis(id, idDetail, h), terima)
          }}
        />
      )}
      {popupSOB && (
        <PilihSumberBisnis
          terpilih={nilai(h, 'Quotation.SourceOfBusiness')}
          sibuk={sibuk}
          onTutup={() => setPopupSOB(false)}
          onPilih={(idAgen, tutup) => {
            if (tutup) setPopupSOB(false)
            void jalankan(() => pilihSumberBisnis(id, idAgen, h), terima)
          }}
        />
      )}
      {konfirmasi && (
        <Modal
          judul={JUDUL.tolak}
          onTutup={() => setKonfirmasi(false)}
          labelBatal={TOMBOL.no}
          aksi={
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setKonfirmasi(false)
                kirim()
              }}
            >
              {TOMBOL.yes}
            </button>
          }
        >
          <p>{KONFIRMASI_TOLAK}</p>
        </Modal>
      )}
      {nomor && (
        <Modal
          judul={JUDUL.nomorPolis}
          onTutup={() => setNomor(null)}
          aksi={
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setNomor(null)
                kirim()
              }}
            >
              {TOMBOL.ok}
            </button>
          }
        >
          {/* Section ShowPolicyNoTreaty_SC: pyWorkPage.pyID, LABEL "telah diaksep menjadi", PolicyNo */}
          <p className="nbti__nomor">
            <strong>{nomor.id}</strong> {NOMOR_DIAKSEP} <strong>{nomor.policyNo}</strong>
          </p>
        </Modal>
      )}
    </div>
  )
}
