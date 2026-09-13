# Roadmap

## Status saat ini

Compiler pseudocode-to-Go (lexer, parser, AST, codegen) sudah stabil dan
tervalidasi end-to-end lewat `go build` sungguhan, termasuk dukungan
record/array bertingkat, dan dua gaya penulisan block-closer (`end if`
maupun `endif`) hasil temuan dari sampel tulisan tangan mahasiswa asli.

## OCR untuk pseudocode tulisan tangan — dieksplorasi, belum diimplementasi

Fitur ini pernah dieksplorasi cukup dalam sebelum diputuskan untuk
ditunda sebagai pengembangan lanjutan, bukan syarat rilis. Catatan di
bawah ini dokumentasi temuan nyata dari eksperimen tersebut, supaya
siapa pun yang melanjutkan tidak mengulang dari nol.

### Yang sudah dicoba

**`microsoft/trocr-base-handwritten`** — gagal total, bukan sekadar
kurang akurat. Model ini adalah *line recognizer* (didesain menerima
satu baris teks yang sudah dipotong rapi oleh model deteksi layout
terpisah), bukan *full-page reader*. Diberi gambar satu halaman penuh
(kolom nama + badan kode berindentasi), modelnya berhalusinasi ke pola
teks Wikipedia dari data latihnya, sama sekali tidak berhubungan dengan
isi gambar.

**PaddleOCR** (`paddleocr.PaddleOCR().predict()`) — jauh lebih baik,
genuinely usable. Kesalahannya lokal per-karakter, bukan halusinasi:

- Substitusi karakter yang cukup dapat ditebak polanya: `l`/`I`,
  `m`/`n`, `0`/`o`, digit `9` terbaca `g`.
- Kata yang sama bisa terbaca beda-beda di baris berbeda (`jumlah`
  pernah muncul sebagai `junlah`, `jinlah`, `jwmlah` di satu dokumen
  yang sama) — ini alasan kenapa cross-reference ke identifier yang
  sudah dideklarasikan di `Kamus` (bukan spell-check generik) adalah
  pendekatan koreksi yang tepat.
- Simbol panah tulisan tangan (`<-`, `->`) terbaca sebagai Unicode asli
  (`←`, `→`) — perlu normalisasi sebelum tokenizing.
- Kegagalan *segmentasi baris* (bukan cuma recognition) di bagian
  tulisan paling padat/bersarang — beberapa baris jadi word-salad yang
  jujur tidak terpulihkan tanpa menebak, dan sebaiknya tidak dipaksa.
- Karakter asing yang tidak seharusnya muncul (CJK, diakritik) walau
  `lang="en"` — anomali yang belum diinvestigasi penyebabnya.
- **`rec_scores` per kata berkorelasi kuat dengan baris bermasalah** —
  baris-baris yang confidence-nya paling rendah (< 0.65) hampir selalu
  bertepatan dengan baris yang secara kualitatif memang tidak
  terpulihkan. Tapi bukan sinyal sempurna (ada satu kasus baris rusak
  dengan confidence 0.78) — confidence sebaiknya jadi sinyal prioritas,
  bukan gerbang biner satu-satunya.

### Perubahan grammar compiler yang lahir dari temuan ini

Dua sampel tulisan tangan mahasiswa nyata (independen satu sama lain)
konsisten memakai pola yang sebelumnya tidak didukung compiler:

- Penutup blok satu-kata (`endif`, `endwhile`, `endfor`, `endfunction`,
  `endprocedure`, `endprogram`) — sekarang diterima sebagai alias dari
  bentuk dua-kata.
- `==` untuk kesetaraan (gaya C), diterima setara dengan `=`.
- `...` (lebih dari dua titik) untuk rentang array, ditoleransi.
- `if` bersarang di dalam `else` dengan penutup `endif` sendiri —
  berbeda dari gaya `else if ... end if` yang satu penutup bersama.
  Parser sekarang mendukung keduanya tanpa ambigu.

### Rancangan pipeline yang direkomendasikan untuk implementasi lanjutan

1. **Pemangkasan pembuka non-kode** — buang baris nama/NIM/judul soal
   sebelum token top-level pertama (`type`/`constant`/`Program`/
   `procedure`/`function`); parser tidak punya cara merepresentasikan
   teks bebas ini.
2. **Normalisasi deterministik** — panah Unicode ke ASCII, dan
   substitusi karakter berisiko rendah lainnya.
3. **Snapping ke vocabulary yang diketahui** — cocokkan identifier/
   keyword asing ke kandidat terdekat (edit-distance) dari: (a) daftar
   keyword grammar, (b) identifier yang sudah dideklarasikan di scope
   yang sama.
4. **Triase berbasis confidence + hasil parse** — baris confidence
   rendah ATAU yang tetap gagal parse setelah langkah 1–3, ditandai
   untuk **tinjauan manual**, bukan ditebak paksa. Ini prinsip
   pegangan: corrector boleh membetulkan typo, tidak boleh mengarang
   logika yang tidak benar-benar tertulis mahasiswa.

### Kandidat mesin transkripsi

Vision-LLM (mis. Claude, GPT) saat ini mengungguli OCR/HTR khusus untuk
akurasi tulisan tangan mentah (per benchmark IAM terkini), dan lebih
murah dari OCR tulisan tangan tingkat spesialis. Risiko utamanya:
model bisa "melengkapi" tulisan ambigu jadi versi yang gramatikal-benar
padahal mahasiswa menulis logika yang salah — mitigasinya lewat prompt
eksplisit ("transkripsikan persis, jangan perbaiki logika") dan tetap
menjaga langkah 3–4 di atas sebagai lapisan terpisah yang auditable.

### Task breakdown (kandidat GitHub Issues)

- [ ] Implementasi tahap pemangkasan header non-kode
- [ ] Implementasi normalisasi simbol (panah, karakter berisiko rendah)
- [ ] Implementasi snapping identifier/keyword berbasis edit-distance
- [ ] Implementasi triase confidence + integrasi dengan parser error
- [ ] Investigasi anomali karakter CJK/diakritik di PaddleOCR
- [ ] Evaluasi vision-LLM sebagai mesin transkripsi alternatif/fallback
- [ ] UI untuk menandai & mengedit baris yang diflag "tinjauan manual"