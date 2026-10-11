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

import { Fragment, useCallback, useRef, useState, type ReactNode } from 'react'

import { Gagal, Memuat, Modal, Panel, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import PanelLampiranReas from '../../../../inti/frontend/lampiran/PanelLampiranReas'
import { useAmbilBatal } from '../ambil'
import {
  POLIS,
  PREFIX_NBTREATYIN,
  ambilAcuan,
  bukaKasus,
  daftar,
  hitung,
  kirimKasus,
  nilai,
  pilihBisnis,
  pilihSumberBisnis,
  setel,
  setelDaftar,
  simpanKasus,
  terbitkanNomor,
  type Acuan,
  type Baris,
  type Halaman,
  type Layar,
  type NomorPolis,
} from '../api'
import DetailNonProp from '../components/DetailNonProp'
import InputAngka from '../components/InputAngka'
import KotakMedan from '../components/KotakMedan'
import Paginasi from '../components/Paginasi'
import PilihBisnis from '../components/PilihBisnis'
import PilihSumberBisnis, { pegangSumberBisnis, tampilTombolSOB } from '../components/PilihSumberBisnis'
import SurveiHistoris from '../components/SurveiHistoris'
import Wadah from '../components/Wadah'
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
  MEDAN_ADMIN_UMUM,
  MEDAN_ATASAN_UMUM,
  SAJIAN_ANGSURAN,
  SAJIAN_SPREADING,
  bukanXOLRetro,
  labelTreatyType,
  medanTampil,
  wadahUangAdmin,
  wadahUangAtasan,
  type Aksi,
  type Medan,
} from '../medan'
import { tampilNonProp } from '../nonprop'
import { pertahankanKetikan } from '../hitungLangsung'
import { BARIS_PER_HALAMAN_USULAN, irisan } from '../paginasi'
import { idTampil, sajikan, sajikanTanggalJam, type Sajian } from '../sajian'
import { TATA_UANG_ADMIN, TATA_UANG_ATASAN, TOTAL_ATASAN, deretLayer, deretQ, tataUmum } from '../tataletak'
import { aktifTombolSurvei, DAFTAR_SURVEI, tampilTombolSurvei } from '../survei'
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
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [info, setInfo] = useState('')
  const [popupBisnis, setPopupBisnis] = useState(false)
  const [popupSOB, setPopupSOB] = useState(false)
  const [popupSurvei, setPopupSurvei] = useState(false)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [nomor, setNomor] = useState<NomorPolis | null>(null)
  const [halUsulan, setHalUsulan] = useState(1)
  // nomor permintaan hitung terakhir - jawaban permintaan yang sudah disusul diabaikan (hitungLangsung.ts)
  const hitungTerakhir = useRef(0)

  const terima = useCallback((ly: Layar) => {
    setLayar(ly)
    setH(ly.halaman)
  }, [])

  useAmbilBatal(
    () => Promise.all([bukaKasus(id), ambilAcuan()]),
    ([ly, a]) => {
      terima(ly)
      setAcuan(a)
    },
    setGalat,
    [id, terima],
  )

  if (galat !== null && layar === null) return <Gagal galat={galat} />
  if (layar === null || h === null)
    return (
      <div className="inbox nbti__akar">
        <Memuat />
      </div>
    )

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
    const dikirim = halaman ?? h
    const ke = ++hitungTerakhir.current
    // Hitung langsung saat mengetik (hitungLangsung.ts): jawaban yang sudah disusul permintaan lain diabaikan, dan
    // medan yang diketik sesudah permintaan dikirim tidak ditimpa nilai lamanya.
    void jalankan(
      () => hitung(id, { urutan, indeks, halaman: dikirim }),
      (ly) => {
        if (ke !== hitungTerakhir.current) return
        setLayar(ly)
        setH((kini) => (kini ? pertahankanKetikan(ly.halaman, dikirim, kini) : ly.halaman))
      },
    )
  }

  const selesai = (m: Medan, v: string) => {
    const baru = setel(h, m.jalur, v)
    setH(baru)
    if (m.aksi && boleh) refresh(m.aksi, undefined, baru)
  }

  /** Sel angka hanya-baca grid: nilai berformat. Kode mata uang tidak diulang di setiap sel - sudah tampil di
   *  medan Currency (perintah work owner 05-10-2026: "banyak sekali tulisan IDR"). */
  const sel = (v: string | undefined, s: Sajian) => sajikan(v ?? '', s)

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

  /** Satu kolom label-kiri (tata letak Pega): deret "Q / U/Y" satu baris; `sesudah[jalur]` dirender tepat
   *  di bawah medannya (tombol Survey Report). */
  const kolom = (ms: Medan[], sesudah: Record<string, ReactNode> = {}) => (
    <div className="nbti__kolom">
      {deretQ(ms).map((x, i) =>
        Array.isArray(x) ? (
          <div key={`deret-${i}`} className={deretLayer(x) ? 'nbti__deret nbti__deret--layer' : 'nbti__deret'}>
            {x.map(kotak)}
          </div>
        ) : (
          <Fragment key={`${x.jalur}-${i}`}>
            {kotak(x, i)}
            {sesudah[x.jalur]}
          </Fragment>
        ),
      )}
    </div>
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

  // sel 22 `DetailPolicyTreatyIn` / `DetailDeptHeadTreatyIn_UW`: showHarness popup Historical Survey Report
  const tombolSurvei = tampilTombolSurvei(h) ? (
    <div className="nbti__tombol-survei">
      <button type="button" className="btn btn--sm" disabled={!aktifTombolSurvei(h)} onClick={() => setPopupSurvei(true)}>
        {TOMBOL.surveyReport}
      </button>
    </div>
  ) : null

  const aksiUmum = ubahAdmin ? (
    <div className="nbti__aksi">
      {/* wadah `.ClaimType != 'XOL Retro'`; click -> showHarness BusinessAndSOBList -> refresh */}
      {bukanXOLRetro(h) && (
        <button type="button" className="btn" onClick={() => setPopupBisnis(true)}>
          {TOMBOL.chooseBusiness}
        </button>
      )}
      {/* `.ClaimType = 'XOL Retro'`; click -> showHarness SOB (paket P3; F4: pilihan dipegang layar) */}
      {tampilTombolSOB(h) && (
        <button type="button" className="btn" onClick={() => setPopupSOB(true)}>
          {TOMBOL.selectSOB}
        </button>
      )}
      {/* `.TreatyType='XOL'`: tombol TAMPIL tetapi NONAKTIF - klik tidak menjalankan pra-DT TreatyEnableDisableInput
          (perintah work owner 06-10-2026: "jangan di hide, tapi di disable aja; muncul untuk XOL, tapi di disable") */}
      {nilai(h, POLIS + 'TreatyType') === 'XOL' && (
        <button type="button" className="btn" disabled>
          {TOMBOL.enableDisable}
        </button>
      )}
    </div>
  ) : null

  return (
    <div className="inbox nbti__layar nbti__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">
          {JUDUL_POSISI[posisi] ?? JUDUL.portal} — {idTampil(layar.kasus.id)}
        </h2>
        <button type="button" className="btn btn--ghost" onClick={() => onKembali()}>
          {TOMBOL.kembali}
        </button>
      </header>
      {!boleh && <div className="alert alert--info">{PORTAL.hanyaBaca}</div>}
      {galat !== null && <Gagal galat={galat} />}
      {info && <div className="alert alert--ok">{info}</div>}
      {(layar.pesan ?? []).length > 0 && (
        <div className="alert alert--warn">
          {(layar.pesan ?? []).map((p) => (
            <div key={p}>{p}</div>
          ))}
        </div>
      )}

      {/* S2 (pyTitle "General", NOHEADER: tanpa judul) */}
      <Wadah>
        {/* tombol aksi paling atas; lalu kolom kiri / kanan, Remark melebar (tata letak Pega) */}
        {aksiUmum}
        {(() => {
          const t = tataUmum(medanTampil(admin ? MEDAN_ADMIN_UMUM : MEDAN_ATASAN_UMUM, h))
          return (
            <>
              <div className="nbti__dua-kolom">
                {kolom(t.kiri, { [POLIS + 'QuotationData.IsSurveyReport']: tombolSurvei })}
                {kolom(t.kanan)}
              </div>
              {t.bawah.length > 0 && <div className="nbti__baris-penuh">{kolom(t.bawah)}</div>}
            </>
          )
        })()}
      </Wadah>

      {/* subsection `DetailPoliciesNonProportional` - wadah `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'` (K8) */}
      {tampilNonProp(h) && (
        <DetailNonProp
          halaman={h}
          // %Share spreading hanya tersunting di layar ADMIN (keputusan work owner 07-10-2026: "HANYA ADMIN YANG
          // BISA EDIT"); XML S88 atasan ber-`pyEditOptions=Auto`, disimpangkan sadar
          sunting={ubahAdmin}
          tempat={layar.tempat}
          // `SpreadingRiskList` .TreatyType: pyListSource reportdefinition BrowseReinsuranceType_RD (.ID / .Note)
          opsiSpreading={acuan?.jenisReas ?? []}
          onUbahBaris={ubahBaris}
          onRefresh={(aksi, indeks) => refresh([{ aksi }], indeks)}
        />
      )}

      {wadahUang && (
        // wadah S19 / S90 NOHEADER; LABEL Heading 4 "OGP" / "ONP" di atas wadah selnya
        <Wadah>
          {/* tata letak Pega: Gross Premium 100% | Claim 100%, lalu kolom OGP dan kolom ONP */}
          {(() => {
            const t = admin ? TATA_UANG_ADMIN : TATA_UANG_ATASAN
            const atasKanan = medanTampil(t.atasKanan, h)
            return (
              <>
                <div className="nbti__dua-kolom">
                  {kolom(medanTampil(t.atasKiri, h))}
                  {atasKanan.length > 0 && kolom(atasKanan)}
                </div>
                <div className="nbti__dua-kolom">
                  {[t.kiri, t.kanan].map((ks, i) => (
                    <div key={i} className="nbti__kolom-grup">
                      {ks.map((k, j) => (
                        <Fragment key={j}>
                          {k.judul && <h4>{k.judul}</h4>}
                          {kolom(medanTampil(k.medan, h))}
                        </Fragment>
                      ))}
                    </div>
                  ))}
                </div>
              </>
            )
          })()}
        </Wadah>
      )}

      {wadahUang && (
        // wadah S29/S30 dan S95/S96 NOHEADER - tanpa judul
        <Wadah>
          {/* perintah work owner 06-10-2026: Add / Delete dibuang, Treaty Type hanya-baca ("ga boleh di ubah lagi") */}
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th scope="col">{KOLOM_SPREADING.treatyType}</th>
                  <th scope="col" className="nbti__angka">{KOLOM_SPREADING.share}</th>
                  <th scope="col" className="nbti__angka">{KOLOM_SPREADING.premium}</th>
                  <th scope="col" className="nbti__angka">{KOLOM_SPREADING.claimPct}</th>
                  <th scope="col" className="nbti__angka">{KOLOM_SPREADING.claim}</th>
                </tr>
              </thead>
              <tbody>
                {spreading.map((b, i) => (
                  <tr key={i}>
                    <td>
                      {/* hanya-baca: label .ID ListSpreading / BrowseReinsuranceType_RD, selain itu .TreatyName */}
                      {labelTreatyType(b, (admin ? acuan?.spreading : acuan?.jenisReas) ?? [])}
                    </td>
                    <td className="nbti__angka">{persenSpreading(b, i, 'SharePercentage', KOLOM_SPREADING.share)}</td>
                    <td className="nbti__angka">{sel(b.PremiumSpreaded, SAJIAN_SPREADING.uang)}</td>
                    <td className="nbti__angka">{persenSpreading(b, i, 'ClaimPercentage', KOLOM_SPREADING.claimPct)}</td>
                    <td className="nbti__angka">{sel(b.ClaimSpreaded, SAJIAN_SPREADING.uang)}</td>
                  </tr>
                ))}
              </tbody>
              <tfoot>
                <tr>
                  <td>{KOLOM_SPREADING.totalShare}</td>
                  <td className="nbti__angka">{sel(nilai(h, POLIS + 'TotalSharePercentagePremium'), SAJIAN_SPREADING.total)}</td>
                  <td className="nbti__angka">{sel(nilai(h, POLIS + 'TotalPremium'), SAJIAN_SPREADING.total)}</td>
                  <td className="nbti__angka">{sel(nilai(h, POLIS + 'TotalSharePercentageClaim'), SAJIAN_SPREADING.total)}</td>
                  <td className="nbti__angka">{sel(nilai(h, POLIS + 'TotalClaim'), SAJIAN_SPREADING.total)}</td>
                </tr>
              </tfoot>
            </table>
          </div>
          {!admin && (
            <div className="nbti__dua-kolom nbti__total">
              {TOTAL_ATASAN.map((ms, i) => (
                <Fragment key={i}>{kolom(medanTampil(ms, h))}</Fragment>
              ))}
            </div>
          )}
        </Wadah>
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
                  <th scope="col" className="nbti__angka">{KOLOM_ANGSURAN.pct}</th>
                  <th scope="col" className="nbti__angka">{KOLOM_ANGSURAN.premium}</th>
                  <th scope="col" className="nbti__angka">{KOLOM_ANGSURAN.total}</th>
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
                    <td className="nbti__angka">{sel(b.InstallmentPercentage, SAJIAN_ANGSURAN.persen)}</td>
                    <td className="nbti__angka">{sel(b.Premium, SAJIAN_ANGSURAN.premium)}</td>
                    <td className="nbti__angka">{sel(b.PaymentTotal, SAJIAN_ANGSURAN.total)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Panel>
      )}

      {/* SUB_SECTION ListSuggest - wadah NOHEADER, tanpa judul */}
      <Wadah>
        {boleh && (
          <div className="nbti__kolom nbti__usulan">
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
              {irisan(usulan, halUsulan, BARIS_PER_HALAMAN_USULAN).map((b, i) => (
                <tr key={i}>
                  {/* tanggal + jam: panel History dibuang (perintah work owner 06-10-2026 "PAKE YG ATAS AJA, TAPI TAMBAHIN JAM NYA") */}
                  <td>{sajikanTanggalJam(b.Date)}</td>
                  <td>{b.OperatorName ?? ''}</td>
                  <td>{PILIHAN_APPROVAL.find((o) => o.value === b.IsApproved)?.label ?? b.IsApproved ?? ''}</td>
                  <td className="nbti__catatan">{b.Suggest ?? ''}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {/* pyGridPaginator - pyPageMode Numeric, pyPageSizeOther 5 */}
        <Paginasi jumlahBaris={usulan.length} ukuran={BARIS_PER_HALAMAN_USULAN} hal={halUsulan} onHal={setHalUsulan} />
      </Wadah>

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

      {/* Lampiran "Reas" (`AttachmentGridReas`) di bawah layar kasus - keputusan work owner 08-10-2026; Upload / Delete
          selama kasus belum Resolve (diputuskan backend). */}
      <PanelLampiranReas dasar={`${PREFIX_NBTREATYIN}/kasus/${encodeURIComponent(id)}/lampiran`} />

      {popupSurvei && (
        <SurveiHistoris
          halaman={h}
          admin={admin}
          sunting={ubahAdmin}
          onSetel={(b) => setH(setelDaftar(h, DAFTAR_SURVEI, b))}
          onTutup={() => setPopupSurvei(false)}
        />
      )}
      {popupBisnis && (
        <PilihBisnis
          id={id}
          halaman={h}
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
            // Klik baris = PostDT TANPA simpan (F4): hasilnya dipegang state layar dan
            // ikut Save/Submit/refresh. XML tidak menjalankan refresh sesudahnya.
            void jalankan(
              () => pilihSumberBisnis(id, idAgen, h),
              (hasil) => setH((x) => (x ? pegangSumberBisnis(x, hasil) : x)),
            )
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
        // Tanpa Cancel / X / Escape / klik luar: nomor polis sudah terbit, jadi hanya OK = kirim (perintah work owner
        // 06-10-2026: "tidak mau ada cancel, tidak ada close ... selalu maju kalo udah klik aksep").
        <Modal
          judul={JUDUL.nomorPolis}
          onTutup={() => setNomor(null)}
          tanpaTutup
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
          {/* Section ShowPolicyNoTreaty_SC: pyWorkPage.pyID, LABEL "telah diaksep menjadi", PolicyNo - bertumpuk di
              tengah, nomor polis menonjol (perintah work owner 06-10-2026 "RAPIHIN") */}
          <div className="nbti__nomor">
            <p className="nbti__nomor-kasus">{idTampil(nomor.id)}</p>
            <p className="nbti__nomor-teks">{NOMOR_DIAKSEP}</p>
            <p className="nbti__nomor-polis">{nomor.policyNo}</p>
          </div>
        </Modal>
      )}
    </div>
  )
}
